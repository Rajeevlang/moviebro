package main

import (
	"context"

	"moviesapi/internal/movie"
	"moviesapi/internal/shared"

	"go.uber.org/zap"
)

type Application struct {
	config shared.Config
}

func main() {
	// 1. Initialize Logging
	shared.InitLogger()

	// 2. Load Configuration
	config, err := shared.LoadConfig("./")
	if err != nil {
		shared.Log.Fatal("cannot load config", zap.Error(err))
	}

	app := &Application{
		config: config,
	}

	// 3. Setup MongoDB
	client, err := shared.ConnectMongoDB(app.config.MongoURI)
	if err != nil {
		shared.Log.Fatal("failed to connect to mongodb", zap.Error(err))
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			shared.Log.Fatal("failed to disconnect mongodb", zap.Error(err))
		}
	}()

	// 4. Setup Repository & Handler
	db := client.Database("movie_catalog")
	collection := db.Collection("movies")
	collectionpl := db.Collection("playlists")
	repo := movie.NewMovieRepository(collection)
	repo2 := movie.PLrepoHandler(collectionpl)

	// Start Redis Consumer for automatic playlist creation
	if err := shared.InitRedis(app.config.RedisURI); err != nil {
		shared.Log.Fatal("failed to initialize Redis", zap.Error(err))
	}
	go movie.StartRedisConsumer(context.Background(), repo2)

	redisDb := movie.NewRedisdb(shared.RedisDB)

	// Automatically build indexes to ensure our queries are fast
	if err := repo.EnsureIndexes(context.Background()); err != nil {
		shared.Log.Fatal("failed to build database indexes", zap.Error(err))
	}

	handler := movie.NewHandler(repo, repo2, redisDb)

	// 5. Setup Gin Router
	router := app.routes()

	// Register movie routes
	movie.RegisterRoutes(router, handler)

	// 6. Start the server
	shared.Log.Info("Starting Movie Service", zap.String("port", app.config.ServerPort))
	err = router.Run(":" + app.config.ServerPort)
	if err != nil {
		shared.Log.Fatal("failed to start server", zap.Error(err))
	}
}
