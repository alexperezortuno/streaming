package media

import (
	"fmt"
	"log/slog"
	"os/exec"
)

type TranscodeJob struct {
	VideoID   string
	InputPath string
	OutputDir string
	OnSuccess func()
	OnError   func(error)
}

type Transcoder struct {
	jobs chan TranscodeJob
}

func NewTranscoder(workers int) *Transcoder {
	t := &Transcoder{
		jobs: make(chan TranscodeJob, 100),
	}
	for i := 0; i < workers; i++ {
		go t.worker(i)
	}
	slog.Info("transcoder started", "workers", workers)
	return t
}

func (t *Transcoder) Enqueue(job TranscodeJob) {
	t.jobs <- job
}

func (t *Transcoder) worker(id int) {
	for job := range t.jobs {
		slog.Info("transcoding started",
			"worker", id,
			"video_id", job.VideoID,
			"input", job.InputPath,
		)

		if err := transcode(job.InputPath, job.OutputDir); err != nil {
			slog.Error("transcoding failed",
				"worker", id,
				"video_id", job.VideoID,
				"error", err,
			)
			if job.OnError != nil {
				job.OnError(err)
			}
			continue
		}

		slog.Info("transcoding completed",
			"worker", id,
			"video_id", job.VideoID,
		)
		if job.OnSuccess != nil {
			job.OnSuccess()
		}

		cleanupInput(job.InputPath)
	}
}

func transcode(inputPath, outputDir string) error {
	outputPath := fmt.Sprintf("%s/index.m3u8", outputDir)
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-profile:v", "baseline",
		"-level", "3.0",
		"-s", "1280x720",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		outputPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %w\noutput: %s", err, string(output))
	}
	return nil
}

func cleanupInput(inputPath string) {
	if err := exec.Command("rm", inputPath).Run(); err != nil {
		slog.Warn("cleanup input file", "path", inputPath, "error", err)
	}
}
