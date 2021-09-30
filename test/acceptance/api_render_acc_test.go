package acceptance

import (
	"context"
	"fmt"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	"github.com/StephanHCB/go-generator-lib/api"
	"github.com/StephanHCB/go-generator-lib/docs"
	"github.com/StephanHCB/go-generator-lib/internal/repository/targetdir"
	"github.com/stretchr/testify/require"
	"os"
	"strings"
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
	fmt.Println("HELLO WORLD")
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
	require.Equal(t, toUnix(expectedContent1), toUnix(string(actual1)))
	actual2, err := dir.ReadFile(context.TODO(), expectedFilename2)
	require.Nil(t, err)
	require.Equal(t, toUnix(expectedContent2), toUnix(string(actual2)))
}

func TestRender_ShouldWriteExpectedFilesForStructured(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-structured"
	targetdirpath := "../output/render-1a"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator main")
	renderspec := `generator: main
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
	expectedFilename1 := "main.txt"
	expectedContent1 := `imagine a European wildcat

then look up something in a list: two

then look up something in a structure in a list: [sub 1 sub 2]
(value is itself a list)
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename1,
			},
		},
	}
	require.Equal(t, expectedResponse, actualResponse)
	actual1, err := dir.ReadFile(context.TODO(), expectedFilename1)
	actual1normalized := strings.ReplaceAll(string(actual1), "\r", "")
	require.Nil(t, err)
	require.Equal(t, expectedContent1, actual1normalized)
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
	expectedErrorMsgPart := "error reading render spec file generated-main.yaml in target directory ../output/render-2: open ../output/render-2/generated-main.yaml: "
	require.Contains(t, actualResponse.Errors[0].Error(), expectedErrorMsgPart)
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
	expectedErrorMsg := "error parsing render spec file generated-something.yaml in target directory ../output/render-3: yaml: unmarshal errors:\n  line 2: field something not found in type api.RenderSpec"
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
	expectedErrorMsgPart := "error reading generator spec file generator-missing.yaml: open ../resources/valid-generator-simple/generator-missing.yaml: "
	require.Contains(t, actualResponse.Errors[0].Error(), expectedErrorMsgPart)
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
	require.Equal(t, 3, len(actualResponse.RenderedFiles))
	require.True(t, actualResponse.RenderedFiles[0].Success)
	require.Empty(t, actualResponse.RenderedFiles[0].Errors)
	require.False(t, actualResponse.RenderedFiles[1].Success)
	require.Equal(t, "failed to parse template src/main.go.tmpl: template: src_main.go.tmpl:9: bad character U+0022 '\"'", actualResponse.RenderedFiles[1].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[2].Success)
	require.Contains(t, actualResponse.RenderedFiles[2].Errors[0].Error(), "failed to load template src/notfound.go.tmpl: open ../resources/valid-generator-syntaxerror-templates/src/notfound.go.tmpl: ")
}

func TestRender_ShouldComplainIfVariableValuesInvalid(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-6"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator main with an invalid variable value")
	renderspec := `generator: main
parameters:
  helloMessage: hello world
  serviceName: 'invalid service name'
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

	docs.Then("appropriate validation errors are returned")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	require.Equal(t, "value for parameter 'serviceName' does not match pattern ^[a-z-]+$", actualResponse.Errors[0].Error())
}

func TestRender_ShouldWriteExpectedFilesForItemized(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-7"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator items, which uses with_items and template directives")
	renderspec := `generator: items
parameters: {}
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-items.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-items.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("the return value is as expected and the correct files are written")
	expectedFilename1 := "first.txt"
	expectedContent1 := "Hi Frank!\n"
	expectedFilename2 := "second.txt"
	expectedContent2 := "Hi John!\n"
	expectedFilename3 := "third.txt"
	expectedContent3 := "Hi Eve!\n"
	expectedFilename4 := "fourth.txt"
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
			{
				Success:          true,
				RelativeFilePath: expectedFilename3,
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
	actual3, err := dir.ReadFile(context.TODO(), expectedFilename3)
	require.Nil(t, err)
	require.Equal(t, expectedContent3, string(actual3))
	// fourth file has a condition and should have been skipped
	_, err = dir.ReadFile(context.TODO(), expectedFilename4)
	require.NotNil(t, err)
}

func TestRender_ShouldComplainIfRequiredParameterValueMissing(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-8"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator main with a missing required parameter value")
	renderspec := `generator: main
parameters:
  helloMessage: hello world
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

	docs.Then("appropriate validation errors are returned")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	require.Equal(t, "parameter 'serviceName' is required but missing", actualResponse.Errors[0].Error())
}

func TestRender_ShouldComplainIfInvalidPattern(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/invalid-generator-specs"
	targetdirpath := "../output/render-9"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator variablepattern")
	renderspec := `generator: variablepattern
