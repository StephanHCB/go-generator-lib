package acceptance

import (
	"context"
	generatorlib "github.com/StephanHCB/go-generator-lib"
	"github.com/StephanHCB/go-generator-lib/docs"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFindGeneratorNames_ShouldReturnCorrectList(t *testing.T) {
	docs.Given("a valid generator source directory")
	sourcedir := "../resources/valid-generator-simple"

	docs.When("FindGeneratorNames is invoked")
	actual, err := generatorlib.FindGeneratorNames(context.TODO(), sourcedir)

	docs.Then("the list of available generators is returned")
	expected := []string{"docker", "items", "main", "templatevars"}
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestFindGeneratorNames_ShouldComplainMissingDirectory(t *testing.T) {
	docs.Given("a nonexistant generator source directory")
	sourcedir := "../resources/invalid-does-not-exist"

	docs.When("FindGeneratorNames is invoked")
	actual, err := generatorlib.FindGeneratorNames(context.TODO(), sourcedir)

	docs.Then("an appropriate error is returned and the resulting list is empty")
	require.Empty(t, actual)
	require.NotNil(t, err)
	expectedErrorMsg := "invalid generator directory: baseDir ../resources/invalid-does-not-exist does not exist"
	require.Equal(t, expectedErrorMsg, err.Error())
}

func TestFindGeneratorNames_ShouldComplainNotDirectory(t *testing.T) {
	docs.Given("a regular file as generator source directory")
	sourcedir := "../resources/valid-generator-simple/generator-docker.yaml"

	docs.When("FindGeneratorNames is invoked")
	actual, err := generatorlib.FindGeneratorNames(context.TODO(), sourcedir)

	docs.Then("an appropriate error is returned and the resulting list is empty")
	require.Empty(t, actual)
	require.NotNil(t, err)
	expectedErrorMsg := "invalid generator directory: baseDir ../resources/valid-generator-simple/generator-docker.yaml must be a directory"
	require.Equal(t, expectedErrorMsg, err.Error())
}
