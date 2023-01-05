package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/arylatt/go-monzo"
)

func main() {
	srv := &http.Server{
		Addr: ":54093",
		Handler: http.HandlerFunc(monzo.WebhookPayloadHandler(func(rw http.ResponseWriter, r *http.Request, payload *monzo.WebhookPayload) {
			rw.WriteHeader(http.StatusOK)

			fmt.Printf("Received request: %s %s\n", r.Method, r.RequestURI)

			data, err := json.MarshalIndent(payload, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
				return
			}

			fmt.Fprintf(os.Stdout, "%s\n\n", data)
		})),
	}

	fmt.Fprint(os.Stderr, srv.ListenAndServe().Error())
}
