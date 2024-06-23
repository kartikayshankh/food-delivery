package requestresponse

import (
	"github.com/kartikayshankh/food-delivery/model"
)

type LocationUpdate struct {
	RiderID  string         `json:"rider_id"`
	Location model.Location `json:"location"`
}
