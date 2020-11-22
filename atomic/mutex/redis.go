package mutex

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

var ErrLockerBeReplace error = errors.New("locker has be replace")

type RedisMutex struct{
	conn redis.Conn
	key string
	uid string
}


func NewRedisMutex(conn redis.Conn, key string) *RedisMutex {
	uid, _ := uuid.NewUUID()
	return &RedisMutex{
		conn: conn,
		key:        key,
		uid:       uid.String(),
	}
}

func (r RedisMutex) Key() string{
	return r.key
}

func (r RedisMutex) Lock(expire int) bool{
	ok, err := redis.String(r.conn.Do("SET", r.key, r.uid, "EX", expire, "NX"))
	if err == nil && ok == "OK"{
		return true
	}
	return false
}

func (r RedisMutex) Unlock() (bool, error) {
	uid, err := redis.String(r.conn.Do("GET", r.key))
	if err == redis.ErrNil {
		return true, nil
	} else if err != nil {
		return false, err
	} else if uid != r.uid {
		return false, ErrLockerBeReplace
	}else{
		i, err := redis.Int(r.conn.Do("DEL", r.key))
		return i==1, err
	}
}

func (r RedisMutex) With(expire int, f func()) (getLock bool){
	getLock = r.Lock(expire)
	if !getLock{
		return
	}
	defer r.Unlock()
	f()
	return
}


