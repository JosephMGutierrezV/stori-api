package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	S3BucketName   string `mapstructure:"S3_BUCKET_NAME"`
	S3Region       string `mapstructure:"S3_REGION"`
	AWSEndpointURL string `mapstructure:"AWS_ENDPOINT_URL"`
	UsePathStyle   bool   `mapstructure:"AWS_S3_USE_PATH_STYLE"`
}

func LoadConfig() (*Config, error) {
	viper.SetDefault("AWS_S3_USE_PATH_STYLE", false)

	for _, k := range []string{

		"S3_BUCKET_NAME", "S3_REGION",
		"AWS_ENDPOINT_URL", "AWS_S3_USE_PATH_STYLE",
	} {
		_ = viper.BindEnv(k)
	}

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		_ = viper.ReadInConfig()
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	missing := []string{}
	req := func(k, v string) {
		if strings.TrimSpace(v) == "" {
			missing = append(missing, k)
		}
	}

	req("S3_BUCKET_NAME", cfg.S3BucketName)
	req("S3_REGION", cfg.S3Region)

	if len(missing) > 0 {
		return nil, fmt.Errorf("faltan variables: %v", missing)
	}

	return &cfg, nil
}
