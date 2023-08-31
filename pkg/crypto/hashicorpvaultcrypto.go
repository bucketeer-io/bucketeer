// Copyright 2023 The Bucketeer Authors.
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
//

package crypto

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/hashicorp/vault/api"
)

type hashicorpvaultCrypto struct {
	client  *api.Client
	keyName string
}

func NewHashicorpvaultCrypto(
	ctx context.Context,
	keyName string, vaulthost string, vaulttoken string,
) (EncrypterDecrypter, error) {
	config := api.DefaultConfig()

	// default is a localenv vault
	if vaulthost == "" {
		vaulthost = "http://localenv-vault.default.svc.cluster.local:8200"
	}

	config.Address = vaulthost
	hcvClient, err := api.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}

	// authentication with token
	hcvClient.SetToken(vaulttoken)

	return hashicorpvaultCrypto{
		client:  hcvClient,
		keyName: keyName,
	}, nil
}

func (c hashicorpvaultCrypto) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	encPlaintext := base64.StdEncoding.EncodeToString(data)
	// transit vault engine is used
	secret, err := c.client.Logical().Write("transit/encrypt/"+c.keyName, map[string]interface{}{
		"plaintext": encPlaintext,
	})
	if err != nil {
		log.Fatal(err)
	}

	// extract its ciphertext
	ct := secret.Data["ciphertext"]
	if ct == nil {
		log.Fatal("after encrypt operation ciphertext was not found in data returned from vault")
	}
	ctStr, ok := ct.(string)
	if !ok {
		log.Fatal("after encrypt operation ciphertext in data returned from vault is not a string")
	}

	return []byte(ctStr), nil
}

func (c hashicorpvaultCrypto) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	// transit vault engine is used
	desecret, err := c.client.Logical().Write("transit/decrypt/"+c.keyName, map[string]interface{}{
		"ciphertext": string(data),
	})
	if err != nil {
		log.Fatal(err)
	}

	// extract the plaintext value
	pt := desecret.Data["plaintext"]
	ptStr, ok := pt.(string)
	if !ok {
		log.Fatal("after decrypt operation plaintext in data returned from vault is not a string")
	}

	deplaintext, base64Err := base64.StdEncoding.DecodeString(ptStr)
	if base64Err != nil {
		log.Fatal(base64Err)
	}
	return deplaintext, nil
}
