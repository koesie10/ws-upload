package main

import (
	"crypto/rand"
	"encoding/hex"
	"net"
	"net/http"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/koesie10/pflagenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"ws-upload/influx"
	"ws-upload/wsupload"
)

var config = struct {
	Addr string `env:"ADDR" flag:"addr" desc:"the address for the HTTP server to listen on"`

	StationPassword string `env:"STATION_PASSWORD" flag:"station-password,p" desc:"the station password that will be accepted"`

	InfluxAddr         string `env:"INFLUX_ADDR" flag:"influx-addr" desc:"InfluxDB HTTP address"`
	InfluxAuthToken    string `env:"INFLUX_AUTH_TOKEN" flag:"influx-auth-token" desc:"InfluxDB auth token, use username:password for InfluxDB 1.8"`
	InfluxOrganization string `env:"INFLUX_ORGANIZATION" flag:"influx-organization" desc:"InfluxDB organization, do not set if using InfluxDB 1.8"`
	InfluxBucket       string `env:"INFLUX_BUCKET" flag:"influx-bucket" desc:"InfluxDB bucket, set to database/retention-policy or database for InfluxDB 1.8"`
	MeasurementName    string `env:"MEASUREMENT_NAME" flag:"measurement-name" desc:"InfluxDB measurement name"`
}{
	Addr: ":9108",

	InfluxAddr:      "http://localhost:8086",
	InfluxBucket:    "weather",
	MeasurementName: "weather",
}

func main() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	flagSet := pflag.NewFlagSet("ws-upload", pflag.ExitOnError)

	if err := pflagenv.Setup(flagSet, &config); err != nil {
		logrus.Fatal(err)
	}

	if err := pflagenv.Parse(&config); err != nil {
		logrus.Fatal(err)
	}

	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	if config.StationPassword == "" {
		key := make([]byte, 8)
		if _, err := rand.Read(key); err != nil {
			return err
		}

		config.StationPassword = hex.EncodeToString(key)
		logrus.Infof("Station password is %s, please set it using the STATION_PASSWORD environment variable or the --station-password/-p flag", config.StationPassword)
	}

	options := influxdb2.DefaultOptions()
	options.SetPrecision(time.Second)
	client := influxdb2.NewClientWithOptions(config.InfluxAddr, config.InfluxAuthToken, options)
	defer client.Close()

	writeAPI := client.WriteAPI(config.InfluxOrganization, config.InfluxBucket)

	l, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.RequestID())
	e.Logger = NewEchoLogger(logrus.WithField("system", "echo"))
	e.Listener = l

	e.GET("/api/v1/observe", func(c echo.Context) error {
		entry := logrus.WithFields(logrus.Fields{
			"scheme":     c.Scheme(),
			"method":     c.Request().Method,
			"path":       c.Path(),
			"remote_ip":  c.RealIP(),
			"request_Id": c.Response().Header().Get(echo.HeaderXRequestID),
		})

		if c.QueryParam("PASSWORD") != config.StationPassword {
			return c.String(http.StatusUnauthorized, "Bad password")
		}

		if c.QueryParam("action") != "updateraw" {
			return c.String(http.StatusBadRequest, "Invalid action")
		}

		obs, err := wsupload.Parse(c.QueryParams(), logrus.WithField("", ""))
		if err != nil {
			entry.WithError(err).Errorf("Failed to parse observation")
			return err
		}

		point, err := influx.CreatePoint(obs, config.MeasurementName)
		if err != nil {
			entry.WithError(err).Error("Failed to convert to point")
			return err
		}

		writeAPI.WritePoint(point)

		return c.String(http.StatusOK, "OK")
	})

	logrus.Infof("Starting HTTP server on %s", e.Listener.Addr().String())

	return e.Start(config.Addr)
}
