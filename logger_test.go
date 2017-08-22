package logger

import (
	"testing"
)

func TestInit(t *testing.T) {
	type args struct {
		logDir   string
		logFile  string
		minLevel LEVEL
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "InitTest",
			args: args{
				logDir:   ".",
				logFile:  "hello",
				minLevel: DEBUG,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Init(tt.args.logDir, tt.args.logFile, tt.args.minLevel)
		})
	}
}

func TestTrace(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Trace(tt.args.v...)
		})
	}
}

func TestTraceln(t *testing.T) {
	type args struct {
		v []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Traceln(tt.args.v...)
		})
	}
}

func TestTracef(t *testing.T) {
	type args struct {
		format string
		v      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Tracef(tt.args.format, tt.args.v...)
		})
	}
}
