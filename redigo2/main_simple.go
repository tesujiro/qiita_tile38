package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type GeoJsonMember struct {
	Type           string          `json:"type"`
	CoordinatesRaw json.RawMessage `json:"coordinates,omitempty"`
	CoordinatesObj interface{}     `json:"-"`
}

type Point [2]float64

type LineString []Point

type Polygon []LineString

func (member *GeoJsonMember) String() string {
	return fmt.Sprintf("%s %v", member.Type, member.CoordinatesObj)
}

func (member *GeoJsonMember) setCoordinates() error {
	var coordinates interface{}
	switch member.Type {
	case "Point":
		coordinates = new(Point)
	case "LineString":
		coordinates = new(LineString)
	case "Polygon":
		coordinates = new(Polygon)
	default:
		return fmt.Errorf("Unknown type: %v", member.Type)
	}
	err := json.Unmarshal(member.CoordinatesRaw, &coordinates)
	if err != nil {
		return fmt.Errorf("json.Unmarshal error: %v", err)
	}
	member.CoordinatesObj = coordinates
	return nil
}

func unmarshalMultiResults(shapes []byte) ([]*GeoJsonMember, error) {
	var members []*GeoJsonMember
	err := json.Unmarshal(shapes, &members)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v", err)
	}

	for i, member := range members {
		err := member.setCoordinates()
		if err != nil {
			return nil, fmt.Errorf("member[%v]:type:%v coordinates:%v err:%v\n", i, member.Type, member.CoordinatesRaw, err)
		}
	}
	return members, nil
}

func unmarshalSingleResult(shapes []byte) (*GeoJsonMember, error) {
	var member GeoJsonMember
	err := json.Unmarshal(shapes, &member)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v", err)
	}

	err = member.setCoordinates()
	if err != nil {
		return nil, fmt.Errorf("type:%v coordinates:%v err:%v\n", member.Type, member.CoordinatesRaw, err)
	}
	return &member, nil
}

func main() {
	// db connect
	c, err := redis.Dial("tcp", ":9851")
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer c.Close()

	// SET fleet
	for _, data := range [][]interface{}{
		{"fleet", "id1", "FIELD", "start", "123456", "FIELD", "end", "789012", "POINT", 35.6581, 139.6975},
		{"fleet", "id2", "OBJECT", `{"type":"Point","coordinates":[139.6975,35.6581]}`},
		{"fleet", "id3", "OBJECT", `{"type":"LineString","coordinates":[[139.6975,35.6581],[1,1],[2,2]]}`},
		{"fleet", "id4", "POINT", 35.6581, 139.6975},
	} {
		ret, err := c.Do("SET", data...)
		if err != nil {
			log.Fatalf("Could not SET: %v\n", err)
		}
		fmt.Printf("SET ret:%#s\n", ret)
	}

	// SCAN fleet
	results, err := redis.Values(c.Do("SCAN", "fleet"))
	if err != nil {
		log.Fatalf("Could not SCAN: %v\n", err)
	}

	var cursor int
	var members []interface{}
	_, err = redis.Scan(results, &cursor, &members)
	if err != nil {
		fmt.Printf("scan result error: %v", err)
		return
	}

	for len(members) > 0 {
		// pick up one record from results as []interface{}
		var object []interface{}
		members, err = redis.Scan(members, &object)
		if err != nil {
			fmt.Printf("scan record error: %v", err)
			return
		}
		// scan columns from one record -> [id,json],fields
		var id []byte
		var json []byte
		others, err := redis.Scan(object, &id, &json)
		if err != nil {
			fmt.Printf("scan columns error: %v", err)
			return
		}

		// unmarshal geojson string to struct
		gjm, err := unmarshalSingleResult(json)
		if err != nil {
			fmt.Printf("unmarshal json error: %v", err)
			return
		}
		fmt.Printf("id:%s  json:%s  others:%s\n", id, gjm, others)
	}
}
