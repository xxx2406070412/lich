// configure
package main

import (
	"fmt"

	"github.com/cihub/seelog"
	"github.com/coocood/freecache"
	"github.com/lxmgo/config"
)

var (
	// 集群名字
	g_clusterName string
	// 节点编号 集群内唯一
	g_id string
	// cache模型 内存设置的大小M
	g_lmemoryCacheSize int64 = 256

	g_dataPath string

	logger seelog.LoggerInterface
)

func init() {

	fmt.Println("配置初始化")

	/* 日志文件配置 */
	log, err := seelog.LoggerFromConfigAsFile("conf/seelog-main.xml")

	if err != nil {
		fmt.Println("log配置文件读取失败 启动失败")
		return
	}

	logger = log

	fmt.Println("日志配置读取完毕")

	config, err := config.NewConfig("conf/conf.ini")

	if err != nil {
		fmt.Println("配置参数读取失败 启动失败")
		return
	}
	g_dataPath = config.String("datapath")
	g_clusterName = config.String("clustername")
	g_id = config.String("id")

	lmemoryCacheSize, err := config.Int64("memorycachesize")
	g_lmemoryCacheSize = lmemoryCacheSize

	logger.Infof("g_lmemoryCacheSize =%d M", g_lmemoryCacheSize)

	// 根据model 初始化cache 部分
	var isize int = int(g_lmemoryCacheSize) * 1024 * 1024

	g_pstcache = freecache.NewCache(isize)

}
