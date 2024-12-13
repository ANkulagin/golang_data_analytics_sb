package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertDirectory(srcDir, destDir string) error {
	// Проверка существования исходной директории
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		log.Errorf("Исходная директория не существует: %s", srcDir)
		return fmt.Errorf("исходная директория не существует: %s", srcDir)
	}

	// Создание целевой директории, если она отсутствует
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		log.Errorf("Не удалось создать целевую директорию: %v", err)
		return fmt.Errorf("не удалось создать целевую директорию: %v", err)
	}

	log.Infof("Начало конвертации директории: %s -> %s", srcDir, destDir)

	// Проход по всем файлам в исходной директории
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Errorf("Ошибка при обходе файла %s: %v", path, err)
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			log.WithFields(log.Fields{
				"file": path,
			}).Info("Начало конвертации файла")
			if err := c.ConvertFile(path, srcDir, destDir); err != nil {
				log.Errorf("Ошибка при конвертации %s: %v", path, err)
				return fmt.Errorf("ошибка при конвертации %s: %v", path, err)
			}
			log.WithFields(log.Fields{
				"file": path,
			}).Info("Файл успешно конвертирован")
		}
		return nil
	})
}
