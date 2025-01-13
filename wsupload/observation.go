package wsupload

import "time"

type Observation struct {
	StationID    string `ws:"ID" json:"station_id" influx:"station_id,tag" homeassistant:"Station ID"`
	SoftwareType string `ws:"softwaretype" json:"software_type" homeassistant:"Software type"`

	ObservationTime time.Time `ws:"dateutc,layout=2006-01-02 15:04:05,location=UTC" json:"observation_time" influx:"ts" homeassistant:"Observation time,device_class=timestamp"`

	OutsideTemperatureCelsius NullFloat64 `ws:"tempf,conversion=fahrenheit_to_celsius" json:"outside_temperature_celsius" homeassistant:"Outside temperature,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`
	IndoorTemperatureCelsius  NullFloat64 `ws:"indoortempf,conversion=fahrenheit_to_celsius" json:"indoor_temperature_celsius" homeassistant:"Indoor temperature,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`
	DewpointCelsius           NullFloat64 `ws:"dewptf,conversion=fahrenheit_to_celsius" json:"dewpoint_celsius" homeassistant:"Dewpoint,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`
	WindchillCelsius          NullFloat64 `ws:"windchillf,conversion=fahrenheit_to_celsius" json:"windchill_celsius" homeassistant:"Windchill,device_class=temperature,unit_of_measurement=°C,state_class=measurement"`

	OutsideRelativeHumidity NullFloat64 `ws:"humidity" json:"outside_relative_humidity" homeassistant:"Outside relative humidity,device_class=humidity,unit_of_measurement=%,state_class=measurement"`
	IndoorRelativeHumidity  NullFloat64 `ws:"indoorhumidity" json:"indoor_relative_humidity" homeassistant:"Indoor relative humidity,device_class=humidity,unit_of_measurement=%,state_class=measurement"`

	RelativeAtmosphericPressurePascal NullFloat64 `ws:"baromin,conversion=inches_of_mercury_to_pascal" json:"relative_atmospheric_pressure_pascal" homeassistant:"Relative atmospheric pressure,device_class=pressure,unit_of_measurement=Pa,state_class=measurement"`
	AbsoluteAtmosphericPressurePascal NullFloat64 `ws:"absbaromin,conversion=inches_of_mercury_to_pascal" json:"absolute_atmospheric_pressure_pascal" homeassistant:"Absolute atmospheric pressure,device_class=pressure,unit_of_measurement=Pa,state_class=measurement"`

	UVIndex                           NullFloat64 `ws:"UV" json:"uv_index" homeassistant:"UV index,state_class=measurement,unit_of_measurement=UV"`
	SolarRadiationWattPerMeterSquared NullFloat64 `ws:"solarradiation" json:"solar_radiation_watt_per_meter_squared" homeassistant:"Solar radiation,state_class=measurement,unit_of_measurement=W/m^2"`

	WindDirectionDegrees     NullInt64   `ws:"winddir" json:"wind_direction_degrees" homeassistant:"Wind direction,state_class=measurement,unit_of_measurement=°"`
	WindSpeedMetersPerSecond NullFloat64 `ws:"windspeedmph,conversion=mph_to_meters_per_second" json:"wind_speed_meters_per_second" homeassistant:"Wind speed,state_class=measurement,unit_of_measurement=m/s"`
	WindGustMetersPerSecond  NullFloat64 `ws:"windgustmph,conversion=mph_to_meters_per_second" json:"wind_gust_meters_per_second" homeassistant:"Wind gust,state_class=measurement,unit_of_measurement=m/s"`

	HourlyRainMillimeters  NullFloat64 `ws:"rainin,conversion=inches_of_rain_to_millimeter" json:"hourly_rain_millimeters" homeassistant:"Hourly rain,unit_of_measurement=mm"`
	DailyRainMillimeters   NullFloat64 `ws:"dailyrainin,conversion=inches_of_rain_to_millimeter" json:"daily_rain_millimeters" homeassistant:"Daily rain,unit_of_measurement=mm"`
	WeeklyRainMillimeters  NullFloat64 `ws:"weeklyrainin,conversion=inches_of_rain_to_millimeter" json:"weekly_rain_millimeters" homeassistant:"Weekly rain,unit_of_measurement=mm"`
	MonthlyRainMillimeters NullFloat64 `ws:"monthlyrainin,conversion=inches_of_rain_to_millimeter" json:"monthly_rain_millimeters" homeassistant:"Monthly rain,unit_of_measurement=mm"`
}
