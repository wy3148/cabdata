package config

import (
	"encoding/json"
	"io"
)

type ServerSetting struct {
	ServerUrl string
}

type SqlSettiing struct {
	Url      string
	Username string
	Password string
	Database string
}

type CacheSetting struct {
	CacheDrive  string
	ElementSize int
}

type AppCfg struct {
	SqlConfig    SqlSettiing
	ServerConfig ServerSetting
	CacheConfig  CacheSetting
}

func NewConfig(reader io.Reader) *AppCfg {
	var appCfg AppCfg
	d := json.NewDecoder(reader)

	if err := d.Decode(&appCfg); err != nil {
		panic(err)
	}

	return &appCfg
}
