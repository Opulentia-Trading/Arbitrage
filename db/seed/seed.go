package main

import (
	"context"
	"fmt"

	database "github.com/Opulentia-Trading/Arbitrage/db"
	"github.com/Opulentia-Trading/Arbitrage/db/services"
	env "github.com/Opulentia-Trading/Arbitrage/env"
	models "github.com/Opulentia-Trading/Arbitrage/models"
)

func main() {
	// Load the env variables
	env.Load_env("./env/.env")

	db := database.Establish_db_connection()

	// create context
	ctx := context.Background()

	err := db.ResetModel(ctx, (*models.Pairing)(nil))
	if err != nil {
		panic(err)
	}
}
