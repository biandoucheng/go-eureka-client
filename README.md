# go-eureka-client
Eureka客户端

## 对象信息
- AddressObject 应用实例的主机信息
- AppObject 应用信息(应用地址维护缓存里面的,是Eureka客户端拉取到的内容的简化)
- EurekaAppCache 应用列表信息缓存,缓存了每个需要的服务的名称及可用主机信息
- EurekaClientConfig Eureka客户端配置,包含所有去请求服务端接口所需的信息
- AppResponse Eureka拉取应用信息时的响应结构
- EurekaAppInfo Eureka拉取应用信息时,其内部的应用信息结构
- EurekaAppInstance 一个完整的Eureka应用实例结构
- EurekaRegisterRequest 注册Eureka应用时的请求结构
## Eureka方法实现
- EurekaRegist 注册应用
- EurekaHeartBeat 心跳续约
- EurekaGetApp 应用拉取
- EurekaGetAppAll 拉取全部应用
- EurekaDeleteApp 应用删除
- GetAllUrls 获取应用的全部服务地址
- SetHTTPTransport 设置Eureka客户端使用的HTTP传输器；传入nil恢复默认传输器

### 自定义HTTP传输器

客户端默认使用保留原有连接池、拨号和超时配置的HTTP传输器，并默认校验TLS证书。如果需要配置HTTPS证书、代理、连接池或其他传输参数，可以在启动Eureka客户端前传入自定义传输器：

```go
transport := &http.Transport{
    TLSClientConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
    },
}
goeurekaclient.SetHTTPTransport(transport)
```

`SetHTTPTransport` 接受 `http.RoundTripper`，传入 `nil` 时恢复为库默认传输器。
## 逻辑流程
- 初始化 EurekaClientConfig 配置信息
- - EurekaServerAddress Eureka服务端地址
- - Authorization Eureka服务端 Http Auth 授权头 如: Basic ZHNkYXM6ZHNkc2FzZGE=
- - AppName 本应用名称
- - InstanceDomain 本应用服务域名,置空将被填充为InstanceIp
- - InstanceIp 实例IP,如果使用外网IP需要填写,否则被填充为内网IP
- - InstancePort 本应用服务端口
- - InstanceHealthCheckUrl 本应用健康检查接口
- - RenewalIntervalInSecs 本应用续命周期 秒
- - DurationInSecs 本应用续命失约后的保留时长 秒
- - AppRefreshSecs 所需应用列表的缓存刷新周期 秒
- 删除删除一次的注册记录
- - EurekaGetApp  获取上次注册的实例信息
- - EurekaDeleteApp 删除上次的注册记录
- 注册本次服务
- - EurekaRegist
- 后台心跳续约
- - EurekaHeartBeat
- 后台缓存所需应用信息
- - EurekaAppCache
- 获取目标应用访问地址
- - GetAppUrl
## 启动方案
### 使用启动方案会自动去注册,续约,应用列表维护 你只需要配置好 EurekaClientConfig 传入启动方法 Start() 即可
- 配置 EurekaClientConfig
- 启动单个Eureka服务 Start
- 启动多个Eureka服务 StartBatch
