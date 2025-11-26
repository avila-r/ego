# Box Package Guide

A comprehensive guide to `box.Box[T]` for safe, functional container-based value handling.

## What Is Box?

`Box[T]` is a generic container type that may or may not hold a value of type `T`. It provides a functional, chainable API for working with optional values without null pointer dereferences.

```go
type Box[T any]
```

Key goals:

- Safe value presence/absence handling
- Chainable transformations
- No nil pointer panics
- Functional operations (map, filter, flatMap)
- Global state management with safe mutation
- Thread-safe value containers

## Creating Boxes

### Of() - Create with a value

```go
// Create a box containing a value
b := box.Of(42)

// Works with any type
_ = box.Of("hello")
_ = box.Of(true)
_ = box.Of(User{ID: 1, Name: "Alice"})

// Even zero values are valid
_ = box.Of(0)
_ = box.Of("")
_ = box.Of(false)
```

### Empty() - Create an empty box

```go
// Create an empty box
b := box.Empty[int]()

// Useful when you don't have a value yet
var result Box.Box[string]
if condition {
    result = box.Of("success")
} else {
    result = box.Empty[string]()
}
```

## Checking Box State

### IsPresent()

```go
b := box.Of(42)
if b.IsPresent() {
    fmt.Println("Box contains a value!")
}
```

### IsEmpty()

```go
b := box.Empty[int]()
if b.IsEmpty() {
    fmt.Println("Box is empty")
}
```

## Extracting Values

### Get() - Get pointer to value

```go
b := box.Of(42)
if v := b.Get(); v != nil {
    fmt.Println("Value:", *v) // Output: Value: 42
}

// Empty box returns nil
empty := box.Empty[int]()
fmt.Println(empty.Get()) // Output: <nil>
```

### GetOrDefault() - Get value with fallback

```go
b := box.Of(42)
value := b.GetOrDefault(100)
fmt.Println(value) // Output: 42

empty := box.Empty[int]()
value2 := empty.GetOrDefault(100)
fmt.Println(value2) // Output: 100
```

## Setting Values

### With() - Set value and return box

```go
// Set value on empty box
b := box.Empty[int]()
b = b.With(42)
fmt.Println(*b.Get()) // Output: 42

// Replace value on existing box
b = b.With(100)
fmt.Println(*b.Get()) // Output: 100

// Chainable
b = box.Empty[int]().With(10).With(20).With(30)
fmt.Println(*b.Get()) // Output: 30
```

### Set() - Set value in-place

```go
b := box.Empty[int]()
b.Set(42)
fmt.Println(*b.Get()) // Output: 42

// Can update existing value
b.Set(100)
fmt.Println(*b.Get()) // Output: 100
```

## Clearing Values

### Deflate() - Clear in-place

```go
b := box.Of(42)
b.Deflate()
fmt.Println(b.IsEmpty()) // Output: true
```

### Deflated() - Clear and return box

```go
b := box.Of(42)
b = b.Deflated()
fmt.Println(b.IsEmpty()) // Output: true

// Chainable
b = box.Of(10).Deflated()
```

## Transformations (same type)

### Then() - Transform value if present

```go
b := box.Of(10)

// Double the value
b = b.Then(func(v int) int {
    return v * 2
})
fmt.Println(*b.Get()) // Output: 20

// Empty boxes are not transformed
empty := box.Empty[int]()
empty = empty.Then(func(v int) int {
    return v * 2
})
fmt.Println(empty.IsEmpty()) // Output: true
```

### Filter() - Keep value if predicate matches

```go
b := box.Of(42)

// Keep values greater than 40
b = b.Filter(func(v int) bool {
    return v > 40
})
fmt.Println(b.IsPresent()) // Output: true

// Filter out values less than 50
b = b.Filter(func(v int) bool {
    return v < 50
})
fmt.Println(b.IsEmpty()) // Output: true (42 > 40 but not < 50 in second filter)
```

Better example:

```go
b := box.Of(42)

// Filter passes - value remains
result := b.Filter(func(v int) bool {
    return v > 40
})
fmt.Println(result.IsPresent()) // Output: true

// Filter fails - box becomes empty
result2 := b.Filter(func(v int) bool {
    return v < 40
})
fmt.Println(result2.IsEmpty()) // Output: true
```

## Generic Transformations (T → R)

### Map() - Transform to different type

