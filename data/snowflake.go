package data

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/choral-io/gommerce-server-core/config"
)

const (
	DefaultIdEpoch = int64(1704067200000) // Defaults to: 2024-01-01T00:00:00Z
)

var (
	defaultIdWorker atomic.Value

	ErrIdEpochOutOfRange       = errors.New("the value of 'idEpoch' must be greater than 0")
	ErrClusterIdBitsOutOfRange = errors.New("the value of 'clusterIdBits' must be greater than 0")
	ErrWorkerIdBitsOutOfRange  = errors.New("the value of 'workerIdBits' must be greater than 0")
	ErrSequenceBitsOutOfRange  = errors.New("the value of 'sequenceBits' must be greater than 0")
	ErrClusterIdOutOfRange     = errors.New("the value of 'clusterId' out of range")
	ErrWorkerIdOutOfRange      = errors.New("the value of 'workerId' out of range")
	ErrTimeMilliBitsOutOfRange = errors.New("the sum of 'clusterIdBits', 'workerIdBits' and 'sequenceBits' must be less than 23")
)

func init() {
	var idw IdWorker = &idWorker{
		idEpoch:       DefaultIdEpoch,
		clusterId:     0,
		workerId:      0,
		clusterIdBits: 5,
		workerIdBits:  5,
		sequenceBits:  12,
		sequenceMask:  4095, // int64(1)<<12 - 1
		sequenceValue: 0,
		lastMillis:    0,
	}
	defaultIdWorker.Store(idw)
}

// SetDefaultIdWorker sets the default IdWorker instance.
func SetDefaultIdWorker(w IdWorker) {
	defaultIdWorker.Store(w)
}

// DefaultIdWorker returns the default IdWorker instance.
func DefaultIdWorker() IdWorker {
	return defaultIdWorker.Load().(IdWorker)
}

// IdWorker is used to generate unique id.
type IdWorker interface {
	NextInt64() int64
	NextBytes() [8]byte
	NextHex() string
}

// idWorker is used to generate unique id.
type idWorker struct {
	sync.Mutex
	idEpoch       int64
	clusterId     int64
	workerId      int64
	clusterIdBits int32
	workerIdBits  int32
	sequenceBits  int32
	sequenceMask  int64
	sequenceValue int64
	lastMillis    int64
}

// NextInt64 returns the next unique id in int64.
func (w *idWorker) NextInt64() int64 {
	w.Lock()
	defer w.Unlock()
	nextMillis := time.Now().UnixMilli()
	if w.lastMillis == nextMillis {
		w.sequenceValue = (w.sequenceValue + 1) & w.sequenceMask
		if w.sequenceValue == 0 {
			nextMillis = time.Now().UnixMilli()
			for w.lastMillis > nextMillis {
				nextMillis = time.Now().UnixMilli()
			}
		}
	} else {
		w.sequenceValue = 0
	}
	w.lastMillis = nextMillis
	slog.Debug("generating new snowflake id", slog.Int64("time.millis", nextMillis), slog.Int64("seq.value", w.sequenceValue))
	return ((nextMillis - w.idEpoch) << int64(w.clusterIdBits+w.workerIdBits+w.sequenceBits)) |
		(w.clusterId << int64(w.workerIdBits+w.sequenceBits)) |
		(w.workerId << int64(w.sequenceBits)) |
		w.sequenceValue
}

// NextBytes returns the next unique id in bytes.
func (w *idWorker) NextBytes() [8]byte {
	v := w.NextInt64()
	return [8]byte{byte(v >> 56), byte(v >> 48), byte(v >> 40), byte(v >> 32), byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

// NextHex returns the next unique id in hex string.
func (w *idWorker) NextHex() string {
	return fmt.Sprintf("%x", w.NextBytes())
}

// NewIdWorker creates a new IdWorker instance with the given config.
// If the workerSeqKey is not empty, the workerId will be generated from the Seq.
func NewIdWorker(cfg config.SnowflakeConfig, seq Seq) (IdWorker, error) {
	idEpoch := cfg.GetIdEpoch()
	clusterId := cfg.GetClusterId()
	workerId := cfg.GetWorkerId()
	clusterIdBits := cfg.GetClusterIdBits()
	workerIdBits := cfg.GetWorkerIdBits()
	sequenceBits := cfg.GetSequenceBits()

	workerSeqKey := cfg.GetWorkerSeqKey()
	if workerSeqKey != "" {
		if seq == nil {
			return nil, errors.New("redis client is not initialized")
		}
		if v, err := seq.Next(workerSeqKey, int64(0), int64(1)<<cfg.GetWorkerIdBits()-1); err != nil {
			return nil, err
		} else {
			workerId = v
		}
	}

	if idEpoch < 0 {
		return nil, ErrIdEpochOutOfRange
	}
	if clusterIdBits <= 0 {
		return nil, ErrClusterIdBitsOutOfRange
	}
	if workerIdBits <= 0 {
		return nil, ErrWorkerIdBitsOutOfRange
	}
	if sequenceBits <= 0 {
		return nil, ErrSequenceBitsOutOfRange
	}
	if clusterIdBits+workerIdBits+sequenceBits >= 23 {
		return nil, ErrTimeMilliBitsOutOfRange
	}
	if m := int64(1)<<clusterIdBits - 1; clusterId < 0 || clusterId > m {
		return nil, fmt.Errorf("the value of clusterId must be in the range 0 to %d: %w", m, ErrClusterIdOutOfRange)
	}
	if m := int64(1)<<workerIdBits - 1; workerId < 0 || workerId > m {
		return nil, fmt.Errorf("the value of workerId must be in the range 0 to %d: %w", m, ErrWorkerIdOutOfRange)
	}

	return &idWorker{
		idEpoch:       idEpoch,
		clusterId:     clusterId,
		workerId:      workerId,
		clusterIdBits: clusterIdBits,
		workerIdBits:  workerIdBits,
		sequenceBits:  sequenceBits,
		sequenceMask:  int64(1)<<sequenceBits - 1,
		sequenceValue: 0,
		lastMillis:    0,
	}, nil
}
