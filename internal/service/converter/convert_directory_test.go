package converter

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestConverterDirectory_SourceDirDoesNotExist(t *testing.T) {
	sut := NewConverter()
	expectedErrMsg := "исходная директория не существует"

	nonExistentDir := filepath.Join(os.TempDir(), "src_dir")

	destDir := filepath.Join(os.TempDir(), "dest_dir")
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(destDir)

	err := sut.ConvertDirectory(nonExistentDir, destDir)

	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErrMsg)
}

func TestConverterDirectory_CreateDestinationDir(t *testing.T) {
	sut := NewConverter()

	srcDir, _ := os.MkdirTemp("", "src_dir")
	defer os.RemoveAll(srcDir)

	destDir := filepath.Join(os.TempDir(), "dest_dir")
	defer os.RemoveAll(destDir)

	err := sut.ConvertDirectory(srcDir, destDir)

	require.NoError(t, err)
	require.DirExists(t, destDir)

}

func TestConverterDirectory_CreateDestinationDir_Failure(t *testing.T) {
	sut := NewConverter()
	expectedErrMsg := "не удалось создать целевую директорию"

	srcDir, _ := os.MkdirTemp("", "src_dir")
	defer os.RemoveAll(srcDir)

	parentDir, _ := os.MkdirTemp("", "parrent_dir")
	defer func(path string) {
		_ = os.Chmod(path, 0755)
		_ = os.RemoveAll(path)
	}(parentDir)

	_ = os.Chmod(parentDir, 0555)

	destDir := filepath.Join(parentDir, "dest_dir")

	err := sut.ConvertDirectory(srcDir, destDir)

	require.Error(t, err)
	require.Contains(t, err.Error(), expectedErrMsg)
}

func TestConverterDirectory_WithMdFile(t *testing.T) {
	sut := NewConverter()

	srcDir, err := os.MkdirTemp("", "src_dir")
	require.NoError(t, err)
	defer os.RemoveAll(srcDir)

	destDir := filepath.Join(os.TempDir(), "dest_dir")
	defer os.RemoveAll(destDir)

	mdFile := filepath.Join(srcDir, "test.md")
	err = os.WriteFile(mdFile, []byte("# Hello World"), 0644)
	require.NoError(t, err)

	err = sut.ConvertDirectory(srcDir, destDir)
	require.NoError(t, err)
}

func TestConverterDirectory_WithNonMdFiles(t *testing.T) {
	sut := NewConverter()

	srcDir, err := os.MkdirTemp("", "src_dir")
	require.NoError(t, err)
	defer os.RemoveAll(srcDir)

	destDir := filepath.Join(os.TempDir(), "dest_dir")
	defer os.RemoveAll(destDir)

	txtFile := filepath.Join(srcDir, "not_markdown.txt")
	err = os.WriteFile(txtFile, []byte("Just text"), 0644)
	require.NoError(t, err)

	subDir := filepath.Join(srcDir, "subfolder")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	err = sut.ConvertDirectory(srcDir, destDir)
	require.NoError(t, err)
}

func TestConverterDirectory_WalkError(t *testing.T) {
	sut := NewConverter()

	srcDir, err := os.MkdirTemp("", "src_dir")
	require.NoError(t, err)
	defer os.RemoveAll(srcDir)

	// Создаём поддиректорию
	subDir := filepath.Join(srcDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	// Забираем права на чтение поддиректории, чтобы Walk не смог её прочитать
	err = os.Chmod(subDir, 0)
	require.NoError(t, err)
	defer os.Chmod(subDir, 0755) // Восстанавливаем права после теста

	destDir := filepath.Join(os.TempDir(), "dest_dir")
	defer os.RemoveAll(destDir)

	err = sut.ConvertDirectory(srcDir, destDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "permission denied")
}

func TestConverterDirectory_ConvertFileError(t *testing.T) {
	// TODO После реалзиации конвертера
}
