package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

const (
	PROGRESS = "watermill-cli-progress"
	EDITED   = "_edited"
)

var videoExt = map[string]struct{}{
	".mp4": {},
	".mov": {},
}

var (
	root             string = "."
	intro            string
	outro            string
	runRemoveFirst   float64
	runRemoveLast    float64
	runNormalize     bool
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Automatically edit the video files from the current directory and subdirectories",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			root = args[0]
		}

		progress := loadProgress()

		progressFile, err := os.OpenFile(root+"/"+PROGRESS, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Println(err)
		}

		defer progressFile.Close()

		videos := make(chan struct{}, 5)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for path := range progress {
			wg.Add(1)
			videos <- struct{}{}

			go func(path string) {
				defer wg.Done()
				defer func() { <-videos }()

				ext := filepath.Ext(path)
				base := strings.TrimSuffix(path, ext)
				outputPath := base + "_edited" + ext
				concatPath := base + "_concat" + ext

				normPath := base + "_norm" + ext

				trim(path, outputPath, runRemoveFirst, runRemoveLast)

				inputForConcat := outputPath
				if runNormalize {
					normalize(outputPath, normPath)
					inputForConcat = normPath
				}

				concatenate([]string{intro, inputForConcat, outro}, concatPath)

				if runNormalize {
					os.Remove(normPath)
				}

				if err := os.Rename(concatPath, outputPath); err != nil {
					log.Fatalf("failed to replace edited file: %v", err)
				}

				mu.Lock()
				fmt.Fprintln(progressFile, path)
				mu.Unlock()
			}(path)
		}

		wg.Wait()
	},
}

func init() {
	runCmd.Flags().StringVarP(&intro, "intro", "i", "intro.mp4", "The video intro file path")
	runCmd.Flags().StringVarP(&outro, "outro", "o", "outro.mp4", "The video outro file path")
	runCmd.Flags().Float64Var(&runRemoveFirst, "removeFirst", 0, "Seconds to remove from the beginning")
	runCmd.Flags().Float64Var(&runRemoveLast, "removeLast", 0, "Seconds to remove from the end")
	runCmd.Flags().BoolVar(&runNormalize, "normalize", true, "Normalize audio levels")

	rootCmd.AddCommand(runCmd)
}

func loadProgress() map[string]struct{} {
	allFilePaths := scanDir()

	data, err := os.ReadFile(PROGRESS)

	if err != nil {
		return allFilePaths
	}

	for line := range strings.SplitSeq(strings.TrimSpace(string(data)), "\n") {
		if line != "" {
			delete(allFilePaths, line)
		}
	}

	return allFilePaths
}

func scanDir() map[string]struct{} {
	filePaths := make(map[string]struct{})

	introBase := filepath.Base(intro)
	outroBase := filepath.Base(outro)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		baseName := info.Name()
		ext := filepath.Ext(baseName)
		fileName := strings.TrimSuffix(baseName, ext)
		notEdited := !strings.HasSuffix(fileName, EDITED)
		_, validVideo := videoExt[ext]

		if notEdited && validVideo && baseName != introBase && baseName != outroBase {
			filePaths[path] = struct{}{}
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

	return filePaths
}
