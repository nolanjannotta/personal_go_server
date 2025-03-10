package main

import (
	"fmt"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("healthy")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("HEALTHY"))
}

func snake(w http.ResponseWriter, r *http.Request) {

}

func calculator(w http.ResponseWriter, r *http.Request) {

}

func notFound(w http.ResponseWriter, r *http.Request) {

}

func tui(w http.ResponseWriter, r *http.Request) {

}
