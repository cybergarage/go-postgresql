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

type groupSum struct {
	sum   float64
	count int
}
type Sum struct {
	// groupBy is the name of the column to group by.
	groupBy string
	// sum is the sum of the values.
	sum map[any]float64
	// count is the number of values.
	count map[any]int
}

// SumOption is a function that configures the Sum aggregator.
type SumOption func(*Sum)

func NewSum(options ...SumOption) *Sum {
	s := &Sum{
		groupBy: "",
		sum:     make(map[any]float64),
		count:   make(map[any]int),
	}
	for _, opt := range options {
		opt(s)
	}
	return s
}

// WithSubGroupBy sets the group by column for the Sum aggregator.
func WithSubGroupBy(group string) SumOption {
	return func(s *Sum) {
		s.groupBy = group
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
	if s.groupBy == "" {
		return nil
	}

	groupValue, ok := row[s.groupBy]
	if !ok {
		return nil
	}

	value, ok := row["value"]
	if !ok {
		return nil
	}

	s.sum[groupValue] += value.(float64)
	s.count[groupValue]++

	return nil
}

// Finalize finalizes the aggregation and returns the result.
func (s *Sum) Finalize() (Result, error) {
	result := make(Result)
	for groupValue, sum := range s.sum {
		count := s.count[groupValue]
		if count > 0 {
			result[groupValue] = sum / float64(count)
		} else {
			result[groupValue] = 0
		}
	}
	return result, nil
}
