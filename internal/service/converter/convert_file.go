package converter

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
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
	_, _, err = c.splitFrontMatter(content)
	if err != nil {
		return fmt.Errorf("ошибка при разборе FrontMatter: %v", err)
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
