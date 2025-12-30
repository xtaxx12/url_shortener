package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/username/url-shortener/internal/domain"
	"github.com/username/url-shortener/internal/service"
)

// URLHandler handles HTTP requests for URL operations (Port)
type URLHandler struct {
	urlService *service.URLService
}

// NewURLHandler creates a new URL handler
func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

// CreateShortURL handles POST /api/shorten
func (h *URLHandler) CreateShortURL(c *gin.Context) {
	var req domain.CreateURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	response, err := h.urlService.CreateShortURL(req.URL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create short URL",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// RedirectURL handles GET /:code
func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortCode := c.Param("code")

	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	originalURL, err := h.urlService.GetOriginalURL(shortCode)
	if err != nil {
		if err.Error() == "url not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "URL not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve URL",
		})
		return
	}

	// Redirect with 301 (Permanent Redirect) for SEO
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetURLStats handles GET /api/stats/:code
func (h *URLHandler) GetURLStats(c *gin.Context) {
	shortCode := c.Param("code")

	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Short code is required",
		})
		return
	}

	stats, err := h.urlService.GetURLStats(shortCode)
	if err != nil {
		if err.Error() == "url not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "URL not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve statistics",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}
