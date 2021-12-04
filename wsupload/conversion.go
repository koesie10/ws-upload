package wsupload

import "fmt"

func getConversionTransformFunc(options map[string]string) (func(value float64) float64, error) {
	conversion, ok := options["conversion"]

	var transformFunc func(value float64) float64

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
			return nil, fmt.Errorf("unsupported conversion %s", conversion)
		}

		delete(options, "conversion")
	} else {
		transformFunc = func(value float64) float64 {
			return value
		}
	}

	return transformFunc, nil
}
