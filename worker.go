package goeurekaclient

import (
	"fmt"
	"sync"
	"time"
)

var (
	// 并发控制
	ch = make(chan int, 10)
)

// StartBatch 批量启动
func StartBatch(cnfs []EurekaClientConfig, debug bool) error {
	for _, cnf := range cnfs {
		eureka := NewEurekaAppInstance(cnf)

		// 单机运行时清除其他旧应用
		if cnf.StandAlone {
			delteOldApp(cnf)
		}

		// 注册新的应用
		err := EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, eureka)
		if err != nil {
			if debug {
				fmt.Println("Eureka Client StartBatch error: " + err.Error())
			}
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

		go func(cf EurekaClientConfig) {
			t := time.NewTicker(time.Second * time.Duration(secs))
			for {
				keepMeAlive(cf, tims)
				<-t.C
			}
		}(cnf)
	}

	// 批量应用列表维护
	keepAppCacheBatch(cnfs, debug)

	return nil
}

// Start 启动
func Start(cnf EurekaClientConfig, debug bool) error {
	eureka := NewEurekaAppInstance(cnf)

	// 删除旧应用
	delteOldApp(cnf)

	// 注册新的应用
	err := EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, eureka)
	if err != nil {
		if debug {
			fmt.Println("Eureka Client Start error: " + err.Error())
		}
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

	go func(cf EurekaClientConfig) {
		t := time.NewTicker(time.Second * time.Duration(secs))
		for {
			keepMeAlive(cf, tims)
			<-t.C
		}
	}(cnf)

	// 启动应用列表缓存表维护
	go func(cf EurekaClientConfig) {
		t := time.NewTicker(time.Second * time.Duration(cnf.AppRefreshSecs))
		for {
			keepAppCache(cf, debug)
			<-t.C
		}
	}(cnf)

	return nil
}

// 仅维护应用列表而启动
func StartForKeeper(cnf EurekaClientConfig, debug bool) {
	// 启动应用列表缓存表维护
	go func(cf EurekaClientConfig) {
		t := time.NewTicker(time.Second * time.Duration(cnf.AppRefreshSecs))
		for {
			keepAppCache(cf, debug)
			<-t.C
		}
	}(cnf)
}

// 	批量仅维护应用列表而启动
func StartForKeeperBatch(cnfs []EurekaClientConfig, debug bool) {
	// 批量应用列表维护
	keepAppCacheBatch(cnfs, debug)
}

// delteOldApp 删除旧应用
func delteOldApp(cnf EurekaClientConfig) {
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
func keepMeAlive(cnf EurekaClientConfig, tm int64) error {
	for i := int64(0); i < tm; i++ {
		err := EurekaHeartBeat(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, cnf.Id())

		if err == nil {
			return nil
		}

		// 续命失败,可能时网络问题,重试
		time.Sleep(time.Second * 2)
	}

	// 频繁注册失败,尝试无果,重新注册
	e := NewEurekaAppInstance(cnf)
	return EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, e)
}

// keepAppCache 应用列表维护
func keepAppCache(cnf EurekaClientConfig, debug bool) {
	for _, name := range cnf.Apps {
		info, err := EurekaGetApp(cnf.EurekaServerAddress, cnf.Authorization, name)
		if err != nil {
			if debug {
				fmt.Println("Eureka Client EurekaGetApp error: " + err.Error())
			}
			continue
		}

		globalEurekaAppCache.Save(cnf.EurekaName, info.Application)
	}
}

// keepAppCacheBatch 批量应用列表维护
func keepAppCacheBatch(cnfs []EurekaClientConfig, debug bool) {
	waitGroup := sync.WaitGroup{}
	for _, cnf := range cnfs {
		ch <- 1
		waitGroup.Add(1)

		go keepAppCache(cnf, debug)

		waitGroup.Done()
		<-ch
	}
	waitGroup.Wait()
}
