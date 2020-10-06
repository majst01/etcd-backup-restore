// Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package restorer

import (
	"fmt"
	"path"

	"github.com/coreos/etcd/pkg/types"
	flag "github.com/spf13/pflag"
)

// NewRestorationConfig returns the restoration config.
func NewRestorationConfig() *RestorationConfig {
	return &RestorationConfig{
		InitialCluster:           initialClusterFromName(defaultName),
		InitialClusterToken:      defaultInitialClusterToken,
		RestoreDataDir:           fmt.Sprintf("%s.etcd", defaultName),
		InitialAdvertisePeerURLs: []string{defaultInitialAdvertisePeerURLs},
		Name:                     defaultName,
		SkipHashCheck:            false,
		MaxFetchers:              defaultMaxFetchers,
		EmbeddedEtcdQuotaBytes:   int64(defaultEmbeddedEtcdQuotaBytes),
	}
}

// AddFlags adds the flags to flagset.
func (c *RestorationConfig) AddFlags(fs *flag.FlagSet) {
	fs.StringVar(&c.InitialCluster, "initial-cluster", c.InitialCluster, "initial cluster configuration for restore bootstrap")
	fs.StringVar(&c.InitialClusterToken, "initial-cluster-token", c.InitialClusterToken, "initial cluster token for the etcd cluster during restore bootstrap")
	fs.StringVarP(&c.RestoreDataDir, "data-dir", "d", c.RestoreDataDir, "path to the data directory")
	fs.StringArrayVar(&c.InitialAdvertisePeerURLs, "initial-advertise-peer-urls", c.InitialAdvertisePeerURLs, "list of this member's peer URLs to advertise to the rest of the cluster")
	fs.StringVar(&c.Name, "name", c.Name, "human-readable name for this member")
	fs.BoolVar(&c.SkipHashCheck, "skip-hash-check", c.SkipHashCheck, "ignore snapshot integrity hash value (required if copied from data directory)")
	fs.UintVar(&c.MaxFetchers, "max-fetchers", c.MaxFetchers, "maximum number of threads that will fetch delta snapshots in parallel")
	fs.Int64Var(&c.EmbeddedEtcdQuotaBytes, "embedded-etcd-quota-bytes", c.EmbeddedEtcdQuotaBytes, "maximum backend quota for the embedded etcd used for applying delta snapshots")
	fs.StringVar(&c.CompressionMethod, "uncompression-method", c.CompressionMethod, "enable compression and define compression method, must be one of auto|none|tar|targz|tarlz4")
}

// Validate validates the config.
func (c *RestorationConfig) Validate() error {
	if _, err := types.NewURLsMap(c.InitialCluster); err != nil {
		return fmt.Errorf("failed creating url map for restore cluster: %v", err)
	}
	if _, err := types.NewURLs(c.InitialAdvertisePeerURLs); err != nil {
		return fmt.Errorf("failed parsing peers urls for restore cluster: %v", err)
	}
	if c.MaxFetchers <= 0 {
		return fmt.Errorf("max fetchers should be greater than zero")
	}
	if c.EmbeddedEtcdQuotaBytes <= 0 {
		return fmt.Errorf("Etcd Quota size for etcd must be greater than 0")
	}
	switch c.CompressionMethod {
	case "", "auto", "none", "tar", "targz", "tarlz4":
	default:
		return fmt.Errorf("Invalid compression-method:%q", c.CompressionMethod)
	}
	c.RestoreDataDir = path.Clean(c.RestoreDataDir)
	return nil
}

func initialClusterFromName(name string) string {
	n := name
	if name == "" {
		n = defaultName
	}
	return fmt.Sprintf("%s=http://localhost:2380", n)
}
