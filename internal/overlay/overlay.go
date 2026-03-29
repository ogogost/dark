package overlay

import (
	"sync"
	"time"

	"dark/internal/vision"
)

// Overlay представляет прозрачное окно поверх игры
type Overlay struct {
	mu         sync.RWMutex
	foundItems []vision.FoundItem
	foundCount int
	totalCount int
	startTime  time.Time
	isVisible  bool
	stopChan   chan struct{}
}

// NewOverlay создает новый оверлей
func NewOverlay() *Overlay {
	return &Overlay{
		foundItems: make([]vision.FoundItem, 0),
		stopChan:   make(chan struct{}),
	}
}

// Start запускает обновление оверлея
func (o *Overlay) Start(updateInterval time.Duration, renderFn func(items []vision.FoundItem, foundCount, totalCount int, elapsed time.Duration)) {
	o.mu.Lock()
	o.startTime = time.Now()
	o.isVisible = true
	o.mu.Unlock()

	go func() {
		ticker := time.NewTicker(updateInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				o.render(renderFn)
			case <-o.stopChan:
				return
			}
		}
	}()
}

// Stop останавливает оверлей
func (o *Overlay) Stop() {
	o.mu.Lock()
	o.isVisible = false
	o.mu.Unlock()
	close(o.stopChan)
}

// Update обновляет данные о найденных предметах
func (o *Overlay) Update(items []vision.FoundItem, totalCount int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.foundItems = items
	o.foundCount = len(items)
	o.totalCount = totalCount
}

// render вызывает функцию отрисовки с текущими данными
func (o *Overlay) render(renderFn func(items []vision.FoundItem, foundCount, totalCount int, elapsed time.Duration)) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	if !o.isVisible {
		return
	}

	elapsed := time.Since(o.startTime)
	renderFn(o.foundItems, o.foundCount, o.totalCount, elapsed)
}

// Reset сбрасывает статистику
func (o *Overlay) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.foundItems = make([]vision.FoundItem, 0)
	o.foundCount = 0
	o.totalCount = 0
	o.startTime = time.Now()
}

// GetStats возвращает текущую статистику
func (o *Overlay) GetStats() (foundCount, totalCount int, elapsed time.Duration) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	return o.foundCount, o.totalCount, time.Since(o.startTime)
}
