package main

import (
	"context"
	"fmt"
	"go-gin-boilerplate/cmd"
	"go-gin-boilerplate/utils"
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
	
	// !-- database

	// web server
	config := cors.DefaultConfig()
	config.AllowOrigins = flags.AllowOrigins
	config.AddAllowHeaders("csrftoken", "session", "authorization")

	gin.SetMode(utils.StrTernary(bool(flags.Dev), gin.DebugMode, gin.ReleaseMode))
	engine := gin.Default()
	engine.Use(cors.New(config))
	// engine.Use(db.Middleware(core))
	// apis.Setup(engine)

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
