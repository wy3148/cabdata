package store

import (
	"github.com/wy3148/cabdata/config"
	"os"
	"testing"
)

func getInitStoreInterface() (*Store, error) {

	//here load the config.json file
	//it will verify mysql connection
	configFile := "../config/config.json"

	f, err := os.Open(configFile)

	if err != nil {
		return nil, err
	}

	cfg := config.NewConfig(f)
	s, err := InitStoreInterface(cfg)

	if err != nil {
		return nil, err
	}

	return s, nil
}

func BenchmarkGetTrips(b *testing.B) {

	key := "FF2C42685FE5822F7A6DE63D32ED8193"

	s, err := getInitStoreInterface()

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		s.GetTrips(key, "2013-12-31", true)
	}
}

func TestInitStoreInterface(t *testing.T) {
	_, err := getInitStoreInterface()

	if err != nil {
		t.Fatal(err)
	}
}

func TestGetTrips(t *testing.T) {

	s, err := getInitStoreInterface()

	if err != nil {
		t.Fatal(err)
	}

	key := "FF2C42685FE5822F7A6DE63D32ED8193"
	expected := 20

	res, err := s.GetTrips(key, "2013-12-31", false)

	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Fatal("read from db, the vaule is wrong")
	}

	//read from cache, should get same result
	res, err = s.GetTrips(key, "2013-12-31", true)

	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Fatal("read from cache second time, value is wrong")
	}

	//different date
	res1, err := s.GetTrips(key, "2013-12-30", false)

	if err != nil {
		t.Fatal(err)
	}

	res2, err := s.GetTrips(key, "2013-12-30", true)

	if err != nil {
		t.Fatal(err)
	}

	if res1 != res2 {
		t.Fatal("values from db and cache are not equal")
	}
}

func TestClearCache(t *testing.T) {

	key := "FF2C42685FE5822F7A6DE63D32ED8193"

	s, err := getInitStoreInterface()

	if err != nil {
		t.Fatal(err)
	}

	expected := 20

	if res, _ := s.GetTrips(key, "2013-12-31", true); res != expected {
		t.Fatal("read from cache first time,value is wrong")
	}

	if s.cache.Len() != 1 {
		t.Fatal("there must be one element in cache now")
	}

	res, err := s.GetTrips(key, "2013-12-31", false)

	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Fatal("got wrong value from db")
	}

	s.ClearCache(&key)

	if s.cache.Len() != 0 {
		t.Fatal("there shouldn't be any element in cache now")
	}

	if res, _ := s.GetTrips(key, "2013-12-31", true); res != expected {
		t.Fatal("read from cache again after clean, got wrong value")
	}
}

func BenchmarkClearCacheLRU(b *testing.B) {

	s, err := getInitStoreInterface()

	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		s.cache.Add(i, i)
	}

	for i := 0; i < b.N; i++ {
		s.cache.Remove(i)
	}
}
