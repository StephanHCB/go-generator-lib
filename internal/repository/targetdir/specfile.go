package targetdir

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type TargetDirectory struct{
	baseDir string
}

func Instance(ctx context.Context, baseDir string) *TargetDirectory {
	return &TargetDirectory{baseDir: baseDir}
}

func (d *TargetDirectory) CheckValid(ctx context.Context) error {
	if strings.HasSuffix(d.baseDir, "/") || strings.HasSuffix(d.baseDir, "\\"){
		return errors.New("baseDir must not contain trailing slash")
	}
	fileInfo, err := os.Stat(d.baseDir)
	if err == nil {
		// path exists, is valid, and we can access it
		if !fileInfo.IsDir() {
			return fmt.Errorf("baseDir must be a directory. %s is not.", d.baseDir)
		}
	}
	return err
}

func (d *TargetDirectory) ReadFile(ctx context.Context, relativePath string) ([]byte, error) {
	if err := d.CheckValid(ctx); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(path.Join(d.baseDir, relativePath))
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

func (d *TargetDirectory) WriteFile(ctx context.Context, relativePath string, contents []byte) error {
	if err := d.CheckValid(ctx); err != nil {
		return err
	}

	if err := d.createDirectoriesForFile(ctx, relativePath); err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(d.baseDir, relativePath), contents, 0644)
}

func (d *TargetDirectory) createDirectoriesForFile(ctx context.Context, relativePathForFile string) error {
	directoryPath := filepath.Dir(path.Join(d.baseDir, relativePathForFile))
	fileInfo, err := os.Stat(directoryPath)
	if err != nil {
		// ok, does not exist, create directories
		return os.MkdirAll(directoryPath, 0755)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("cannot create path up to %s, something is in the way", directoryPath)
	}
	return nil
}