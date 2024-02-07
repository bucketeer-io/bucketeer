// Copyright 2024 The Bucketeer Authors.
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

package crypto

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type awsKMSCrypto struct {
	client *kms.Client
	keyID  string
}

func NewAwsKMSCrypto(
	ctx context.Context,
	keyID, region string,
) (EncrypterDecrypter, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	client := kms.NewFromConfig(cfg)
	return awsKMSCrypto{
		client: client,
		keyID:  keyID,
	}, nil
}

func (c awsKMSCrypto) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	resp, err := c.client.Encrypt(ctx, &kms.EncryptInput{
		Plaintext: data,
		KeyId:     &c.keyID,
	})
	if err != nil {
		return nil, err
	}
	return resp.CiphertextBlob, nil
}

func (c awsKMSCrypto) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	resp, err := c.client.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: data,
		KeyId:          &c.keyID,
	})
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}
