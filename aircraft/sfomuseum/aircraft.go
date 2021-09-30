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
	IsCurrent      int64  `json:"mz:is_current"`
}

func (a *Aircraft) String() string {
	return fmt.Sprintf("%d %s \"%s\" (%d)", a.WhosOnFirstId, a.ICAODesignator, a.Name, a.IsCurrent)
}
