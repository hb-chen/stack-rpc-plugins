/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package component

import (
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/protocol/http"
	"testing"

	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/cluster/roundrobin"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env/config"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env/config/json"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/extension"
	. "github.com/tevid/gohamcrest"

	json2 "encoding/json"
)

func init() {
	extension.SetLoadBalance(&roundrobin.RoundRobin{})
}

const servicesConfigResponseStr = `[{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.15.128.102:apollo-configservice:8080",
"homepageUrl": "http://10.15.128.102:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.15.88.125:apollo-configservice:8080",
"homepageUrl": "http://10.15.88.125:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.14.0.11:apollo-configservice:8080",
"homepageUrl": "http://10.14.0.11:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.14.0.193:apollo-configservice:8080",
"homepageUrl": "http://10.14.0.193:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.15.128.101:apollo-configservice:8080",
"homepageUrl": "http://10.15.128.101:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.14.0.192:apollo-configservice:8080",
"homepageUrl": "http://10.14.0.192:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.15.88.124:apollo-configservice:8080",
"homepageUrl": "http://10.15.88.124:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.15.128.103:apollo-configservice:8080",
"homepageUrl": "http://10.15.128.103:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "localhost:apollo-configservice:8080",
"homepageUrl": "http://10.14.0.12:8080/"
},
{
"appName": "APOLLO-CONFIGSERVICE",
"instanceId": "10.14.0.194:apollo-configservice:8080",
"homepageUrl": "http://10.14.0.194:8080/"
}
]`

var (
	jsonConfigFile = &json.ConfigFile{}
)

func TestSelectOnlyOneHost(t *testing.T) {
	appConfig := env.InitFileConfig()
	trySyncServerIPList(appConfig)
	host := "http://localhost:8888/"
	Assert(t, host, Equal(appConfig.GetHost()))
	load := extension.GetLoadBalance().Load(*appConfig.GetServers())
	Assert(t, load, NotNilVal())
	Assert(t, host, NotEqual(load.HomepageURL))

	appConfig.IP = host
	Assert(t, host, Equal(appConfig.GetHost()))
	load = extension.GetLoadBalance().Load(*appConfig.GetServers())
	Assert(t, load, NotNilVal())
	Assert(t, host, NotEqual(load.HomepageURL))

	appConfig.IP = "https://localhost:8888"
	https := "https://localhost:8888/"
	Assert(t, https, Equal(appConfig.GetHost()))
	load = extension.GetLoadBalance().Load(*appConfig.GetServers())
	Assert(t, load, NotNilVal())
	Assert(t, host, NotEqual(load.HomepageURL))
}

type testComponent struct {
}

//Start 启动同步服务器列表
func (s *testComponent) Start() {
}

func TestStartRefreshConfig(t *testing.T) {
	StartRefreshConfig(&testComponent{})
}

func TestName(t *testing.T) {

}

func trySyncServerIPList(appConfig *config.AppConfig) {
	SyncServerIPListSuccessCallBack([]byte(servicesConfigResponseStr), http.CallBack{AppConfig: appConfig})
}

//SyncServerIPListSuccessCallBack 同步服务器列表成功后的回调
func SyncServerIPListSuccessCallBack(responseBody []byte, callback http.CallBack) (o interface{}, err error) {
	log.Debug("get all server info:", string(responseBody))

	tmpServerInfo := make([]*config.ServerInfo, 0)

	err = json2.Unmarshal(responseBody, &tmpServerInfo)

	if err != nil {
		log.Error("Unmarshal json Fail,Error:", err)
		return
	}

	if len(tmpServerInfo) == 0 {
		log.Info("get no real server!")
		return
	}

	for _, server := range tmpServerInfo {
		if server == nil {
			continue
		}
		callback.AppConfig.GetServers().Store(server.HomepageURL, server)
	}
	return
}
