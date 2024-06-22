package requestresponse

import (
	"assignment/model"
)

type LocationUpdate struct {
	RiderID  string         `json:"rider_id"`
	Location model.Location `json:"location"`
}
