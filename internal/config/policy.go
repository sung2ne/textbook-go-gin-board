package config

import (
    "sync/atomic"

    "gorm.io/gorm"
)

// WeightedPolicy - 가중치 기반 정책
type WeightedPolicy struct {
    weights []int
    current uint64
}

func NewWeightedPolicy(weights []int) *WeightedPolicy {
    return &WeightedPolicy{weights: weights}
}

func (p *WeightedPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
    // 간단한 라운드 로빈 + 가중치
    total := 0
    for _, w := range p.weights {
        total += w
    }

    current := atomic.AddUint64(&p.current, 1)
    target := int(current % uint64(total))

    sum := 0
    for i, w := range p.weights {
        sum += w
        if target < sum {
            return connPools[i]
        }
    }

    return connPools[0]
}
