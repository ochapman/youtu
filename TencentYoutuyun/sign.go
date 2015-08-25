/*
* File Name:	sign.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-08-25
 */
package youtu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

func (y *Youtu) orignalSign() string {
	as := y.appSign
	now := time.Now().Unix()
	rand.Seed(int64(now))
	rnd := rand.Int31()
	sign := fmt.Sprintf("a=%d&k=%s&e=%d&t=%d&r=%d&u=%s&f=",
		as.appID,
		as.secretID,
		now+expiredInterval,
		now,
		rnd,
		as.userID)

	if y.debug {
		fmt.Printf("orignal sign: %s\n", sign)
	}
	return sign
}

func (y *Youtu) sign() string {
	origSign := y.orignalSign()
	h := hmac.New(sha1.New, []byte(y.appSign.secretKey))
	h.Write([]byte(origSign))
	hm := h.Sum(nil)
	//attach orig_sign to hm
	dstSign := []byte(string(hm) + origSign)
	b64 := base64.StdEncoding.EncodeToString(dstSign)
	return b64
}
