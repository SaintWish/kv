package kv1

import (
	"fmt"
	"testing"
	"time"
)

func TestSetGet_KeyString(t *testing.T) {
	cache := New[string, string](time.Minute, 2048, 32)
	cache.Set("unicorns", "are cool")

	if res := cache.Get("unicorns"); res != "are cool" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", res, "are cool")
	}
}

func TestSetGet_KeyInt(t *testing.T) {
	cache := New[int, string](time.Minute, 2048, 32)
	cache.Set(1337, "leet haxiors")

	if res := cache.Get(1337); res != "leet haxiors" {
		t.Errorf("Result was incorrect, got: %s, want: %s.", res, "leet haxiors")
	}
}

func TestFlush(t *testing.T) {
	cache := New[int, string](time.Minute, 2048, 32)
	cache.SetOnEvicted(func(k int, v string){
		fmt.Println(k)
	})

	cache.Set(1337, "leet haxiors")
	cache.Set(1338, "leet haxiors1")
	cache.Set(3434, "leet haxiors2")
	cache.Set(5465, "leet haxiors3")

	cache.Flush()
}