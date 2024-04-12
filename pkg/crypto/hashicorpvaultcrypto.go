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
	"errors"
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
		return nil, err
	}

	// authentication with token
	hcvClient.SetToken(vaulttoken)

	// check and create key
	err = checkAndCreateKey(hcvClient, keyName)
	if err != nil {
		return nil, err
	}

	return hashicorpvaultCrypto{
		client:  hcvClient,
		keyName: keyName,
	}, nil
}

func checkAndCreateKey(client *api.Client, keyName string) error {
	logical := client.Logical()

	// check if key exists
	path := "transit/keys/" + keyName
	_, err := logical.Read(path)
	if err != nil {
		var respErr *api.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			// create key
			data := map[string]interface{}{
				"type": "aes256-gcm96",
			}
			_, err := logical.Write(path, data)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c hashicorpvaultCrypto) Encrypt(ctx context.Context, data []byte) ([]byte, error) {
	encPlaintext := base64.StdEncoding.EncodeToString(data)
	// transit vault engine is used
	secret, err := c.client.Logical().Write("transit/encrypt/"+c.keyName, map[string]interface{}{
		"plaintext": encPlaintext,
	})
	if err != nil {
		return nil, err
	}

	// extract its ciphertext
	ct := secret.Data["ciphertext"]
	if ct == nil {
		return nil, errors.New("after encrypt operation ciphertext in data returned from vault is nil")
	}
	ctStr, ok := ct.(string)
	if !ok {
		return nil, errors.New("after encrypt operation ciphertext in data returned from vault is not a string")
	}

	return []byte(ctStr), nil
}

func (c hashicorpvaultCrypto) Decrypt(ctx context.Context, data []byte) ([]byte, error) {
	// transit vault engine is used
	desecret, err := c.client.Logical().Write("transit/decrypt/"+c.keyName, map[string]interface{}{
		"ciphertext": string(data),
	})
	if err != nil {
		return nil, err
	}

	// extract the plaintext value
	pt := desecret.Data["plaintext"]
	ptStr, ok := pt.(string)
	if !ok {
		return nil, errors.New("after decrypt operation plaintext in data returned from vault is not a string")
	}

	deplaintext, base64Err := base64.StdEncoding.DecodeString(ptStr)
	if base64Err != nil {
		log.Fatal(base64Err)
	}
	return deplaintext, nil
}
