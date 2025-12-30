package domain

import "time"

// URL represents the core domain entity
type URL struct {
	ID          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// MaxURLLength defines the maximum allowed URL length to prevent abuse
const MaxURLLength = 2048

// CreateURLRequest represents the request to create a short URL
// 🔒 SECURITY: Limit URL length to prevent payload attacks
type CreateURLRequest struct {
	URL string `json:"url" binding:"required,url,max=2048"`
}

// CreateURLResponse represents the response after creating a short URL
type CreateURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// URLStatsResponse represents the statistics of a short URL
type URLStatsResponse struct {
	ShortCode   string    `json:"short_code"`
	OriginalURL string    `json:"original_url"`
	Clicks      int64     `json:"clicks"`
	CreatedAt   time.Time `json:"created_at"`
}

// URLRepository defines the interface for URL storage operations (Port)
type URLRepository interface {
	Create(url *URL) error
	FindByShortCode(shortCode string) (*URL, error)
	IncrementClicks(shortCode string) error
	GetStats(shortCode string) (*URL, error)
}

// CacheRepository defines the interface for caching operations (Port)
type CacheRepository interface {
	Get(key string) (string, error)
	Set(key string, value string) error
	Delete(key string) error
}
