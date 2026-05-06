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
	TMP      = "_tmp"
	INTRO    = "intro.mp4"
	OUTRO    = "outro.mp4"
)

var videoExt = map[string]struct{}{
	".mp4": {},
	".mov": {},
}

var (
	root           string
	intro          string
	outro          string
	runRemoveFirst float64
	runRemoveLast  float64
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Automatically edit the video files from the current directory and subdirectories",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			root = args[0]
		}

		progress := loadProgress()

		if len(progress) == 0 {
			log.Println("No pending files were found")
			os.Exit(0)
		}

		progressFile, err := os.OpenFile(root+"/"+PROGRESS, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln(err)
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
				concatPath := base + "_tmp" + ext

				if err := trim(path, outputPath, runRemoveFirst, runRemoveLast); err != nil {
					message := fmt.Sprintf("trim failed on path %v", path)

					if !verbose {
						fmt.Fprintf(os.Stderr, WATERMIL+": %v", message)
					}

					log.Println(message)

					os.Remove(outputPath)

					return
				}

				introPath := intro
				if intro == INTRO {
					introPath = root + "/" + intro
				}

				outroPath := outro
				if outro == OUTRO {
					outroPath = root + "/" + outro
				}

				if err := concatenate([]string{introPath, outputPath, outroPath}, concatPath); err != nil {
					message := fmt.Sprintf("concatenate failed on path %v", path)

					if !verbose {
						fmt.Fprintf(os.Stderr, WATERMIL+": %v", message)
					}

					log.Println(message)

					os.Remove(outputPath)
					os.Remove(concatPath)

					return
				}

				if err := os.Rename(concatPath, outputPath); err != nil {
					log.Printf("failed to replace edited file: %v", err)

					os.Remove(outputPath)
					os.Remove(concatPath)

					return
				}

				if verbose {
					log.Printf("renamed %v to %v", concatPath, outputPath)
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
	runCmd.Flags().StringVarP(&intro, "intro", "i", INTRO, "The video intro file path")
	runCmd.Flags().StringVarP(&outro, "outro", "o", OUTRO, "The video outro file path")
	runCmd.Flags().Float64Var(&runRemoveFirst, "removeFirst", 0, "Seconds to remove from the beginning")
	runCmd.Flags().Float64Var(&runRemoveLast, "removeLast", 0, "Seconds to remove from the end")

	rootCmd.AddCommand(runCmd)
}

func loadProgress() map[string]struct{} {
	allFilePaths := scanDir()

	data, err := os.ReadFile(root + "/" + PROGRESS)

	if err != nil {
		return allFilePaths
	}

	for line := range strings.SplitSeq(strings.TrimSpace(string(data)), "\n") {
		if verbose {
			log.Println(line, "found")
		}

		if line != "" {
			delete(allFilePaths, line)
			if verbose {
				log.Println(line, "already edited")
			}
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

		isEdited := strings.HasSuffix(fileName, EDITED)
		isTmp := strings.HasSuffix(fileName, TMP)
		isIntro := baseName == introBase
		isOutro := baseName == outroBase
		_, isValid := videoExt[ext]

		if !isEdited && !isTmp && isValid && !isIntro && !isOutro {
			filePaths[path] = struct{}{}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return filePaths
}
