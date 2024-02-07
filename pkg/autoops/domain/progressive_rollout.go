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

package domain

import (
	"errors"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/bucketeer-io/bucketeer/pkg/uuid"
	autoopsproto "github.com/bucketeer-io/bucketeer/proto/autoops"
)

var (
	ErrProgressiveRolloutScheduleNotFound  = errors.New("progressiveRollout: schedule not found")
	ErrProgressiveRolloutInvalidType       = errors.New("progressiveRollout: invalid type")
	ErrProgressiveRolloutStoopedByRequired = errors.New("progressiveRollout: stopped by is required")
)

type ProgressiveRollout struct {
	*autoopsproto.ProgressiveRollout
}

func NewProgressiveRollout(
	featureID string,
	manual *autoopsproto.ProgressiveRolloutManualScheduleClause,
	template *autoopsproto.ProgressiveRolloutTemplateScheduleClause,
) (*ProgressiveRollout, error) {
	now := time.Now().Unix()
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}
	p := &ProgressiveRollout{&autoopsproto.ProgressiveRollout{
		Id:        id.String(),
		FeatureId: featureID,
		Status:    autoopsproto.ProgressiveRollout_WAITING,
		StoppedBy: autoopsproto.ProgressiveRollout_UNKNOWN,
		Clause:    nil,
		CreatedAt: now,
		UpdatedAt: now,
	}}
	if manual != nil {
		if err := p.addManualScheduleClause(manual); err != nil {
			return nil, err
		}
	}
	if template != nil {
		if err := p.addTemplatelScheduleClause(template); err != nil {
			return nil, err
		}
	}
	return p, nil
}

func (p *ProgressiveRollout) addManualScheduleClause(
	c *autoopsproto.ProgressiveRolloutManualScheduleClause,
) error {
	p.Type = autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE
	for _, s := range c.Schedules {
		id, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		s.ScheduleId = id.String()
	}
	return p.setClause(c)
}

func (p *ProgressiveRollout) addTemplatelScheduleClause(
	c *autoopsproto.ProgressiveRolloutTemplateScheduleClause,
) error {
	p.Type = autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE
	for _, s := range c.Schedules {
		id, err := uuid.NewUUID()
		if err != nil {
			return err
		}
		s.ScheduleId = id.String()
	}
	return p.setClause(c)
}

func (p *ProgressiveRollout) setClause(c protoiface.MessageV1) error {
	ac, err := ptypes.MarshalAny(c)
	if err != nil {
		return err
	}
	p.Clause = ac
	p.UpdatedAt = time.Now().Unix()
	return nil
}

func (p *ProgressiveRollout) IsWaiting() bool {
	return p.Status == autoopsproto.ProgressiveRollout_WAITING
}

func (p *ProgressiveRollout) IsRunning() bool {
	return p.Status == autoopsproto.ProgressiveRollout_RUNNING
}

func (p *ProgressiveRollout) IsStopped() bool {
	return p.Status == autoopsproto.ProgressiveRollout_STOPPED
}

func (p *ProgressiveRollout) IsFinished() bool {
	return p.Status == autoopsproto.ProgressiveRollout_FINISHED
}

func (p *ProgressiveRollout) AlreadyTriggered(scheduleID string) (bool, error) {
	switch p.Type {
	case autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE:
		c, err := unmarshalProgressiveRolloutManualClause(p.Clause)
		if err != nil {
			return false, err
		}
		s, err := findTargetSchedule(c.Schedules, scheduleID)
		if err != nil {
			return false, err
		}
		return s.TriggeredAt > 0, nil
	case autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE:
		c, err := unmarshalProgressiveRolloutTemplateClause(p.Clause)
		if err != nil {
			return false, err
		}
		s, err := findTargetSchedule(c.Schedules, scheduleID)
		if err != nil {
			return false, err
		}
		return s.TriggeredAt > 0, nil
	}
	return false, ErrProgressiveRolloutInvalidType
}

func (p *ProgressiveRollout) ExtractSchedules() ([]*autoopsproto.ProgressiveRolloutSchedule, error) {
	switch p.Type {
	case autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE:
		c, err := unmarshalProgressiveRolloutManualClause(p.Clause)
		if err != nil {
			return nil, err
		}
		return c.Schedules, nil
	case autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE:
		c, err := unmarshalProgressiveRolloutTemplateClause(p.Clause)
		if err != nil {
			return nil, err
		}
		return c.Schedules, nil
	}
	return nil, ErrProgressiveRolloutInvalidType
}

func unmarshalProgressiveRolloutManualClause(
	clause *anypb.Any,
) (*autoopsproto.ProgressiveRolloutManualScheduleClause, error) {
	c := &autoopsproto.ProgressiveRolloutManualScheduleClause{}
	if err := ptypes.UnmarshalAny(clause, c); err != nil {
		return nil, err
	}
	return c, nil
}

func unmarshalProgressiveRolloutTemplateClause(
	clause *anypb.Any,
) (*autoopsproto.ProgressiveRolloutTemplateScheduleClause, error) {
	c := &autoopsproto.ProgressiveRolloutTemplateScheduleClause{}
	if err := ptypes.UnmarshalAny(clause, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (p *ProgressiveRollout) SetTriggeredAt(scheduleID string) error {
	now := time.Now().Unix()
	switch p.Type {
	case autoopsproto.ProgressiveRollout_MANUAL_SCHEDULE:
		c, err := unmarshalProgressiveRolloutManualClause(p.Clause)
		if err != nil {
			return err
		}
		s, err := findTargetSchedule(c.Schedules, scheduleID)
		if err != nil {
			return err
		}
		s.TriggeredAt = now
		if err := p.setClause(c); err != nil {
			return err
		}
		p.Status = autoopsproto.ProgressiveRollout_RUNNING
		if c.Schedules[len(c.Schedules)-1].ScheduleId == scheduleID {
			p.Status = autoopsproto.ProgressiveRollout_FINISHED
		}
	case autoopsproto.ProgressiveRollout_TEMPLATE_SCHEDULE:
		c, err := unmarshalProgressiveRolloutTemplateClause(p.Clause)
		if err != nil {
			return err
		}
		s, err := findTargetSchedule(c.Schedules, scheduleID)
		if err != nil {
			return err
		}
		s.TriggeredAt = now
		if err := p.setClause(c); err != nil {
			return err
		}
		p.Status = autoopsproto.ProgressiveRollout_RUNNING
		if c.Schedules[len(c.Schedules)-1].ScheduleId == scheduleID {
			p.Status = autoopsproto.ProgressiveRollout_FINISHED
		}
	default:
		return ErrProgressiveRolloutInvalidType
	}
	p.ProgressiveRollout.UpdatedAt = now
	return nil
}

func findTargetSchedule(
	schedules []*autoopsproto.ProgressiveRolloutSchedule,
	scheduleID string,
) (*autoopsproto.ProgressiveRolloutSchedule, error) {
	for _, s := range schedules {
		if s.ScheduleId == scheduleID {
			return s, nil
		}
	}
	return nil, ErrProgressiveRolloutScheduleNotFound
}

func (p *ProgressiveRollout) Stop(stoppedBy autoopsproto.ProgressiveRollout_StoppedBy) error {
	if stoppedBy == autoopsproto.ProgressiveRollout_UNKNOWN {
		return ErrProgressiveRolloutStoopedByRequired
	}
	now := time.Now().Unix()
	p.StoppedBy = stoppedBy
	p.Status = autoopsproto.ProgressiveRollout_STOPPED
	p.StoppedAt = now
	p.UpdatedAt = now
	return nil
}
