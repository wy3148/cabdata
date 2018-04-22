package cabapp

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Result       string
	ErrorMessage string `json:",omitempty"`
}

type Trip struct {
	Id           string
	Trips        int
	Result       string `json:",omitempty"`
	ErrorMessage string `json:",omitempty"`
}

//CaculateTrips get the number of trips
func (app *CabDataApp) CaculateTrips(w http.ResponseWriter, r *http.Request) {

	cache := false
	date := r.URL.Query().Get("date")

	if len(date) == 0 {
		log.Println("date is missing")
		w.WriteHeader(http.StatusBadRequest)
		r := &Response{ErrorMessage: "date is missing", Result: "failure"}
		b, _ := json.Marshal(r)
		w.Write(b)
		return
	}

	ids := r.URL.Query()["id"]

	if len(ids) == 0 {
		log.Println("id is missing")
		w.WriteHeader(http.StatusBadRequest)
		r := &Response{ErrorMessage: "id is missing", Result: "failure"}
		b, _ := json.Marshal(r)
		w.Write(b)
		return
	}

	readCache := r.URL.Query().Get("cache")

	if readCache == "true" {
		cache = true
	}

	var res []Trip

	for _, id := range ids {
		r, err := app.store.GetTrips(id, date, cache)

		if err != nil {
			log.Printf("failed to data from cache:%s\n", err.Error())
			t := Trip{Id: id, Result: "failure", ErrorMessage: err.Error()}
			res = append(res, t)
		} else {
			t := Trip{Id: id, Trips: r}
			res = append(res, t)
		}
	}

	b, err := json.Marshal(res)

	if err != nil {
		log.Println("faile to encode the object:", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		r := &Response{ErrorMessage: err.Error(), Result: "failure"}
		b, _ := json.Marshal(r)
		w.Write(b)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

//ClearTripsCache clear the cache
func (app *CabDataApp) ClearTripsCache(w http.ResponseWriter, r *http.Request) {

	ids := r.URL.Query()["id"]

	if len(ids) == 0 {
		app.store.ClearCache(nil)
	} else {
		for _, id := range ids {
			idV := id
			app.store.ClearCache(&idV)
		}
	}

	w.WriteHeader(http.StatusOK)
	res := Response{Result: "success"}
	b, _ := json.Marshal(res)
	w.Write(b)
	return
}
