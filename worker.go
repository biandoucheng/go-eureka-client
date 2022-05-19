package goeurekaclient

import "time"

// Start 启动
func Start(cnf *EurekaClientConfig) error {
	if cnf == nil {
		cnf = &DefaultEurekaClientConf
	}

	eureka := NewEurekaAppInstance(cnf)

	// 删除旧应用
	delteOldApp(cnf)

	// 注册新的应用
	err := EurekaRegist(cnf.EurekaServerAddress, &eureka)
	if err != nil {
		return err
	}

	// 启动心跳续约
	secs := cnf.RenewalIntervalInSecs
	if secs > 10 {
		secs -= 5
	}

	go func() {
		t := time.NewTicker(time.Second * time.Duration(secs))
		for {
			KeepMeAlive(cnf, 3)
			<-t.C
		}
	}()

	// 启动应用列表缓存表维护
	go func() {
		t := time.NewTicker(time.Second * time.Duration(cnf.AppRefreshSecs))
		for {
			KeepAppCache(cnf)
			<-t.C
		}
	}()

	return nil
}

// delteOldApp 删除旧应用
func delteOldApp(cnf *EurekaClientConfig) {
	info, err := EurekaGetApp(cnf.EurekaServerAddress, cnf.AppName)
	if err != nil {
		return
	}

	app := info.Application

	for _, ins := range app.Instance {
		EurekaDelteApp(cnf.EurekaServerAddress, cnf.AppName, ins.InstanceId)
	}
}

// KeepMeAlive 本服务保活
func KeepMeAlive(cnf *EurekaClientConfig, tm int64) error {
	err := EurekaHeartBeat(cnf.EurekaServerAddress, cnf.AppName, cnf.Id())

	if err == nil {
		return nil
	}

	if tm > 0 {
		tm -= 1
		time.Sleep(time.Second * 3)
		return KeepMeAlive(cnf, tm)
	}

	return err
}

// KeepAppCache 应用列表维护
func KeepAppCache(cnf *EurekaClientConfig) {
	for _, name := range cnf.Apps {
		info, err := EurekaGetApp(cnf.EurekaServerAddress, name)
		if err != nil {
			continue
		}

		GlobalEurekaAppCache.Save(info.Application)
	}
}
