package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var (
	removeFirst float64
	removeLast  float64
	trimOutput  string
)

var trimCmd = &cobra.Command{
	Use:   "trim [input]",
	Short: "Trim the begining and/or end of a video",
	Run: func(cmd *cobra.Command, args []string) {
		trim(args[0], trimOutput, removeFirst, removeLast)
	},
}

func init() {
	trimCmd.Flags().Float64Var(&removeFirst, "removeFirst", 0, "Seconds to remove from the beginning")
	trimCmd.Flags().Float64Var(&removeLast, "removeLast", 0, "Seconds to remove from the end")
	trimCmd.Flags().StringVarP(&trimOutput, "output", "o", "output.mp4", "Output file path")

	rootCmd.AddCommand(trimCmd)
}

func trim(input, output string, removeFirst, removeLast float64) {
	duration, err := getVideoDuration(input)
	if err != nil {
		log.Fatalf("failed to probe video: %v", err)
	}

	start := removeFirst + 1
	end := duration - removeLast

	if end <= start {
		log.Fatalf("resulting duration is zero or negative (start=%.2f, end=%.2f)", start-1, end)
	}

	trimmer := ffmpeg.Input(input, ffmpeg.KwArgs{
		"ss": start,
		"to": end,
	}).
		Output(trimOutput, ffmpeg.KwArgs{
			"c": "copy",
		}).
		OverWriteOutput()

	if verbose {
		trimmer.ErrorToStdOut()
	}

	err = trimmer.Run()

	if err != nil {
		log.Fatalf("ffmpeg error: %v", err)
	}

	fmt.Printf("saved trimmed video to %s\n", trimOutput)
}

type probeFormat struct {
	Duration string `json:"duration"`
}

type probeData struct {
	Format probeFormat `json:"format"`
}

func getVideoDuration(path string) (float64, error) {
	data, err := ffmpeg.Probe(path)
	if err != nil {
		return 0, err
	}

	var probe probeData
	if err := json.Unmarshal([]byte(data), &probe); err != nil {
		return 0, err
	}

	return strconv.ParseFloat(probe.Format.Duration, 64)
}
