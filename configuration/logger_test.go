package configuration_test

import (
	"dforum-app/configuration"
	"testing"
)

func TestLoggerOutput(t *testing.T) {
	configuration.InitLogger("/home/henri/DoNotBackup/coding-projects/dforums-app/configuration/dfd.log")
}
