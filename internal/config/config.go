package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sethvargo/go-envconfig"
)

// Configuration is basic structure that contains configuration.
type Configuration struct {
	Log         LogConfig         `env:",prefix=LOG_"`
	Runtime     RuntimeConfig     `env:",prefix=RUNTIME_"`
	Debug       DebugConfig       `env:",prefix=DEBUG_"`
	HTTPServer  HTTPServerConfig  `env:",prefix=HTTP_"`
	Postgres    Postgres          `env:",prefix=PG_"`
	ProgressAPI ProgressAPIConfig `env:",prefix=PROGRESS_API_SERVICE_"`
}

type LogConfig struct {
	Level             string        `env:"LEVEL,default=info"`
	Batch             bool          `env:"BATCH,default=true"`
	BatchSize         int           `env:"BATCH_SIZE,default=1000"`
	BatchPollInterval time.Duration `env:"BATCH_POLL_INTERVAL,default=5s"`
}

type RuntimeConfig struct {
	UseCPUs    int `env:"USE_CPUS,default=0"`
	MaxThreads int `env:"MAX_THREADS,default=0"`
}

type HTTPServerConfig struct {
	CORSEnabled                bool          `env:"CORS_ENABLED,default=false"`
	RequestLoggingEnabled      bool          `env:"REQUEST_LOGGING_ENABLED,default=true"`
	ResponseTimeLoggingEnabled bool          `env:"RESPONSE_TIME_LOGGING_ENABLED,default=false"`
	ReadTimeout                time.Duration `env:"READ_TIMEOUT,default=30s"`
	WriteTimeout               time.Duration `env:"WRITE_TIMEOUT,default=30s"`
	IdleTimeout                time.Duration `env:"IDLE_TIMEOUT,default=30s"`
	MaxRequestBodySize         int           `env:"MAX_REQUEST_BODY_SIZE,default=4194304"`
	Network                    string        `env:"NETWORK,default=tcp"`
	Address                    string        `env:"ADDRESS,default=:8080"`
	AccessTokenSecret          string        `env:"ACCESS_TOKEN_SECRET,default=secret"`
	AccessTokenExpiryHour      int           `env:"ACCESS_TOKEN_EXPIRY_HOUR,default=2"`
}

type Postgres struct {
	DSN             string        `env:"DSN" json:"-"` // Hide in zap logs
	ConnMaxLifetime time.Duration `env:"CONN_MAX_LIFETIME" envDefault:"5m"`
}

type DebugConfig struct {
	InsecureSkipVerify bool `env:"INSECURE_SKIP_VERIFY"`
	WebDebugEnabled    bool `env:"WEB_DEBUG_ENABLED,default=false"`
}

type ProgressAPIConfig struct {
	GrpcURI            string        `env:"URI,default=campaign-progress-api.campaign.svc.cluster.local:18080"`
	IdleTimeout        time.Duration `env:"IDLE_TIMEOUT,default=30s"`
	MaxCallRecvMsgSize int           `env:"MAX_CALL_RECV_MSG_SIZE,default=4194304"`
	MaxCallSendMsgSize int           `env:"MAX_CALL_SEND_MSG_SIZE,default=4194304"`
}

func NewConfig() (*Configuration, error) {
	var envFiles []string

	if _, err := os.Stat(".env"); err == nil {
		log.Println("found .env file, adding it to env config files list")

		envFiles = append(envFiles, ".env")
	}

	if os.Getenv("APP_ENV") != "" {
		appEnvName := fmt.Sprintf(".env.%s", os.Getenv("APP_ENV"))
		if _, err := os.Stat(appEnvName); err == nil {
			log.Println("found", appEnvName, "file, adding it to env config files list")
			envFiles = append(envFiles, appEnvName)
		}
	}

	if len(envFiles) > 0 {
		err := godotenv.Overload(envFiles...)
		if err != nil {
			return nil, errors.Wrapf(err, "error while opening env config: %s", err)
		}
	}

	cfg := &Configuration{} //nolint:exhaustruct
	ctx := context.Background()

	err := envconfig.Process(ctx, cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "error while parsing env config: %s", err)
	}

	return cfg, nil
}
