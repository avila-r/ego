# Migration Guide

## From std slices/maps

Replace manual loops with `slice.Map`, `slice.Filter`, `iterator.Of` for traversal.

## From channels for simple async

Use `promise.Of` + `Then` instead of spinning goroutines + channels for single-result flows.

## From manual (value,error) pairs

Use `result.Ok` / `result.Error` and chain state checks; panic on `Join()` only when intentional.
