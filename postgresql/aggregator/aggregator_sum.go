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
	"fmt"
)

type Sum struct {
	args []string
	// groupBy is the name of the column to group by.
	groupBy string
	// sum is the sum of the values.
	sum map[any]float64
	// count is the number of values.
	count map[any]int
}

// SumOption is a function that configures the Sum aggregator.
type SumOption func(*Sum) error

// NewSum creates a new Sum aggregator with the given options.
func NewSum(options ...SumOption) (*Sum, error) {
	s := &Sum{
		args:    make([]string, 0),
		groupBy: "",
		sum:     make(map[any]float64),
		count:   make(map[any]int),
	}
	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	// Validate the arguments
	if len(s.args) == 0 {
		return nil, fmt.Errorf("no argument %w", ErrNotSupported)
	}
	if len(s.args) > 1 {
		return nil, fmt.Errorf("multiple argument %w", ErrNotSupported)
	}

	return s, nil
}

// WithSubArguments sets the arguments for the Sum aggregator.
func WithSubArguments(args ...string) SumOption {
	return func(s *Sum) error {
		s.args = args
		return nil
	}
}

// WithSubGroupBy sets the group by column for the Sum aggregator.
func WithSubGroupBy(group string) SumOption {
	return func(s *Sum) error {
		s.groupBy = group
		return nil
	}
}

// Name returns the name of the aggregator.
func (s *Sum) Name() string {
	return "sum"
}

// Reset resets the aggregator to its initial state.
func (s *Sum) Reset() error {
	s.sum = make(map[any]float64)
	s.count = make(map[any]int)
	return nil
}

// Aggregate aggregates a row of data.
func (s *Sum) Aggregate(row Row) error {
	return nil
}

// Finalize finalizes the aggregation and returns the result.
func (s *Sum) Finalize() (ResultSet, error) {
	return NewResultSet(), nil
}
