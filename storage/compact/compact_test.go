package compact

import (
	"testing"
	"io/ioutil"
	"os"
)

var c *CompactCache
var tempf *os.File

func init() {
	dir, _ := ioutil.TempDir("", "CompactCache_testing")
	tempf, _ = ioutil.TempFile("", "CompactCache_testing")
	cfg := CompactCacheConfig{root: dir}
	c = NewCompactCache(cfg)
	c.Set("test", []byte{1})
}

func TestCompactCache_Set(t *testing.T) {
	err := c.Set("test1", []byte{1})
	if err != nil {
		t.Error(err)
	}
}

func TestCompactCache_Get(t *testing.T) {
	v, err := c.Get("test")
	if err != nil {
		t.Error(err)
	}
	if v[0] != 1 {
		t.Error("value != 1")
	}
}

func BenchmarkCompactCache_Get(b *testing.B) {
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

func BenchmarkCompactCache_GetParallel(b *testing.B) {
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

func BenchmarkCompactCache_SetAndGet(b *testing.B) {
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
