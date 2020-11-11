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

package remote

import (
	log "github.com/stack-labs/stack-rpc/logger"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/env/config"
	"github.com/stack-labs/stack-rpc/plugins/config/source/apollo/agollo/protocol/http"
)

// AbsApolloConfig 抽象 apollo 配置
type AbsApolloConfig struct {
	remoteApollo ApolloConfig
}

func (a *AbsApolloConfig) SyncWithNamespace(namespace string, appConfig *config.AppConfig) *config.ApolloConfig {
	if appConfig == nil {
		panic("can not find apollo config!please confirm!")
	}

	urlSuffix := a.remoteApollo.GetSyncURI(*appConfig, namespace)

	callback := a.remoteApollo.CallBack(namespace)
	apolloConfig, err := http.RequestRecovery(appConfig, &env.ConnectConfig{
		URI:    urlSuffix,
		AppID:  appConfig.AppID,
		Secret: appConfig.Secret,
	}, &callback)
	if err != nil {
		log.Errorf("request %s fail, error:%v", urlSuffix, err)
		return nil
	}

	if apolloConfig == nil {
		log.Warn("apolloConfig is nil")
		return nil
	}

	return apolloConfig.(*config.ApolloConfig)
}
