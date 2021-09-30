package sfomuseum

import (
	"fmt"
)

type Airline struct {
	WhosOnFirstId int64  `json:"wof:id"`
	Name          string `json:"wof:name"`
	SFOMuseumID   int64  `json:"sfomuseum:airline_id"`
	IATACode      string `json:"iata:code,omitempty"`
	ICAOCode      string `json:"icao:code,omitempty"`
	ICAOCallsign  string `json:"icao:callsign,omitempty"`
	WikidataID    string `json:"wd:id,omitempty"`
}

func (a *Airline) String() string {
	return fmt.Sprintf("%s %s %s \"%s\" %d", a.IATACode, a.ICAOCode, a.ICAOCallsign, a.Name, a.WhosOnFirstId)
}
