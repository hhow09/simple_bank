package main

import (
	"github.com/hhow09/simple_bank/api"
	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(func() util.ConfigPath {
			return "."
		}),
		fx.Provide(util.LoadConfig),
		token.Module,
		db.Module,
		api.Module,
	).Run()
}
