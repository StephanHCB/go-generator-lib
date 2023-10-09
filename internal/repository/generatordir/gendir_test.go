package generatordir

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGlobInvalid(t *testing.T) {
	ctx := context.TODO()
	cut := Instance(ctx, "../../../test/resources/valid-generator-simple")

	actual, err := cut.Glob(ctx, "[]a]")
	require.Empty(t, actual)
	require.NotNil(t, err)
	require.Equal(t, "syntax error in pattern", err.Error())
}

func TestGlobForbidden(t *testing.T) {
	ctx := context.TODO()
	cut := Instance(ctx, "../../../test/resources/valid-generator-simple")

	// this could be ../../../etc/passwd
	actual, err := cut.Glob(ctx, "src/sub/../../../valid-generator-structured/*.tmpl")
	require.Empty(t, actual)
	require.NotNil(t, err)
	require.Equal(t, "file glob src/sub/../../../valid-generator-structured/*.tmpl leads to file that is not inside base directory ../../../test/resources/valid-generator-simple - this is forbidden", err.Error())
}
