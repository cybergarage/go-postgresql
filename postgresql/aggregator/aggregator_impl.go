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

// AggrAggregateFunc is a function that performs aggregation on the given values.
type AggrAggregateFunc func(aggr *Aggr, accumulatedValue float64, inputValue float64) (float64, error)

type Aggr struct {
	name        string
	args        []string
	colums      []string
	aggrs       []float64
	counts      int
	groupBy     string
	groupAggrs  map[any][]float64
	groupCounts map[any]int
	aggFunc     AggrAggregateFunc
}

// AggrOption is a function that configures the Aggr aggregator.
type AggrOption func(*Aggr) error

// NewAggr creates a new Aggr aggregator with the given options.
func NewAggr(options ...AggrOption) *Aggr {
	aggr := &Aggr{
		name:        "",
		args:        make([]string, 0),
		colums:      make([]string, 0),
		groupBy:     "",
		aggrs:       make([]float64, 0),
		counts:      0,
		groupAggrs:  make(map[any][]float64),
		groupCounts: make(map[any]int),
		aggFunc:     nil,
	}
	return aggr
}

// WithAggrName sets the name of the Aggr aggregator.
func WithAggrName(name string) AggrOption {
	return func(aggr *Aggr) error {
		aggr.name = name
		return nil
	}
}

// WithAggrArguments sets the arguments for the Aggr aggregator.
func WithAggrArguments(args ...string) AggrOption {
	return func(aggr *Aggr) error {
		if 1 < len(aggr.args) {
			return fmt.Errorf("multiple argument %w : %v", ErrNotSupported, aggr.args)
		}
		aggr.args = args
		return nil
	}
}

// WithAggrAggreateFunc sets the aggregation function for the Aggr aggregator.
func WithAggrAggreateFunc(aggFunc AggrAggregateFunc) AggrOption {
	return func(aggr *Aggr) error {
		aggr.aggFunc = aggFunc
		return nil
	}
}

// WithAggrGroupBy sets the group by column for the Aggr aggregator.
func WithAggrGroupBy(group string) AggrOption {
	return func(aggr *Aggr) error {
		aggr.groupBy = group
		return nil
	}
}

// Name returns the name of the aggregator.
func (aggr *Aggr) Name() string {
	return aggr.name
}

// GroupBy returns the group by column name and a boolean indicating if it is set.
func (aggr *Aggr) GroupBy() (string, bool) {
	if len(aggr.groupBy) == 0 {
		return "", false
	}
	return aggr.groupBy, true
}

// Reset resets the aggregator to its initial state.
func (aggr *Aggr) Reset() error {
	aggr.colums = []string{}
	if groupBy, ok := aggr.GroupBy(); ok {
		aggr.colums = append(aggr.colums, groupBy)
	}
	for _, arg := range aggr.args {
		aggr.colums = append(aggr.colums, fmt.Sprintf("%s(%s)", aggr.Name(), arg))
	}

	// Validate the arguments

	if len(aggr.colums) == 0 {
		return fmt.Errorf("no argument %w", ErrNotSupported)
	}

	// Reset aggregator variables

	aggr.aggrs = make([]float64, len(aggr.colums))
	aggr.counts = 0

	aggr.groupAggrs = make(map[any][]float64)
	aggr.groupCounts = make(map[any]int)

	return nil
}

// Aggregate aggregates a row of data.
func (aggr *Aggr) Aggregate(row Row) error {
	if len(aggr.colums) != len(row) {
		return fmt.Errorf("%w column count (%d != %d)", ErrInvalid, len(aggr.colums), len(row))
	}

	if _, ok := aggr.GroupBy(); ok {
		group := row[0]
		if _, ok := aggr.groupAggrs[group]; !ok {
			aggr.groupAggrs[group] = make([]float64, (len(aggr.colums) - 1))
			aggr.groupCounts[group] = 0
		}
		for n, rv := range row[1:] {
			var fv float64
			if err := safecast.ToFloat64(rv, &fv); err != nil {
				return fmt.Errorf("[%d] %w row : %s", n, ErrInvalid, err)
			}
			nv, err := aggr.aggFunc(aggr, aggr.groupAggrs[group][n], fv)
			if err != nil {
				return err
			}
			aggr.groupAggrs[group][n] = nv
		}
		aggr.groupCounts[group]++
	} else {
		for n, rv := range row {
			var fv float64
			if err := safecast.ToFloat64(rv, &fv); err != nil {
				return fmt.Errorf("[%d] %w row : %s", n, ErrInvalid, err)
			}
			nv, err := aggr.aggFunc(aggr, aggr.aggrs[n], fv)
			if err != nil {
				return err
			}
			aggr.aggrs[n] = nv
		}
		aggr.counts++
	}

	return nil
}

// Finalize finalizes the aggregation and returns the result.
func (aggr *Aggr) Finalize() (ResultSet, error) {
	rows := make([]Row, 0)
	if _, ok := aggr.GroupBy(); ok {
		for group, values := range aggr.groupAggrs {
			row := make([]any, 0)
			row = append(row, group)
			for _, value := range values {
				row = append(row, value)
			}
			rows = append(rows, row)
		}
	} else {
		row := make([]any, 0)
		for _, value := range aggr.aggrs {
			row = append(row, value)
		}
		rows = append(rows, row)
	}
	return NewResultSet(
		WithRows(rows),
		WithColumns(aggr.colums),
	), nil
}
