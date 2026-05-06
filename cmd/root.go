package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"
)

const WATERMIL = "watermill-cli"

var verbose bool

var rootCmd = &cobra.Command{
	Use:   "watermill-cli",
	Short: "A simple video editor CLI",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !verbose {
			log.SetOutput(io.Discard)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	log.SetPrefix(WATERMIL + ": ")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show ffmpeg output")
}
