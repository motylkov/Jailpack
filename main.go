// Package main is the entry point for Jailpack.
package main

import (
	"jailpack/cmd"
)

// VERSION is the current version of the FreeBSD Command Manager.
var VERSION = "0.02" //nolint: gochecknoglobals

func main() {
	cmd.Execute()
}
