package main

import (
	"encoding/json"
	"log"
)

type GeoJsonMember struct {
	Type           string          `json:"type"`
	CoordinatesRaw json.RawMessage `json:"coordinates,omitempty"`
}

type Point [2]float64

type LineString []Point

type Polygon []LineString

func main() {
	var shapes = []byte(`[
	{"type": "Point", "coordinates": [1.23, 4.56]},
	{"type": "LineString", "coordinates": [[1.23, 4.56],[7.89,10.12]]},
	{"type": "Polygon", "coordinates": [[[1.23, 4.56],[7.89,10.12],[3.45,6.78],[1.23,4.56]]]}
	]`)

	var members []*GeoJsonMember
	err := json.Unmarshal(shapes, &members)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}

	for i, member := range members {
		var coordinates interface{}
		switch member.Type {
		case "Point":
			coordinates = new(Point)
		case "LineString":
			coordinates = new(LineString)
		case "Polygon":
			coordinates = new(Polygon)
		default:
			log.Fatalf("Unknown type: %v", member.Type)
		}
		err := json.Unmarshal(member.CoordinatesRaw, &coordinates)
		if err != nil {
			log.Fatalf("Unmarshal error:%v coordinates:%s", err, member.CoordinatesRaw)
		}
		log.Printf("member[%v]:type:%v coordinates:%v\n", i, member.Type, coordinates)
	}
}
