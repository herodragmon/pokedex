package main 

import (
	"time"
	"github.com/herodragmon/pokedex/internal/pokecache"
)

func main() {
  cfg := &config{
		Cache: pokecache.NewCache(5 * time.Second),
	}
  startRepl(cfg)
}
