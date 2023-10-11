package human

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

// 结构体转为Map, 按json结构转换
func Decode(data interface{}) (map[string]interface{}, error) {
	info := map[string]interface{}{}
	res, err := json.Marshal(data)
	if err != nil {
		logrus.Warning(err)
		return nil, err
	}
	err = json.Unmarshal(res, &info)
	if err != nil {
		logrus.Warning(err)
		return nil, err
	}
	return info, nil
}

// map转为结构体, 按json结构转换
func Encode(data map[string]interface{}, result interface{}) error {
	res, err := json.Marshal(data)
	if err != nil {
		logrus.Warning(err)
		return err
	}
	err = json.Unmarshal(res, result)
	if err != nil {
		logrus.Warning(err)
	}
	return err
}
