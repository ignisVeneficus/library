package config

import (
	"errors"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	ENV_PREFIX   = "LIBRARY"
	ENV_USERNAME = ENV_PREFIX + "_DB_USERNAME"
	ENV_PASSWORD = ENV_PREFIX + "_DB_PASSWORD"
	ENV_URL      = ENV_PREFIX + "_DB_HOST"
	ENV_DATABASE = ENV_PREFIX + "_DB_DATABASE"
)

type DatabaseConfig struct {
	Url      string
	User     string
	Pass     string
	Database string
}
type FileSystemConfig struct {
	BookSource  string
	CoverSource string
}

var (
	databaseConfig   DatabaseConfig
	filesystemConfig FileSystemConfig
)

func init() {
	databaseConfig = initDatabase()
	filesystemConfig = initFilesystem()
}

func initDatabase() DatabaseConfig {
	log.Logger.Debug().Msg("Load Database config Start")
	userName := os.Getenv(ENV_USERNAME)
	password := os.Getenv(ENV_PASSWORD)
	url := os.Getenv(ENV_URL)
	database := os.Getenv(ENV_DATABASE)

	if len(strings.TrimSpace(userName)) == 0 {
		log.Logger.Error().Msg("No username for database configurated! (" + ENV_USERNAME + ")")
		panic(errors.New("missing username config item"))
	} else {
		log.Logger.Info().Str("Database username", userName).Msg("Config item loaded")
	}
	if len(strings.TrimSpace(password)) == 0 {
		log.Logger.Error().Msg("No password for database configurated!  (" + ENV_PASSWORD + ")")
		panic(errors.New("missing password config item"))
	} else {
		log.Logger.Info().Str("Database password", "****").Msg("Config item loaded")
	}
	if len(strings.TrimSpace(url)) == 0 {
		log.Logger.Error().Msg("No url for database configurated! (" + ENV_URL + ")")
		panic(errors.New("missing database url config item"))
	} else {
		log.Logger.Info().Str("Database url", url).Msg("Config item loaded")
	}
	if len(strings.TrimSpace(database)) == 0 {
		log.Logger.Error().Msg("No database for database configurated! (" + ENV_DATABASE + ")")
		panic(errors.New("missing database config item"))
	} else {
		log.Logger.Info().Str("Database database", database).Msg("Config item loaded")
	}
	log.Logger.Info().Msg("Load Database config finished")
	return DatabaseConfig{User: userName, Url: url, Pass: password, Database: database}
}

func initFilesystem() FileSystemConfig {
	log.Logger.Debug().Msg("Load Filesystem config")
	bookDir := os.Getenv(ENV_PREFIX + "_BOOKS")
	coverDir := os.Getenv(ENV_PREFIX + "_COVERS")

	if len(strings.TrimSpace(coverDir)) == 0 {
		log.Logger.Warn().Msg("No directory for covers configurated!")
		panic(errors.New("missing cover directory config item"))
	} else {
		log.Logger.Info().Str("Cover folder", coverDir).Msg("Config item loaded")
	}
	if len(strings.TrimSpace(bookDir)) == 0 {
		log.Logger.Error().Msg("No directory for covers configurated!")
		panic(errors.New("missing book directory config item"))
	} else {
		log.Logger.Info().Str("Books folder", bookDir).Msg("Config item loaded")
	}
	log.Logger.Info().Msg("Load Filesystem config finished")

	return FileSystemConfig{BookSource: bookDir, CoverSource: coverDir}
}

func GetDatabaseConfig() DatabaseConfig {
	return databaseConfig
}
func GetFilesystemConfig() FileSystemConfig {
	return filesystemConfig
}
