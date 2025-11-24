# Getting Started

## Install

```bash
go get github.com/avila-r/ego@latest
```

## Import

```go
import (
    "github.com/avila-r/ego/box"
    "github.com/avila-r/ego/optional"
    "github.com/avila-r/ego/result"
    "github.com/avila-r/ego/promise"
)
```

## Quick Tour

### Box

```go
b := box.Empty[int]().With(42).Then(box.Square()).Peek(func(v int){ fmt.Println(v) })
fmt.Println(b.GetOrDefault(0)) // 1764
```

### Optional

```go
o := optional.Of(10)
val := o.GetOrDefault(0)
```

### Result

```go
r := result.Ok("ok")
if r.IsSuccess() { fmt.Println(r.Unwrap()) }
```

### Promise

```go
p := promise.Of(func() (int,error){ return 5,nil }).Then(func(v int) int { return v*2 })
fmt.Println(p.Join()) // 10
```

## Next

See `packages/` for deeper docs and `examples/usage.md` for integrated snippets.
