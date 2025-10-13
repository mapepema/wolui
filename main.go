package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

const (
	targetMAC = "00:11:22:33:44:55"
)

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wake", handleWake)

	fmt.Println("server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}

func handleWake(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	// Run the wol command
	cmd := exec.Command("wakeonlan", targetMAC)
	err := cmd.Run()
	if err != nil {
		log.Printf("error executing wol: %v", err)
		http.Error(w, "Failed to send WOL", 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("WOL packet sent"))
}
