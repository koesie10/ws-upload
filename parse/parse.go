package parse

import (
	"fmt"
	"strconv"
)

func ParseFahrenheitToCelsius(value string) (float64, error) {
	fahrenheit, err := ParseFloat(value)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse temperature: %w", err)
	}

	return (fahrenheit - 32.0) * 5.0 / 9.0, nil
}

func ParseMphToMetersPerSecond(value string) (float64, error) {
	mph, err := ParseFloat(value)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse mph: %w", err)
	}

	return mph * 0.44704, nil
}

func ParseInchesOfMercuryToPascal(value string) (float64, error) {
	inHg, err := ParseFloat(value)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse inHg: %w", err)
	}

	return inHg * 3386, nil
}

func ParseInchesOfRainToMillimeter(value string) (float64, error) {
	inHg, err := ParseFloat(value)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse inches of rain: %w", err)
	}

	return inHg * 25.4, nil
}

func ParseInt(value string) (int64, error) {
	return strconv.ParseInt(value, 10, 64)
}

func ParseFloat(value string) (float64, error) {
	return strconv.ParseFloat(value, 64)
}

func IntFunc(f func (value string) (int64, error)) (func (value string) (interface{}, error)) {
	return func(value string) (interface{}, error) {
		return f(value)
	}
}

func FloatFunc(f func (value string) (float64, error)) (func (value string) (interface{}, error)) {
	return func(value string) (interface{}, error) {
		return f(value)
	}
}

func StringIdentity(value string) (interface{}, error) {
	return value, nil
}

func StringStringIdentity(value string) (string, error) {
	return value, nil
}
