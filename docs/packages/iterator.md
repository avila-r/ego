# iterator

Stateful iteration over slice-backed collections.

## Basic

```go
it := iterator.Of(1,2,3)
for it.HasNext() { 
    fmt.Println(it.Next()) 
}
```

## Map

```go
mapped := iterator.Map(it, func(v int) string { 
    return fmt.Sprintf("n=%d", v) 
})
```
