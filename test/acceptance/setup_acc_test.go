package acceptance

import (
	aulogging "github.com/StephanHCB/go-autumn-logging"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	beforeTest()
	code := m.Run()
	os.Exit(code)
}

func beforeTest() {
	aulogging.SetupNoLoggerForTesting()
}

