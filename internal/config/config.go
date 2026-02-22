package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/netbill/organizations-svc/pkg/log"
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
	Port     int `mapstructure:"port"`
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
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		SessionToken    string `mapstructure:"session_token"`
	} `mapstructure:"aws"`

	Media struct {
		Link struct {
			TTL time.Duration `mapstructure:"ttl"`
		} `mapstructure:"link"`

		Resources struct {
			Organization struct {
				Icon struct {
					AllowedFormats []string `mapstructure:"allowed_formats" required:"true"`
					MaxWidth       int      `mapstructure:"max_width" required:"true"`
					MinWidth       int      `mapstructure:"min_width" required:"true"`
					MaxHeight      int      `mapstructure:"max_height" required:"true"`
					MinHeight      int      `mapstructure:"min_height" required:"true"`
					ContentSizeMax int64    `mapstructure:"content_size_max" required:"true"`
				} `mapstructure:"icon"`

				Banner struct {
					AllowedFormats []string `mapstructure:"allowed_formats" required:"true"`
					MaxWidth       int      `mapstructure:"max_width" required:"true"`
					MinWidth       int      `mapstructure:"min_width" required:"true"`
					MaxHeight      int      `mapstructure:"max_height" required:"true"`
					MinHeight      int      `mapstructure:"min_height" required:"true"`
					ContentSizeMax int64    `mapstructure:"content_size_max" required:"true"`
				} `mapstructure:"banner"`
			} `mapstructure:"organization"`
		} `mapstructure:"resources"`
	} `mapstructure:"media"`
}

type KafkaConfig struct {
	Brokers  []string `mapstructure:"brokers"`
	Identity string   `mapstructure:"identity"`

	Consume struct {
		Backoff struct {
			Min time.Duration `mapstructure:"min"`
			Max time.Duration `mapstructure:"max"`
		} `mapstructure:"backoff"`

		Topics struct {
			ProfilesV1 struct {
				Instances      int           `mapstructure:"instances"`
				MinBytes       int           `mapstructure:"min_bytes"`
				MaxBytes       int           `mapstructure:"max_bytes"`
				MaxWait        time.Duration `mapstructure:"max_wait"`
				CommitInterval time.Duration `mapstructure:"commit_interval"`
				QueueCapacity  int           `mapstructure:"queue_capacity"`
			} `mapstructure:"profiles_v1"`
		} `mapstructure:"topics"`
	} `mapstructure:"consume"`

	Produce struct {
		Topics struct {
			OrganizationV1 struct {
				RequiredAcks string        `mapstructure:"required_acks"`
				Compression  string        `mapstructure:"compression"`
				Balancer     string        `mapstructure:"balancer"`
				BatchSize    int           `mapstructure:"batch_size"`
				BatchTimeout time.Duration `mapstructure:"batch_timeout"`
				DialTimeout  time.Duration `mapstructure:"dial_timeout"`
				IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
			} `mapstructure:"organizations_v1"`

			OrgMemberV1 struct {
				RequiredAcks string        `mapstructure:"required_acks"`
				Compression  string        `mapstructure:"compression"`
				Balancer     string        `mapstructure:"balancer"`
				BatchSize    int           `mapstructure:"batch_size"`
				BatchTimeout time.Duration `mapstructure:"batch_timeout"`
				DialTimeout  time.Duration `mapstructure:"dial_timeout"`
				IdleTimeout  time.Duration `mapstructure:"idle_timeout"`
			} `mapstructure:"organization_members_v1"`
		} `mapstructure:"topics"`
	} `mapstructure:"produce"`

	Inbox struct {
		Routines       int           `mapstructure:"routines"`
		Slots          int           `mapstructure:"slots"`
		BatchSize      int           `mapstructure:"batch_size"`
		Sleep          time.Duration `mapstructure:"sleep"`
		MinNextAttempt time.Duration `mapstructure:"min_next_attempt"`
		MaxNextAttempt time.Duration `mapstructure:"max_next_attempt"`
		MaxAttempts    int32         `mapstructure:"max_attempts"`
	} `mapstructure:"inbox"`

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
	Service  ServiceCfg     `mapstructure:"service"`
	Database DatabaseConfig `mapstructure:"database"`
	Rest     RestConfig     `mapstructure:"rest"`
	Log      LogConfig      `mapstructure:"log"`
	Auth     AuthConfig     `mapstructure:"auth"`
	S3       S3Config       `mapstructure:"s3"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
}

func LoadConfig() *Config {
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

	return &config
}

func (cfg *Config) Logger() *log.Logger {
	return log.New(cfg.Log.Level, cfg.Log.Format, cfg.Service.Name)
}

func (cfg *Config) PoolDB(ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, cfg.Database.SQL.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}
