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
}

func (y *Youtu) AppId() string {
	return strconv.Itoa(int(y.app_sign.app_id))
}

func Init(appSign AppSign, host string) *Youtu {
	return &Youtu{
		app_sign: appSign,
		host:     host,
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
	Face_id    string  `json:"face_id"`    //人脸标识
	X          int32   `json:"x"`          //人脸框左上角x
	Y          int32   `json:"y"`          //人脸框左上角y
	Width      float32 `json:"width"`      //人脸框宽度
	Height     float32 `json:"height"`     //人脸框高度
	Gender     int32   `json:"gender"`     //性别 [0/(female)~100(male)]
	Age        int32   `json:"age"`        //年龄 [0~100]
	Expression int32   `json:"expression"` //object 	微笑[0(normal)~50(smile)~100(laugh)]
	Glass      bool    `json:"glass"`      //是否有眼镜 [true,false]
	Pitch      int32   `json:"pitch"`      //上下偏移[-30,30]
	Yaw        int32   `json:"yaw"`        //左右偏移[-30,30]
	Roll       int32   `json:"roll"`       //平面旋转[-180,180]
}

type DetectFaceRsp struct {
	Session_id   string `json:"session_id"`   //相应请求的session标识符，可用于结果查询
	Image_id     string `json:"image_id"`     //系统中的图片标识符，用于标识用户请求中的图片
	Image_width  int32  `json:"image_width"`  //请求图片的宽度
	Image_height int32  `json:"image_height"` //请求图片的高度
	Face         []Face `json:"face"`         //被检测出的人脸Face的列表
	ErrorCode    int    `json:"errorcode"`    //返回状态值
	ErrorMsg     string `json:"errormsg"`     //返回错误消息
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

type FaceCompareReq struct {
	App_id string `json:"app_id"`
	ImageA string `json:"imageA"` //使用base64编码的二进制图片数据A
	ImageB string `json:"imageB"` //使用base64编码的二进制图片数据B
}

type FaceCompareRsp struct {
	Eyebrow_sim float32 `json:"eyebrow_sim"` //眉毛的相似度。
	Eye_sim     float32 `json:"eye_sim"`     //眼睛的相似度
	Nose_sim    float32 `json:"nose_sim"`    //鼻子的相似度
	Mouth_sim   float32 `json:"mouth_sim"`   //嘴巴的相似度
	Similarity  float32 `json:"similarity"`  //两个face的相似度
	Errorcode   int32   `json:"errorcode"`   //返回状态码
	Errormsg    string  `json:"errormsg"`    //返回错误消息
}

func (y *Youtu) FaceCompare(imageA, imageB string) (fcr FaceCompareRsp, err error) {
	url := "http://" + y.host + "/youtu/api/facecompare"
	req := FaceCompareReq{
		App_id: y.AppId(),
		ImageA: imageA,
		ImageB: imageB,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return
	}
	rsp, err := y.get(url, string(data))
	if err != nil {
		return
	}
	err = json.Unmarshal(rsp, &fcr)
	if err != nil {
		return fcr, fmt.Errorf("json.Unmarshal() rsp: %s failed: %s\n", rsp, err)
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
