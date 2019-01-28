package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

type ObjectType int

const (
	GeometryObject = iota
	FeatureObject
	FeatureCollectionObject
)

/*
type GeoJsonType int
const (
	Point = iota
	MultiPoint
	LineString
	MultiLineString
	Polygon
	MultiPolygon
	GeometryCollection
	Feature
	FeatureCollection
)
*/

type GeoJsonMember struct {
	ObjectType        ObjectType        `json:"-"`
	Type              string            `json:"type"`
	CoordinatesRaw    json.RawMessage   `json:"coordinates,omitempty"`
	CoordinatesObject interface{}       `json:"-"`
	Geometry          json.RawMessage   `json:"geometry,omitempty"`
	Properties        map[string]string `json:"properties,omitempty"`
}

type Point [2]float64

type LineString []Point

type Polygon []LineString

func NewGeoJsonMember(b []byte) (*GeoJsonMember, error) {
	var member GeoJsonMember
	err := json.Unmarshal(b, &member)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal error: %v", err)
	}
	err = member.setObjectType()
	if err != nil {
		return nil, err
	}
	if member.ObjectType == GeometryObject {
		err := member.setCoordinatesObject()
		if err != nil {
			return nil, err
		}
	}
	return &member, nil
}

func NewGeoJsonMembers(b []byte) ([]*GeoJsonMember, error) {
	var members []*GeoJsonMember
	err := json.Unmarshal(b, &members)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		err := member.setObjectType()
		if err != nil {
			return nil, fmt.Errorf("%v:%v", err, member)
		}
		if member.ObjectType == GeometryObject {
			err := member.setCoordinatesObject()
			if err != nil {
				return nil, err
			}
		}
	}
	return members, nil
}

func (member *GeoJsonMember) setObjectType() error {
	switch member.Type {
	case "Point", "LineString", "Polygon":
		member.ObjectType = GeometryObject
	case "Feature":
		member.ObjectType = FeatureObject
	case "FeatureCollection":
		member.ObjectType = FeatureCollectionObject
	default:
		return fmt.Errorf("Unknown type: %v", member.Type)
	}
	return nil
}

func (member *GeoJsonMember) setCoordinatesObject() error {
	var object interface{}
	switch member.Type {
	case "Point":
		object = new(Point)
	case "LineString":
		object = new(LineString)
	case "Polygon":
		object = new(Polygon)
	default:
		return fmt.Errorf("Unknown type: %v", member.Type)
	}
	err := json.Unmarshal(member.CoordinatesRaw, &object)
	if err != nil {
		return fmt.Errorf("Unmarshal error:%v coordinates:%s", err, member.CoordinatesRaw)
	}
	//fmt.Printf("object:%v\n", object)
	member.CoordinatesObject = object
	return nil
}

func (member *GeoJsonMember) String() string {
	switch member.ObjectType {
	case GeometryObject:
		return fmt.Sprintf("type:%v coordinates:%v", member.Type, member.CoordinatesObject)
	case FeatureObject:
		geometry, err := NewGeoJsonMember(member.Geometry)
		if err != nil {
			return fmt.Sprintf("GeoJsonMember.String() NewGeoJsonMember error:%v geometry:%s", err, member.Geometry)
		}
		return fmt.Sprintf("type:%v geometry:%v properties:%v", member.Type, geometry, member.Properties)
	default:
		return "Unknown Object Type"
	}
}

func unmarshalGeoJson() {
	var shapes = []byte(`[
	{"type": "Point", "coordinates": [1.23, 4.56]},
	{"type": "LineString", "coordinates": [[1.23, 4.56],[7.89,10.12]]},
	{"type": "Polygon", "coordinates": [[[1.23, 4.56],[7.89,10.12],[3.45,6.78],[1.23,4.56]]]},
	{"type": "Feature", "geometry": {"type": "Point", "coordinates": [1.23, 4.56]}, "properties": {"name": "point:A"}}
	]`)

	members, err := NewGeoJsonMembers(shapes)
	if err != nil {
		fmt.Printf("err:%v\n", err)
	}

	for i, m := range members {
		fmt.Printf("member[%v]:%v\n", i, m)
	}
}

func main() {
	// db connect
	c, err := redis.Dial("tcp", ":9851")
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	defer c.Close()

	// SET location
	ret, err := c.Do("SET", "location", "me", "POINT", 35.6581, 139.6975)
	if err != nil {
		log.Fatalf("Could not SET: %v\n", err)
	}
	fmt.Printf("SET ret:%s\n", ret)

	// GET location
	ret, err = c.Do("GET", "location", "me")
	if err != nil {
		log.Fatalf("Could not GET: %v\n", err)
	}
	fmt.Printf("GET ret:%s\n", ret)

	unmarshalGeoJson()
}
