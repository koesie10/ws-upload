package wsupload

import "time"

type Observation struct {
	StationID    string `ws:"ID" json:"station_id" influx:"station_id,tag" homeassistant:"Station ID"`
	SoftwareType string `ws:"softwaretype" json:"software_type" homeassistant:"Software type"`

	ObservationTime time.Time `ws:"dateutc,layout=2006-01-02 15:04:05,location=UTC" json:"observation_time" influx:"ts" homeassistant:"Observation time,device_class=timestamp,unit_of_measurement=ISO8601"`

	OutsideTemperatureCelsius float64 `ws:"tempf,conversion=fahrenheit_to_celsius" json:"outside_temperature_celsius" homeassistant:"Outside temperature,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`
	IndoorTemperatureCelsius  float64 `ws:"indoortempf,conversion=fahrenheit_to_celsius" json:"indoor_temperature_celsius" homeassistant:"Indoor temperature,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`
	DewpointCelsius           float64 `ws:"dewptf,conversion=fahrenheit_to_celsius" json:"dewpoint_celsius" homeassistant:"Dewpoint,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`
	WindchillCelsius          float64 `ws:"windchillf,conversion=fahrenheit_to_celsius" json:"windchill_celsius" homeassistant:"Windchill,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`

	OutsideRelativeHumidity float64 `ws:"humidity" json:"outside_relative_humidity" homeassistant:"Outside relative humidity,device_class=humidity,unit_of_measurement=%,state_class=measurement"`
	IndoorRelativeHumidity  float64 `ws:"indoorhumidity" json:"indoor_relative_humidity" homeassistant:"Indoor relative humidity,device_class=humidity,unit_of_measurement=%,state_class=measurement"`

	RelativeAtmosphericPressurePascal float64 `ws:"baromin,conversion=inches_of_mercury_to_pascal" json:"relative_atmospheric_pressure_pascal" homeassistant:"Relative atmospheric pressure,device_class=pressure,unit_of_measurement=Pa,state_class=measurement"`
	AbsoluteAtmosphericPressurePascal float64 `ws:"absbaromin,conversion=inches_of_mercury_to_pascal" json:"absolute_atmospheric_pressure_pascal" homeassistant:"Absolute atmospheric pressure,device_class=pressure,unit_of_measurement=Pa,state_class=measurement"`

	UVIndex                           int64   `ws:"UV" json:"uv_index" homeassistant:"UV index,state_class=measurement,unit_of_measurement=UV"`
	SolarRadiationWattPerMeterSquared float64 `ws:"solarradiation" json:"solar_radiation_watt_per_meter_squared" homeassistant:"Solar radiation,state_class=measurement,unit_of_measurement=W/m^2"`

	WindDirectionDegrees     int64   `ws:"winddir" json:"wind_direction_degrees" homeassistant:"Wind direction,state_class=measurement,unit_of_measurement=°"`
	WindSpeedMetersPerSecond float64 `ws:"windspeedmph,conversion=mph_to_meters_per_second" json:"wind_speed_meters_per_second" homeassistant:"Wind speed,state_class=measurement,unit_of_measurement=°"`
	WindGustMetersPerSecond  float64 `ws:"windgustmph,conversion=mph_to_meters_per_second" json:"wind_gust_meters_per_second" homeassistant:"Wind gust,state_class=measurement,unit_of_measurement=°"`

	HourlyRainMillimeters  float64 `ws:"rainin,conversion=inches_of_rain_to_millimeter" json:"hourly_rain_millimeters" homeassistant:"Hourly rain,unit_of_measurement=mm"`
	DailyRainMillimeters   float64 `ws:"dailyrainin,conversion=inches_of_rain_to_millimeter" json:"daily_rain_millimeters" homeassistant:"Daily rain,unit_of_measurement=mm"`
	WeeklyRainMillimeters  float64 `ws:"weeklyrainin,conversion=inches_of_rain_to_millimeter" json:"weekly_rain_millimeters" homeassistant:"Weekly rain,unit_of_measurement=mm"`
	MonthlyRainMillimeters float64 `ws:"monthlyrainin,conversion=inches_of_rain_to_millimeter" json:"monthly_rain_millimeters" homeassistant:"Monthly rain,unit_of_measurement=mm"`
}
