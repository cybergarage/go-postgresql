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
	"sort"
	"testing"

	"github.com/cybergarage/go-postgresql/postgresql/aggregator"
	"github.com/cybergarage/go-safecast/safecast"

	"github.com/cybergarage/go-logger/log"
)

func TestAggregators(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	tests := []struct {
		orderBy          string
		args             []string
		rows             []aggregator.Row
		expectedSumRows  [][]float64
		expectedAvgRows  [][]float64
		expectedRowCount int
	}{
		{
			orderBy: "",
			args:    []string{"foo"},
			rows: []aggregator.Row{
				{1},
				{2},
				{3},
			},
			expectedSumRows:  [][]float64{{6}},
			expectedAvgRows:  [][]float64{{2}},
			expectedRowCount: 1,
		},
		{
			orderBy: "",
			args:    []string{"foo"},
			rows: []aggregator.Row{
				{1},
				{2},
				{3},
				{4},
			},
			expectedSumRows:  [][]float64{{10}},
			expectedAvgRows:  [][]float64{{2.5}},
			expectedRowCount: 1,
		},
		{
			orderBy: "bar",
			args:    []string{"foo"},
			rows: []aggregator.Row{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
			},
			expectedSumRows:  [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedAvgRows:  [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedRowCount: 4,
		},
		{
			orderBy: "bar",
			args:    []string{"foo"},
			rows: []aggregator.Row{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
			},
			expectedSumRows:  [][]float64{{1, 2}, {2, 4}, {3, 6}, {4, 8}},
			expectedAvgRows:  [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedRowCount: 4,
		},
	}

	for _, test := range tests {

		aggrsFunc := []func() (aggregator.Aggregator, error){
			func() (aggregator.Aggregator, error) {
				return aggregator.NewSum(
					aggregator.WithSumGroupBy(test.orderBy),
					aggregator.WithSumArguments(test.args...),
				)
			},
			func() (aggregator.Aggregator, error) {
				return aggregator.NewAvg(
					aggregator.WithAvgGroupBy(test.orderBy),
					aggregator.WithAvgArguments(test.args...),
				)
			},
		}

		for _, aggrs := range aggrsFunc {

			// Aggregate

			aggr, err := aggrs()
			if err != nil {
				t.Error(err)
				continue
			}

			t.Run(aggr.Name(), func(t *testing.T) {

				for _, row := range test.rows {
					if err := aggr.Aggregate(row); err != nil {
						t.Errorf("Error adding row: %v", err)
						continue
					}
				}

				rs, err := aggr.Finalize()
				if err != nil {
					t.Errorf("Error finalizing Sum: %v", err)
					return
				}

				rsRows := []aggregator.Row{}
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

				switch aggr.(type) {
				case *aggregator.Sum:
					expectedRows = test.expectedSumRows
				case *aggregator.Avg:
					expectedRows = test.expectedAvgRows
				default:
					t.Errorf("Unexpected aggregator type: %T", aggr)
					return
				}

				for n, expectedRow := range expectedRows {
					if len(rsRows[n]) != len(expectedRow) {
						t.Errorf("Expected %d columns, got %d", len(expectedRow), len(rsRows[n]))
						continue
					}
					for i, expectedSum := range expectedRow {
						var rowValue float64
						if err := safecast.ToFloat64(rsRows[n][i], &rowValue); err != nil {
							t.Errorf("Error converting row value to int: %v", err)
							continue
						}
						if rowValue != expectedSum {
							t.Errorf("Expected %f, got %f", expectedSum, rowValue)
							continue
						}
					}
				}
			})
		}
	}
}
