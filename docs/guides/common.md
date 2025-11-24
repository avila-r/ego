# Common Patterns

## Transform and Chain

```go
value := box.Of(3).
    Then(box.Cube()).
    Then(box.Increment()).
    GetOrDefault(0)
```

## Result Check

```go
r := result.Ok(42)
if r.IsSuccess() { 
    fmt.Println(r.Unwrap()) 
}
```

## Promise Composition

```go
p := promise.Of(func() (int,error) { 
    return 2,nil 
})

q := promise.Compose(p, func(v int) *promise.Promise[int] { 
    return promise.Completed(v * 5)
})

fmt.Println(q.Join()) // 10
```
