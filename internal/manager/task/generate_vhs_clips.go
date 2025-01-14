package task

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/stashapp/stash/pkg/job"
	"github.com/stashapp/stash/pkg/logger"
	"github.com/stashapp/stash/pkg/models"
)

type GenerateVHSClipsJob struct {
	Database models.ReaderWriter
	Paths    *models.Paths
}

func (j *GenerateVHSClipsJob) Execute(ctx context.Context, progress *job.Progress) error {
	markers, err := j.Database.Marker.All(ctx)
	if err != nil {
		return fmt.Errorf("error getting markers: %v", err)
	}

	progress.SetTotal(len(markers))

	for _, marker := range markers {
		if job.IsCancelled(ctx) {
			logger.Info("Cancelled generating VHS clips")
			return nil
		}

		taskDesc := fmt.Sprintf("Processing marker %d", marker.ID)
		progress.ExecuteTask(taskDesc, func() {
			scene, err := j.Database.Scene.Find(ctx, marker.SceneID)
			if err != nil {
				logger.Errorf("error finding scene for marker %d: %v", marker.ID, err)
				return
			}

			outputPath := filepath.Join(j.Paths.Generated, "vhs_clips",
				strconv.Itoa(marker.SceneID), strconv.Itoa(marker.ID)+".mp4")

			if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
				logger.Errorf("error creating directory for marker %d: %v", marker.ID, err)
				return
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
				return
			}
		})
	}

	logger.Info("Finished generating VHS clips")
	return nil
}
