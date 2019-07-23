package pbx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	method        = "GET"
	signedHeaders = "date;x-version"
	algorithm     = "hmac-sha256"
)

type Request struct {
	Config   *viper.Viper
	endpoint string
}

//Get fetches results from remote url as string
func (r *Request) Get(endpoint string) ([]byte, error) {
	r.endpoint = endpoint
	req, err := http.NewRequest(method, r.getURL(), nil)

	req.Header.Set("Authorization", r.AuthHeader())
	req.Header.Set("X-Version", r.Config.GetString("api_version"))
	req.Header.Set("X-Getbox-Id", r.Config.GetString("getbox_id"))
	req.Header.Set("Date", r.curDate())

	if err != nil {
		log.Println("Error on request.\n[ERROR] -", err)
	}

	// Send req using http Client
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	//log.Println(string([]byte(body)))
	return body, err
}

//AuthHeader creates authentication header for request
func (r *Request) AuthHeader() string {
	var header strings.Builder

	header.WriteString("Signature ApiKey=\"" + r.Config.GetString("auth_access_key") + "\", ")
	header.WriteString("Algorithm=\"" + algorithm + "\", ")
	header.WriteString("SignedHeaders=\"" + signedHeaders + "\", ")
	header.WriteString("Signature=\"" + r.Sign() + "\"")

	return header.String()
}

//Sign creates signature
func (r *Request) Sign() string {
	var sign strings.Builder

	sign.WriteString(method + "\n")
	sign.WriteString(r.apiURL() + "\n")
	sign.WriteString("\n")
	d := r.curDate()
	fmt.Println(d)
	sign.WriteString("date:" + d + "\n")
	sign.WriteString("x-version:" + r.Config.GetString("api_version") + "\n")
	sign.WriteString("\n")
	sign.WriteString(signedHeaders + "\n")
	sign.WriteString(r.hashPayload(""))

	fmt.Println(sign.String())

	return base64.StdEncoding.EncodeToString(
		r.hmacHash(sign.String(), r.Config.GetString("auth_secret_key")))
}

func (r *Request) curDate() string {
	t := time.Now()
	return t.Format(time.RFC3339)
}

//creates sha256 hash of given payload
//atm payload is empty
func (r *Request) hashPayload(payload string) string {
	h := sha256.Sum256([]byte(payload))
	return hex.EncodeToString(h[:])
}

//creates hmac has
func (r *Request) hmacHash(sign string, secret string) []byte {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(sign))
	return h.Sum(nil)
}

//creates request url
func (r *Request) getURL() string {
	var url strings.Builder

	url.WriteString(r.Config.GetString("schema") + "://")
	url.WriteString(r.Config.GetString("host"))
	url.WriteString(r.apiURL())

	return url.String()
}

//creates api url with endpoint
func (r *Request) apiURL() string {
	var url strings.Builder

	url.WriteString("/" + r.Config.GetString("api_url") + "/")
	url.WriteString(r.endpoint + "/")

	return url.String()
}
