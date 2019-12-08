package main

import (
	"gopkg.in/alecthomas/kingpin.v2"
)

// Flags : command line flags for service
type Flags struct {
	Preseed bool
	Verbose bool
}

var (
	preseed = kingpin.Flag("preseed", "Enables preseeding of database.").Short('p').Bool()
	verbose = kingpin.Flag("verbose", "Enables verbose output.").Short('v').Bool()
)

// GetFlags : gets command line flags at runtime
func GetFlags() Flags {
	kingpin.Parse()
	return Flags{Preseed: *preseed, Verbose: *verbose}
}
