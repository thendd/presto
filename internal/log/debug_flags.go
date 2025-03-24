package log

type DebugLevel uint8

const (
	DebugLevelNone  DebugLevel = 1 << iota
	DebugLevelError DebugLevel = 1 << DebugLevelNone
	DebugLevelWarn  DebugLevel = 1 << DebugLevelError
	DebugLevelInfo  DebugLevel = 1 << DebugLevelWarn
	DebugAllFlags   DebugLevel = DebugLevelError | DebugLevelWarn | DebugLevelInfo
	DebugStdFlags   DebugLevel = DebugLevelNone | DebugLevelError
)