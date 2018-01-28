package compact

import (
	"github.com/edsrzf/mmap-go"
	"os"
	"sync"
	"path/filepath"
	"strconv"
	"time"
	"runtime"
)

func openFile(p string, flags int) *os.File {
	f, err := os.OpenFile(p, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

type CompactCache struct {
	cfg    CompactCacheConfig
	blockf *os.File
	blockm *mmap.MMap
	rwLock sync.RWMutex
}

type CompactCacheConfig struct {
	root string
}

func NewCompactCache(cfg CompactCacheConfig) *CompactCache {
	fpath := filepath.Join(cfg.root, "block")
	f := openFile(fpath, os.O_RDWR | os.O_CREATE | os.O_TRUNC)
	f.Truncate(1e7)
	f.Sync()
	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic(err.Error())
	}
	c := &CompactCache{cfg:cfg, blockf: f, blockm: &m}
	runtime.SetFinalizer(c, (*c).Close)
	return c
}

func (c *CompactCache) Get(key string) ([]byte, error) {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	k, _ := strconv.ParseInt(key, 10, 32)
	return []byte{(*c.blockm)[int(k)]}, nil
}

func (c *CompactCache) Set(key string, value []byte) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	k, _ := strconv.ParseInt(key, 10, 32)
	(*c.blockm)[int(k)] = value[0]
	return nil
}

func (c *CompactCache) SetWithMaxAge(key string, value []byte, duration time.Duration) error {
	k, _ := strconv.ParseInt(key, 10, 32)
	(*c.blockm)[int(k)] = value[0]
	return nil
}

func (c *CompactCache) Close(p *CompactCache) {
	if p != nil {
		p.blockm.Unmap()
		p.blockm.Flush()
		p.blockf.Close()
	} else {
		c.blockm.Unmap()
		c.blockm.Flush()
		c.blockf.Close()
	}
}
