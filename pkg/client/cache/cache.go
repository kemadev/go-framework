// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	"fmt"

	"github.com/kemadev/go-framework/pkg/config"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeyotel"
)

func NewClient(conf config.CacheConfig) (valkey.Client, error) {
	clientAddresses := []string{}
	for _, addr := range conf.ClientAddress {
		clientAddresses = append(clientAddresses, addr.String())
	}

	client, err := valkeyotel.NewClient(valkey.ClientOption{
		InitAddress:         clientAddresses,
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
		return nil, fmt.Errorf("error creating cache client: %w", err)
	}

	return client, nil
}
