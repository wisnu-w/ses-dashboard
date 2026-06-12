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
	valueArgs := make([]interface{}, 0, len(suppressions)*8)

	for i, s := range suppressions {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8))
		valueArgs = append(valueArgs, s.Email, s.Reason, s.Source, s.SuppressionType, s.AWSStatus, true, s.CreatedAt, s.UpdatedAt)
	}

	query := fmt.Sprintf(`
		INSERT INTO suppressions (email, reason, source, suppression_type, aws_status, is_active, created_at, updated_at) 
		VALUES %s
		ON CONFLICT (email) 
		DO UPDATE SET 
			reason = EXCLUDED.reason,
			source = EXCLUDED.source,
			suppression_type = EXCLUDED.suppression_type,
			aws_status = EXCLUDED.aws_status,
			is_active = EXCLUDED.is_active,
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

// GetAll mengembalikan suppressions dengan pagination (hanya yang aktif)
func (r *SuppressionRepository) GetAll(limit, offset int) ([]*models.Suppression, error) {
	query := `
		SELECT id, email, reason, source, suppression_type, aws_status, is_active, added_by, created_at, updated_at 
		FROM suppressions 
		WHERE is_active = true
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
		err := rows.Scan(&s.ID, &s.Email, &s.Reason, &s.Source, &s.SuppressionType, &s.AWSStatus, &s.IsActive, &s.AddedBy, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppressions = append(suppressions, s)
	}

	return suppressions, nil
}

// GetAllCount mengembalikan total count suppressions yang aktif
func (r *SuppressionRepository) GetAllCount() (int, error) {
	query := `SELECT COUNT(*) FROM suppressions WHERE is_active = true`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

// SearchSuppressions mencari suppressions dengan pagination (hanya yang aktif)
func (r *SuppressionRepository) SearchSuppressions(searchTerm string, limit, offset int) ([]*models.Suppression, error) {
	query := `
		SELECT id, email, reason, source, suppression_type, aws_status, is_active, added_by, created_at, updated_at 
		FROM suppressions 
		WHERE is_active = true 
		AND (email ILIKE $1 OR reason ILIKE $1 OR source ILIKE $1)
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
		err := rows.Scan(&s.ID, &s.Email, &s.Reason, &s.Source, &s.SuppressionType, &s.AWSStatus, &s.IsActive, &s.AddedBy, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppressions = append(suppressions, s)
	}

	return suppressions, nil
}

// GetSearchCount mengembalikan count hasil search (hanya yang aktif)
func (r *SuppressionRepository) GetSearchCount(searchTerm string) (int, error) {
	query := `
		SELECT COUNT(*) FROM suppressions 
		WHERE is_active = true 
		AND (email ILIKE $1 OR reason ILIKE $1 OR source ILIKE $1)
	`
	search := "%" + searchTerm + "%"
	var count int
	err := r.db.QueryRow(query, search).Scan(&count)
	return count, err
}

// GetByEmail mencari suppression berdasarkan email
func (r *SuppressionRepository) GetByEmail(email string) (*models.Suppression, error) {
	query := `
		SELECT id, email, reason, source, suppression_type, aws_status, is_active, added_by, created_at, updated_at 
		FROM suppressions 
		WHERE email = $1
	`

	s := &models.Suppression{}
	err := r.db.QueryRow(query, email).Scan(&s.ID, &s.Email, &s.Reason, &s.Source, &s.SuppressionType, &s.AWSStatus, &s.IsActive, &s.AddedBy, &s.CreatedAt, &s.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return s, nil
}

// Delete menghapus suppression berdasarkan email (soft delete - set is_active = false)
func (r *SuppressionRepository) Delete(email string) error {
	query := `UPDATE suppressions SET is_active = false, updated_at = NOW() WHERE email = $1`
	_, err := r.db.Exec(query, email)
	return err
}

// HardDelete menghapus suppression secara permanen berdasarkan email
func (r *SuppressionRepository) HardDelete(email string) error {
	query := `DELETE FROM suppressions WHERE email = $1`
	_, err := r.db.Exec(query, email)
	return err
}

// GetBySource mengembalikan suppressions berdasarkan source
func (r *SuppressionRepository) GetBySource(source string) ([]*models.Suppression, error) {
	query := `
		SELECT id, email, reason, source, suppression_type, aws_status, is_active, added_by, created_at, updated_at 
		FROM suppressions 
		WHERE source = $1 AND is_active = true
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
		err := rows.Scan(&s.ID, &s.Email, &s.Reason, &s.Source, &s.SuppressionType, &s.AWSStatus, &s.IsActive, &s.AddedBy, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		suppressions = append(suppressions, s)
	}

	return suppressions, nil
}

// BulkDelete menghapus multiple suppressions berdasarkan email list (soft delete)
func (r *SuppressionRepository) BulkDelete(emails []string) error {
	if len(emails) == 0 {
		return nil
	}

	log.Printf("Starting bulk soft delete for %d emails", len(emails))

	// Process in batches
	batchSize := 100
	totalDeleted := 0

	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}

		batch := emails[i:end]
		log.Printf("Soft deleting batch %d-%d (%d emails)", i+1, end, len(batch))

		// Create placeholders for IN clause
		placeholders := make([]string, len(batch))
		args := make([]interface{}, len(batch))
		for j, email := range batch {
			placeholders[j] = fmt.Sprintf("$%d", j+1)
			args[j] = email
		}

		query := fmt.Sprintf("UPDATE suppressions SET is_active = false, updated_at = NOW() WHERE email IN (%s)", strings.Join(placeholders, ","))
		result, err := r.db.Exec(query, args...)
		if err != nil {
			return fmt.Errorf("failed to soft delete batch %d-%d: %w", i+1, end, err)
		}

		rowsAffected, _ := result.RowsAffected()
		totalDeleted += int(rowsAffected)
		log.Printf("Soft deleted %d rows in batch", rowsAffected)
	}

	log.Printf("Bulk soft delete completed: %d total emails soft deleted", totalDeleted)
	return nil
}

