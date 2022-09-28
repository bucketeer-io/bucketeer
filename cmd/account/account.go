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

package main

import (
	"log"

	"github.com/bucketeer-io/bucketeer/pkg/account/cmd/apikeycacher"
	"github.com/bucketeer-io/bucketeer/pkg/account/cmd/server"
	"github.com/bucketeer-io/bucketeer/pkg/cli"
)

var (
	name    = "bucketeer-account"
	version = ""
	build   = ""
)

func main() {
	app := cli.NewApp(name, "A/B Testing Microservice", version, build)
	registerCommands(app)
	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func registerCommands(app *cli.App) {
	server.RegisterCommand(app, app)
	apikeycacher.RegisterCommand(app, app)
}
