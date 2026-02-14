package boot

import (
	"fmt"
	"os"

	_ "github.com/lib/pq" // postgres driver don`t delete
	"github.com/netbill/organizations-svc/internal/bucket"
	"github.com/netbill/organizations-svc/internal/messenger"
	"github.com/netbill/organizations-svc/internal/rest"
	"github.com/netbill/organizations-svc/internal/tokenmanager"
	"github.com/spf13/viper"
)

const ServiceName = "organization-svc"

type DatabaseConfig struct {
	SQL struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"sql"`
}

type AuthConfig struct {
	Tokens tokenmanager.Config `mapstructure:"tokens"`
}

type S3Config struct {
	AWS   AwsConfig     `mapstructure:"aws"`
	Media bucket.Config `mapstructure:"media"`
}

type Config struct {
	Log      LogConfig        `mapstructure:"log"`
	Rest     rest.Config      `mapstructure:"rest"`
	Auth     AuthConfig       `mapstructure:"auth"`
	Kafka    messenger.Config `mapstructure:"kafka"`
	Database DatabaseConfig   `mapstructure:"database"`
	S3       S3Config         `mapstructure:"s3"`
}

func LoadConfig() Config {
	configPath := os.Getenv("KV_VIPER_FILE")
	if configPath == "" {
		panic(fmt.Errorf("KV_VIPER_FILE env var is not set"))
	}
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %s", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("error unmarshalling config: %s", err))
	}

	return config
}
