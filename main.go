package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

// targetMAC is the machine to wake, read from WOL_MAC at startup so the
// address of a real device never ends up in the source code.
var targetMAC string

func main() {
	mac, err := net.ParseMAC(os.Getenv("WOL_MAC"))
	if err != nil {
		log.Fatal("set WOL_MAC to the MAC address of the machine to wake, e.g. WOL_MAC=00:11:22:33:44:55")
	}
	targetMAC = mac.String()

	addr := os.Getenv("WOL_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wake", handleWake)

	fmt.Printf("server started on %s, waking %s\n", addr, targetMAC)
	log.Fatal(http.ListenAndServe(addr, nil))
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
