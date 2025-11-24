# optional

Represents presence/absence of a value.

## Basic
```go
o := optional.Of("user")
name := o.GetOrDefault("guest")
```
## Take
```go
val, err := o.Take() // failure.Error if empty
```
