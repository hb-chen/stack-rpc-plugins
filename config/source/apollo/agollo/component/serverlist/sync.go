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
	"encoding/json"
	"time"

	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/component"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env/config"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/protocol/http"
)

const (
	//refresh ip list
	refreshIPListInterval = 20 * time.Minute //20m
)

func init() {

}

//InitSyncServerIPList 初始化同步服务器信息列表
func InitSyncServerIPList(appConfig *config.AppConfig) {
	go component.StartRefreshConfig(&SyncServerIPListComponent{appConfig})
}

//SyncServerIPListComponent set timer for update ip list
//interval : 20m
type SyncServerIPListComponent struct {
	appConfig *config.AppConfig
}

//Start 启动同步服务器列表
func (s *SyncServerIPListComponent) Start() {
	SyncServerIPList(s.appConfig)
	log.Debug("syncServerIpList started")

	t2 := time.NewTimer(refreshIPListInterval)
	for {
		select {
		case <-t2.C:
			SyncServerIPList(s.appConfig)
			t2.Reset(refreshIPListInterval)
		}
	}
}

//SyncServerIPList sync ip list from server
//then
//1.update agcache
//2.store in disk
func SyncServerIPList(appConfig *config.AppConfig) error {
	if appConfig == nil {
		panic("can not find apollo config!please confirm!")
	}

	_, err := http.Request(appConfig.GetServicesConfigURL(), &env.ConnectConfig{
		AppID:  appConfig.AppID,
		Secret: appConfig.Secret,
	}, &http.CallBack{
		SuccessCallBack: SyncServerIPListSuccessCallBack,
		AppConfig:       appConfig,
	})

	return err
}

//SyncServerIPListSuccessCallBack 同步服务器列表成功后的回调
func SyncServerIPListSuccessCallBack(responseBody []byte, callback http.CallBack) (o interface{}, err error) {
	log.Debug("get all server info:", string(responseBody))

	tmpServerInfo := make([]*config.ServerInfo, 0)

	err = json.Unmarshal(responseBody, &tmpServerInfo)

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
