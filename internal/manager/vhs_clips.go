package manager

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
)

type GenerateVHSClipsJob struct {
	BaseJob
}

func (j *GenerateVHSClipsJob) Execute(ctx context.Context) error {
	logger.Infof("Starting VHS clip generation")

	markers, err := db.Marker.All(ctx)
	if err != nil {
		return fmt.Errorf("error getting markers: %v", err)
	}

	total := len(markers)
	for i, marker := range markers {
		if job.IsCancelled() {
			return nil
		}

		scene, err := db.Scene.Find(ctx, marker.SceneID)
		if err != nil {
			logger.Errorf("error finding scene for marker %d: %v", marker.ID, err)
			continue
		}

		outputPath := filepath.Join(instance.Paths.Generated, "vhs_clips",
			strconv.Itoa(marker.SceneID), strconv.Itoa(marker.ID)+".mp4")

		if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
			logger.Errorf("error creating directory for marker %d: %v", marker.ID, err)
			continue
		}

		args := []string{
			"-ss", fmt.Sprintf("%.3f", marker.Seconds),
			"-t", "5",
			"-i", scene.Path,
			"-c:v", "libx264",
			"-preset", "fast",
			"-tune", "fastdecode",
			"-force_key_frames", "expr:gte(t,0)",
			"-y",
			outputPath,
		}

		cmd := exec.Command("ffmpeg", args...)
		if err := cmd.Run(); err != nil {
			logger.Errorf("error generating clip for marker %d: %v", marker.ID, err)
			continue
		}

		j.progress = float64(i+1) / float64(total)
	}

	logger.Info("Finished generating VHS clips")
	return nil
}

func (j *GenerateVHSClipsJob) Description() string {
	return "Generating VHS Clips"
}
