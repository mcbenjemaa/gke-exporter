package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/common/version"

	collector "github.com/mcbenjemaa/gke-metadata-exporter/collector/ip_allocation"
	"github.com/mcbenjemaa/gke-metadata-exporter/pkg/gke"
	"github.com/mcbenjemaa/gke-metadata-exporter/pkg/k8s"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	projectID   = flag.String("project-id", "", "The GCP project id")
	zone        = flag.String("zone", "", "The GCP zone")
	endpoint    = flag.String("endpoint", ":9905", "The endpoint of the HTTP server")
	metricsPath = flag.String("path", "/metrics", "The path on which Prometheus metrics will be served")
)

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	fmt.Fprint(w, "<h2>GKE info Exporter</h2>")
	fmt.Fprint(w, "<ul>")
	fmt.Fprintf(w, "<li><a href=\"%s\">metrics</a></li>", *metricsPath)
	fmt.Fprintf(w, "<li><a href=\"/healthz\">healthz</a></li>")
	fmt.Fprint(w, "</ul>")
}
func main() {
	flag.Parse()

	if projectID == nil || *projectID == "" {
		log.Fatal("ERROR [projectID] is not set: Please set it via argument --project-id")
	}

	// get kubernetes client
	clientset, err := k8s.GetClient() 
	if err != nil {
		log.Fatalf("unable to get k8s client %v", err)
	}

	// create new gke client
	client, err := gke.NewGKEClient(context.Background(), *projectID)
	if err != nil {
		log.Fatalf("unable to create gke client %v", err)
	}
	gke := gke.GKEClient{Client: client,
		Zone:        *zone,
		ProjectID: *projectID}
	

	registry := prometheus.NewRegistry()

	registry.MustRegister(version.NewCollector("gke_info_exporter"))

	registry.MustRegister(collector.NewGKEIpAllocationCollector(gke, clientset))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/healthz", http.HandlerFunc(handleHealthz))
	mux.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	log.Printf("[main] Server starting (%s)", *endpoint)
	log.Printf("[main] metrics served on: %s", *metricsPath)
	log.Fatal(http.ListenAndServe(*endpoint, mux))
}
