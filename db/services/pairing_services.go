package services

import (
	"context"
	models "github.com/Opulentia-Trading/Arbitrage/models"
	"github.com/uptrace/bun"
)

func Get_platform_pairs(platform1_ID int64, asset1 string, database *bun.DB, ctx context.Context) []models.Pairing {
	var pairings []models.Pairing
	err := database.NewSelect().Model(&pairing).Where("platform_asset1 = ? AND asset1 = ?", platform_asset, asset1).Scan(ctx)
	if err != nil {
		panic(err)
	}

	return pairings
}
