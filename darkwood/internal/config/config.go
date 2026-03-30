package config

// Location представляет одну локацию в базе знаний
type Location struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	AnchorHash  uint64   `json:"anchor_hash"`  // Хеш якоря (угол экрана)
	WordHashes  []uint64 `json:"word_hashes"`  // Хеши слов из списка заданий
	Items       []Item   `json:"items"`        // Координаты предметов (ROI)
}

// Item представляет предмет на локации
type Item struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	X, Y     int    `json:"x"`         // Координаты центра предмета
	Width    int    `json:"width"`     // Ширина ROI
	Height   int    `json:"height"`    // Высота ROI
	Found    bool   `json:"-"`         // Статус нахождения (не сохраняется в JSON)
}

// Config хранит всю базу знаний
type Config struct {
	Locations []Location `json:"locations"`
	CurrentLocationID string `json:"current_location_id"`
}

// DefaultConfig возвращает конфигурацию по умолчанию
func DefaultConfig() *Config {
	return &Config{
		Locations: make([]Location, 0, 200),
	}
}
