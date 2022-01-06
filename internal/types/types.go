package types

type JsonResult struct {
	Code    int    `json:"code"`
	Message string `json:"meassage"`
	//Data    string `json:"data"`
	//Data interface{} `json:"data"`
	Data map[string]interface{} `json:"data"`
}

type ResponseBase struct {
	Code    int    `json:"code"`
	Message string `json:"meassage"`
}

type DataGetIP struct {
	IP string `json:"ip"`
}

type ResponseGetIP struct {
	ResponseBase
	Data DataGetIP `json:"data"`
}

type DataGetConfig struct {
	FileName string `json:"filename"`
	Content  string `json:"content"`
}

type ResponseGetConfig struct {
	ResponseBase
	Data DataGetConfig `json:"data"`
}

type JsonIPResult struct {
	Code    int    `json:"code"`
	Message string `json:"meassage"`
	Data    IpData `json:"data"`
}

type JsonIPData struct {
	Ip string `json:"ip"`
}

type IpData struct {
	Ip string `json:"ip"`
}

type JsonConfigResult struct {
	Code    int            `json:"code"`
	Message string         `json:"meassage"`
	Data    JsonConfigData `json:"data"`
}

type JsonConfigRequest struct {
	Data JsonConfigData `json:"data"`
}

type JsonConfigData struct {
	FileName string `json:"filename"`
	Content  string `json:"content"`
}
