package httpx_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

func Test_Get(t *testing.T) {
	


}

/*
1. SIMPLE GET REQUEST
*/
func TestSimpleGet(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Status:", resp.Status)
	t.Log("Body:", string(body))
}

/*
2. GET WITH QUERY PARAMETERS
*/
func TestGetWithQuery(t *testing.T) {
	u, _ := url.Parse("https://httpbin.org/get")
	q := u.Query()
	q.Set("search", "golang")
	q.Set("page", "2")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("URL:", u.String())
}

/*
3. CUSTOM CLIENT (timeouts)
*/
func TestCustomClient(t *testing.T) {
	client := &http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get("https://httpbin.org/delay/2")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("Success with custom client")
}

/*
4. MANUAL REQUEST + HEADERS
*/
func TestManualRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "https://httpbin.org/headers", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "Go-HttpClient/1.0")
	req.Header.Set("X-Test", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("Sent custom headers!")
}

/*
5. POST JSON
*/
func TestPostJSON(t *testing.T) {
	payload := map[string]any{
		"name": "Renato",
		"age":  25,
	}

	b, _ := json.Marshal(payload)

	resp, err := http.Post("https://httpbin.org/post", "application/json", bytes.NewReader(b))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("JSON POST sent")
}

/*
6. POST FORM DATA
*/
func TestPostForm(t *testing.T) {
	data := url.Values{}
	data.Set("username", "john")
	data.Set("password", "1234")

	resp, err := http.PostForm("https://httpbin.org/post", data)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("POST form sent")
}

/*
7. MULTIPART (SIMULATED)
*/
func TestPostMultipart(t *testing.T) {
	body := &bytes.Buffer{}
	writer := io.MultiWriter(body)
	writer.Write([]byte("fake file data"))

	req, err := http.NewRequest("POST", "https://httpbin.org/post", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary=---ABC")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("Multipart upload simulated")
}

/*
8. PARSE JSON RESPONSE
*/
func TestParseJSONResponse(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/json")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		t.Fatal(err)
	}

	t.Logf("Parsed JSON: %#v\n", data)
}

/*
9. REQUEST WITH CONTEXT TIMEOUT
*/
func TestRequestWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", "https://httpbin.org/delay/2", nil)

	client := &http.Client{}
	_, err := client.Do(req)

	if err != nil {
		t.Log("Request timed out:", err)
		return
	}

	t.Fatal("expected timeout but request succeeded")
}

/*
10. CHECKING STATUS CODE
*/
func TestCheckStatusCode(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/status/404")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Logf("Expected non-200: got %d", resp.StatusCode)
		return
	}

	t.Fatal("expected 404 but got 200")
}

/*
11. STREAM RESPONSE
*/
func TestStreamResponse(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/stream/5")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	buf := make([]byte, 128)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			t.Log("Chunk:", string(buf[:n]))
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

/*
12. SEND COOKIE
*/
func TestSendCookie(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://httpbin.org/cookies", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: "abc123"})

	client := &http.Client{}
	_, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Cookie sent")
}

/*
13. READ COOKIES
*/
func TestReadCookies(t *testing.T) {
	resp, err := http.Get("https://httpbin.org/cookies/set?test=123")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("Cookies:", resp.Cookies())
}

/*
14. CUSTOM TRANSPORT
*/
func TestCustomTransport(t *testing.T) {
	transport := &http.Transport{
		MaxIdleConns:        50,
		MaxConnsPerHost:     10,
		MaxIdleConnsPerHost: 10,
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Get("https://httpbin.org/get")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("Custom transport used")
}

/*
15. POST RAW TEXT BODY
*/
func TestPostRawBody(t *testing.T) {
	body := strings.NewReader("Hello world")

	req, _ := http.NewRequest("POST", "https://httpbin.org/post", body)
	req.Header.Set("Content-Type", "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	t.Log("Raw body POST sent")
}
