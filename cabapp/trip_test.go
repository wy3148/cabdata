package cabapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestCaculateTripsSingle(t *testing.T) {
	a := NewCabDataApp("../config/config.json")
	req := httptest.NewRequest("GET", "http://localhost:8080/trips?id=FFB1D86301D7CCD329B6AEFF09546C4D&date=2013-12-31", nil)
	w := httptest.NewRecorder()

	//from db
	expected := 21

	a.CaculateTrips(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("got the error response code, expecting 200, got:%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var trips []Trip
	err = json.Unmarshal(body, &trips)

	if err != nil {
		t.Fatal(err)
	}

	if len(trips) == 0 {
		t.Fatal("got the empty trips,should be one element")
	}

	if trips[0].Trips != expected {
		t.Fatalf("got wrong trips for id:%s,expected:%d, got:%d", trips[0].Id, expected, trips[0].Trips)
	}
}

func TestCaculateTripsMulitple(t *testing.T) {

	a := NewCabDataApp("../config/config.json")

	req := httptest.NewRequest("GET", "http://localhost:8080/trips?id=FFB1D86301D7CCD329B6AEFF09546C4D&id=F4FA02D140DE01950D4691AAFC9AAC8F&date=2013-12-31", nil)
	w := httptest.NewRecorder()

	//from db:21 trips
	key := []string{`FFB1D86301D7CCD329B6AEFF09546C4D`, `F4FA02D140DE01950D4691AAFC9AAC8F`}
	expected := []int{21, 0}

	a.CaculateTrips(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatalf("got the error response code, expecting 200, got:%d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Fatal(err)
	}

	var trips []Trip
	err = json.Unmarshal(body, &trips)

	if err != nil {
		t.Fatal(err)
	}

	if len(trips) != len(key) {
		t.Fatal("got the wrong trips, number of elements are wrong")
	}

	for _, v := range trips {
		for i := 0; i < len(key); i++ {

			if v.Id == key[i] {
				if expected[i] != v.Trips {
					t.Fatalf("got wrong trips for id:%s,expected:%d, got:%d", v.Id, expected[i], v.Trips)
				}
			}
		}
	}
}

func TestClearCache(t *testing.T) {
	a := NewCabDataApp("../config/config.json")

	req := httptest.NewRequest("GET", "http://localhost:8080/trips?id=FFB1D86301D7CCD329B6AEFF09546C4D&id=F4FA02D140DE01950D4691AAFC9AAC8F&date=2013-12-31", nil)
	w := httptest.NewRecorder()
	a.CaculateTrips(w, req)

	if a.store.Len() != 2 {
		t.Fatal("the number of elements in cache is wrong, expected 2 but got:", a.store.Len())
	}

	req = httptest.NewRequest("DELETE", "http://localhost:8080/trips/cache", nil)
	w = httptest.NewRecorder()
	a.ClearTripsCache(w, req)
	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Fatal("got wrong response code, expected 200, got:", resp.StatusCode)
	}

	if a.store.Len() != 0 {
		t.Fatal("cache is not cleared correctly, got:", a.store.Len())
	}
}
