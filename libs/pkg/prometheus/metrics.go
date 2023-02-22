package prometheus

import (
	"net/http"

	"github.com/LSDXXX/libs/config"
	"github.com/LSDXXX/libs/pkg/container"
	"github.com/LSDXXX/libs/pkg/util"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var conf *config.Config

// Init prometheus server
func Init() error {
	util.PanicWhenError(container.Resolve(&conf))
	go listenServer()
	//cleanCronTaskInit(false)
	return nil
}

// listenServer 监听
func listenServer() {
	if len(conf.PrometheusBindURL) > 0 {
		http.Handle("/metrics", promhttp.Handler())
		logrus.Infof("metrics: prometheus listenServer:%s", conf.PrometheusBindURL)
		logrus.Warn(http.ListenAndServe(conf.PrometheusBindURL, nil))
	}
}

/*// cleanCronTaskInit clean metrics 定时任务重置指标
func cleanCronTaskInit(enable bool) {
	if enable && len(conf.Cron) > 0 {
		duration := conf.Cron
		if len(duration) == 0 {
			return
		}
		logrus.Infof("metrics: cleanCronTask: %s", conf.Cron)

		task := cron.New()
		_, err := task.AddFunc(duration, func() {
			FlowReceiveCounterReset()
			FlowErrorCounterReset()
			NodeReceiveCounterReset()
			NodeErrorCounterReset()
		})
		if err != nil {
			panic(fmt.Errorf("metrics: cleanCronTask AddFunc fail, error: %s", err.Error()))
		}
		task.Start()
	}
}*/
