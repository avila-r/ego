# Maps Package Guide

A comprehensive guide to `maps` package for generic, type-safe map implementations with insertion-order preservation and functional operations.

## What Are Ego Maps?

The `maps` package provides two generic map implementations plus utility functions:

- `HashMap[K, V]`: Standard hash map (unordered, fast)
- `LinkedHashMap[K, V]`: Insertion-order preserving map
- Utility functions: `Clone`, `Copy`

Both implement the `collection.Map[K, V]` interface and support functional operations like filtering, cloning, and iteration.

```go
type Map[K comparable, V any] interface {
    Get(K) (V, bool)
    Put(K, V)
    Delete(K)
    // ... more methods
}
```

Key benefits:

- Type safety with generics
- Consistent API across implementations
- Insertion-order preservation (LinkedHashMap)
- Functional utilities (Filter, Clone)
- Interop with standard `map[K]V`

## Creating Maps

### HashMap (Unordered)

```go
// Empty map
m := maps.New[string, int]()
m := maps.Empty[string, int]()
m := maps.NewHashMap[string, int]()

// From std map
m := maps.From(maps.Map[string, int]{
    "a": 1,
    "b": 2,
})

m := maps.From(map[string]int{
    "a": 1,
    "b": 2,
})

// From entries
m := maps.Of(
    collection.Entry[string, int]{Key: "a", Value: 1},
    collection.Entry[string, int]{Key: "b", Value: 2},
)

m := maps.Of([]collection.Entry[string, int]{
    {Key: "a", Value: 1},
    {Key: "b", Value: 2},
}...)
```

### LinkedHashMap (Insertion-Order)

```go
// Empty linked map
m := maps.NewLinked[string, int]()
m := maps.EmptyLinked[string, int]()
m := maps.NewLinkedHashMap[string, int]()

// From std map (preserves order)
m := maps.LinkedFrom(maps.Map[string, int]{
    "a": 1,
    "b": 2,
})

m := maps.LinkedFrom(map[string]int{
    "a": 1,
    "b": 2,
})

// From entries (preserves order)
m := maps.LinkedOf([]collection.Entry[string, int]{
    {Key: "a", Value: 1},
    {Key: "b", Value: 2},
}...)

m := maps.LinkedOf([]collection.Entry[string, int]{
    {Key: "a", Value: 1},
    {Key: "b", Value: 2},
}...)
```

## Basic Operations

### Put & Get

```go
m := maps.New[string, int]()

// Put adds or updates
m.Put("port", 8080)
m.Put("timeout", 30)

// Get retrieves value
if val, ok := m.Get("port"); ok {
    fmt.Println(val) // 8080
}

// Get non-existent returns zero value
val, ok := m.Get("missing") // val=0, ok=false
```

### PutIfAbsent

```go
m := maps.New[string, int]()

// Returns true if inserted
ok := m.PutIfAbsent("port", 8080) // true

// Returns false if key exists
ok = m.PutIfAbsent("port", 9000) // false

fmt.Println(m.Get("port"))        // 8080 (unchanged)
```

### Delete

```go
m := maps.From(maps.Map[string, int]{
    "a": 1,
    "b": 2,
})

m.Delete("a")

fmt.Println(m.ContainsKey("a")) // false
```

### Clear

```go
m := maps.Of(collection.Entry[string, int]{
    Key: "a", 
    Value: 1
})

m.Clear()

fmt.Println(m.IsEmpty()) // true

fmt.Println(m.Len())     // 0
```

## Checking State

```go
m := maps.New[string, int]()

// Size
fmt.Println(m.Len())     // 0
fmt.Println(m.IsEmpty()) // true

m.Put("key", 100)
fmt.Println(m.Len())     // 1
fmt.Println(m.IsEmpty()) // false

// Containment
fmt.Println(m.ContainsKey("key"))   // true
fmt.Println(m.ContainsValue(100))   // true
fmt.Println(m.ContainsValue(999))   // false
```

## Extracting Data

### Keys, Values, Entries

```go
m := maps.From(maps.Map[string, int]{
    "a": 1,
    "b": 2,
})

// As slices
keys := m.KeySlice()     // []string{"a", "b"} (unordered)
vals := m.ValueSlice()   // []int{1, 2} (unordered)
entries := m.ToSlice()   // []Entry{{Key:"a",Value:1}, ...}

// As collections
kColl := m.Keys()        // collection.Collection[string]
vColl := m.Values()      // collection.Collection[int]
eColl := m.Entries()     // collection.Collection[Entry[string,int]]

// As standard map (copy)
stdMap := m.Elements()   // map[string]int
```

