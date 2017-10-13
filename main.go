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
)

// Prom prometheus metrics struct
type Prom struct{}

// Init prometheus metrics and register all metrics objects
func (p *Prom) Init() {
	prometheus.MustRegister(systemdNodeStartupDuration)
}

func (p *Prom) serv(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

func getTimeFromSystemdAnalyse(input string) (time.Duration, error) {
	split := strings.Split(input, "=")
	onlytime := strings.TrimSpace(split[len(split)-1])
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
	source := "Startup finished in 1.238s (kernel) + 1.677s (initrd) + 7.406s (userspace) = 10.322s"
	parsedTime, err := getTimeFromSystemdAnalyse(source)
	if err != nil {
		panic(err)
	}
	systemdNodeStartupDuration.Set(parsedTime.Seconds())
}
