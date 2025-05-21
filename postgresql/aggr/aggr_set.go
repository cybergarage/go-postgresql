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

type AggregatorSet []Aggregator

// NewAggregatorSetForNames creates a new AggregatorSet for the given names.
func NewAggregatorSetForNames(names []string, opts ...aggrOption) (AggregatorSet, error) {
	aggrSet := make(AggregatorSet, len(names))
	for i, name := range names {
		aggr, err := NewAggregatorForName(name, opts...)
		if err != nil {
			return nil, err
		}
		aggrSet[i] = aggr
	}
	return aggrSet, nil
}

// Names returns the names of the aggregators in the set.
func (aggrSet *AggregatorSet) Names() []string {
	names := make([]string, len(*aggrSet))
	for i, aggr := range *aggrSet {
		names[i] = aggr.Name()
	}
	return names
}

// Reset resets all aggregators in the set.
func (aggrSet *AggregatorSet) Reset() error {
	for _, aggr := range *aggrSet {
		if err := aggr.Reset(); err != nil {
			return err
		}
	}
	return nil
}

// Aggregate aggregates a row of data using all aggregators in the set.
func (aggrSet *AggregatorSet) Aggregate(row Row) error {
	for _, aggr := range *aggrSet {
		if err := aggr.Aggregate(row); err != nil {
			return err
		}
	}
	return nil
}

// Finalize finalizes the aggregation and returns the result set.
func (aggrSet *AggregatorSet) Finalize() ([]ResultSet, error) {
	resultSet := make([]ResultSet, len(*aggrSet))
	for i, aggr := range *aggrSet {
		result, err := aggr.Finalize()
		if err != nil {
			return nil, err
		}
		resultSet[i] = result
	}
	return resultSet, nil
}
