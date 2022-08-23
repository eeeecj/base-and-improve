package redis

import (
	"context"
	"errors"
	"github.com/balance/boolfilter"
	"github.com/balance/boolfilter/global"
	"github.com/balance/boolfilter/memeory"
	redis "github.com/go-redis/redis/v8"
)

func NewCacheFilter(ctx context.Context, hashTableName, key string, cli *redis.Client, blen uint64, hashes ...global.HashFunc) (boolfilter.Adapter, error) {
	c := &CacheFilter{
		Filter: &memeory.Filter{
			Bits:      make([]byte, blen, blen),
			Hashes:    hashes,
			IsChanged: false,
		},
		hashTableName: hashTableName,
		key:           key,
		cli:           cli,
		ctx:           ctx,
	}
	err := c.Init(c.ctx, c.cli, c.hashTableName, c.key)
	if err != nil {
		return nil, err
	}
	return c, nil
}

type CacheFilter struct {
	*memeory.Filter
	hashTableName string
	key           string
	cli           *redis.Client
	ctx           context.Context
}

func (f *CacheFilter) Clear() error {
	err := f.cli.Del(f.ctx, f.key).Err()
	if err != nil {
		return err
	}
	return f.Filter.Clear()
}

func (f *CacheFilter) Init(ctx context.Context, cli *redis.Client, hashTableName string, key string) error {
	f.cli = cli
	f.ctx = ctx
	f.key = key
	f.hashTableName = hashTableName
	cmd := f.cli.HGet(f.ctx, f.hashTableName, f.key)
	if err := cmd.Err(); err != nil {
		if err == redis.Nil {
			f.cli.HSet(f.ctx, f.key, f.Bits)
		} else {
			return err
		}
	} else {
		b, _ := cmd.Bytes()
		if len(b) == len(f.Bits) {
			f.Bits = b
		} else {
			return errors.New("length not match")
		}
	}
	return nil
}

func (f *CacheFilter) Write() error {
	if f.IsChanged {
		cmd := f.cli.HSet(f.ctx, f.hashTableName, f.key, f.Bits)
		if err := cmd.Err(); err != nil {
			return err
		}
		f.IsChanged = false
	}
	return nil
}

func (f *CacheFilter) Close() error {
	return f.cli.Close()
}
