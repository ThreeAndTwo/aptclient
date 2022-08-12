package client

import (
	"encoding/json"
	"fmt"
	"github.com/deng00/req"
	"strings"
)

type Net struct {
	Url    string
	Header map[string]string
	Params map[string]interface{}
	IsJson bool
}

type netType string

const (
	GetTy    netType = "get"
	PostTy   netType = "post"
	DeleteTy netType = "delete"
	PutTy    netType = "put"
)

func NewNet(url string, header map[string]string, params map[string]interface{}) *Net {
	return &Net{Url: url, Header: header, Params: params}
}

func (n *Net) Request(netType netType) (string, error) {
	reqHeader, hasJson := n.initHeader()
	reqParams := n.initParam()

	if hasJson {
		n.IsJson = hasJson
	}

	switch netType {
	case GetTy:
		return n.get(reqHeader)
	case PostTy:
		return n.post(reqHeader, reqParams)
	case DeleteTy:
		return n.delete(reqHeader)
	case PutTy:
		return n.put(reqHeader, reqParams)
	default:
		return n.get(reqHeader)
	}
}

func (n *Net) initParam() req.Param {
	reqParams := req.Param{}
	for k, v := range n.Params {
		reqParams[k] = v
	}
	return reqParams
}

func (n *Net) get(header req.Header) (string, error) {
	return checkResp(req.Get(n.Url, header))
}

func (n *Net) post(header req.Header, param req.Param) (string, error) {
	if n.IsJson {
		jsonParam, _ := json.Marshal(param)
		return checkResp(req.Post(n.Url, header, jsonParam))
	}
	return checkResp(req.Post(n.Url, header, param))
}

func checkResp(res *req.Resp, err error) (string, error) {
	if err != nil || res == nil {
		return "", fmt.Errorf("request rpc error: %s", err)
	}
	return res.String(), err
}

func (n *Net) delete(header req.Header) (string, error) {
	return checkResp(req.Delete(n.Url, header))
}

func (n *Net) put(header req.Header, param req.Param) (string, error) {
	if n.IsJson {
		jsonParam, _ := json.Marshal(param)
		return checkResp(req.Put(n.Url, header, jsonParam))
	}
	return checkResp(req.Put(n.Url, header, param))
}

func (n *Net) initHeader() (req.Header, bool) {
	authHeader := req.Header{}
	hasJson := false

	for k, v := range n.Header {
		authHeader[k] = v
		if hasJsonInHeader(k, v) {
			hasJson = true
		}
	}
	return authHeader, hasJson
}

func hasJsonInHeader(key, value string) bool {
	return strings.ToLower(key) == "content-type" && strings.Contains(strings.ToLower(value), "json")
}
