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
	"fmt"
	"sort"
	"testing"

	"github.com/cybergarage/go-postgresql/postgresql/aggr"
	"github.com/cybergarage/go-safecast/safecast"

	"github.com/cybergarage/go-logger/log"
)

func TestAggregators(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	tests := []struct {
		orderBy           string
		args              []string
		rows              []aggr.Row
		expectedSumRows   [][]float64
		expectedAvgRows   [][]float64
		expectedMinRows   [][]float64
		expectedMaxRows   [][]float64
		expectedCountRows [][]float64
		expectedRowCount  int
	}{
		{
			orderBy: "",
			args:    []string{"foo"},
			rows: []aggr.Row{
				{1},
				{2},
				{3},
			},
			expectedSumRows:   [][]float64{{6}},
			expectedAvgRows:   [][]float64{{2}},
			expectedMinRows:   [][]float64{{1}},
			expectedMaxRows:   [][]float64{{3}},
			expectedCountRows: [][]float64{{3}},
			expectedRowCount:  1,
		},
		{
			orderBy: "",
			args:    []string{"foo"},
			rows: []aggr.Row{
				{1},
				{2},
				{3},
				{4},
			},
			expectedSumRows:   [][]float64{{10}},
			expectedAvgRows:   [][]float64{{2.5}},
			expectedMinRows:   [][]float64{{1}},
			expectedMaxRows:   [][]float64{{4}},
			expectedCountRows: [][]float64{{4}},
			expectedRowCount:  1,
		},
		{
			orderBy: "bar",
			args:    []string{"foo"},
			rows: []aggr.Row{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
			},
			expectedSumRows:   [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedAvgRows:   [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedMinRows:   [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedMaxRows:   [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedCountRows: [][]float64{{1, 1}, {2, 1}, {3, 1}, {4, 1}},
			expectedRowCount:  4,
		},
		{
			orderBy: "bar",
			args:    []string{"foo"},
			rows: []aggr.Row{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{1, 2},
				{2, 4},
				{3, 6},
				{4, 8},
			},
			expectedSumRows:   [][]float64{{1, 3}, {2, 6}, {3, 9}, {4, 12}},
			expectedAvgRows:   [][]float64{{1, 1.5}, {2, 3}, {3, 4.5}, {4, 6}},
			expectedMinRows:   [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedMaxRows:   [][]float64{{1, 2}, {2, 4}, {3, 6}, {4, 8}},
			expectedCountRows: [][]float64{{1, 2}, {2, 2}, {3, 2}, {4, 2}},
			expectedRowCount:  4,
		},
		{
			orderBy: "bar",
			args:    []string{"foo"},
			rows: []aggr.Row{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{1, 2},
				{2, 4},
				{3, 6},
				{4, 8},
				{1, 3},
				{2, 6},
				{3, 9},
				{4, 12},
			},
			expectedSumRows:   [][]float64{{1, 6}, {2, 12}, {3, 18}, {4, 24}},
			expectedAvgRows:   [][]float64{{1, 2}, {2, 4}, {3, 6}, {4, 8}},
			expectedMinRows:   [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedMaxRows:   [][]float64{{1, 3}, {2, 6}, {3, 9}, {4, 12}},
			expectedCountRows: [][]float64{{1, 3}, {2, 3}, {3, 3}, {4, 3}},
			expectedRowCount:  4,
		},
	}

	for n, test := range tests {

		t.Run(fmt.Sprintf("%02d", n), func(t *testing.T) {

			aggrFuncs := []func() (aggr.Aggregator, error){
				func() (aggr.Aggregator, error) {
					return aggr.NewSum(
						aggr.WithSumGroupBy(test.orderBy),
						aggr.WithSumArguments(test.args...),
					)
				},
				func() (aggr.Aggregator, error) {
					return aggr.NewAvg(
						aggr.WithAvgGroupBy(test.orderBy),
						aggr.WithAvgArguments(test.args...),
					)
				},
				func() (aggr.Aggregator, error) {
					return aggr.NewMin(
						aggr.WithMinGroupBy(test.orderBy),
						aggr.WithMinArguments(test.args...),
					)
				},
				func() (aggr.Aggregator, error) {
					return aggr.NewMax(
						aggr.WithMaxGroupBy(test.orderBy),
						aggr.WithMaxArguments(test.args...),
					)
				},
				func() (aggr.Aggregator, error) {
					return aggr.NewCount(
						aggr.WithCountGroupBy(test.orderBy),
						aggr.WithCountArguments(test.args...),
					)
				},
			}

			for _, aggrFunc := range aggrFuncs {

				// Aggregate

				testAggr, err := aggrFunc()
				if err != nil {
					t.Error(err)
					continue
				}

				t.Run(testAggr.Name(), func(t *testing.T) {

					for _, row := range test.rows {
						if err := testAggr.Aggregate(row); err != nil {
							t.Errorf("Error adding row: %v", err)
							continue
						}
					}

					rs, err := testAggr.Finalize()
					if err != nil {
						t.Errorf("Error finalizing Sum: %v", err)
						return
					}

					rsRows := []aggr.Row{}
					for rs.Next() {
						row, err := rs.Row()
						if err != nil {
							t.Errorf("Error getting row: %v", err)
							continue
						}
						rsRows = append(rsRows, row)
					}

					sort.Slice(rsRows, func(i, j int) bool {
						var ii, ij int
						if err := safecast.ToInt(rsRows[i][0], &ii); err != nil {
							return false
						}
						if err := safecast.ToInt(rsRows[j][0], &ij); err != nil {
							return false
						}
						return ii < ij
					})

					if len(rsRows) != test.expectedRowCount {
						t.Errorf("Expected %d rows, got %d", test.expectedRowCount, len(rsRows))
						return
					}

					// Compare the result set with the expected rows

					var expectedRows [][]float64

					switch testAggr.(type) {
					case *aggr.Sum:
						expectedRows = test.expectedSumRows
					case *aggr.Avg:
						expectedRows = test.expectedAvgRows
					case *aggr.Min:
						expectedRows = test.expectedMinRows
					case *aggr.Max:
						expectedRows = test.expectedMaxRows
					case *aggr.Count:
						expectedRows = test.expectedCountRows
					default:
						t.Errorf("Unexpected aggr type: %T", testAggr)
						return
					}

					for n, expectedRow := range expectedRows {
						rsRow := rsRows[n]
						if len(rsRow) != len(expectedRow) {
							t.Errorf("%s(%v): Expected %d columns, got %d", testAggr.Name(), test.rows, len(expectedRow), len(rsRow))
							continue
						}
						for i, expectedValue := range expectedRow {
							var rsRowValue float64
							if err := safecast.ToFloat64(rsRow[i], &rsRowValue); err != nil {
								t.Errorf("Error converting row value to float64: %v", err)
								continue
							}
							if rsRowValue != expectedValue {
								t.Errorf("%s(%v): Expected %v, got %v", testAggr.Name(), test.rows, expectedRow, rsRow)
								continue
							}
						}
					}
				})
			}
		})
	}
}