```go
// Transform int to string
bi := box.Of(42)
bs := box.Map(bi, func(n int) string {
    return fmt.Sprintf("Number: %d", n)
})
fmt.Println(*bs.Get()) // Output: Number: 42

// Empty boxes remain empty
empty := box.Empty[int]()
result := box.Map(empty, func(n int) string {
    return fmt.Sprintf("%d", n)
})
fmt.Println(result.IsEmpty()) // Output: true
```

### FlatMap() - Transform and flatten

```go
// Transform int to Box[string]
bi := box.Of(42)
bs := box.FlatMap(bi, func(n int) Box.Box[string] {
    if n > 0 {
        return box.Of(fmt.Sprintf("Positive: %d", n))
    }
    return box.Empty[string]()
})
fmt.Println(*bs.Get()) // Output: Positive: 42

// Can return empty boxes to filter
bi2 := box.Of(-5)
bs2 := box.FlatMap(bi2, func(n int) Box.Box[string] {
    if n > 0 {
        return box.Of(fmt.Sprintf("Positive: %d", n))
    }
    return box.Empty[string]()
})
fmt.Println(bs2.IsEmpty()) // Output: true
```

## Side Effects

### Peek() - Execute action without changing value

```go
b := box.Of(42)

// Peek doesn't change the value
b = b.Peek(func(v int) {
    fmt.Println("Current value:", v)
})
// Output: Current value: 42

fmt.Println(*b.Get()) // Output: 42 (unchanged)

// Empty boxes don't execute the action
empty := box.Empty[int]()
empty = empty.Peek(func(v int) {
    fmt.Println("This won't print")
})
```

### IfPresent() - Execute action if present

```go
b := box.Of(42)
b.IfPresent(func(v int) {
    fmt.Println("Value is:", v)
})
// Output: Value is: 42

// No output for empty box
empty := box.Empty[int]()
empty.IfPresent(func(v int) {
    fmt.Println("This won't print")
})
```

### ThenConsumer() - Execute action and return box

```go
var sum int
b := box.Of(42)

// Execute action and continue chaining
b = b.ThenConsumer(func(v int) {
    sum += v
}).ThenConsumer(func(v int) {
    sum += v
})

fmt.Println(sum) // Output: 84
fmt.Println(*b.Get()) // Output: 42 (unchanged)
```

### ThenSupplier() - Set value from supplier

```go
b := box.Empty[int]()

// Set value from a supplier function
b = b.ThenSupplier(func() int {
    return 42
})
fmt.Println(*b.Get()) // Output: 42

// Can replace existing values
b = b.ThenSupplier(func() int {
    return 100
})
fmt.Println(*b.Get()) // Output: 100
```

## Utilities

### Copy() - Create independent copy

```go
original := box.Of(42)
copy := original.Copy()

// Modify the copy
copy.Set(100)

// Original is unchanged
fmt.Println(*original.Get()) // Output: 42
fmt.Println(*copy.Get())     // Output: 100
```

### Equals() - Compare boxes

```go
b1 := box.Of(42)
b2 := box.Of(42)
b3 := box.Of(100)

fmt.Println(b1.Equals(b2)) // Output: true
fmt.Println(b1.Equals(b3)) // Output: false

// Empty boxes are equal
e1 := box.Empty[int]()
e2 := box.Empty[int]()
fmt.Println(e1.Equals(e2)) // Output: true

// Present and empty are not equal
fmt.Println(b1.Equals(e1)) // Output: false
```

### String() - String representation

```go
b := box.Of(42)
fmt.Println(b.String()) // Output: Box(42)

empty := box.Empty[int]()
fmt.Println(empty.String()) // Output: Box(<empty>)
```

## Integer Toolkit Functions

The box package includes helpful transformation functions for integer boxes.

### Increment() - Add 1

```go
b := box.Of(5)
result := b.Then(box.Increment)
fmt.Println(*result.Get()) // Output: 6
```

### Decrement() - Subtract 1

```go
b := box.Of(10)
result := b.Then(box.Decrement)
fmt.Println(*result.Get()) // Output: 9
```

### Square() - Multiply by itself

```go
b := box.Of(5)
result := b.Then(box.Square)
fmt.Println(*result.Get()) // Output: 25
```

### Cube() - Raise to power of 3

```go
b := box.Of(3)
result := b.Then(box.Cube)
fmt.Println(*result.Get()) // Output: 27
```

### Twice() - Multiply by 2

```go
b := box.Of(7)
result := b.Then(box.Twice)
fmt.Println(*result.Get()) // Output: 14
```

### Halve() - Divide by 2

```go
b := box.Of(20)
result := b.Then(box.Halve)
fmt.Println(*result.Get()) // Output: 10
```