parameters:
  helloMessage: hello world
  serviceName: something
  serviceUrl: github.com/StephanHCB/temp
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-variablepattern.yaml", []byte(renderspec)))

	docs.Given("the generator spec contains an invalid regex pattern for one of the variables")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-variablepattern.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("appropriate validation errors are returned")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	require.Equal(t, "variable declaration serviceName has invalid pattern (this is an error in the generator spec, not the render request): error parsing regexp: missing closing ]: `[a-z-+$`", actualResponse.Errors[0].Error())
}

func TestRender_ShouldComplainIfTemplateSyntaxErrorsInGenSpec(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/invalid-generator-specs"
	targetdirpath := "../output/render-10"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator itemssyntax")
	renderspec := `generator: itemssyntax
parameters: {}
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-itemssyntax.yaml", []byte(renderspec)))

	docs.Given("the generator spec contains a template language expression in the target: field with syntax errors")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-itemssyntax.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("appropriate errors are returned")
	require.False(t, actualResponse.Success)
	require.Equal(t, 6, len(actualResponse.RenderedFiles))
	require.False(t, actualResponse.RenderedFiles[0].Success)
	require.Equal(t, "error evaluating target path from '{{ .item.file .txt' for item #1: template: item.txt.tmpl_path_1:1: unclosed action", actualResponse.RenderedFiles[0].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[1].Success)
	require.Equal(t, "error evaluating target path from '{{ .item.file .txt' for item #2: template: item.txt.tmpl_path_2:1: unclosed action", actualResponse.RenderedFiles[1].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[2].Success)
	require.Equal(t, "error evaluating target path from '{{ .item.file .txt' for item #3: template: item.txt.tmpl_path_3:1: unclosed action", actualResponse.RenderedFiles[2].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[3].Success)
	require.Equal(t, "error evaluating target path from '{{ .something .txt': template: item.txt.tmpl_path:1: unclosed action", actualResponse.RenderedFiles[3].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[4].Success)
	require.Equal(t, "error evaluating condition from '{{ .something .txt': template: item.txt.tmpl_condition:1: unclosed action", actualResponse.RenderedFiles[4].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[5].Success)
	require.Equal(t, "error evaluating condition from '{{ .item.file ' for item #1: template: item.txt.tmpl_condition_1:1: unclosed action", actualResponse.RenderedFiles[5].Errors[0].Error())
}

func TestRender_ShouldComplainIfSyntaxErrorsInTemplateWithItems(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/invalid-generator-specs"
	targetdirpath := "../output/render-11"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator items")
	renderspec := `generator: items
parameters: {}
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-items.yaml", []byte(renderspec)))

	docs.Given("a template used with with_items contains syntax errors")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-items.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("appropriate errors are returned")
	require.False(t, actualResponse.Success)
	require.Equal(t, 4, len(actualResponse.RenderedFiles))
	require.False(t, actualResponse.RenderedFiles[0].Success)
	require.Equal(t, "failed to parse template itemerror.txt.tmpl: template: itemerror.txt.tmpl:1: unexpected \"!\" in operand", actualResponse.RenderedFiles[0].Errors[0].Error())
	require.True(t, actualResponse.RenderedFiles[1].Success)
	require.Empty(t, actualResponse.RenderedFiles[1].Errors)
	require.True(t, actualResponse.RenderedFiles[2].Success)
	require.Empty(t, actualResponse.RenderedFiles[2].Errors)
	require.True(t, actualResponse.RenderedFiles[3].Success)
	require.Empty(t, actualResponse.RenderedFiles[3].Errors)
}

