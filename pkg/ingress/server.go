package ingress

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-http-utils/logger"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

//RunHTTP starts the webserver to listen for inbound webhooks from external services
func RunHTTP(port int, natsEndpoint string) {
	connectNats(natsEndpoint)

	r := mux.NewRouter()
	r.HandleFunc("/{stream}/{subject}", HandlePub).Methods("GET")

	r.Use(cors.AllowAll().Handler, ValidateAuth)

	handler := logger.Handler(r, os.Stdout, logger.CommonLoggerType)

	listeningEndpoint := fmt.Sprintf(":%d", port)
	log.Printf("Starting HTTP server (port %d)\n", port)

	srv := &http.Server{
		Addr:         listeningEndpoint,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      handler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("shutting down...")
	os.Exit(0)
}
