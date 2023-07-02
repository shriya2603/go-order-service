package observablity

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	upTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "uptime",
		Help: "Uptime of service"},
		[]string{"service_name"})
)

func init() {
	prometheus.MustRegister(upTime)
}
func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func StartObserver(wg *sync.WaitGroup, observerAddr string, serverName string, bootTime float64) {
	defer wg.Done()
	wg.Add(1)
	go func(serverName string) {
		defer wg.Done()
		for {
			select {
			case <-time.After(30 * time.Second):
				bootTime += 30.0
				upTime.WithLabelValues(serverName).Set(bootTime)
			}

		}
	}(serverName)
	upTime.WithLabelValues(serverName).Set(bootTime)
	r := gin.New()
	r.GET("/metrics", prometheusHandler())
	r.Run(observerAddr)
}
