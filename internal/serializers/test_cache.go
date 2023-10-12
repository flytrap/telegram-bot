package serializers

import "testing"

func TestParseLocation(t *testing.T) {
	dl := DataCacheLocation{}
	dl.ParseLocation("  四川省 - 成都市 -  其他地区 ")
	if len(dl.Tags) != 4 {
		t.Error(dl.Tags)
	}
}
