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

	// Buat query untuk bulk upsert
	valueStrings := make([]string, 0, len(suppressions))
	valueArgs := make([]interface{}, 0, len(suppressions)*4)
	
	for i, s := range suppressions {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", 
			i*4+1, i*4+2, i*4+3, i*4+4))
		valueArgs = append(valueArgs, s.Email, s.Reason, s.Source, s.UpdatedAt)
	}

	query := fmt.Sprintf(`
		INSERT INTO suppressions (email, reason, source, updated_at) 
		VALUES %s
		ON CONFLICT (email) 
		DO UPDATE SET 
			reason = EXCLUDED.reason,
			source = EXCLUDED.source,
			updated_at = EXCLUDED.updated_at
	`, strings.Join(valueStrings, ","))

	log.Printf("Executing query with %d parameters", len(valueArgs))

	result, err := r.db.Exec(query, valueArgs...)
	if err != nil {
		log.Printf("Bulk upsert failed: %v", err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Bulk upsert completed: %d rows affected", rowsAffected)

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