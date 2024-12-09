package converter

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestConvertFile_Success(t *testing.T) {
	sut := NewConverter()

	testCases := []struct {
		name  string
		setup func() (path, srcDir, destDir string, cleanup func())
	}{
		{
			name: "Успешное чтение файла без FrontMatter",
			setup: func() (path, srcDir, destDir string, cleanup func()) {
				srcDir, _ = os.MkdirTemp("", "src_dir")
				destDir, _ = os.MkdirTemp("", "dest_dir")
				path = filepath.Join(srcDir, "test_ok.md")
				f, _ := os.Create(path)
				_ = f.Close()

				return path, srcDir, destDir, func() {
					_ = os.RemoveAll(srcDir)
					_ = os.RemoveAll(destDir)
				}
			},
		},
		{
			name: "Успешное чтение файла с корректным FrontMatter",
			setup: func() (path, srcDir, destDir string, cleanup func()) {
				srcDir, _ = os.MkdirTemp("", "src_dir")
				destDir, _ = os.MkdirTemp("", "dest_dir")
				path = filepath.Join(srcDir, "test_ok_fm.md")
				content := `---
date: 2024-12-09
author: "ANkulagin"
tags:
  - "#daily"
  - "#notes"
closed: false
---
# Заголовок

Контент...
`
				err := os.WriteFile(path, []byte(content), 0644)
				require.NoError(t, err)

				return path, srcDir, destDir, func() {
					_ = os.RemoveAll(srcDir)
					_ = os.RemoveAll(destDir)
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path, srcDir, destDir, cleanup := tc.setup()
			defer cleanup()

			err := sut.ConvertFile(path, srcDir, destDir)

			require.NoError(t, err)
		})
	}
}

func TestConvertFile_Error(t *testing.T) {
	sut := NewConverter()

	testCases := []struct {
		name           string
		setup          func() (path, srcDir, destDir string, cleanup func())
		expectedErrMsg string
	}{
		{
			name: "The file does not exist",
			setup: func() (path, srcDir, destDir string, cleanup func()) {
				path = filepath.Join(os.TempDir(), "test.md")
				//Не создаём файл, чтобы гарантировать ошибку чтения
				return path, "", "", func() {
					_ = os.RemoveAll(path)
				}
			},
			expectedErrMsg: "не удалось прочитать файл",
		},
		{
			name: "No read permissions",
			setup: func() (path, srcDir, destDir string, cleanup func()) {
				path = filepath.Join(os.TempDir(), "test.md")
				f, _ := os.Create(path)
				_ = f.Close()
				_ = os.Chmod(path, 0000)
				return path, "", "", func() {
					// Советуют восстанавливать права перед удалением
					_ = os.Chmod(path, 0644)
					_ = os.RemoveAll(path)
				}
			},
			expectedErrMsg: "не удалось прочитать файл",
		},
		{
			name: "Invalid FrontMatter format",
			setup: func() (path, srcDir, destDir string, cleanup func()) {
				path = filepath.Join(os.TempDir(), "test_bad_fm.md")
				content := `---
date: 2024-12-09
author: ANkulagin
tags:
  - "#daily"
  - "#notes
closed: false
---
# Заголовок
Контент
`
				err := os.WriteFile(path, []byte(content), 0644)
				require.NoError(t, err)

				return path, "", "", func() {
					_ = os.RemoveAll(path)
				}
			},
			expectedErrMsg: "ошибка при разборе FrontMatter: не удалось распарсить FrontMatter:",
		},
		{
			name: "Ошибка определения относительного пути (srcDir не является родителем filePath)",
			setup: func() (string, string, string, func()) {
				// Создадим файл вне srcDir
				srcDir, _ := os.MkdirTemp("", "/a/a/src_dir")
				fPath := filepath.Join(os.TempDir(), "outside.md")
				require.NoError(t, os.WriteFile(fPath, []byte("# Outside"), 0644))

				cleanup := func() {
					_ = os.RemoveAll(srcDir)
					_ = os.RemoveAll(fPath)
				}
				return fPath, srcDir, os.TempDir(), cleanup
			},
			expectedErrMsg: "не удалось определить относительный путь",
		},
		{
			name: "Ошибка записи итогового файла (destDir не существует)",
			setup: func() (string, string, string, func()) {
				srcDir, _ := os.MkdirTemp("", "src_dir")
				fPath := filepath.Join(srcDir, "test.md")
				require.NoError(t, os.WriteFile(fPath, []byte("# Content"), 0644))

				// destDir не создаём специально
				destDir := filepath.Join(os.TempDir(), "non_existing_dest")

				cleanup := func() {
					_ = os.RemoveAll(srcDir)
					// destDir не существует, значит ничего не удаляем
				}
				return fPath, srcDir, destDir, cleanup
			},
			expectedErrMsg: "не удалось записать HTML файл",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			path, srcDir, destDir, cleanup := tc.setup()
			defer cleanup()

			err := sut.ConvertFile(path, srcDir, destDir)

			require.Error(t, err)
			require.Contains(t, err.Error(), tc.expectedErrMsg)

		})
	}
}

func TestConvertFile_Success_FrontMatterVariants(t *testing.T) {
	sut := NewConverter()

	testCases := []struct {
		name         string
		content      string
		checkContent func(t *testing.T, htmlStr string)
	}{
		{
			name: "FrontMatter только с датой",
			content: `---
date: 2024-12-10
---
# Title
Some content
`,
			checkContent: func(t *testing.T, htmlStr string) {
				require.Contains(t, htmlStr, "Date: 2024-12-10")
			},
		},
		{
			name: "FrontMatter только с автором",
			content: `---
author: ANkulagin
---
# Title
`,
			checkContent: func(t *testing.T, htmlStr string) {
				require.Contains(t, htmlStr, "Author: ANkulagin")
			},
		},
		{
			name: "FrontMatter только с тегами",
			content: `---
tags:
  - "#go"
  - "#test"
---
# Title
`,
			checkContent: func(t *testing.T, htmlStr string) {
				require.Contains(t, htmlStr, "Tags: #go, #test")
			},
		},
		{
			name: "FrontMatter с пустым значением",
			content: `---
date:
---
# Title
`,
			checkContent: func(t *testing.T, htmlStr string) {
				require.NotContains(t, htmlStr, "Date:")
				require.NotContains(t, htmlStr, "Author:")
				require.NotContains(t, htmlStr, "Tags:")
				require.NotContains(t, htmlStr, "Closed:")
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			srcDir, err := os.MkdirTemp("", "src_dir")
			require.NoError(t, err)
			destDir, err := os.MkdirTemp("", "dest_dir")
			require.NoError(t, err)

			path := filepath.Join(srcDir, "test.md")
			require.NoError(t, os.WriteFile(path, []byte(tc.content), 0644))

			defer func() {
				_ = os.RemoveAll(srcDir)
				_ = os.RemoveAll(destDir)
			}()

			err = sut.ConvertFile(path, srcDir, destDir)
			require.NoError(t, err)

			htmlPath := filepath.Join(destDir, "test.html")
			require.FileExists(t, htmlPath)

			htmlData, err := os.ReadFile(htmlPath)
			require.NoError(t, err)

			tc.checkContent(t, string(htmlData))
		})
	}
}
