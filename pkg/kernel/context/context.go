package context

import (
	"context"
	"github.com/PaulBarrie/infra-worker/pkg/auth"
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

func GetUser() auth.User {
	return CurrentContext.Value(UserKey).(auth.User)
}

func SetUser(user *auth.User) {
	CurrentContext = context.WithValue(CurrentContext, UserKey, user)
}

func Clear() {
	CurrentContext = nil
}

func init() {
	if CurrentContext == nil {
		CurrentContext = context.Background()
	}
}
