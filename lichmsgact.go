// lichmsgact
package main

import (
	"strconv"

	"github.com/coocood/freecache"
)

/* 这个是对外的主要消息接口 */
type msgact struct {
	lichmsg
	iscacheHit bool //如果是查询 并且是混合模式情况下 是否击中cache
}

func (p *msgact) dealProcess() {

	var paramMap map[string]string = *p.param

	act, _ := paramMap["act"]
	prefix, _ := paramMap["prefix"]
	key, _ := paramMap["key"]
	value, _ := paramMap["value"]
	option, _ := paramMap["option"]
	expiresec, _ := paramMap["expiresec"]

	var iexpiresec int = 60 * 30 //默认半小时

	if "" != expiresec {
		temp, err := strconv.Atoi(expiresec)
		if nil == err {
			if temp > 0 {
				iexpiresec = temp
			}
		}
	}

	switch act {
	case "1001":
		if "" == prefix || "" == key || "" == value || "" == option {
			p.code = PARAMLOST
			p.cause = "request param error please check prefix key value option"
			return
		}

		err := actInsert(prefix, key, &value, iexpiresec, option)

		if nil != err {
			p.code = INSERTERROR
			p.cause = "insert fail"
		}

	case "1002":
		if "" == prefix || "" == key || "" == option {
			p.code = PARAMLOST
			p.cause = "request param error please check prefix key option"
			return
		}

		value, err := actSearch(prefix, key, option)
		if nil != err {

			if freecache.ErrNotFound == err {
				p.code = NOTEXISTKEY
				p.cause = "no exist key"
			} else {
				p.code = SEARCHERROR
				p.cause = "search fail"
			}

		} else {
			p.result = string(value)
		}
	case "1003":
		if "" == prefix || "" == key || "" == option {
			p.code = PARAMLOST
			p.cause = "request param error please check prefix key option"
			return
		}
		actDelete(prefix, key, option)
	default:
		//不会走入此分支
		return
	}

	return
}

func (p *msgact) getResult() *map[string]string {
	/* 采用 标准返回 除非特殊接口 例如 返回html 需要重写此方法 */
	return p.lichmsg.getResult()
}

// option 1 代表仅仅插入内存 memcache类似
//        2 代表redis 落地缓存
//        3 代表直接落地文件 就是mysql类似
func actInsert(prefix string, key string, ps_value *string, iexpireSeconds int, option string) error {

	if option == "1" {
		return lich_cache_insert([]byte(prefix+key), []byte(*ps_value), iexpireSeconds)
	} else if option == "3" {
		return lich_bolt_insert([]byte(prefix), []byte(key), []byte(*ps_value))
	} else {
		err := lich_bolt_insert([]byte(prefix), []byte(key), []byte(*ps_value))
		if nil != err {
			return err
		}
		err = lich_cache_insert([]byte(prefix+key), []byte(*ps_value), 60*60)
		if nil != err {
			lich_bolt_delete([]byte(prefix), []byte(key))
			return err
		}
		return nil
	}
}

func actDelete(bulkName string, key string, option string) error {

	if option == "1" {
		lich_cache_delete([]byte(bulkName + key))
	} else if option == "3" {
		lich_bolt_delete([]byte(bulkName), []byte(key))
	} else {
		lich_bolt_delete([]byte(bulkName), []byte(key))
		lich_cache_delete([]byte(bulkName + key))
	}

	return nil
}

func actSearch(prefix string, key string, option string) ([]byte, error) {

	if option == "1" {
		return lich_cache_get([]byte(prefix + key))
	} else if option == "3" {
		return lich_bolt_get([]byte(prefix), []byte(key))
	} else {
		result, err := lich_cache_get([]byte(prefix + key))
		if nil == err {
			// 缓存找到了 直接返回
			return result, err
		}

		result, err = lich_bolt_get([]byte(prefix), []byte(key))
		if nil != err {
			// 硬盘中也没找到
			return result, err
		}

		// 将硬盘中的插入缓存
		lich_bolt_insert([]byte(prefix), []byte(key), result)
		return result, err
	}
}
