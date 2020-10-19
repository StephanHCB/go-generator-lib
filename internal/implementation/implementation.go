package implementation

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/internal/repository/generatordir"
	"github.com/StephanHCB/go-generator-lib/internal/repository/targetdir"
	"regexp"
	"strings"
	"text/template"
)

type GeneratorImpl struct {
}

func (i *GeneratorImpl) FindGeneratorNames(ctx context.Context, sourceBaseDir string) ([]string, error) {
	sourceDir := generatordir.Instance(ctx, sourceBaseDir)
	return sourceDir.FindGeneratorNames(ctx)
}

func (i *GeneratorImpl) ObtainGeneratorSpec(ctx context.Context, sourceBaseDir string, generatorName string) (*api.GeneratorSpec, error) {
	sourceDir := generatordir.Instance(ctx, sourceBaseDir)
	return sourceDir.ObtainGeneratorSpec(ctx, generatorName)
}

func (i *GeneratorImpl) WriteRenderSpecWithDefaults(ctx context.Context, request *api.Request, generatorName string) *api.Response {
	sourceDir := generatordir.Instance(ctx, request.SourceBaseDir)
	targetDir := targetdir.Instance(ctx, request.TargetBaseDir)

	genSpec, err := sourceDir.ObtainGeneratorSpec(ctx, generatorName)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	renderSpec := i.constructRenderSpecWithDefaults(ctx, generatorName, genSpec)

	targetFile := "generated-" + generatorName + ".yaml"
	err = targetDir.WriteRenderSpec(ctx, renderSpec, targetFile)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}
	return i.successResponse(ctx, []api.FileResult{i.successFileResult(ctx, targetFile)})
}

func (i *GeneratorImpl) Render(ctx context.Context, request *api.Request) *api.Response {
	sourceDir := generatordir.Instance(ctx, request.SourceBaseDir)
	targetDir := targetdir.Instance(ctx, request.TargetBaseDir)

	renderSpec, err := targetDir.ObtainRenderSpec(ctx, request.RenderSpecFile)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	genSpec, err := sourceDir.ObtainGeneratorSpec(ctx, renderSpec.GeneratorName)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	parameters, err := i.constructParameterMap(ctx, genSpec, renderSpec)
	if err != nil {
		return i.errorResponseToplevel(ctx, err)
	}

	renderedFiles, allSuccessful := i.renderAllTemplates(ctx, genSpec, parameters, sourceDir, targetDir)
	if allSuccessful {
		return i.successResponse(ctx, renderedFiles)
	} else {
		return i.errorResponseRender(ctx, renderedFiles)
	}
}

// helper functions

func (i *GeneratorImpl) constructRenderSpecWithDefaults(_ context.Context, generatorName string, genSpec *api.GeneratorSpec) *api.RenderSpec {
	renderSpec := &api.RenderSpec{
		GeneratorName: generatorName,
		Parameters:    map[string]string{},
	}
	for k, v := range genSpec.Variables {
		renderSpec.Parameters[k] = v.DefaultValue
	}
	return renderSpec
}

func (i *GeneratorImpl) constructParameterMap(_ context.Context, genSpec *api.GeneratorSpec, renderSpec *api.RenderSpec) (map[string]interface{}, error) {
	parameters := make(map[string]interface{})
	for varName, varSpec := range genSpec.Variables {
		val, ok := renderSpec.Parameters[varName]
		if !ok {
			val = varSpec.DefaultValue
		}

		if val == "" {
			return nil, fmt.Errorf("parameter %s is required but missing or empty", varName)
		}
		if varSpec.ValidationPattern != "" {
			matches, err := regexp.MatchString(varSpec.ValidationPattern, val)
			if err != nil {
				return nil, fmt.Errorf("variable declaration %s has invalid pattern: %s", varName, err.Error())
			}
			if !matches {
				return nil, fmt.Errorf("value for parameter %s does not match pattern %s", varName, varSpec.ValidationPattern)
			}
		}
		parameters[varName] = val
	}
	return parameters, nil
}

func (i *GeneratorImpl) renderAllTemplates(ctx context.Context, genSpec *api.GeneratorSpec, parameters map[string]interface{}, sourceDir *generatordir.GeneratorDirectory, targetDir *targetdir.TargetDirectory) ([]api.FileResult, bool) {
	renderedFiles := []api.FileResult{}
	allSuccessful := true
	for _, tplSpec := range genSpec.Templates {
		rendered, success := i.renderSingleTemplate(ctx, &tplSpec, parameters, sourceDir, targetDir)
		renderedFiles = append(renderedFiles, rendered...)
		allSuccessful = allSuccessful && success
	}
	return renderedFiles, allSuccessful
}

