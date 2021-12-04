package wsupload

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/fatih/structtag"
	"github.com/koesie10/ws-upload/x"
	"github.com/sirupsen/logrus"
)

var timeType = reflect.TypeOf(time.Time{})
var nullFloat64Type = reflect.TypeOf(NullFloat64{})
var nullInt64Type = reflect.TypeOf(NullInt64{})

func Parse(params url.Values, entry *logrus.Entry) (*Observation, error) {
	obs := Observation{}

	v := reflect.ValueOf(&obs)
	reflectValue := v.Elem()

	var optional bool

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
			transformFunc, err := getConversionTransformFunc(options)
			if err != nil {
				return nil, fmt.Errorf("failed to get conversion transform func for %s: %w", field.Name, err)
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
			} else if field.Type.AssignableTo(nullFloat64Type) {
				optional = true

				transformFunc, err := getConversionTransformFunc(options)
				if err != nil {
					return nil, fmt.Errorf("failed to get conversion transform func for %s: %w", field.Name, err)
				}

				setFunc = func(value string, fieldValue reflect.Value) error {
					if value == "-9999" {
						fieldValue.Set(reflect.ValueOf(NullFloat64{Valid: false}))
						return nil
					}

					v, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return fmt.Errorf("failed to parse value: %w", err)
					}

					fieldValue.Set(reflect.ValueOf(NullFloat64{Valid: true, Float64: transformFunc(v)}))

					return nil
				}
			} else if field.Type.AssignableTo(nullInt64Type) {
				optional = true

				setFunc = func(value string, fieldValue reflect.Value) error {
					if value == "-9999" {
						fieldValue.Set(reflect.ValueOf(NullInt64{Valid: false}))
						return nil
					}

					v, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return fmt.Errorf("failed to parse value: %w", err)
					}

					fieldValue.Set(reflect.ValueOf(NullInt64{Valid: true, Int64: v}))

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
			var level = logrus.WarnLevel
			if optional {
				level = logrus.DebugLevel
			}

			entry.Logf(level, "Missing query param '%s' for field '%s'", wsTag.Name, field.Name)
			continue
		}

		if err := setFunc(queryValue, fieldValue); err != nil {
			entry.WithError(err).Errorf("Failed to parse query param '%s' for field '%s' with value %q", wsTag.Name, field.Name, queryValue)
			continue
		}
	}

	return &obs, nil
}
