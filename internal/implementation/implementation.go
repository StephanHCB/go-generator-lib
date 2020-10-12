package implementation

import (
	"bytes"
	"context"
	"errors"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/internal/repository/generatordir"
	"github.com/StephanHCB/go-generator-lib/internal/repository/targetdir"
	"gopkg.in/yaml.v2"
	"strings"
	"text/template"
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
		// unreachable with current feature set
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
	if request.RenderSpecFile == "" {
		aulogging.Logger.Ctx(ctx).Debug().Print("using default renderSpec generated-main.yaml")
		request.RenderSpecFile = "generated-main.yaml"
	}

	targetDir := targetdir.Instance(ctx, request.TargetBaseDir)
	renderSpecYaml, err := targetDir.ReadFile(ctx, request.RenderSpecFile)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}
	renderSpec, err := i.parseRenderSpec(ctx, renderSpecYaml)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	genSpec, err := i.ObtainGeneratorSpec(ctx, request.SourceBaseDir, renderSpec.GeneratorName)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	// TODO: validate variable values here - not part of current acceptance test

	sourceDir := generatordir.Instance(ctx, request.SourceBaseDir)
	renderedFiles := []api.FileResult{}
	allSuccessful := true
	for _, tplSpec := range genSpec.Templates {
		err := i.renderSingleTemplate(ctx, tplSpec, renderSpec, sourceDir, targetDir)
		if err != nil {
			renderedFiles = append(renderedFiles, i.errorFileResult(ctx, tplSpec.RelativeTargetPath, err))
			allSuccessful = false
		} else {
			renderedFiles = append(renderedFiles, i.successFileResult(ctx, tplSpec.RelativeTargetPath))
		}
	}

	if allSuccessful {
		return i.successResponse(ctx, renderedFiles)
	} else {
		return i.errorResponseRender(ctx, renderedFiles)
	}
}

// helper functions

func (i *GeneratorImpl) renderSingleTemplate(ctx context.Context, tplSpec api.TemplateSpec, renderSpec *api.RenderSpec, sourceDir *generatordir.GeneratorDirectory, targetDir *targetdir.TargetDirectory) error {
	templateName := strings.ReplaceAll(tplSpec.RelativeSourcePath, "/", "_")
	templateContents, err := sourceDir.ReadFile(ctx, tplSpec.RelativeSourcePath)
	if err != nil {
		return err
	}

	tmpl, err := template.New(templateName).Parse(string(templateContents))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templateName, renderSpec.Variables)
	if err != nil {
		// unsure if this is reachable. All errors I've been able to produce are found during template parse
		return err
	}

	err = targetDir.WriteFile(ctx, tplSpec.RelativeTargetPath, buf.Bytes())
	return err
}

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

func (i *GeneratorImpl) parseRenderSpec(ctx context.Context, specYaml []byte) (*api.RenderSpec, error) {
	spec := &api.RenderSpec{}
	err := yaml.UnmarshalStrict(specYaml, spec)
	if err != nil {
		return &api.RenderSpec{}, err
	}
	return spec, nil
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

func (i *GeneratorImpl) errorResponseRender(ctx context.Context, renderedFiles []api.FileResult) *api.Response {
	return &api.Response{
		Success: false,
		RenderedFiles: renderedFiles,
		Errors: []error{errors.New("an error occurred during rendering, see individual files")},
	}
}

func (i *GeneratorImpl) successFileResult(ctx context.Context, relativeFilePath string) api.FileResult {
	return api.FileResult{
		Success:          true,
		RelativeFilePath: relativeFilePath,
	}
}

func (i *GeneratorImpl) errorFileResult(ctx context.Context, relativeFilePath string, err error) api.FileResult {
	return api.FileResult{
		Success:          false,
		RelativeFilePath: relativeFilePath,
		Errors:           []error{err},
	}
}
