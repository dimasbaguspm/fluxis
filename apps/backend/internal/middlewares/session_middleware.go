package middlewares

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/dimasbaguspm/fluxis/internal/repositories"
)

var authRepo = repositories.AuthRepository{}

func SessionMiddleware(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cat := ctx.Header("Authorization")

		if len(cat) > 7 && cat[:7] == "Bearer " {
			cat = cat[7:]
		}

		isValid := authRepo.IsTokenValid(cat)

		if !isValid {
			huma.WriteErr(api, ctx, repositories.AuthErrorInvalidAccessToken.GetStatus(), repositories.AuthErrorInvalidAccessToken.Error())
			return
		}

		next(ctx)
	}
}