### Iterator

```go
m := maps.From(maps.Map[string, int]{
    "x": 10,
    "y": 20,
})

it := m.Iterator()
for it.HasNext() {
    entry := it.Next()
    fmt.Printf("%s=%d\n", entry.Key, entry.Value)
}
```

## Functional Operations

### Filter

```go
m := maps.From(maps.Map[string, int]{
    "a": 1,
    "b": 2,
    "c": 3,
})

// Keep only even values
evens := m.Filter(func(k string, v int) bool {
    return v % 2 == 0
})

// evens contains {b:2}
```

### Clone

```go
original := maps.Of(collection.Entry[string, int]{
    Key: "a", 
    Value: 1
})

// Deep copy
copy := original.Clone()

copy.Put("b", 2)

fmt.Println(original.Len()) // 1
fmt.Println(copy.Len())     // 2
```

## Utility Functions

### Clone (Package-Level)

```go
std := map[string]int{
    "a": 1, 
    "b": 2,
}

// Clone standard map
cloned := maps.Clone(std)

cloned["c"] = 3

fmt.Println(len(std))   // 2
fmt.Println(len(cloned))   // 3
```

### Copy

```go
src := map[string]int{"a": 1, "b": 2}
dst := map[string]int{"c": 3}

// Copy all entries from src to dst
maps.Copy(dst, src)

fmt.Println(dst) // map[a:1 b:2 c:3]
```

## LinkedHashMap: Insertion-Order

The key difference with `LinkedHashMap` is **iteration order preservation**.

```go
m := maps.NewLinkedHashMap[string, int]()

m.Put("first", 1)
m.Put("second", 2)
m.Put("third", 3)

// Iteration maintains insertion order
for _, entry := range m.ToSlice() {
    fmt.Println(entry.Key) // first, second, third (in order)
}

// Update doesn't change position
m.Put("second", 20)

// Still: first, second, third
```

### When to Use LinkedHashMap

Use `LinkedHashMap` when:

- Order matters (config files, user preferences)
- Reproducible iteration is required
- Building ordered outputs (JSON, logs)

Use `HashMap` when:

- Pure key-value lookup (order irrelevant)
- Maximum performance needed
- Memory is a concern (linked list overhead)

## Real-World Examples

### Configuration Store

```go
type Config struct {
    settings collection.Map[string, string]
}

func NewConfig() *Config {
    return &Config{
        settings: maps.NewLinkedHashMap[string, string](),
    }
}

func (c *Config) Set(key, value string) {
    c.settings.Put(key, value)
}

func (c *Config) Get(key string) (string, bool) {
    return c.settings.Get(key)
}

func (c *Config) Export() []string {
    lines := []string{}
    for _, e := range c.settings.ToSlice() {
        lines = append(lines, fmt.Sprintf("%s=%s", e.Key, e.Value))
    }
    return lines
}
```

### Caching Layer

```go
type Cache[K comparable, V any] struct {
    data collection.Map[K, V]
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{
        data: maps.New[K, V](),
    }
}

func (c *Cache[K, V]) GetOrCompute(key K, compute func() V) V {
    if val, ok := c.data.Get(key); ok {
        return val
    }
    val := compute()
    c.data.Put(key, val)
    return val
}

func (c *Cache[K, V]) Invalidate(predicate func(K, V) bool) {
    c.data = c.data.Filter(func(k K, v V) bool {
        return !predicate(k, v)
    })
}
```

### Counting / Grouping

```go
words := []string{"apple", "banana", "apple", "cherry", "banana"}

counts := maps.New[string, int]()
for _, word := range words {
    if val, ok := counts.Get(word); ok {
        counts.Put(word, val+1)
    } else {
        counts.Put(word, 1)
    }
}

// Filter words appearing more than once
frequent := counts.Filter(func(word string, count int) bool {
    return count > 1
})
```

### Ordered JSON-like Output

```go
func BuildOrderedJSON(data map[string]any) string {
    ordered := maps.NewLinkedHashMap[string, any]()
    
    // Insert in specific order
    keys := []string{"id", "name", "email", "created_at"}
    for _, k := range keys {
        if v, ok := data[k]; ok {
            ordered.Put(k, v)
        }
    }
    
    var buf strings.Builder
    buf.WriteString("{")
    entries := ordered.ToSlice()
    for i, e := range entries {
        if i > 0 {
            buf.WriteString(", ")
        }
        buf.WriteString(fmt.Sprintf(`"%s": %v`, e.Key, e.Value))
    }
    buf.WriteString("}")
    return buf.String()
}
```

