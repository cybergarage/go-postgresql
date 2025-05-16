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
	"testing"

	"github.com/cybergarage/go-postgresql/postgresql/aggregator"

	"github.com/cybergarage/go-logger/log"
)

func TestSum(t *testing.T) {
	log.SetStdoutDebugEnbled(true)

	tests := []struct {
		orderBy          string
		args             []string
		rows             []aggregator.Row
		expectedSums     [][]float64
		expectedRowCount int
	}{
		// {
		// 	orderBy: "",
		// 	args:    []string{"foo"},
		// 	rows: []aggregator.Row{
		// 		{1},
		// 		{2},
		// 		{3},
		// 	},
		// 	expectedSums:     [][]float64{{6}},
		// 	expectedRowCount: 1,
		// },
		// {
		// 	orderBy: "",
		// 	args:    []string{"foo"},
		// 	rows: []aggregator.Row{
		// 		{1},
		// 		{2},
		// 		{3},
		// 		{4},
		// 	},
		// 	expectedSums:     [][]float64{{10}},
		// 	expectedRowCount: 1,
		// },
		{
			orderBy: "bar",
			args:    []string{"foo"},
			rows: []aggregator.Row{
				{1, 1},
				{2, 2},
				{3, 3},
				{4, 4},
			},
			expectedSums:     [][]float64{{1, 1}, {2, 2}, {3, 3}, {4, 4}},
			expectedRowCount: 4,
		},
	}

	for _, test := range tests {
		aggr, err := aggregator.NewSum(
			aggregator.WithSubGroupBy(test.orderBy),
			aggregator.WithSubArguments(test.args...),
		)
		if err != nil {
			t.Errorf("Error creating Sum: %v", err)
			continue
		}

		for _, row := range test.rows {
			if err := aggr.Aggregate(row); err != nil {
				t.Errorf("Error adding row: %v", err)
				continue
			}
		}

		rs, err := aggr.Finalize()
		if err != nil {
			t.Errorf("Error finalizing Sum: %v", err)
			continue
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

		if len(rsRows) != test.expectedRowCount {
			t.Errorf("Expected %d rows, got %d", test.expectedRowCount, len(rsRows))
			continue
		}

		for i, expectedSum := range test.expectedSums {
			if len(rsRows[i]) != len(expectedSum) {
				t.Errorf("Expected %d columns, got %d", len(expectedSum), len(rsRows[i]))
				continue
			}
			for j, expected := range expectedSum {
				if rsRows[i][j] != expected {
					t.Errorf("Expected %f, got %f", expected, rsRows[i][j])
					continue
				}
			}
		}
	}
}
