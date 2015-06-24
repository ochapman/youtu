##youtu
`youtu`是http://open.youtu.qq.com 提供接口的Go实现版本


### 下载:

```bash
go get https://github.com/ochapman/youtu
```

### 使用:
1. 到http://open.youtu.qq.com 注册成为开发者
2. 注册应用，获取开发者密钥

### 使用例子:

```go
package main

import (
	"fmt"
	"os"

	"github.com/ochapman/youtu"
)

func main() {
	//Register your app on http://open.youtu.qq.com
	//Get the following details
	app_id := uint32(12345678)
	secret_id := "your_secret_id"
	secret_key := "your_secret_key"
	expired := uint32(1436353609)
	user_id := "your_qq_id"

	as, err := youtu.NewAppSign(app_id, secret_id, secret_key, expired, user_id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "NewAppSign() failed: %s\n", err)
		return
	}
	imgData, err := youtu.EncodeImage("testdata/imageA.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "EncodeImage() failed: %s\n", err)
		return
	}

	yt := youtu.Init(as, youtu.DefaultHost)
	df, err := yt.DetectFace(imgData, youtu.DetectMode_Normal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DetectFace() failed: %s", err)
		return
	}
	fmt.Printf("df: %#v\n", df)
}
```

###文档
[![GoDoc](https://godoc.org/github.com/ochapman/youtu?status.svg)](https://godoc.org/github.com/ochapman/youtu)

###构建
[![Build Status](https://travis-ci.org/ochapman/youtu.svg?branch=master)](https://travis-ci.org/ochapman/youtu)
