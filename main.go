package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"flag"
	sloghttp "github.com/samber/slog-http"
	"google.golang.org/api/option"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var fGoogleProjectID = flag.String("google-project-id", "", "Google project id")
var fGooglePubsubTopic = flag.String("google-pubsub-topic", "", "Google PubSub topic")
var fGoogleCredentialsFile = flag.String("google-service-account", "", "Google service account file")
var fCertificatePublic = flag.String("https-cert", "", "HTTP certificate public key")
var fCertificateKey = flag.String("https-key", "", "HTTP certificate private key")
var fBindAddr = flag.String("bind", ":8000", "TCP bind address")

func handler(log *slog.Logger, topic *pubsub.Topic) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rb, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		msg := &pubsub.Message{Data: rb}

		for hk, kv := range r.Header {
			switch strings.ToLower(hk) {
			case "content-type":
			default:
				continue
			}
			msg.Attributes[hk] = kv[0]
		}

		messageID, err := topic.Publish(r.Context(), msg).Get(r.Context())
		if err != nil {
			log.Error("Unable to get pubsub message id", "error", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		log.Info("published message", "id", messageID)
		_, _ = io.WriteString(w, messageID)
	})
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	nl, err := net.Listen("tcp", *fBindAddr)
	if err != nil {
		log.Error("Unable to bind to supplied address", "addr", *fBindAddr, "error", err)
		os.Exit(1)
	}

	psInitContext, psInitCancel := context.WithTimeout(ctx, 30*time.Second)
	var pubsubClientOpts []option.ClientOption
	if *fGoogleCredentialsFile != "" {
		pubsubClientOpts = append(pubsubClientOpts, option.WithCredentialsFile(*fGoogleCredentialsFile))
	}
	psClient, err := pubsub.NewClient(psInitContext, *fGoogleProjectID, pubsubClientOpts...)
	if err != nil {
		log.Error("Unable to create pubsub client", "error", err)
		os.Exit(1)
	}

	psTopic := psClient.Topic(*fGooglePubsubTopic)
	if ok, err := psTopic.Exists(psInitContext); err != nil {
		log.Error("Unable to check if pubsub topic exists", "error", err)
		os.Exit(1)
	} else if !ok {
		log.Error("Pubsub topic doesn't exist", "topic", *fGooglePubsubTopic)
		os.Exit(1)
	}
	psInitCancel()

	mux := http.NewServeMux()
	mux.Handle("POST /publish", handler(log, psTopic))

	srv := &http.Server{
		Handler:           sloghttp.New(log)(mux),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       time.Minute,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	var serveErr error
	if *fCertificatePublic != "" || *fCertificateKey != "" {
		serveErr = srv.ServeTLS(nl, *fCertificatePublic, *fCertificateKey)
	} else {
		log.Warn("Running without HTTPS")
		serveErr = srv.Serve(nl)
	}

	if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
		log.Error("Unable to start server", "error", serveErr)
		os.Exit(1)
	}
}
