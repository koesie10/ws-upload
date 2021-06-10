package wsupload

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/structtag"
	"github.com/sirupsen/logrus"
	"github.com/koesie10/ws-upload/x"
)

var timeType = reflect.TypeOf(time.Time{})

func Parse(params url.Values, entry *logrus.Entry) (*Observation, error) {
	obs := Observation{}

	v := reflect.ValueOf(&obs)
	reflectValue := v.Elem()

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		field := reflectValue.Type().Field(i)

		tag, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return nil, fmt.Errorf("failed to parse struct tag for %s: %w", field.Name, err)
		}

		wsTag, err := tag.Get("ws")
		if err != nil {
			continue
		}

		var setFunc func(value string, fieldValue reflect.Value) error

		options := x.ParseStructTagOptions(wsTag.Options)

		switch field.Type.Kind() {
		case reflect.String:
			setFunc = func(value string, fieldValue reflect.Value) error {
				fieldValue.SetString(value)
				return nil
			}
		case reflect.Float64:
			conversion, ok := options["conversion"]

			var transformFunc func (value float64) float64

			if ok {
				switch conversion {
				case "fahrenheit_to_celsius":
					transformFunc = func(fahrenheit float64) float64 {
						return (fahrenheit - 32.0) * 5.0 / 9.0
					}
				case "inches_of_mercury_to_pascal":
					transformFunc = func(inHg float64) float64 {
						return inHg * 3386
					}
				case "mph_to_meters_per_second":
					transformFunc = func(mph float64) float64 {
						return mph * 0.44704
					}
				case "inches_of_rain_to_millimeter":
					transformFunc = func(inRain float64) float64 {
						return inRain * 25.4
					}
				default:
					return nil, fmt.Errorf("unsupported conversion %s for %s")
				}

				delete(options, "conversion")
			} else {
				transformFunc = func(value float64) float64 {
					return value
				}
			}

			setFunc = func(value string, fieldValue reflect.Value) error {
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return fmt.Errorf("failed to parse value: %w", err)
				}

				fieldValue.SetFloat(transformFunc(v))

				return nil
			}
		case reflect.Int64:
			setFunc = func(value string, fieldValue reflect.Value) error {
				v, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return fmt.Errorf("failed to parse value: %w", err)
				}

				fieldValue.SetInt(v)

				return nil
			}
		case reflect.Struct:
			if field.Type.AssignableTo(timeType) {
				var location *time.Location
				if locationOption, ok := options["location"]; ok {
					if locationOption == "UTC" {
						location = time.UTC
					} else {
						return nil, fmt.Errorf("unsupported location %s for field %s", locationOption, field.Name)
					}

					delete(options, "location")
				}

				layout := time.RFC3339
				if formatOption, ok := options["layout"]; ok {
					layout = formatOption

					delete(options, "layout")
				}

				setFunc = func(value string, fieldValue reflect.Value) error {
					t, err := time.ParseInLocation(layout, value, location)
					if err != nil {
						return fmt.Errorf("failed to parse date: %w", err)
					}

					fieldValue.Set(reflect.ValueOf(t))

					return nil
				}
			}
		}

		if setFunc == nil {
			return nil, fmt.Errorf("unsupported field type %s for %s", field.Type, field.Name)
		}

		if len(options) > 0 {
			return nil, fmt.Errorf("unused options %s for %s", options, field.Name)
		}

		queryValue := params.Get(wsTag.Name)
		if queryValue == "" {
			entry.Warnf("Missing query param '%s' for field '%s'", wsTag.Name, field.Name)
			continue
		}

		if err := setFunc(queryValue, fieldValue); err != nil {
			entry.WithError(err).Errorf("Failed to parse query param '%s' for field '%s' with value %q", wsTag.Name, field.Name, queryValue)
			continue
		}
	}

	return &obs, nil
}
