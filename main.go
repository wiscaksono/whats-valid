package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal/v3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

//go:embed all:frontend/build
var embeddedFrontend embed.FS

type apiServer struct {
	client *whatsmeow.Client
	log    waLog.Logger
}

type CheckResponse struct {
	Number       string `json:"number"`
	IsOnWhatsApp bool   `json:"isOnWhatsApp"`
	Status       string `json:"status"`
}

func newAPIServer(client *whatsmeow.Client, logger waLog.Logger) *apiServer {
	return &apiServer{
		client: client,
		log:    logger,
	}
}

func (s *apiServer) writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		s.client.Log.Errorf("Failed to write JSON response", "error", err)
	}
}

func (s *apiServer) checkNumberHandler(w http.ResponseWriter, r *http.Request) {
	number := r.URL.Query().Get("number")
	if number == "" {
		s.writeJSON(w, http.StatusBadRequest, CheckResponse{
			Status: "error: number parameter is missing",
		})
		return
	}

	if !s.client.IsConnected() {
		s.writeJSON(w, http.StatusServiceUnavailable, CheckResponse{
			Number: number,
			Status: "error: client not connected",
		})
		return
	}

	resp, err := s.client.IsOnWhatsApp([]string{number})
	if err != nil {
		s.writeJSON(w, http.StatusInternalServerError, CheckResponse{
			Number: number,
			Status: fmt.Sprintf("error checking number: %v", err),
		})
		return
	}

	if len(resp) > 0 && resp[0].IsIn {
		s.writeJSON(w, http.StatusOK, CheckResponse{
			Number:       resp[0].Query,
			IsOnWhatsApp: true,
			Status:       "success",
		})
	} else {
		s.writeJSON(w, http.StatusOK, CheckResponse{
			Number:       number,
			IsOnWhatsApp: false,
			Status:       "success",
		})
	}
}

func main() {
	clientLog := waLog.Stdout("Client", "INFO", true)

	ctx := context.Background()
	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New(ctx, "sqlite3", "file:store.db?_foreign_keys=on", dbLog)
	if err != nil {
		clientLog.Errorf("Failed to create SQL container", "error", err)
		os.Exit(1)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		clientLog.Errorf("Failed to get first device", "error", err)
		os.Exit(1)
	}

	client := whatsmeow.NewClient(deviceStore, clientLog)

	if client.Store.ID == nil {
		clientLog.Infof("Starting new login process...")
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			clientLog.Errorf("Failed to connect for QR code generation", "error", err)
			os.Exit(1)
		}
		for evt := range qrChan {
			switch evt.Event {
			case "code":
				client.Log.Infof("⬇️  Scan the QR code below with WhatsApp on your phone ⬇️")
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			case "timeout":
				clientLog.Errorf("QR code scan timed out")
				os.Exit(1)
			case "success":
				clientLog.Infof("Login successful!")
			default:
				clientLog.Infof("Login event", "event", evt.Event)
			}
		}
	} else {
		clientLog.Infof("Attempting to reconnect...")
		if err := client.Connect(); err != nil {
			clientLog.Errorf("Failed to connect WhatsApp client", "error", err)
			os.Exit(1)
		}
	}
	defer client.Disconnect()

	frontendFS, err := fs.Sub(embeddedFrontend, "frontend/build")
	if err != nil {
		clientLog.Errorf("Failed to create frontend file system", "error", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	server := newAPIServer(client, waLog.Stdout("API", "INFO", true))
	mux.Handle("/check", http.HandlerFunc(server.checkNumberHandler))
	fileHandler := http.FileServer(http.FS(frontendFS))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Get the clean path for the requested file.
		filePath := strings.TrimPrefix(r.URL.Path, "/")

		// If the root is requested, serve the main index.html file.
		if filePath == "" {
			filePath = "index.html"
		}

		// Check if the requested file actually exists in the embedded filesystem.
		_, err := frontendFS.Open(filePath)
		if os.IsNotExist(err) {
			// The file does not exist. Redirect the request to the root path.
			// This allows the SPA's router to handle the path.
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		// The file exists, so serve it using the file server.
		fileHandler.ServeHTTP(w, r)
	})

	httpServer := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	go func() {
		clientLog.Infof("🚀 Server starting on port 3000...")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			clientLog.Errorf("HTTP server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	clientLog.Infof("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		clientLog.Errorf("HTTP server shutdown error", "error", err)
	} else {
		clientLog.Infof("✅ Server gracefully stopped")
	}
}
