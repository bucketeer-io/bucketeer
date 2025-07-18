// Copyright 2025 The Bucketeer Authors.
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

	"github.com/bucketeer-io/bucketeer/pkg/cli"
)

var (
	name    = "create-mysql-event-tables"
	version = ""
	build   = ""
)

func main() {
	app := cli.NewApp(name, "Bucketeer tool to create MySQL event tables for data warehouse", version, build)
	registerCommand(app, app)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
