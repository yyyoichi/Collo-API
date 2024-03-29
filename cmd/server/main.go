package main

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"yyyoichi/Collo-API/internal/api/v3/apiv3connect"
	"yyyoichi/Collo-API/internal/handler"
	logger "yyyoichi/Collo-API/pkg/logger"

	"github.com/google/uuid"
	"github.com/rs/cors"
	"github.com/shogo82148/go-mecab"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	slog.SetDefault(slog.New(logger.NewLogHandler(os.Stdout, nil)))
	handler := getHandler()
	log.Printf("start connectPRC server: %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Panic(err)
	}
}

func getHandler() http.Handler {
	rpc := http.NewServeMux()
	rpc.Handle(apiv3connect.NewMintGreenServiceHandler(&handler.V3Handler{}))

	mux := http.NewServeMux()
	mux.HandleFunc("/", notFoundHandler)
	mux.Handle("/rpc/", http.StripPrefix("/rpc", rpc))
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
	handler := corsHandler.Handler(logginHandler(mux))
	h2cHandler := h2c.NewHandler(handler, &http2.Server{})
	return h2cHandler
}

func logginHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l logger.LogContext
		l.RequestID = uuid.NewString()
		ctx := l.Set(r.Context())
		req := r.Clone(ctx)
		slog.InfoContext(ctx, "request",
			slog.String("host", req.Host),
			slog.String("path", req.URL.Path),
			slog.String("method", req.Method),
		)
		next.ServeHTTP(w, req)
	})
}

//go:embed all:out
var assets embed.FS

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	err := tryRead(r.URL.Path, w)
	if err == nil {
		return
	}
	err = tryRead("404.html", w)
	if err != nil {
		return
	}
}

func tryRead(requestedPath string, w http.ResponseWriter) error {
	reqPath := path.Join("out", requestedPath)
	if reqPath == "out" {
		reqPath = "out/index"
	}
	extension := strings.LastIndex(reqPath, ".")
	if extension == -1 {
		reqPath = fmt.Sprintf("%s.html", reqPath)
	}

	// read file
	f, err := assets.Open(reqPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// dir check
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return errors.ErrUnsupported
	}

	// content type check
	ext := filepath.Ext(requestedPath)
	var contentType string

	if m := mime.TypeByExtension(ext); m != "" {
		contentType = m
	} else {
		contentType = "text/html"
	}

	w.Header().Set("Content-Type", contentType)
	io.Copy(w, f)

	return nil
}

func init() {
	tagger, err := mecab.New(map[string]string{})
	if err != nil {
		panic(err)
	}
	defer tagger.Destroy()
	result, err := tagger.Parse("こんにちは世界")
	if err != nil {
		panic(err)
	}
	for i, s := range strings.Split(result, "\n") {
		fmt.Printf("%d: %s\n\n", i, s)
	}
}
