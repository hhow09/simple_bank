package api

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/hhow09/simple_bank/db/sqlc"
	"github.com/hhow09/simple_bank/lib"
	"github.com/hhow09/simple_bank/token"
	"github.com/hhow09/simple_bank/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
)

// set gin into TestMode to get cleaner logs
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func newTestServer(t *testing.T, mockstore db.Store) *Server {
	var s *Server
	fx.New(
		fx.Provide(func() util.ConfigPath {
			return "../"
		}),
		fx.Provide(func(path util.ConfigPath) util.Config {
			config, err := util.LoadConfig(path)
			require.NoError(t, err)
			config.TokenSymmetricKey = util.RandomString(32)
			config.AccessTokenDuration = time.Minute
			return config
		}),
		token.Module,
		// user mock store
		fx.Provide(func() db.Store {
			return mockstore
		}),
		lib.Module,
		Module,
		fx.Populate(&s),
	)

	return s
}
