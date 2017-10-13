package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	systemdNodeStartupDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "systemd_analyse_total_startup_time_secounds",
		Help: "systemd total startup time",
	})
	systemdNodeStartupDurationHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "systemd_analyse_total_startup_time_histogram_secounds",
		Help:    "systemd total startup time histogram",
		Buckets: prometheus.LinearBuckets(1.0, 30.0, 30),
	})
)

// Prom prometheus metrics struct
type Prom struct{}

// Init prometheus metrics and register all metrics objects
func (p *Prom) Init() {
	prometheus.MustRegister(systemdNodeStartupDuration)
	prometheus.MustRegister(systemdNodeStartupDurationHistogram)
}

func (p *Prom) serv(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getTimeFromSystemdAnalyse(input string) (time.Duration, error) {
	split := strings.Split(input, "=")
	onlytime := strings.TrimSpace(split[len(split)-1])

	if strings.Contains(onlytime, "min") {
		onlytime = strings.Replace(onlytime, "min ", "m", -1)
	}

	parsedTime, err := time.ParseDuration(onlytime)
	if err != nil {
		return time.Duration(0), fmt.Errorf("failed to parse duration: %v", err.Error())
	}
	return parsedTime, nil
}

func getSystemdAnalyseOutput() (string, error) {
	cmd := exec.Command("systemd-analyze")
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run systemd analyse : %v", err.Error())
	}
	outString := string(out)
	return outString, nil
}

func main() {
	output, err := getSystemdAnalyseOutput()
	if err != nil {
		panic(err)
	}
	parsedTime, err := getTimeFromSystemdAnalyse(output)
	if err != nil {
		panic(err)
	}
	systemdNodeStartupDuration.Set(parsedTime.Seconds())
	systemdNodeStartupDurationHistogram.Observe(parsedTime.Seconds())
	prom := Prom{}
	prom.Init()
	prom.serv("0.0.0.0:9011")
}
