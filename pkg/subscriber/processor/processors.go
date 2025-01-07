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

package processor

import (
	"errors"

	"github.com/bucketeer-io/bucketeer/pkg/metrics"
	"github.com/bucketeer-io/bucketeer/pkg/subscriber"
)

const (
	AuditLogPersisterName                = "auditLogPersister"
	DomainEventInformerName              = "domainEventInformer"
	EvaluationCountEventDWHPersisterName = "evaluationCountEventDWHPersister"
	EvaluationCountEventOPSPersisterName = "evaluationCountEventOPSPersister"
	EvaluationCountEventPersisterName    = "evaluationCountEventPersister"
	GoalCountEventDWHPersisterName       = "goalCountEventDWHPersister"
	GoalCountEventOPSPersisterName       = "goalCountEventOPSPersister"
	MetricsEventPersisterName            = "metricsEventPersister"
	PushSenderName                       = "pushSender"
	SegmentUserPersisterName             = "segmentUserPersister"
	UserEventPersisterName               = "userEventPersister"
)

var (
	unsupportedProcessorErr = errors.New("subscriber: unsupported processor")
)

type Processors struct {
	processorMap map[string]subscriber.Processor
}

func NewProcessors(r metrics.Registerer) *Processors {
	registerMetrics(r)
	return &Processors{
		processorMap: make(map[string]subscriber.Processor),
	}
}

func (p *Processors) RegisterProcessor(name string, processor subscriber.Processor) {
	p.processorMap[name] = processor
}

func (p *Processors) GetProcessorByName(name string) (subscriber.Processor, error) {
	if p, ok := p.processorMap[name]; ok {
		return p, nil
	}
	return nil, unsupportedProcessorErr
}
