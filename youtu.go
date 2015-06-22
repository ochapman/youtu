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
	App_id string `json:"app_id"` //App的 API ID
	Image  string `json:"image"`  //base64编码的二进制图片数据
	//	Mode   int    `omiempty,json:"mode"` //检测模式 0/1 正常/大脸模式
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
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	rsp, err := y.get(url, string(data))
	if err != nil {
		fmt.Println("y.get failed", err)
		return
	}
	err = json.Unmarshal(rsp, &dfr)
	if err != nil {
		fmt.Printf("json.Unmarshal() rsp: %s failed: %s\n", rsp, err)
	}
	return
}

func (y *Youtu) orignalSign() string {
	as := y.app_sign
	now := time.Now().Unix()
	//now := 1434900916
	rand.Seed(int64(now))
	//rnd := rand.Int31n(99999999)
	rnd := 1882422330
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
	//orig_sign := "a=10000089&k=AKIDZgHaklNpo7vrztgYLl1LEyHO3I1GCcvz&e=1436353609&t=1434929048&r=1752593874&u=3041722595&f="
	fmt.Println("orig_sign:", orig_sign)
	h := hmac.New(sha1.New, []byte(y.app_sign.secret_key))
	h.Write([]byte(orig_sign))
	hm := h.Sum(nil)
	for i := range hm {
		fmt.Printf("%2x ", hm[i])
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(string(hm) + orig_sign))
	fmt.Println("b64:", b64)
	return b64
}

func (y *Youtu) get(addr string, req string) (rsp []byte, err error) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	//fmt.Printf("req: %s\n", req)
	httpreq, err := http.NewRequest("POST", addr, strings.NewReader(req))
	if err != nil {
		return
	}
	httpreq.Header.Add("Authorization", y.sign())
	httpreq.Header.Add("Content-Type", "text/json")
	httpreq.Header.Add("User-Agent", "")
	httpreq.Header.Add("Accept", "*/*")
	httpreq.Header.Add("Expect", "100-continue")
	fmt.Printf("httpreq: %s\n", httpreq.URL)
	resp, err := client.Do(httpreq)
	if err != nil {
		return
	}
	rsp, err = ioutil.ReadAll(resp.Body)
	return
}
