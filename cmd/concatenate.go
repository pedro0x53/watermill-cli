package cmd

import (
	"log"

	"github.com/spf13/cobra"
	ffmpeg "github.com/u2takey/ffmpeg-go"
)

var (
	concatOutput string
	dimensions   string
)

var concatenateCmd = &cobra.Command{
	Use:   "concatenate [inputs...]",
	Short: "Concatenate multiple videos into one",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		err := concatenate(args, concatOutput)
		if err != nil {
			log.Println(err)
			log.Printf("concatenate failed on paths: %v", args)
		}
	},
}

func init() {
	concatenateCmd.Flags().StringVarP(&concatOutput, "output", "o", "output.mp4", "Output file path")
	concatenateCmd.Flags().StringVarP(&dimensions, "dimensions", "d", "1920:1080", "Output dimensions (WxH format, e.g. 1920:1080)")

	rootCmd.AddCommand(concatenateCmd)
}

func concatenate(inputs []string, output string) error {
	streams := make([]*ffmpeg.Stream, len(inputs)*2)
	for i, input := range inputs {
		in := ffmpeg.Input(input)
		streams[i*2] = in.Video().
			Filter("scale", ffmpeg.Args{dimensions + ":force_original_aspect_ratio=decrease"}).
			Filter("pad", ffmpeg.Args{dimensions + ":(ow-iw)/2:(oh-ih)/2"})
		streams[i*2+1] = in.Audio()
	}

	concatenator := ffmpeg.Concat(streams, ffmpeg.KwArgs{
		"v": 1,
		"a": 1,
	}).
		Output(output).
		OverWriteOutput()

	if verbose {
		concatenator.ErrorToStdOut()
	}

	err := concatenator.Run()

	if err != nil {
		log.Printf("ffmpeg error: %v", err)
		return err
	}

	if verbose { 
		log.Printf("concatenated %d videos into %s", len(inputs), output)
	}

	return nil
}