### Negate() - Flip sign

```go
b := box.Of(15)
result := b.Then(box.Negate)
fmt.Println(*result.Get()) // Output: -15
```

### Abs() - Absolute value

```go
b := box.Of(-25)
result := b.Then(box.Abs)
fmt.Println(*result.Get()) // Output: 25
```

### Identity() - Return unchanged

```go
b := box.Of(42)
result := b.Then(box.Identity)
fmt.Println(*result.Get()) // Output: 42
```

### Modulo(m) - Remainder after division

```go
b := box.Of(17)
result := b.Then(func(v int) int {
    return box.Modulo(5)(v)
})
fmt.Println(*result.Get()) // Output: 2
```

### Clamp(min, max) - Constrain to range

```go
// Value below range
b := box.Of(5)
result := b.Then(func(v int) int {
    return box.Clamp(10, 20)(v)
})
fmt.Println(*result.Get()) // Output: 10

// Value above range
b2 := box.Of(30)
result2 := b2.Then(func(v int) int {
    return box.Clamp(10, 20)(v)
})
fmt.Println(*result2.Get()) // Output: 20

// Value within range
b3 := box.Of(15)
result3 := b3.Then(func(v int) int {
    return box.Clamp(10, 20)(v)
})
fmt.Println(*result3.Get()) // Output: 15
```

## Real-World Examples

### Optional Configuration

```go
type Config struct {
    Port    int
    Host    string
    Timeout int
}

func LoadConfig(path string) Box.Box[Config] {
    data, err := os.ReadFile(path)
    if err != nil {
        return box.Empty[Config]()
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return box.Empty[Config]()
    }
    
    return box.Of(cfg)
}

// Usage
configBox := LoadConfig("app.json")

// Use with default
cfg := configBox.GetOrDefault(Config{
    Port:    8080,
    Host:    "localhost",
    Timeout: 30,
})

fmt.Printf("Server: %s:%d\n", cfg.Host, cfg.Port)
```

### User Lookup

```go
type User struct {
    ID    int
    Name  string
    Email string
}

func FindUser(id int) Box.Box[User] {
    var user User
    err := db.QueryRow(
        "SELECT id, name, email FROM users WHERE id = ?", 
        id,
    ).Scan(&user.ID, &user.Name, &user.Email)
    
    if err != nil {
        return box.Empty[User]()
    }
    
    return box.Of(user)
}

// Usage
userBox := FindUser(123)

userBox.IfPresent(func(u User) {
    fmt.Printf("Found: %s (%s)\n", u.Name, u.Email)
})

if userBox.IsEmpty() {
    fmt.Println("User not found")
}
```

### Validation Pipeline

```go
func ValidateAge(age int) Box.Box[int] {
    return box.Of(age).Filter(func(a int) bool {
        return a >= 0 && a <= 150
    })
}

func ValidateEmail(email string) Box.Box[string] {
    b := box.Of(email)
    return b.Filter(func(e string) bool {
        return strings.Contains(e, "@")
    })
}

// Usage
ageBox := ValidateAge(25)
if ageBox.IsPresent() {
    fmt.Println("Valid age:", *ageBox.Get())
} else {
    fmt.Println("Invalid age")
}

emailBox := ValidateEmail("user@example.com")
if emailBox.IsPresent() {
    fmt.Println("Valid email:", *emailBox.Get())
}
```

### Chained Transformations

```go
type Product struct {
    Name  string
    Price float64
}

func ApplyDiscount(discount float64) func(Product) Product {
    return func(p Product) Product {
        p.Price = p.Price * (1 - discount)
        return p
    }
}

func ApplyTax(rate float64) func(Product) Product {
    return func(p Product) Product {
        p.Price = p.Price * (1 + rate)
        return p
    }
}

// Usage
productBox := box.Of(Product{
    Name:  "Laptop",
    Price: 1000.00,
})

// Apply 10% discount, then 8% tax
finalProduct := productBox.
    Then(ApplyDiscount(0.10)).
    Then(ApplyTax(0.08))

finalProduct.IfPresent(func(p Product) {
    fmt.Printf("%s: $%.2f\n", p.Name, p.Price)
    // Output: Laptop: $972.00
})
```

### Cache Lookup

