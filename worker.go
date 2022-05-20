package goeurekaclient

import (
	"fmt"
	"time"
)

// Start 启动
func Start(cnf *EurekaClientConfig) error {
	if cnf == nil {
		cnf = &DefaultEurekaClientConf
	}

	eureka := NewEurekaAppInstance(cnf)

	// 删除旧应用
	delteOldApp(cnf)

	// 注册新的应用
	err := EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, eureka)
	if err != nil {
		return err
	}

	// 启动心跳续约
	secs := cnf.RenewalIntervalInSecs
	if secs > 10 {
		secs -= 5
	}

	// 计算心跳续约失败后的重试次数
	tims := cnf.RenewalIntervalInSecs / 2
	if tims <= 0 {
		tims = 1
	}

	go func() {
		t := time.NewTicker(time.Second * time.Duration(secs))
		for {
			keepMeAlive(cnf, tims)
			<-t.C
		}
	}()

	// 启动应用列表缓存表维护
	go func() {
		t := time.NewTicker(time.Second * time.Duration(cnf.AppRefreshSecs))
		for {
			keepAppCache(cnf)
			<-t.C
		}
	}()

	return nil
}

// delteOldApp 删除旧应用
func delteOldApp(cnf *EurekaClientConfig) {
	info, err := EurekaGetApp(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName)
	if err != nil {
		return
	}

	app := info.Application

	for _, ins := range app.Instance {
		EurekaDelteApp(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, ins.InstanceId)
	}
}

// KeepMeAlive 本服务保活
func keepMeAlive(cnf *EurekaClientConfig, tm int64) error {
	err := EurekaHeartBeat(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, cnf.Id())

	if err == nil {
		return nil
	}

	// 续命失败,可能时网络问题,重试
	if tm > 0 {
		tm -= 1
		time.Sleep(time.Second * 2)
		return keepMeAlive(cnf, tm)
	}

	// 频繁注册失败,尝试无果,重新注册
	e := NewEurekaAppInstance(cnf)
	return EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, e)
}

// keepAppCache 应用列表维护
func keepAppCache(cnf *EurekaClientConfig) {
	for _, name := range cnf.Apps {
		info, err := EurekaGetApp(cnf.EurekaServerAddress, cnf.Authorization, name)
		if err != nil {
			fmt.Println("拉取应用信息错误 >>>", err)
			continue
		}

		globalEurekaAppCache.Save(info.Application)
	}
}
