package httpServer

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"
	"strconv"

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
	mux.HandleFunc("GET /ipdistance", handleIpAddressDistance)
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

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type Response struct {
	Miles string `json:"miles"`
}
type ClientIP struct {
	IP string `json:"ip"`
}

func checkAndReturnServerError(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleIpAddressDistance(w http.ResponseWriter, r *http.Request) {
	IPDistanceURL := "http://ip-api.com/json/"
	laLat, laLon := 34.052235, -118.243683 // latitude and longitude of Los Angeles

	// get clients public IPv4 address
	var clientIP ClientIP
	cIP, err := http.Get("https://api.ipify.org?format=json")
	checkAndReturnServerError(w, err)
	err = json.NewDecoder(cIP.Body).Decode(&clientIP)
	checkAndReturnServerError(w, err)

	// latitude and longitude of clients IP address
	locationResp, err := http.Get(fmt.Sprint(IPDistanceURL, clientIP.IP))
	checkAndReturnServerError(w, err)
	var location Location
	err = json.NewDecoder(locationResp.Body).Decode(&location)
	checkAndReturnServerError(w, err)

	// calculate distance between lats and lons
	distance := getDistance(laLat, laLon, location.Lat, location.Lon)

	response := Response{Miles: strconv.Itoa(int(distance * 0.621371))}

	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func getDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert degrees to radians
	lat1 = lat1 * math.Pi / 180
	lon1 = lon1 * math.Pi / 180
	lat2 = lat2 * math.Pi / 180
	lon2 = lon2 * math.Pi / 180

	// Haversine formula to calculate the great-circle distance
	distance := math.Acos(math.Sin(lat1)*math.Sin(lat2)+math.Cos(lat1)*math.Cos(lat2)*math.Cos(lon2-lon1)) * 6371

	return distance
}
