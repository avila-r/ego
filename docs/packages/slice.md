# slice

Functional helpers for slices.

## Map / Filter
```go
values := []int{1,2,3,4}
even := slice.Filter(values, func(v int) bool { return v%2==0 })
squared := slice.Map(even, func(v int) int { return v*v })
```
## Reduce
```go
sum := slice.Reduce(values, 0, func(acc, v int) int { return acc+v })
```
