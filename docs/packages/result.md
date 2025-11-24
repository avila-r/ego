# Result Package Guide

A comprehensive guide to using the `result` package for elegant error handling in Go.

## Table of Contents

1. [What is Result?](#what-is-result)
2. [Creating Results](#creating-results)
3. [Checking Result State](#checking-result-state)
4. [Extracting Values](#extracting-values)
5. [Mutating Results](#mutating-results)
6. [Named Returns Pattern](#named-returns-pattern)
7. [Real-World Examples](#real-world-examples)
8. [Best Practices](#best-practices)

## What is Result?

`Result[T]` is a container type that represents either a successful value of type `T` or an error. It's similar to Rust's `Result` type and helps you handle errors in a more functional way.

```go
type Result[T any]
```

Currently, the API is simple, but we are massively extending this container into the vanguard of Go error handling.

**Why use Result?**

- Makes error handling explicit
- Prevents forgotten error checks
- Enables functional error handling patterns
- Great for named return values

## Creating Results

### Of() - Create from value and error

```go
// When you have both a value and potential error
file, err := os.Open("config.json")

r := result.Of(file, err)
```

### Ok() - Create with just a value

```go
// When you know the operation succeeded
user := User{
    ID: 1, 
    Name: "Alice",
}

r := result.Ok(user)

// Works with any type
_ = result.Ok(42)
_ = result.Ok("hello")
_ = result.Ok(true)

// Even zero values are valid
_ = result.Ok(0)
_ = result.Ok("")
_ = result.Ok(false)
```

### Error() - Create with an error

```go
// When the operation failed
err := failure.New("connection timeout")

r := result.Error[User](err)
```

### Failure() - Alias for Error()

```go
err := failure.New("invalid input")

// Same as Error(), more expressive in some contexts
r := result.Failure[int](err)
```

## Checking Result State

### IsSuccess()

```go
r := result.Ok(42)
if r.IsSuccess() {
    fmt.Println("Operation succeeded!")
    // Safe to extract value
}
```

### IsError()

```go
err := failure.New("failed")

r := result.Error[int](err)
if r.IsError() {
    fmt.Println("Operation failed:", r.Error())
    // Handle error case
}
```

### IsEmpty()

```go
// Empty result has no value and no error
// This happens with forgotten named returns!
var r result.Result[int]
if r.IsEmpty() {
    fmt.Println("Result was never initialized")
}
```

## Extracting Values

### Value() - Get pointer to value

```go
r := result.Ok(42)
if v := r.Value(); v != nil {
    fmt.Println("Value:", *v) // Output: Value: 42
}
```

### Error() - Get error if present

```go
err := failure.New("timeout")

r := result.Error[int](err)
if err := r.Error(); err != nil {
    fmt.Println("Error:", err) // Output: Error: timeout
}
```

### Unwrap() - Get value (panics on error)

```go
// Use when you're SURE it's a success
r := result.Ok(42)

value := r.Unwrap() // value = 42

// DANGER: Panics if result is an error!
r2 := result.Error[int](failure.New("failed"))

value2 := r2.Unwrap() // PANIC!
```

### Take() - Get value safely

```go
r := result.Ok(42)

value, err := r.Take()
if err != nil {
    // Handle error
    return
}

fmt.Println("Value:", *value)
```

### Join() - Unwrap or panic with error message

```go
// Panics with the actual error if result is an error
r := result.Error[int](failure.New("database down"))

value := r.Join() // PANIC: database down
```

### Expect() - Unwrap with custom panic message

```go
r := result.Error[int](failure.New("original error"))

// Panic with custom message
value := r.Expect("Failed to get configuration")
// PANIC: Failed to get configuration

// Panic with original error message if no custom message
value2 := r.Expect()
// PANIC: original error
```

## Mutating Results

### Success() / Ok() - Set success value

```go
var r result.Result[int]

// Set the result to success
r.Success(42) // or r = r.Success(42)

fmt.Println(r.Unwrap()) // 42

// Can change from error to success
r = result.Error[int](failure.New("failed"))

r.Ok(100) // or r = r.Ok(42)

fmt.Println(r.Unwrap()) // 100
```

### Failure() - Set error

```go
var r result.Result[int]

// Set the result to error
r.Failure(failure.New("something went wrong"))

fmt.Println(r.Error()) // something went wrong

// Can change from success to error
result.Ok(42)

r.Failure(failure.New("changed my mind"))

fmt.Println(r.Error()) // changed my mind
```

### Chaining mutations

```go
r := result.Ok(10)

r.Success(20).Success(30).Success(40)

fmt.Println(r.Unwrap()) // 40
```

## Named Returns Pattern

One of the most powerful features of Result is using it with named returns.

### Basic Pattern

```go
func GetUser(id int) (r result.Result[User]) {
    user, err := database.FindUser(id)
    if err != nil {
        return r.Failure(err)
    }
    return r.Success(user)
}
```

### Early Returns

```go
var (
    ErrEmptyInput    = failure.New("input is empty")
    ErrInvalidNumber = failure.New("invalid number")
)

func ValidateAndCreate(input string) (r result.Result[int]) {
    // Validation
    if input == "" {
        return r.Failure(ErrEmptyInput)
    }
    
    // Parse
    value, err := strconv.Atoi(input)
    if err != nil {
        return r.Failure(ErrInvalidNumber)
    }
    
    // Success case
    return r.Success(value)
}
```

### Complex Flow

```go
func ProcessOrder(orderID int) (r result.Result[Order]) {
    // Step 1: Fetch order
    order, err := database.GetOrder(orderID)
    if err != nil {
        return r.Failure(failure.New("fetch failed: %w", err))
    }
    
    // Step 2: Validate
    if order.Total <= 0 {
        return r.Failure(failure.New("invalid order total"))
    }
    
    // Step 3: Process payment
    if err := payment.Charge(order.Total); err != nil {
        return r.Failure(failure.New("payment failed: %w", err))
    }
    
    // Step 4: Update status
    order.Status = "completed"
    if err := database.UpdateOrder(order); err != nil {
        return r.Failure(failure.New("update failed: %w", err))
    }
    
    return r.Success(order)
}
```

### ⚠️ Forgotten Attribution

```go
func DangerousFunction(id int) (r result.Result[int]) {
    if id > 0 {
        return r.Success(id)
    }
    
    // OOPS! Forgot to set Failure or Success
    return // r is now empty (no value, no error)
}

// Usage
r := DangerousFunction(-1)
r.IsEmpty()  // true - dangerous state!
r.IsError()  // true - treated as error
r.Error()    // returns ErrEmptyResult
```

## Real-World Examples

### Database Operations

```go
type User struct {
    ID   int
    Name string
    Email string
}

var (
    ErrEmptyEmail   = failure.New("email cannot be empty")
    ErrUserNotFound = failure.New("user not found")
)

func FindUserByEmail(email string) (r result.Result[User]) {
    // Validate input
    if email == "" {
        return r.Failure(ErrEmptyEmail)
    }
    
    // Query database
    var user User
    err := db.QueryRow(
        "SELECT id, name, email FROM users WHERE email = ?",
        email,
    ).Scan(&user.ID, &user.Name, &user.Email)
    
    if err == sql.ErrNoRows {
        return r.Failure(ErrUserNotFound)
    }

    if err != nil {
        return r.Failure(failure.New("database error: %w", err))
    }
    
    return r.Success(user)
}

// Usage
r := FindUserByEmail("alice@example.com")
if r.IsError() {
    fmt.Printf("Error: %v\n", r.Error())
    return
}

user := r.Unwrap()

fmt.Printf("Found user: %s\n", user.Name)
```

### HTTP API Call

```go
type Json map[string]any

func FetchUserData(id int) (r result.Result[Json]) {
    // Build request
    url := fmt.Sprintf("https://api.example.com/users/%d", id)

    resp, err := http.Get(url)
    if err != nil {
        return r.Failure(failure.Decorate(err, "request failed"))
    }
    defer resp.Body.Close()

    // Check status
    if resp.StatusCode != 200 {
        return r.Failure(failure.New("api returned %d", resp.StatusCode))
    }

    // Parse response
    var data Json
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return r.Failure(failure.Decorate(err, "invalid json"))
    }

    return r.Success(data)
}
```

### File Operations

```go
var (
    ErrInvalidPortNumber = failure.New("invalid port number")
)

func ReadConfig(path string) (r result.Result[Config]) {
    // Read file
    data, err := os.ReadFile(path)
    if err != nil {
        return r.Failure(failure.Decorate(err, "unable to read file"))
    }
    
    // Parse JSON
    var config Config
    if err := json.Unmarshal(data, &config); err != nil {
        return r.Failure(failure.New("invalid config: %w", err))
    }
    
    // Validate
    if config.Port <= 0 {
        return r.Failure(ErrInvalidPortNumber)
    }
    
    return r.Success(config)
}
```

### Pipeline Pattern

```go
func ProcessUserPipeline(id int) (r result.Result[User]) {
    // Each step can fail, but we handle it elegantly
    results := []result.Result[any]{
        ValidateUserID(id),
        FetchUser(id),
        EnrichUserData(id),
    }

    // Check if any step failed
    for i, res := range results {
        if res.IsError() {
            return r.Failure(
                failure.New("step %d failed: %w", i+1, res.Error()),
            )
        }
    }

    // All steps succeeded
    user := results[1].(result.Result[User])

    return r.Success(user.Unwrap())
}
```

### Batch Operations

```go
func ProcessBatch(ids []int) []result.Result[User] {
    results := make([]result.Result[User], len(ids))

    for i, id := range ids {
        results[i] = GetUser(id)
    }

    return results
}

// Usage
results := ProcessBatch([]int{1, 2, 3, 4, 5})

successCount := 0
errorCount := 0

for _, r := range results {
    if r.IsSuccess() {
        successCount++
    } else {
        errorCount++
        log.Printf("Error: %v", r.Error())
    }
}

fmt.Printf("Success: %d, Errors: %d\n", successCount, errorCount)
```

## Best Practices

### ✅ DO: Use named returns with Result

```go
func DoSomething() (r result.Result[int]) {
    if err := validate(); err != nil {
        return r.Failure(err)
    }
    return r.Success(42)
}
```

### ✅ DO: Check state before extracting values

```go
r := DoSomething()
if r.IsError() {
    return
}

value := r.Unwrap()
```

### ✅ DO: Use Expect() with clear messages

```go
// Good: Clear context
config := ReadConfig("app.json").Expect("Failed to load configuration")

// Better: Check before unwrapping in production
r := ReadConfig("app.json")
if r.IsError() {
    log.Fatal("Configuration error:", r.Error())
}

config := r.Unwrap()
```

### ❌ DON'T: Forget to set Success or Failure

```go
// BAD!
func BadFunction() (r result.Result[int]) {
    if condition {
        return r.Success(42)
    }
    return // OOPS! r is empty
}
```

### ❌ DON'T: Use Unwrap() without checking

```go
// BAD!
r := MightFail()

value := r.Unwrap() // Might panic!

// GOOD!
r := MightFail()
if r.IsSuccess() {
    value := r.Unwrap()
}
```

### ❌ DON'T: Ignore the error state

```go
// BAD!
r := DoSomething()

value := *r.Value() // Might be nil!

// GOOD!
r := DoSomething()
if v := r.Value(); v != nil {
    value := *v
}
```

## Quick Reference

| Method | Returns | Panics? | Use Case |
|--------|---------|---------|----------|
| `Of(v, e)` | `Result[T]` | No | Create from value + error |
| `Ok(v)` | `Result[T]` | No | Create success result |
| `Error(e)` | `Result[T]` | No | Create error result |
| `Failure(e)` | `Result[T]` | No | Alias for Error() |
| `Value()` | `*T` | No | Get pointer to value |
| `Error()` | `error` | No | Get error if present |
| `IsSuccess()` | `bool` | No | Check if successful |
| `IsError()` | `bool` | No | Check if error |
| `IsEmpty()` | `bool` | No | Check if uninitialized |
| `Unwrap()` | `T` | Yes | Get value (unsafe) |
| `Take()` | `(*T, error)` | No | Get value safely |
| `Join()` | `T` | Yes | Unwrap or panic with error |
| `Expect(msg)` | `T` | Yes | Unwrap or panic with custom message |
| `Success(v)` | `Result[T]` | No | Set success value |
| `Ok(v)` | `Result[T]` | No | Alias for Success() |
| `Failure(e)` | `Result[T]` | No | Set error |

---

**Remember:** Result is about making error handling explicit and preventing mistakes. Use it when you want clear, functional error handling patterns!
