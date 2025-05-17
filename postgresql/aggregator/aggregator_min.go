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

import (
	"math"
)

// Min is an aggregator that calculates the minimum of values.
type Min struct {
	*Aggr
}

// MinOption is a function that configures the Min aggregator.
type MinOption = AggrOption

// NewMin creates a new Min aggregator with the given options.
func NewMin(opts ...MinOption) (*Min, error) {
	aggr := &Min{
		Aggr: NewAggr(),
	}

	opts = append(opts,
		WithAggrName("MIN"),
		WithAggrResetFunc(
			func(aggr *Aggr) (float64, error) {
				return math.MaxFloat64, nil
			},
		),
		WithAggrAggreateFunc(
			func(aggr *Aggr, accumulatedValue float64, inputValue float64) (float64, error) {
				if inputValue < accumulatedValue {
					return inputValue, nil
				}
				return accumulatedValue, nil
			},
		),
		WithAggrFinalizeFunc(
			func(aggr *Aggr, accumulatedValue float64, accumulatedCount int) (float64, error) {
				return accumulatedValue, nil
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

// WithMinArguments sets the arguments for the Min aggregator.
func WithMinArguments(args ...string) MinOption {
	return WithAggrArguments(args...)
}

// WithMinGroupBy sets the group by column for the Min aggregator.
func WithMinGroupBy(group string) MinOption {
	return WithAggrGroupBy(group)
}
