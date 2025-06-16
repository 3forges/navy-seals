package main

import (
	"errors"
	"fmt"
	api "navy-seals/api"
	"os"

	"github.com/gin-gonic/gin"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/joho/godotenv/autoload"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

func main() {

	/**
	 * Command Line start GNU Options parsing with "github.com/jessevdk/go-flags"
	 **/
	var err error
	if _, err = flags.Parse(api.ApiConfig); err != nil {
		var ferr *flags.Error
		if errors.As(err, &ferr) && ferr.Type == flags.ErrHelp {
			os.Exit(0)
		}
		os.Exit(1)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	/**
	 * Albums Endpoints
	 **/
	// router.GET("/albums", api.GetAlbums)
	// router.GET("/albums/:id", api.GetAlbumByID)
	// router.POST("/albums", api.AddAlbum)
	/**
	 * Unseal Keys Endpoints
	 **/
	router.GET("/albums", api.GetAlbums)
	router.GET("/albums/:id", api.GetAlbumByID)
	router.POST("/albums", api.AddAlbum)

	// router.Run("localhost:8765")

	router.Run(fmt.Sprintf("%v:%v", api.ApiConfig.BindAddress, api.ApiConfig.Port))
}
