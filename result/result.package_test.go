package result_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/avila-r/ego/result"
	"github.com/avila-r/failure"
	"github.com/stretchr/testify/assert"
)

func Test_Map_Generic(t *testing.T) {
	t.Run("map int to string", func(t *testing.T) {
		r := result.Ok(42)
		mapped := result.Map(r, func(v int) string {
			return fmt.Sprintf("Value: %d", v)
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, "Value: 42", mapped.Unwrap())
	})

	t.Run("map string to int", func(t *testing.T) {
		r := result.Ok("123")
		mapped := result.Map(r, func(v string) int {
			i, _ := strconv.Atoi(v)
			return i
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, 123, mapped.Unwrap())
	})

	t.Run("map int to user struct", func(t *testing.T) {
		r := result.Ok(1)
		mapped := result.Map(r, func(id int) User {
			return User{ID: id, Name: "Alice"}
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, User{ID: 1, Name: "Alice"}, mapped.Unwrap())
	})

	t.Run("map user to string", func(t *testing.T) {
		r := result.Ok(User{ID: 1, Name: "Alice"})
		mapped := result.Map(r, func(u User) string {
			return u.Name
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, "Alice", mapped.Unwrap())
	})

	t.Run("map error preserves error type", func(t *testing.T) {
		r := result.Error[int](failure.New("original error"))
		mapped := result.Map(r, func(v int) string {
			return "should not be called"
		})

		assert.True(t, mapped.IsError())
		assert.Equal(t, "original error", mapped.Error().Error())
	})

	t.Run("map with complex transformation", func(t *testing.T) {
		r := result.Ok([]int{1, 2, 3, 4, 5})
		mapped := result.Map(r, func(numbers []int) int {
			sum := 0
			for _, n := range numbers {
				sum += n
			}
			return sum
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, 15, mapped.Unwrap())
	})

	t.Run("map bool to string", func(t *testing.T) {
		r := result.Ok(true)
		mapped := result.Map(r, func(b bool) string {
			if b {
				return "yes"
			}
			return "no"
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, "yes", mapped.Unwrap())
	})

	t.Run("map pointer type", func(t *testing.T) {
		user := &User{ID: 1, Name: "Alice"}
		r := result.Ok(user)
		mapped := result.Map(r, func(u *User) int {
			return u.ID
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, 1, mapped.Unwrap())
	})
}

func Test_MapOr_Generic(t *testing.T) {
	t.Run("map or int to string with success", func(t *testing.T) {
		r := result.Ok(42)
		value := result.MapOr(r, "default", func(v int) string {
			return fmt.Sprintf("%d", v)
		})

		assert.Equal(t, "42", value)
	})

	t.Run("map or int to string with error returns default", func(t *testing.T) {
		r := result.Error[int](failure.New("error"))
		value := result.MapOr(r, "default", func(v int) string {
			return fmt.Sprintf("%d", v)
		})

		assert.Equal(t, "default", value)
	})

	t.Run("map or string to int with success", func(t *testing.T) {
		r := result.Ok("123")
		value := result.MapOr(r, 0, func(v string) int {
			i, _ := strconv.Atoi(v)
			return i
		})

		assert.Equal(t, 123, value)
	})

	t.Run("map or user to string with error", func(t *testing.T) {
		r := result.Error[User](failure.New("user not found"))
		value := result.MapOr(r, "Unknown", func(u User) string {
			return u.Name
		})

		assert.Equal(t, "Unknown", value)
	})

	t.Run("map or with different default type", func(t *testing.T) {
		r := result.Ok(5)
		value := result.MapOr(r, []string{}, func(n int) []string {
			result := make([]string, n)
			for i := 0; i < n; i++ {
				result[i] = fmt.Sprintf("item%d", i)
			}
			return result
		})

		assert.Len(t, value, 5)
		assert.Equal(t, "item0", value[0])
	})

	t.Run("map or bool to struct", func(t *testing.T) {
		r := result.Ok(true)
		value := result.MapOr(r, User{ID: 0, Name: "Guest"}, func(b bool) User {
			if b {
				return User{ID: 1, Name: "Admin"}
			}
			return User{ID: 2, Name: "User"}
		})

		assert.Equal(t, User{ID: 1, Name: "Admin"}, value)
	})
}

func Test_MapOrElse_Generic(t *testing.T) {
	t.Run("map or else int to string with success", func(t *testing.T) {
		r := result.Ok(42)
		value := result.MapOrElse(
			r,
			func(err error) string { return "error occurred" },
			func(v int) string { return fmt.Sprintf("value: %d", v) },
		)

		assert.Equal(t, "value: 42", value)
	})

	t.Run("map or else with error uses fallback", func(t *testing.T) {
		r := result.Error[int](failure.New("computation failed"))
		value := result.MapOrElse(
			r,
			func(err error) string { return "error: " + err.Error() },
			func(v int) string { return fmt.Sprintf("value: %d", v) },
		)

		assert.Equal(t, "error: computation failed", value)
	})

	t.Run("map or else user to int", func(t *testing.T) {
		r := result.Ok(User{ID: 100, Name: "Alice"})
		value := result.MapOrElse(
			r,
			func(err error) int { return -1 },
			func(u User) int { return u.ID },
		)

		assert.Equal(t, 100, value)
	})

	t.Run("map or else string to bool with error", func(t *testing.T) {
		r := result.Error[string](failure.New("parsing failed"))
		value := result.MapOrElse(
			r,
			func(err error) bool { return false },
			func(s string) bool { return s == "true" },
		)

		assert.False(t, value)
	})

	t.Run("map or else with complex fallback", func(t *testing.T) {
		r := result.Error[int](failure.New("network timeout"))
		value := result.MapOrElse(
			r,
			func(err error) User {
				return User{ID: -1, Name: "Error: " + err.Error()}
			},
			func(id int) User {
				return User{ID: id, Name: "Success"}
			},
		)

		assert.Equal(t, -1, value.ID)
		assert.Contains(t, value.Name, "network timeout")
	})

	t.Run("map or else slice to count", func(t *testing.T) {
		r := result.Ok([]string{"a", "b", "c"})
		value := result.MapOrElse(
			r,
			func(err error) int { return 0 },
			func(slice []string) int { return len(slice) },
		)

		assert.Equal(t, 3, value)
	})
}

func Test_FlatMap_Generic(t *testing.T) {
	t.Run("flatmap int to string result", func(t *testing.T) {
		r := result.Ok(42)
		mapped := result.FlatMap(r, func(v int) result.Result[string] {
			return result.Ok(fmt.Sprintf("Number: %d", v))
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, "Number: 42", mapped.Unwrap())
	})

	t.Run("flatmap string to int result", func(t *testing.T) {
		r := result.Ok("123")
		mapped := result.FlatMap(r, func(v string) result.Result[int] {
			i, err := strconv.Atoi(v)
			if err != nil {
				return result.Error[int](err)
			}
			return result.Ok(i)
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, 123, mapped.Unwrap())
	})

	t.Run("flatmap with validation that fails", func(t *testing.T) {
		r := result.Ok("abc")
		mapped := result.FlatMap(r, func(v string) result.Result[int] {
			i, err := strconv.Atoi(v)
			if err != nil {
				return result.Error[int](err)
			}
			return result.Ok(i)
		})

		assert.True(t, mapped.IsError())
		assert.Error(t, mapped.Error())
	})

	t.Run("flatmap user to user id", func(t *testing.T) {
		r := result.Ok(User{ID: 100, Name: "Alice"})
		mapped := result.FlatMap(r, func(u User) result.Result[int] {
			if u.ID <= 0 {
				return result.Error[int](failure.New("invalid user id"))
			}
			return result.Ok(u.ID)
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, 100, mapped.Unwrap())
	})

	t.Run("flatmap error propagates", func(t *testing.T) {
		r := result.Error[int](failure.New("initial error"))
		mapped := result.FlatMap(r, func(v int) result.Result[string] {
			return result.Ok("should not be called")
		})

		assert.True(t, mapped.IsError())
		assert.Equal(t, "initial error", mapped.Error().Error())
	})

	t.Run("flatmap chaining different types", func(t *testing.T) {
		r := result.Ok(5)

		// int -> string -> User -> int
		final := result.FlatMap(
			result.FlatMap(
				result.FlatMap(r, func(n int) result.Result[string] {
					return result.Ok(fmt.Sprintf("ID-%d", n))
				}),
				func(s string) result.Result[User] {
					return result.Ok(User{ID: 5, Name: s})
				},
			),
			func(u User) result.Result[int] {
				return result.Ok(u.ID * 10)
			},
		)

		assert.True(t, final.IsSuccess())
		assert.Equal(t, 50, final.Unwrap())
	})

	t.Run("flatmap with complex transformation", func(t *testing.T) {
		r := result.Ok([]int{1, 2, 3})
		mapped := result.FlatMap(r, func(numbers []int) result.Result[string] {
			if len(numbers) == 0 {
				return result.Error[string](failure.New("empty slice"))
			}
			sum := 0
			for _, n := range numbers {
				sum += n
			}
			return result.Ok(fmt.Sprintf("sum: %d", sum))
		})

		assert.True(t, mapped.IsSuccess())
		assert.Equal(t, "sum: 6", mapped.Unwrap())
	})
}

func Test_AndThen_Generic(t *testing.T) {
	t.Run("andthen is alias for flatmap", func(t *testing.T) {
		r := result.Ok(42)

		viaFlatMap := result.FlatMap(r, func(v int) result.Result[string] {
			return result.Ok(fmt.Sprintf("%d", v))
		})

		viaAndThen := result.AndThen(r, func(v int) result.Result[string] {
			return result.Ok(fmt.Sprintf("%d", v))
		})

		assert.Equal(t, viaFlatMap.Unwrap(), viaAndThen.Unwrap())
	})

	t.Run("andthen chains multiple transformations", func(t *testing.T) {
		r := result.Ok(10)

		final := result.AndThen(
			result.AndThen(r, func(n int) result.Result[string] {
				return result.Ok(fmt.Sprintf("Value: %d", n))
			}),
			func(s string) result.Result[bool] {
				return result.Ok(len(s) > 5)
			},
		)

		assert.True(t, final.IsSuccess())
		assert.True(t, final.Unwrap())
	})
}

func Test_Bind_Generic(t *testing.T) {
	t.Run("bind is alias for flatmap", func(t *testing.T) {
		r := result.Ok(42)

		viaBind := result.Bind(r, func(v int) result.Result[string] {
			return result.Ok(fmt.Sprintf("%d", v))
		})

		viaFlatMap := result.FlatMap(r, func(v int) result.Result[string] {
			return result.Ok(fmt.Sprintf("%d", v))
		})

		assert.Equal(t, viaFlatMap.Unwrap(), viaBind.Unwrap())
	})

	t.Run("bind monadic composition", func(t *testing.T) {
		// Simulating: parse string -> validate -> transform
		parseString := func(s string) result.Result[int] {
			i, err := strconv.Atoi(s)
			if err != nil {
				return result.Error[int](err)
			}
			return result.Ok(i)
		}

		validatePositive := func(n int) result.Result[int] {
			if n <= 0 {
				return result.Error[int](failure.New("must be positive"))
			}
			return result.Ok(n)
		}

		r := result.Ok("42")
		final := result.Bind(
			result.Bind(r, parseString),
			validatePositive,
		)

		assert.True(t, final.IsSuccess())
		assert.Equal(t, 42, final.Unwrap())
	})
}

func Test_Generic_Map_RealWorld_Scenarios(t *testing.T) {
	t.Run("user registration flow", func(t *testing.T) {
		type RegistrationForm struct {
			EmailRequest    string
			PasswordRequest string
		}

		type ValidatedForm struct {
			Email    string
			Password string
		}

		type UserAccount struct {
			ID       int
			Email    string
			Password string
		}

		validateForm := func(form RegistrationForm) result.Result[ValidatedForm] {
			if form.EmailRequest == "" {
				return result.Error[ValidatedForm](failure.New("email required"))
			}
			if len(form.PasswordRequest) < 8 {
				return result.Error[ValidatedForm](failure.New("password too short"))
			}
			return result.Ok(ValidatedForm{
				Email:    form.EmailRequest,
				Password: form.PasswordRequest,
			})
		}

		createAccount := func(vf ValidatedForm) result.Result[UserAccount] {
			return result.Ok(UserAccount{
				ID:       1,
				Email:    vf.Email,
				Password: vf.Password,
			})
		}

		r := result.Ok(RegistrationForm{
			EmailRequest:    "user@example.com",
			PasswordRequest: "securepass123",
		})

		account := result.FlatMap(
			result.FlatMap(r, validateForm),
			createAccount,
		)

		assert.True(t, account.IsSuccess())
		assert.Equal(t, 1, account.Unwrap().ID)
		assert.Equal(t, "user@example.com", account.Unwrap().Email)
	})

	t.Run("api response transformation", func(t *testing.T) {
		type APIResponse struct {
			StatusCode int
			Body       string
		}

		type ParsedData struct {
			Message string
		}

		parseResponse := func(resp APIResponse) result.Result[ParsedData] {
			if resp.StatusCode != 200 {
				return result.Error[ParsedData](fmt.Errorf("bad status: %d", resp.StatusCode))
			}
			return result.Ok(ParsedData{Message: resp.Body})
		}

		r := result.Ok(APIResponse{
			StatusCode: 200,
			Body:       "Success",
		})

		parsed := result.FlatMap(r, parseResponse)

		assert.True(t, parsed.IsSuccess())
		assert.Equal(t, "Success", parsed.Unwrap().Message)
	})

	t.Run("database query transformation", func(t *testing.T) {
		type QueryResult struct {
			Rows []map[string]interface{}
		}

		type ProcessedData struct {
			Count int
			IDs   []int
		}

		processQuery := func(qr QueryResult) result.Result[ProcessedData] {
			if len(qr.Rows) == 0 {
				return result.Error[ProcessedData](failure.New("no results"))
			}

			ids := make([]int, len(qr.Rows))
			for i, row := range qr.Rows {
				if id, ok := row["id"].(int); ok {
					ids[i] = id
				}
			}

			return result.Ok(ProcessedData{
				Count: len(qr.Rows),
				IDs:   ids,
			})
		}

		r := result.Ok(QueryResult{
			Rows: []map[string]interface{}{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
		})

		processed := result.FlatMap(r, processQuery)

		assert.True(t, processed.IsSuccess())
		assert.Equal(t, 2, processed.Unwrap().Count)
		assert.Equal(t, []int{1, 2}, processed.Unwrap().IDs)
	})

	t.Run("configuration loading pipeline", func(t *testing.T) {
		type RawConfig struct {
			Port string
			Host string
		}

		type ParsedConfig struct {
			Port int
			Host string
		}

		parseConfig := func(raw RawConfig) result.Result[ParsedConfig] {
			port, err := strconv.Atoi(raw.Port)
			if err != nil {
				return result.Error[ParsedConfig](failure.New("invalid port"))
			}
			if port < 1 || port > 65535 {
				return result.Error[ParsedConfig](failure.New("port out of range"))
			}
			return result.Ok(ParsedConfig{
				Port: port,
				Host: raw.Host,
			})
		}

		r := result.Ok(RawConfig{
			Port: "8080",
			Host: "localhost",
		})

		config := result.FlatMap(r, parseConfig)

		assert.True(t, config.IsSuccess())
		assert.Equal(t, 8080, config.Unwrap().Port)
	})
}
