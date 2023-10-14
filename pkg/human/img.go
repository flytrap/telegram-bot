package human

import (
	"encoding/base64"
	"errors"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
)

// base64转io对象
func Base64ToIoReader(text string) (io.Reader, error) {
	li := strings.Split(text, ",")
	if len(li) != 2 {
		return nil, errors.New(text)
	}
	imgStr := strings.TrimSpace(li[len(li)-1])
	if len(imgStr) == 0 {
		return nil, errors.New(text)
	}
	imgData, err := base64.StdEncoding.DecodeString(imgStr)
	if err != nil {
		logrus.Warning(err)
		return nil, err
	}
	return strings.NewReader(string(imgData)), nil
}
