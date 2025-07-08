package repositories

import (
	"log"
	"os"
)

// Config はアプリケーション全体の設定を保持します。
type Config struct {
	Env              string // "development", "staging", "production"など
	CollectionPrefix string
}

// Cfg はアプリケーションで唯一の設定インスタンスです。
// main関数で一度だけ初期化されます。
var Cfg *Config

func LoadConfig() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development" // デフォルトは開発環境とする
		log.Println("APP_ENV is not set. Defaulting to 'development'")
	}

	var prefix string
	switch env {
	case "production":
		prefix = "" // 本番環境ではプレフィックスなし
	case "staging":
		prefix = "stg_"
	case "test":
		prefix = "test_"
	case "development":
		fallthrough // 開発環境もプレフィックスあり
	default:
		prefix = "dev_"
	}

	Cfg = &Config{
		Env:              env,
		CollectionPrefix: prefix,
	}

	log.Printf("Application running in '%s' environment. Collection prefix: '%s'", Cfg.Env, Cfg.CollectionPrefix)
}

// GetCollectionName はベース名に環境に応じたプレフィックスを付与して返します。
// アプリケーション内でコレクション名を取得する際は、必ずこの関数を使用します。
func GetCollectionName(baseName string) string {
	if Cfg == nil {
		// Load() が呼ばれる前にこの関数が使われた場合のエラーハンドリング
		log.Fatal("Configuration is not loaded. Call config.Load() at startup.")
	}
	return Cfg.CollectionPrefix + baseName
}
