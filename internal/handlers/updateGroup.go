package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *HandlerManagerImp) UpdateGroupInfo(num int64) error {
	i := int64(1)
	res, err := s.dataService.GetNeedUpdateCode(10, i, num)
	if err != nil {
		logrus.Warn(err)
		return err
	}
	for _, item := range res {
		s.updateInfo(item)
	}
	return nil
}

func (s *HandlerManagerImp) updateInfo(code string) {
	n, err := s.GetChatMembers(fmt.Sprintf("@%s", code))
	if err != nil {
		logrus.Warn(err)
		if strings.Contains(err.Error(), "chat not found") {
			s.dataService.Delete([]string{code})
		} else if strings.Contains(err.Error(), "retry after") {
			time.Sleep(time.Second * 999)
		}
		return
	}

	res, err := s.GetChatInfo(fmt.Sprintf("@%s", code))
	if err != nil {
		logrus.Warn(err)
		return
	}
	desc := ""
	if _, ok := res["description"]; ok {
		desc = res["description"].(string)
	}
	activeNum := len(res["active_usernames"].([]interface{}))
	s.dataService.Update(code, int64(res["id"].(float64)), res["title"].(string), desc, n, activeNum, "", 0)
}

// 获取群人数
func (s *HandlerManagerImp) GetChatMembers(code string) (uint32, error) {
	checkSleep()
	params := map[string]string{"chat_id": code}
	res, err := s.Bot.Raw("getChatMemberCount", params)
	if err != nil {
		return 0, err
	}
	result := map[string]interface{}{}
	if err := json.Unmarshal(res, &result); err != nil {
		return 0, err
	}
	if result["ok"].(bool) {
		return uint32(result["result"].(float64)), nil
	}
	return 0, nil
}

// 获取群信息
func (s *HandlerManagerImp) GetChatInfo(code string) (map[string]interface{}, error) {
	checkSleep()
	params := map[string]string{"chat_id": code}
	res, err := s.Bot.Raw("getChat", params)
	if err != nil {
		return nil, err
	}
	result := map[string]interface{}{}
	if err := json.Unmarshal(res, &result); err != nil {
		return nil, err
	}
	return result["result"].(map[string]interface{}), nil
}

// 检查是否需要暂停一下
func checkSleep() {
	i := (rand.Intn(10) + 1) + rand.Intn(1000)/1000
	time.Sleep(time.Second * time.Duration(i))
}
