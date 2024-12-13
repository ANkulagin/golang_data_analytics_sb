package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/ANkulagin/golang_markdown_converter_sb/internal/config"
	"github.com/ANkulagin/golang_markdown_converter_sb/internal/service/converter"

	log "github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Путь к конфигурационному файлу")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}

	// Настройка уровня логирования
	level, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Fatalf("Не удалось установить уровень логирования: %v", err)
	}
	log.SetLevel(level)

	// Настройка формата логирования
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		// Можно добавить другие настройки
		ForceColors: true,
	})

	// Настройка вывода логов (можно перенаправить в файл, если нужно)
	log.SetOutput(os.Stdout)

	// Преобразование относительных путей в абсолютные
	absSrcDir, err := filepath.Abs(cfg.SrcDir)
	if err != nil {
		log.Fatalf("Не удалось определить абсолютный путь для исходной директории: %v", err)
	}

	absDestDir, err := filepath.Abs(cfg.DestDir)
	if err != nil {
		log.Fatalf("Не удалось определить абсолютный путь для целевой директории: %v", err)
	}

	log.Infof("Конвертация заметок из %s в %s", absSrcDir, absDestDir)
	log.Infof("Уровень логирования: %s", cfg.LogLevel)

	conv := converter.NewConverter()

	if err := conv.ConvertDirectory(absSrcDir, absDestDir); err != nil {
		log.Fatalf("Конвертация не удалась: %v", err)
	}

	log.Info("Конвертация завершена успешно.")
}
