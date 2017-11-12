package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

func roomWatcherKey(roomID int64, token string) string {
	return fmt.Sprintf("room_watcher:%d:%s", roomID, token)
}

func getWatcherCount(roomID int64) (int, error) {
	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", roomWatcherKey(roomID, "*")))
	return len(keys), err
}

func updateRoomWatcher(roomID int64, token string) error {
	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SETEX", roomWatcherKey(roomID, token), 3, token)
	return err
}
