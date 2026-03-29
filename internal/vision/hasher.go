package vision

import (
	"image"
	"image/color"
)

// Hasher предоставляет методы для вычисления хешей изображений
type Hasher struct{}

// NewHasher создает новый экземпляр Hasher
func NewHasher() *Hasher {
	return &Hasher{}
}

// ComputeHash вычисляет простой perceptual hash для изображения
// Возвращает 64-битный хеш для быстрого сравнения
func (h *Hasher) ComputeHash(img image.Image) uint64 {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width == 0 || height == 0 {
		return 0
	}

	// Упрощенный хеш: усредняем яркость по квадрантам
	var hash uint64
	
	// Разбиваем на 8x8 grid для получения 64 бит
	stepX := width / 8
	stepY := height / 8
	
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			px := x * stepX + stepX/2
			py := y * stepY + stepY/2
			
			if px >= width || py >= height {
				continue
			}
			
			c := img.At(px, py)
			r, g, b, _ := c.RGBA()
			
			// Вычисляем яркость
			luminance := (0.299*float64(r>>8) + 0.587*float64(g>>8) + 0.114*float64(b>>8)) / 255.0
			
			if luminance > 0.5 {
				hash |= (1 << (y*8 + x))
			}
		}
	}
	
	return hash
}

// CompareHashes сравнивает два хеша и возвращает процент схожести (0.0 - 1.0)
func (h *Hasher) CompareHashes(hash1, hash2 uint64) float64 {
	diff := hash1 ^ hash2
	
	// Считаем количество одинаковых битов
	sameBits := 64 - countBits(diff)
	
	return float64(sameBits) / 64.0
}

// countBits считает количество установленных битов в числе
func countBits(n uint64) int {
	count := 0
	for n != 0 {
		count += int(n & 1)
		n >>= 1
	}
	return count
}

// Crop вырезает область интереса из изображения
func Crop(img image.Image, x, y, width, height int) image.Image {
	bounds := img.Bounds()
	
	// Ограничиваем координаты границами изображения
	x = max(x, bounds.Min.X)
	y = max(y, bounds.Min.Y)
	width = min(width, bounds.Max.X-x)
	height = min(height, bounds.Max.Y-y)
	
	return img.(interface {
		SubImage(image.Rectangle) image.Image
	}).SubImage(image.Rect(x, y, x+width, y+height))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Color представляет цвет в формате RGBA
type Color struct {
	R, G, B, A uint8
}

// ToColor конвертирует color.Color в наш Color
func ToColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(a >> 8),
	}
}
