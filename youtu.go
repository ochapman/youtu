/*
* File Name:	youtu.go
* Description:  http://open.youtu.qq.com API
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-06-19
 */

package youtu

import (
	"encoding/base64"
	"errors"
	"strconv"
)

const (
	//UserIDMaxLen 用户ID的最大长度
	UserIDMaxLen = 110
)

const expiredInterval = 1000

var (
	//ErrUserIDTooLong 用户ID过长错误
	ErrUserIDTooLong = errors.New("user id too long")
)

var (
	//DefaultHost 默认host
	DefaultHost = "api.youtu.qq.com"
)

//AppSign 应用签名鉴权
type AppSign struct {
	appID     uint32 //接入优图服务时,生成的唯一id, 用于唯一标识接入业务
	secretID  string //标识api鉴权调用者的密钥身份
	secretKey string //用于加密签名字符串和服务器端验证签名字符串的密钥，secret_key 必须严格保管避免泄露
	userID    string //接入业务自行定义的用户id，用于唯一标识一个用户, 登陆开发者账号的QQ号码
}

//NewAppSign 新建应用签名
func NewAppSign(appID uint32, secretID string, secretKey string, userID string) (as AppSign, err error) {
	if len(userID) > UserIDMaxLen {
		err = ErrUserIDTooLong
		return
	}
	as = AppSign{
		appID:     appID,
		secretID:  secretID,
		secretKey: secretKey,
		userID:    userID,
	}
	return
}

//Youtu 存储签名和host
type Youtu struct {
	appSign AppSign
	host    string
	debug   bool //Default false
}

func (y *Youtu) appID() string {
	return strconv.Itoa(int(y.appSign.appID))
}

//Init Youtu初始化
func Init(appSign AppSign, host string) *Youtu {
	return &Youtu{
		appSign: appSign,
		host:    host,
		debug:   false,
	}
}

//detectMode 检测模式，分正常和大脸
type detectMode int

const (
	//detectModeNormal 正常模式
	detectModeNormal detectMode = iota
	//detectModeBigFace 大脸模式
	detectModeBigFace
)

func mode(isBigFace bool) detectMode {
	if isBigFace {
		return detectModeBigFace
	}
	return detectModeNormal
}

// SetDebug For Debug
func (y *Youtu) SetDebug(isDebug bool) {
	y.debug = isDebug
}

type detectFaceReq struct {
	AppID string     `json:"app_id"`         //App的 API ID
	Image string     `json:"image"`          //base64编码的二进制图片数据
	Mode  detectMode `json:"mode,omitempty"` //检测模式 0/1 正常/大脸模式
}

