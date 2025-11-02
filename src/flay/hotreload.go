package flay

import (
	"fmt"
	"log"
	"path/filepath"
	"sync"

	"fyne.io/fyne/v2"
	"github.com/fsnotify/fsnotify"
)

// HotReloadConfig configura il comportamento del hot reload
type HotReloadConfig struct {
	Enabled     bool
	LayoutPath  string
	OnReload    func(fyne.CanvasObject)
	OnError     func(error)
	DebugLog    bool
	watchMutex  sync.Mutex
	watcher     *fsnotify.Watcher
	stopChannel chan bool
}

// NewHotReloadConfig crea una nuova configurazione di hot reload
func NewHotReloadConfig(layoutPath string) *HotReloadConfig {
	return &HotReloadConfig{
		Enabled:     true,
		LayoutPath:  layoutPath,
		DebugLog:    false,
		stopChannel: make(chan bool),
	}
}

// EnableHotReload attiva il hot reload per un layout file
func (b *Builder) EnableHotReload(config *HotReloadConfig) error {
	if !config.Enabled {
		return nil
	}

	// Resolve absolute path
	absPath, err := filepath.Abs(config.LayoutPath)
	if err != nil {
		return fmt.Errorf("invalid layout path: %w", err)
	}

	// Create file watcher
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}

	config.watcher = watcher
	config.watchMutex.Lock()
	defer config.watchMutex.Unlock()

	// Watch the layout file
	if err := watcher.Add(absPath); err != nil {
		_ = watcher.Close() //nolint:errcheck // Close error can be ignored
		return fmt.Errorf("failed to watch file: %w", err)
	}

	if config.DebugLog {
		log.Printf("[HotReload] Watching: %s\n", absPath)
	}

	// Start watching in goroutine
	go b.watchFileChanges(config)

	return nil
}

// watchFileChanges monitors files for changes and triggers callbacks
//
//nolint:gocognit // Complex but necessary for hot reload implementation
func (b *Builder) watchFileChanges(config *HotReloadConfig) {
	defer func() {
		if config.watcher != nil {
			_ = config.watcher.Close() //nolint:errcheck // Close error can be ignored
		}
	}()

	for {
		select {
		case event, ok := <-config.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
				if config.DebugLog {
					log.Printf("[HotReload] File changed: %s\n", event.Name)
				}

				// Reload layout
				if err := b.reloadLayout(config); err != nil {
					if config.OnError != nil {
						config.OnError(err)
					} else if config.DebugLog {
						log.Printf("[HotReload] Error reloading: %v\n", err)
					}
				}
			}

		case err, ok := <-config.watcher.Errors:
			if !ok {
				return
			}

			if config.OnError != nil {
				config.OnError(fmt.Errorf("watcher error: %w", err))
			} else if config.DebugLog {
				log.Printf("[HotReload] Watcher error: %v\n", err)
			}

		case <-config.stopChannel:
			if config.DebugLog {
				log.Println("[HotReload] Stopping watcher")
			}
			return
		}
	}
}

// reloadLayout reloads the layout file and calls the callback
func (b *Builder) reloadLayout(config *HotReloadConfig) error {
	// Load new layout
	newBuilder, content, err := LoadLayoutFromFile(config.LayoutPath)
	if err != nil {
		return fmt.Errorf("failed to load layout: %w", err)
	}

	// Copy event handlers and contexts from current builder
	newBuilder.eventHandler = b.eventHandler
	newBuilder.eventCallbacks = b.eventCallbacks
	newBuilder.entryCallbacks = b.entryCallbacks
	newBuilder.bindingContext = b.bindingContext
	newBuilder.templateContext = b.templateContext

	// Call reload callback
	if config.OnReload != nil {
		config.OnReload(content)
	}

	// Update current builder reference
	*b = *newBuilder

	if config.DebugLog {
		log.Println("[HotReload] Layout reloaded successfully")
	}

	return nil
}

// StopHotReload ferma il hot reload
func (config *HotReloadConfig) Stop() {
	if config.stopChannel != nil {
		config.stopChannel <- true
		close(config.stopChannel)
	}
}

// SimpleHotReload è una funzione helper per casi d'uso semplici
func (b *Builder) SimpleHotReload(layoutPath string, window fyne.Window) error {
	config := NewHotReloadConfig(layoutPath)
	config.DebugLog = true
	config.OnReload = func(content fyne.CanvasObject) {
		window.SetContent(content)
		log.Println("[HotReload] UI updated")
	}
	config.OnError = func(err error) {
		log.Printf("[HotReload] Error: %v\n", err)
	}

	return b.EnableHotReload(config)
}
