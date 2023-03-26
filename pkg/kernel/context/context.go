package context

import (
	"context"
	"sync"
)

var Current context.Context
var lock = &sync.Mutex{}

const (
	WorkingDirectoryKey = "working_directory"
	PluginKey           = "plugin"
	UserKey             = "user"
)

func init() {
	lock.Lock()
	defer lock.Unlock()
	if Current == nil {
		Current = context.Background()
	}
}

func SetWorkingDirectory(workDir string) {
	Current = context.WithValue(Current, WorkingDirectoryKey, workDir)
}

func GetWorkingDirectory() string {
	return Current.Value(WorkingDirectoryKey).(string)
}

func Clear() {
	Current = nil
}

func init() {
	if Current == nil {
		Current = context.Background()
	}
}
