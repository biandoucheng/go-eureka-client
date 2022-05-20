package goeurekaclient

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestAddress(t *testing.T) {
	addr := NewAddress("TEST", "http", "127.0.0.1", "8080", "/health")
	fmt.Println("-----------服务地址测试----------------")
	fmt.Println(addr.Check())
	fmt.Println(addr.HealthUrl)
	fmt.Println(addr.GetUrl())
	fmt.Println(addr.AppName)
}

func TestApp(t *testing.T) {
	app := NewApp("TEST")
	app.AddHost("http", "127.0.0.1", "8080", "/health")
	fmt.Println("-----------应用信息测试----------------")
	fmt.Println(app.Name)
	fmt.Println(app.Hosts)
	fmt.Println(app.GetAnUrl())
	fmt.Println(app.HasActiveInstance())
	fmt.Println(app.Hosts)
}

func TestEurekaConf(t *testing.T) {
	cnf := NewEurekaConf()
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090
	fmt.Println("-----------Eureka配置测试----------------")
	fmt.Println(cnf)
	fmt.Println(cnf.Id())
	fmt.Println(cnf.HostName())
}

func TestEurekaRegister(t *testing.T) {
	cnf := NewEurekaConf()
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090

	ins := NewEurekaAppInstance(cnf)
	err := EurekaRegist(cnf.EurekaServerAddress, cnf.Authorization, ins)
	fmt.Println("-----------Eureka注册测试----------------")
	fmt.Println("注册错误信息 >>", err)
}

func TestEurekaHeartBeat(t *testing.T) {
	cnf := NewEurekaConf()
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090

	fmt.Println("-----------Eureka续命测试----------------")
	fmt.Println(time.Now())
	time.Sleep(time.Second * 5)
	err := EurekaHeartBeat(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, cnf.Id())
	fmt.Println("续命错误信息 >>", err)
}

func TestEurekaApps(t *testing.T) {
	cnf := NewEurekaConf()
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090
	fmt.Println("-----------Eureka拉取应用列表测试----------------")
	resp, err := EurekaGetApp(cnf.EurekaServerAddress, cnf.Authorization, "DEFAULT-EUREKA-APP0")
	fmt.Println("拉取应用错误信息 >>", err)
	res, _ := json.Marshal(resp)
	fmt.Println("拉取应用信息 >>", string(res))
}

func TestEurekaDeleteApp(t *testing.T) {
	cnf := NewEurekaConf()
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090
	err := EurekaDelteApp(cnf.EurekaServerAddress, cnf.Authorization, cnf.AppName, cnf.Id())
	fmt.Println("-----------Eureka删除应用测试----------------")
	fmt.Println("删除应用错误信息 >>", err)
}

func TestEurekaAppsCache(t *testing.T) {
	cnf := NewEurekaConf()
	cnf.AppName = "DEFAULT-EUREKA-APP8"
	cnf.Authorization = "Basic cm9vdDpyb290"
	cnf.EurekaServerAddress = "http://127.0.0.1:8080"
	cnf.Apps = []string{"DEFAULT-EUREKA-APP0", "DEFAULT-EUREKA-APP1"}
	cnf.InstanceIp = "127.0.0.1"
	cnf.InstancePort = 8090

	err := Start(cnf)
	fmt.Println("-----------Eureka客户端服务启动测试----------------")
	fmt.Println("Eureka客户端服务启动错误信息 >>", err)

	time.Sleep(time.Second * 10)
	ul, err := GetAppUrl("DEFAULT-EUREKA-APP0")
	fmt.Println("获取应用服务地址错误信息 >>", err)
	fmt.Println("获取应用服务地址 >>", ul)
	fmt.Println("获取应用服务地址 >>", ul)
	fmt.Println("获取应用服务地址 >>", ul)
}
