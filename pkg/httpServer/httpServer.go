package httpServer

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/charmbracelet/log"
)

func Start(s *http.Server, done chan<- os.Signal) {
	log.Info("Starting HTTP server", "host", "localhost", "port", 8080)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Could not start server", "error", err)
		done <- nil
	}

}

//go:embed static
var static embed.FS

func SetUp() *http.Server {

	mux := http.NewServeMux()

	var staticFS, err = fs.Sub(static, "static")

	if err != nil {
		panic(err)
	}

	fs := http.FileServer(http.FS(staticFS))

	mux.Handle("/", fs)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("healthy")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("HEALTHY"))
	})
	mux.HandleFunc("POST /email", handleEmail)
	// mux.HandleFunc("GET /{page}", subPageHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// ReadTimeout:    10 * time.Second,
		// WriteTimeout:   10 * time.Second,
		// MaxHeaderBytes: 1 << 20,
	}

	return s

}

// func homePageHandler(w http.ResponseWriter, r *http.Request) {

// }

// func subPageHandler(w http.ResponseWriter, r *http.Request) {
// 	page := r.PathValue("page")

// }

type Email struct {
	From string `json:"from"`
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func handleEmail(w http.ResponseWriter, r *http.Request) {
	var email Email

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&email)

	if err != nil {
		fmt.Println("Error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
	_, err = mail.ParseAddress(email.From)

	if err != nil || email.Msg == "" {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg := []byte("Subject: ***SENT FROM WEBSITE***\r\n" + "\r\n" + "SENDER: " + email.From + "\r\n\n\n\n" + "NAME: " + email.Name + "\r\n\n\n\n" + "MESSAGE: " + email.Msg)

	password := os.Getenv("EMAIL_PASSWORD")
	fromAddress := os.Getenv("EMAIL_ADDRESS")
	personalEmail := os.Getenv("PERSONAL_EMAIL")

	auth := smtp.PlainAuth("", fromAddress, password, "smtp.gmail.com")

	to := []string{personalEmail}

	err = smtp.SendMail("smtp.gmail.com:587", auth, fromAddress, to, msg)

	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}
}
