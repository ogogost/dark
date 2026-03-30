package game

import (
	"encoding/json"
	"os"
	"sync"

	"dark/internal/config"
)

// Manager управляет состоянием игры и режимами работы
type Manager struct {
	mu            sync.RWMutex
	cfg           *config.Config
	configPath    string
	currentMode   Mode
	isRunning     bool
	locationIndex map[string]*config.Location
}

// Mode представляет режим работы программы
type Mode int

const (
	ModeNone  Mode = iota
	ModeGame       // F1 - Игра: бот работает в фоне, подсвечивает предметы
	ModeLearn      // F2 - Обучение: запоминание новых предметов
)

// NewManager создает новый менеджер игры
func NewManager(configPath string) (*Manager, error) {
	m := &Manager{
		configPath:    configPath,
		currentMode:   ModeNone,
		locationIndex: make(map[string]*config.Location),
	}

	// Загружаем или создаем конфигурацию
	if err := m.loadConfig(); err != nil {
		return nil, err
	}

	// Строим индекс локаций
	for i := range m.cfg.Locations {
		m.locationIndex[m.cfg.Locations[i].ID] = &m.cfg.Locations[i]
	}

	return m, nil
}

// loadConfig загружает конфигурацию из файла
func (m *Manager) loadConfig() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Создаем новую конфигурацию
			m.cfg = config.DefaultConfig()
			return nil
		}
		return err
	}

	m.cfg = &config.Config{}
	return json.Unmarshal(data, m.cfg)
}

// SaveConfig сохраняет конфигурацию в файл
func (m *Manager) SaveConfig() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data, err := json.MarshalIndent(m.cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.configPath, data, 0644)
}

// SetMode устанавливает режим работы
func (m *Manager) SetMode(mode Mode) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.currentMode = mode
}

// GetMode возвращает текущий режим работы
func (m *Manager) GetMode() Mode {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.currentMode
}

// SetCurrentLocation устанавливает текущую локацию
func (m *Manager) SetCurrentLocation(locationID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.locationIndex[locationID]; !exists {
		return false
	}

	m.cfg.CurrentLocationID = locationID
	return true
}

// GetCurrentLocation возвращает текущую локацию
func (m *Manager) GetCurrentLocation() *config.Location {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.cfg.CurrentLocationID == "" {
		return nil
	}

	return m.locationIndex[m.cfg.CurrentLocationID]
}

// GetLocations возвращает все локации
func (m *Manager) GetLocations() []config.Location {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.cfg.Locations
}

// AddLocation добавляет новую локацию
func (m *Manager) AddLocation(location config.Location) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cfg.Locations = append(m.cfg.Locations, location)
	m.locationIndex[location.ID] = &m.cfg.Locations[len(m.cfg.Locations)-1]
}

// AddItemToLocation добавляет предмет к локации
func (m *Manager) AddItemToLocation(locationID string, item config.Item) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	location, exists := m.locationIndex[locationID]
	if !exists {
		return false
	}

	location.Items = append(location.Items, item)
	return true
}

// IsRunning проверяет, запущен ли бот
func (m *Manager) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.isRunning
}

// SetRunning устанавливает статус запуска
func (m *Manager) SetRunning(running bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isRunning = running
}

// StartGameMode запускает режим игры
func (m *Manager) StartGameMode() {
	m.SetMode(ModeGame)
	m.SetRunning(true)
}

// StartLearnMode запускает режим обучения
func (m *Manager) StartLearnMode() {
	m.SetMode(ModeLearn)
	m.SetRunning(true)
}

// Stop останавливает текущий режим
func (m *Manager) Stop() {
	m.SetRunning(false)
	m.SetMode(ModeNone)
}
