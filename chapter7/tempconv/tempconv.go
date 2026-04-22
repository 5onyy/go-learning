package tempconv

import (
	"flag"
	"fmt"
)

type Celsius float64
type Farenheit float64
type Kelvin float64

const (
	AbsolueZeroC Celsius = -273.15
	FreezingC    Celsius = 0
	BoilingC     Celsius = 100
)

func CToF(c Celsius) Farenheit {
	return Farenheit(c*9/5 + 32)
}

func FToC(f Farenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

func KToC(f Kelvin) Celsius {
	return Celsius(f - 273.15)
}

type celsiusFlag struct{ Celsius }

func (f *celsiusFlag) Set(s string) error {
	var unit string
	var value float64
	fmt.Sscanf(s, "%f%s", &value, &unit)
	switch unit {
	case "C", "°C":
		f.Celsius = Celsius(value)
		return nil
	case "F", "°F":
		f.Celsius = FToC(Farenheit(value))
		return nil
	case "K", "°K":
		f.Celsius = KToC(Kelvin(value))
		return nil
		// Exercise 7.6: Add support for Kelvin temperatures to tempflag
	}
	return fmt.Errorf("invalid temperature %q", s)
}

func (f *celsiusFlag) String() string {
	return fmt.Sprintf("%g°C", f.Celsius)
}

func CelsiusFlag(name string, value Celsius, usage string) *Celsius {
	f := celsiusFlag{value}
	flag.Var(&f, name, usage)
	return &f.Celsius
}
