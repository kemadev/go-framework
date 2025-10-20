// Copyright 2025 kemadev
// SPDX-License-Identifier: MPL-2.0

package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/ristretto/v2"
	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/monitoring"
	"github.com/valkey-io/valkey-go"
	"github.com/valkey-io/valkey-go/valkeyotel"
)

type ristrettoCache[R any] struct {
	cache *ristretto.Cache[string, R]
}

func (r *ristrettoCache[V]) Get(key string) (V, bool) {
	result, found := r.cache.Get(key)
	if found {
		return result, true
	}
	return *(new(V)), false
}

func (r *ristrettoCache[V]) Set(key string, value V) {
	r.cache.Set(key, value, 0)
}

func NewFailsafeLocal[V any](config ristretto.Config[string, V]) (*ristrettoCache[V], error) {
	c, err := NewLocal(config)
	if err != nil {
		return nil, fmt.Errorf("error creating ristretto cache: %w", err)
	}

	return &ristrettoCache[V]{
		cache: c,
	}, nil
}

func NewLocal[K ristretto.Key, V any](config ristretto.Config[K, V]) (*ristretto.Cache[K, V], error) {
	cache, err := ristretto.NewCache(&config)
	if err != nil {
		return nil, fmt.Errorf("error creating cache: %w", err)
	}

	return cache, nil
}

type valkeyCache[R any] struct {
	client valkey.Client
}

func (v *valkeyCache[R]) Get(key string) (R, bool) {
	ctx := context.Background()

	result, err := v.client.Do(ctx, v.client.B().Get().Key(key).Build()).AsBytes()
	if err != nil {
		return *(new(R)), false
	}

	var typedResult R
	buf := bytes.NewBuffer(result)

	err = gob.NewDecoder(buf).Decode(&typedResult)
	if err != nil {
		return *(new(R)), false
	}

	return typedResult, true
}

func (v *valkeyCache[R]) Set(key string, value R) {
	ctx := context.Background()
	var buf bytes.Buffer

	if err := gob.NewEncoder(&buf).Encode(value); err != nil {
		return
	}

	v.client.Do(ctx, v.client.B().Set().Key(key).Value(buf.String()).Build())
}

func NewFailsafeShared[V any](client valkey.Client) *valkeyCache[V] {
	return &valkeyCache[V]{
		client: client,
	}
}

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
		// Let user define a failsafe strategy when calling
		DisableRetry: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating cache client: %w", err)
	}

	return client, nil
}

func Check(c valkey.Client, statusOnPingFail monitoring.Status) monitoring.StatusCheck {
	err := c.Do(context.Background(), c.B().Ping().Build()).Error()
	if err != nil {
		return monitoring.StatusCheck{
			Status:  statusOnPingFail,
			Message: "ping failed",
		}
	}

	return monitoring.StatusCheck{
		Status:  monitoring.StatusOK,
		Message: monitoring.StatusOK.String(),
	}
}
