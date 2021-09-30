package sfomuseum

import (
	"fmt"
)

type Aircraft struct {
	WhosOnFirstId  int64  `json:"wof:id"`
	Name           string `json:"wof:name"`
	SFOMuseumId    int64  `json:"sfomuseum:aircraft_id"`
	ICAODesignator string `json:"icao:designator,omitempty"`
	WikidataId     string `json:"wd:id,omitempty"`
	IsCurrent      int64  `json:"mz:is_current"`
}

func (a *Aircraft) String() string {
	return fmt.Sprintf("%s \"%s\" %d (%d) Is current: %d", a.ICAODesignator, a.Name, a.WhosOnFirstId, a.SFOMuseumId, a.IsCurrent)
}
