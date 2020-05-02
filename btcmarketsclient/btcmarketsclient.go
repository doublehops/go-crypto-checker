package btcmarketsclient

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Config struct {
	Key    string
	Secret string
}

func Configuration(Key, Secret string) *Config {

	c := Config{
		Key:    Key,
		Secret: Secret,
	}

	return &c
}

func (c Config) decodeSecret() []byte {
	decoded, err := base64.StdEncoding.DecodeString(c.Secret)

	if err != nil {
		panic(err)
	}

	return decoded
}

func (c Config) MakeRequest(path string) ([]byte, error) {
	timestamp := time.Now().Unix() * 1000
	timestampStr := strconv.FormatInt(timestamp, 10)

	decodedSecret := c.decodeSecret()
	signature := buildSignature(path, decodedSecret, timestampStr)

	client := &http.Client{}
	request, err := http.NewRequest("GET", "https://api.btcmarkets.net"+path, nil)
	if err != nil {
		panic(err)
	}

	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "btc markets Go client")
	request.Header.Set("Accept-Charset", "UTF-8")
	request.Header.Set("Apikey", c.Key)
	request.Header.Set("Signature", signature)
	request.Header.Set("Timestamp", timestampStr)

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("Failed to make request")
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func buildSignature(path string, secret []byte, timestamp string) string {
	stringToSign := fmt.Sprintf("%s\n%s\n", path, timestamp)

	mac := hmac.New(sha512.New, secret)
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	return signature
}
