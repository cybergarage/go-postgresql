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
	"strings"
)

// AggregatorOption is a function that configures the Aggregator.
type AggregatorOption = aggrOption

// Aggregator is an interface that defines the methods for an aggregator.
func WithAggregatorGroupBy(group string) AggregatorOption {
	return withAggrGroupBy(group)
}

// NewAggregator creates a new Aggregator with the given options.
func NewAggregatorForName(name string, opts ...aggrOption) (Aggregator, error) {
	switch strings.ToUpper(name) {
	case "SUM":
		return NewSum(opts...)
	case "AVG":
		return NewAvg(opts...)
	case "MIN":
		return NewMin(opts...)
	case "MAX":
		return NewMax(opts...)
	case "COUNT":
		return NewCount(opts...)
	}
	return nil, fmt.Errorf("%w aggregator: %s", ErrNotSupported, name)
}
