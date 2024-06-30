package main

import (
    "log"
    "net/http"
)

func main() {
	port := `:8080`
    server, err := NewServer()
    if err != nil {
        log.Fatalf("Server run error >>>>>>> %v", err)
    }

    log.Println("Server running and listening on port", port)
    if err := http.ListenAndServe(port, server); err != nil {
        log.Fatalf("Server listen failed error >>>>>>>> %v", err)
    }
}