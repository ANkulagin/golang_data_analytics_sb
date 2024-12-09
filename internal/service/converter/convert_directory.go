package converter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertDirectory(srcDir, destDir string) error {
	// Проверка существования исходной директории
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return fmt.Errorf("исходная директория не существует: %s", srcDir)
	}

	// Создание целевой директории, если она отсутствует
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("не удалось создать целевую директорию: %v", err)
	}

	// Проход по всем файлам в исходной директории
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			if err := c.ConvertFile(path, srcDir, destDir); err != nil {
				return fmt.Errorf("ошибка при конвертации %s: %v", path, err)
			}
		}
		return nil
	})

}
