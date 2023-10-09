package generatordir

import (
	"context"
	"fmt"
	"github.com/StephanHCB/go-generator-lib/api"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type GeneratorDirectory struct {
	baseDir string
}

func Instance(_ context.Context, baseDir string) *GeneratorDirectory {
	return &GeneratorDirectory{baseDir: baseDir}
}

func (d *GeneratorDirectory) CheckValid(_ context.Context) error {
	if strings.HasSuffix(d.baseDir, "/") || strings.HasSuffix(d.baseDir, "\\") {
		return fmt.Errorf("invalid generator directory: baseDir %s must not contain trailing slash", d.baseDir)
	}
	fileInfo, err := os.Stat(d.baseDir)
	if err == nil {
		// path exists, is valid, and we can access it
		if !fileInfo.IsDir() {
			return fmt.Errorf("invalid generator directory: baseDir %s must be a directory", d.baseDir)
		}
		return nil
	}
	return fmt.Errorf("invalid generator directory: baseDir %s does not exist", d.baseDir)
}

func (d *GeneratorDirectory) FindGeneratorNames(ctx context.Context) ([]string, error) {
	if err := d.CheckValid(ctx); err != nil {
		return []string{}, err
	}

	files, err := ioutil.ReadDir(d.baseDir)
	if err != nil {
		// not sure this is even reachable given we check for file stats in CheckValid
		return []string{}, fmt.Errorf("error reading generator directory: %s", err.Error())
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

func (d *GeneratorDirectory) ObtainGeneratorSpec(ctx context.Context, generatorName string) (*api.GeneratorSpec, error) {
	if err := d.CheckValid(ctx); err != nil {
		return &api.GeneratorSpec{}, err
	}

	fileName := "generator-" + generatorName + ".yaml"
	generatorSpecYaml, err := d.ReadFile(ctx, fileName)
	if err != nil {
		return &api.GeneratorSpec{}, fmt.Errorf("error reading generator spec file %s: %s", fileName, err.Error())
	}

	generatorSpec, err := d.parseGenSpec(ctx, generatorSpecYaml)
	if err != nil {
		return &api.GeneratorSpec{}, fmt.Errorf("error parsing generator spec from file %s: %s", fileName, err.Error())
	}
	return generatorSpec, nil
}

// --- public low level methods ---

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

func (d *GeneratorDirectory) Glob(ctx context.Context, relativeGlob string) ([]string, error) {
	if err := d.CheckValid(ctx); err != nil {
		return []string{}, err
	}

	filenames, err := filepath.Glob(path.Join(d.baseDir, relativeGlob))
	if err != nil {
		return []string{}, err
	}

	relativeFilenames := make([]string, 0)
	for _, fn := range filenames {
		rel, err := filepath.Rel(d.baseDir, fn)
		if err != nil {
			// unreachable as far as I'm aware
			return []string{}, fmt.Errorf("file glob %s leads to file that is not inside base directory %s - this is forbidden", relativeGlob, d.baseDir)
		}

		rel = strings.ReplaceAll(rel, `\`, `/`) // sanitize under Windows

		if strings.HasPrefix(rel, "../") {
			return []string{}, fmt.Errorf("file glob %s leads to file that is not inside base directory %s - this is forbidden", relativeGlob, d.baseDir)
		}

		relativeFilenames = append(relativeFilenames, rel)
	}

	return relativeFilenames, nil
}

// --- helper methods ---

func (d *GeneratorDirectory) parseGenSpec(_ context.Context, specYaml []byte) (*api.GeneratorSpec, error) {
	spec := &api.GeneratorSpec{}
	err := yaml.UnmarshalStrict(specYaml, spec)
	if err != nil {
		return &api.GeneratorSpec{}, err
	}
	return spec, nil
}