func TestRender_ShouldComplainIfInvalidTargetFiles(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/invalid-generator-specs"
	targetdirpath := "../output/render-12"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator invalidtargets")
	renderspec := `generator: invalidtargets
parameters: {}
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-invalidtargets.yaml", []byte(renderspec)))

	docs.Given("the generator spec contains expressions that result in empty target file names")

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-invalidtargets.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("appropriate errors are returned")
	require.False(t, actualResponse.Success)
	require.Equal(t, 2, len(actualResponse.RenderedFiles))
	require.False(t, actualResponse.RenderedFiles[0].Success)
	require.Equal(t, "error evaluating template for target '' for item #1: open ../output/render-12: is a directory", actualResponse.RenderedFiles[0].Errors[0].Error())
	require.False(t, actualResponse.RenderedFiles[1].Success)
	require.Equal(t, "error evaluating template for target '': open ../output/render-12: is a directory", actualResponse.RenderedFiles[1].Errors[0].Error())
}

func TestRender_ShouldWriteExpectedFilesForTemplatedDefault(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := "../output/render-13"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator templatevars")
	renderspec := `generator: templatevars
parameters:
  serviceName: 'temp-service'
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-templatevars.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-templatevars.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("the return value is as expected and the correct files are written")
	expectedFilename1 := "sub/orig.go.txt"
	expectedContent1 := `package sub

import "fmt"

func PrintMessage() {
	fmt.Println("heya")
}
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename1,
			},
		},
	}
	require.Equal(t, expectedResponse, actualResponse)
	actual1, err := dir.ReadFile(context.TODO(), expectedFilename1)
	require.Nil(t, err)
	require.Equal(t, toUnix(expectedContent1), toUnix(string(actual1)))
}

func TestRender_ShouldComplainIfSyntaxErrorInTemplatedDefault(t *testing.T) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/invalid-generator-specs"
	targetdirpath := "../output/render-14"
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file with a syntax error in a templated default")
	renderspec := `generator: templatevars
parameters:
  serviceName: 'temp-service'
`
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-templatevars.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-templatevars.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("the return value is as expected and contains the correct error message")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	require.Equal(t, 1, len(actualResponse.Errors))
	expectedErrorPart := "variable declaration helloMessage has invalid default (this is an error in the generator spec): template: __defaultvalue_helloMessage:1:"
	require.Contains(t, actualResponse.Errors[0].Error(), expectedErrorPart)
}

func _testRender_emptyDefaultsSuccessTestCase(t *testing.T, testcase uint, renderspec string) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := fmt.Sprintf("../output/render-%d", testcase)
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator emptydefaults")
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-emptydefaults.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-emptydefaults.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("the return value is as expected and the correct files are written")
	expectedFilename1 := "strings.txt"
	expectedContent1 := `emptyStringDefault: ##
missingDefault: ##
`
	expectedResponse := &api.Response{
		Success: true,
		RenderedFiles: []api.FileResult{
			{
				Success:          true,
				RelativeFilePath: expectedFilename1,
			},
		},
	}
	require.Equal(t, expectedResponse, actualResponse)
	actual1, err := dir.ReadFile(context.TODO(), expectedFilename1)
	require.Nil(t, err)
	require.Equal(t, toUnix(expectedContent1), toUnix(string(actual1)))
}

func TestRender_ShouldWriteExpectedForEmptyDefaultsBothSet(t *testing.T) {
	renderspec := `generator: emptydefaults
parameters:
  emptyStringDefault: ''
  missingDefault: ''
`
	_testRender_emptyDefaultsSuccessTestCase(t, 15, renderspec)
}

func TestRender_ShouldWriteExpectedForEmptyDefaultsMissingSet(t *testing.T) {
	renderspec := `generator: emptydefaults
parameters:
  missingDefault: ''
`
	_testRender_emptyDefaultsSuccessTestCase(t, 16, renderspec)
}

func _testRender_emptyDefaultsErrorTestCase(t *testing.T, testcase uint, renderspec string, expectedContainedInErr string) {
	docs.Given("a valid generator source directory and a valid target directory")
	sourcedirpath := "../resources/valid-generator-simple"
	targetdirpath := fmt.Sprintf("../output/render-%d", testcase)
	require.Nil(t, os.RemoveAll(targetdirpath))
	require.Nil(t, os.Mkdir(targetdirpath, 0755))

	docs.Given("a valid render spec file for generator emptydefaults")
	dir := targetdir.Instance(context.TODO(), targetdirpath)
	require.Nil(t, dir.WriteFile(context.TODO(), "generated-emptydefaults.yaml", []byte(renderspec)))

	docs.When("Render is invoked")
	request := &api.Request{
		SourceBaseDir: sourcedirpath,
		TargetBaseDir: targetdirpath,
		RenderSpecFile: "generated-emptydefaults.yaml",
	}
	actualResponse := generatorlib.Render(context.TODO(), request)

	docs.Then("the return value is as expected and the correct error is raised")
	require.False(t, actualResponse.Success)
	require.Empty(t, actualResponse.RenderedFiles)
	require.Equal(t, 1, len(actualResponse.Errors))
	require.Contains(t, actualResponse.Errors[0].Error(), expectedContainedInErr)
}

func TestRender_ShouldWriteExpectedForEmptyDefaultsNoneSet(t *testing.T) {
	renderspec := `generator: emptydefaults
parameters: {}
`
	_testRender_emptyDefaultsErrorTestCase(t, 17, renderspec, "parameter 'missingDefault' is required but missing")
}

func TestRender_ShouldWriteExpectedForEmptyDefaultsEmptyStringSet(t *testing.T) {
	renderspec := `generator: emptydefaults
parameters:
  emptyStringDefault: ''
`
	_testRender_emptyDefaultsErrorTestCase(t, 18, renderspec, "parameter 'missingDefault' is required but missing")
}
