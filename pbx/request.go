package pbx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	method        = "GET"
	signedHeaders = "date;x-version"
	algorithm     = "hmac-sha256"
)

type Requester interface {
	Get(endpoint string) ([]byte, error)
	DownloadFile(filepath string, url string) error
	GetConfig() *viper.Viper
}

type Request struct {
	Config   *viper.Viper
	endpoint string
	date     string
}

// NewRequest starts new request instance
func NewRequest(config *viper.Viper) Requester {
	return &Request{config, "", ""}
}

// Get fetches results from remote url as string
func (r *Request) Get(endpoint string) ([]byte, error) {
	r.endpoint = endpoint
	r.date = r.curDate()
	req, err := http.NewRequest(method, r.getURL(), nil)

	req.Header.Set("Authorization", r.AuthHeader())
	req.Header.Set("X-Version", r.Config.GetString("api_version"))
	req.Header.Set("X-Getbox-Id", r.Config.GetString("getbox_id"))
	req.Header.Set("Date", r.date)

	if err != nil {
		log.Printf("Error on request.\n[ERROR] - %v", err)
		return nil, err
	}

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Error on response.\n[ERROR] - %v", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	//log.Println(string([]byte(body)))
	return body, err
}

// AuthHeader creates authentication header for request
func (r *Request) AuthHeader() string {
	var header strings.Builder

	header.WriteString("Signature ApiKey=\"" + r.Config.GetString("auth_access_key") + "\", ")
	header.WriteString("Algorithm=\"" + algorithm + "\", ")
	header.WriteString("SignedHeaders=\"" + signedHeaders + "\", ")
	header.WriteString("Signature=\"" + r.Sign() + "\"")

	return header.String()
}

// Sign creates signature
func (r *Request) Sign() string {
	var sign strings.Builder

	sign.WriteString(method + "\n")
	sign.WriteString(r.apiURL() + "\n")
	sign.WriteString("\n")
	sign.WriteString("date:" + r.date + "\n")
	sign.WriteString("x-version:" + r.Config.GetString("api_version") + "\n")
	sign.WriteString("\n")
	sign.WriteString(signedHeaders + "\n")
	sign.WriteString(r.hashPayload(""))

	return base64.StdEncoding.EncodeToString(
		r.hmacHash(sign.String(), r.Config.GetString("auth_secret_key")))
}

// curDate returns current date in ISO 8601 format
func (r *Request) curDate() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

// creates sha256 hash of given payload
// atm payload is empty
func (r *Request) hashPayload(payload string) string {
	h := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(h[:])
}

// creates hmac has
func (r *Request) hmacHash(sign string, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(sign))
	return h.Sum(nil)
}

// creates request url
func (r *Request) getURL() string {
	var url strings.Builder

	url.WriteString(r.Config.GetString("schema") + "://")
	url.WriteString(r.Config.GetString("host"))
	url.WriteString(r.apiURL())

	return url.String()
}

// creates api url with endpoint
func (r *Request) apiURL() string {
	var url strings.Builder

	url.WriteString("/" + r.Config.GetString("api_url") + "/")
	url.WriteString(r.endpoint + "/")

	return url.String()
}

// GetConfig returns config instance
func (r *Request) GetConfig() *viper.Viper {
	return r.Config
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func (r *Request) DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
