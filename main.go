package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	rabbitmqapiwrapper "rabbitmq-exporter/rabbitmq_api_wrapper"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Struct describing the metrics
type RabbitMqMetrics struct {
	AllMsg   *prometheus.GaugeVec
	ReadyMsg *prometheus.GaugeVec
	UnackMsg *prometheus.GaugeVec
}

var mu sync.Mutex

var RABBITMQ_HOST string
var RABBITMQ_USER string
var RABBITMQ_PASSWORD string

var rabbitMqMetric *RabbitMqMetrics

func main() {
	loadEnvVariables()

	reg := prometheus.NewRegistry()

	rabbitMqMetric = createNewMetric(reg)

	// Start a ticker in a goroutine
	go startTicker()

	// Create a server and export /metrics endpoint to prometheus
	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	http.Handle("/metrics", promHandler)

	fmt.Println("Starting server on port 9051...")
	log.Fatal(http.ListenAndServe(":9051", nil))
}

/**
* Starts a ticker, queries RabbitMQ API then updates the metrics
 */
func startTicker() {
	ticker := time.NewTicker(time.Second * time.Duration(2)) // every 2 seconds

	for {
		select {
		case <-ticker.C:
			ticker.Reset(time.Second * time.Duration(2))
			getMetrics(rabbitMqMetric)
		default:
			continue
		}
	}
}

/**
* loads .env file
 */
func loadEnvVariables() {
	godotenv.Load(".env")
	RABBITMQ_PASSWORD = os.Getenv("RABBITMQ_PASSWORD")
	RABBITMQ_HOST = os.Getenv("RABBITMQ_HOST")
	RABBITMQ_USER = os.Getenv("RABBITMQ_USER")

	// if any of the variables are not set, exit
	if RABBITMQ_USER == "" {
		log.Fatal("RABBITMQ_USER env variable is not provided")
	}

	if RABBITMQ_HOST == "" {
		log.Fatal("RABBITMQ_HOST env variable is not provided")
	}

	if RABBITMQ_PASSWORD == "" {
		log.Fatal("RABBITMQ_PASSWORD env variable is not provided")
	}
}

/*
* Creates and registeres a new metric
 */
func createNewMetric(reg prometheus.Registerer) (newMet *RabbitMqMetrics) {
	metr := &RabbitMqMetrics{
		AllMsg: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "rabbitmq_individual_queue_messages",
				Help: "Total messages in this queue",
			},
			[]string{"host", "vhost", "name"},
		),
		ReadyMsg: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "rabbitmq_individual_queue_messages_ready",
				Help: "Ready Messages fo rspecified queue",
			},
			[]string{"host", "vhost", "name"},
		),
		UnackMsg: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "rabbitmq_individual_queue_messages_unacknowledged",
				Help: "Unacknowledged messages",
			},
			[]string{"host", "vhost", "name"},
		),
	}

	reg.MustRegister(metr.AllMsg)
	reg.MustRegister(metr.ReadyMsg)
	reg.MustRegister(metr.UnackMsg)

	return metr
}

/*
* Get metrics from RabbitMQ API and update the RabbitMqMetrics variable in a concurrency safe way
 */
func getMetrics(m *RabbitMqMetrics) {
	defer mu.Unlock()
	err, metrics := rabbitmqapiwrapper.GetAllQueueMetrics(RABBITMQ_HOST)

	if err != nil {
		log.Println("Unable to get metric")
		return
	}

	mu.Lock()
	for _, metric := range metrics {
		m.AllMsg.WithLabelValues(RABBITMQ_HOST, metric.Vhost, metric.Name).Set(float64(metric.Messages))
		m.ReadyMsg.WithLabelValues(RABBITMQ_HOST, metric.Vhost, metric.Name).Set(float64(metric.MessagesReady))
		m.UnackMsg.WithLabelValues(RABBITMQ_HOST, "/", metric.Name).Set(float64(metric.MessagesUnack))
	}
}
