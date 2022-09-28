// Copyright 2022 The Bucketeer Authors.
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

package kafka

import (
	"context"

	"github.com/Shopify/sarama"
)

type ClusterAdmin struct {
	url    string
	config *sarama.Config
}

func NewClusterAdmin(
	ctx context.Context,
	url,
	username,
	password string,
) (*ClusterAdmin, error) {
	config := sarama.NewConfig()
	config.Metadata.Full = true
	config.Version = sarama.V2_6_0_0
	config.ClientID = "sasl_scram_client"

	config.Net.SASL.Enable = true
	config.Net.SASL.User = username
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
	}
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

	return &ClusterAdmin{
		url:    url,
		config: config,
	}, nil
}

func (c *ClusterAdmin) CreateTopic(topicName string, detail *sarama.TopicDetail) error {
	ca, err := sarama.NewClusterAdmin([]string{c.url}, c.config)
	if err != nil {
		return err
	}
	defer ca.Close()
	return ca.CreateTopic(topicName, detail, false)
}
