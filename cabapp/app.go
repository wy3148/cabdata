package cabapp

import (
	"github.com/gorilla/mux"
	"github.com/wy3148/cabdata/config"
	"github.com/wy3148/cabdata/store"
	"log"
	"net/http"
	"os"
)

type CabDataApp struct {
	r     *mux.Router
	store store.StoreInterface
	cfg   *config.AppCfg
}

func NewCabDataApp(configFile string) *CabDataApp {

	f, err := os.Open(configFile)

	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig(f)

	s, err := store.InitStoreInterface(cfg)

	if err != nil {
		panic(err)
	}

	app := &CabDataApp{cfg: cfg, store: s}
	app.initAPI()
	return app
}

func (app *CabDataApp) initAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/trips", app.CaculateTrips).Methods("GET")
	r.HandleFunc("/trips/cache", app.ClearTripsCache).Methods("DELETE")
	app.r = r
}

func (app *CabDataApp) Start() {
	log.Println("start http server on:", app.cfg.ServerConfig.ServerUrl)
	log.Fatal(http.ListenAndServe(app.cfg.ServerConfig.ServerUrl, app.r))
}
