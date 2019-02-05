// Copyright (c) 2018 Cisco and/or its affiliates.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package data

import (
	"github.com/ligato/cn-infra/db/keyval"
	"github.com/ligato/cn-infra/logging"
	"github.com/ligato/cn-infra/logging/logrus"
	"github.com/ligato/cn-infra/servicelabel"
)

// VppAgentCtl is test tool for testingVPP Agent plugins. In addition to testing, the vpp-agent-ctl tool can
// be used to demonstrate the usage of VPP Agent plugins and their APIs.
type VppAgentCtl interface {
	// GetCommands returns provided command set
	GetCommands() []string

	// DB access
	DbCtl
	// Other interfaces with configuration related methods
	ACLCtl
	InterfacesCtl
	IPSecCtl
	L2Ctl
	L3Ctl
	NatCtl
	PuntCtl
	StnCtl
}

// VppAgentCtlImpl is a ctl context
type VppAgentCtlImpl struct {
	Log             logging.Logger
	commands        []string
	serviceLabel    servicelabel.Plugin
	bytesConnection keyval.BytesBroker
	broker          keyval.ProtoBroker
}

// NewVppAgentCtl creates new VppAgentCtl object with initialized fields
func NewVppAgentCtl(dbCfg string, cmdSet []string, etcd bool) (*VppAgentCtlImpl, error) {
	var err error
	ctl := initVppAgentCtl(cmdSet)

	if etcd {
		// Establish ETCD connection
		ctl.bytesConnection, ctl.broker, err = CreateEtcdClient(dbCfg, ctl.serviceLabel.GetAgentPrefix(), ctl.Log)
	} else {
		// Establish Redis connection
		ctl.bytesConnection, ctl.broker, err = CreateRedisClient(dbCfg, ctl.serviceLabel.GetAgentPrefix(), ctl.Log)
	}

	return ctl, err
}

// GetCommands returns origin al vpp-agent-ctl commands
func (ctl *VppAgentCtlImpl) GetCommands() []string {
	return ctl.commands
}

func initVppAgentCtl(cmdSet []string) *VppAgentCtlImpl {
	ctl := &VppAgentCtlImpl{
		Log:      logrus.DefaultLogger(),
		commands: cmdSet,
	}

	if err := ctl.serviceLabel.Init(); err != nil {
		ctl.Log.Panic("failed to init servicvice label plugin")
	}
	return ctl
}