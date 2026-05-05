package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var normalizeCmd = &cobra.Command{
	Use:   "normalize [inputs...]",
	Short: "Normalize the audio level of one or more video files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, input := range args {
			ext := filepath.Ext(input)
			base := strings.TrimSuffix(input, ext)
			tmp := base + "_norm" + ext

			normalize(input, tmp)

			if err := os.Rename(tmp, input); err != nil {
				log.Fatalf("failed to replace file: %v", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(normalizeCmd)
}

func normalize(input, output string) {
	normalizer := ffmpeg.Input(input).
		Output(output, ffmpeg.KwArgs{"af": "loudnorm"}).
		OverWriteOutput()

	if verbose {
		normalizer.ErrorToStdOut()
	}

	if err := normalizer.Run(); err != nil {
		log.Fatalf("ffmpeg error: %v", err)
	}

	fmt.Printf("saved normalized video to %s\n", output)
}
