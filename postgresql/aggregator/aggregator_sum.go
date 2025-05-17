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
	*Aggr
}

// SumOption is a function that configures the Sum aggregator.
type SumOption = AggrOption

// NewSum creates a new Sum aggregator with the given options.
func NewSum(options ...SumOption) (*Sum, error) {
	aggr := &Sum{
		Aggr: NewAggr(),
	}

	for _, opt := range options {
		if err := opt(aggr.Aggr); err != nil {
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
	return WithAggrArguments(args...)
}

// WithSumGroupBy sets the group by column for the Sum aggregator.
func WithSumGroupBy(group string) SumOption {
	return WithAggrGroupBy(group)
}

// Finalize finalizes the aggregation and returns the result.
func (aggr *Sum) Finalize() (ResultSet, error) {
	rows := make([]Row, 0)
	if _, ok := aggr.GroupBy(); ok {
		for group, values := range aggr.groupSums {
			row := make([]any, 0)
			row = append(row, group)
			for _, value := range values {
				row = append(row, value)
			}
			rows = append(rows, row)
		}
	} else {
		row := make([]any, 0)
		for _, value := range aggr.sums {
			row = append(row, value)
		}
		rows = append(rows, row)
	}
	return NewResultSet(
		WithRows(rows),
		WithColumns(aggr.colums),
	), nil
}
