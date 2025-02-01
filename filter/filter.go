package filter

import (
	"github.com/flint92/webapp/context"
	"log"
	"time"
)

type Filter func(ctx *context.Context)

type FilterBuilder func(next Filter) Filter

var _ FilterBuilder = MetricFilterBuilder

func MetricFilterBuilder(next Filter) Filter {
	return func(ctx *context.Context) {
		start := time.Now().Nanosecond()
		next(ctx)
		end := time.Now().Nanosecond()
		log.Printf("用时: %d纳秒!", end-start)
	}
}
