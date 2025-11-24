# Usage Examples

## Combined Value Flow

```go
r := result.Ok(5)

value := box.Of(r.Unwrap()).
    Then(box.Square()).
    GetOrDefault(0)

fmt.Println(value)
```

## Async + Result

```go
p := promise.Of(func() (int,error) { 
    return 7,nil 
})

res := promise.Map(p, func(v int) int { 
    return v * 3 
})

fmt.Println(res.Join()) // 21
```
