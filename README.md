# Simple cabdata API service

###Last change

## Design

Features: <br />
1 Cabdata application provide restful API service. <br />
2 Users can call service to get trips information for NY drivers.<br />
3 Uses can choose to get those information from cache or from db directly<br />
4 Users can also clear the cached data<br />

Backend cache:<br />
1 So far cache is implemented with 'LRU' cache. To support restful API and go cocurrency,
We implement a thread-safe LRU cache with golang opensource code.<br />

2 We design system correctly so we could use other cache system like redis in future.<br />

## Install And the API description

This application is written in Golang, to get it, 

```
go get github.com/wy3148/cabdata

```

API:
'swagger_cabdata_api.json' in the project has described the APIs (swagger platform)


## Configuration file

The default configuration is located in './config/config.json', users have to configure the
mysql DSN settings before run the application

```
{
    "SqlConfig":{
        "Url":"localhost:3306",
        "Username":"test",
        "Password":"test",
        "Database":"test"
    },
    "ServerConfig":{
        "ServerUrl":"localhost:8080"
    },
    "CacheConfig":{
        "ElementSize":1000000
    }
}
```


## Run the application
```
go run main.go -config=./config/config.json
```
by default server is running on local, you will see similar output like following,
```
2018/04/22 18:20:05 start http server on: localhost:8080
```

## Use curl to test

Examples:

Get trips for a single id
```
curl -X GET 'http://localhost:8080/trips?id=FF2C42685FE5822F7A6DE63D32ED8193&date=2013-12-31&cache=true'
[{"Id":"FF2C42685FE5822F7A6DE63D32ED8193","Trips":20}]
```

Get trips for multiple ids
```
curl -X GET 'http://localhost:8080/trips?id=FF2C42685FE5822F7A6DE63D32ED8193&date=2013-12-31&id=FD631F3F8981584EA408223CA3BE6F26&cache=true'

[{"Id":"FF2C42685FE5822F7A6DE63D32ED8193","Trips":20},{"Id":"FD631F3F8981584EA408223CA3BE6F26","Trips":6}]
```

Clear the cache for a single id
```
curl -X DELETE 'http://localhost:8080/trips/cache?id=FF2C42685FE5822F7A6DE63D32ED8193'
```

Clear cache for mulitple id
```
curl -X DELETE 'http://localhost:8080/trips/cache?id=FF2C42685FE5822F7A6DE63D32ED8193&id=F4FA02D140DE01950D4691AAFC9AAC8F'
```

Clear all cached data
```
curl -X DELETE 'http://localhost:8080/trips/cache'
```
