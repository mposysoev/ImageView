package main

import (
	"flag"
	"fmt"
	"github.com/mattn/go-sixel"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
)

type Config struct {
	width  int
	height int
	scale  float64
}

func main() {
	// Настройка параметров командной строки
	config := Config{}
	flag.IntVar(&config.width, "width", 0, "Output width (0 for auto)")
	flag.IntVar(&config.height, "height", 0, "Output height (0 for auto)")
	flag.Float64Var(&config.scale, "scale", 1.0, "Scale factor (default: 1.0)")

	// Парсинг аргументов
	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] image_file\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	imagePath := flag.Arg(0)
	if err := convertToSixel(imagePath, config); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func convertToSixel(filename string, config Config) error {
	// Проверка существования файла
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filename)
	}

	// Открытие файла
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Декодирование изображения
	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	// Создание энкодера
	enc := sixel.NewEncoder(os.Stdout)

	// Настройка параметров
	enc.Width = config.width
	enc.Height = config.height

	// Если указан масштаб, отличный от 1.0, применяем его
	if config.scale != 1.0 {
		bounds := img.Bounds()
		newWidth := int(float64(bounds.Dx()) * config.scale)
		newHeight := int(float64(bounds.Dy()) * config.scale)
		if newWidth > 0 && newHeight > 0 {
			enc.Width = newWidth
			enc.Height = newHeight
		}
	}

	// Кодирование и вывод
	if err := enc.Encode(img); err != nil {
		return fmt.Errorf("failed to encode sixel: %v", err)
	}

	return nil
}
