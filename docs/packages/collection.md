# collection

Interfaces for Lists, Sets, Maps plus default slice-backed collection.

## DefaultCollection

```go
c := collection.New[int](1,2,3)
c.Add(4)
c.ForEach(func(v int){ fmt.Println(v) })
```

## Clone

```go
copy := c.Clone()
```
