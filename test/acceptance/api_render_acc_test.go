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
