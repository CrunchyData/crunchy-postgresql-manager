package collect

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	uniformDomain     = 200
	normDomain        = 200
	normMean          = 10
	oscillationPeriod = 10 * time.Minute
)

var guages = make([]prometheus.Gauge, 3)

func init() {

	cnt := 3
	x := 0
	for x < cnt {
		if x == 0 {
			guages[x] = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "cpm_server_cpu",
				Help:        "Current temperature of the CPU.",
				ConstLabels: prometheus.Labels{"server": "espresso"},
			})
		} else {
			guages[x] = prometheus.NewGauge(prometheus.GaugeOpts{
				Name:        "cpm_server_cpu",
				Help:        "Current temperature of the CPU.",
				ConstLabels: prometheus.Labels{"server": "server" + strconv.Itoa(x)},
			})
		}
		fmt.Println("registered ")
		prometheus.MustRegister(guages[x])
		x++
	}
}

func main() {

	go func() {
		for {
			fmt.Println("collecting server cpu goroutine ")
			v := rand.Float64() * 100.00
			CollectServerCPU()
			guages[0].Set(v)
			time.Sleep(time.Duration(10000 * time.Millisecond))
		}
	}()
	/*
		go func() {
			for {
				fmt.Println("firing up goroutine ")
				v := rand.Float64() * 100.00
				guages[1].Set(v)
				time.Sleep(time.Duration(1000 * time.Millisecond))
			}
		}()
		go func() {
			for {
				fmt.Println("firing up goroutine ")
				v := rand.Float64() * 100.00
				guages[2].Set(v)
				time.Sleep(time.Duration(1000 * time.Millisecond))
			}
		}()
	*/

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(":8080", nil)
}

func CollectServerCPU() error {
	var err error
	return err
}
