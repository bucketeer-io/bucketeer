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

//go:generate mockgen -source=$GOFILE -package=mock -destination=./mock/$GOFILE
package migrate

import (
	"fmt"
	"os"
	"strings"

	libmigrate "github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/github"
)

const mysqlParams = "collation=utf8mb4_bin"

type Client interface {
	Up() error
	Steps(n int) error
}

type client struct {
	*libmigrate.Migrate
}

type ClientFactory interface {
	New(string) (Client, error)
}

type clientFactory struct {
	githubUser        string
	githubAccessToken string
	githubSourcePath  string
	mysqlUser         string
	mysqlPass         string
	mysqlHost         string
	mysqlPort         int
	mysqlDBName       string
}

func NewClientFactory(
	githubUser, githubAccessTokenPath, githubSourcePath string,
	mysqlUser, mysqlPass, mysqlHost string, mysqlPort int, mysqlDBName string,
) (ClientFactory, error) {
	githubAccessToken := ""
	if githubAccessTokenPath != "" {
		data, err := os.ReadFile(githubAccessTokenPath)
		if err != nil {
			return nil, err
		}
		githubAccessToken = strings.TrimSpace(string(data))
	}
	return &clientFactory{
		githubUser:        githubUser,
		githubAccessToken: githubAccessToken,
		githubSourcePath:  githubSourcePath,
		mysqlUser:         mysqlUser,
		mysqlPass:         mysqlPass,
		mysqlHost:         mysqlHost,
		mysqlPort:         mysqlPort,
		mysqlDBName:       mysqlDBName,
	}, nil
}

func (cf *clientFactory) New(ref string) (Client, error) {
	//sourceURL := fmt.Sprintf("github://bucketeer:password@%s", cf.githubSourcePath)
	sourceURL := fmt.Sprintf("github://%s", cf.githubSourcePath)
	if cf.githubUser != "" && cf.githubAccessToken != "" {
		sourceURL = fmt.Sprintf(
			"github://%s:%s@%s",
			cf.githubUser, cf.githubAccessToken, cf.githubSourcePath,
		)
	}
	if ref != "" {
		sourceURL = fmt.Sprintf("%s#%s", sourceURL, ref)
	}
	databaseURL := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%d)/%s?%s",
		cf.mysqlUser, cf.mysqlPass, cf.mysqlHost, cf.mysqlPort, cf.mysqlDBName, mysqlParams,
	)
	m, err := libmigrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, err
	}
	return &client{m}, nil
}
