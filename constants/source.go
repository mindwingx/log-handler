package constants

import "time"

const (
	EnvFile                        = "env"
	SlowSqlThreshold time.Duration = 5
	TmpLockFile                    = "/tmp/log-handler.lock"
	TimestampLayout  string        = "2006-01-02 03:04:05"
)