## Comparison: HashMap vs LinkedHashMap

| Feature | HashMap | LinkedHashMap |
|---------|---------|---------------|
| Lookup speed | O(1) average | O(1) average |
| Insertion order | ❌ No | ✅ Yes |
| Memory overhead | Lower | Higher (linked list) |
| Iteration order | Random | Insertion order |
| Use case | Pure lookup | Ordered iteration |

## Best Practices

### ✅ DO: Use type inference

```go
// Good
m := maps.New[string, int]()

// Verbose (but valid)
var m collection.Map[string, int] = maps.NewHashMap[string, int]()
```

### ✅ DO: Check existence before assuming

```go
// Good
if val, ok := m.Get("key"); ok {
    process(val)
}

// Bad
val, _ := m.Get("key")
process(val) // might process zero value
```

### ✅ DO: Use PutIfAbsent for conditional insert

```go
// Good
m.PutIfAbsent("default_port", 8080)

// Verbose
if !m.ContainsKey("default_port") {
    m.Put("default_port", 8080)
}
```

### ✅ DO: Use LinkedHashMap for reproducible outputs

```go
// Good: config serialization
cfg := maps.NewLinkedHashMap[string, string]()
cfg.Put("host", "localhost")
cfg.Put("port", "8080")

// Bad: HashMap iteration order is random
cfg := maps.NewHashMap[string, string]()
// ... iteration order is unpredictable
```

### ❌ DON'T: Mutate during iteration

```go
// BAD! May cause issues
it := m.Iterator()
for it.HasNext() {
    entry := it.Next()
    if entry.Value > 10 {
        m.Delete(entry.Key) // concurrent modification!
    }
}

// GOOD: Collect keys first
toDelete := []string{}
for _, e := range m.ToSlice() {
    if e.Value > 10 {
        toDelete = append(toDelete, e.Key)
    }
}
for _, k := range toDelete {
    m.Delete(k)
}
```

### ❌ DON'T: Rely on HashMap iteration order

```go
// BAD: Assumes order
m := maps.NewHashMap[string, int]()
m.Put("a", 1)
m.Put("b", 2)
// ToSlice() order is undefined!

// GOOD: Use LinkedHashMap if order matters
m := maps.NewLinkedHashMap[string, int]()
```

## Quick Reference

| Operation | HashMap | LinkedHashMap | Notes |
|-----------|---------|---------------|-------|
| `New[K, V]()` | ✅ | via `NewLinked[K, V]()` | Constructor |
| `Of(entries...)` | ✅ | via `OfLinked(...)` | From entries |
| `Put(k, v)` | ✅ | ✅ | Insert/update |
| `Get(k)` | ✅ | ✅ | Returns (V, bool) |
| `PutIfAbsent(k, v)` | ✅ | ✅ | Conditional insert |
| `Delete(k)` | ✅ | ✅ | Remove entry |
| `Clear()` | ✅ | ✅ | Remove all |
| `Len()` | ✅ | ✅ | Size |
| `IsEmpty()` | ✅ | ✅ | Check empty |
| `ContainsKey(k)` | ✅ | ✅ | Key exists |
| `ContainsValue(v)` | ✅ | ✅ | Value exists (slow) |
| `Filter(predicate)` | ✅ | ✅ | New filtered map |
| `Clone()` | ✅ | ✅ | Deep copy |
| `ToSlice()` | ✅ | ✅ (ordered) | Entry slice |
| `KeySlice()` | ✅ | ✅ (ordered) | Keys |
| `ValueSlice()` | ✅ | ✅ (ordered) | Values |
| `Elements()` | ✅ | ✅ | Copy as `map[K]V` |
| `Iterator()` | ✅ | ✅ (ordered) | Entry iterator |

### Package Utilities

| Function | Signature | Description |
|----------|-----------|-------------|
| `Clone` | `Clone[M ~map[K]V, K comparable, V any](M) M` | Clone std map |
| `Copy` | `Copy[L, R ~map[K]V, K comparable, V any](dst L, src R)` | Copy entries |

**Tip:** Use `LinkedHashMap` when order matters; `HashMap` for pure lookup performance.
