package main

type config struct {
	Version float32           `yml:"version"`
	EnvFile string            `yml:"envfile"`
	Scripts map[string]string `yml:"scripts"`
}
