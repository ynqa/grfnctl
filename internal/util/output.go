package util

import (
	"github.com/spf13/pflag"
)

// Output represents the output format for displaying data.
type Output string

const (
	OutputJSON  Output = "json"
	OutputTable Output = "table"
)

// Ensure Output implements the pflag.Value interface
var _ pflag.Value = (*Output)(nil)

func (o *Output) String() string {
	return string(*o)
}

func (o *Output) Set(s string) error {
	switch s {
	case string(OutputJSON), string(OutputTable):
		*o = Output(s)
		return nil
	default:
		return nil
	}
}

func (o *Output) Type() string {
	return "output"
}
