package store

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/wy3148/cabdata/config"
	"log"
	"strings"
)

type StoreInterface interface {

	//if cache is false, get the value from mysqldb directly
	GetTrips(id string, date string, cache bool) (int, error)

	//if key is nil,clear all cache
	ClearCache(id *string)

	//return the number of elements in cache, used for testing
	Len() int
}

type Store struct {
	cache *Cache
	db    *sql.DB
}

var querySql = "SELECT count(*) as cnt FROM cab_trip_data WHERE medallion=? AND Date(pickup_datetime)=?"

const (
	entitySize = 1000000
)

func InitStoreInterface(cfg *config.AppCfg) (*Store, error) {
	s := &Store{}

	var err error
	size := entitySize

	if cfg.CacheConfig.ElementSize > 0 {
		size = cfg.CacheConfig.ElementSize
	}

	c, err := NewCache(size)

	if err != nil {
		return nil, err
	}

	s.cache = c

	dsn := cfg.SqlConfig.Username + ":" + cfg.SqlConfig.Password + "@tcp(" + cfg.SqlConfig.Url + `)/` + cfg.SqlConfig.Database
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal("failed to open db:", dsn)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect db:", err.Error())
		return nil, err
	}

	s.db = db
	return s, nil
}

func (s *Store) GetTrips(id string, date string, cache bool) (int, error) {
	//if read from cache, and there is no data
	//we will read from mysql and write back the data to cache

	//key in cache is  id:date string
	key := id + `:` + date

	if cache {
		if v, ok := s.cache.Get(key); ok {
			return v.(int), nil
		}
	}

	res := 0
	if rows, err := s.db.Query(querySql, id, date); err == nil {
		for rows.Next() {
			if err := rows.Scan(&res); err != nil {
				return 0, err
			}
			break
		}
	}

	s.cache.Add(key, res)
	return res, nil
}

func (s *Store) ClearCache(id *string) {
	if id == nil {
		s.cache.Purge()
	} else {

		//we remove all related cache for this id, which is id:date
		keys := s.cache.Keys()
		prefix := *id + `:`

		for _, v := range keys {
			if strings.HasPrefix(v.(string), prefix) {
				s.cache.Remove(v)
			}
		}
	}
}

func (s *Store) Len() int {
	return s.cache.Len()
}
