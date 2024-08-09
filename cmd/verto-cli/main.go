package main

import (
	"fmt"
	"os"
)

var (
	AppName string
	Version string
)

func main() {
	fmt.Fprintf(os.Stdout, "%s %s\n", AppName, Version)
	os.Exit(0)
}
