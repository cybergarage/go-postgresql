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

package aggr

import (
	"math"
)

// Max is an aggregator that calculates the sum of values.
type Max struct {
	*aggrImpl
}

// MaxOption is a function that configures the Max aggregator.
type MaxOption = aggrOption

// NewMax creates a new Max aggregator with the given options.
func NewMax(opts ...MaxOption) (*Max, error) {
	aggr := &Max{
		aggrImpl: newAggr(),
	}

	opts = append(opts,
		withAggrName("MAX"),
		withAggrResetFunc(
			func(aggr *aggrImpl) (float64, error) {
				return math.Inf(-1), nil
			},
		),
		withAggrAggreateFunc(
			func(aggr *aggrImpl, accumulatedValue float64, inputValue float64) (float64, error) {
				if accumulatedValue < inputValue {
					return inputValue, nil
				}
				return accumulatedValue, nil
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

// WithMaxArguments sets the arguments for the Max aggregator.
func WithMaxArguments(args ...string) MaxOption {
	return withAggrArguments(args...)
}

// WithMaxGroupBy sets the group by column for the Max aggregator.
func WithMaxGroupBy(group string) MaxOption {
	return WithAggrGroupBy(group)
}
