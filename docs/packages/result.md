# result

Represents a value coupled with possible error.

## Ok / Error
```go
r := result.Ok(5)
if r.IsSuccess() { fmt.Println(r.Unwrap()) }
```
## Expect
```go
v := r.Expect() // panics if error present
```
