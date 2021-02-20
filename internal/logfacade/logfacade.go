package logfacade

import (
	"context"
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/internal/implementation"
)

type GeneratorLogfacade struct{
	Wrapped *implementation.GeneratorImpl
}

func (i *GeneratorLogfacade) FindGeneratorNames(ctx context.Context, sourceBaseDir string) ([]string, error) {
	aulogging.Logger.Ctx(ctx).Debug().Printf("entering FindGeneratorNames sourceBaseDir=%s", sourceBaseDir)
	result, err := i.Wrapped.FindGeneratorNames(ctx, sourceBaseDir)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print("error in FindGeneratorNames")
	}
	return result, err
}

func (i *GeneratorLogfacade) ObtainGeneratorSpec(ctx context.Context, sourceBaseDir string, generatorName string) (*api.GeneratorSpec, error) {
	aulogging.Logger.Ctx(ctx).Debug().Printf("entering ObtainGeneratorSpec sourceBaseDir=%s generatorName=%s", sourceBaseDir, generatorName)
	result, err := i.Wrapped.ObtainGeneratorSpec(ctx, sourceBaseDir, generatorName)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print("error in ObtainGeneratorSpec")
	}
	return result, err
}

func (i *GeneratorLogfacade) WriteRenderSpecWithDefaults(ctx context.Context, request *api.Request, generatorName string) *api.Response {
	aulogging.Logger.Ctx(ctx).Debug().Printf("entering WriteRenderSpecWithDefaults sourceBaseDir=%s targetBaseDir=%s renderSpecFile=%s generatorName=%s", request.SourceBaseDir, request.TargetBaseDir, request.RenderSpecFile, generatorName)
	result := i.Wrapped.WriteRenderSpecWithDefaults(ctx, request, generatorName)
	if len(result.Errors) > 0 {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(result.Errors[0]).Printf("%d error(s) in WriteRenderSpecWithDefaults: first error was %s", len(result.Errors), result.Errors[0].Error())
	}
	return result
}

func (i *GeneratorLogfacade) WriteRenderSpecWithValues(ctx context.Context, request *api.Request, generatorName string, parameters map[string]string) *api.Response {
	aulogging.Logger.Ctx(ctx).Debug().Printf("entering WriteRenderSpecWithValues sourceBaseDir=%s targetBaseDir=%s renderSpecFile=%s generatorName=%s", request.SourceBaseDir, request.TargetBaseDir, request.RenderSpecFile, generatorName)
	result := i.Wrapped.WriteRenderSpecWithValues(ctx, request, generatorName, parameters)
	if len(result.Errors) > 0 {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(result.Errors[0]).Printf("%d error(s) in WriteRenderSpecWithValues: first error was %s", len(result.Errors), result.Errors[0].Error())
	}
	return result
}

func (i *GeneratorLogfacade) Render(ctx context.Context, request *api.Request) *api.Response {
	aulogging.Logger.Ctx(ctx).Debug().Printf("entering Render sourceBaseDir=%s targetBaseDir=%s renderspec=%s", request.SourceBaseDir, request.TargetBaseDir, request.RenderSpecFile)
	result := i.Wrapped.Render(ctx, request)
	if len(result.Errors) > 0 || !result.Success {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(result.Errors[0]).Printf("%d top level error(s) in Render: first error was %s", len(result.Errors), result.Errors[0].Error())
		for _, f := range result.RenderedFiles {
			if len(f.Errors) > 0 || !f.Success {
				aulogging.Logger.Ctx(ctx).Warn().Printf("%s %s %d errors, first is:", "ERR", f.RelativeFilePath, len(f.Errors), f.Errors[0].Error())
			} else {
				aulogging.Logger.Ctx(ctx).Info().Printf("%s %s", "OK", f.RelativeFilePath)
			}
		}
	} else {
		aulogging.Logger.Ctx(ctx).Info().Printf("successfully rendered %d files", len(result.RenderedFiles))
		for _, f := range result.RenderedFiles {
			aulogging.Logger.Ctx(ctx).Debug().Printf("%s %s", "OK", f.RelativeFilePath)
		}
	}
	return result
}
