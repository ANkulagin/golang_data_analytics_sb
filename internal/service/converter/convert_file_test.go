package converter

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestConverterFile_ReadFile(t *testing.T) {
	sut := NewConverter()
	srcDir := ""
	destDir := ""
	path := filepath.Join(os.TempDir(), "test.md")
	_, _ = os.Create(path)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(path)

	err := sut.ConvertFile(path, srcDir, destDir)

	require.NoError(t, err)
	require.FileExists(t, path)
}

func TestConverterFile_ReadFile_Error(t *testing.T) {
	sut := NewConverter()
	expected := "не удалось прочитать файл"
	srcDir := ""
	destDir := ""
	path := filepath.Join(os.TempDir(), "test.md")
	_, _ = os.Create(path)
	os.Chmod(path, 0000)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(path)

	err := sut.ConvertFile(path, srcDir, destDir)

	require.Error(t, err)
	require.Contains(t, err.Error(), expected)
}

func TestConverterFile_SplitFrontMatter(t *testing.T) {
	sut := NewConverter()
	expectedContent := `
---
date: 2024-12-09
author: ANkulagin
tags:
  - "#daily"
  - "#notes"
closed: false
---
***
## 🧾 Expenses

| Category      | Андрей | Юля |
| ------------- |:------:|:---:|
| Food          |   0    |  0  |
| Deliveries    |   0    |  0  |
| Pharmacy      |   0    |  0  |
| Entertainment |   0    |  0  |
| Gifts         |   0    |  0  |
| Wants         |   0    |  0  |
| Transport     |   0    |  0  |
| Clothing      |   0    |  0  |
| Education     |   0    |  0  |
| Home          |   0    |  0  |
| Other         |   0    |  0  |


## 🧾 Income
`
	srcDir := ""
	destDir := ""
	path := filepath.Join(os.TempDir(), "test.md")
	_, _ = os.Create(path)
	os.WriteFile(path, []byte(expectedContent), 0777)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(path)

	err := sut.ConvertFile(path, srcDir, destDir)

	require.FileExists(t, path)
}

func TestConverterFile_SplitFrontMatter_YmlIsMissing(t *testing.T) {
	sut := NewConverter()
	srcDir := ""
	destDir := ""
	path := filepath.Join(os.TempDir(), "test.md")
	_, _ = os.Create(path)

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(path)

	err := sut.ConvertFile(path, srcDir, destDir)

	require.Error(t, err)
}
