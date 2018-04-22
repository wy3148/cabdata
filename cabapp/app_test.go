package cabapp

import (
	"testing"
)

func TestNewCabDataApp(t *testing.T) {
	a := NewCabDataApp("../config/config.json")
	a.initAPI()
}
