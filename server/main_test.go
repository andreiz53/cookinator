package server

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	database "github.com/andreiz53/cookinator/database/handlers"
	"github.com/andreiz53/cookinator/util"
)

func newTestServer(t *testing.T, store database.Store) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		TokenDuration:     time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	require.NotEmpty(t, server)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
