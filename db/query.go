package db

import (
	"context"
	"github.com/petrostrak/picwise-ai/types"
)

func CreateAccount(account *types.Account) error {
	_, err := Bun.NewInsert().Model(account).Exec(context.Background())
	return err
}