func (i *GeneratorImpl) renderSingleTemplate(ctx context.Context, tplSpec *api.TemplateSpec, parameters map[string]interface{}, sourceDir *generatordir.GeneratorDirectory, targetDir *targetdir.TargetDirectory) ([]api.FileResult, bool) {
	templateName := strings.ReplaceAll(tplSpec.RelativeSourcePath, "/", "_")
	templateContents, err := sourceDir.ReadFile(ctx, tplSpec.RelativeSourcePath)
	if err != nil {
		return []api.FileResult{i.errorFileResult(ctx, tplSpec.RelativeTargetPath, err)}, false
	}

	tmpl, err := template.New(templateName).Parse(string(templateContents))
	if err != nil {
		return []api.FileResult{i.errorFileResult(ctx, tplSpec.RelativeTargetPath, err)}, false
	}

	renderedFiles := []api.FileResult{}
	allSuccessful := true
	if len(tplSpec.WithItems) > 0 {
		for counter, item := range tplSpec.WithItems {
			parameters["item"] = item
			targetPath, err := i.renderString(ctx, parameters, fmt.Sprintf("%s_path_%d", templateName, counter), tplSpec.RelativeTargetPath)
			if err != nil {
				renderedFiles = append(renderedFiles, i.errorFileResult(ctx, targetPath, err))
				allSuccessful = false
			} else {
				err := i.renderAndWriteFile(ctx, parameters, tmpl, templateName, targetDir, targetPath)
				if err != nil {
					renderedFiles = append(renderedFiles, i.errorFileResult(ctx, targetPath, err))
					allSuccessful = false
				} else {
					renderedFiles = append(renderedFiles, i.successFileResult(ctx, targetPath))
				}
			}
		}
	} else {
		targetPath, err := i.renderString(ctx, parameters, fmt.Sprintf("%s_path", templateName), tplSpec.RelativeTargetPath)
		if err != nil {
			renderedFiles = append(renderedFiles, i.errorFileResult(ctx, targetPath, err))
			allSuccessful = false
		} else {
			err := i.renderAndWriteFile(ctx, parameters, tmpl, templateName, targetDir, targetPath)
			if err != nil {
				renderedFiles = append(renderedFiles, i.errorFileResult(ctx, targetPath, err))
				allSuccessful = false
			} else {
				renderedFiles = append(renderedFiles, i.successFileResult(ctx, targetPath))
			}
		}
	}

	return renderedFiles, allSuccessful
}

func (i *GeneratorImpl) renderAndWriteFile(ctx context.Context, parameters map[string]interface{}, tmpl *template.Template, templateName string, targetDir *targetdir.TargetDirectory, targetPath string) error {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, templateName, parameters)
	if err != nil {
		// unsure if this is reachable. All errors I've been able to produce are found during template parse
		return err
	}

	err = targetDir.WriteFile(ctx, targetPath, buf.Bytes())
	return err
}

func (i *GeneratorImpl) renderString(_ context.Context, parameters map[string]interface{}, templateName string, templateContents string) (string, error) {
	tmpl, err := template.New(templateName).Parse(templateContents)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, templateName, parameters)
	if err != nil {
		// unsure if this is reachable. All errors I've been able to produce are found during template parse
		return "", err
	}

	return buf.String(), nil
}

// --- response helpers

func (i *GeneratorImpl) errorResponseToplevel(_ context.Context, err error) *api.Response {
	return &api.Response{
		Errors: []error{err},
	}
}

func (i *GeneratorImpl) successResponse(_ context.Context, renderedFiles []api.FileResult) *api.Response {
	return &api.Response{
		Success:       true,
		RenderedFiles: renderedFiles,
	}
}

func (i *GeneratorImpl) errorResponseRender(_ context.Context, renderedFiles []api.FileResult) *api.Response {
	return &api.Response{
		Success:       false,
		RenderedFiles: renderedFiles,
		Errors:        []error{errors.New("an error occurred during rendering, see individual files")},
	}
}

func (i *GeneratorImpl) successFileResult(_ context.Context, relativeFilePath string) api.FileResult {
	return api.FileResult{
		Success:          true,
		RelativeFilePath: relativeFilePath,
	}
}

func (i *GeneratorImpl) errorFileResult(_ context.Context, relativeFilePath string, err error) api.FileResult {
	return api.FileResult{
		Success:          false,
		RelativeFilePath: relativeFilePath,
		Errors:           []error{err},
	}
}
