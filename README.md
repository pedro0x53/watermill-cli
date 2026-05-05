# watermill-cli

A simple CLI for batch video editing - built as a learning project for Go and [Cobra](https://github.com/spf13/cobra).

Wraps [ffmpeg-go](https://github.com/u2takey/ffmpeg-go) to automate common edits: trimming, concatenation, and adding intros/outros to a batch of videos.

## Commands

- `run` — scan a directory for videos and apply trim + intro/outro to each one in parallel
- `trim [input]` — trim seconds from the beginning and/or end of a video
- `concatenate [inputs...]` — concatenate multiple videos into one

## Requirements

`ffmpeg` must be installed and available in your `PATH`.

## Usage

```sh
# Trim 2s from the start and 3s from the end
watermill-cli trim input.mp4 --removeFirst 2 --removeLast 3 -o output.mp4

# Concatenate videos
watermill-cli concatenate a.mp4 b.mp4 c.mp4 -o final.mp4

# Batch process all videos in a directory
watermill-cli run ./videos --intro intro.mp4 --outro outro.mp4 --removeFirst 1
```
