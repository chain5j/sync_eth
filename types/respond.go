// description: sync_eth
//
// @author: xwc1125
// @date: 2020/10/05
package types

type Response struct {
	Code int         `json:"code" example:"200"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type Page struct {
	Details interface{} `json:"details"`
	Total   int         `json:"total"`
	Page    int         `json:"page"`
	Limit   int         `json:"limit"`
}

type PageResponse struct {
	Code int    `json:"code" example:"200"`
	Data Page   `json:"data"`
	Msg  string `json:"msg"`
}

func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	return res
}

func (res *PageResponse) ReturnOK() *PageResponse {
	res.Code = 200
	return res
}
