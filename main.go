package main

import (
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
)

func Start(s *http.Server) {
	log.Info("Starting HTTP server", "host", "localhost", "port", 8080)

	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Could not start server", "error", err)
	}

}

func newServer() *http.Server {

	mux := http.NewServeMux()

	// var staticFS, err = fs.Sub(static, "static")

	// if err != nil {
	// 	panic(err)
	// }

	// fs := http.FileServer(http.FS(staticFS))

	// mux.Handle("/", fs)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	mux.HandleFunc("GET /health", health)

	// mux.HandleFunc("POST /email", handleEmail)
	// mux.HandleFunc("GET /{page}", subPageHandler)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return s

}

func main() {
	srv := newServer()
	Start(srv)
}

// type Email struct {
// 	From string `json:"from"`
// 	Name string `json:"name"`
// 	Msg  string `json:"msg"`
// }

// func handleEmail(w http.ResponseWriter, r *http.Request) {
// 	var email Email

// 	decoder := json.NewDecoder(r.Body)
// 	decoder.DisallowUnknownFields()

// 	err := decoder.Decode(&email)

// 	if err != nil {
// 		fmt.Println("Error", err)
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return

// 	}
// 	_, err = mail.ParseAddress(email.From)

// 	if err != nil || email.Msg == "" {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}
// 	msg := []byte("Subject: ***SENT FROM WEBSITE***\r\n" + "\r\n" + "SENDER: " + email.From + "\r\n\n\n\n" + "NAME: " + email.Name + "\r\n\n\n\n" + "MESSAGE: " + email.Msg)

// 	password := os.Getenv("EMAIL_PASSWORD")
// 	fromAddress := os.Getenv("EMAIL_ADDRESS")
// 	personalEmail := os.Getenv("PERSONAL_EMAIL")

// 	auth := smtp.PlainAuth("", fromAddress, password, "smtp.gmail.com")

// 	to := []string{personalEmail}

// 	err = smtp.SendMail("smtp.gmail.com:587", auth, fromAddress, to, msg)

// 	if err != nil {

// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return

// 	}
// }
