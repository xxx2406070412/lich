// lichmsgreturn
package main

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	LICH_OK    int = 0
	URLUNMATCH int = 5000 + iota
	LOSTACT
	INSERTERROR // 插入错误
	SEARCHERROR // 查找失败
	NOTEXISTKEY //key 不存在
	PARAMLOST   //参数缺失
)

type lichmsg struct {
	// 接受参数
	param *map[string]string

	// 返回参数
	code   int
	cause  string
	result string
	lcost  int64
}

type msgcommon interface {
	dealProcess()
	getResult() *map[string]string
}

func lichmsg_exception(code int, cause string) []byte {
	resultMap := make(map[string]string)
	resultMap["code"] = strconv.Itoa(code)
	resultMap["cause"] = cause

	strbyte, _ := json.Marshal(resultMap)

	return strbyte
}

func (p *lichmsg) dealProcess() {

	p.code = LOSTACT
	p.cause = "can not deal"
}

func (p *lichmsg) getResult() *map[string]string {

	//默认以json 方式返回
	resultMap := make(map[string]string)

	resultMap["code"] = strconv.Itoa(p.code)

	if p.code == 0 {
		resultMap["result"] = p.result
	} else {
		resultMap["cause"] = p.cause
	}

	resultMap["cost"] = strconv.FormatInt(time.Now().UnixNano()/1e6-p.lcost, 10)

	return &resultMap
}
