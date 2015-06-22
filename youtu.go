/*
* File Name:	youtu.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-06-19
 */

package youtu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	USER_ID_MAX_LEN = 110
)

var (
	ErrUserIDTooLong = errors.New("user id too long")
)

type AppSign struct {
	app_id     uint32
	secret_id  string
	secret_key string
	expired    uint32
	user_id    string
}

type Youtu struct {
	app_sign AppSign
	host     string
	app_id   string
}

func Init(appSign AppSign, host string) *Youtu {
	app_id := strconv.Itoa(int(appSign.app_id))
	return &Youtu{
		app_sign: appSign,
		host:     host,
		app_id:   app_id,
	}
}

type DetectMode int

const (
	DetectMode_Normal DetectMode = iota
	DetectMode_BigFace
)

type DetectFaceReq struct {
	App_id string     `json:"app_id"`         //App的 API ID
	Image  string     `json:"image"`          //base64编码的二进制图片数据
	Mode   DetectMode `json:"mode,omitempty"` //检测模式 0/1 正常/大脸模式
}

type Face struct {
	Face_id    string  `json:"face_id"`
	X          int32   `json:"x"`
	Y          int32   `json:"y"`
	Width      float32 `json:"width"`
	Height     float32 `json:"height"`
	Gender     int32   `json:"gender"`
	Age        int32   `json:"age"`
	Expression int32   `json:"expression"`
	Glass      bool    `json:"glass"`
	Pitch      int32   `json:"pitch"`
	Yaw        int32   `json:"yaw"`
	Roll       int32   `json:"roll"`
}

type DetectFaceRsp struct {
	Session_id   string `json:"session_id"`
	Image_id     string `json:"image_id"`
	Image_width  int32  `json:"image_width"`
	Image_height int32  `json:"image_height"`
	Face         []Face `json:"face"`
	ErrorCode    int    `json:"errorcode"`
	ErrorMsg     string `json:"errormsg"`
}

func (y *Youtu) DetectFace(imageData string, mode DetectMode) (dfr DetectFaceRsp, err error) {
	url := "http://" + y.host + "/youtu/api/detectface"
	req := DetectFaceReq{
		App_id: strconv.Itoa(int(y.app_sign.app_id)),
		Image:  imageData,
		Mode:   mode,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	rsp, err := y.get(url, string(data))
	if err != nil {
		return
	}
	err = json.Unmarshal(rsp, &dfr)
	if err != nil {
		return dfr, fmt.Errorf("json.Unmarshal() rsp: %s failed: %s\n", rsp, err)
	}
	return
}

func (y *Youtu) orignalSign() string {
	as := y.app_sign
	now := time.Now().Unix()
	rand.Seed(int64(now))
	rnd := rand.Int31()
	return fmt.Sprintf("a=%d&k=%s&e=%d&t=%d&r=%d&u=%s&f=",
		as.app_id,
		as.secret_id,
		as.expired,
		now,
		rnd,
		as.user_id)
}

func EncodeImage(file string) (imgData string, err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	imgData = base64.StdEncoding.EncodeToString(buf)
	return
}

func (y *Youtu) sign() string {
	orig_sign := y.orignalSign()
	h := hmac.New(sha1.New, []byte(y.app_sign.secret_key))
	h.Write([]byte(orig_sign))
	hm := h.Sum(nil)
	b64 := base64.StdEncoding.EncodeToString([]byte(string(hm) + orig_sign))
	return b64
}

func (y *Youtu) get(addr string, req string) (rsp []byte, err error) {
	tr := &http.Transport{
		DisableCompression: false,
	}
	client := &http.Client{Transport: tr}
	httpreq, err := http.NewRequest("POST", addr, strings.NewReader(req))
	if err != nil {
		return
	}
	httpreq.Header.Add("Authorization", y.sign())
	httpreq.Header.Add("Content-Type", "text/json")
	httpreq.Header.Add("User-Agent", "")
	httpreq.Header.Add("Accept", "*/*")
	httpreq.Header.Add("Expect", "100-continue")
	resp, err := client.Do(httpreq)
	if err != nil {
		return
	}
	rsp, err = ioutil.ReadAll(resp.Body)
	return
}
