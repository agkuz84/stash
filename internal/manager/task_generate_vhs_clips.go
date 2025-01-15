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
	"github.com/stashapp/stash/pkg/models"
)

type generateVHSClipsJob struct {
	repository *models.Repository
}

func (j *generateVHSClipsJob) Execute(ctx context.Context, progress *job.Progress) error {
	logger.Infof("Starting VHS clip generation")

	return j.repository.WithTxn(ctx, func(ctx context.Context) error {
		markers, err := j.repository.SceneMarker.All(ctx)
		if err != nil {
			return fmt.Errorf("error getting markers: %v", err)
		}

		total := len(markers)
		progress.SetTotal(total) // Changed from float64(total)

		for _, marker := range markers {
			if job.IsCancelled(ctx) {
				return nil
			}

			scene, err := j.repository.Scene.Find(ctx, marker.SceneID)
			if err != nil {
				logger.Errorf("error finding scene for marker %d: %v", marker.ID, err)
				continue
			}

			outputPath := filepath.Join(GetInstance().Config.GetGeneratedPath(), "vhs_clips",
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

			progress.Increment()
		}
		return nil
	})
}

func (s *Manager) GenerateVHSClips(ctx context.Context) int {
	j := &generateVHSClipsJob{
		repository: &s.Repository, // Added & to get pointer to Repository
	}

	return s.JobManager.Add(ctx, "Generating VHS clips...", j)
}
