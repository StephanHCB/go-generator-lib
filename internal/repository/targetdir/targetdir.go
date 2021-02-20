package targetdir

import (
	"context"
	"fmt"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-generator-lib/api"
	"gopkg.in/yaml.v2"
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
		return fmt.Errorf("error invalid target directory: baseDir %s must not contain trailing slash", d.baseDir)
	}
	fileInfo, err := os.Stat(d.baseDir)
	if err == nil {
		// path exists, is valid, and we can access it
		if !fileInfo.IsDir() {
			return fmt.Errorf("error invalid target directory: baseDir %s must be a directory", d.baseDir)
		}
		return nil
	}
	return fmt.Errorf("error invalid target directory: baseDir %s does not exist", d.baseDir)
}

func (d *TargetDirectory) ObtainRenderSpec(ctx context.Context, renderSpecFilenameOrEmptyString string) (*api.RenderSpec, error) {
	specFile := d.RenderSpecFilenameOrDefault(ctx, renderSpecFilenameOrEmptyString)

	renderSpecYaml, err := d.ReadFile(ctx, specFile)
	if err != nil {
		return &api.RenderSpec{}, fmt.Errorf("error reading render spec file %s in target directory %s: %s", specFile, d.baseDir, err.Error())
	}
	renderSpec, err := d.parseRenderSpec(ctx, renderSpecYaml)
	if err != nil {
		return renderSpec, fmt.Errorf("error parsing render spec file %s in target directory %s: %s", specFile, d.baseDir, err.Error())
	}
	return renderSpec, nil
}

func (d *TargetDirectory) WriteRenderSpec(ctx context.Context, renderSpec *api.RenderSpec, renderSpecFilenameOrEmptyString string) (string, error) {
	targetFile := d.RenderSpecFilenameOrDefaultForGenerator(ctx, renderSpecFilenameOrEmptyString, renderSpec.GeneratorName)

	renderSpecYaml, err := d.renderRenderSpec(ctx, renderSpec)
	if err != nil {
		// unreachable with current feature set as far as I'm aware
		return targetFile, fmt.Errorf("error preparing render spec: %s", err.Error())
	}

	err = d.WriteFile(ctx, targetFile, renderSpecYaml)
	if err != nil {
		return targetFile, fmt.Errorf("error writing render spec file %s in target dir %s: %s", targetFile, d.baseDir, err.Error())
	}
	return targetFile, nil
}

// --- low level methods, public so they can be used in tests ---

func (d *TargetDirectory) RenderSpecFilenameOrDefault(ctx context.Context, renderSpecFilename string) string {
	return d.RenderSpecFilenameOrDefaultForGenerator(ctx, renderSpecFilename, "main")
}

func (d *TargetDirectory) RenderSpecFilenameOrDefaultForGenerator(ctx context.Context, renderSpecFilename string, generatorName string) string {
	if renderSpecFilename == "" {
		result := "generated-" + generatorName + ".yaml"
		aulogging.Logger.Ctx(ctx).Debug().Printf("using default renderSpec %s", result)
		return result
	}
	return renderSpecFilename
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
		err2 := os.MkdirAll(directoryPath, 0755)
		if err2 != nil {
			return fmt.Errorf("cannot create path up to %s, something is in the way or invalid path: %s", strings.ReplaceAll(directoryPath, "\\", "/"), err2.Error())
		} else {
			return nil
		}
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("cannot create path up to %s, something is in the way", strings.ReplaceAll(directoryPath, "\\", "/"))
	}
	return nil
}

// --- helper methods ---

func (d *TargetDirectory) parseRenderSpec(ctx context.Context, specYaml []byte) (*api.RenderSpec, error) {
	spec := &api.RenderSpec{}
	err := yaml.UnmarshalStrict(specYaml, spec)
	if err != nil {
		return &api.RenderSpec{}, err
	}
	return spec, nil
}

func (d *TargetDirectory) renderRenderSpec(ctx context.Context, renderSpec *api.RenderSpec) ([]byte, error) {
	return yaml.Marshal(renderSpec)
}

