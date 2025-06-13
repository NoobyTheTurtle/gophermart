package config

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	RunAddress           string `mapstructure:"run_address"`
	DatabaseURI          string `mapstructure:"database_uri"`
	AccrualSystemAddress string `mapstructure:"accrual_system_address"`
	JWTSecret            string `mapstructure:"jwt_secret"`
}

func New() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using defaults")
		} else {
			log.Printf("Error reading config file: %v", err)
		}
	} else {
		log.Printf("Using config file: %s", v.ConfigFileUsed())
	}

	setupFlags(v)

	setupEnv(v)

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	if err := validateConfig(&cfg); err != nil {
		log.Fatalf("Config validation failed: %v", err)
	}

	return &cfg
}

func setupFlags(v *viper.Viper) {
	pflag.StringP("address", "a", v.GetString("run_address"), "Server address")
	pflag.StringP("database", "d", v.GetString("database_uri"), "Database URI")
	pflag.StringP("accrual", "r", v.GetString("accrual_system_address"), "Accrual system address")
	pflag.Parse()

	v.BindPFlag("run_address", pflag.Lookup("address"))
	v.BindPFlag("database_uri", pflag.Lookup("database"))
	v.BindPFlag("accrual_system_address", pflag.Lookup("accrual"))
}

func setupEnv(v *viper.Viper) {
	v.BindEnv("run_address", "RUN_ADDRESS")
	v.BindEnv("database_uri", "DATABASE_URI")
	v.BindEnv("accrual_system_address", "ACCRUAL_SYSTEM_ADDRESS")
	v.BindEnv("jwt_secret", "JWT_SECRET")
}

func validateConfig(cfg *Config) error {
	if cfg.RunAddress == "" {
		return fmt.Errorf("run_address is required")
	}
	if cfg.DatabaseURI == "" {
		return fmt.Errorf("database_uri is required")
	}
	if cfg.AccrualSystemAddress == "" {
		return fmt.Errorf("accrual_system_address is required")
	}
	if cfg.JWTSecret == "" {
		return fmt.Errorf("jwt_secret is required")
	}
	return nil
}
