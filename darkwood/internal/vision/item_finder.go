package vision

import (
	"sync"
	"time"

	"darkwood/internal/config"
)

// ItemFinder ищет предметы по заранее известным координатам
type ItemFinder struct {
	hasher        *Hasher
	checkInterval time.Duration
	mu            sync.RWMutex
	targetItemIDs []string
	foundItems    map[string]config.Item
	stopChan      chan struct{}
}

// FoundItem представляет найденный предмет
type FoundItem struct {
	ID       string
	Name     string
	X, Y     int
	FoundAt  time.Time
}

// NewItemFinder создает новый поисковик предметов
func NewItemFinder(checkInterval time.Duration) *ItemFinder {
	return &ItemFinder{
		hasher:        NewHasher(),
		checkInterval: checkInterval,
		foundItems:    make(map[string]config.Item),
		stopChan:      make(chan struct{}),
	}
}

// Start запускает поиск предметов в отдельной горутине
func (f *ItemFinder) Start(captureFn func() ImageProvider, locations []config.Location, currentLocationID string) {
	go func() {
		ticker := time.NewTicker(f.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				f.scan(captureFn, locations, currentLocationID)
			case <-f.stopChan:
				return
			}
		}
	}()
}

// Stop останавливает поиск предметов
func (f *ItemFinder) Stop() {
	close(f.stopChan)
}

// SetTargetItems устанавливает список ID предметов для поиска
func (f *ItemFinder) SetTargetItems(itemIDs []string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	f.targetItemIDs = itemIDs
	// Сбрасываем статус найденных предметов при смене цели
	f.foundItems = make(map[string]config.Item)
}

// scan выполняет одно сканирование предметов
func (f *ItemFinder) scan(captureFn func() ImageProvider, locations []config.Location, currentLocationID string) {
	f.mu.RLock()
	targetIDs := make([]string, len(f.targetItemIDs))
	copy(targetIDs, f.targetItemIDs)
	f.mu.RUnlock()

	if len(targetIDs) == 0 {
		return
	}

	// Захватываем изображение экрана
	screen := captureFn()
	if screen == nil {
		return
	}

	// Ищем текущую локацию
	var currentLocation *config.Location
	for i := range locations {
		if locations[i].ID == currentLocationID {
			currentLocation = &locations[i]
			break
		}
	}

	if currentLocation == nil {
		return
	}

	// Проверяем каждый целевой предмет
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, itemID := range targetIDs {
		// Ищем предмет в базе знаний
		var targetItem *config.Item
		for i := range currentLocation.Items {
			if currentLocation.Items[i].ID == itemID {
				targetItem = &currentLocation.Items[i]
				break
			}
		}

		if targetItem == nil {
			continue
		}

		// Если уже найден, пропускаем
		if _, exists := f.foundItems[itemID]; exists {
			continue
		}

		// Вырезаем ROI предмета
		roiX := targetItem.X - targetItem.Width/2
		roiY := targetItem.Y - targetItem.Height/2
		
		roi := Crop(screen, roiX, roiY, targetItem.Width, targetItem.Height)
		
		// Вычисляем хеш текущей области
		currentHash := f.hasher.ComputeHash(roi)
		
		// Для первого раза просто сохраняем хеш как эталон
		// В реальной реализации здесь должно быть сравнение с эталонным хешем
		if currentHash != 0 {
			f.foundItems[itemID] = *targetItem
		}
	}
}

// GetFoundItems возвращает список найденных предметов
func (f *ItemFinder) GetFoundItems() []FoundItem {
	f.mu.RLock()
	defer f.mu.RUnlock()

	result := make([]FoundItem, 0, len(f.foundItems))
	for id, item := range f.foundItems {
		result = append(result, FoundItem{
			ID:      id,
			Name:    item.Name,
			X:       item.X,
			Y:       item.Y,
			FoundAt: time.Now(),
		})
	}

	return result
}

// IsItemFound проверяет, найден ли конкретный предмет
func (f *ItemFinder) IsItemFound(itemID string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	_, exists := f.foundItems[itemID]
	return exists
}

// Reset сбрасывает статус найденных предметов
func (f *ItemFinder) Reset() {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	f.foundItems = make(map[string]config.Item)
}
