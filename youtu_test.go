/*
* File Name:	youtu_test.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-06-21
 */

package youtu

import (
	"io/ioutil"
	"testing"
)

//Update as if you want to test your own app
var as = AppSign{
	appID:     1000061,
	secretID:  "AKID4Bhs9vqYT6mHa9TkIrAe7w5oijOCEjql",
	secretKey: "P2VTKNvTAnYNwBrqXbgxRSFQs6FTEhNJ",
	//expired:   1440207436 + 5000,
	userID: "3041722595",
}

const testDataDir = "./testdata/"

var yt = Init(as, DefaultHost)

func TestDetectFace(t *testing.T) {
	imgData, err := ioutil.ReadFile(testDataDir + "imageA.jpg")
	if err != nil {
		t.Errorf("ReadFile failed: %s", err)
		return
	}
	rsp, err := yt.DetectFace(imgData, false)
	if err != nil {
		t.Errorf("Detect face faild: %s", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestFaceShape(t *testing.T) {
	imgData, err := ioutil.ReadFile(testDataDir + "faceshape.jpg")
	if err != nil {
		t.Errorf("ReadFile failed: %s\n", err)
		return
	}
	rsp, err := yt.FaceShape(imgData, false)
	if err != nil {
		t.Errorf("FaceShape failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestFaceCompare(t *testing.T) {
	imageA, err := ioutil.ReadFile(testDataDir + "imageA.jpg")
	if err != nil {
		t.Errorf("Encode imageA failed: %s\n", err)
		return
	}
	imageB, err := ioutil.ReadFile(testDataDir + "imageB.jpg")
	if err != nil {
		t.Errorf("Encode imageB failed: %s\n", err)
		return
	}
	rsp, err := yt.FaceCompare(imageA, imageB)
	if err != nil {
		t.Errorf("FaceCompare failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestFaceVerify(t *testing.T) {
	image, err := ioutil.ReadFile(testDataDir + "imageA.jpg")
	if err != nil {
		t.Errorf("ioutil.ReadFile failed: %s\n", err)
		return
	}
	personID := "1045684262752288767"
	rsp, err := yt.FaceVerify(personID, image)
	if err != nil {
		t.Errorf("FaceVerify failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestFaceIdentify(t *testing.T) {
	image, err := ioutil.ReadFile(testDataDir + "imageA.jpg")
	if err != nil {
		t.Errorf("ioutil.ReadFile failed: %s\n", err)
		return
	}
	groupID := "tencent"
	rsp, err := yt.FaceIdentify(groupID, image)
	if err != nil {
		t.Errorf("FaceIdentify failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestNewPerson(t *testing.T) {
	image, err := ioutil.ReadFile(testDataDir + "imageA.jpg")
	if err != nil {
		t.Errorf("ioutil.ReadFile failed: %s\n", err)
		return
	}
	groupIDs := []string{"tencent"}
	rsp, err := yt.NewPerson("ochapman", "ochapman", groupIDs, image, "person tag")
	if err != nil && rsp.ErrorMsg != "ERROR_PERSON_EXISTED" {
		t.Errorf("NewPerson failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestDelPerson(t *testing.T) {
	rsp, err := yt.DelPerson("ochapman")
	if err != nil {
		t.Errorf("DelPerson failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestAddFace(t *testing.T) {
	image, err := ioutil.ReadFile(testDataDir + "imageA.jpg")
	if err != nil {
		t.Errorf("ioutil.ReadFile failed: %s\n", err)
		return
	}
	personID := "ochapman"
	images := [][]byte{image}
	tag := "face tag"
	rsp, err := yt.AddFace(personID, images, tag)
	if err != nil {
		t.Errorf("AddFace failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestDelFace(t *testing.T) {
	personID := "ochapman"
	faceIDs := []string{"123456"}
	rsp, err := yt.DelFace(personID, faceIDs)
	if err != nil {
		t.Errorf("DelFace failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestSetInfo(t *testing.T) {
	personID := "ochapman"
	personName := "ochapman_new"
	tag := "SetInfo tag"
	rsp, err := yt.SetInfo(personID, personName, tag)
	if err != nil {
		t.Errorf("SetInfo failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestGetInfo(t *testing.T) {
	personID := "ochapman"
	rsp, err := yt.GetInfo(personID)
	if err != nil {
		t.Errorf("GetInfo failed: %s\n", err)
		return
	}
	t.Logf("rsp %#v\n", rsp)
}

func TestGetGroupIDs(t *testing.T) {
	rsp, err := yt.GetGroupIDs()
	if err != nil {
		t.Errorf("GetGroupIDs failed: %s\n", err)
		return
	}
	t.Logf("rsp %#v\n", rsp)

}

func TestGetPersonIDs(t *testing.T) {
	rsp, err := yt.GetPersonIDs("12345")
	if err != nil {
		t.Errorf("GetPersonIDs failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestGetFaceIDs(t *testing.T) {
	rsp, err := yt.GetFaceIDs("12345")
	if err != nil {
		t.Errorf("GetFaceIDs failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}

func TestGetFaceInfo(t *testing.T) {
	rsp, err := yt.GetFaceInfo("12345")
	if err != nil {
		t.Errorf("GetFaceInfo failed: %s\n", err)
		return
	}
	t.Logf("rsp: %#v\n", rsp)
}
