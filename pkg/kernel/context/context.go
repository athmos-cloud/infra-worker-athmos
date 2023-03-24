package context

import (
	"context"
)

var CurrentContext context.Context

const (
	WorkingDirectoryKey = "working_directory"
	PluginKey           = "plugin"
	UserKey             = "user"
)

func SetWorkingDirectory(workDir string) {
	CurrentContext = context.WithValue(CurrentContext, WorkingDirectoryKey, workDir)
}

func GetWorkingDirectory() string {
	return CurrentContext.Value(WorkingDirectoryKey).(string)
}

func Clear() {
	CurrentContext = nil
}

func init() {
	if CurrentContext == nil {
		CurrentContext = context.Background()
	}
}
