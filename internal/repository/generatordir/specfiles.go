package generatordir

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type GeneratorDirectory struct{
	baseDir string
}

func Instance(ctx context.Context, baseDir string) *GeneratorDirectory {
	return &GeneratorDirectory{baseDir: baseDir}
}

func (d *GeneratorDirectory) CheckValid(ctx context.Context) error {
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

func (d *GeneratorDirectory) FindGeneratorNames(ctx context.Context) ([]string, error) {
	if err := d.CheckValid(ctx); err != nil {
		return []string{}, err
	}
	files, err := ioutil.ReadDir(d.baseDir)
	if err != nil {
		return []string{}, err
	}

	regex, _ := regexp.Compile("^generator-(.*).yaml$")
	result := []string{}
	for _, f := range files {
		if f.Mode().IsRegular() {
			if matchInfo := regex.FindStringSubmatch(f.Name()); matchInfo != nil {
				result = append(result, matchInfo[1])
			}
		}
	}

	sort.Strings(result)

	return result, nil
}

func (d *GeneratorDirectory) ReadFile(ctx context.Context, relativePath string) ([]byte, error) {
	if err := d.CheckValid(ctx); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(path.Join(d.baseDir, relativePath))
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}
