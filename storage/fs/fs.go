package fs

import (
	"os"
	"path/filepath"
	"io/ioutil"
	"sync"
	"time"
	"errors"
)

func openFile(p string, flags int) *os.File {
	f, err := os.OpenFile(p, flags, 0644)
	if err != nil {
		panic(err.Error())
	}
	return f
}

type FSCache struct {
	cfg    FSCacheConfig
	rwLock sync.RWMutex
}

type FSCacheConfig struct {
	root string
}

func NewFSCache(cfg FSCacheConfig) *FSCache {
	return &FSCache{cfg:cfg}
}

func (c *FSCache) Get(key string) ([]byte, error) {
	fpath := filepath.Join(c.cfg.root, key)
	f := openFile(fpath, os.O_RDONLY)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *FSCache) Set(key string, value []byte) error {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	fpath := filepath.Join(c.cfg.root, key)
	return ioutil.WriteFile(fpath, value, 0644)
}

func (c *FSCache) SetWithMaxAge(key string, value []byte, duration time.Duration) error {
	return errors.New("not implemented")
}

func (c *FSCache) Remove(key string) error {
	fpath := filepath.Join(c.cfg.root, key)
	return os.Remove(fpath)
}