package urlsigner

import (
	"fmt"
	goalone "github.com/bwmarrin/go-alone"
	"strings"
	"time"
)

type Signer struct {
	Secret []byte
}

func (signer *Signer) GenerateTokenFromString(url string) string {
	var urlToSign string

	crypt := goalone.New(signer.Secret, goalone.Timestamp)
	if strings.Contains(url, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", url)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", url)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))

	return string(tokenBytes)
}

func (signer *Signer) VerifyToken(token string) bool {
	crypt := goalone.New(signer.Secret, goalone.Timestamp)
	_, err := crypt.Unsign([]byte(token))
	if err != nil {
		return false
	}

	return true
}

func (signer *Signer) Expired(token string, minutesExpire int) bool {
	crypt := goalone.New(signer.Secret, goalone.Timestamp)
	ts := crypt.Parse([]byte(token))

	return time.Since(ts.Timestamp) > time.Duration(minutesExpire)*time.Minute
}
