package pbx

import (
	"encoding/hex"
	"fmt"
	"go-getbox/config"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHashPayloadEmpty(t *testing.T) {
	r := &Request{}
	res := r.hashPayload("")
	v := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestHashPayload(t *testing.T) {
	r := &Request{}
	res := r.hashPayload("123")
	v := "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestHmacHash(t *testing.T) {
	r := &Request{}
	res := hex.EncodeToString(r.hmacHash("test", "test123"))
	v := "9e1cdea55a6add8dc6688fbabfd6dd28b1b7896fa39aa36a0bef8f5e6c06c680"
	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestGetUrl(t *testing.T) {
	config := config.Init("../tests")
	config.Set("host", "localhost")
	config.Set("port", 80)
	config.Set("schema", "https")
	r := &Request{}
	r.Config = config
	res := r.getURL()
	v := "https://localhost/api/route//"

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestSign(t *testing.T) {
	config := config.Init("../tests")
	r := &Request{}
	r.Config = config
	r.date = "2016-01-01T12:00:00+00:00"
	res := r.Sign()
	v := "Xg+SwsEGV0KfdJDYRCd773ZJQZFDNQG/7JGLBtJVN8U="

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestAuthHeader(t *testing.T) {
	config := config.Init("../tests")
	r := &Request{}
	r.Config = config
	r.date = "2016-01-01T12:00:00+00:00"
	res := r.AuthHeader()
	v := "Signature ApiKey=\"asd123\", Algorithm=\"hmac-sha256\", SignedHeaders=\"date;x-version\", Signature=\"Xg+SwsEGV0KfdJDYRCd773ZJQZFDNQG/7JGLBtJVN8U=\""

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestApiUrl(t *testing.T) {
	config := config.Init("../tests")
	r := &Request{}
	r.Config = config
	res := r.apiURL()
	v := "/api/route//"

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}

func TestDownloadFile(t *testing.T) {
	r := &Request{}
	err := r.DownloadFile("../tests/downloaded", "http://localhost")
	if err != nil {
		t.Errorf("Couldn't download file %v", err)
	}
}

func TestGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"fake response"}`)
	}))
	defer ts.Close()

	config := config.Init("../tests")
	config.Set("host", extractHost(ts.URL))
	config.Set("Schema", "http")
	r := &Request{}
	r.Config = config

	body, err := r.Get("asd")
	if err != nil {
		t.Errorf("Couldn't make get request to %v", ts.URL)
	}

	v := `{"fake response"}`
	res := string(body)
	if strings.TrimRight(res, "\n") != v {
		t.Errorf("want %v; got %v", v, res)
	}

}

func extractHost(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return u.Host
}
