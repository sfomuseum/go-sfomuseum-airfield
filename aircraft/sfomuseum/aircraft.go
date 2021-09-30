package sfomuseum

import (
	"fmt"
)

type Aircraft struct {
	WhosOnFirstId  int64  `json:"wof:id"`
	Name           string `json:"wof:name"`
	SFOMuseumID    int64  `json:"sfomuseum:aircraft_id"`
	ICAODesignator string `json:"icao:designator,omitempty"`
	WikidataID     string `json:"wd:id,omitempty"`
}

func (a *Aircraft) String() string {
	return fmt.Sprintf("%d %s \"%s\"", a.WhosOnFirstId, a.ICAODesignator, a.Name)
}
