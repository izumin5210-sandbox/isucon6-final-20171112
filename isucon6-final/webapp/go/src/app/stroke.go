package main

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

func strokeIDsKey(roomID int64) string {
	return fmt.Sprintf("room:%d:stroke_ids", roomID)
}

func strokeKey(strokeID int64) string {
	return fmt.Sprintf("stroke:%d", strokeID)
}

func roomUpdatedAtKey() string {
	return "room_updated_at"
}

func appendStroke(roomID int64, s *Stroke) error {
	conn := pool.Get()
	defer conn.Close()

	if s.ID == 0 {
		s.ID = time.Now().UTC().UnixNano()
	}
	s.RoomID = roomID
	var err error
	_, err = conn.Do("ZADD", strokeIDsKey(roomID), s.ID, s.ID)
	if err != nil {
		return err
	}
	_, err = conn.Do("HMSET", strokeKey(s.ID),
		"id", s.ID,
		"roomID", s.RoomID,
		"width", s.Width,
		"red", s.Red,
		"green", s.Green,
		"blue", s.Blue,
		"alpha", s.Alpha,
	)
	if err != nil {
		return err
	}
	values := make([]float64, 0, len(s.Points)*2)
	for _, p := range s.Points {
		values = append(values, p.X, p.Y)
	}
	if len(s.Points) > 0 {
		err = appendPoints(s.ID, values...)
		if err != nil {
			return err
		}
	}
	_, err = conn.Do("ZADD", roomUpdatedAtKey(), time.Now().UTC().Unix(), roomID)
	if err != nil {
		return err
	}

	return nil
}

func getStrokes(roomID int64, greaterThanID int64) ([]*Stroke, error) {
	conn := pool.Get()
	defer conn.Close()

	idsKey := strokeIDsKey(roomID)
	i := 0
	var err error

	if greaterThanID > 0 {
		i, err = redis.Int(conn.Do("ZRANK", idsKey, greaterThanID))
		if err != nil {
			if err == redis.ErrNil {
				return []*Stroke{}, nil
			}
			return nil, err
		}
	}

	ids, err := redis.Int64s(conn.Do("ZRANGE", idsKey, i, -1))
	if err != nil {
		return nil, err
	}

	strokes := make([]*Stroke, 0, len(ids))
	for _, id := range ids {
		s := &Stroke{}
		reply, err := redis.Values(conn.Do("HGETALL", strokeKey(id)))
		if err != nil {
			return nil, err
		}
		err = redis.ScanStruct(reply, s)
		if err != nil {
			return nil, err
		}
		strokes = append(strokes, s)
	}

	return strokes, nil
}

func getStrokeCount(roomID int64) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("ZCARD", strokeIDsKey(roomID)))
}

func getRecentUpdatedRoomIDs() ([]int64, error) {
	conn := pool.Get()
	defer conn.Close()

	return redis.Int64s(conn.Do("ZREVRANGE", roomUpdatedAtKey(), 0, 100))
}
