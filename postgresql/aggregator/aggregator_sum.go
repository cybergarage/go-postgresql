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

	"github.com/cybergarage/go-safecast/safecast"
)

type Sum struct {
	// colums is the list of columns to sum.
	colums []string
	// sums is the sums of the values.
	sums []float64
	// counts is the number of values.
	counts int
	// groupBy is the name of the column to group by.
	groupBy string
	// groupSums is the sum of the values for each group.
	groupSums map[any][]float64
	// groupCounts is the count of the values for each group.
	groupCounts map[any]int
}

// SumOption is a function that configures the Sum aggregator.
type SumOption func(*Sum) error

// NewSum creates a new Sum aggregator with the given options.
func NewSum(options ...SumOption) (*Sum, error) {
	s := &Sum{
		colums:      make([]string, 0),
		groupBy:     "",
		sums:        make([]float64, 0),
		counts:      0,
		groupSums:   make(map[any][]float64),
		groupCounts: make(map[any]int),
	}

	for _, opt := range options {
		if err := opt(s); err != nil {
			return nil, err
		}
	}

	if err := s.Reset(); err != nil {
		return nil, err
	}

	return s, nil
}

// WithSubArguments sets the arguments for the Sum aggregator.
func WithSubArguments(args ...string) SumOption {
	return func(s *Sum) error {
		s.colums = args
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

// GroupBy returns the group by column name and a boolean indicating if it is set.
func (s *Sum) GroupBy() (string, bool) {
	if len(s.groupBy) == 0 {
		return "", false
	}
	return s.groupBy, true
}

// Reset resets the aggregator to its initial state.
func (s *Sum) Reset() error {
	s.colums = []string{}
	if groupBy, ok := s.GroupBy(); ok {
		s.colums = append(s.colums, groupBy)
	}
	for _, arg := range s.colums {
		s.colums = append(s.colums, fmt.Sprintf("%s(%s)", s.Name(), arg))
	}

	// Validate the arguments

	if len(s.colums) == 0 {
		return fmt.Errorf("no argument %w", ErrNotSupported)
	}
	if 1 < len(s.colums) {
		return fmt.Errorf("multiple argument %w", ErrNotSupported)
	}

	// Reset aggregator variables

	s.sums = make([]float64, len(s.colums))
	s.counts = 0

	s.groupSums = make(map[any][]float64)
	s.groupCounts = make(map[any]int)

	return nil
}

// Aggregate aggregates a row of data.
func (s *Sum) Aggregate(row Row) error {
	if len(s.colums) != len(row) {
		return fmt.Errorf("%w column count (%d != %d)", ErrInvalid, len(s.colums), len(row))
	}

	if _, ok := s.GroupBy(); ok {
		group := row[0]
		if _, ok := s.groupSums[group]; !ok {
			s.groupSums[group] = make([]float64, (len(s.colums) - 1))
			s.groupCounts[group] = 0
		}
		for n, rv := range row[1:] {
			var fv float64
			if err := safecast.ToFloat64(rv, &fv); err != nil {
				return fmt.Errorf("[%d] %w row : %s", n, ErrInvalid, err)
			}
			s.groupSums[group][n-1] += fv
		}
		s.groupCounts[group]++
	} else {
		for n, rv := range row {
			var fv float64
			if err := safecast.ToFloat64(rv, &fv); err != nil {
				return fmt.Errorf("[%d] %w row : %s", n, ErrInvalid, err)
			}
			s.sums[n] += fv
		}
		s.counts++
	}

	return nil
}

// Finalize finalizes the aggregation and returns the result.
func (s *Sum) Finalize() (ResultSet, error) {
	rows := make([]Row, 0)
	if _, ok := s.GroupBy(); ok {
		for key, value := range s.groupSums {
			row := make([]any, 0)
			row = append(row, key)
			row = append(row, value)
			rows = append(rows, row)
		}
	} else {
		row := make([]any, 0)
		row = append(row, s.sums)
		rows = append(rows, row)
	}
	return NewResultSet(
		WithRows(rows),
		WithColumns(s.colums),
	), nil
}
