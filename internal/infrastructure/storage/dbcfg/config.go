package dbcfg

import "time"

type Config interface {
	Int(string, int, string) *int
	Bool(string, bool, string) *bool
	String(string, string, string) *string
	Duration(string, time.Duration, string) *time.Duration
}
