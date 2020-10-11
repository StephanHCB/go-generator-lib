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
	aulogging.Logger.Ctx(ctx).Debug().Printf("FindGeneratorNames sourceBaseDir=%s", sourceBaseDir)
	result, err := i.Wrapped.FindGeneratorNames(ctx, sourceBaseDir)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print("Error in FindGeneratorNames")
	}
	return result, err
}

func (i *GeneratorLogfacade) ObtainGeneratorSpec(ctx context.Context, sourceBaseDir string, generatorName string) (*api.GeneratorSpec, error) {
	aulogging.Logger.Ctx(ctx).Debug().Printf("ObtainGeneratorSpec sourceBaseDir=%s generatorName=%s", sourceBaseDir, generatorName)
	result, err := i.Wrapped.ObtainGeneratorSpec(ctx, sourceBaseDir, generatorName)
	if err != nil {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(err).Print("Error in ObtainGeneratorSpec")
	}
	return result, err
}

func (i *GeneratorLogfacade) WriteRenderSpecWithDefaults(ctx context.Context, request *api.Request, generatorName string) *api.Response {
	aulogging.Logger.Ctx(ctx).Debug().Printf("WriteRenderSpecWithDefaults sourceBaseDir=%s targetBaseDir=%s renderSpecFile=%s generatorName=%s", request.SourceBaseDir, request.TargetBaseDir, request.RenderSpecFile, generatorName)
	result := i.Wrapped.WriteRenderSpecWithDefaults(ctx, request, generatorName)
	if len(result.Errors) > 0 {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(result.Errors[0]).Printf("%n Error(s) in WriteRenderSpecWithDefaults: first error was %s", len(result.Errors), result.Errors[0].Error())
	}
	return result
}

func (i *GeneratorLogfacade) Render(ctx context.Context, request *api.Request) *api.Response {
	aulogging.Logger.Ctx(ctx).Debug().Printf("Render sourceBaseDir=%s targetBaseDir=%s renderspec=%s", request.SourceBaseDir, request.TargetBaseDir, request.RenderSpecFile)
	result := i.Wrapped.Render(ctx, request)
	if len(result.Errors) > 0 || !result.Success {
		aulogging.Logger.Ctx(ctx).Warn().WithErr(result.Errors[0]).Printf("%n top level Error(s) in Render: first error was %s", len(result.Errors), result.Errors[0].Error())
		for _, f := range result.RenderedFiles {
			if len(f.Errors) > 0 || !f.Success {
				aulogging.Logger.Ctx(ctx).Warn().Printf("%s %s %n errors, first is:", "ERR", f.RelativeFilePath, len(f.Errors), f.Errors[0].Error())
			} else {
				aulogging.Logger.Ctx(ctx).Info().Printf("%s %s", "OK", f.RelativeFilePath)
			}
		}
	} else {
		aulogging.Logger.Ctx(ctx).Info().Printf("successfully rendered %n files", len(result.RenderedFiles))
		for _, f := range result.RenderedFiles {
			aulogging.Logger.Ctx(ctx).Debug().Printf("%s %s", "OK", f.RelativeFilePath)
		}
	}
	return result
}
