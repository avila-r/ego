# list

Array-backed list implementation.

## Create
```go
l := list.Of(1,2,3)
l.Add(4)
fmt.Println(l.Size())
```
## Iterate
```go
l.Iterator().ForEach(func(v int){ fmt.Println(v) })
```
