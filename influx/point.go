package influx

import (
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/structtag"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/koesie10/ws-upload/wsupload"
	"github.com/koesie10/ws-upload/x"
)

func CreatePoint(obs *wsupload.Observation, measurementName string) (*write.Point, error) {
	fields := make(map[string]interface{})
	tags := make(map[string]string)

	var ts time.Time

	v := reflect.ValueOf(obs)
	reflectValue := v.Elem()

	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		field := reflectValue.Type().Field(i)

		tag, err := structtag.Parse(string(field.Tag))
		if err != nil {
			return nil, fmt.Errorf("failed to parse struct tag for %s: %w", field.Name, err)
		}

		influxTag, _ := tag.Get("influx")
		jsonTag, _ := tag.Get("json")

		if influxTag != nil && influxTag.Name == "ts" {
			if !ts.IsZero() {
				return nil, fmt.Errorf("multiple timestamps found for %s", field.Name)
			}

			ts = fieldValue.Interface().(time.Time)

			continue
		}

		var fieldName string
		var options map[string]string
		if jsonTag != nil {
			fieldName = jsonTag.Name
		}
		if influxTag != nil {
			fieldName = influxTag.Name
			options = x.ParseStructTagOptions(influxTag.Options)
		}

		if fieldName == "" || fieldName == "-" {
			fmt.Println("skipping")
			continue
		}

		if _, ok := options["tag"]; ok {
			tags[fieldName] = fieldValue.String()
			continue
		}

		fields[fieldName] = fieldValue.Interface()
	}

	if ts.IsZero() {
		ts = time.Now()
	}

	return influxdb2.NewPoint(measurementName, tags, fields, ts), nil
}
