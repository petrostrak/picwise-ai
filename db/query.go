package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/petrostrak/picwise-ai/types"
)

func CreateAccount(account *types.Account) error {
	_, err := Bun.
		NewInsert().
		Model(account).
		Exec(context.Background())
	return err
}

func GetAccountByUserID(userID uuid.UUID) (types.Account, error) {
	var account types.Account
	err := Bun.
		NewInsert().
		Model(&account).
		Where("user_id = ?", userID).
		Scan(context.Background())
	return account, err
}

func UpdateAccount(account *types.Account) error {
	_, err := Bun.
		NewUpdate().
		Model(account).
		WherePK("user_id").
		Exec(context.Background())
	return err
}

func GetImagesByUserID(userID uuid.UUID) ([]types.Image, error) {
	var images []types.Image
	err := Bun.
		NewSelect().
		Model(&images).
		Where("deleted = ?", false).
		Where("user_id = ?", userID).
		Scan(context.Background())
	return images, err
}
