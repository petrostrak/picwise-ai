package view

import (
	"context"

	"github.com/petrostrak/picwise-ai/types"
)

func AuthenticatedUser(ctx context.Context) types.AuthenticatedUser {
	user, ok := ctx.Value("user").(types.AuthenticatedUser)
	if !ok {
		return types.AuthenticatedUser{}
	}
	return user
}
