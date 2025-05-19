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

// Aggregator is an interface for aggregating data.
type Aggregator interface {
	// Name returns the name of the aggregator.
	Name() string
	// Reset resets the aggregator to its initial state.
	Reset() error
	// Aggregate aggregates a row of data.
	Aggregate(row Row) error
	// Finalize finalizes the aggregation and returns the result.
	Finalize() (ResultSet, error)
}
