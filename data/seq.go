package data

import (
	"context"
	"strconv"

	"github.com/redis/rueidis"
)

var (
	luaSeq = rueidis.NewLuaScript(`
        local r = redis.call('incr', KEYS[1])
        if (r < ARGV[1] + 1 or (ARGV[2] > ARGV[1] and r > ARGV[2] + 1)) then
          redis.call('set', KEYS[1], ARGV[1] + 1)
        end
        return redis.call('get', KEYS[1]) - 1
    `)
)

// Seq is a sequence generator.
type Seq interface {
	// Next returns the next value in the sequence.
	Next(key string, min int64, max int64) (int64, error)
}

type redisSeq struct {
	rdb rueidis.Client
}

// NewRedisSeq creates a new redis sequence generator.
func NewRedisSeq(rdb rueidis.Client) Seq {
	return &redisSeq{rdb: rdb}
}

// Next returns the next value in the sequence.
func (s *redisSeq) Next(key string, min int64, max int64) (int64, error) {
	if v, err := luaSeq.Exec(context.Background(), s.rdb, []string{key}, []string{strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)}).AsInt64(); err != nil {
		return 0, err
	} else {
		return v, nil
	}
}
