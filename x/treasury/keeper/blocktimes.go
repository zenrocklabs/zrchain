package keeper

import (
	"context"
	"encoding/binary"
	"fmt"
	"strconv"

	"time"

	"cosmossdk.io/core/store"
)

const (
	blockWindow                   = 20
	defaultBlockTime      float64 = 5
	lastBlockTimestampKey         = "lastBlockTimestamp"
	blockTimeKeyPrefix            = "blockTimeKey"
)

func getBlockTimes(store store.KVStore) ([]time.Duration, error) {
	var blockTimes []time.Duration
	for i := 0; i < blockWindow; i++ {
		btBytes, err := store.Get([]byte(blockTimeKeyPrefix + strconv.Itoa(i)))
		if err != nil {
			return nil, err
		}

		if btBytes != nil {
			blockTimes = append(blockTimes, time.Duration(binary.LittleEndian.Uint64(btBytes)))
		}
	}

	return blockTimes, nil
}

func setBlockTimes(store store.KVStore, blockTimes []time.Duration) error {
	for i, blockTime := range blockTimes {
		bt := make([]byte, 8)
		binary.LittleEndian.PutUint64(bt, uint64(blockTime))
		if err := store.Set([]byte(blockTimeKeyPrefix+strconv.Itoa(i)), bt); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) AddBlockTime(ctx context.Context) error {
	now := time.Now()

	var (
		lastBlockTime time.Time
		blockTimes    []time.Duration
	)

	s := k.memStore.OpenMemoryStore(ctx)
	blockTimes, err := getBlockTimes(s)
	if err != nil {
		return err
	}

	lastBlockTimeBytes, err := s.Get([]byte(lastBlockTimestampKey))
	if err != nil {
		return err
	}

	if lastBlockTimeBytes != nil {
		err = lastBlockTime.UnmarshalBinary(lastBlockTimeBytes)
		if err != nil {
			return err
		}
	}

	if !lastBlockTime.IsZero() {
		lastDuration := now.Sub(lastBlockTime)

		blockTimes = append(blockTimes, lastDuration)
		if len(blockTimes) > blockWindow {
			blockTimes = blockTimes[1:]
		}

		err = setBlockTimes(s, blockTimes)
		if err != nil {
			return err
		}
	}

	ts, _ := now.MarshalBinary()
	err = s.Set([]byte(lastBlockTimestampKey), ts)
	if err != nil {
		return err
	}

	fmt.Printf("avg block time: %f\n", k.AverageBlockTime(ctx))
	return nil
}

func (k Keeper) AverageBlockTime(ctx context.Context) float64 {
	s := k.memStore.OpenMemoryStore(ctx)
	blockTimes, err := getBlockTimes(s)
	if err != nil {
		k.Logger().Error("error getting block times", err)
	}

	n := len(blockTimes)
	if n == 0 {
		bti, err := k.mintKeeper.GetDefaultBlockTime(ctx)
		bt := float64(bti)
		if err != nil {
			k.Logger().Error("Get default block time", err)
			bt = defaultBlockTime
		}
		return bt
	}

	var sum float64
	for _, bt := range blockTimes {
		sum += bt.Seconds()
	}

	blockTime := sum / float64(n)

	return blockTime
}
