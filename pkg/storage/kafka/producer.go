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

type Producer struct {
	sarama.SyncProducer
}

func NewProducer(
	ctx context.Context,
	project,
	url,
	userName,
	password string,
) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	config.Metadata.Full = true
	config.Version = sarama.V2_6_0_0
	config.ClientID = "sasl_scram_client"

	config.Net.SASL.Enable = true
	config.Net.SASL.User = userName
	config.Net.SASL.Password = password
	config.Net.SASL.Handshake = true
	config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
		return &XDGSCRAMClient{HashGeneratorFcn: SHA512}
	}
	config.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512

	prd, err := sarama.NewSyncProducer([]string{url}, config)
	if err != nil {
		return nil, err
	}
	return &Producer{
		SyncProducer: prd,
	}, nil
}
