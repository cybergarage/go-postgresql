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

// Sum is an aggregator that calculates the sum of values.
type Sum struct {
	*aggrImpl
}

// SumOption is a function that configures the Sum aggregator.
type SumOption = aggrOption

// NewSum creates a new Sum aggregator with the given options.
func NewSum(opts ...SumOption) (*Sum, error) {
	aggr := &Sum{
		aggrImpl: newAggr(),
	}

	opts = append(opts,
		withAggrName("SUM"),
		withAggrResetFunc(
			func(aggr *aggrImpl) (float64, error) {
				return 0, nil
			},
		),
		withAggrAggreateFunc(
			func(aggr *aggrImpl, accumulatedValue float64, inputValue float64) (float64, error) {
				return accumulatedValue + inputValue, nil
			},
		),
		withAggrFinalizeFunc(
			func(aggr *aggrImpl, accumulatedValue float64, accumulatedCount int) (float64, error) {
				return accumulatedValue, nil
			},
		),
	)

	for _, opt := range opts {
		if err := opt(aggr.aggrImpl); err != nil {
			return nil, err
		}
	}

	if err := aggr.Reset(); err != nil {
		return nil, err
	}

	return aggr, nil
}

// WithSumArguments sets the arguments for the Sum aggregator.
func WithSumArguments(args ...string) SumOption {
	return withAggrArguments(args...)
}

// WithSumGroupBy sets the group by column for the Sum aggregator.
func WithSumGroupBy(group string) SumOption {
	return withAggrGroupBy(group)
}
