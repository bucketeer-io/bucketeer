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

package trace

import (
	"contrib.go.opencensus.io/exporter/stackdriver"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
)

func NewStackdriverExporter(service, version string, logger *zap.Logger) (*stackdriver.Exporter, error) {
	return stackdriver.NewExporter(stackdriver.Options{
		OnError: func(err error) {
			logger.Warn("Failed to upload tracing data to Stackdriver", zap.Error(err))
		},
		DefaultTraceAttributes: map[string]interface{}{
			"Service": service,
			"Version": version,
		},
	})
}

type sampler struct {
	probability       float64
	filteringSamplers map[string]trace.Sampler
}

type SamplerOption func(*sampler)

func WithDefaultProbability(p float64) SamplerOption {
	return func(s *sampler) {
		s.probability = p
	}
}

func WithFilteringSampler(name string, fs trace.Sampler) SamplerOption {
	return func(s *sampler) {
		s.filteringSamplers[name] = fs
	}
}

func NewSampler(options ...SamplerOption) trace.Sampler {
	s := &sampler{
		probability:       0.01,
		filteringSamplers: make(map[string]trace.Sampler),
	}
	for _, opt := range options {
		opt(s)
	}
	return s.sampler
}

func (s *sampler) sampler(p trace.SamplingParameters) trace.SamplingDecision {
	if fs, ok := s.filteringSamplers[p.Name]; ok {
		return fs(p)
	}
	return trace.ProbabilitySampler(s.probability)(p)
}
