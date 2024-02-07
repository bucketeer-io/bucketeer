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
//

package api

import (
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	envclient "github.com/bucketeer-io/bucketeer/pkg/environment/client"
	ecclient "github.com/bucketeer-io/bucketeer/pkg/eventcounter/client"
	experimentclient "github.com/bucketeer-io/bucketeer/pkg/experiment/client"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/experimentcalc"
	"github.com/bucketeer-io/bucketeer/pkg/experimentcalculator/stan"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/rpc"
	"github.com/bucketeer-io/bucketeer/pkg/storage/v2/mysql"
	calculator "github.com/bucketeer-io/bucketeer/proto/experimentcalculator"
)

type calculatorService struct {
	calculator *experimentcalc.ExperimentCalculator
	logger     *zap.Logger
}

func NewCalculatorService(
	stan *stan.Stan,
	environmentClient envclient.Client,
	eventCounterClient ecclient.Client,
	experimentClient experimentclient.Client,
	mysqlClient mysql.Client,
	metrics metrics.Registerer,
	loc *time.Location,
	logger *zap.Logger,
) rpc.Service {
	return &calculatorService{
		calculator: experimentcalc.NewExperimentCalculator(
			stan,
			environmentClient,
			eventCounterClient,
			experimentClient,
			mysqlClient,
			metrics,
			loc,
			logger,
		),
		logger: logger,
	}
}

func (c calculatorService) Register(server *grpc.Server) {
	calculator.RegisterExperimentCalculatorServiceServer(server, c)
}
