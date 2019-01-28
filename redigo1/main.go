package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
)

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
	fmt.Printf("SET ret:%#v\n", ret)

	// GET location
	ret, err = c.Do("GET", "location", "me")
	if err != nil {
		log.Fatalf("Could not GET: %v\n", err)
	}
	fmt.Printf("GET ret:%s\n", ret)

}
