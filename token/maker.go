package token

import (
	"time"

	"go.uber.org/fx"
)

type Maker interface {
	CreateToken(username string, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

var Module = fx.Options(
	fx.Provide(NewPasetoMaker),
)
