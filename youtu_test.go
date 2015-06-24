/*
* File Name:	youtu_test.go
* Description:
* Author:	Chapman Ou <ochapman.cn@gmail.com>
* Created:	2015-06-21
 */

package youtu

import "testing"

//Update as if you want to test your own app
var as = AppSign{
	app_id:     12345678,
	secret_id:  "your_secret_id",
	secret_key: "your_secret_key",
	expired:    1436353609,
	user_id:    "your_qq_id",
}

func TestDetectFace(t *testing.T) {
	yt := Init(as, DefaultHost)
	imgData, err := EncodeImage("testdata/imageA.jpg")
	if err != nil {
		t.Errorf("EncodeImage failed: %s", err)
		return
	}
	dfr, err := yt.DetectFace(imgData, DetectMode_Normal)
	if err != nil {
		t.Errorf("Detect face faild: %s", err)
		return
	}
	t.Logf("dfr: %#v\n", dfr)
}

func TestFaceCompare(t *testing.T) {
	yt := Init(as, DefaultHost)
	imageA, err := EncodeImage("testdata/imageA.jpg")
	if err != nil {
		t.Errorf("Encode imageA failed: %s\n", err)
		return
	}
	imageB, err := EncodeImage("testdata/imageB.jpg")
	if err != nil {
		t.Errorf("Encode imageB failed: %s\n", err)
		return
	}
	fcr, err := yt.FaceCompare(imageA, imageB)
	if err != nil {
		t.Errorf("FaceCompare failed: %s\n", err)
		return
	}
	t.Logf("fcr: %#v\n", fcr)
}

func TestFaceVerify(t *testing.T) {
	yt := Init(as, DefaultHost)
	image, err := EncodeImage("testdata/imageA.jpg")
	if err != nil {
		t.Errorf("EncodeImage failed: %s\n", err)
		return
	}
	person_id := "1045684262752288767"
	fvr, err := yt.FaceVerify(image, person_id)
	if err != nil {
		t.Errorf("FaceVerify failed: %s\n", err)
		return
	}
	t.Logf("fvr: %#v\n", fvr)
}

func TestFaceIdentify(t *testing.T) {
	yt := Init(as, DefaultHost)
	image, err := EncodeImage("testdata/imageA.jpg")
	if err != nil {
		t.Errorf("EncodeImage failed: %s\n", err)
		return
	}
	group_id := "tencent"
	fir, err := yt.FaceIdentify(image, group_id)
	if err != nil {
		t.Errorf("FaceIdentify failed: %s\n", err)
		return
	}
	t.Logf("fir: %#v\n", fir)
}

func TestNewPerson(t *testing.T) {
	yt := Init(as, DefaultHost)
	image, err := EncodeImage("testdata/imageA.jpg")
	if err != nil {
		t.Errorf("EncodeImage failed: %s\n", err)
		return
	}
	group_ids := []string{"tencent"}
	npr, err := yt.NewPerson(image, "ochapman", group_ids, "ochapman", "person tag")
	if err != nil && npr.Errormsg != "ERROR_PERSON_EXISTED" {
		t.Errorf("NewPerson failed: %s\n", err)
		return
	}
	t.Logf("npr: %#v\n", npr)
}

func TestDelPerson(t *testing.T) {
	yt := Init(as, DefaultHost)
	dpr, err := yt.DelPerson("ochapman")
	if err != nil {
		t.Errorf("DelPerson failed: %s\n", err)
		return
	}
	t.Logf("dpr: %#v\n", dpr)
}

func TestAddFace(t *testing.T) {
	yt := Init(as, DefaultHost)
	image, err := EncodeImage("testdata/imageA.jpg")
	if err != nil {
		t.Errorf("EncodeImage failed: %s\n", err)
		return
	}
	person_id := "ochapman"
	images := []string{image}
	tag := "face tag"
	afr, err := yt.AddFace(images, person_id, tag)
	if err != nil {
		t.Errorf("AddFace failed: %s\n", err)
		return
	}
	t.Logf("afr: %#v\n", afr)
}

func TestDelFace(t *testing.T) {
	yt := Init(as, DefaultHost)
	person_id := "ochapman"
	face_ids := []string{"123456"}
	dfr, err := yt.DelFace(person_id, face_ids)
	if err != nil {
		t.Errorf("DelFace failed: %s\n", err)
		return
	}
	t.Logf("dfr: %#v\n", dfr)
}

func TestSetInfo(t *testing.T) {
	yt := Init(as, DefaultHost)
	person_id := "ochapman"
	person_name := "ochapman_new"
	tag := "SetInfo tag"
	sir, err := yt.SetInfo(person_id, person_name, tag)
	if err != nil {
		t.Errorf("SetInfo failed: %s\n", err)
		return
	}
	t.Logf("sir: %#v\n", sir)
}

func TestGetInfo(t *testing.T) {
	yt := Init(as, DefaultHost)
	person_id := "ochapman"
	gir, err := yt.GetInfo(person_id)
	if err != nil {
		t.Errorf("GetInfo failed: %s\n", err)
		return
	}
	t.Logf("sir %#v\n", gir)
}

func TestGetGroupIDs(t *testing.T) {
	yt := Init(as, DefaultHost)
	ggr, err := yt.GetGroupIDs()
	if err != nil {
		t.Errorf("GetGroupIDs failed: %s\n", err)
		return
	}
	t.Logf("ggr %#v\n", ggr)

}

func TestGetPersonIDs(t *testing.T) {
	yt := Init(as, DefaultHost)
	gpr, err := yt.GetPersonIDs("12345")
	if err != nil {
		t.Errorf("GetPersonIDs failed: %s\n", err)
		return
	}
	t.Logf("gpr: %#v\n", gpr)
}

func TestGetPersonIDs(t *testing.T) {
	yt := Init(as, DefaultHost)
	gfr, err := yt.GetFaceIDs("12345")
	if err != nil {
		t.Errorf("GetFaceIDs failed: %s\n", err)
		return
	}
	t.Logf("gfr: %#v\n", gfr)
}

func TestGetFaceInfo(t *testing.T) {
	yt := Init(as, DefaultHost)
	gfr, err := yt.GetFaceInfo("12345")
	if err != nil {
		t.Errorf("GetFaceInfo failed: %s\n", err)
		return
	}
	t.Logf("gfr: %#v\n", gfr)
}
