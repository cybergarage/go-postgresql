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

// Count is an aggregator that calculates the sum of values.
type Count struct {
	*Aggr
}

// CountOption is a function that configures the Count aggregator.
type CountOption = AggrOption

// NewCount creates a new Count aggregator with the given options.
func NewCount(opts ...CountOption) (*Count, error) {
	aggr := &Count{
		Aggr: NewAggr(),
	}

	opts = append(opts,
		WithAggrName("COUNT"),
		WithAggrResetFunc(
			func(aggr *Aggr) (float64, error) {
				return 0, nil
			},
		),
		WithAggrAggreateFunc(
			func(aggr *Aggr, accumulatedValue float64, inputValue float64) (float64, error) {
				return 0, nil
			},
		),
		WithAggrFinalizeFunc(
			func(aggr *Aggr, accumulatedValue float64, accumulatedCount int) (float64, error) {
				return float64(accumulatedCount), nil
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

// WithCountArguments sets the arguments for the Count aggregator.
func WithCountArguments(args ...string) CountOption {
	return WithAggrArguments(args...)
}

// WithCountGroupBy sets the group by column for the Count aggregator.
func WithCountGroupBy(group string) CountOption {
	return WithAggrGroupBy(group)
}
