/*
 * MIT License
 *
 * Copyright (c) 2020 Dan <https://github.com/delivrance>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package limiter

import (
	"sort"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type Limit struct {
	Amount int
	Period time.Duration
}

type Limiter struct {
	limiters map[string][]*rate.Limiter
	limits   []Limit
	mu       *sync.Mutex
}

func NewLimiter(limits ...Limit) *Limiter {
	sort.SliceStable(limits, func(i, j int) bool {
		a := limits[i].Period * time.Duration(limits[i].Amount)
		b := limits[j].Period * time.Duration(limits[j].Amount)

		return a < b
	})

	return &Limiter{
		limiters: make(map[string][]*rate.Limiter),
		limits:   limits,
		mu:       &sync.Mutex{},
	}
}

func (lim *Limiter) add(key string) {
	for _, limit := range lim.limits {
		lim.limiters[key] = append(lim.limiters[key], rate.NewLimiter(
			rate.Limit(float64(limit.Amount)/float64(limit.Period)*float64(time.Second)),
			limit.Amount,
		))
	}
}

func (lim *Limiter) check(key string) bool {
	for _, keyLim := range lim.limiters[key] {
		if !keyLim.Allow() {
			return false
		}
	}

	return true
}

func (lim *Limiter) IsAllowed(key string) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	_, exists := lim.limiters[key]

	if !exists {
		lim.add(key)
	}

	return lim.check(key)
}
