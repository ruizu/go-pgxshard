package pgxshard

import (
	"context"
	"errors"
	"fmt"
	"hash/crc32"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var defaultShardIndexFunc = func(key any, numShards int) (int, error) {
	switch v := key.(type) {
	case int:
		return v % numShards, nil
	case int32:
		return int(v) % numShards, nil
	case int64:
		return int(v) % numShards, nil
	case string:
		return int(crc32.ChecksumIEEE([]byte(v))) % numShards, nil
	}

	return 0, errors.New("shard key type not supported")
}

type ShardManager struct {
	mu             sync.Mutex
	shards         []*pgxpool.Pool
	numShards      int
	shardIndexFunc func(key any, numShards int) (int, error)
}

func New(ctx context.Context, connectionStrings []string) (*ShardManager, error) {
	shards := make([]*pgxpool.Pool, len(connectionStrings))

	for i, connStr := range connectionStrings {
		db, err := pgxpool.New(ctx, connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to shard %d: %v", i, err)
		}
		shards[i] = db
	}

	return &ShardManager{
		shards:         shards,
		numShards:      len(shards),
		shardIndexFunc: defaultShardIndexFunc,
	}, nil
}

func (s *ShardManager) SetShardIndexFunc(ctx context.Context, f func(key any, count int) (int, error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shardIndexFunc = f
}

func (s *ShardManager) Shard(ctx context.Context, key any) (*pgxpool.Pool, error) {
	index, err := s.shardIndexFunc(key, s.numShards)
	if err != nil {
		return nil, err
	}

	if index < 0 || index > len(s.shards)-1 {
		return nil, fmt.Errorf("shard index %d is out of range", index)
	}

	return s.shards[index], nil
}

func (s *ShardManager) Shards(ctx context.Context) ([]*pgxpool.Pool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.shards, nil
}

func (s *ShardManager) Ping(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, shard := range s.shards {
		if err := shard.Ping(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (s *ShardManager) Close(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, shard := range s.shards {
		shard.Close()
	}

	return nil
}
