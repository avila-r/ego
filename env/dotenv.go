package env

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/avila-r/ego/dotenv"
)

var (
	IsLoaded         bool
	HasCustomProfile bool
)

func Get(key string) string {
	var (
		once sync.Once
	)

	if !IsLoaded {
		if err := Load(); err != nil {
			once.Do(func() {
				log.Print("[warning] there is no .env file to load")
			})
		}
	}

	return os.Getenv(key)
}

func Load(path ...string) error {
	if IsLoaded {
		return nil
	}

	var (
		file   = ".env"
		target string
	)

	if len(path) > 0 {
		target = path[0]
	}

	def := filepath.Join(target, file)

	// Loads default .env
	if err := dotenv.Load(def); err != nil {
		return err
	}

	profile := os.Getenv("PROFILE")

	// No custom profile to load
	if profile == "" {
		IsLoaded, HasCustomProfile = true, false

		return nil
	}

	custom := filepath.Join(target, file+"."+profile)

	if err := dotenv.Load(custom); err != nil {
		return err
	}

	IsLoaded, HasCustomProfile = true, true

	return nil
}
