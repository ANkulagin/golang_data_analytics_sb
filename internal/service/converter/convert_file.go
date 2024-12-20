package converter

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

type FrontMatter struct {
	Date   string   `yaml:"date"`
	Author string   `yaml:"author"`
	Tags   []string `yaml:"tags"`
	Closed bool     `yaml:"closed"`
}

func (c *Converter) ConvertFile(filePath, srcDir, destDir string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Errorf("Не удалось прочитать файл %s: %v", filePath, err)
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}
	fm, mdContent, err := c.splitFrontMatter(content)
	if err != nil {
		log.Errorf("Ошибка при разборе FrontMatter для файла %s: %v", filePath, err)
		return fmt.Errorf("ошибка при разборе FrontMatter: %v", err)
	}

	htmlContent := blackfriday.Run(mdContent)

	if len(fm.Date) > 0 || len(fm.Author) > 0 || len(fm.Tags) > 0 {
		meta := fmt.Sprintf(
			"<!-- Date: %s | Author: %s | Tags: %s | Closed: %t -->\n",
			fm.Date, fm.Author, strings.Join(fm.Tags, ", "), fm.Closed,
		)
		htmlContent = append([]byte(meta), htmlContent...)
	}

	// Определение относительного пути к файлу из исходной директории к целевой директории
	relPath, err := filepath.Rel(srcDir, filePath)
	if err != nil {
		log.Errorf("Не удалось определить относительный путь для файла %s: %v", filePath, err)
		return fmt.Errorf("не удалось определить относительный путь: %v", err)
	}

	// Замена расширения на .html с сохранением названия исходного файла
	htmlFileName := fmt.Sprintf("%s.html", filepath.Base(relPath[:len(relPath)-len(filepath.Ext(relPath))]))
	htmlFilePath := filepath.Join(destDir, htmlFileName)

	// Проверка существования HTML-файла
	if info, err := os.Stat(htmlFilePath); err == nil {
		// Получение времени последнего изменения исходного файла
		srcInfo, err := os.Stat(filePath)
		if err != nil {
			log.Errorf("Не удалось получить информацию о исходном файле %s: %v", filePath, err)
			return fmt.Errorf("не удалось получить информацию о исходном файле: %v", err)
		}

		// Сравнение времени модификации
		if !srcInfo.ModTime().After(info.ModTime()) {
			log.WithFields(log.Fields{
				"file": filePath,
			}).Info("Файл не изменился, пропуск конвертации")
			return nil
		}
	}

	// Запись HTML содержимого в файл | Создание файла если не было | Переписывание если был
	if err := os.WriteFile(htmlFilePath, htmlContent, 0644); err != nil {
		log.Errorf("Не удалось записать HTML файл %s: %v", htmlFilePath, err)
		return fmt.Errorf("не удалось записать HTML файл: %v", err)
	}

	log.WithFields(log.Fields{
		"file": htmlFilePath,
	}).Info("HTML файл успешно записан")

	return nil
}

func (c *Converter) splitFrontMatter(content []byte) (*FrontMatter, []byte, error) {
	delimiter := []byte("---")

	parts := bytes.SplitN(content, delimiter, 3)
	if len(parts) < 3 {
		return &FrontMatter{}, content, nil
	}

	var fm FrontMatter
	if err := yaml.Unmarshal(parts[1], &fm); err != nil {
		return nil, nil, fmt.Errorf("не удалось распарсить FrontMatter: %v", err)
	}

	return &fm, parts[2], nil
}
