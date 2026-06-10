package service

import (
	"fmt"
	"sync"
	"time"

	"qc-system/internal/store"
)

type LotGenerator struct {
	mu       sync.Mutex
	counters map[string]int
}

var lotGenerator = &LotGenerator{
	counters: make(map[string]int),
}

func GetLotGenerator() *LotGenerator {
	return lotGenerator
}

func (g *LotGenerator) GenerateLotNo(materialCode string) string {
	g.mu.Lock()
	defer g.mu.Unlock()

	yearMonth := time.Now().Format("200601")
	key := materialCode + "-" + yearMonth

	if _, exists := g.counters[key]; !exists {
		g.counters[key] = g.getMaxSequenceFromStore(materialCode, yearMonth)
	}

	g.counters[key]++

	return fmt.Sprintf("%s-%s-%04d", materialCode, yearMonth, g.counters[key])
}

func (g *LotGenerator) getMaxSequenceFromStore(materialCode string, yearMonth string) int {
	maxSeq := 0
	prefix := materialCode + "-" + yearMonth + "-"

	lots := store.GlobalStore.ListLots()
	for _, lot := range lots {
		if len(lot.LotNo) > len(prefix) && lot.LotNo[:len(prefix)] == prefix {
			seqStr := lot.LotNo[len(prefix):]
			var seq int
			fmt.Sscanf(seqStr, "%d", &seq)
			if seq > maxSeq {
				maxSeq = seq
			}
		}
	}

	return maxSeq
}
