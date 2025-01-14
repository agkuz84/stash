package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) GenerateVHSClips(ctx context.Context) (bool, error) {
	job := &manager.GenerateVHSClipsJob{
		BaseJob: manager.BaseJob{},
	}

	manager.GetInstance().JobManager.Add(ctx, job)
	return true, nil
}
