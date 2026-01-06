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

// GetAll mengembalikan semua suppressions
func (r *SuppressionRepository) GetAll() ([]*models.Suppression, error) {
	query := `
		SELECT email, reason, source, created_at, updated_at 
		FROM suppressions 
		ORDER BY updated_at DESC
	`

	rows, err := r.db.Query(query)
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
