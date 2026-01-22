package middleware

import (
	"errors"

	"github.com/kataras/golog"
	"github.com/wenlng/go-captcha-assets/resources/imagesv2"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/slide"
)

type CaptchaResult struct {
	CaptchaId string `json:"captchaId"`
	Image     string `json:"image"`
	Tile      string `json:"tile"`
	X         int    `json:"-"` // X值不返回前端
	Y         int    `json:"y"`
	Width     int    `json:"width"`
}

type CaptchaCache struct {
	X int `json:"x"`
	Y int `json:"y"`
}

var Captcha slide.Captcha

// 初始化验证码组件
func InitCaptcha() error {
	builder := slide.NewBuilder()

	imgs, err := imagesv2.GetImages()
	if err != nil {
		golog.Error("验证码图片加载失败：", err)
		return err
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		golog.Error("验证码图片拼图加载失败：", err)
		return err
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := range graphs {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	Captcha = builder.Make()

	return nil
}

// 获取验证码
func CaptchaGet() (CaptchaResult, error) {
	captchaResult := CaptchaResult{}

	if Captcha == nil {
		return captchaResult, errors.New("验证码组件未初始化")
	}
	captData, err := Captcha.Generate()
	if err != nil {
		return captchaResult, err
	}

	image, err := captData.GetMasterImage().ToBase64()
	if err != nil {
		return captchaResult, err
	}

	tile, err := captData.GetTileImage().ToBase64()
	if err != nil {
		return captchaResult, err
	}

	captchaResult.Image = image
	captchaResult.Tile = tile
	captchaResult.X = captData.GetData().X
	captchaResult.Y = captData.GetData().Y
	captchaResult.Width = captData.GetData().Width

	return captchaResult, err
}
