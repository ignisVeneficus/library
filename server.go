package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go.mau.fi/zeroconfig"
	"gopkg.in/yaml.v3"

	"github.com/ignisVeneficus/library/api"
	"github.com/ignisVeneficus/library/config"
	"github.com/ignisVeneficus/library/db"
	"github.com/ignisVeneficus/library/db/dao"
	"github.com/ignisVeneficus/library/file"
	"github.com/ignisVeneficus/library/scraper"
	"github.com/ignisVeneficus/library/status"
)

var (
	forceUpdate bool
	importFile  string
	exportFile  string
	noServer    bool
	noUpdate    bool
)

func init() {
	loadLogging("log.config")
}
func loadLogging(logConfig string) error {
	var f *os.File
	f, err := os.Open(logConfig)
	if err != nil {
		return err
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	var cfg zeroconfig.Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return err
	}
	logger, err := cfg.Compile()
	if err != nil {
		return err
	}
	log.Logger = *logger
	return nil
}

func handleFile() {
	status := status.GetStatus()
	base := config.GetFilesystemConfig().BookSource
	cover := config.GetFilesystemConfig().CoverSource
	scraper.Parse(status, base, cover, forceUpdate)
}
func help() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()

}
func main() {
	flag.Usage = help

	flag.BoolVar(&forceUpdate, "fu", false, "Force to update images from ebooks")
	flag.BoolVar(&forceUpdate, "forceUpdate", false, "Force to update images from ebooks")
	flag.BoolVar(&noServer, "ns", false, "Not start the server")
	flag.BoolVar(&noServer, "noServer", false, "Not start the server")
	flag.BoolVar(&noUpdate, "nc", false, "Not initial file checking")
	flag.BoolVar(&noUpdate, "noCheck", false, "Not initial file checking")
	flag.StringVar(&exportFile, "export", "", "Export the database into JSON file")
	flag.StringVar(&exportFile, "e", "", "Export the database into JSON file")

	flag.Parse()

	database := db.GetDatabaseMulti()
	defer database.Close()
	ctx := context.Background()
	if err := dao.CreateDatabase(database, ctx); err != nil {
		log.Logger.Fatal().Err(err).Msg("Database Error")
		log.Logger.Info().Msg("Stopping")
		panic(err)
	}

	if exportFile != "" {
		log.Logger.Info().Str("Filename", exportFile).Msg("Export to file")
		if err := file.WriteAllBookToFile(exportFile); err != nil {
			log.Logger.Fatal().Err(err).Msg("Export error")
		}
	}

	log.Logger.Info().Bool("Check files", !noUpdate).Msg("Files")

	var wg sync.WaitGroup
	if !noUpdate {
		log.Logger.Info().Msg("File checking start")
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleFile()
		}()
	}
	log.Logger.Info().Bool("Start server", !noServer).Msg("Server")
	if !noServer {
		log.Logger.Info().Msg("Server starting")
		router := gin.Default()
		router.Use(static.Serve("/", static.LocalFile("./web", false)))
		router.Use(static.Serve("/cover", static.LocalFile(config.GetFilesystemConfig().CoverSource, false)))
		router.Use(static.Serve("/books", static.LocalFile(config.GetFilesystemConfig().BookSource, false)))
		router.GET("/api/book", api.GetAllBook)
		router.GET("/api/book/:id", api.GetBook)
		router.POST("/api/book", api.PostBook)
		router.GET("/api/author", api.GetAllAuthor)
		router.GET("/api/series", api.GetAllSeries)

		router.GET("/api/scraper/", api.Scrape)
		router.GET("/api/export", api.DownloadAllBook)

		// Listen and serve on 0.0.0.0:8080
		router.Run(":8888")
	}
	wg.Wait()
}
