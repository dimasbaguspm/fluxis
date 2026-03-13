package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

type ProjectCache struct {
	c       cache.Cache
	cfg     cache.Config
	hmacKey string
}

func New(c cache.Cache) *ProjectCache {
	return &ProjectCache{
		c:       c,
		cfg:     c.GetConfig(),
		hmacKey: c.GetConfig().HMACKey,
	}
}

func (pc *ProjectCache) GetSingleProject(ctx context.Context, projectID pgtype.UUID, fetch func(context.Context) (domain.ProjectModel, error)) (domain.ProjectModel, error) {
	key := cache.KeySingleProject(pc.hmacKey, projectID)
	return cache.ReadThrough(ctx, pc.c, key, pc.cfg.DefaultTTL, func() (domain.ProjectModel, error) {
		return fetch(ctx)
	})
}

func (pc *ProjectCache) GetPagedProjects(ctx context.Context, params interface{}, fetch func(context.Context) (domain.ProjectsPagedModel, error)) (domain.ProjectsPagedModel, error) {
	key := cache.KeyPagedProjects(pc.hmacKey, params)
	return cache.ReadOrWrite(ctx, pc.c, key, pc.cfg.DefaultTTL, func(ctx context.Context) (domain.ProjectsPagedModel, error) {
		return fetch(ctx)
	})
}

func (pc *ProjectCache) InvalidateSingleProject(ctx context.Context, projectID pgtype.UUID) {
	_ = pc.c.Delete(ctx, cache.KeySingleProject(pc.hmacKey, projectID))
}

func (pc *ProjectCache) InvalidateSingleProjectByKey(ctx context.Context, orgID pgtype.UUID, key string) {
	_ = pc.c.Delete(ctx, cache.KeySingleProjectByKey(pc.hmacKey, orgID, key))
}
