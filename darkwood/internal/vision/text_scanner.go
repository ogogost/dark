package vision

import (
	"image"
	"sync"
	"time"

	"dark/internal/config"
)

// TextScanner сканирует текст в нижней части экрана (подвал игры)
type TextScanner struct {
	hasher            *Hasher
	checkInterval     time.Duration
	mu                sync.RWMutex
	currentWordHashes []uint64
	stopChan          chan struct{}
}

// NewTextScanner создает новый сканер текста
func NewTextScanner(checkInterval time.Duration) *TextScanner {
	return &TextScanner{
		hasher:        NewHasher(),
		checkInterval: checkInterval,
		stopChan:      make(chan struct{}),
	}
}

// Start запускает сканер текста в отдельной горутине
func (s *TextScanner) Start(captureFn func() ImageProvider) {
	go func() {
		ticker := time.NewTicker(s.checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.scan(captureFn)
			case <-s.stopChan:
				return
			}
		}
	}()
}

// Stop останавливает сканер текста
func (s *TextScanner) Stop() {
	close(s.stopChan)
}

// scan выполняет одно сканирование
func (s *TextScanner) scan(captureFn func() ImageProvider) {
	// Захватываем изображение экрана
	screen := captureFn()
	if screen == nil {
		return
	}

	// Получаем размеры экрана
	bounds := screen.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Определяем область "подвала" (нижние 15% экрана)
	footerHeight := height / 7
	footerY := height - footerHeight

	// Разбиваем подвал на несколько зон для проверки слов
	// Предполагаем, что слова расположены горизонтально
	zoneWidth := width / 5
	zoneHeight := footerHeight / 2

	var detectedHashes []uint64

	// Сканируем 5 зон в подвале
	for i := 0; i < 5; i++ {
		x := i * zoneWidth
		y := footerY + zoneHeight/4

		// Вырезаем зону
		zone := Crop(screen, x, y, zoneWidth, zoneHeight/2)

		// Вычисляем хеш
		hash := s.hasher.ComputeHash(zone)

		if hash != 0 {
			detectedHashes = append(detectedHashes, hash)
		}
	}

	// Обновляем текущие найденные хеши
	s.mu.Lock()
	s.currentWordHashes = detectedHashes
	s.mu.Unlock()
}

// GetCurrentHashes возвращает текущие обнаруженные хеши слов
func (s *TextScanner) GetCurrentHashes() []uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]uint64, len(s.currentWordHashes))
	copy(result, s.currentWordHashes)
	return result
}

// MatchItems сравнивает обнаруженные хеши с базой знаний и возвращает ID предметов для поиска
func (s *TextScanner) MatchItems(locations []config.Location, currentLocationID string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var itemIDs []string

	// Ищем текущую локацию
	var currentLocation *config.Location
	for i := range locations {
		if locations[i].ID == currentLocationID {
			currentLocation = &locations[i]
			break
		}
	}

	if currentLocation == nil {
		return itemIDs
	}

	// Сравниваем хеши
	for _, detectedHash := range s.currentWordHashes {
		for _, knownHash := range currentLocation.WordHashes {
			if detectedHash == knownHash {
				// Нашли совпадение - добавляем все предметы этого типа
				for _, item := range currentLocation.Items {
					itemIDs = append(itemIDs, item.ID)
				}
				break
			}
		}
	}

	return itemIDs
}

// ImageProvider интерфейс для захвата изображения
type ImageProvider interface {
	image.Image
}
