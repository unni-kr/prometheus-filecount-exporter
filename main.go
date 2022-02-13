package main

import (
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type fileCountCollector struct {
	fileCountMetric *prometheus.Desc
}

func newFileCountCollector() *fileCountCollector {
	return &fileCountCollector{
		fileCountMetric: prometheus.NewDesc("ls_metric",
			"Shows the count of files the given path",
			nil, nil),
	}
}

func (collector *fileCountCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.fileCountMetric
}

func (collector *fileCountCollector) Collect(ch chan<- prometheus.Metric) {
	path := "."
	fileCountResult, _ := checkFileCount(path)
	fileCountResultString := strings.TrimSpace(string(fileCountResult))
	metricValue, err := strconv.ParseFloat(fileCountResultString, 64)
	if err != nil {
		log.Fatal(err)
	}
	m1 := prometheus.MustNewConstMetric(collector.fileCountMetric, prometheus.GaugeValue, metricValue)
	m1 = prometheus.NewMetricWithTimestamp(time.Now(), m1)
	ch <- m1
}

func checkFileCount(path string) (result []byte, err error) {
	cmd1 := exec.Command("ls", path)
	cmd2 := exec.Command("wc", "-l")
	outPipe, err := cmd1.StdoutPipe()
	if err != nil {
		return nil, err
	}
	defer outPipe.Close()
	cmd1.Start()
	cmd2.Stdin = outPipe
	out, err := cmd2.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func main() {

	mux := http.NewServeMux()

	fileCount := newFileCountCollector()
	prometheus.MustRegister(fileCount)

	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Listening...")
	http.ListenAndServe(":3000", mux)
}
