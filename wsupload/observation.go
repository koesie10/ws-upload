package wsupload

import "time"

type Observation struct {
	StationID    string `ws:"ID" json:"station_id" influx:"station_id,tag"`
	SoftwareType string `ws:"softwaretype" json:"software_type"`

	ObservationTime time.Time `ws:"dateutc,layout=2006-01-02 15:04:05,location=UTC" json:"observation_time" influx:"ts"`

	OutsideTemperatureCelsius float64 `ws:"tempf,conversion=fahrenheit_to_celsius" json:"outside_temperature_celsius"`
	IndoorTemperatureCelsius  float64 `ws:"indoortempf,conversion=fahrenheit_to_celsius" json:"indoor_temperature_celsius"`
	DewpointCelsius           float64 `ws:"dewptf,conversion=fahrenheit_to_celsius" json:"dewpoint_celsius"`
	WindchillCelsius          float64 `ws:"windchillf,conversion=fahrenheit_to_celsius" json:"windchill_celsius"`

	OutsideRelativeHumidity float64 `ws:"humidity" json:"outside_relative_humidity"`
	IndoorRelativeHumidity  float64 `ws:"indoorhumidity" json:"indoor_relative_humidity"`

	RelativeAtmosphericPressurePascal float64 `ws:"baromin,conversion=inches_of_mercury_to_pascal" json:"relative_atmospheric_pressure_pascal"`
	AbsoluteAtmosphericPressurePascal float64 `ws:"absbaromin,conversion=inches_of_mercury_to_pascal" json:"absolute_atmospheric_pressure_pascal"`

	UVIndex                           int64   `ws:"UV" json:"uv_index"`
	SolarRadiationWattPerMeterSquared float64 `ws:"solarradiation" json:"solar_radiation_watt_per_meter_squared"`

	WindDirectionDegrees     int64   `ws:"winddir" json:"wind_direction_degrees"`
	WindSpeedMetersPerSecond float64 `ws:"windspeedmph,conversion=mph_to_meters_per_second" json:"wind_speed_meters_per_second"`
	WindGustMetersPerSecond  float64 `ws:"windgustmph,conversion=mph_to_meters_per_second" json:"wind_gust_meters_per_second"`

	HourlyRainMillimeters  float64 `ws:"rainin,conversion=inches_of_rain_to_millimeter" json:"hourly_rain_millimeters"`
	DailyRainMillimeters   float64 `ws:"dailyrainin,conversion=inches_of_rain_to_millimeter" json:"daily_rain_millimeters"`
	WeeklyRainMillimeters  float64 `ws:"weeklyrainin,conversion=inches_of_rain_to_millimeter" json:"weekly_rain_millimeters"`
	MonthlyRainMillimeters float64 `ws:"monthlyrainin,conversion=inches_of_rain_to_millimeter" json:"monthly_rain_millimeters"`
}
