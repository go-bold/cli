package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	version = "0.0.1"
	rootCmd = &cobra.Command{
		Use:   "bold",
		Short: "Go Bold - framework for Go",
		Long: `Go Bold Framework CLI

A rapid development framework that combines developer experience 
with Go's performance.

Get started:
  bold new myapp     Create a new Go Bold application
  bold serve         Start the development server
  bold help			 Show this help
  bold version 		 Show current version`,
		Version: version,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{printf "%s\n" .Version}}`)
}
