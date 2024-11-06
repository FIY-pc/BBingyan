package util

import (
	"github.com/FIY-pc/BBingyan/internal/config"
	"github.com/FIY-pc/BBingyan/internal/util"
	"testing"
)

func TestGenerateCaptcha(t *testing.T) {
	for i := 0; i < 5; i++ {
		captcha := util.GenerateCaptcha()
		t.Log(captcha)
	}
}

func TestSendCaptcha(t *testing.T) {
	config.InitConfig()
	_ = util.SendCaptcha("2953997758@qq.com", util.GenerateCaptcha())
}
