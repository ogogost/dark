// dark - сверхбыстрый визуальный помощник для поиска скрытых предметов
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dark/internal/config"
	"dark/internal/game"
	"dark/internal/overlay"
	"dark/internal/vision"
)

func main() {
	log.Println("🌲 dark - Visual Helper starting...")

	// Инициализация менеджера игры
	cfgPath := "dark_config.json"
	manager, err := game.NewManager(cfgPath)
	if err != nil {
		log.Fatalf("Failed to initialize game manager: %v", err)
	}

	// Создаем компоненты
	textScanner := vision.NewTextScanner(time.Second)
	itemFinder := vision.NewItemFinder(500 * time.Millisecond)
	overlay := overlay.NewOverlay()

	// Фейковая функция захвата экрана (в реальной версии здесь будет скриншот)
	mockCapture := func() vision.ImageProvider {
		// В реальной реализации: возвращаем скриншот экрана
		return nil
	}

	// Фейковая функция отрисовки (в реальной версии здесь будет рендеринг оверлея)
	mockRender := func(items []vision.FoundItem, foundCount, totalCount int, elapsed time.Duration) {
		fmt.Printf("\r⏱️  %s | 📦 %d/%d found", elapsed.String(), foundCount, totalCount)
	}

	// Обработка сигналов
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск в режиме игры
	log.Println("🎮 Starting Game Mode (F1)...")
	manager.StartGameMode()

	// Получаем текущую локацию
	location := manager.GetCurrentLocation()
	if location == nil {
		log.Println("⚠️  No current location set. Please select a location first.")
	} else {
		log.Printf("📍 Current location: %s", location.Name)
	}

	// Запускаем сканер текста
	textScanner.Start(mockCapture)

	// Запускаем поиск предметов
	if location != nil {
		itemFinder.Start(mockCapture, manager.GetLocations(), location.ID)
	}

	// Запускаем оверлей
	overlay.Start(100*time.Millisecond, mockRender)

	log.Println("✅ dark is running! Press Ctrl+C to stop.")

	// Ожидание сигнала остановки
	<-sigChan

	log.Println("\n🛑 Shutting down...")

	// Остановка всех компонентов
	textScanner.Stop()
	itemFinder.Stop()
	overlay.Stop()
	manager.Stop()

	// Сохраняем конфигурацию
	if err := manager.SaveConfig(); err != nil {
		log.Printf("⚠️  Failed to save config: %v", err)
	} else {
		log.Println("💾 Configuration saved.")
	}

	log.Println("👋 Goodbye!")
}

// Пример создания тестовой конфигурации
func createSampleConfig() *config.Config {
	return &config.Config{
		Locations: []config.Location{
			{
				ID:         "forest_01",
				Name:       "Dark Forest",
				AnchorHash: 0x1234567890ABCDEF,
				WordHashes: []uint64{0xAABBCCDD, 0x11223344},
				Items: []config.Item{
					{ID: "apple", Name: "Red Apple", X: 800, Y: 600, Width: 50, Height: 50},
					{ID: "dagger", Name: "Old Dagger", X: 1200, Y: 400, Width: 60, Height: 40},
				},
			},
		},
		CurrentLocationID: "forest_01",
	}
}
