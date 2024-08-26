package internal

const (
	limitDefVal  = 5
	offsetDefVal = 0
)

type RespError struct {
	Error string `json:"error"`
}

type RespOk struct {
	Resp interface{} `json:"response"`
}

type RespInner struct {
	Records []map[string]interface{} `json:"records,omitempty"`
	Record  map[string]interface{}   `json:"record,omitempty"`
	Tables  []string                 `json:"tables,omitempty"`
	Updated *int                     `json:"updated,omitempty"`
	Deleted *int                     `json:"deleted,omitempty"`
}
