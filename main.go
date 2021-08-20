package main

import (
	"context"
	"fmt"
	"go-gin-boilerplate/cmd"
	"go-gin-boilerplate/db"
	srv "go-gin-boilerplate/server"
	"go-gin-boilerplate/utils"
	"go-gin-boilerplate/web"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	backgroundCtx := context.Background()
	flags := cmd.InitFlags()

	// database
	core := new(db.Core)
	core.DB = db.OpenDB("mysql", flags.DSN)
	defer core.DB.Close()
	dbCtxRoot, stop := context.WithCancel(backgroundCtx)
	defer stop()
	core.Ctx = &dbCtxRoot
	if err := core.Ping(3 * time.Second); err != nil {
		log.Fatalf("%v", err)
	}
	log.Println("db contexts are ready for roll")
	// !-- database

	// web server
	config := cors.DefaultConfig()
	config.AllowHeaders = flags.AllowHeaders
	config.AllowOrigins = flags.AllowOrigins

	gin.SetMode(utils.StrTernary(bool(flags.Dev), gin.DebugMode, gin.ReleaseMode))
	engine := gin.Default()
	engine.Use(cors.New(config))
	engine.Use(srv.Middleware(core))
	web.Setup(engine)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", flags.Port),
		Handler: engine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen and Serve: %v\n", err)
		}
	}()
	// !-- web server

	// graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(backgroundCtx, 5*time.Second)
	defer cancel()

	log.Println("Server Existing In 5 Seconds.")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server Shutdown: %v\n", err)
	}
	// !-- graceful Shutdown
}
