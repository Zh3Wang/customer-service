package service

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	zh2 "gopkg.in/go-playground/validator.v9/translations/zh"
	"net/http"
)

const (
	API_SUCCESS = 0
	API_FAIL    = -1
)

type RespData struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

// 验证参数
func Validate(in interface{}) error {
	validate := validator.New()
	z := zh.New()

	uni := ut.New(z, z)
	trans, _ := uni.GetTranslator("zh")

	_ = zh2.RegisterDefaultTranslations(validate, trans)

	if err := validate.Struct(in); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(err.Translate(trans))
		}
	}
	return nil
}

func Respond(w http.ResponseWriter, code int, msg string, data interface{}) {
	var resp = RespData{
		Code: code,
		Msg:  msg,
		Data: data,
	}
	r, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(r)

	return
}

