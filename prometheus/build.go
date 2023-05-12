package prometheus

import (
	"github.com/koesie10/ws-upload/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	infoGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name:      "info",
		Help:      "Information about the ws-upload build",
		Namespace: "ws_upload",
		ConstLabels: map[string]string{
			"version":    version.Version,
			"commit":     version.Commit,
			"build_date": version.BuildDate,
		},
	})
)

func init() {
	infoGauge.Set(1)
}
