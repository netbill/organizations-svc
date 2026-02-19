package boot

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

type RestConfig struct {
	Port     int `mapstructure:"port"`
	Timeouts struct {
		Read       time.Duration `mapstructure:"read"`
		ReadHeader time.Duration `mapstructure:"read_header"`
		Write      time.Duration `mapstructure:"write"`
		Idle       time.Duration `mapstructure:"idle"`
	} `mapstructure:"timeouts"`
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

type AuthConfig struct {
	Tokens struct {
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

	InboxConfig struct {
		Routines       int           `mapstructure:"routines"`
		Slots          int           `mapstructure:"slots"`
		BatchSize      int           `mapstructure:"batch_size"`
		Sleep          time.Duration `mapstructure:"sleep"`
		MinNextAttempt time.Duration `mapstructure:"min_next_attempt"`
		MaxNextAttempt time.Duration `mapstructure:"max_next_attempt"`
		MaxAttempts    int32         `mapstructure:"max_attempts"`
	} `mapstructure:"inbox"`

	Reader struct {
		Topics struct {
			Profiles struct {
				Instances      int           `mapstructure:"instances"`
				MinBytes       int           `mapstructure:"min_bytes"`
				MaxBytes       int           `mapstructure:"max_bytes"`
				MaxWait        time.Duration `mapstructure:"max_wait"`
				CommitInterval time.Duration `mapstructure:"commit_interval"`
				StartOffset    string        `mapstructure:"start_offset"`
				QueueCapacity  int           `mapstructure:"queue_capacity"`
			} `mapstructure:"profiles"`
		} `mapstructure:"profiles"`
	} `mapstructure:"readers"`

	Writer struct {
		RequiredAcks string        `mapstructure:"required_acks"`
		Compression  string        `mapstructure:"compression"`
		Balancer     string        `mapstructure:"balancer"`
		BatchSize    int           `mapstructure:"batch_size"`
		BatchTimeout time.Duration `mapstructure:"batch_timeout"`
		DialTimeout  time.Duration `mapstructure:"dial_timeout"`
		IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
	} `mapstructure:"sender"`

	Outbox struct {
		Routines       int           `mapstructure:"routines"`
		Slots          int           `mapstructure:"slots"`
		BatchSize      int           `mapstructure:"batch_size"`
		Sleep          time.Duration `mapstructure:"sleep"`
		MinNextAttempt time.Duration `mapstructure:"min_next_attempt"`
		MaxNextAttempt time.Duration `mapstructure:"max_next_attempt"`
		MaxAttempts    int32         `mapstructure:"max_attempts"`
	} `mapstructure:"outbox"`
}

type Config struct {
	Service ServiceCfg  `mapstructure:"service"`
	Rest    RestConfig  `mapstructure:"rest"`
	Log     LogConfig   `mapstructure:"log"`
	Auth    AuthConfig  `mapstructure:"auth"`
	S3      S3Config    `mapstructure:"s3"`
	Kafka   KafkaConfig `mapstructure:"kafka"`
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
