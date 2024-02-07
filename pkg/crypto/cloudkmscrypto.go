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

	cloudkms "cloud.google.com/go/kms/apiv1"
	kms "cloud.google.com/go/kms/apiv1"
	kmsproto "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type cloudKMSCrypto struct {
	client  *cloudkms.KeyManagementClient
	keyName string
}

func NewCloudKMSCrypto(
	ctx context.Context,
	keyName string,
) (EncrypterDecrypter, error) {
	kmsClient, err := kms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}
	return cloudKMSCrypto{
		client:  kmsClient,
		keyName: keyName,
	}, nil
}

func (c cloudKMSCrypto) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	resp, err := c.client.Encrypt(ctx, &kmsproto.EncryptRequest{
		Name:      c.keyName,
		Plaintext: data,
	})
	if err != nil {
		return nil, err
	}
	return resp.Ciphertext, nil
}

func (c cloudKMSCrypto) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	resp, err := c.client.Decrypt(ctx, &kmsproto.DecryptRequest{
		Name:       c.keyName,
		Ciphertext: data,
	})
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}
