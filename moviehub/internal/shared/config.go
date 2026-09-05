package shared

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Environment        string `mapstructure:"ENVIRONMENT"`
	ServerPort         string `mapstructure:"SERVER_PORT"` // Used for whatever service is running
	GatewayPort        string `mapstructure:"GATEWAY_PORT"`
	MovieServiceURL    string `mapstructure:"MOVIE_SERVICE_URL"`
	UserServiceURL     string `mapstructure:"USER_SERVICE_URL"`
	MongoURI           string `mapstructure:"MONGO_URI"`
	SECERETKEY         string `mapstructure:"SECERET_KEY"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL"`
	RedisURI           string `mapstructure:"REDIS_URI"`
	PostgresURI        string `mapstructure:"POSTGRES_URI"`
	CORSAllowedOrigins string `mapstructure:"CORS_ALLOWED_ORIGINS"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// Looks for a file named app.env in the specified path
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	// Bind all known environment variables so viper discovers them even without app.env
	viper.BindEnv("ENVIRONMENT")
	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("GATEWAY_PORT")
	viper.BindEnv("MOVIE_SERVICE_URL")
	viper.BindEnv("USER_SERVICE_URL")
	viper.BindEnv("MONGO_URI")
	viper.BindEnv("SECERET_KEY")
	viper.BindEnv("GOOGLE_CLIENT_ID")
	viper.BindEnv("GOOGLE_CLIENT_SECRET")
	viper.BindEnv("GOOGLE_REDIRECT_URL")
	viper.BindEnv("REDIS_URI")
	viper.BindEnv("POSTGRES_URI")
	viper.BindEnv("CORS_ALLOWED_ORIGINS")

	// Automatically override values with environment variables if they exist
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("No app.env file found. Falling back to environment variables.")
			err = nil // Ignore error if file doesn't exist but env vars are used
		} else {
			return
		}
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	// Fallback directly to os.Getenv for absolute safety in Docker and Kubernetes
	if config.MongoURI == "" {
		config.MongoURI = os.Getenv("MONGO_URI")
	}
	if config.PostgresURI == "" {
		config.PostgresURI = os.Getenv("POSTGRES_URI")
	}
	if config.RedisURI == "" {
		config.RedisURI = os.Getenv("REDIS_URI")
	}
	if config.ServerPort == "" {
		config.ServerPort = os.Getenv("SERVER_PORT")
	}
	if config.GatewayPort == "" {
		config.GatewayPort = os.Getenv("GATEWAY_PORT")
	}
	if config.MovieServiceURL == "" {
		config.MovieServiceURL = os.Getenv("MOVIE_SERVICE_URL")
	}
	if config.UserServiceURL == "" {
		config.UserServiceURL = os.Getenv("USER_SERVICE_URL")
	}
	if config.SECERETKEY == "" {
		config.SECERETKEY = os.Getenv("SECERET_KEY")
	}
	if config.CORSAllowedOrigins == "" {
		config.CORSAllowedOrigins = os.Getenv("CORS_ALLOWED_ORIGINS")
	}

	return
}
