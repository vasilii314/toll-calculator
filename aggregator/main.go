package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"tolling/types"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	store := NewMemoryStore()
	var (
		svc Aggregator = NewInvoiceAggregator(store)
	)
	svc = NewMetricsMiddleware(NewLogMiddleware(svc))
	go func() {
		log.Fatal(makeGrpcTransport(":5051", svc))
	}()
	log.Fatal(makeHttpTransport(":5050", svc))
}

func makeHttpTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(listenAddr, nil)
}

func makeGrpcTransport(listenAddr string, svc Aggregator) error {
	// make a tcp listener
	fmt.Println("GRPC transport running on port", listenAddr)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer lis.Close()
	// make a new GRPC native server with options
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// register custom GRPC server implementation
	types.RegisterAggregatorServer(server, NewAggregatorGrpcServer(svc))
	return server.Serve(lis)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJson(
				w,
				http.StatusBadRequest,
				map[string]string{
					"error": err.Error(),
				},
			)
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJson(
				w,
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
			return
		}
	}
}

func writeJson(rw http.ResponseWriter, status int, v any) error {
	rw.WriteHeader(status)
	rw.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "missing obu id"})
			return
		}
		obuId, err := strconv.Atoi(values[0])
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid obu id"})
			return
		}
		invoice, err := svc.CalcualteInvoice(obuId)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, http.StatusOK, invoice)
	}
}
