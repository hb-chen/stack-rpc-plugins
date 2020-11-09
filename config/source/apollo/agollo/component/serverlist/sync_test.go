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

package serverlist

import (
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/protocol/http"
	"testing"

	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env/config"
	. "github.com/tevid/gohamcrest"
)

func TestSyncServerIPList(t *testing.T) {
	trySyncServerIPList(t)
}

func trySyncServerIPList(t *testing.T) {
	server := runMockServicesConfigServer()
	defer server.Close()

	newAppConfig := getTestAppConfig()
	newAppConfig.IP = server.URL
	err := SyncServerIPList(newAppConfig)

	Assert(t, err, NilVal())

	serverLen := 0
	newAppConfig.GetServers().Range(func(k, v interface{}) bool {
		serverLen++
		return true
	})

	Assert(t, 10, Equal(serverLen))

}

func getTestAppConfig() *config.AppConfig {
	jsonStr := `{
    "appId": "test",
    "cluster": "dev",
    "namespaceName": "application",
    "ip": "localhost:8888",
    "releaseKey": "1"
	}`
	c, _ := env.Unmarshal([]byte(jsonStr))

	return c.(*config.AppConfig)
}

func TestSyncServerIpListSuccessCallBack(t *testing.T) {
	appConfig := getTestAppConfig()
	SyncServerIPListSuccessCallBack([]byte(servicesConfigResponseStr), http.CallBack{AppConfig: appConfig})
	Assert(t, appConfig.GetServersLen(), Equal(10))
}

func TestSetDownNode(t *testing.T) {
	t.SkipNow()
	appConfig := getTestAppConfig()
	SyncServerIPListSuccessCallBack([]byte(servicesConfigResponseStr), http.CallBack{AppConfig: appConfig})

	downNode := "10.15.128.102:8080"
	appConfig.SetDownNode(downNode)

	value, ok := appConfig.GetServers().Load("http://10.15.128.102:8080/")
	info := value.(*config.ServerInfo)
	Assert(t, ok, Equal(true))
	Assert(t, info.IsDown, Equal(true))
}
