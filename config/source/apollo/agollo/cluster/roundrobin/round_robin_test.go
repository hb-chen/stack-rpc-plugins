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

package roundrobin

import (
	"github.com/stack-labs/stack-rpc-plugins/config/source/apollo/agollo/component/serverlist"
	"github.com/stack-labs/stack-rpc-plugins/config/source/apollo/agollo/env/config"
	"github.com/stack-labs/stack-rpc-plugins/config/source/apollo/agollo/protocol/http"
	"testing"

	"github.com/stack-labs/stack-rpc-plugins/config/source/apollo/agollo/env"
	. "github.com/tevid/gohamcrest"
)

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

func TestSelectHost(t *testing.T) {
	balanace := &RoundRobin{}

	appConfig := env.InitFileConfig()
	//mock ip data
	trySyncServerIPList(appConfig)

	servers := appConfig.GetServers()
	t.Log("appconfig host:" + appConfig.GetHost())
	t.Log("appconfig select host:", balanace.Load(*appConfig.GetServers()).HomepageURL)

	host := "http://localhost:8888/"
	Assert(t, host, Equal(appConfig.GetHost()))
	Assert(t, host, NotEqual(balanace.Load(*appConfig.GetServers()).HomepageURL))

	//check select next time
	appConfig.SetNextTryConnTime(5)
	Assert(t, host, NotEqual(balanace.Load(*appConfig.GetServers()).HomepageURL))

	//check servers
	appConfig.SetNextTryConnTime(5)
	firstHost := balanace.Load(*appConfig.GetServers())
	Assert(t, host, NotEqual(firstHost.HomepageURL))
	appConfig.SetDownNode(firstHost.HomepageURL)

	secondHost := balanace.Load(*appConfig.GetServers()).HomepageURL
	Assert(t, host, NotEqual(secondHost))
	Assert(t, firstHost, NotEqual(secondHost))
	appConfig.SetDownNode(secondHost)

	thirdHost := balanace.Load(*appConfig.GetServers()).HomepageURL
	Assert(t, host, NotEqual(thirdHost))
	Assert(t, firstHost, NotEqual(thirdHost))
	Assert(t, secondHost, NotEqual(thirdHost))

	servers.Range(func(k, v interface{}) bool {
		appConfig.SetDownNode(k.(string))
		return true
	})

	Assert(t, balanace.Load(*appConfig.GetServers()), NilVal())

	//no servers
	//servers = make(map[string]*serverInfo, 0)
	deleteServers(appConfig)
	Assert(t, balanace.Load(*appConfig.GetServers()), NilVal())
}

func deleteServers(appConfig *config.AppConfig) {
	servers := appConfig.GetServers()
	servers.Range(func(k, v interface{}) bool {
		servers.Delete(k)
		return true
	})
}

func trySyncServerIPList(appConfig *config.AppConfig) {
	serverlist.SyncServerIPListSuccessCallBack([]byte(servicesConfigResponseStr), http.CallBack{AppConfig: appConfig})
}
