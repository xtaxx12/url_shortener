package service

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/username/url-shortener/internal/domain"
	"github.com/username/url-shortener/pkg/shortener"
)

// 🔒 SECURITY: Blocked URL schemes to prevent javascript: and data: attacks
var blockedSchemes = []string{"javascript", "data", "vbscript", "file"}

// URLService handles URL shortening business logic (Application Layer)
type URLService struct {
	urlRepo   domain.URLRepository
	cacheRepo domain.CacheRepository
	baseURL   string
}

// NewURLService creates a new URL service
func NewURLService(urlRepo domain.URLRepository, cacheRepo domain.CacheRepository) *URLService {
	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost"
	}

	return &URLService{
		urlRepo:   urlRepo,
		cacheRepo: cacheRepo,
		baseURL:   baseURL,
	}
}

// validateURL checks if the URL is safe to redirect to
// 🔒 SECURITY: Prevents Open Redirect and malicious scheme attacks
func validateURL(rawURL string) error {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.New("invalid URL format")
	}

	// Check for blocked schemes
	scheme := strings.ToLower(parsedURL.Scheme)
	for _, blocked := range blockedSchemes {
		if scheme == blocked {
			return errors.New("URL scheme not allowed")
		}
	}

	// Ensure URL has a valid scheme
	if scheme != "http" && scheme != "https" {
		return errors.New("only http and https URLs are allowed")
	}

	// Ensure URL has a host
	if parsedURL.Host == "" {
		return errors.New("URL must have a valid host")
	}

	return nil
}

// CreateShortURL creates a new shortened URL
func (s *URLService) CreateShortURL(originalURL string) (*domain.CreateURLResponse, error) {
	// 🔒 SECURITY: Validate URL before storing
	if err := validateURL(originalURL); err != nil {
		return nil, err
	}

	// Generate unique short code
	shortCode := shortener.GenerateCode(7)

	// Check for collision (unlikely but possible)
	existing, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	// Regenerate if collision detected
	attempts := 0
	for existing != nil && attempts < 5 {
		shortCode = shortener.GenerateCode(7)
		existing, err = s.urlRepo.FindByShortCode(shortCode)
		if err != nil {
			return nil, err
		}
		attempts++
	}

	if existing != nil {
		return nil, errors.New("failed to generate unique short code")
	}

	// Create URL record
	url := &domain.URL{
		ShortCode:   shortCode,
		OriginalURL: originalURL,
	}

	if err := s.urlRepo.Create(url); err != nil {
		return nil, err
	}

	// Cache the URL for fast redirects
	if err := s.cacheRepo.Set(shortCode, originalURL); err != nil {
		log.Printf("Failed to cache URL: %v", err)
		// Don't fail the request, just log the cache error
	}

	return &domain.CreateURLResponse{
		ShortCode:   shortCode,
		ShortURL:    s.baseURL + "/s/" + shortCode,
		OriginalURL: originalURL,
	}, nil
}

// GetOriginalURL retrieves the original URL for redirection
func (s *URLService) GetOriginalURL(shortCode string) (string, error) {
	// Try cache first (fast path)
	if cachedURL, err := s.cacheRepo.Get(shortCode); err == nil && cachedURL != "" {
		// Increment click count asynchronously
		go func() {
			if err := s.urlRepo.IncrementClicks(shortCode); err != nil {
				log.Printf("Failed to increment clicks: %v", err)
			}
		}()
		return cachedURL, nil
	}

	// Cache miss, query database
	url, err := s.urlRepo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	if url == nil {
		return "", errors.New("url not found")
	}

	// Populate cache for next request
	if err := s.cacheRepo.Set(shortCode, url.OriginalURL); err != nil {
		log.Printf("Failed to cache URL: %v", err)
	}

	// Increment click count
	if err := s.urlRepo.IncrementClicks(shortCode); err != nil {
		log.Printf("Failed to increment clicks: %v", err)
	}

	return url.OriginalURL, nil
}

// GetURLStats retrieves statistics for a shortened URL
func (s *URLService) GetURLStats(shortCode string) (*domain.URLStatsResponse, error) {
	url, err := s.urlRepo.GetStats(shortCode)
	if err != nil {
		return nil, err
	}

	if url == nil {
		return nil, errors.New("url not found")
	}

	return &domain.URLStatsResponse{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		Clicks:      url.Clicks,
		CreatedAt:   url.CreatedAt,
	}, nil
}
