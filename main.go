package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"yyyoichi/Collo-API/gen/proto/collo/v1/collov1connect"

	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	certPath := os.Getenv("CERT_PATH")
	keyPath := os.Getenv("KEY_PATH")

	log.Printf("start gPRC server: %s", port)
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%s", port), certPath, keyPath, getHandler()); err != nil {
		log.Panic(err)
	}
}

func getHandler() http.Handler {
	svc := &ColloServer{}
	mux := http.NewServeMux()
	mux.Handle(collov1connect.NewColloServiceHandler(svc))
	corsHandler := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodPost,
		},
		AllowedOrigins: []string{os.Getenv("CLIENT_HOST")},
		AllowedHeaders: []string{
			"Accept-Encoding",
			"Content-Encoding",
			"Content-Type",
			"Connect-Protocol-Version",
			"Connect-Timeout-Ms",
		},
		ExposedHeaders: []string{},
	})
	handler := corsHandler.Handler(mux)
	h2cHandler := h2c.NewHandler(handler, &http2.Server{})
	return h2cHandler
}
