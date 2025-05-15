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

type Sum struct {
	// Sum is the sum of the values.
	Sum float64
	// Count is the number of values.
	Count int
}

// NewSum creates a new Sum aggregator.
func NewSum() *Sum {
	return &Sum{}
}

// Name returns the name of the aggregator.
func (s *Sum) Name() string {
	return "sum"
}

// Reset resets the aggregator to its initial state.
func (s *Sum) Reset() error {
	s.Sum = 0
	s.Count = 0
	return nil
}

// Aggregate aggregates a row of data.
func (s *Sum) Aggregate(row Row) error {
	if value, ok := row["value"].(float64); ok {
		s.Sum += value
		s.Count++
	}
	return nil
}

// Finalize finalizes the aggregation and returns the result.
func (s *Sum) Finalize() (Result, error) {
	return Result{
		"sum":   s.Sum,
		"count": s.Count,
	}, nil
}
