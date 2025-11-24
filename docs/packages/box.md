# box

Generic mutable container with fluent transforms.

## Core

```go
b := box.Empty[int]().With(4).Then(func(v int) int { return v*10 })
fmt.Println(b.GetOrDefault(0)) // 40
```

## Mapping Another Type

```go
b := box.Of(5)
c := box.Map(b, func(v int) string { return fmt.Sprintf("n=%d", v) })
```

## Filter

```go
keep := b.Filter(func(v int) bool { return v%2==0 })
```
