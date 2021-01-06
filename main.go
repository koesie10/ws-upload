package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/koesie10/pflagenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"ws-upload/parse"
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

		observationMap := make(map[string]interface{})
		for key, values := range c.QueryParams() {
			if len(values) == 1 {
				observationMap[key] = values[0]
			} else {
				observationMap[key] = values
			}
		}

		point, err := observationToPoint(entry, c.QueryParams(), config.MeasurementName)
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

func observationToPoint(entry *logrus.Entry, params url.Values, measurementName string) (*write.Point, error) {
	fields := make(map[string]interface{})
	tags := make(map[string]string)

	parseTag(entry, tags, params, "station_id", "ID", parse.StringStringIdentity)

	dateUTC := params.Get("dateutc")
	if dateUTC == "" {
		return nil, fmt.Errorf("missing dateutc in params")
	}

	observationTime, err := time.ParseInLocation("2006-01-02 15:04:05", dateUTC, time.UTC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %w", err)
	}

	parseField(entry, fields, params, "software_type", "softwaretype", parse.StringIdentity)
	parseField(entry, fields, params, "outside_temperature_celsius", "tempf", parse.FloatFunc(parse.ParseFahrenheitToCelsius))
	parseField(entry, fields, params, "indoor_temperature_celsius", "indoortempf", parse.FloatFunc(parse.ParseFahrenheitToCelsius))
	parseField(entry, fields, params, "dewpoint_celsius", "dewptf", parse.FloatFunc(parse.ParseFahrenheitToCelsius))
	parseField(entry, fields, params, "windchill_celsius", "windchillf", parse.FloatFunc(parse.ParseFahrenheitToCelsius))

	parseField(entry, fields, params, "outside_relative_humidity", "humidity", parse.FloatFunc(parse.ParseFloat))
	parseField(entry, fields, params, "indoor_relative_humidity", "indoorhumidity", parse.FloatFunc(parse.ParseFloat))

	parseField(entry, fields, params, "relative_atmospheric_pressure_pascal", "baromin", parse.FloatFunc(parse.ParseInchesOfMercuryToPascal))
	parseField(entry, fields, params, "absolute_atmospheric_pressure_pascal", "absbaromin", parse.FloatFunc(parse.ParseInchesOfMercuryToPascal))

	parseField(entry, fields, params, "uv_index", "UV", parse.IntFunc(parse.ParseInt))
	parseField(entry, fields, params, "solar_radiation_watt_per_meter_squared", "solarradiation", parse.FloatFunc(parse.ParseFloat))

	parseField(entry, fields, params, "wind_direction_degrees", "winddir", parse.IntFunc(parse.ParseInt))
	parseField(entry, fields, params, "wind_speed_meters_per_second", "windspeedmph", parse.FloatFunc(parse.ParseMphToMetersPerSecond))
	parseField(entry, fields, params, "wind_gust_meters_per_second", "windgustmph", parse.FloatFunc(parse.ParseMphToMetersPerSecond))

	parseField(entry, fields, params, "hourly_rain_millimeters", "rainin", parse.FloatFunc(parse.ParseInchesOfRainToMillimeter))
	parseField(entry, fields, params, "daily_rain_millimeters", "dailyrainin", parse.FloatFunc(parse.ParseInchesOfRainToMillimeter))
	parseField(entry, fields, params, "weekly_rain_millimeters", "weeklyrainin", parse.FloatFunc(parse.ParseInchesOfRainToMillimeter))
	parseField(entry, fields, params, "monthly_rain_millimeters", "monthlyrainin", parse.FloatFunc(parse.ParseInchesOfRainToMillimeter))

	return influxdb2.NewPoint(measurementName, tags, fields, observationTime), nil
}

func parseField(entry *logrus.Entry, fields map[string]interface{}, params url.Values, field, queryParam string, parseFunc func(value string) (interface{}, error)) {
	value := params.Get(queryParam)
	if value == "" {
		entry.Warnf("Missing query param '%s' for field '%s'", queryParam, field)
		return
	}

	fieldValue, err := parseFunc(value)
	if err != nil {
		entry.WithError(err).Errorf("Failed to parse query param '%s' for field '%s' with value %q", queryParam, field, value)
		return
	}

	fields[field] = fieldValue
}

func parseTag(entry *logrus.Entry, fields map[string]string, params url.Values, tag, queryParam string, parseFunc func(value string) (string, error)) {
	value := params.Get(queryParam)
	if value == "" {
		entry.Warnf("Missing query param '%s' for field '%s'", queryParam, tag)
		return
	}

	fieldValue, err := parseFunc(value)
	if err != nil {
		entry.WithError(err).Errorf("Failed to parse query param '%s' for tag '%s' with value %q", queryParam, tag, value)
		return
	}

	fields[tag] = fieldValue
}
