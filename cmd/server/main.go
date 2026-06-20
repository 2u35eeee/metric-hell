package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"metric-hell/internal/api"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := api.MustNewDefaultHandler()
	addr := ":" + port
	fmt.Printf("WorkflowBench / Metric Hell running at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(addr, handler))
}
