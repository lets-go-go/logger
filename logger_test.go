package logger

import "testing"

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
			name: "test",
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
