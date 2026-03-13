package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

func (pc *ProjectCache) InvalidateSingleProject(ctx context.Context, projectID pgtype.UUID) {
	_ = pc.c.Delete(ctx, cache.KeySingleProject(pc.hmacKey, projectID))
}

func (pc *ProjectCache) InvalidateSingleProjectByKey(ctx context.Context, orgID pgtype.UUID, key string) {
	_ = pc.c.Delete(ctx, cache.KeySingleProjectByKey(pc.hmacKey, orgID, key))
}

func (pc *ProjectCache) InvalidatePagedProjects(ctx context.Context) {
	_ = pc.c.Delete(ctx, cache.KeyPagedProjects(pc.hmacKey, nil))
}
