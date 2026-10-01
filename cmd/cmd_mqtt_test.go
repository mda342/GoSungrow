package cmd

import (
	"testing"

	"github.com/MickMake/GoUnify/cmdLog"
)

// --debug used to be parsed but never reach the logger, because the loggers
// are built in init() before flags are known. SetLogLevel is the fix, so verify
// it actually changes the level that was constructed earlier.
func TestSetLogLevelChangesLoggerLevel(t *testing.T) {
	c := NewCmdMqtt("")
	if c == nil {
		t.Fatal("NewCmdMqtt returned nil")
	}
	if got := c.log.GetLogLevel(); got != cmdLog.LogLevelInfoStr {
		t.Fatalf("default log level = %q, want %q", got, cmdLog.LogLevelInfoStr)
	}

	c.SetLogLevel(cmdLog.LogLevelDebugStr)
	if got := c.log.GetLogLevel(); got != cmdLog.LogLevelDebugStr {
		t.Errorf("after SetLogLevel(debug), log level = %q, want %q", got, cmdLog.LogLevelDebugStr)
	}
}

// SetLogLevel is called from ProcessArgs before the MQTT client exists, so it
// must not panic on a nil client.
func TestSetLogLevelBeforeClientExists(t *testing.T) {
	c := NewCmdMqtt("")
	c.Client = nil
	c.SetLogLevel(cmdLog.LogLevelDebugStr)
	if got := c.log.GetLogLevel(); got != cmdLog.LogLevelDebugStr {
		t.Errorf("log level = %q, want %q", got, cmdLog.LogLevelDebugStr)
	}
}

// An empty level must leave the logger alone rather than resetting it.
func TestSetLogLevelEmptyIsNoOp(t *testing.T) {
	c := NewCmdMqtt("")
	c.SetLogLevel(cmdLog.LogLevelDebugStr)
	c.SetLogLevel("")
	if got := c.log.GetLogLevel(); got != cmdLog.LogLevelDebugStr {
		t.Errorf("log level = %q, want %q (empty should be a no-op)", got, cmdLog.LogLevelDebugStr)
	}
}
