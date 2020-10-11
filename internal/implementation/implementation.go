package implementation

import (
	"context"
	"errors"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/internal/repository/generatordir"
	"github.com/StephanHCB/go-generator-lib/internal/repository/targetdir"
	"gopkg.in/yaml.v2"
)

type GeneratorImpl struct{
}

func (i *GeneratorImpl) FindGeneratorNames(ctx context.Context, sourceBaseDir string) ([]string, error) {
	sourceDir := generatordir.Instance(ctx, sourceBaseDir)
	return sourceDir.FindGeneratorNames(ctx)
}

func (i *GeneratorImpl) ObtainGeneratorSpec(ctx context.Context, sourceBaseDir string, generatorName string) (*api.GeneratorSpec, error) {
	sourceDir := generatordir.Instance(ctx, sourceBaseDir)
	fileName := "generator-" + generatorName + ".yaml"
	generatorSpecYaml, err := sourceDir.ReadFile(ctx, fileName)
	if err != nil {
		return &api.GeneratorSpec{}, err
	}

	return i.parseGenSpec(ctx, generatorSpecYaml)
}

func (i *GeneratorImpl) WriteRenderSpecWithDefaults(ctx context.Context, request *api.Request, generatorName string) *api.Response {
	genSpec, err := i.ObtainGeneratorSpec(ctx, request.SourceBaseDir, generatorName)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	renderSpec := &api.RenderSpec{
		GeneratorName: generatorName,
		Variables: map[string]string{},
	}
	for k, v := range genSpec.Variables {
		renderSpec.Variables[k] = v.DefaultValue
	}

	renderSpecYaml, err := i.renderRenderSpec(ctx, renderSpec)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	targetDir := targetdir.Instance(ctx, request.TargetBaseDir)
	targetFile := "generated-" + generatorName + ".yaml"
	err = targetDir.WriteFile(ctx, targetFile, renderSpecYaml)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	return i.successResponse(ctx, []api.FileResult{i.successFileResult(ctx, targetFile)})
}

func (i *GeneratorImpl) Render(ctx context.Context, request *api.Request) *api.Response {
	return i.errorResponseToplevel(ctx, errors.New("Not implemented"))
}

// helper functions

func (i *GeneratorImpl) parseGenSpec(ctx context.Context, specYaml []byte) (*api.GeneratorSpec, error) {
	spec := &api.GeneratorSpec{}
	err := yaml.UnmarshalStrict(specYaml, spec)
	if err != nil {
		return &api.GeneratorSpec{}, err
	}
	return spec, nil
}

func (i *GeneratorImpl) renderRenderSpec(ctx context.Context, renderSpec *api.RenderSpec) ([]byte, error) {
	return yaml.Marshal(renderSpec)
}

func (i *GeneratorImpl) errorResponseToplevel(ctx context.Context, err error) *api.Response {
	return &api.Response{
		Errors:  []error{err},
	}
}

func (i *GeneratorImpl) successResponse(ctx context.Context, renderedFiles []api.FileResult) *api.Response {
	return &api.Response{
		Success: true,
		RenderedFiles: renderedFiles,
	}
}

func (i *GeneratorImpl) successFileResult(ctx context.Context, relativeFilePath string) api.FileResult {
	return api.FileResult{
		Success:          true,
		RelativeFilePath: relativeFilePath,
	}
}
