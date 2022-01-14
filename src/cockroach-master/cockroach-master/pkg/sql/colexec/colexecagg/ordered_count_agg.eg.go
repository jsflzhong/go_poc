// Code generated by execgen; DO NOT EDIT.
// Copyright 2018 The Cockroach Authors.
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package colexecagg

import (
	"unsafe"

	"github.com/cockroachdb/cockroach/pkg/col/coldata"
	"github.com/cockroachdb/cockroach/pkg/sql/colmem"
)

func newCountRowsOrderedAggAlloc(
	allocator *colmem.Allocator, allocSize int64,
) aggregateFuncAlloc {
	return &countRowsOrderedAggAlloc{aggAllocBase: aggAllocBase{
		allocator: allocator,
		allocSize: allocSize,
	}}
}

// countRowsOrderedAgg supports either COUNT(*) or COUNT(col) aggregate.
type countRowsOrderedAgg struct {
	orderedAggregateFuncBase
	col    []int64
	curAgg int64
}

var _ AggregateFunc = &countRowsOrderedAgg{}

func (a *countRowsOrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Int64()
}

func (a *countRowsOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, startIdx, endIdx int, sel []int,
) {
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		if sel == nil {
			_, _ = groups[endIdx-1], groups[startIdx]
			{
				for i := startIdx; i < endIdx; i++ {
					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							a.col[a.curIdx] = a.curAgg
							a.curIdx++
							a.curAgg = int64(0)
						}
						a.isFirstGroup = false
					}

					var y int64
					y = int64(1)
					a.curAgg += y
				}
			}
		} else {
			{
				for _, i := range sel[startIdx:endIdx] {
					if groups[i] {
						if !a.isFirstGroup {
							a.col[a.curIdx] = a.curAgg
							a.curIdx++
							a.curAgg = int64(0)
						}
						a.isFirstGroup = false
					}

					var y int64
					y = int64(1)
					a.curAgg += y
				}
			}
		}
	},
	)
}

func (a *countRowsOrderedAgg) Flush(outputIdx int) {
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	a.col[outputIdx] = a.curAgg
}

func (a *countRowsOrderedAgg) HandleEmptyInputScalar() {
	// COUNT aggregates are special because they return zero in case of an
	// empty input in the scalar context.
	a.col[0] = 0
}

func (a *countRowsOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = 0
}

type countRowsOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []countRowsOrderedAgg
}

var _ aggregateFuncAlloc = &countRowsOrderedAggAlloc{}

const sizeOfCountRowsOrderedAgg = int64(unsafe.Sizeof(countRowsOrderedAgg{}))
const countRowsOrderedAggSliceOverhead = int64(unsafe.Sizeof([]countRowsOrderedAgg{}))

func (a *countRowsOrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(countRowsOrderedAggSliceOverhead + sizeOfCountRowsOrderedAgg*a.allocSize)
		a.aggFuncs = make([]countRowsOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}

func newCountOrderedAggAlloc(
	allocator *colmem.Allocator, allocSize int64,
) aggregateFuncAlloc {
	return &countOrderedAggAlloc{aggAllocBase: aggAllocBase{
		allocator: allocator,
		allocSize: allocSize,
	}}
}

// countOrderedAgg supports either COUNT(*) or COUNT(col) aggregate.
type countOrderedAgg struct {
	orderedAggregateFuncBase
	col    []int64
	curAgg int64
}

var _ AggregateFunc = &countOrderedAgg{}

func (a *countOrderedAgg) SetOutput(vec coldata.Vec) {
	a.orderedAggregateFuncBase.SetOutput(vec)
	a.col = vec.Int64()
}

func (a *countOrderedAgg) Compute(
	vecs []coldata.Vec, inputIdxs []uint32, startIdx, endIdx int, sel []int,
) {
	// If this is a COUNT(col) aggregator and there are nulls in this batch,
	// we must check each value for nullity. Note that it is only legal to do a
	// COUNT aggregate on a single column.
	nulls := vecs[inputIdxs[0]].Nulls()
	a.allocator.PerformOperation([]coldata.Vec{a.vec}, func() {
		// Capture groups to force bounds check to work. See
		// https://github.com/golang/go/issues/39756
		groups := a.groups
		if sel == nil {
			_, _ = groups[endIdx-1], groups[startIdx]
			if nulls.MaybeHasNulls() {
				for i := startIdx; i < endIdx; i++ {
					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							a.col[a.curIdx] = a.curAgg
							a.curIdx++
							a.curAgg = int64(0)
						}
						a.isFirstGroup = false
					}

					var y int64
					y = int64(0)
					if !nulls.NullAt(i) {
						y = 1
					}
					a.curAgg += y
				}
			} else {
				for i := startIdx; i < endIdx; i++ {
					//gcassert:bce
					if groups[i] {
						if !a.isFirstGroup {
							a.col[a.curIdx] = a.curAgg
							a.curIdx++
							a.curAgg = int64(0)
						}
						a.isFirstGroup = false
					}

					var y int64
					y = int64(1)
					a.curAgg += y
				}
			}
		} else {
			if nulls.MaybeHasNulls() {
				for _, i := range sel[startIdx:endIdx] {
					if groups[i] {
						if !a.isFirstGroup {
							a.col[a.curIdx] = a.curAgg
							a.curIdx++
							a.curAgg = int64(0)
						}
						a.isFirstGroup = false
					}

					var y int64
					y = int64(0)
					if !nulls.NullAt(i) {
						y = 1
					}
					a.curAgg += y
				}
			} else {
				for _, i := range sel[startIdx:endIdx] {
					if groups[i] {
						if !a.isFirstGroup {
							a.col[a.curIdx] = a.curAgg
							a.curIdx++
							a.curAgg = int64(0)
						}
						a.isFirstGroup = false
					}

					var y int64
					y = int64(1)
					a.curAgg += y
				}
			}
		}
	},
	)
}

func (a *countOrderedAgg) Flush(outputIdx int) {
	// Go around "argument overwritten before first use" linter error.
	_ = outputIdx
	outputIdx = a.curIdx
	a.curIdx++
	a.col[outputIdx] = a.curAgg
}

func (a *countOrderedAgg) HandleEmptyInputScalar() {
	// COUNT aggregates are special because they return zero in case of an
	// empty input in the scalar context.
	a.col[0] = 0
}

func (a *countOrderedAgg) Reset() {
	a.orderedAggregateFuncBase.Reset()
	a.curAgg = 0
}

type countOrderedAggAlloc struct {
	aggAllocBase
	aggFuncs []countOrderedAgg
}

var _ aggregateFuncAlloc = &countOrderedAggAlloc{}

const sizeOfCountOrderedAgg = int64(unsafe.Sizeof(countOrderedAgg{}))
const countOrderedAggSliceOverhead = int64(unsafe.Sizeof([]countOrderedAgg{}))

func (a *countOrderedAggAlloc) newAggFunc() AggregateFunc {
	if len(a.aggFuncs) == 0 {
		a.allocator.AdjustMemoryUsage(countOrderedAggSliceOverhead + sizeOfCountOrderedAgg*a.allocSize)
		a.aggFuncs = make([]countOrderedAgg, a.allocSize)
	}
	f := &a.aggFuncs[0]
	f.allocator = a.allocator
	a.aggFuncs = a.aggFuncs[1:]
	return f
}
