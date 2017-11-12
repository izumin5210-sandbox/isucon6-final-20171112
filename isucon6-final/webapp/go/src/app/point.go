package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func pointsKey(strokeID int64) string {
	return fmt.Sprintf("stroke:%d:points", strokeID)
}

func appendPoints(strokeID int64, values ...float64) error {
	conn := pool.Get()
	defer conn.Close()

	args := make([]interface{}, 0, len(values)+1)
	args = append(args, pointsKey(strokeID))
	for _, v := range values {
		args = append(args, fmt.Sprint(v))
	}

	_, err := conn.Do("RPUSH", args...)
	return err
}

func getStrokePoints(strokeID int64) ([]*Point, error) {
	conn := pool.Get()
	defer conn.Close()

	values, err := redis.Float64s(conn.Do("LRANGE", pointsKey(strokeID), 0, -1))
	if err != nil {
		return nil, err
	}
	n := len(values) / 2
	ps := make([]*Point, 0, n)
	for i := 0; i < n; i++ {
		ps[i] = &Point{
			X:        values[2*i],
			Y:        values[2*i+1],
			StrokeID: strokeID,
		}
	}
	return ps, nil
}
