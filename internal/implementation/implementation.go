package implementation

import (
	"context"
	"errors"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/internal/repository/generatordir"
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

	return i.parseSpec(ctx, generatorSpecYaml)
}

func (i *GeneratorImpl) WriteRenderSpecWithDefaults(ctx context.Context, request *api.Request, generatorName string) *api.Response {
	return &api.Response{
		Errors:  []error{errors.New("Not implemented")},
	}
}

func (i *GeneratorImpl) Render(ctx context.Context, request *api.Request) *api.Response {
	return &api.Response{
		Errors:  []error{errors.New("Not implemented")},
	}
}

// helper functions

func (i *GeneratorImpl) parseSpec(ctx context.Context, specYaml []byte) (*api.GeneratorSpec, error) {
	spec := &api.GeneratorSpec{}
	err := yaml.UnmarshalStrict(specYaml, spec)
	if err != nil {
		return &api.GeneratorSpec{}, err
	}
	return spec, nil
}
