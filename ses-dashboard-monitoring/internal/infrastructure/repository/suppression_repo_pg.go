package repository

import (
	"context"
	"database/sql"
	"strings"

	"ses-monitoring/internal/domain/suppression"
)

type suppressionRepo struct {
	db *sql.DB
}

func NewSuppressionRepository(db *sql.DB) suppression.Repository {
	return &suppressionRepo{db: db}
}

func (r *suppressionRepo) Add(ctx context.Context, entry *suppression.SuppressionEntry) error {
	query := `
		INSERT INTO suppression_list (email, suppression_type, reason, aws_status, added_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (email) 
		DO UPDATE SET 
			suppression_type = $2, 
			reason = $3, 
			aws_status = $4, 
			is_active = true,
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, 
		entry.Email, 
		entry.SuppressionType, 
		entry.Reason, 
		entry.AWSStatus, 
		entry.AddedBy,
	)
	return err
}

func (r *suppressionRepo) Remove(ctx context.Context, email string) error {
	query := `UPDATE suppression_list SET is_active = false, updated_at = NOW() WHERE email = $1`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}

func (r *suppressionRepo) GetAll(ctx context.Context, limit, offset int) ([]*suppression.SuppressionEntry, error) {
	query := `
		SELECT s.id, s.email, s.suppression_type, s.reason, s.aws_status, s.is_active, 
		       s.added_by, COALESCE(u.username, 'System') as added_by_name, 
		       s.created_at, s.updated_at
		FROM suppression_list s
		LEFT JOIN users u ON s.added_by = u.id
		WHERE s.is_active = true
		ORDER BY s.created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*suppression.SuppressionEntry
	for rows.Next() {
		e := &suppression.SuppressionEntry{}
		err := rows.Scan(
			&e.ID, &e.Email, &e.SuppressionType, &e.Reason, &e.AWSStatus, &e.IsActive,
			&e.AddedBy, &e.AddedByName, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *suppressionRepo) GetCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM suppression_list WHERE is_active = true`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *suppressionRepo) Search(ctx context.Context, query string, limit, offset int) ([]*suppression.SuppressionEntry, error) {
	sqlQuery := `
		SELECT s.id, s.email, s.suppression_type, s.reason, s.aws_status, s.is_active, 
		       s.added_by, COALESCE(u.username, 'System') as added_by_name, 
		       s.created_at, s.updated_at
		FROM suppression_list s
		LEFT JOIN users u ON s.added_by = u.id
		WHERE s.is_active = true 
		AND (s.email ILIKE $1 OR s.reason ILIKE $1)
		ORDER BY s.created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.QueryContext(ctx, sqlQuery, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*suppression.SuppressionEntry
	for rows.Next() {
		e := &suppression.SuppressionEntry{}
		err := rows.Scan(
			&e.ID, &e.Email, &e.SuppressionType, &e.Reason, &e.AWSStatus, &e.IsActive,
			&e.AddedBy, &e.AddedByName, &e.CreatedAt, &e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *suppressionRepo) GetSearchCount(ctx context.Context, query string) (int, error) {
	sqlQuery := `
		SELECT COUNT(*) FROM suppression_list 
		WHERE is_active = true 
		AND (email ILIKE $1 OR reason ILIKE $1)
	`
	searchTerm := "%" + strings.ToLower(query) + "%"
	var count int
	err := r.db.QueryRowContext(ctx, sqlQuery, searchTerm).Scan(&count)
	return count, err
}

func (r *suppressionRepo) IsSupressed(ctx context.Context, email string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM suppression_list WHERE email = $1 AND is_active = true)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *suppressionRepo) UpdateAWSStatus(ctx context.Context, email string, status suppression.AWSStatus) error {
	query := `UPDATE suppression_list SET aws_status = $1, updated_at = NOW() WHERE email = $2`
	_, err := r.db.ExecContext(ctx, query, status, email)
	return err
}

func (r *suppressionRepo) GetUnsyncedEntries(ctx context.Context) ([]*suppression.SuppressionEntry, error) {
	query := `
		SELECT id, email, suppression_type, reason, aws_status, is_active, 
		       added_by, created_at, updated_at, synced_at
		FROM suppression_list 
		WHERE is_active = true AND (synced_at IS NULL OR updated_at > synced_at)
		ORDER BY created_at ASC
	`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*suppression.SuppressionEntry
	for rows.Next() {
		e := &suppression.SuppressionEntry{}
		err := rows.Scan(
			&e.ID, &e.Email, &e.SuppressionType, &e.Reason, &e.AWSStatus, &e.IsActive,
			&e.AddedBy, &e.CreatedAt, &e.UpdatedAt, &e.SyncedAt,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func (r *suppressionRepo) MarkAsSynced(ctx context.Context, email string) error {
	query := `UPDATE suppression_list SET synced_at = NOW() WHERE email = $1`
	_, err := r.db.ExecContext(ctx, query, email)
	return err
}