package main

import (
	"flag"
	"fmt"
	"github.com/ANkulagin/golang_markdown_converter_sb/internal/config"
	"log"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "Путь к конфигурационному файлу")
	flag.Parse()
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Не удалось загрузить конфигурацию: %v", err)
	}
	fmt.Printf("src_dir: %s\n", cfg.SrcDir)
	fmt.Printf("dest_dir: %s\n", cfg.DestDir)
	fmt.Printf("log_level: %s\n", cfg.LogLevel)
}
