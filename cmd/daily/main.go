package main

import (
	"flag"
	"github.com/ANkulagin/golang_markdown_converter_sb/internal/config"
	"github.com/ANkulagin/golang_markdown_converter_sb/internal/service/converter"
	"log"
	"path/filepath"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Путь к конфигурационному файлу")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Преобразование относительных путей в абсолютные
	absSrcDir, err := filepath.Abs(cfg.SrcDir)
	if err != nil {
		log.Fatalf("Не удалось определить абсолютный путь для исходной директории: %v", err)
	}

	absDestDir, err := filepath.Abs(cfg.DestDir)
	if err != nil {
		log.Fatalf("Не удалось определить абсолютный путь для целевой директории: %v", err)
	}

	log.Printf("Конвертация заметок из %s в %s\n", absSrcDir, absDestDir)
	log.Printf("Уровень логирования: %s\n", cfg.LogLevel)

	conv := converter.NewConverter()

	if err := conv.ConvertDirectory(absSrcDir, absDestDir); err != nil {
		log.Fatalf("Конвертация не удалась: %v", err)
	}

	log.Println("Конвертация завершена успешно.")
}
