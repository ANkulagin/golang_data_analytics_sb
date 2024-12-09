package converter

import (
	"bytes"
	"fmt"
	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
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
		return fmt.Errorf("не удалось прочитать файл: %v", err)
	}
	fm, mdContent, err := c.splitFrontMatter(content)
	if err != nil {
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
		return fmt.Errorf("не удалось определить относительный путь: %v", err)
	}

	// Замена расширения на .html с сохранением названия исходного файла
	htmlFileName := fmt.Sprintf("%s.html", filepath.Base(relPath[:len(relPath)-len(filepath.Ext(relPath))]))
	htmlFilePath := filepath.Join(destDir, htmlFileName)

	// Запись HTML содержимого в файл | Создания файла если не было |Переписывание если был
	if err := os.WriteFile(htmlFilePath, htmlContent, 0644); err != nil {
		return fmt.Errorf("не удалось записать HTML файл: %v", err)
	}

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
