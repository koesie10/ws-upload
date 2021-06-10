package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/koesie10/pflagenv"
	"github.com/koesie10/ws-upload/influx"
	"github.com/koesie10/ws-upload/wsupload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

var config = struct {
	Addr string `env:"ADDR" flag:"addr" desc:"the address for the HTTP server to listen on"`

	StationPassword string `env:"STATION_PASSWORD" flag:"station-password,p" desc:"the station password that will be accepted"`

	Influx influx.PublisherOptions `env:",squash"`
}{
	Addr: ":9108",

	Influx: influx.PublisherOptions{
		Addr:            "http://localhost:8086",
		Bucket:          "weather",
		MeasurementName: "weather",
	},
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

	var publishers []wsupload.Publisher

	if config.Influx.Addr != "" {
		publisher, err := influx.NewPublisher(config.Influx)
		if err != nil {
			return fmt.Errorf("failed to create influx publisher: %w", err)
		}
		defer publisher.Close()
		publishers = append(publishers, publisher)

		logrus.Info("Influx publisher enabled")
	}

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

		for _, publisher := range publishers {
			if err := publisher.Publish(obs); err != nil {
				logrus.WithError(err).Errorf("Failed to publish")
			}
		}

		return c.String(http.StatusOK, "OK")
	})

	logrus.Infof("Starting HTTP server on %s", e.Listener.Addr().String())

	return e.Start(config.Addr)
}
