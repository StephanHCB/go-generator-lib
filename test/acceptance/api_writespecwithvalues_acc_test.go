package acceptance

import (
	"context"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/docs"
	"github.com/StephanHCB/go-generator-lib/internal/repository/targetdir"
	"github.com/stretchr/testify/require"
	"os"
	"path"
	"testing"
)

// --- happy cases

func TestWriteRenderSpecWithValues_ShouldCreateMainSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-1"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "main"

	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-main.yaml",
	}
	docs.When("WriteRenderSpecWithValues is invoked with valid parameters")
	parameters := map[string]interface{}{
		"helloMessage": "hello nice world",
		"serviceName": "something-valid",
		"serviceUrl": "github.com/StephanHCB/scratch",
	}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the correct spec file is written and the return value is as expected")
	expectedFilename := "generated-main.yaml"
	expectedContent := `generator: main
parameters:
  helloMessage: hello nice world
  serviceName: something-valid
  serviceUrl: github.com/StephanHCB/scratch
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename,
			},
		},
	}
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	actual, err := dir.ReadFile(context.TODO(), expectedFilename)
	require.Nil(t, err)
	require.Equal(t, expectedContent, string(actual))
	require.Equal(t, expectedResponse, actualResponse)
}

func TestWriteRenderSpecWithValues_ShouldCreateNonstandardSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-2"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "main"

	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generator-values.dat",
	}
	docs.When("WriteRenderSpecWithValues is invoked with valid parameters and a nonstandard render spec filename")
	parameters := map[string]interface{}{
		"helloMessage": "hello nice world",
		"serviceName": "something-valid",
		"serviceUrl": "github.com/StephanHCB/scratch",
	}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the correct spec file is written and the return value is as expected")
	expectedFilename := "generator-values.dat"
	expectedContent := `generator: main
parameters:
  helloMessage: hello nice world
  serviceName: something-valid
  serviceUrl: github.com/StephanHCB/scratch
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename,
			},
		},
	}
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	actual, err := dir.ReadFile(context.TODO(), expectedFilename)
	require.Nil(t, err)
	require.Equal(t, expectedContent, string(actual))
	require.Equal(t, expectedResponse, actualResponse)
}

func TestWriteRenderSpecWithValues_ShouldCreateDefaultSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-3"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "main"

	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	docs.When("WriteRenderSpecWithValues is invoked with valid parameters, but empty render spec filename")
	parameters := map[string]interface{}{
		"helloMessage": "hello nice world",
		"serviceName": "something-valid",
		"serviceUrl": "github.com/StephanHCB/scratch",
	}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the correct spec file is written and the return value is as expected")
	expectedFilename := "generated-main.yaml"
	expectedContent := `generator: main
parameters:
  helloMessage: hello nice world
  serviceName: something-valid
  serviceUrl: github.com/StephanHCB/scratch
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename,
			},
		},
	}
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	actual, err := dir.ReadFile(context.TODO(), expectedFilename)
	require.Nil(t, err)
	require.Equal(t, expectedContent, string(actual))
	require.Equal(t, expectedResponse, actualResponse)
}

func TestWriteRenderSpecWithValues_ShouldOverwriteSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-4"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "docker"

	docs.Given("a render spec for this generator already exists")
	expectedFilename := "generated-docker.yaml"
	originalContent := `generator: docker
parameters:
  somethingElse: "contents from the first write"
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), expectedFilename, []byte(originalContent)))

	docs.When("WriteRenderSpecWithValues is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-docker.yaml",
	}
	parameters := map[string]interface{}{
		"serviceName": "docker-is-great",
	}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the render spec file is silently overwritten and the return value is as expected")
	expectedContent := `generator: docker
parameters:
  serviceName: docker-is-great
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename,
			},
		},
	}
	actual, err := dir.ReadFile(context.TODO(), expectedFilename)
	require.Nil(t, err)
	require.Equal(t, expectedContent, string(actual))
	require.Equal(t, expectedResponse, actualResponse)
}

// TODO test with a structured default value

// --- error cases

func TestWriteRenderSpecWithValues_ShouldComplainMissingSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-5"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a generator name for which no spec exists")
	name := "notpresent"

	docs.When("WriteRenderSpecWithValues is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	parameters := map[string]interface{}{}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the response reports an appropriate error")
	expectedErrorMessagePart := "error reading generator spec file generator-notpresent.yaml: open ../resources/valid-generator-simple/generator-notpresent.yaml: "
	require.False(t, actualResponse.Success)
	require.Contains(t, actualResponse.Errors[0].Error(), expectedErrorMessagePart)
}

func TestWriteRenderSpecWithValues_ShouldComplainTargetExistsIsDir(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-6"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("the target filename is taken by a directory")
	require.Nil(t, os.Mkdir(path.Join(targetdirpath, "generated-docker.yaml"), 0755))

	docs.Given("a valid generator name")
	name := "docker"

	docs.When("WriteRenderSpecWithValues is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	parameters := map[string]interface{}{
		"serviceName": "docker-is-great",
	}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the response reports an appropriate error")
	expectedErrorMessage := "error writing render spec file generated-docker.yaml in target dir ../output/write-render-spec-values-6: open ../output/write-render-spec-values-6/generated-docker.yaml: is a directory"
	require.False(t, actualResponse.Success)
	require.Equal(t, expectedErrorMessage, actualResponse.Errors[0].Error())
}

func TestWriteRenderSpecWithValues_ShouldComplainMissingParameter(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-7"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "docker"

	docs.When("WriteRenderSpecWithValues is invoked with a missing required parameter")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	parameters := map[string]interface{}{}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the response reports an appropriate error")
	expectedErrorMessage := "parameter 'serviceName' is required but missing or empty"
	require.False(t, actualResponse.Success)
	require.Equal(t, expectedErrorMessage, actualResponse.Errors[0].Error())
}

func TestWriteRenderSpecWithValues_ShouldComplainUnknownParameter(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-values-7"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "docker"

	docs.When("WriteRenderSpecWithValues is invoked with an unexpected additional parameter (not in the spec)")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	parameters := map[string]interface{}{
		"serviceName":         "docker-is-great",
		"somethingUnexpected": "huh, what's this",
	}
	actualResponse := generatorlib.WriteRenderSpecWithValues(context.TODO(), request, name, parameters)

	docs.Then("the response reports an appropriate error")
	expectedErrorMessage := "parameter 'somethingUnexpected' is not allowed according to generator spec"
	require.False(t, actualResponse.Success)
	require.Equal(t, expectedErrorMessage, actualResponse.Errors[0].Error())
}