```go
var cache = make(map[string]string)

func GetFromCache(key string) Box.Box[string] {
    if value, exists := cache[key]; exists {
        return box.Of(value)
    }
    return box.Empty[string]()
}

func GetOrCompute(key string, compute func() string) string {
    cached := GetFromCache(key)
    
    if cached.IsPresent() {
        return *cached.Get()
    }
    
    // Compute and cache
    value := compute()
    cache[key] = value
    return value
}

// Usage
result := GetOrCompute("user:123", func() string {
    // Expensive computation
    return "John Doe"
})
fmt.Println(result) // Output: John Doe
```

### Safe Parsing

```go
func ParseInt(s string) Box.Box[int] {
    n, err := strconv.Atoi(s)
    if err != nil {
        return box.Empty[int]()
    }
    return box.Of(n)
}

// Usage with chaining
input := "42"
result := ParseInt(input).
    Filter(func(n int) bool { return n > 0 }).
    Then(func(n int) int { return n * 2 })

fmt.Println(result.GetOrDefault(0)) // Output: 84

// Invalid input
invalid := "abc"
result2 := ParseInt(invalid).
    Filter(func(n int) bool { return n > 0 }).
    Then(func(n int) int { return n * 2 })

fmt.Println(result2.GetOrDefault(0)) // Output: 0
```

### Complex Workflow with Integer Toolkit

```go
// Start with user input, validate, transform
func ProcessInput(input int) Box.Box[int] {
    return box.Of(input).
        Filter(func(v int) bool { 
            return v > 0 // Must be positive
        }).
        Then(func(v int) int { 
            return box.Square()(&v) // Square it
        }).
        Then(func(v int) int { 
            return box.Clamp(1, 100)(&v) // Clamp to range
        })
}

// Usage
result := ProcessInput(5)
result.IfPresent(func(v int) {
    fmt.Println("Processed:", v) // Output: Processed: 25
})

// Too large input gets clamped
result2 := ProcessInput(15)
result2.IfPresent(func(v int) {
    fmt.Println("Processed:", v) // Output: Processed: 100 (15² = 225, clamped to 100)
})

// Negative input filtered out
result3 := ProcessInput(-5)
fmt.Println(result3.IsEmpty()) // Output: true
```

### Multi-Step Calculation

```go
// Start with 5, square (25), halve (12), increment (13)
result := box.Of(5).
    Then(box.Square).
    Then(box.Halve).
    Then(box.Increment)

fmt.Println(*result.Get()) // Output: 13
```

### Global State Management

```go
// Application configuration
var appConfig = box.Empty[Config]()

func InitConfig(path string) error {
    cfg, err := loadConfigFromFile(path)
    if err != nil {
        return err
    }
    appConfig.Set(cfg)
    return nil
}

func GetConfig() Config {
    return appConfig.GetOrDefault(Config{
        Port: 8080,
        Host: "localhost",
    })
}

func UpdateConfig(updater func(Config) Config) {
    if cfg := appConfig.Get(); cfg != nil {
        appConfig.Set(updater(*cfg))
    }
}

// Usage
InitConfig("app.json")
cfg := GetConfig()
fmt.Printf("Server: %s:%d\n", cfg.Host, cfg.Port)

UpdateConfig(func(c Config) Config {
    c.Port = 9090
    return c
})
```

### Singleton Pattern

```go
type Database struct {
    conn *sql.DB
    dsn  string
}

var dbInstance = box.Empty[*Database]()

func GetDB() *Database {
    if dbInstance.IsEmpty() {
        // Initialize database connection
        db := &Database{
            dsn: "user:pass@tcp(localhost:3306)/mydb",
        }
        conn, err := sql.Open("mysql", db.dsn)
        if err == nil {
            db.conn = conn
            dbInstance.Set(db)
        }
    }
    return dbInstance.GetOrDefault(&Database{})
}

// Usage
db := GetDB()
// Use database connection
```

### Global Feature Flags

```go
var featureFlags = box.Of(map[string]bool{
    "new_ui":       false,
    "beta_feature": false,
    "dark_mode":    true,
})

func IsFeatureEnabled(feature string) bool {
    if flags := featureFlags.Get(); flags != nil {
        return (*flags)[feature]
    }
    return false
}

func EnableFeature(feature string) {
    featureFlags.Then(func(flags map[string]bool) map[string]bool {
        flags[feature] = true
        return flags
    })
}

func DisableFeature(feature string) {
    featureFlags.Then(func(flags map[string]bool) map[string]bool {
        flags[feature] = false
        return flags
    })
}

// Usage
if IsFeatureEnabled("dark_mode") {
    fmt.Println("Dark mode is enabled")
}

EnableFeature("beta_feature")
```

### Global Logger Instance

