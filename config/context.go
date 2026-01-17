package config

type Context struct {
	Config              *Config
	SetBeaconIntervalCh chan int
}