// BulkHardDelete menghapus multiple suppressions secara permanen berdasarkan email list
func (r *SuppressionRepository) BulkHardDelete(emails []string) error {
	if len(emails) == 0 {
		return nil
	}

	log.Printf("Starting bulk hard delete for %d emails", len(emails))

	// Process in batches
	batchSize := 100
	totalDeleted := 0

	for i := 0; i < len(emails); i += batchSize {
		end := i + batchSize
		if end > len(emails) {
			end = len(emails)
		}

		batch := emails[i:end]
		log.Printf("Hard deleting batch %d-%d (%d emails)", i+1, end, len(batch))

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
			return fmt.Errorf("failed to hard delete batch %d-%d: %w", i+1, end, err)
		}

		rowsAffected, _ := result.RowsAffected()
		totalDeleted += int(rowsAffected)
		log.Printf("Hard deleted %d rows in batch", rowsAffected)
	}

	log.Printf("Bulk hard delete completed: %d total emails hard deleted", totalDeleted)
	return nil
}

// CountBySource menghitung jumlah suppressions berdasarkan source
func (r *SuppressionRepository) CountBySource(source string) (int, error) {
	query := `SELECT COUNT(*) FROM suppressions WHERE source = $1 AND is_active = true`
	var count int
	err := r.db.QueryRow(query, source).Scan(&count)
	return count, err
}

// Add menambahkan suppression baru
func (r *SuppressionRepository) Add(suppression *models.Suppression) error {
	query := `
		INSERT INTO suppressions (email, reason, source, suppression_type, aws_status, is_active, added_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		ON CONFLICT (email) 
		DO UPDATE SET 
			reason = $2,
			source = $3,
			suppression_type = $4,
			aws_status = $5,
			is_active = $6,
			added_by = $7,
			updated_at = NOW()
	`
	_, err := r.db.Exec(query, suppression.Email, suppression.Reason, suppression.Source, suppression.SuppressionType, suppression.AWSStatus, suppression.IsActive, suppression.AddedBy)
	return err
}

// UpdateAWSStatus mengupdate AWS status suppression
func (r *SuppressionRepository) UpdateAWSStatus(email string, status string) error {
	query := `UPDATE suppressions SET aws_status = $1, updated_at = NOW() WHERE email = $2`
	_, err := r.db.Exec(query, status, email)
	return err
}

// MarkAsSynced menandai suppression sebagai synced
func (r *SuppressionRepository) MarkAsSynced(email string) error {
	query := `UPDATE suppressions SET updated_at = NOW() WHERE email = $1`
	_, err := r.db.Exec(query, email)
	return err
}

// IsSuppressed mengecek apakah email ada di suppression list dan aktif
func (r *SuppressionRepository) IsSuppressed(email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppressions WHERE email = $1 AND is_active = true)`
	var exists bool
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}
