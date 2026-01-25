/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"clockify-app/cmd"
	"fmt"
	"os"
)

var version = "dev" // default value

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
