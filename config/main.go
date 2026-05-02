package config

import (
	"bot-engine/utils"
	"fmt"
	"strings"

	"github.com/google/wire"
)

type DatabaseSecretsConfig struct {
	MongoURI string
}

// type RedisSecretsConfig struct {
// 	Password string
// }

// type AuthSecretsConfig struct {
// 	InternalServiceToken string
// }

// type FilesSecretsConfig struct {
// 	PublicHMACSecret string
// }

// 1. Added the configuration for our newly created Encryption Service
type BotSecretsConfig struct {
	AESEncryptionKey string
}

type AuthSecretsConfig struct {
	JWTSecret string
}

type SecretsConfig struct {
	Database DatabaseSecretsConfig
	Auth     AuthSecretsConfig
	// Redis    RedisSecretsConfig
	// Auth     AuthSecretsConfig
	// Files    FilesSecretsConfig
	Bot BotSecretsConfig
}

func NewSecretsConfig() (*SecretsConfig, error) {
	// internalToken, err := utils.RequireEnv("INTERNAL_SERVICE_TOKEN")
	// if err != nil {
	// 	return nil, err
	// }

	config := &SecretsConfig{
		Database: DatabaseSecretsConfig{
			MongoURI: utils.GetEnvOrDefault("MONGO_URI", "mongodb://rule-engine-mongo:27017"),
		},
		// Redis: RedisSecretsConfig{
		// 	Password: utils.GetEnvOrDefault("REDIS_PASSWORD", ""),
		// },
		// Auth: AuthSecretsConfig{
		// 	InternalServiceToken: internalToken,
		// },
		// Files: FilesSecretsConfig{
		// 	PublicHMACSecret: utils.GetEnvOrDefault("FILES_PUBLIC_HMAC_SECRET", ""),
		// },
		Auth: AuthSecretsConfig{
			JWTSecret: utils.GetEnvOrDefault("JWT_SECRET", "8baeedf13dcb12dabe8c05cd361df0268"),
		},
		Bot: BotSecretsConfig{
			// 2. Load the 32-byte AES key for encrypting Telegram tokens
			AESEncryptionKey: utils.GetEnvOrDefault("BOT_AES_ENCRYPTION_KEY", "1234567890"),
		},
	}

	if utils.GetEnvOrDefault("GO_ENV", "") == "production" {
		requiredVars := []string{
			"MONGO_URI",
			// "REDIS_PASSWORD",
			// "INTERNAL_SERVICE_TOKEN",
			// "FILES_PUBLIC_HMAC_SECRET",
			"BOT_AES_ENCRYPTION_KEY", // 3. Ensure it's required in production
		}

		var missing []string
		for _, key := range requiredVars {
			if _, err := utils.RequireEnv(key); err != nil {
				missing = append(missing, key)
			}
		}

		if len(missing) > 0 {
			return nil, fmt.Errorf("missing required secret environment variables in production: %s", strings.Join(missing, ", "))
		}
	}

	return config, nil
}

func ProvideBotSecrets(config *SecretsConfig) *BotSecretsConfig {
	return &config.Bot
}

func ProvideDatabaseSecrets(config *SecretsConfig) *DatabaseSecretsConfig {
	return &config.Database
}

func ProvideJWTSecrets(config *SecretsConfig) *AuthSecretsConfig {
	return &config.Auth
}

var ConfigSet = wire.NewSet(
	NewSecretsConfig,
	ProvideBotSecrets,
	ProvideDatabaseSecrets,
	ProvideJWTSecrets,
)
