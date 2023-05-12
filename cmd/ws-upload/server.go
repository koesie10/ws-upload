package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/brpaz/echozap"
	"github.com/koesie10/pflagenv"
	"github.com/koesie10/ws-upload/influx"
	"github.com/koesie10/ws-upload/jsondebug"
	"github.com/koesie10/ws-upload/mqtt"
	"github.com/koesie10/ws-upload/wsupload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var serverConfig = struct {
	Addr string `env:"ADDR" flag:"addr" desc:"the address for the HTTP server to listen on"`

	StationPassword string `env:"STATION_PASSWORD" flag:"station-password,p" desc:"the station password that will be accepted"`

	Influx influx.PublisherOptions `env:",squash"`
	MQTT   mqtt.PublisherOptions   `env:",squash"`

	EnableJSONDebug   bool `env:"ENABLE_JSON_DEBUG" flag:"enable-json-debug" desc:"enable json debug output"`
	EnableInfluxDebug bool `env:"ENABLE_INFLUX_DEBUG" flag:"enable-influx-debug" desc:"enable influx debug output"`
}{
	Addr: ":9108",

	Influx: influx.PublisherOptions{
		Addr:            "http://localhost:8086",
		Bucket:          "weather",
		MeasurementName: "weather",
	},

	MQTT: mqtt.PublisherOptions{
		Brokers: []string{"tcp://127.0.0.1:1883"},
		Topic:   "homeassistant/sensor/sensorWeatherStation/state",
		HomeAssistant: mqtt.HomeAssistantOptions{
			DiscoveryEnabled:  true,
			DiscoveryInterval: 30 * time.Second,
			DiscoveryQoS:      1, // At least once
			DiscoveryPrefix:   "homeassistant",
			DevicePrefix:      "weatherstation_",
		},
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the ws-upload server",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := pflagenv.Parse(&serverConfig); err != nil {
			return err
		}

		return nil
	},
	RunE: RunServer,
}

func RunServer(cmd *cobra.Command, args []string) error {
	if serverConfig.StationPassword == "" {
		key := make([]byte, 8)
		if _, err := rand.Read(key); err != nil {
			return err
		}

		serverConfig.StationPassword = hex.EncodeToString(key)
		logger.Info("Station password has been generated automatically, please set it using the STATION_PASSWORD environment variable or the --station-password/-p flag", zap.String("ws_upload.station_password", serverConfig.StationPassword))
	}

	var publishers []wsupload.Publisher

	if serverConfig.EnableJSONDebug {
		publisher, err := jsondebug.NewDebugPublisher()
		if err != nil {
			return fmt.Errorf("failed to create JSON debug publisher: %w", err)
		}
		defer publisher.Close()
		publishers = append(publishers, publisher)

		logger.Info("JSON debug publisher enabled")
	}

	if serverConfig.Influx.Addr != "" {
		publisher, err := influx.NewPublisher(serverConfig.Influx)
		if err != nil {
			return fmt.Errorf("failed to create Influx publisher: %w", err)
		}
		defer publisher.Close()
		publishers = append(publishers, publisher)

		logger.Info("Influx publisher enabled")
	}

	if serverConfig.EnableInfluxDebug {
		publisher, err := influx.NewDebugPublisher(influx.DebugPublisherOptions{
			MeasurementName: "weather",
		})
		if err != nil {
			return fmt.Errorf("failed to create Influx debug publisher: %w", err)
		}
		defer publisher.Close()
		publishers = append(publishers, publisher)

		logger.Info("Influx debug publisher enabled")
	}

	if len(serverConfig.MQTT.Brokers) > 0 {
		publisher, err := mqtt.NewPublisher(logger, serverConfig.MQTT)
		if err != nil {
			return fmt.Errorf("failed to create MQTT publisher: %w", err)
		}
		defer publisher.Close()
		publishers = append(publishers, publisher)

		logger.Info("MQTT publisher enabled")
	}

	l, err := net.Listen("tcp", serverConfig.Addr)
	if err != nil {
		return err
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(middleware.Recover())
	e.Use(echozap.ZapLogger(logger.With(zap.String("component", "echo"))))
	e.Use(middleware.RequestID())
	e.Listener = l

	e.GET("/api/v1/observe", func(c echo.Context) error {
		entry := logger.With(
			zap.String("http.scheme", c.Scheme()),
			zap.String("http.method", c.Request().Method),
			zap.String("http.path", c.Path()),
			zap.String("http.remote_ip", c.RealIP()),
			zap.String("http.request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
		)

		if c.QueryParam("PASSWORD") != serverConfig.StationPassword {
			return c.String(http.StatusUnauthorized, "Bad password")
		}

		if c.QueryParam("action") != "updateraw" {
			return c.String(http.StatusBadRequest, "Invalid action")
		}

		obs, err := wsupload.Parse(c.QueryParams(), logger)
		if err != nil {
			entry.Error("Failed to parse observation", zap.Error(err))
			return err
		}

		if !obs.IndoorTemperatureCelsius.Valid || obs.IndoorTemperatureCelsius.Float64 < -50 || obs.IndoorTemperatureCelsius.Float64 > 80 {
			entry.Error("Invalid indoor temperature", zap.Bool("ws_upload.indoor_temperature_celsius_valid", obs.IndoorTemperatureCelsius.Valid), zap.Float64("ws_upload.indoor_temperature_celsius", obs.IndoorTemperatureCelsius.Float64))
			return c.String(http.StatusOK, "OK")
		}

		for _, publisher := range publishers {
			if err := publisher.Publish(obs); err != nil {
				entry.Error("Failed to publish observation", zap.Error(err))
			}
		}

		return c.String(http.StatusOK, "OK")
	})

	e.POST("/api/v1/mqtt/homeassistant/delete-all-devices", func(c echo.Context) error {
		if err := mqtt.DeleteAllDevices(serverConfig.MQTT); err != nil {
			return err
		}

		return c.String(http.StatusOK, "OK")
	})

	logger.Info("Starting HTTP server", zap.String("net.addr", serverConfig.Addr))

	return e.Start(serverConfig.Addr)
}

func init() {
	rootCmd.AddCommand(serverCmd)

	if err := pflagenv.Setup(serverCmd.Flags(), &serverConfig); err != nil {
		log.Fatal(err)
	}
}