//Face 脸参数
type Face struct {
	FaceID     string  `json:"face_id"`    //人脸标识
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

//DetectFaceRsp 脸检测返回
type DetectFaceRsp struct {
	SessionID   string `json:"session_id"`   //相应请求的session标识符，可用于结果查询
	ImageID     string `json:"image_id"`     //系统中的图片标识符，用于标识用户请求中的图片
	ImageWidth  int32  `json:"image_width"`  //请求图片的宽度
	ImageHeight int32  `json:"image_height"` //请求图片的高度
	Face        []Face `json:"face"`         //被检测出的人脸Face的列表
	ErrorCode   int    `json:"errorcode"`    //返回状态值
	ErrorMsg    string `json:"errormsg"`     //返回错误消息
}

//DetectFace 检测给定图片(Image)中的所有人脸(Face)的位置和相应的面部属性。
//位置包括(x, y, w, h)，面部属性包括性别(gender), 年龄(age),
//表情(expression), 眼镜(glass)和姿态(pitch，roll，yaw).
func (y *Youtu) DetectFace(imageData []byte, isBigFace bool) (rsp DetectFaceRsp, err error) {
	b64Image := base64.StdEncoding.EncodeToString(imageData)
	req := detectFaceReq{
		AppID: strconv.Itoa(int(y.appSign.appID)),
		Image: b64Image,
		Mode:  mode(isBigFace),
	}
	err = y.interfaceRequest("detectface", req, &rsp)
	return
}

type faceShapeReq struct {
	AppID string     `json:"app_id"`         //App的 API ID
	Image string     `json:"image"`          //base64编码的二进制图片数据
	Mode  detectMode `json:"mode,omitempty"` //检测模式 0/1 正常/大脸模式
}

type pos struct {
	X int `json:"x"`
	Y int `json:"y"`
}

//FaceShape 五官定位
type FaceShape struct {
	FaceProfile  []pos `json:"face_profile"`  //描述脸型轮廓的21点
	LeftEye      []pos `json:"left_eye"`      //描述左眼轮廓的8点
	RightEye     []pos `json:"right_eye"`     //描述右眼轮廓的8点
	LeftEyebrow  []pos `json:"left_eyebrow"`  //描述左眉轮廓的8点
	RightEyebrow []pos `json:"right_eyebrow"` //描述右眉轮廓的8点
	Mouth        []pos `json:"mouth"`         //描述嘴巴轮廓的22点
	Nose         []pos `json:"nose"`          //描述鼻子轮廓的13点
}

// FaceShapeRsp 返回
type FaceShapeRsp struct {
	SessionID   string      `json:"session_id"`   //相应请求的session标识符，可用于结果查询
	FaceShape   []FaceShape `json:"face_shape"`   //人脸轮廓结构体，包含所有人脸的轮廓点
	ImageWidth  int         `json:"image_width"`  //请求图片的宽度
	ImageHeight int         `json:"image_height"` //请求图片的高度
	ErrorCode   int         `json:"errorcode"`    //返回状态值
	ErrorMsg    string      `json:"errormsg"`     //返回错误消息
}

//FaceShape 对请求图片进行五官定位，计算构成人脸轮廓的88个点，包括眉毛（左右各8点）、眼睛（左右各8点）、鼻子（13点）、嘴巴（22点）、脸型轮廓（21点）
func (y *Youtu) FaceShape(image []byte, isBigFace bool) (rsp FaceShapeRsp, err error) {
	b64Image := base64.StdEncoding.EncodeToString(image)
	req := faceShapeReq{
		AppID: strconv.Itoa(int(y.appSign.appID)),
		Image: b64Image,
		Mode:  mode(isBigFace),
	}
	err = y.interfaceRequest("faceshape", req, &rsp)
	return
}

type faceCompareReq struct {
	AppID  string `json:"app_id"`
	ImageA string `json:"imageA"` //使用base64编码的二进制图片数据A
	ImageB string `json:"imageB"` //使用base64编码的二进制图片数据B
}

//FaceCompareRsp 脸比较返回
type FaceCompareRsp struct {
	EyebrowSim float32 `json:"eyebrow_sim"` //眉毛的相似度。
	EyeSim     float32 `json:"eye_sim"`     //眼睛的相似度
	NoseSim    float32 `json:"nose_sim"`    //鼻子的相似度
	MouthSim   float32 `json:"mouth_sim"`   //嘴巴的相似度
	Similarity float32 `json:"similarity"`  //两个face的相似度
	ErrorCode  int32   `json:"errorcode"`   //返回状态码
	ErrorMsg   string  `json:"errormsg"`    //返回错误消息
}

//FaceCompare 计算两个Face的相似性以及五官相似度
func (y *Youtu) FaceCompare(imageA, imageB []byte) (rsp FaceCompareRsp, err error) {
	b64ImageA := base64.StdEncoding.EncodeToString(imageA)
	b64ImageB := base64.StdEncoding.EncodeToString(imageB)
	req := faceCompareReq{
		AppID:  y.appID(),
		ImageA: b64ImageA,
		ImageB: b64ImageB,
	}
	err = y.interfaceRequest("facecompare", req, &rsp)
	return
}

type faceVerifyReq struct {
	AppID    string `json:"app_id"`    //App的 API ID
	Image    string `json:"image"`     //使用base64编码的二进制图片数据
	PersonID string `json:"person_id"` //待验证的Person
}

//FaceVerifyRsp 脸验证返回
type FaceVerifyRsp struct {
	Ismatch    bool    `json:"ismatch"`    //两个输入是否为同一人的判断
	Confidence float32 `json:"confidence"` //系统对这个判断的置信度。
	SessionID  string  `json:"session_id"` //相应请求的session标识符，可用于结果查询
	ErrorCode  int32   `json:"errorcode"`  //返回状态码
	ErrorMsg   string  `json:"errormsg"`   //返回错误消息
}

//FaceVerify 给定一个Face和一个Person，返回是否是同一个人的判断以及置信度。
func (y *Youtu) FaceVerify(personID string, image []byte) (rsp FaceVerifyRsp, err error) {
	b64Image := base64.StdEncoding.EncodeToString(image)
	req := faceVerifyReq{
		AppID:    y.appID(),
		Image:    b64Image,
		PersonID: personID,
	}
	err = y.interfaceRequest("faceverify", req, &rsp)
	return
}

type faceIdentifyReq struct {
	AppID   string `json:"app_id"`   //App的 API ID
	GroupID string `json:"group_id"` //候选人组id
	Image   string `json:"image"`    //使用base64编码的二进制图片数据
}

//FaceIdentifyRsp 脸识别返回
type FaceIdentifyRsp struct {
	SessionID  string  `json:"session_id"` //相应请求的session标识符，可用于结果查询
	PersonID   string  `json:"person_id"`  //识别结果，person_id
	FaceID     string  `json:"face_id"`    //识别的face_id
	Confidence float32 `json:"confidence"` //置信度
	ErrorCode  int     `json:"errorcode"`  //返回状态码
	ErrorMsg   string  `json:"errormsg"`   //返回错误消息
}

//FaceIdentify 对于一个待识别的人脸图片，在一个Group中识别出最相似的Person作为其身份返回
func (y *Youtu) FaceIdentify(groupID string, image []byte) (rsp FaceIdentifyRsp, err error) {
	b64Image := base64.StdEncoding.EncodeToString(image)
	req := faceIdentifyReq{
		AppID:   y.appID(),
		GroupID: groupID,
		Image:   b64Image,
	}
	err = y.interfaceRequest("faceidentify", req, &rsp)
	return
}

type newPersonReq struct {
	AppID      string   `json:"app_id"` //App的 API ID
	Image      string   `json:"image"`  //使用base64编码的二进制图片数据
	PersonID   string   `json:"person_id"`
	GroupIDs   []string `json:"group_ids"`             // 	加入到组的列表
	PersonName string   `json:"person_name,omitempty"` //名字
	Tag        string   `json:"tag,omitempty"`         //备注信息
}

//NewPersonRsp 个体创建返回
type NewPersonRsp struct {
	SessionID  string `json:"session_id"`  //相应请求的session标识符
	SucGroup   int    `json:"suc_group"`   //成功被加入的group数量
	SucFace    int    `json:"suc_face"`    //成功加入的face数量
	PersonName string `json:"person_name"` //相应person的name
	PersonID   string `json:"person_id"`   //相应person的id
	FaceID     string `json:"face_id"`     //创建所用图片生成的face_id
	ErrorCode  int    `json:"errorcode"`   //返回码
	ErrorMsg   string `json:"errormsg"`    //返回错误消息
}

//NewPerson 创建一个Person，并将Person放置到group_ids指定的组当中
func (y *Youtu) NewPerson(personID string, personName string, groupIDs []string, image []byte, tag string) (rsp NewPersonRsp, err error) {
	b64Image := base64.StdEncoding.EncodeToString(image)
	req := newPersonReq{
		AppID:      y.appID(),
		PersonID:   personID,
		Image:      b64Image,
		GroupIDs:   groupIDs,
		PersonName: personName,
		Tag:        tag,
	}
	err = y.interfaceRequest("newperson", req, &rsp)
	return
}

type delPersonReq struct {
	AppID    string `json:"app_id"`
	PersonID string `json:"person_id"` //待删除个体ID
}

//DelPersonRsp 删除个体返回
type DelPersonRsp struct {
	SessionID string `json:"session_id"` //相应请求的session标识符
	Deleted   int    `json:"deleted"`    //成功删除的Person数量
	ErrorCode int    `json:"errorcode"`  //返回状态码
	ErrorMsg  string `json:"errormsg"`   //返回错误消息
}

//DelPerson 删除一个Person
func (y *Youtu) DelPerson(personID string) (rsp DelPersonRsp, err error) {
	req := delPersonReq{
		AppID:    y.appID(),
		PersonID: personID,
	}
	err = y.interfaceRequest("delperson", req, &rsp)
	return
}

type addFaceReq struct {
	AppID    string   `json:"app_id"`        //App的 API ID
	PersonID string   `json:"person_id"`     //String 	待增加人脸的个体id
	Images   []string `json:"images"`        //base64编码的二进制图片数据构成的数组
	Tag      string   `json:"tag,omitempty"` //备注信息
}

//AddFaceRsp 增加人脸返回
type AddFaceRsp struct {
	SessionID string   `json:"session_id"` //相应请求的session标识符
	Added     int      `json:"added"`      //成功加入的face数量
	FaceIDs   []string `json:"face_ids"`   //增加的人脸ID列表
	ErrorCode int      `json:"errorcode"`  //返回状态码
	ErrorMsg  string   `json:"errormsg"`   //返回错误消息
}

//AddFace 将一组Face加入到一个Person中。注意，一个Face只能被加入到一个Person中。
//一个Person最多允许包含10000个Face
func (y *Youtu) AddFace(personID string, images [][]byte, tag string) (rsp AddFaceRsp, err error) {
	b64Images := make([]string, len(images))
	for i, img := range images {
		b64Images[i] = base64.StdEncoding.EncodeToString([]byte(img))
	}
	req := addFaceReq{
		AppID:    y.appID(),
		Images:   b64Images,
		PersonID: personID,
		Tag:      tag,
	}
	err = y.interfaceRequest("addface", req, &rsp)
	return
}

type delFaceReq struct {
	AppID    string   `json:"app_id"`    //App的 API ID
	PersonID string   `json:"person_id"` //待删除人脸的person ID
	FaceIDs  []string `json:"face_ids"`  //删除人脸id的列表
}

//DelFaceRsp 删除人脸返回
type DelFaceRsp struct {
	SessonID  string `json:"session_id"` //相应请求的session标识符
	Deleted   int32  `json:"deleted"`    //成功删除的face数量
	ErrorCode int32  `json:"errorcode"`  //返回状态码
	ErrorMsg  string `json:"errormsg"`   //返回错误消息
}

//DelFace 删除一个person下的face，包括特征，属性和face_id.
func (y *Youtu) DelFace(personID string, faceIDs []string) (rsp DelFaceRsp, err error) {
	req := delFaceReq{
		AppID:    y.appID(),
		PersonID: personID,
		FaceIDs:  faceIDs,
	}
	err = y.interfaceRequest("delface", req, &rsp)
	return
}

type setInfoReq struct {
	AppID      string `json:"app_id"` //App的 API ID
	PersonID   string `json:"person_id"`
	PersonName string `json:"person_name,omitempty"` //新的name
	Tag        string `json:"tag,omitempty"`         //备注信息
}

//SetInfoRsp 设置信息返回
type SetInfoRsp struct {
	sessionID string `json:"session_id"` //相应请求的session标识符
	personID  string `json:"person_id"`  //相应person的id
	errorcode int32  `json:"errorcode"`  //返回状态码
	errormsg  string `json:"errormsg"`   //返回错误消息
}

//SetInfo 设置Person的name.
func (y *Youtu) SetInfo(personID string, personName string, tag string) (rsp SetInfoRsp, err error) {
	req := setInfoReq{
		AppID:      y.appID(),
		PersonID:   personID,
		PersonName: personName,
		Tag:        tag,
	}
	err = y.interfaceRequest("setinfo", req, &rsp)
	return
}

type getInfoReq struct {
	AppID    string `json:"app_id"`    //App的 API ID
	PersonID string `json:"person_id"` //待查询个体的ID
}

//GetInfoRsp 获取信息返回
type GetInfoRsp struct {
	PersonName string   `json:"person_name"` //相应person的name
	PersonID   string   `json:"person_id"`   //相应person的id
	GroupIDs   []string `json:"group_ids"`   //包含此个体的组列表
	FaceIDs    []string `json:"face_ids"`    //包含的人脸列表
	SessionID  string
	ErrorCode  int    `json:"errorcode"` //返回状态码
	ErrorMsg   string `json:"errormsg"`  //返回错误消息
}

//GetInfo 获取一个Person的信息, 包括name, id, tag, 相关的face, 以及groups等信息。
func (y *Youtu) GetInfo(personID string) (rsp GetInfoRsp, err error) {
	req := getInfoReq{
		AppID:    y.appID(),
		PersonID: personID,
	}
	err = y.interfaceRequest("getinfo", req, &rsp)
	return
}

type getGroupIDsReq struct {
	AppID string `json:"app_id"` //App的 API ID
}

//GetGroupIDsRsp 获取组ID返回
type GetGroupIDsRsp struct {
	GroupIDs  []string `json:"group_ids"` //相应app_id的group_id列表
	ErrorCode int32    `json:"errorcode"` //返回状态码
	ErrorMsg  string   `json:"errormsg"`  //返回错误消息
}

//GetGroupIDs 获取一个appId下所有group列表
func (y *Youtu) GetGroupIDs() (rsp GetGroupIDsRsp, err error) {
	req := getGroupIDsReq{
		AppID: y.appID(),
	}
	err = y.interfaceRequest("getgroupids", req, &rsp)
	return
}

type getPersonIDsReq struct {
	AppID   string `json:"app_id"`   //App的 API ID
	GroupID string `json:"group_id"` //组id
}

//GetPersonIDsRsp 获取个人ID返回
type GetPersonIDsRsp struct {
	PersonIDs []string `json:"person_ids"` //相应person的id列表
	ErrorCode int32    `json:"errorcode"`  //返回状态码
	ErrorMsg  string   `json:"errormsg"`   //返回错误消息
}

//GetPersonIDs 获取一个组Group中所有person列表
func (y *Youtu) GetPersonIDs(groupID string) (rsp GetPersonIDsRsp, err error) {
	req := getPersonIDsReq{
		AppID:   y.appID(),
		GroupID: groupID,
	}
	err = y.interfaceRequest("getpersonids", req, &rsp)
	return
}

type getFaceIDsReq struct {
	AppID    string `json:"app_id"`    //App的 API ID
	PersonID string `json:"person_id"` //个体id
}

//GetFaceIDsRsp 获取脸ID返回
type GetFaceIDsRsp struct {
	FaceIDs   []string `json:"face_ids"`  //相应face的id列表
	ErrorCode int32    `json:"errorcode"` //返回状态码
	ErrorMsg  string   `json:"errormsg"`  //返回错误消息
}

//GetFaceIDs 获取一个组person中所有face列表
func (y *Youtu) GetFaceIDs(personID string) (rsp GetFaceIDsRsp, err error) {
	req := getFaceIDsReq{
		AppID:    y.appID(),
		PersonID: personID,
	}
	err = y.interfaceRequest("getfaceids", req, &rsp)
	return
}

type getFaceInfoReq struct {
	AppID  string `json:"app_id"`  //App的 API ID
	FaceID string `json:"face_id"` //人脸id
}

//GetFaceInfoRsp 获取脸部信息返回
type GetFaceInfoRsp struct {
	FaceInfo  Face   `json:"face_info"` //人脸信息
	ErrorCode int32  `json:"errorcode"` //返回状态码
	ErrorMsg  string `json:"errormsg"`  //返回错误消息
}

//GetFaceInfo 获取一个face的相关特征信息
func (y *Youtu) GetFaceInfo(faceID string) (rsp GetFaceInfoRsp, err error) {
	req := getFaceInfoReq{
		AppID:  y.appID(),
		FaceID: faceID,
	}
	err = y.interfaceRequest("getfaceinfo", req, &rsp)
	return
}
