package sfomuseum

import (
	"fmt"
)

type Airport struct {
	WhosOnFirstId int64  `json:"wof:id"`
	Name          string `json:"wof:name"`
	SFOMuseumID   int64  `json:"sfomuseum:airport_id"`
	IATACode      string `json:"iata:code"`
	ICAOCode      string `json:"icao:code"`
	IsCurrent     int64  `json:"mz:is_current"`
}

func (a *Airport) String() string {
	return fmt.Sprintf("%s %s \"%s\" %d (%d)", a.IATACode, a.ICAOCode, a.Name, a.WhosOnFirstId, a.IsCurrent)
}
