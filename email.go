package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

type Email struct {
	From string `json:"from"`
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

func HandleEmail(w http.ResponseWriter, r *http.Request) {
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

	sendEmail(msg)
}

func sendEmail(msg []byte) {
	godotenv.Load(".env")

	password := os.Getenv("EMAIL_PASSWORD")
	fromAddress := os.Getenv("EMAIL_ADDRESS")
	personalEmail := os.Getenv("PERSONAL_EMAIL")

	// fmt.Println(password)
	// fmt.Println(fromAddress)

	auth := smtp.PlainAuth("", fromAddress, password, "smtp.gmail.com")

	to := []string{personalEmail}
	// msg := []byte("Subject: ***SENT FROM WEBSITE***\r\n" + "\r\n" + "whats good shorty")

	err := smtp.SendMail("smtp.gmail.com:587", auth, fromAddress, to, msg)

	if err != nil {

		log.Fatal(err)

	}

}
