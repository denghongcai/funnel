package fs

import (
	"testing"
	"io/ioutil"
	"os"
)

var c *FSCache
var tempf *os.File

func init() {
	dir, _ := ioutil.TempDir("", "fscache_testing")
	tempf, _ = ioutil.TempFile("", "fscache_testing")
	cfg := FSCacheConfig{root: dir}
	c = NewFSCache(cfg)
	c.Set("test", []byte{1})
}

func TestFSCache_Set(t *testing.T) {
	err := c.Set("test1", []byte{1})
	if err != nil {
		t.Error(err)
	}
}

func TestFSCache_Get(t *testing.T) {
	v, err := c.Get("test")
	if err != nil {
		t.Error(err)
	}
	if v[0] != 1 {
		t.Error("value != 1")
	}
}

func BenchmarkFSCache_Get(b *testing.B) {
	for n := 0; n < b.N; n++ {
		v, err := c.Get("test")
		if err != nil {
			b.Error(err)
		}
		if v[0] != 1 {
			b.Error("value != 1")
		}
	}
}

func BenchmarkFSCache_GetParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			v, err := c.Get("test")
			if err != nil {
				b.Error(err)
			}
			if v[0] != 1 {
				b.Error("value != 1")
			}
		}
	})
}

func BenchmarkFSCache_SetAndGet(b *testing.B) {
	for n := 0; n < b.N; n++ {
		err := c.Set("test", []byte{1})
		if err != nil {
			b.Error(err)
		}
		v, err := c.Get("test")
		if err != nil {
			b.Error(err)
		}
		if v[0] != 1 {
			b.Error("value != 1")
		}
	}
}

func BenchmarkFS_Get(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := ioutil.ReadAll(tempf)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkFS_GetParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := ioutil.ReadAll(tempf)
			if err != nil {
				b.Error(err)
			}
		}
	})
}