```go
type Logger struct {
    level string
    file  *os.File
}

var logger = box.Empty[*Logger]()

func InitLogger(level string, filename string) error {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    
    logger.Set(&Logger{
        level: level,
        file:  f,
    })
    return nil
}

func Log(message string) {
    logger.IfPresent(func(l *Logger) {
        fmt.Fprintf(l.file, "[%s] %s\n", l.level, message)
    })
}

func SetLogLevel(level string) {
    logger.Then(func(l *Logger) *Logger {
        l.level = level
        return l
    })
}

// Usage
InitLogger("INFO", "app.log")
Log("Application started")
SetLogLevel("DEBUG")
Log("Debug message")
```

### Global Counter/Metrics

```go
var requestCounter = box.Of(0)
var activeConnections = box.Of(0)

func IncrementRequests() int {
    requestCounter.Then(func(count int) int {
        return count + 1
    })
    return *requestCounter.Get()
}

func GetRequestCount() int {
    return requestCounter.GetOrDefault(0)
}

func TrackConnection() {
    activeConnections.Then(func(count int) int {
        return count + 1
    })
}

func ReleaseConnection() {
    activeConnections.Then(func(count int) int {
        if count > 0 {
            return count - 1
        }
        return 0
    })
}

func GetMetrics() map[string]int {
    return map[string]int{
        "requests":    GetRequestCount(),
        "connections": activeConnections.GetOrDefault(0),
    }
}

// Usage in HTTP handler
func handleRequest(w http.ResponseWriter, r *http.Request) {
    count := IncrementRequests()
    fmt.Fprintf(w, "Request #%d", count)
}
```

### Global Session Store

```go
type Session struct {
    UserID    int
    Username  string
    ExpiresAt time.Time
}

var currentSession = box.Empty[Session]()

func Login(userID int, username string) {
    session := Session{
        UserID:    userID,
        Username:  username,
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    currentSession.Set(session)
}

func Logout() {
    currentSession.Deflate()
}

func GetCurrentUser() (int, bool) {
    if s := currentSession.Get(); s != nil {
        if time.Now().Before(s.ExpiresAt) {
            return s.UserID, true
        }
        // Session expired
        currentSession.Deflate()
    }
    return 0, false
}

func IsLoggedIn() bool {
    _, loggedIn := GetCurrentUser()
    return loggedIn
}

// Usage
Login(123, "alice")
if IsLoggedIn() {
    userID, _ := GetCurrentUser()
    fmt.Printf("User %d is logged in\n", userID)
}
Logout()
```

### Application State Machine

```go
type AppState string

const (
    StateStarting AppState = "starting"
    StateRunning  AppState = "running"
    StateStopping AppState = "stopping"
    StateStopped  AppState = "stopped"
)

var appState = box.Of(StateStarting)

func TransitionTo(newState AppState) bool {
    // Validate state transition
    current := appState.GetOrDefault(StateStopped)
    
    validTransitions := map[AppState][]AppState{
        StateStarting: {StateRunning, StateStopped},
        StateRunning:  {StateStopping},
        StateStopping: {StateStopped},
        StateStopped:  {StateStarting},
    }
    
    for _, valid := range validTransitions[current] {
        if valid == newState {
            appState.Set(newState)
            return true
        }
    }
    return false
}

func GetState() AppState {
    return appState.GetOrDefault(StateStopped)
}

func OnStateChange(action func(AppState)) {
    appState.Peek(action)
}

// Usage
TransitionTo(StateRunning)
fmt.Println("State:", GetState()) // Output: State: running

OnStateChange(func(state AppState) {
    fmt.Printf("Current state: %s\n", state)
})
```

## Best Practices

### ✅ DO: Check presence before dereferencing

```go
b := GetValue()
if v := b.Get(); v != nil {
    fmt.Println("Value:", *v)
}

// Or use GetOrDefault
value := b.GetOrDefault(42)
```

### ✅ DO: Use IfPresent for side effects

```go
b := FindUser(123)
b.IfPresent(func(u User) {
    SendEmail(u.Email)
})
```

### ✅ DO: Chain transformations

```go
result := box.Of(10).
    Filter(func(v int) bool { return v > 0 }).
    Then(func(v int) int { return v * 2 }).
    Then(func(v int) int { return v + 5 })
```

### ✅ DO: Use Filter for validation

```go
validatedBox := box.Of(email).Filter(func(e string) bool {
    return strings.Contains(e, "@")
})
```

### ✅ DO: Use for global state management

