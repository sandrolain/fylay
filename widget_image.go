package fylay

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2/canvas"
	"github.com/coocood/freecache"
)

var (
	// Global image cache - 10MB
	imageCache = freecache.NewCache(10 * 1024 * 1024)
	// HTTP client with timeout
	httpClient = &http.Client{
		Timeout: 30 * time.Second,
	}
)

// ImageWidget wraps a canvas.Image with caching support
type ImageWidget struct {
	*canvas.Image
	src string
}

// NewImageWidget creates a new image widget
func NewImageWidget(src string) (*ImageWidget, error) {
	img := &ImageWidget{
		Image: canvas.NewImageFromFile(""),
		src:   src,
	}

	if err := img.Load(); err != nil {
		return nil, err
	}

	return img, nil
}

// Load loads the image from source (file or URL)
func (iw *ImageWidget) Load() error {
	// Check if it's a URL
	if strings.HasPrefix(iw.src, "http://") || strings.HasPrefix(iw.src, "https://") {
		return iw.loadFromURL()
	}

	// Load from file
	return iw.loadFromFile()
}

// loadFromFile loads image from local file
func (iw *ImageWidget) loadFromFile() error {
	absPath, err := filepath.Abs(iw.src)
	if err != nil {
		return fmt.Errorf("invalid file path: %w", err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", absPath)
	}

	iw.File = absPath
	iw.Refresh()

	return nil
}

// loadFromURL loads image from HTTP(S) URL with caching
func (iw *ImageWidget) loadFromURL() error {
	// Try to get from cache first
	cacheKey := []byte(iw.src)
	if cached, err := imageCache.Get(cacheKey); err == nil {
		// Cache hit
		img := canvas.NewImageFromReader(bytes.NewReader(cached), filepath.Base(iw.src))
		iw.Image = img
		return nil
	}

	// Cache miss - download from URL
	resp, err := httpClient.Get(iw.src)
	if err != nil {
		return fmt.Errorf("failed to download image: %w", err)
	}
	defer func() {
		_ = resp.Body.Close() //nolint:errcheck // Defer close error can be ignored
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read image data: %w", err)
	}

	// Store in cache (1 hour TTL)
	_ = imageCache.Set(cacheKey, data, 3600) //nolint:errcheck // Cache error can be ignored

	// Create image from data
	img := canvas.NewImageFromReader(bytes.NewReader(data), filepath.Base(iw.src))
	iw.Image = img

	return nil
}

// SetSource updates the image source and reloads
func (iw *ImageWidget) SetSource(src string) error {
	iw.src = src
	return iw.Load()
}

// GetSource returns the current image source
func (iw *ImageWidget) GetSource() string {
	return iw.src
}

// ClearImageCache clears the global image cache
func ClearImageCache() {
	imageCache.Clear()
}

// GetImageCacheStats returns cache statistics
func GetImageCacheStats() (entryCount int64, hitCount, missCount int64) {
	return imageCache.EntryCount(), imageCache.HitCount(), imageCache.MissCount()
}
