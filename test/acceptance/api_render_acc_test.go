package acceptance

import (
	"context"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/docs"
	"github.com/StephanHCB/go-generator-lib/internal/repository/targetdir"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestRender_ShouldWriteExpectedFilesForDefault(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-1"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator main")
	renderspec := `generator: main
parameters:
  helloMessage: hello world
  serviceName: 'temp-service'
  serviceUrl: github.com/StephanHCB/temp
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-main.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("the return value is as expected and the correct files are written")
	expectedFilename1 := "sub/sub.go.txt"
	expectedContent1 := `package sub

import "fmt"

func PrintMessage() {
	fmt.Println("hello world")
}
`
	expectedFilename2 := "main.go.txt"
	expectedContent2 := `package src

import (
	"fmt"
	"github.com/StephanHCB/temp/sub"
)

func main() {
	fmt.Println("temp-service started")
	sub.PrintMessage()
}
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename1,
			},
			{
				Success:          true,
				RelativeFilePath: expectedFilename2,
			},
		},
	}
	require.Equal(t, expectedResponse, actualResponse)
	actual1, err := dir.ReadFile(context.TODO(), expectedFilename1)
	require.Nil(t, err)
	require.Equal(t, expectedContent1, string(actual1))
	actual2, err := dir.ReadFile(context.TODO(), expectedFilename2)
	require.Nil(t, err)
	require.Equal(t, expectedContent2, string(actual2))
}

func TestRender_ShouldComplainIfRenderSpecNotFound(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-2"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("the render spec file for generator main is missing")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("an appropriate error is returned")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	expectedErrorMsg := "open ../output/render-2/generated-main.yaml: The system cannot find the file specified."
	require.Equal(t, expectedErrorMsg, actualResponse.Errors[0].Error())
}

func TestRender_ShouldComplainIfRenderInvalid(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-3"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("an invalid render spec file")
	renderspec := `generator: something
something: weird
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-something.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-something.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("an appropriate error is returned")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	expectedErrorMsg := "yaml: unmarshal errors:\n  line 2: field something not found in type api.RenderSpec"
	require.Equal(t, expectedErrorMsg, actualResponse.Errors[0].Error())
}

func TestRender_ShouldComplainIfGenSpecNotFound(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-4"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator 'missing'")
	renderspec := `generator: missing
parameters:
  helloMessage: hello world
  serviceName: 'temp-service'
  serviceUrl: github.com/StephanHCB/temp
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-missing.yaml", []byte(renderspec)))

	docs.Given("the generator does not declare the generator name 'missing'")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-missing.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("an appropriate error is returned")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	expectedErrorMsg := "open ../resources/valid-generator-simple/generator-missing.yaml: The system cannot find the file specified."
	require.Equal(t, expectedErrorMsg, actualResponse.Errors[0].Error())
}

func TestRender_ShouldComplainIfTemplateSyntaxErrors(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-syntaxerror-templates"
	targetdirpath := "../output/render-5"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator main")
	renderspec := `generator: main
parameters:
  helloMessage: hello world
  serviceName: 'temp-service'
  serviceUrl: github.com/StephanHCB/temp
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-main.yaml", []byte(renderspec)))

	docs.Given("the generator templates contain syntax errors")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("appropriate errors are returned")
	require.False(t, actualResponse.Success)
	require.Equal(t, 2, len(actualResponse.RenderedFiles))
	require.True(t, actualResponse.RenderedFiles[0].Success)
	require.Empty(t, actualResponse.RenderedFiles[0].Errors)
	require.False(t, actualResponse.RenderedFiles[1].Success)
	require.Equal(t, "template: src_main.go.tmpl:9: bad character U+0022 '\"'", actualResponse.RenderedFiles[1].Errors[0].Error())
}
