package pbx

import (
	"encoding/hex"
	"getbox/config"
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
	r := &Request{config, ""}
	res := r.getURL()
	v := "https://localhost/api/route//"

	if res != v {
		t.Errorf("want %v; got %v", v, res)
	}
}
