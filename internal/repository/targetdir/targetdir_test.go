package targetdir

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

// these tests add coverage for some internal error conditions only

func TestCheckValid_TrailingSlash(t *testing.T) {
	cut := Instance(context.TODO(), "./has-a-slash/")
	actualErr := cut.CheckValid(context.TODO())
	expected := "error invalid target directory: baseDir ./has-a-slash/ must not contain trailing slash"
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestCheckValid_TrailingBackslash(t *testing.T) {
	cut := Instance(context.TODO(), "./has-a-backslash\\")
	actualErr := cut.CheckValid(context.TODO())
	expected := "error invalid target directory: baseDir ./has-a-backslash\\ must not contain trailing slash"
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestCheckValid_NotADirectory(t *testing.T) {
	cut := Instance(context.TODO(), "./targetdir.go")
	actualErr := cut.CheckValid(context.TODO())
	expected := "error invalid target directory: baseDir ./targetdir.go must be a directory"
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestCheckValid_DoesNotExist(t *testing.T) {
	cut := Instance(context.TODO(), "./targetdir/does-not-exist")
	actualErr := cut.CheckValid(context.TODO())
	expected := "error invalid target directory: baseDir ./targetdir/does-not-exist does not exist"
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestReadFile_Invalid(t *testing.T) {
	cut := Instance(context.TODO(), "./targetdir/does-not-exist")
	actualBytes, actualErr := cut.ReadFile(context.TODO(), "pointless")
	expected := "error invalid target directory: baseDir ./targetdir/does-not-exist does not exist"
	require.Empty(t, actualBytes)
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestWriteFile_Invalid(t *testing.T) {
	cut := Instance(context.TODO(), "./targetdir/does-not-exist")
	actualErr := cut.WriteFile(context.TODO(), "pointless", []byte{})
	expected := "error invalid target directory: baseDir ./targetdir/does-not-exist does not exist"
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestWriteFile_InTheWay1(t *testing.T) {
	cut := Instance(context.TODO(), ".")
	actualErr := cut.WriteFile(context.TODO(), "intheway_test.go/pointlessdir/pointless", []byte{})
	expected := "cannot create path up to intheway_test.go/pointlessdir, something is in the way or invalid path: mkdir intheway_test.go: The system cannot find the path specified."
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}

func TestWriteFile_InTheWay2(t *testing.T) {
	cut := Instance(context.TODO(), ".")
	actualErr := cut.WriteFile(context.TODO(), "intheway_test.go/pointless", []byte{})
	expected := "cannot create path up to intheway_test.go, something is in the way"
	require.NotNil(t, expected, actualErr)
	require.Equal(t, expected, actualErr.Error())
}
