package cache

import (
	"context"

	"github.com/dimasbaguspm/fluxis/pkg/cache"
	"github.com/jackc/pgx/v5/pgtype"
)

func (oc *OrgCache) InvalidateSingleOrg(ctx context.Context, orgID pgtype.UUID) {
	_ = oc.c.Delete(ctx, cache.KeySingleOrg(oc.hmacKey, orgID))
}

func (oc *OrgCache) InvalidatePagedOrganizations(ctx context.Context) {
	_ = oc.c.Delete(ctx, cache.KeyPagedOrganizations(oc.hmacKey, nil))
}
