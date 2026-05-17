package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	Providers ProvidersConfig
	Logger    LoggerConfig
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type DatabaseConfig struct {
	Type   string        `mapstructure:"type"`
	SQLite SQLiteConfig  `mapstructure:"sqlite"`
}

type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

type ProvidersConfig struct {
	OpenAI OpenAIConfig  `mapstructure:"openai"`
	Claude ClaudeConfig  `mapstructure:"claude"`
	Ollama OllamaConfig  `mapstructure:"ollama"`
}

type OpenAIConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	APIKey  string `mapstructure:"api_key"`
	BaseURL string `mapstructure:"base_url"`
}

type ClaudeConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	APIKey  string `mapstructure:"api_key"`
}

type OllamaConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	BaseURL string `mapstructure:"base_url"`
}

type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

func Load() (*Config, error) {
	// 加载.env文件
	godotenv.Load()

	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AddConfigPath(".")

	// 设置默认值
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.sqlite.path", "./data/app.db")
	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")

	// 绑定环境变量
	v.BindEnv("providers.openai.api_key", "OPENAI_API_KEY")
	v.BindEnv("providers.claude.api_key", "CLAUDE_API_KEY")
	v.BindEnv("providers.ollama.base_url", "OLLAMA_BASE_URL")

	err := v.ReadInConfig()
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 从环境变量覆盖
	if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
		cfg.Providers.OpenAI.APIKey = apiKey
	}
	if apiKey := os.Getenv("CLAUDE_API_KEY"); apiKey != "" {
		cfg.Providers.Claude.APIKey = apiKey
	}
	if baseURL := os.Getenv("OLLAMA_BASE_URL"); baseURL != "" {
		cfg.Providers.Ollama.BaseURL = baseURL
	}

	return &cfg, nil
}
