package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/username/url-shortener/internal/domain"
)

// PostgresRepository implements domain.URLRepository (Adapter)
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a new URL record
func (r *PostgresRepository) Create(url *domain.URL) error {
	query := `
		INSERT INTO urls (short_code, original_url, clicks, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now().UTC()
	url.CreatedAt = now
	url.UpdatedAt = now
	url.Clicks = 0

	err := r.db.QueryRow(
		query,
		url.ShortCode,
		url.OriginalURL,
		url.Clicks,
		url.CreatedAt,
		url.UpdatedAt,
	).Scan(&url.ID)

	return err
}

// FindByShortCode retrieves a URL by its short code
func (r *PostgresRepository) FindByShortCode(shortCode string) (*domain.URL, error) {
	query := `
		SELECT id, short_code, original_url, clicks, created_at, updated_at
		FROM urls
		WHERE short_code = $1
	`

	url := &domain.URL{}
	err := r.db.QueryRow(query, shortCode).Scan(
		&url.ID,
		&url.ShortCode,
		&url.OriginalURL,
		&url.Clicks,
		&url.CreatedAt,
		&url.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return url, nil
}

// IncrementClicks increases the click counter for a URL
func (r *PostgresRepository) IncrementClicks(shortCode string) error {
	query := `
		UPDATE urls
		SET clicks = clicks + 1, updated_at = $1
		WHERE short_code = $2
	`

	result, err := r.db.Exec(query, time.Now().UTC(), shortCode)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("url not found")
	}

	return nil
}

// GetStats retrieves URL statistics
func (r *PostgresRepository) GetStats(shortCode string) (*domain.URL, error) {
	return r.FindByShortCode(shortCode)
}
