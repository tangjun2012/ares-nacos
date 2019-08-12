package naming

import (
	"encoding/json"
	"errors"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/swift9/ares-nacos/config"
	"log"
	"os"
)

var namingClient naming_client.INamingClient

func init() {
	serverConfigList := arraylist.New()
	nacosServerConfigs := config.GetLocalConfig().Get("nacos.serverConfigs").Array()
	for _, serverConfig := range nacosServerConfigs {
		serverConfigList.Add(constant.ServerConfig{
			IpAddr:      serverConfig.Get("ipAddr").String(),
			ContextPath: serverConfig.Get("contextPath").String(),
			Port:        serverConfig.Get("port").Uint(),
		})
	}

	var serverConfigs []constant.ServerConfig
	serverConfigsJson, _ := serverConfigList.ToJSON()
	json.Unmarshal(serverConfigsJson, &serverConfigs)
	properties := map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig": constant.ClientConfig{
			TimeoutMs:            10 * 1000, //http请求超时时间，单位毫秒
			ListenInterval:       15 * 1000, //监听间隔时间，单位毫秒（仅在ConfigClient中有效）
			BeatInterval:         15 * 1000, //心跳间隔时间，单位毫秒（仅在ServiceClient中有效）
			NamespaceId:          config.GetLocalConfig().Get("nacos.namespaceId").String(),
			UpdateThreadNum:      20,   //更新服务的线程数
			NotLoadCacheAtStart:  true, //在启动时不读取本地缓存数据，true--不读取，false--读取
			UpdateCacheWhenEmpty: true, //当服务列表为空时是否更新本地缓存，true--更新,false--不更新
		},
	}

	namingClient, _ = clients.CreateNamingClient(properties)

	if namingClient == nil {
		log.Println("load namingClient error")
		os.Exit(1)
	}

	log.Println("create nacos naming client finished")
}

/* 注册服务 */
func RegisterService(ip string, port uint64, serviceName string, clusterName string, metadata map[string]string) (bool, error) {
	rst, err := checkArgs0(ip, port, serviceName, clusterName)
	if !rst {
		return rst, err
	}
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		ClusterName: clusterName,
		Metadata:    metadata,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
	}
	return namingClient.RegisterInstance(param)
}

/* 注销服务 */
func LogoutService(ip string, port uint64, serviceName string, clusterName string) (bool, error) {
	rst, err := checkArgs0(ip, port, serviceName, clusterName)
	if !rst {
		return rst, err
	}
	param := vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: serviceName,
		Cluster:     clusterName,
		Ephemeral:   false,
	}
	return namingClient.DeregisterInstance(param)
}

/* 监听服务 */
func AddListener(serviceName string, clusters []string, callback func(services []model.SubscribeService, err error)) error {
	rst, err := checkArgs1(serviceName, clusters, callback)
	if !rst {
		return err
	}
	param := vo.SubscribeParam{
		ServiceName:       serviceName,
		Clusters:          clusters,
		SubscribeCallback: callback,
	}
	return namingClient.Subscribe(&param)
}

/* 获取注册服务列表 */
func SelectServices(serviceName string, clusters []string, healthy bool) ([]model.Instance, error) {
	rst, err := checkArgs2(serviceName, clusters)
	if !rst {
		return nil, err
	}
	param := vo.SelectInstancesParam{
		Clusters:    clusters,
		ServiceName: serviceName,
		HealthyOnly: healthy,
	}
	return namingClient.SelectInstances(param)
}

/* 获取所有注册服务 */
func SelectAllServices(serviceName string, clusters []string) ([]model.Instance, error) {
	rst, err := checkArgs2(serviceName, clusters)
	if !rst {
		return nil, err
	}
	param := vo.SelectAllInstancesParam{
		Clusters:    clusters,
		ServiceName: serviceName,
	}
	return namingClient.SelectAllInstances(param)
}

// 参数校验
func checkArgs0(ip string, port uint64, serviceName string, clusterName string) (bool, error) {
	if len(ip) <= 0 {
		return false, errors.New("ip not be empty")
	}
	if port == 0 {
		return false, errors.New("port is illegal")
	}
	if len(serviceName) <= 0 {
		return false, errors.New("serviceName not be empty")
	}
	return true, nil
}
func checkArgs1(serviceName string, clusters []string, callback func(services []model.SubscribeService, err error)) (bool, error) {
	if len(serviceName) <= 0 {
		return false, errors.New("serviceName not be empty")
	}
	if len(clusters) <= 0 {
		return false, errors.New("clusters not be empty")
	}
	if callback == nil {
		return false, errors.New("callback not be empty")
	}
	return true, nil
}
func checkArgs2(serviceName string, clusters []string) (bool, error) {
	if len(serviceName) <= 0 {
		return false, errors.New("serviceName not be empty")
	}
	if len(clusters) <= 0 {
		return false, errors.New("clusters not be empty")
	}
	return true, nil
}
