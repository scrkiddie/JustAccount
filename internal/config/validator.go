package config

import (
	"github.com/go-playground/validator/v10"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func NewValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	err := v.RegisterValidation("image", ValidateImage)
	if err != nil {
		panic(err.Error())
	}
	err = v.RegisterValidation("not_contain_space", ValidateNotContainSpace)
	if err != nil {
		panic(err.Error())
	}
	err = v.RegisterValidation("is_password_format", ValidatePasswordFormat)
	if err != nil {
		panic(err.Error())
	}
	return v
}

func ValidateImage(fl validator.FieldLevel) bool {
	file := fl.Field().Interface().(multipart.FileHeader)

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return false
	}

	if file.Size > 1024*1024 {
		return false
	}

	dim := fl.Param()
	parts := strings.Split(dim, "x")
	if dim != "" {
		if len(parts) != 2 {
			return false
		}
		width, errW := strconv.Atoi(parts[0])
		height, errH := strconv.Atoi(parts[1])
		if errW != nil || errH != nil {
			return false
		}

		imgFile, err := file.Open()
		if err != nil {
			return false
		}
		defer imgFile.Close()
		cfg, _, err := image.DecodeConfig(imgFile)
		if err != nil {
			return false
		}
		if cfg.Width != width || cfg.Height != height {
			return false
		}
	}
	return true
}

func ValidateNotContainSpace(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	return !strings.Contains(field, " ")
}

func ValidatePasswordFormat(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#\$%\^&\*\(\)_\+\-=\[\]\{\};:'",.<>\/?\\|~]`).MatchString(password)
	return hasLower && hasUpper && hasNumber && hasSpecial
}
