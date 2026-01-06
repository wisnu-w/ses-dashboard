package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"ses-monitoring/internal/domain/models"
)

type SuppressionRepository struct {
	db *sql.DB
}

func NewSuppressionRepository(db *sql.DB) *SuppressionRepository {
	return &SuppressionRepository{db: db}
}

// BulkUpsert melakukan bulk insert/update suppressions
func (r *SuppressionRepository) BulkUpsert(suppressions []*models.Suppression) error {
	if len(suppressions) == 0 {
		return nil
	}

	log.Printf("Starting bulk upsert for %d suppressions", len(suppressions))

	// Process in batches to avoid parameter limit
	batchSize := 100
	totalProcessed := 0

	for i := 0; i < len(suppressions); i += batchSize {
		end := i + batchSize
		if end > len(suppressions) {
			end = len(suppressions)
		}

		batch := suppressions[i:end]
		log.Printf("Processing batch %d-%d (%d items)", i+1, end, len(batch))

		if err := r.processBatch(batch); err != nil {
			return fmt.Errorf("failed to process batch %d-%d: %w", i+1, end, err)
		}

		totalProcessed += len(batch)
	}

	log.Printf("Bulk upsert completed: %d total suppressions processed", totalProcessed)
	return nil
}

func (r *SuppressionRepository) processBatch(suppressions []*models.Suppression) error {
	valueStrings := make([]string, 0, len(suppressions))
	valueArgs := make([]interface{}, 0, len(suppressions)*5)

	for i, s := range suppressions {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)",
			i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		valueArgs = append(valueArgs, s.Email, s.Reason, s.Source, s.CreatedAt, s.UpdatedAt)
	}

	query := fmt.Sprintf(`
		INSERT INTO suppressions (email, reason, source, created_at, updated_at) 
		VALUES %s
		ON CONFLICT (email) 
		DO UPDATE SET 
			reason = EXCLUDED.reason,
			source = EXCLUDED.source,
			updated_at = EXCLUDED.updated_at
	`, strings.Join(valueStrings, ","))

	result, err := r.db.Exec(query, valueArgs...)
	if err != nil {
		log.Printf("Batch upsert failed: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Batch upsert completed: %d rows affected", rowsAffected)

	return nil
}

// GetAll mengembalikan suppressions dengan pagination
func (r *SuppressionRepository) GetAll(limit, offset int) ([]*models.Suppression, error) {
	query := `
		SELECT email, reason, source, created_at, updated_at 
		FROM suppressions 
		ORDER BY updated_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppressions []*models.Suppression
	for rows.Next() {
		s := &models.Suppression{}
		err := rows.Scan(&s.Email, &s.Reason, &s.Source, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppressions = append(suppressions, s)
	}

	return suppressions, nil
}

// GetAllCount mengembalikan total count suppressions
func (r *SuppressionRepository) GetAllCount() (int, error) {
	query := `SELECT COUNT(*) FROM suppressions`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// SearchSuppressions mencari suppressions dengan pagination
func (r *SuppressionRepository) SearchSuppressions(searchTerm string, limit, offset int) ([]*models.Suppression, error) {
	query := `
		SELECT email, reason, source, created_at, updated_at 
		FROM suppressions 
		WHERE email ILIKE $1 OR reason ILIKE $1 OR source ILIKE $1
		ORDER BY updated_at DESC
		LIMIT $2 OFFSET $3
	`

	search := "%" + searchTerm + "%"
	rows, err := r.db.Query(query, search, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppressions []*models.Suppression
	for rows.Next() {
		s := &models.Suppression{}
		err := rows.Scan(&s.Email, &s.Reason, &s.Source, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppressions = append(suppressions, s)
	}

	return suppressions, nil
}

// GetSearchCount mengembalikan count hasil search
func (r *SuppressionRepository) GetSearchCount(searchTerm string) (int, error) {
	query := `
		SELECT COUNT(*) FROM suppressions 
		WHERE email ILIKE $1 OR reason ILIKE $1 OR source ILIKE $1
	`
	search := "%" + searchTerm + "%"
	var count int
	err := r.db.QueryRow(query, search).Scan(&count)
	return count, err
}

// GetByEmail mencari suppression berdasarkan email
func (r *SuppressionRepository) GetByEmail(email string) (*models.Suppression, error) {
	query := `
		SELECT email, reason, source, created_at, updated_at 
		FROM suppressions 
		WHERE email = $1
	`

	s := &models.Suppression{}
	err := r.db.QueryRow(query, email).Scan(&s.Email, &s.Reason, &s.Source, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

// Delete menghapus suppression berdasarkan email
func (r *SuppressionRepository) Delete(email string) error {
	query := `DELETE FROM suppressions WHERE email = $1`
	_, err := r.db.Exec(query, email)
	return err
}

// GetBySource mengembalikan suppressions berdasarkan source
func (r *SuppressionRepository) GetBySource(source string) ([]*models.Suppression, error) {
	query := `
		SELECT email, reason, source, created_at, updated_at 
		FROM suppressions 
		WHERE source = $1
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(query, source)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suppressions []*models.Suppression
	for rows.Next() {
		s := &models.Suppression{}
		err := rows.Scan(&s.Email, &s.Reason, &s.Source, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppressions = append(suppressions, s)
	}

	return suppressions, nil
}

// BulkDelete menghapus multiple suppressions berdasarkan email list
func (r *SuppressionRepository) BulkDelete(emails []string) error {
	if len(emails) == 0 {
		return nil
	}

	log.Printf("Starting bulk delete for %d emails", len(emails))

	// Process in batches
	batchSize := 100
	totalDeleted := 0

	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}

		batch := emails[i:end]
		log.Printf("Deleting batch %d-%d (%d emails)", i+1, end, len(batch))

		// Create placeholders for IN clause
		placeholders := make([]string, len(batch))
		args := make([]interface{}, len(batch))
		for j, email := range batch {
			placeholders[j] = fmt.Sprintf("$%d", j+1)
			args[j] = email
		}

		query := fmt.Sprintf("DELETE FROM suppressions WHERE email IN (%s)", strings.Join(placeholders, ","))
		result, err := r.db.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to delete batch %d-%d: %w", i+1, end, err)
		}

		rowsAffected, _ := result.RowsAffected()
		totalDeleted += int(rowsAffected)
		log.Printf("Deleted %d rows in batch", rowsAffected)
	}

	log.Printf("Bulk delete completed: %d total emails deleted", totalDeleted)
	return nil
}

// CountBySource menghitung jumlah suppressions berdasarkan source
func (r *SuppressionRepository) CountBySource(source string) (int, error) {
	query := `SELECT COUNT(*) FROM suppressions WHERE source = $1`
	var count int
	err := r.db.QueryRow(query, source).Scan(&count)
	return count, err
}
