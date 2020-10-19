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

func TestWriteRenderSpecWithDefaults_ShouldCreateMainSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-1"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid generator name")
	name := "main"

	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	docs.When("WriteRenderSpecWithDefaults is invoked")
	actualResponse := generatorlib.WriteRenderSpecWithDefaults(context.TODO(), request, name)

	docs.Then("the correct spec file is written and the return value is as expected")
	expectedFilename := "generated-main.yaml"
	expectedContent := `generator: main
parameters:
  helloMessage: hello world
  serviceName: ""
  serviceUrl: github.com/StephanHCB/temp
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

func TestWriteRenderSpecWithDefaults_ShouldOverwriteMainSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-2"
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

	docs.When("WriteRenderSpecWithDefaults is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	actualResponse := generatorlib.WriteRenderSpecWithDefaults(context.TODO(), request, name)

	docs.Then("the render spec file is silently overwritten and the return value is as expected")
	expectedContent := `generator: docker
parameters:
  serviceName: ""
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

func TestWriteRenderSpecWithDefaults_ShouldComplainMissingSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-3"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a generator name for which no spec exists")
	name := "notpresent"

	docs.When("WriteRenderSpecWithDefaults is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	actualResponse := generatorlib.WriteRenderSpecWithDefaults(context.TODO(), request, name)

	docs.Then("the response reports an appropriate error")
	expectedErrorMessage := "error reading generator spec file generator-notpresent.yaml: open ../resources/valid-generator-simple/generator-notpresent.yaml: The system cannot find the file specified."
	require.False(t, actualResponse.Success)
	require.Equal(t, expectedErrorMessage, actualResponse.Errors[0].Error())
}

func TestWriteRenderSpecWithDefaults_ShouldComplainTargetExistsIsDir(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/write-render-spec-4"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("the target filename is taken by a directory")
	require.Nil(t, os.Mkdir(path.Join(targetdirpath, "generated-docker.yaml"), 0755))

	docs.Given("a valid generator name")
	name := "docker"

	docs.When("WriteRenderSpecWithDefaults is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	actualResponse := generatorlib.WriteRenderSpecWithDefaults(context.TODO(), request, name)

	docs.Then("the response reports an appropriate error")
	expectedErrorMessage := "error writing render spec file generated-docker.yaml in target dir ../output/write-render-spec-4: open ../output/write-render-spec-4/generated-docker.yaml: is a directory"
	require.False(t, actualResponse.Success)
	require.Equal(t, expectedErrorMessage, actualResponse.Errors[0].Error())
}