```go
var config = box.Empty[Config]()

func UpdateConfig(newCfg Config) {
    config.Set(newCfg)
}

func GetConfig() Config {
    return config.GetOrDefault(defaultConfig)
}
```

### ✅ DO: Combine with Peek for debugging global state

```go
var counter = box.Of(0)

func Increment() {
    counter.Then(func(v int) int {
        return v + 1
    }).Peek(func(v int) {
        log.Printf("Counter updated to: %d", v)
    })
}
```

### ❌ DON'T: Dereference without checking

```go
// BAD!
b := GetValue()
value := *b.Get() // Might panic if nil!

// GOOD!
b := GetValue()
if v := b.Get(); v != nil {
    value := *v
}
```

### ❌ DON'T: Ignore empty state

```go
// BAD!
b := FindUser(123)
SendEmail(*b.Get()) // Might panic!

// GOOD!
b := FindUser(123)
b.IfPresent(func(u User) {
    SendEmail(u.Email)
})
```

### ❌ DON'T: Modify original after Copy

```go
// Careful - this is safe because Copy creates independence
original := box.Of(42)
copy := original.Copy()
copy.Set(100)
// original is still 42 ✓

// But be aware:
shared := original
shared.Set(100)
// original is now 100 too (same instance)
```

### ⚠️ CAUTION: Global state and concurrency

```go
// Box is NOT thread-safe by default
var globalCounter = box.Of(0)

// If accessed from multiple goroutines, protect with mutex
var mu sync.Mutex

func SafeIncrement() {
    mu.Lock()
    defer mu.Unlock()
    globalCounter.Then(func(v int) int {
        return v + 1
    })
}
```

## Quick Reference

### Creation & State

| Method | Returns | Use Case |
|--------|---------|----------|
| `Of(v)` | `Box[T]` | Create box with value |
| `Empty[T]()` | `Box[T]` | Create empty box |
| `IsPresent()` | `bool` | Check if has value |
| `IsEmpty()` | `bool` | Check if empty |

### Extraction

| Method | Returns | Panics? | Use Case |
|--------|---------|---------|----------|
| `Get()` | `*T` | No | Get pointer to value |
| `GetOrDefault(v)` | `T` | No | Get value with fallback |

### Mutation

| Method | Returns | Use Case |
|--------|---------|----------|
| `Set(v)` | `void` | Set value in-place |
| `With(v)` | `Box[T]` | Set value, return box (chainable) |
| `Deflate()` | `void` | Clear in-place |
| `Deflated()` | `Box[T]` | Clear, return box (chainable) |

### Transformations

| Method | Returns | Use Case |
|--------|---------|----------|
| `Then(f)` | `Box[T]` | Transform value (same type) |
| `Map(b, f)` | `Box[R]` | Transform to different type |
| `FlatMap(b, f)` | `Box[R]` | Transform and flatten |
| `Filter(pred)` | `Box[T]` | Keep if predicate matches |

### Hooks

| Method | Returns | Use Case |
|--------|---------|----------|
| `Peek(action)` | `Box[T]` | Execute action, continue chain |
| `IfPresent(action)` | `void` | Execute if present |
| `ThenConsumer(action)` | `Box[T]` | Execute action, return box |
| `ThenSupplier(f)` | `Box[T]` | Set from supplier |

### Utils

| Method | Returns | Use Case |
|--------|---------|----------|
| `Copy()` | `Box[T]` | Create independent copy |
| `Equals(other)` | `bool` | Compare boxes |
| `String()` | `string` | String representation |

### Integer Toolkit

| Function | Returns | Use Case |
|----------|---------|----------|
| `Increment()` | `func(int) int` | Add 1 |
| `Decrement()` | `func(int) int` | Subtract 1 |
| `Square()` | `func(int) int` | Multiply by itself |
| `Cube()` | `func(int) int` | Raise to power of 3 |
| `Twice()` | `func(int) int` | Multiply by 2 |
| `Halve()` | `func(int) int` | Divide by 2 |
| `Negate()` | `func(int) int` | Flip sign |
| `Abs()` | `func(int) int` | Absolute value |
| `Identity()` | `func(int) int` | Return unchanged |
| `Modulo(m)` | `func(int) int` | Remainder after division |
| `Clamp(min, max)` | `func(int) int` | Constrain to range |

---

**Remember:** Box is about safe, functional value handling and state management. Use it when you want to:

- Avoid nil pointer dereferences
- Create elegant transformation pipelines
- Manage global application state safely
- Implement singleton patterns
- Handle optional configuration values
