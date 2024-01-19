package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
	"tolling/client"

	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	var (
		client         = client.NewHttpClient("http://localhost:5050")
		invoiceHandler = NewInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeApiFunc(invoiceHandler.handleGetInvoice))
	logrus.Info("gateway HTTP server running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.URL.Query().Get("obu"))
	if err != nil {
		return err
	}
	invoice, err := h.client.GetInvoice(context.Background(), id)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, invoice)
}

func writeJson(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeApiFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri": r.RequestURI,
			}).Info("REQ")
		}(time.Now()) 
		if err := fn(w, r); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
