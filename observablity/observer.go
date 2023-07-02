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
	CreateOrderApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "create_order_api_count",
		Help: "Counter for create order api",
	})

	GetOrderApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "get_order_api_count",
		Help: "Counter for get order api",
	})

	UpdateOrderedProductsApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_ordered_products_api_count",
		Help: "Counter for update ordered products api",
	})
	UpdateOrderStatusApiCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "update_order_status_api_count",
		Help: "Counter for create order api",
	})
	OrderServiceApiElapsedTime = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "order_service_apis_elapsed_time",
		Help: "Time taken to execute an API",
	}, []string{"api_name"})
)

func init() {
	prometheus.MustRegister(upTime)
	prometheus.MustRegister(CreateOrderApiCounter)
	prometheus.MustRegister(GetOrderApiCounter)
	prometheus.MustRegister(UpdateOrderedProductsApiCounter)
	prometheus.MustRegister(UpdateOrderStatusApiCounter)
	prometheus.MustRegister(OrderServiceApiElapsedTime)

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
			case <-time.After(10 * time.Second):
				bootTime += 10.0
				upTime.WithLabelValues(serverName).Set(bootTime)
			}

		}
	}(serverName)
	upTime.WithLabelValues(serverName).Set(bootTime)
	r := gin.New()
	r.GET("/metrics", prometheusHandler())
	r.Run(observerAddr)
}
