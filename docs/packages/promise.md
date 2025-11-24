# promise

Async computation primitive with chaining & recovery.

## Supply

```go
p := promise.Of(func() (int,error){ return 3,nil })
fmt.Println(p.Join())
```

## Chain

```go
p2 := p.Then(func(v int) int { return v*2 }).Exceptionally(func(err error) int { return 0 })
```

## Compose

```go
final := promise.Compose(p2, func(v int) *promise.Promise[int]{ return promise.Of(func() (int,error){ return v+1,nil }) })
fmt.Println(final.Join())
```
