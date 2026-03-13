package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/dimasbaguspm/fluxis/pkg/domain"
	"github.com/jackc/pgx/v5/pgtype"
)

func (oc *OrgCache) GetSingleOrg(ctx context.Context, orgID pgtype.UUID, fetch func(context.Context) (domain.OrganisationModel, error)) (domain.OrganisationModel, error) {
	key := cache.KeySingleOrg(oc.hmacKey, orgID)
	return cache.ReadThrough(ctx, oc.c, key, oc.cfg.DefaultTTL, func() (domain.OrganisationModel, error) {
		return fetch(ctx)
	})
}

func (oc *OrgCache) GetPagedOrganizations(ctx context.Context, params interface{}, fetch func(context.Context) (domain.OrganisationPagedModel, error)) (domain.OrganisationPagedModel, error) {
	key := cache.KeyPagedOrganizations(oc.hmacKey, params)
	return cache.ReadOrWrite(ctx, oc.c, key, oc.cfg.DefaultTTL, func(ctx context.Context) (domain.OrganisationPagedModel, error) {
		return fetch(ctx)
	})
}
