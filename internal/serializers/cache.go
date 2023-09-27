package serializers

type DataCache struct {
	Category string  `json:"category"`
	Language string  `json:"language"`
	Name     string  `json:"name"`
	Code     string  `json:"code" `
	Type     int8    `json:"type"`
	Number   uint32  `json:"number"`
	Desc     string  `json:"desc"`
	Weight   float32 `json:"weight"`
}
