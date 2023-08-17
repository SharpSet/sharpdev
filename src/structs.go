package main

type config struct {
	Version float32           `yml:"version"`
	EnvFile string            `yml:"envfile"`
	Setup   string            `yml:"setup"`
	Scripts map[string]string `yml:"scripts"`
	Values  map[string]string `yml:"values"`
}

// Version Number
var Version float32 = 1.8
