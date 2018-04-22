package config

import (
	"os"
	"testing"
)

func TestNewConfig(t *testing.T) {

	configFile := "./default_test.json"

	f, err := os.Open(configFile)

	if err != nil {
		t.Fatalf("failed to open the config file:%s", configFile)
	}

	expectedCfg := &AppCfg{}
	expectedCfg.ServerConfig.ServerUrl = "localhost:8080"
	expectedCfg.SqlConfig.Url = "localhost:3306"
	expectedCfg.SqlConfig.Username = "hdb"
	expectedCfg.SqlConfig.Password = "iop123hdb"
	expectedCfg.SqlConfig.Database = "nf_central"

	cfg := NewConfig(f)

	if expectedCfg.ServerConfig.ServerUrl != cfg.ServerConfig.ServerUrl ||
		expectedCfg.SqlConfig.Username != cfg.SqlConfig.Username ||
		expectedCfg.SqlConfig.Password != cfg.SqlConfig.Password ||
		expectedCfg.SqlConfig.Url != cfg.SqlConfig.Url ||
		expectedCfg.SqlConfig.Database != cfg.SqlConfig.Database {
		t.Fatal("config file parse failed, got the wrong configuration")
	}

	expectedCfg.CacheConfig.ElementSize = 1000000

	if expectedCfg.CacheConfig.ElementSize != cfg.CacheConfig.ElementSize {
		t.Fatal("config file parse failed, the cache size is wrong", cfg.CacheConfig.ElementSize)
	}
}
