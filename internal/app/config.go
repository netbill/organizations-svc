package app

import (
	"fmt"
	"os"
	"time"

	"github.com/netbill/awsx"
	"github.com/spf13/viper"
)

type ServiceCfg struct {
	Name string `mapstructure:"name"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type DatabaseConfig struct {
	SQL struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"sql"`
}

type RestConfig struct {
	Port     string `mapstructure:"port"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
}

type AuthConfig struct {
	Tokens struct {
		Issuer        string `mapstructure:"issuer"`
		AccountAccess struct {
			SecretKey string `mapstructure:"secret_key"`
		} `mapstructure:"account_access"`
	} `mapstructure:"tokens"`
}

type S3Config struct {
	Aws struct {
		BucketName      string `mapstructure:"bucket_name"`
		Region          string `json:"region"`
		AccessKeyID     string `json:"access_key_id"`
		SecretAccessKey string `json:"secret_access_key"`
		SessionToken    string `json:"session_token"`
	} `mapstructure:"aws"`

	Media struct {
		Link struct {
			TTL time.Duration `mapstructure:"ttl"`
		} `mapstructure:"link"`

		Organization struct {
			Icon   awsx.ImageValidator `mapstructure:"icon"`
			Banner awsx.ImageValidator `mapstructure:"banner"`
		} `mapstructure:"organization"`
	} `mapstructure:"media"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Identity string   `mapstructure:"identity"`

	Reader struct {
		Topics struct {
			ProfilesV1 ReaderKafkaConfig `mapstructure:"profiles_v1"`
		} `mapstructure:"profiles"`
	} `mapstructure:"readers"`

	Writer struct {
		Topics struct {
			OrganizationV1 WriterKafkaConfig `mapstructure:"organization_v1"`
			OrgMemberV1    WriterKafkaConfig `mapstructure:"organization_member_v1"`
		} `mapstructure:"topics"`
	} `mapstructure:"write"`

	Inbox KafkaInboxConfig `mapstructure:"inbox"`

	Outbox OutboxConfig `mapstructure:"outbox"`
}

type Config struct {
	Service  ServiceCfg     `mapstructure:"service"`
	Database DatabaseConfig `mapstructure:"database"`
	Rest     RestConfig     `mapstructure:"rest"`
	Log      LogConfig      `mapstructure:"log"`
	Auth     AuthConfig     `mapstructure:"auth"`
	S3       S3Config       `mapstructure:"s3"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
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
