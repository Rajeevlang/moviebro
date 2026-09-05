package main

import (
	"context"
	"moviesapi/internal/shared"
	"moviesapi/internal/user"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	shared.InitLogger()
	defer shared.Log.Sync()

	// Load config
	Mconfig, err := shared.LoadConfig(".")
	if err != nil {
		shared.Log.Fatal("cannot load config", zap.Error(err))
	}

	if err := user.InitOpenID(Mconfig.GoogleClientID, Mconfig.GoogleClientSecret, Mconfig.GoogleRedirectURL); err != nil {
		shared.Log.Fatal("failed to initialize OpenID", zap.Error(err))
	}

	if err := shared.InitRedis(Mconfig.RedisURI); err != nil {
		shared.Log.Fatal("failed to initialize Redis", zap.Error(err))
	}

	ctx := context.Background()

	connStr := Mconfig.PostgresURI
	if connStr == "" {
		// Fallback for local development if not in Docker
		connStr = "postgres://app_user:secure_pass_123@localhost:5434/app_db"
	}

	// 1. Parse the standard configuration
	config, err := pgxpool.ParseConfig(connStr)

	if err != nil {
		shared.Log.Fatal("failed to connect the connstring", zap.Error(err))
	}

	// 2. Fine-tune production settings (Optional but recommended)
	config.MaxConns = 25                        // Maximum active connections
	config.MaxConnLifetime = 1 * time.Hour      // Recycle connections after an hour
	config.HealthCheckPeriod = 30 * time.Second // Check for dead connections

	getconn, err := pgx(config, ctx)
	if err != nil {
		shared.Log.Fatal("failed to connect the connstring", zap.Error(err))

	}
	nwrepo := user.Newpoolstruct(getconn)
	nwserv := user.NewService(nwrepo, Mconfig)

	handler := NewHandler(nwserv)

	r := gin.Default()
	RegisterRoutes(r, handler)

	shared.Log.Info("Starting User Service", zap.String("port", Mconfig.ServerPort))
	r.Run(":" + Mconfig.ServerPort)
}

func pgx(poolconn *pgxpool.Config, ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.NewWithConfig(ctx, poolconn)

	if err != nil {
		shared.Log.Fatal("failed to connect a pool ", zap.Error(err))
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil

}
