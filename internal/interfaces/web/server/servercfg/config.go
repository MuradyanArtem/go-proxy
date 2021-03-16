package servercfg

import "time"

type Config interface {
	Duration(string, time.Duration, string) *time.Duration
	String(string, string, string) *string
	Int(string, int, string) *int
}
