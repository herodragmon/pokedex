package main 

import (
	"time"
	"github.com/pokedex/internal/pokecache"
)

func main() {
  cfg := &config{
		Cache: pokecache.NewCache(5 * time.Second)
	}
  startRepl(cfg)
}
