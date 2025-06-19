package main

import (
	"fmt"
	api "navy-seals/api"
	"navy-seals/config"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

var (
	version = "master"
	commit  = "latest"
	date    = "-"
)

func main() {

	config.LoadConfig()
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	/**
	 * Albums Endpoints
	 **/
	// router.GET("/albums", api.GetAlbums)
	// router.GET("/albums/:id", api.GetAlbumByID)
	// router.POST("/albums", api.AddAlbum)
	router.GET("/vault-status", api.GetVaultStatus)
	router.POST("/vault-init", api.InitVault)
	router.POST("/vault-unseal", api.UnsealVault)
	/**
	 * Unseal Keys Endpoints
	 **/
	router.GET("/albums", api.GetAlbums)
	router.GET("/albums/:id", api.GetAlbumByID)
	router.POST("/albums", api.AddAlbum)

	// router.Run("192.168.1.12:8751")
	// router.Run()
	var listen_on string = fmt.Sprintf("%v:%v", config.ApiConfig.BindAddress, config.ApiConfig.Port)
	fmt.Printf(" Welcome to navy seals ")
	fmt.Printf(" listen_on [%v]", listen_on)

	fmt.Printf(" Navy seals VAULT ADDRESS IS [%v]", config.ApiConfig.VaultAddress)
	fmt.Printf(" Navy seals BIND ADDRESS IS [%v]", config.ApiConfig.BindAddress)
	fmt.Printf(" Navy seals PORT IS [%v]", config.ApiConfig.Port)
	// router.Run(listen_on)
	//router.Run("0.0.0.0:8751")
	router.RunTLS("0.0.0.0:8751", "./navyseals.pesto.io.pem", "./navyseals.pesto.io-key.pem")

}
