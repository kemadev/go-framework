// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	"fmt"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/valkey-io/valkey-go"
)

func NewClient(conf config.CacheConfig) (valkey.Client, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress:         conf.ClientAddress,
		ShuffleInit:         true,
		EnableReplicaAZInfo: true,
		SendToReplicas: func(cmd valkey.Completed) bool {
			return cmd.IsReadOnly()
		},
		ClusterOption: valkey.ClusterOption{
			ShardsRefreshInterval: conf.ShardsRefreshInterval,
		},
		Sentinel: valkey.SentinelOption{
			MasterSet: conf.SentinelMasterSet,
		},
		Username: conf.Username,
		Password: conf.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating valkey client: %w", err)
	}

	return client, nil
}
