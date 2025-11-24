# constraint

Reusable generic constraints for numeric, comparable and arithmetic sets.

```go
func Sum[T constraint.Integer](vals []T) T {
    var s T
    for _, v := range vals { s += v }
    return s
}
```
