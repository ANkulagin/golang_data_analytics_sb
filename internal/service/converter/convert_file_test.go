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
				path = filepath.Join(os.TempDir(), "test_ok.md")
				f, _ := os.Create(path)
				_ = f.Close()

				return path, "", "", func() {
					_ = os.RemoveAll(path)
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
