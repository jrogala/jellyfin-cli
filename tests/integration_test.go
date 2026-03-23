package tests

import (
	"os"
	"testing"

	"github.com/cucumber/godog"
)

var globalMock *MockServer

func TestMain(m *testing.M) {
	globalMock = NewMockServer()
	code := m.Run()
	globalMock.Close()
	os.Exit(code)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: initializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			TestingT: t,
		},
	}
	if suite.Run() != 0 {
		t.Fatal("non-zero exit status from godog")
	}
}
