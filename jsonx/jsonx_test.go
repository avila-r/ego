package jsonx_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

/*
1. MARSHAL A BASIC STRUCT
*/
func TestMarshalStruct(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	u := User{"Alice", 30}
	b, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(b)) // {"Name":"Alice","Age":30}
}

/*
2. MARSHAL WITH TAGS
*/
func TestMarshalWithTags(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age,omitempty"` // omit if 0
	}

	u := User{Name: "Alice"}
	b, _ := json.Marshal(u)

	t.Log(string(b)) // {"name":"Alice"}
}

/*
3. UNMARSHAL INTO STRUCT
*/
func TestUnmarshalStruct(t *testing.T) {
	jsonStr := `{"name":"Bob","age":40}`

	type User struct {
		Name string
		Age  int
	}

	var u User
	if err := json.Unmarshal([]byte(jsonStr), &u); err != nil {
		t.Fatal(err)
	}

	t.Log(u)
}

/*
4. UNMARSHAL INTO MAP
*/
func TestUnmarshalMap(t *testing.T) {
	jsonStr := `{"score": 99, "active": true}`

	var m map[string]any
	json.Unmarshal([]byte(jsonStr), &m)

	t.Log(m["score"], m["active"])
}

/*
5. MARSHAL INDENT (PRETTY PRINT)
*/
func TestMarshalIndent(t *testing.T) {
	obj := map[string]any{
		"name": "Charlie",
		"age":  28,
	}

	b, _ := json.MarshalIndent(obj, "", "  ")
	t.Log("\n" + string(b))
}

/*
6. RAW MESSAGE (DEFER PARSING)
*/
func TestRawMessage(t *testing.T) {
	type Event struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	jsonStr := `{"type":"user","data":{"id":1,"name":"Dave"}}`

	var e Event
	json.Unmarshal([]byte(jsonStr), &e)

	t.Log("Type:", e.Type)
	t.Log("Raw data:", string(e.Data))

	// Parse Data later
	var user map[string]any
	json.Unmarshal(e.Data, &user)

	t.Log(user)
}

/*
7. STREAM DECODING (DECODER)
*/
func TestDecodeStream(t *testing.T) {
	stream := `{"id":1}{"id":2}{"id":3}`

	type Item struct {
		ID int `json:"id"`
	}

	dec := json.NewDecoder(bytes.NewBufferString(stream))

	for dec.More() {
		var it Item
		if err := dec.Decode(&it); err != nil {
			t.Fatal(err)
		}
		t.Log("Decoded item:", it)
	}
}

/*
8. STREAM ENCODING (ENCODER)
*/
func TestEncodeStream(t *testing.T) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	enc.Encode(map[string]any{"a": 1})
	enc.Encode(map[string]any{"b": 2})
	enc.Encode(map[string]any{"c": 3})

	t.Log(buf.String())
}

type Color struct {
	R, G, B int
}

func (c Color) MarshalJSON() ([]byte, error) {
	hex := fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B)
	return json.Marshal(hex)
}

func (c *Color) UnmarshalJSON(b []byte) error {
	var hex string
	if err := json.Unmarshal(b, &hex); err != nil {
		return err
	}

	_, err := fmt.Sscanf(hex, "#%02x%02x%02x", &c.R, &c.G, &c.B)

	return err
}

/*
9. CUSTOM MARSHALER
*/
func TestCustomMarshalJSON(t *testing.T) {
	c := Color{255, 0, 128}

	b, _ := json.Marshal(c)

	t.Log(string(b)) // "#ff0080"
}

/*
10. CUSTOM UNMARSHALER
*/
func TestCustomUnmarshalJSON(t *testing.T) {
	jsonStr := `" #00ff7f"`
	var c Color
	json.Unmarshal([]byte(jsonStr), &c)

	t.Logf("Parsed color: %+v", c)
}

/*
11. HANDLE UNKNOWN FIELDS SAFELY
*/
func TestUnknownFields(t *testing.T) {
	jsonStr := `{"name":"Eve","extra":123}`

	type User struct {
		Name string
	}

	var u User

	dec := json.NewDecoder(bytes.NewBufferString(jsonStr))
	dec.DisallowUnknownFields()

	err := dec.Decode(&u)
	if err != nil {
		t.Log("Unknown fields error:", err)
		return
	}

	t.Fatal("expected error from unknown field")
}

/*
12. DECODE NUMBER AS json.Number (avoid float64 loss)
*/
func TestUseJSONNumber(t *testing.T) {
	jsonStr := `{"value":1234567890123456789}`

	dec := json.NewDecoder(strings.NewReader(jsonStr))
	dec.UseNumber()

	var obj map[string]json.Number
	dec.Decode(&obj)

	// Keep number as string
	t.Log("As string:", obj["value"].String())
}

/*
13. OMIT EMPTY FIELDS
*/
func TestOmitEmpty(t *testing.T) {
	type User struct {
		Name string `json:"name"`
		Age  int    `json:"age,omitempty"` // omitted if zero
	}

	u := User{Name: "Frank"}
	b, _ := json.Marshal(u)

	t.Log(string(b)) // {"name":"Frank"}
}

/*
14. INLINE EMBEDDED STRUCTS
*/
func TestInlineStructs(t *testing.T) {
	type Info struct {
		Active bool   `json:"active"`
		Note   string `json:"note"`
	}

	type User struct {
		Name string           `json:"name"`
		Info `json:",inline"` // Go does not actually support ',inline' like YAML
	}

	// Note: JSON doesn't support inline tagging.
	// But embedded structs *are* flattened automatically.

	u := User{
		Name: "Grace",
		Info: Info{Active: true, Note: "test"},
	}

	b, _ := json.Marshal(u)
	t.Log(string(b)) // {"name":"Grace","active":true,"note":"test"}
}

/*
15. PARTIAL DECODING (READ ONLY SOME FIELDS)
*/
func TestPartialDecoding(t *testing.T) {
	jsonStr := `{"user":"Henry","age":50,"secret":"ignore me"}`

	var partial struct {
		User string `json:"user"`
	}

	json.Unmarshal([]byte(jsonStr), &partial)
	t.Log(partial.User) // Henry
}
