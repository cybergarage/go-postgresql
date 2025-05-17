// Copyright (C) 2019 The go-postgresql Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aggregator

type Avg struct {
	*Aggr
}

// AvgOption is a function that configures the Avg aggregator.
type AvgOption = AggrOption

// NewAvg creates a new Avg aggregator with the given options.
func NewAvg(opts ...AvgOption) (*Avg, error) {
	aggr := &Avg{
		Aggr: NewAggr(),
	}

	opts = append(opts,
		WithAggrName("SUM"),
		WithAggrAggreateFunc(
			func(aggr *Aggr, accumulatedValue float64, inputValue float64) (float64, error) {
				return accumulatedValue + inputValue, nil
			},
		),
		WithAggrFinalizeFunc(
			func(aggr *Aggr, accumulatedValue float64, accumulatedCount int) (float64, error) {
				if accumulatedCount == 0 {
					return 0, nil
				}
				return accumulatedValue / float64(accumulatedCount), nil
			},
		),
	)

	for _, opt := range opts {
		if err := opt(aggr.Aggr); err != nil {
			return nil, err
		}
	}

	if err := aggr.Reset(); err != nil {
		return nil, err
	}

	return aggr, nil
}

// WithAvgArguments sets the arguments for the Avg aggregator.
func WithAvgArguments(args ...string) AvgOption {
	return WithAggrArguments(args...)
}

// WithAvgGroupBy sets the group by column for the Avg aggregator.
func WithAvgGroupBy(group string) AvgOption {
	return WithAggrGroupBy(group)
}
