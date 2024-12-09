package converter

import (
	"bytes"
	"fmt"
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
	_, _, _ = c.splitFrontMatter(content)

	return nil
}

func (c *Converter) splitFrontMatter(content []byte) (*FrontMatter, []byte, error) {
	delimiter := []byte("---")

	parts := bytes.SplitN(content, delimiter, 3)
	if len(parts) != 3 {
		return &FrontMatter{}, content, nil
	}

	return &FrontMatter{}, nil, nil
}
