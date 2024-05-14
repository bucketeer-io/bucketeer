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

package processor

import (
	"errors"

	"github.com/bucketeer-io/bucketeer/pkg/batch/subscriber"
	"github.com/bucketeer-io/bucketeer/pkg/metrics"
)

const (
	DomainEventInformerName           = "domainEventInformer"
	EvaluationCountEventPersisterName = "evaluationCountEventPersister"
	SegmentUserPersisterName          = "segmentUserPersister"
	UserEventPersisterName            = "userEventPersister"
)

const (
	TypeNormal     = "normal"
	TypeServerless = "serverless"
)

var (
	unsupportedProcessorErr = errors.New("subscriber: unsupported processor")
)

type Processors struct {
	processorMap map[string]map[string]subscriber.Processor
}

func NewProcessors(r metrics.Registerer) *Processors {
	registerMetrics(r)
	return &Processors{
		processorMap: map[string]map[string]subscriber.Processor{
			TypeNormal:     make(map[string]subscriber.Processor),
			TypeServerless: make(map[string]subscriber.Processor),
		},
	}
}

func (p *Processors) RegisterProcessor(
	processorType, name string,
	processor subscriber.Processor,
) {
	p.processorMap[processorType][name] = processor
}

func (p *Processors) GetProcessorByName(processorType, name string) (subscriber.Processor, error) {
	if p, ok := p.processorMap[processorType][name]; ok {
		return p, nil
	}
	return nil, unsupportedProcessorErr
}
