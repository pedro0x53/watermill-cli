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
		err := trim(args[0], trimOutput, removeFirst, removeLast)
		if err != nil {
			log.Println(err)
			log.Printf("trim failed on path: %v", args[0])
		}
	},
}

func init() {
	trimCmd.Flags().Float64Var(&removeFirst, "removeFirst", 0, "Seconds to remove from the beginning")
	trimCmd.Flags().Float64Var(&removeLast, "removeLast", 0, "Seconds to remove from the end")
	trimCmd.Flags().StringVarP(&trimOutput, "output", "o", "output.mp4", "Output file path")

	rootCmd.AddCommand(trimCmd)
}

func trim(input, output string, removeFirst, removeLast float64) error {
	duration, err := getVideoDuration(input)
	if err != nil {
		log.Printf("failed to probe video: %v", err)
		return err
	}

	start := removeFirst
	end := duration - removeLast

	if end <= start {
		err := fmt.Errorf("resulting duration is zero or negative (start=%.2f, end=%.2f)", start, end)
		log.Print(err)
		return err
	}

	trimmer := ffmpeg.Input(input, ffmpeg.KwArgs{
		"ss": start,
		"to": end,
	}).
		Output(output, ffmpeg.KwArgs{
			"c": "copy",
		}).
		OverWriteOutput()

	if verbose {
		trimmer.ErrorToStdOut()
	}

	err = trimmer.Run()

	if err != nil {
		log.Printf("ffmpeg error: %v", err)
		return err
	}

	if verbose {
		log.Printf("saved trimmed video to %s", output)
	}

	return nil
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
