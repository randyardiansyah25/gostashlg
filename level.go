package gostashlg

type Level string

const (
	LOG   Level = "LOG"
	WARN  Level = "WARN"
	ERROR Level = "ERROR"
	INFO  Level = "INFO"
	DEBUG Level = "DEBUG"
	PRINT Level = "PRINT"
	TRACE Level = "TRACE"
)
