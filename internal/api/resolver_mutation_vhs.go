package api

import (
	"context"

	"github.com/stashapp/stash/internal/manager"
)

func (r *mutationResolver) GenerateVHSClips(ctx context.Context) (bool, error) {
	manager.GetInstance().GenerateVHSClips(ctx)
	return true, nil
}
