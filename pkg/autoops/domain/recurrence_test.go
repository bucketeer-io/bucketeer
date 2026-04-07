// Copyright 2026 The Bucketeer Authors.
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
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	autoopsproto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

func TestCalculateNextExecution_NilClause(t *testing.T) {
	t.Parallel()
	nextExec, shouldContinue := CalculateNextExecution(nil, time.Now())
	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_NilRecurrence(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
	}
	nextExec, shouldContinue := CalculateNextExecution(clause, time.Now())
	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_InvalidTimeOfDay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		time int64
	}{
		{"negative", -1},
		{"exactly 86400", 86400},
		{"large value", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			clause := &autoopsproto.DatetimeClause{
				Time: tt.time,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency: autoopsproto.RecurrenceRule_DAILY,
					Timezone:  "UTC",
					StartDate: time.Now().Add(-24 * time.Hour).Unix(),
				},
			}
			nextExec, shouldContinue := CalculateNextExecution(clause, time.Now())
			assert.False(t, shouldContinue)
			assert.Equal(t, int64(0), nextExec)
		})
	}
}

func TestCalculateNextExecution_OnceFrequency(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_ONCE,
		},
	}
	nextExec, shouldContinue := CalculateNextExecution(clause, time.Now())
	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_Daily(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_DAILY,
			Timezone:  "Asia/Tokyo",
			StartDate: time.Date(2026, 2, 10, 10, 0, 0, 0, jst).Unix(),
		},
		ExecutionCount: 0,
	}

	executedAt := time.Date(2026, 2, 10, 10, 0, 0, 0, jst)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	expected := time.Date(2026, 2, 11, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_Weekly_SingleDay(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1}, // Monday
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Date(2026, 2, 9, 10, 0, 0, 0, jst).Unix(),
		},
		ExecutionCount: 0,
	}

	// Executed on Monday Feb 9 at 10:00 JST
	executedAt := time.Date(2026, 2, 9, 10, 0, 0, 0, jst)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	// Next Monday is Feb 16
	expected := time.Date(2026, 2, 16, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_Weekly_MultipleDays(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1, 2}, // Monday, Tuesday
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Date(2026, 2, 9, 10, 0, 0, 0, jst).Unix(),
		},
		ExecutionCount: 0,
	}

	// Executed on Monday Feb 9 at 10:00 JST
	executedAt := time.Date(2026, 2, 9, 10, 0, 0, 0, jst)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	// Next scheduled day is Tuesday Feb 10
	expected := time.Date(2026, 2, 10, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_Weekly_WrapToNextWeek(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1, 5}, // Monday, Friday
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Date(2026, 2, 9, 10, 0, 0, 0, jst).Unix(),
		},
		ExecutionCount: 1,
	}

	// Executed on Friday Feb 13 at 10:00 JST
	executedAt := time.Date(2026, 2, 13, 10, 0, 0, 0, jst)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	// Next scheduled day wraps to Monday Feb 16
	expected := time.Date(2026, 2, 16, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_Weekly_EmptyDaysOfWeek(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{},
			Timezone:   "UTC",
			StartDate:  time.Now().Add(-24 * time.Hour).Unix(),
		},
	}

	nextExec, shouldContinue := CalculateNextExecution(clause, time.Now())
	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_Monthly(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
			DayOfMonth: 15,
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Date(2026, 1, 15, 10, 0, 0, 0, jst).Unix(),
		},
		ExecutionCount: 0,
	}

	executedAt := time.Date(2026, 1, 15, 10, 0, 0, 0, jst)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	expected := time.Date(2026, 2, 15, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_Monthly_DayOverflow(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
			DayOfMonth: 31,
			Timezone:   "UTC",
			StartDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
		ExecutionCount: 0,
	}

	// Execute on Jan 31 - next should skip Feb (no 31st) and go to Mar 31
	executedAt := time.Date(2026, 1, 31, 10, 0, 0, 0, time.UTC)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	expected := time.Date(2026, 3, 31, 10, 0, 0, 0, time.UTC)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_Monthly_InvalidDay(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
			DayOfMonth: 0,
			Timezone:   "UTC",
			StartDate:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}

	nextExec, shouldContinue := CalculateNextExecution(clause, time.Now())
	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_MaxOccurrences(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:      autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek:     []int32{1},
			Timezone:       "Asia/Tokyo",
			StartDate:      time.Date(2026, 2, 9, 10, 0, 0, 0, time.UTC).Unix(),
			MaxOccurrences: 5,
		},
		ExecutionCount: 5,
	}

	executedAt := time.Date(2026, 3, 10, 1, 0, 0, 0, time.UTC)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_MaxOccurrences_NotReached(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:      autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek:     []int32{1},
			Timezone:       "Asia/Tokyo",
			StartDate:      time.Date(2026, 2, 9, 10, 0, 0, 0, jst).Unix(),
			MaxOccurrences: 5,
		},
		ExecutionCount: 3,
	}

	executedAt := time.Date(2026, 2, 23, 10, 0, 0, 0, jst) // Monday Feb 23
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	expected := time.Date(2026, 3, 2, 10, 0, 0, 0, jst) // Next Monday Mar 2
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_EndDate_Passed(t *testing.T) {
	t.Parallel()
	endDate := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Date(2026, 2, 9, 10, 0, 0, 0, time.UTC).Unix(),
			EndDate:    endDate.Unix(),
		},
		ExecutionCount: 3,
	}

	executedAt := time.Date(2026, 3, 5, 1, 0, 0, 0, time.UTC)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_EndDate_NextExecAfterEnd(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	endDate := time.Date(2026, 2, 20, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1}, // Monday
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Date(2026, 2, 9, 10, 0, 0, 0, jst).Unix(),
			EndDate:    endDate.Unix(),
		},
		ExecutionCount: 1,
	}

	// Executed on Monday Feb 16 - next would be Feb 23 which is past end
	executedAt := time.Date(2026, 2, 16, 10, 0, 0, 0, jst)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestCalculateNextExecution_DifferentTimezones(t *testing.T) {
	t.Parallel()
	ny, err := time.LoadLocation("America/New_York")
	require.NoError(t, err)

	clause := &autoopsproto.DatetimeClause{
		Time: 32400, // 9:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_DAILY,
			Timezone:  "America/New_York",
			StartDate: time.Date(2026, 3, 10, 9, 0, 0, 0, ny).Unix(),
		},
		ExecutionCount: 0,
	}

	executedAt := time.Date(2026, 3, 10, 9, 0, 0, 0, ny)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	expected := time.Date(2026, 3, 11, 9, 0, 0, 0, ny)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestCalculateNextExecution_InvalidTimezone(t *testing.T) {
	t.Parallel()

	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_DAILY,
			Timezone:  "Invalid/Timezone",
			StartDate: time.Now().Add(-24 * time.Hour).Unix(),
		},
		ExecutionCount: 0,
	}

	// Should fall back to UTC
	executedAt := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	nextExec, shouldContinue := CalculateNextExecution(clause, executedAt)

	assert.True(t, shouldContinue)
	expected := time.Date(2026, 2, 11, 10, 0, 0, 0, time.UTC)
	assert.Equal(t, expected.Unix(), nextExec)
}

func TestInitializeRecurringClause_NonRecurring(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
	}
	err := InitializeRecurringClause(clause)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), clause.NextExecutionAt)
}

func TestInitializeRecurringClause_OnceFrequency(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_ONCE,
		},
	}
	err := InitializeRecurringClause(clause)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), clause.NextExecutionAt)
}

func TestInitializeRecurringClause_MissingStartDate(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
		},
	}
	err := InitializeRecurringClause(clause)
	assert.Error(t, err)
}

func TestInitializeRecurringClause_Weekly(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	// Start date is Sunday Feb 8
	startDate := time.Date(2026, 2, 8, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1}, // Monday
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.NoError(t, err)

	// First execution should be Monday Feb 9 at 10:00 JST
	expected := time.Date(2026, 2, 9, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), clause.NextExecutionAt)
	assert.Equal(t, int32(0), clause.ExecutionCount)
	assert.Equal(t, int64(0), clause.LastExecutedAt)
}

func TestInitializeRecurringClause_Daily(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 10, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_DAILY,
			Timezone:  "Asia/Tokyo",
			StartDate: startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.NoError(t, err)

	// First execution should be start date at 10:00 AM JST
	expected := time.Date(2026, 2, 10, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), clause.NextExecutionAt)
}

func TestInitializeRecurringClause_Daily_TimeOfDayAlreadyPassed(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	// Start date is Feb 10 at 15:00 JST, but scheduled time is 10:00 AM (already past)
	startDate := time.Date(2026, 2, 10, 15, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_DAILY,
			Timezone:  "Asia/Tokyo",
			StartDate: startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.NoError(t, err)

	// Should advance to next day: Feb 11 at 10:00 AM JST
	expected := time.Date(2026, 2, 11, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), clause.NextExecutionAt)
}

func TestInitializeRecurringClause_Weekly_SameDayButTimePassed(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	// Start date is Monday Feb 9 at 15:00 JST, scheduled for Monday 10:00 AM
	startDate := time.Date(2026, 2, 9, 15, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1}, // Monday
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.NoError(t, err)

	// Monday 10:00 AM is before 15:00, so should wrap to next Monday Feb 16
	expected := time.Date(2026, 2, 16, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), clause.NextExecutionAt)
}

func TestInitializeRecurringClause_EndDateBeforeFirstExecution(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	// Start date is Feb 8 (Sunday), scheduled for Monday, but end_date is Feb 9
	// First execution would be Monday Feb 9 10:00, but end_date is Feb 9 00:00
	startDate := time.Date(2026, 2, 8, 0, 0, 0, 0, jst)
	endDate := time.Date(2026, 2, 9, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1}, // Monday
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
			EndDate:    endDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.NoError(t, err)
	// First execution (Mon Feb 9 10:00) is at or after end_date (Feb 9 00:00)
	// so NextExecutionAt should be 0 (already finished)
	assert.Equal(t, int64(0), clause.NextExecutionAt)
}

func TestInitializeRecurringClause_Monthly_InvalidDayOfMonth(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
			DayOfMonth: 0, // invalid
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "day_of_month must be 1-31")
}

func TestInitializeRecurringClause_Weekly_EmptyDaysOfWeek(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{},
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "days_of_week must be non-empty")
}

func TestInitializeRecurringClause_Weekly_InvalidDayValue(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1, 8}, // 8 is invalid
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "days_of_week must be non-empty")
}

func TestInitializeRecurringClause_Monthly(t *testing.T) {
	t.Parallel()
	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	startDate := time.Date(2026, 2, 1, 0, 0, 0, 0, jst)
	clause := &autoopsproto.DatetimeClause{
		Time: 36000, // 10:00 AM
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
			DayOfMonth: 15,
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	}
	err = InitializeRecurringClause(clause)
	assert.NoError(t, err)

	// First execution should be Feb 15 at 10:00 AM JST
	expected := time.Date(2026, 2, 15, 10, 0, 0, 0, jst)
	assert.Equal(t, expected.Unix(), clause.NextExecutionAt)
}

func TestIsRecurring(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc     string
		clause   *autoopsproto.DatetimeClause
		expected bool
	}{
		{
			desc:     "nil recurrence",
			clause:   &autoopsproto.DatetimeClause{Time: 36000},
			expected: false,
		},
		{
			desc: "ONCE frequency",
			clause: &autoopsproto.DatetimeClause{
				Time: 36000,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency: autoopsproto.RecurrenceRule_ONCE,
				},
			},
			expected: false,
		},
		{
			desc: "FREQUENCY_UNSPECIFIED",
			clause: &autoopsproto.DatetimeClause{
				Time: 36000,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency: autoopsproto.RecurrenceRule_FREQUENCY_UNSPECIFIED,
				},
			},
			expected: false,
		},
		{
			desc: "DAILY",
			clause: &autoopsproto.DatetimeClause{
				Time: 36000,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency: autoopsproto.RecurrenceRule_DAILY,
				},
			},
			expected: true,
		},
		{
			desc: "WEEKLY",
			clause: &autoopsproto.DatetimeClause{
				Time: 36000,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
				},
			},
			expected: true,
		},
		{
			desc: "MONTHLY",
			clause: &autoopsproto.DatetimeClause{
				Time: 36000,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_MONTHLY,
					DayOfMonth: 15,
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expected, IsRecurring(tt.clause))
		})
	}
}

func TestAutoOpsRule_IsRecurringSchedule(t *testing.T) {
	t.Parallel()

	t.Run("non-schedule rule returns false", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_EVENT_RATE,
			[]*autoopsproto.OpsEventRateClause{
				{
					GoalId:          "goalid01",
					MinCount:        10,
					ThreadsholdRate: 0.5,
					Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
					ActionType:      autoopsproto.ActionType_DISABLE,
				},
			},
			[]*autoopsproto.DatetimeClause{},
		)
		require.NoError(t, err)
		assert.False(t, aor.IsRecurringSchedule())
	})

	t.Run("one-time schedule returns false", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
			},
		)
		require.NoError(t, err)
		assert.False(t, aor.IsRecurringSchedule())
	})

	t.Run("recurring schedule returns true", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "Asia/Tokyo",
						StartDate:  time.Now().Add(24 * time.Hour).Unix(),
					},
				},
			},
		)
		require.NoError(t, err)
		assert.True(t, aor.IsRecurringSchedule())
	})
}

func TestAutoOpsRule_AllClausesFinished(t *testing.T) {
	t.Parallel()

	t.Run("one-time clause not executed", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
			},
		)
		require.NoError(t, err)
		assert.False(t, aor.AllClausesFinished())
	})

	t.Run("one-time clause executed", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
			},
		)
		require.NoError(t, err)
		aor.Clauses[0].ExecutedAt = time.Now().Unix()
		assert.True(t, aor.AllClausesFinished())
	})

	t.Run("recurring clause with future execution", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{
					Time:            36000,
					ActionType:      autoopsproto.ActionType_ENABLE,
					NextExecutionAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "Asia/Tokyo",
						StartDate:  time.Now().Add(-24 * time.Hour).Unix(),
					},
				},
			},
		)
		require.NoError(t, err)
		assert.False(t, aor.AllClausesFinished())
	})

	t.Run("recurring clause completed", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{
					Time:            36000,
					ActionType:      autoopsproto.ActionType_ENABLE,
					NextExecutionAt: 0,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:      autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek:     []int32{1},
						Timezone:       "Asia/Tokyo",
						StartDate:      time.Now().Add(-24 * time.Hour).Unix(),
						MaxOccurrences: 5,
					},
					ExecutionCount: 5,
				},
			},
		)
		require.NoError(t, err)
		// The clause is recurring but NextExecutionAt is 0, meaning completed
		assert.True(t, aor.AllClausesFinished())
	})
}

func TestAutoOpsRule_GetNextExecutionTime(t *testing.T) {
	t.Parallel()

	t.Run("one-time not yet executed", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
			},
		)
		require.NoError(t, err)
		nextExec, err := aor.GetNextExecutionTime()
		assert.NoError(t, err)
		assert.Equal(t, int64(1000000000), nextExec)
	})

	t.Run("recurring with next execution", func(t *testing.T) {
		futureTime := time.Now().Add(7 * 24 * time.Hour).Unix()
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{
					Time:            36000,
					ActionType:      autoopsproto.ActionType_ENABLE,
					NextExecutionAt: futureTime,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "Asia/Tokyo",
						StartDate:  time.Now().Add(-24 * time.Hour).Unix(),
					},
				},
			},
		)
		require.NoError(t, err)
		nextExec, err := aor.GetNextExecutionTime()
		assert.NoError(t, err)
		assert.Equal(t, futureTime, nextExec)
	})

	t.Run("all executed returns zero", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
			},
		)
		require.NoError(t, err)
		aor.Clauses[0].ExecutedAt = time.Now().Unix()
		nextExec, err := aor.GetNextExecutionTime()
		assert.NoError(t, err)
		assert.Equal(t, int64(0), nextExec)
	})

	t.Run("recurring with NextExecutionAt=0 is skipped not treated as one-time", func(t *testing.T) {
		aor, err := NewAutoOpsRule(
			"feature-id",
			autoopsproto.OpsType_SCHEDULE,
			[]*autoopsproto.OpsEventRateClause{},
			[]*autoopsproto.DatetimeClause{
				{
					Time:            36000,
					ActionType:      autoopsproto.ActionType_ENABLE,
					NextExecutionAt: 0,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:      autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek:     []int32{1},
						Timezone:       "Asia/Tokyo",
						StartDate:      time.Now().Add(-24 * time.Hour).Unix(),
						MaxOccurrences: 5,
					},
					ExecutionCount: 5,
				},
			},
		)
		require.NoError(t, err)
		nextExec, err := aor.GetNextExecutionTime()
		assert.NoError(t, err)
		// Must return 0, NOT 36000 (seconds-since-midnight)
		assert.Equal(t, int64(0), nextExec)
	})
}

func TestAddDatetimeClause_SetsIsRecurring(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)

	t.Run("non-recurring clause has is_recurring=false", func(t *testing.T) {
		clause, err := aor.AddDatetimeClause(&autoopsproto.DatetimeClause{
			Time:       1000000000,
			ActionType: autoopsproto.ActionType_ENABLE,
		})
		require.NoError(t, err)
		assert.False(t, clause.IsRecurring)
	})

	t.Run("recurring clause has is_recurring=true", func(t *testing.T) {
		clause, err := aor.AddDatetimeClause(&autoopsproto.DatetimeClause{
			Time:       36000,
			ActionType: autoopsproto.ActionType_ENABLE,
			Recurrence: &autoopsproto.RecurrenceRule{
				Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
				DaysOfWeek: []int32{1},
				Timezone:   "Asia/Tokyo",
				StartDate:  time.Now().Add(24 * time.Hour).Unix(),
			},
		})
		require.NoError(t, err)
		assert.True(t, clause.IsRecurring)
	})
}

func TestAddDatetimeClause_InitializesRecurringClause(t *testing.T) {
	t.Parallel()
	aor := createAutoOpsRule(t)

	startDate := time.Now().Add(24 * time.Hour)
	clause, err := aor.AddDatetimeClause(&autoopsproto.DatetimeClause{
		Time:       36000,
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
			StartDate:  startDate.Unix(),
		},
	})
	require.NoError(t, err)
	assert.True(t, clause.IsRecurring)

	dtClauses, err := aor.ExtractDatetimeClauses()
	require.NoError(t, err)
	dtClause := dtClauses[clause.Id]
	require.NotNil(t, dtClause)
	assert.True(t, dtClause.NextExecutionAt > 0, "NextExecutionAt should be initialized")
	assert.Equal(t, int32(0), dtClause.ExecutionCount)
}

func TestCalculateNextExecution_Weekly_InvalidDaysOfWeek(t *testing.T) {
	t.Parallel()
	clause := &autoopsproto.DatetimeClause{
		Time: 36000,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{7}, // invalid: must be 0-6
			Timezone:   "UTC",
			StartDate:  time.Now().Add(-24 * time.Hour).Unix(),
		},
	}
	nextExec, shouldContinue := CalculateNextExecution(clause, time.Now())
	assert.False(t, shouldContinue)
	assert.Equal(t, int64(0), nextExec)
}

func TestIsRecurring_NilClause(t *testing.T) {
	t.Parallel()
	assert.False(t, IsRecurring(nil))
}

func TestSortDatetimeClause_MixedRecurringAndOneTime(t *testing.T) {
	t.Parallel()

	futureTimestamp := time.Now().Add(48 * time.Hour).Unix()       // one-time: 48h from now
	nextExecTimestamp := time.Now().Add(7 * 24 * time.Hour).Unix() // recurring: 7 days from now

	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{
				Time:            36000,
				ActionType:      autoopsproto.ActionType_ENABLE,
				NextExecutionAt: nextExecTimestamp,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
					Timezone:   "Asia/Tokyo",
					StartDate:  time.Now().Add(-24 * time.Hour).Unix(),
				},
			},
		},
	)
	require.NoError(t, err)

	// Add a one-time clause that should execute sooner
	_, err = aor.AddDatetimeClause(&autoopsproto.DatetimeClause{
		Time:       futureTimestamp,
		ActionType: autoopsproto.ActionType_DISABLE,
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(aor.Clauses))

	// The one-time clause (48h) should sort before the recurring clause (7 days)
	dtClauses, err := aor.ExtractDatetimeClauses()
	require.NoError(t, err)

	firstDt := dtClauses[aor.Clauses[0].Id]
	require.NotNil(t, firstDt)
	assert.Equal(t, futureTimestamp, firstDt.Time)
	assert.False(t, aor.Clauses[0].IsRecurring)

	secondDt := dtClauses[aor.Clauses[1].Id]
	require.NotNil(t, secondDt)
	assert.True(t, aor.Clauses[1].IsRecurring)
}

func TestSortDatetimeClause_RecurringClauseSortsByNextExecution(t *testing.T) {
	t.Parallel()

	now := time.Now()
	earlyNext := now.Add(1 * 24 * time.Hour).Unix()
	laterNext := now.Add(3 * 24 * time.Hour).Unix()

	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{
				Time:            64800, // 6:00 PM
				ActionType:      autoopsproto.ActionType_DISABLE,
				NextExecutionAt: laterNext,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{5},
					Timezone:   "UTC",
					StartDate:  now.Add(-48 * time.Hour).Unix(),
				},
			},
		},
	)
	require.NoError(t, err)

	_, err = aor.AddDatetimeClause(&autoopsproto.DatetimeClause{
		Time:            32400, // 9:00 AM
		ActionType:      autoopsproto.ActionType_ENABLE,
		NextExecutionAt: earlyNext,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "UTC",
			StartDate:  now.Add(-48 * time.Hour).Unix(),
		},
	})
	require.NoError(t, err)

	// The clause with earlyNext should be first despite having a higher time-of-day
	dtClauses, err := aor.ExtractDatetimeClauses()
	require.NoError(t, err)

	firstDt := dtClauses[aor.Clauses[0].Id]
	require.NotNil(t, firstDt)
	assert.Equal(t, earlyNext, firstDt.NextExecutionAt)

	secondDt := dtClauses[aor.Clauses[1].Id]
	require.NotNil(t, secondDt)
	assert.Equal(t, laterNext, secondDt.NextExecutionAt)
}

func TestSortDatetimeClause_CompletedRecurringSortsLast(t *testing.T) {
	t.Parallel()

	futureTimestamp := time.Now().Add(48 * time.Hour).Unix()

	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{
				Time:            36000,
				ActionType:      autoopsproto.ActionType_ENABLE,
				NextExecutionAt: 0, // completed
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:      autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek:     []int32{1},
					Timezone:       "UTC",
					StartDate:      time.Now().Add(-720 * time.Hour).Unix(),
					MaxOccurrences: 5,
				},
				ExecutionCount: 5,
			},
		},
	)
	require.NoError(t, err)

	_, err = aor.AddDatetimeClause(&autoopsproto.DatetimeClause{
		Time:       futureTimestamp,
		ActionType: autoopsproto.ActionType_DISABLE,
	})
	require.NoError(t, err)
	require.Equal(t, 2, len(aor.Clauses))

	// Completed recurring clause (NextExecutionAt=0) must sort AFTER the one-time clause
	assert.False(t, aor.Clauses[0].IsRecurring)
	assert.True(t, aor.Clauses[1].IsRecurring)
}

func TestUpdateAutoOpsRule_GranularCreate_SetsIsRecurring(t *testing.T) {
	t.Parallel()

	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
		},
	)
	require.NoError(t, err)
	assert.False(t, aor.Clauses[0].IsRecurring)

	updated, err := aor.Update(nil, nil, []*autoopsproto.DatetimeClauseChange{
		{
			ChangeType: autoopsproto.ChangeType_CREATE,
			Clause: &autoopsproto.DatetimeClause{
				Time:       36000,
				ActionType: autoopsproto.ActionType_DISABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
					DaysOfWeek: []int32{1},
					Timezone:   "UTC",
					StartDate:  time.Now().Add(24 * time.Hour).Unix(),
				},
			},
		},
	})
	require.NoError(t, err)

	var recurringClause *autoopsproto.Clause
	for _, c := range updated.Clauses {
		if c.IsRecurring {
			recurringClause = c
		}
	}
	require.NotNil(t, recurringClause, "granular CREATE should set is_recurring")

	dtClauses, err := updated.ExtractDatetimeClauses()
	require.NoError(t, err)
	dtClause := dtClauses[recurringClause.Id]
	require.NotNil(t, dtClause)
	assert.True(t, dtClause.NextExecutionAt > 0,
		"granular CREATE should initialize NextExecutionAt")
}

func TestGetNextExecutionTime_NonScheduleRule(t *testing.T) {
	t.Parallel()

	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_EVENT_RATE,
		[]*autoopsproto.OpsEventRateClause{
			{
				GoalId:          "goalid01",
				MinCount:        10,
				ThreadsholdRate: 0.5,
				Operator:        autoopsproto.OpsEventRateClause_GREATER_OR_EQUAL,
				ActionType:      autoopsproto.ActionType_DISABLE,
			},
		},
		[]*autoopsproto.DatetimeClause{},
	)
	require.NoError(t, err)

	_, err = aor.GetNextExecutionTime()
	assert.Error(t, err)
}

func TestChangeDatetimeClause_SetsIsRecurring(t *testing.T) {
	t.Parallel()
	aor, err := NewAutoOpsRule(
		"feature-id",
		autoopsproto.OpsType_SCHEDULE,
		[]*autoopsproto.OpsEventRateClause{},
		[]*autoopsproto.DatetimeClause{
			{Time: 1000000000, ActionType: autoopsproto.ActionType_ENABLE},
		},
	)
	require.NoError(t, err)
	clauseID := aor.Clauses[0].Id
	assert.False(t, aor.Clauses[0].IsRecurring)

	err = aor.ChangeDatetimeClause(clauseID, &autoopsproto.DatetimeClause{
		Time:       36000,
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "Asia/Tokyo",
			StartDate:  time.Now().Add(24 * time.Hour).Unix(),
		},
	})
	require.NoError(t, err)
	assert.True(t, aor.Clauses[0].IsRecurring)

	dtClauses, err := aor.ExtractDatetimeClauses()
	require.NoError(t, err)
	dtClause := dtClauses[clauseID]
	require.NotNil(t, dtClause)
	assert.True(t, dtClause.NextExecutionAt > 0,
		"ChangeDatetimeClause should initialize NextExecutionAt")
}

func TestInitializeRecurringClause_InvalidTimeOfDay(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc string
		time int64
	}{
		{"negative", -1},
		{"exactly 86400", 86400},
		{"large value", 100000},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			clause := &autoopsproto.DatetimeClause{
				Time: tt.time,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency: autoopsproto.RecurrenceRule_DAILY,
					Timezone:  "UTC",
					StartDate: time.Now().Add(24 * time.Hour).Unix(),
				},
			}
			err := InitializeRecurringClause(clause)
			assert.Error(t, err)
		})
	}
}

func TestAdvanceRecurringClause(t *testing.T) {
	t.Parallel()

	jst, err := time.LoadLocation("Asia/Tokyo")
	require.NoError(t, err)

	tests := []struct {
		desc                   string
		clauseID               string
		datetimeClauses        []*autoopsproto.DatetimeClause
		now                    time.Time
		expectErr              bool
		expectExecutionCount   int32
		expectNextExecutionAt  func(t *testing.T, nextExec int64)
		expectClauseExecutedAt int64
	}{
		{
			desc:     "success: advance weekly recurring clause",
			clauseID: "clause-1",
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000, // 10:00 AM
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1}, // Monday
						Timezone:   "Asia/Tokyo",
						StartDate:  time.Date(2026, 2, 9, 0, 0, 0, 0, jst).Unix(),
					},
					NextExecutionAt: time.Date(2026, 2, 9, 10, 0, 0, 0, jst).Unix(),
					ExecutionCount:  0,
					LastExecutedAt:  0,
				},
			},
			now:                  time.Date(2026, 2, 9, 10, 0, 1, 0, jst),
			expectErr:            false,
			expectExecutionCount: 1,
			expectNextExecutionAt: func(t *testing.T, nextExec int64) {
				expected := time.Date(2026, 2, 16, 10, 0, 0, 0, jst).Unix()
				assert.Equal(t, expected, nextExec, "next execution should be next Monday")
			},
			expectClauseExecutedAt: 0, // Still active
		},
		{
			desc:     "success: advance daily recurring clause",
			clauseID: "clause-1",
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{
					Time:       32400, // 9:00 AM
					ActionType: autoopsproto.ActionType_DISABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency: autoopsproto.RecurrenceRule_DAILY,
						Timezone:  "UTC",
						StartDate: time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC).Unix(),
					},
					NextExecutionAt: time.Date(2026, 3, 1, 9, 0, 0, 0, time.UTC).Unix(),
					ExecutionCount:  0,
				},
			},
			now:                  time.Date(2026, 3, 1, 9, 0, 5, 0, time.UTC),
			expectErr:            false,
			expectExecutionCount: 1,
			expectNextExecutionAt: func(t *testing.T, nextExec int64) {
				expected := time.Date(2026, 3, 2, 9, 0, 0, 0, time.UTC).Unix()
				assert.Equal(t, expected, nextExec)
			},
			expectClauseExecutedAt: 0,
		},
		{
			desc:     "success: exhaust max occurrences",
			clauseID: "clause-1",
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:      autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek:     []int32{1},
						Timezone:       "UTC",
						StartDate:      time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix(),
						MaxOccurrences: 3,
					},
					NextExecutionAt: time.Date(2026, 1, 19, 10, 0, 0, 0, time.UTC).Unix(),
					ExecutionCount:  2, // Will become 3, reaching max
				},
			},
			now:                  time.Date(2026, 1, 19, 10, 0, 1, 0, time.UTC),
			expectErr:            false,
			expectExecutionCount: 3,
			expectNextExecutionAt: func(t *testing.T, nextExec int64) {
				assert.Equal(t, int64(0), nextExec, "should be exhausted")
			},
			expectClauseExecutedAt: time.Date(2026, 1, 19, 10, 0, 1, 0, time.UTC).Unix(),
		},
		{
			desc:     "success: exhaust by end date",
			clauseID: "clause-1",
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix(),
						EndDate:    time.Date(2026, 1, 20, 0, 0, 0, 0, time.UTC).Unix(),
					},
					NextExecutionAt: time.Date(2026, 1, 19, 10, 0, 0, 0, time.UTC).Unix(),
					ExecutionCount:  1,
				},
			},
			now:                  time.Date(2026, 1, 19, 10, 0, 1, 0, time.UTC),
			expectErr:            false,
			expectExecutionCount: 2,
			expectNextExecutionAt: func(t *testing.T, nextExec int64) {
				assert.Equal(t, int64(0), nextExec, "should be exhausted past end date")
			},
			expectClauseExecutedAt: time.Date(2026, 1, 19, 10, 0, 1, 0, time.UTC).Unix(),
		},
		{
			desc:     "error: clause not found",
			clauseID: "nonexistent",
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{
					Time:       36000,
					ActionType: autoopsproto.ActionType_ENABLE,
					Recurrence: &autoopsproto.RecurrenceRule{
						Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
						DaysOfWeek: []int32{1},
						Timezone:   "UTC",
						StartDate:  time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix(),
					},
					NextExecutionAt: time.Now().Unix(),
					ExecutionCount:  0,
				},
			},
			now:       time.Now(),
			expectErr: true,
		},
		{
			desc:     "error: clause is not recurring",
			clauseID: "clause-1",
			datetimeClauses: []*autoopsproto.DatetimeClause{
				{
					Time:       time.Now().Add(24 * time.Hour).Unix(),
					ActionType: autoopsproto.ActionType_ENABLE,
				},
			},
			now:       time.Now(),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			rule, e := NewAutoOpsRule(
				"feature-1",
				autoopsproto.OpsType_SCHEDULE,
				nil,
				tt.datetimeClauses,
			)
			require.NoError(t, e)

			clauseID := tt.clauseID
			if clauseID == "clause-1" && len(rule.Clauses) > 0 {
				clauseID = rule.Clauses[0].Id
			}

			e = rule.AdvanceRecurringClause(clauseID, tt.now)

			if tt.expectErr {
				assert.Error(t, e)
				return
			}
			require.NoError(t, e)

			dateClauses, e := rule.ExtractDatetimeClauses()
			require.NoError(t, e)

			clause := rule.Clauses[0]
			dtClause := dateClauses[clause.Id]

			assert.Equal(t, tt.expectExecutionCount, dtClause.ExecutionCount)
			assert.Equal(t, tt.now.Unix(), dtClause.LastExecutedAt)
			assert.Equal(t, tt.expectClauseExecutedAt, clause.ExecutedAt)

			if tt.expectNextExecutionAt != nil {
				tt.expectNextExecutionAt(t, dtClause.NextExecutionAt)
			}
		})
	}
}

func TestAdvanceRecurringClause_MultipleClauses(t *testing.T) {
	t.Parallel()

	enableClause := &autoopsproto.DatetimeClause{
		Time:       36000,
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "UTC",
			StartDate:  time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}
	disableClause := &autoopsproto.DatetimeClause{
		Time:       64800,
		ActionType: autoopsproto.ActionType_DISABLE,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency:  autoopsproto.RecurrenceRule_WEEKLY,
			DaysOfWeek: []int32{1},
			Timezone:   "UTC",
			StartDate:  time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}

	rule, err := NewAutoOpsRule(
		"feature-1",
		autoopsproto.OpsType_SCHEDULE,
		nil,
		[]*autoopsproto.DatetimeClause{enableClause, disableClause},
	)
	require.NoError(t, err)
	require.Len(t, rule.Clauses, 2)

	now := time.Date(2026, 1, 5, 10, 0, 1, 0, time.UTC)
	err = rule.AdvanceRecurringClause(rule.Clauses[0].Id, now)
	require.NoError(t, err)

	assert.False(t, rule.AllClausesFinished(), "rule should still have active clauses")

	dateClauses, err := rule.ExtractDatetimeClauses()
	require.NoError(t, err)

	dtFirst := dateClauses[rule.Clauses[0].Id]
	assert.Equal(t, int32(1), dtFirst.ExecutionCount)
	assert.True(t, dtFirst.NextExecutionAt > 0, "should have a next execution")

	dtSecond := dateClauses[rule.Clauses[1].Id]
	assert.Equal(t, int32(0), dtSecond.ExecutionCount)
}

// TestAdvanceRecurringClause_UpdateAndJSONRoundTrip reproduces the exact
// ExecuteAutoOps flow: AdvanceRecurringClause → Update() (copier.Copy) →
// JSON marshal/unmarshal (storage layer). This identifies whether copier.Copy
// or the JSON round-trip loses the re-marshaled Any bytes.
func TestAdvanceRecurringClause_UpdateAndJSONRoundTrip(t *testing.T) {
	t.Parallel()

	dc := &autoopsproto.DatetimeClause{
		Time:       36000,
		ActionType: autoopsproto.ActionType_ENABLE,
		Recurrence: &autoopsproto.RecurrenceRule{
			Frequency: autoopsproto.RecurrenceRule_DAILY,
			Timezone:  "UTC",
			StartDate: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		},
	}
	rule, err := NewAutoOpsRule(
		"feature-1",
		autoopsproto.OpsType_SCHEDULE,
		nil,
		[]*autoopsproto.DatetimeClause{dc},
	)
	require.NoError(t, err)
	require.Len(t, rule.Clauses, 1)
	assert.True(t, rule.Clauses[0].IsRecurring)

	now := time.Date(2026, 1, 1, 10, 0, 1, 0, time.UTC)
	err = rule.AdvanceRecurringClause(rule.Clauses[0].Id, now)
	require.NoError(t, err)

	// Step 1: Verify the in-memory state is correct after AdvanceRecurringClause
	dtBefore, err := rule.ExtractDatetimeClauses()
	require.NoError(t, err)
	dtc := dtBefore[rule.Clauses[0].Id]
	assert.Equal(t, int32(1), dtc.ExecutionCount, "in-memory: ExecutionCount should be 1")
	assert.Equal(t, now.Unix(), dtc.LastExecutedAt, "in-memory: LastExecutedAt should be set")
	assert.True(t, dtc.NextExecutionAt > 0, "in-memory: NextExecutionAt should be set")

	// Step 2: Call Update() which uses copier.Copy — this is what ExecuteAutoOps does
	opsStatus := autoopsproto.AutoOpsStatus_RUNNING
	updated, err := rule.Update(&opsStatus, nil, nil)
	require.NoError(t, err)

	// Step 2a: Check the copier.Copy result
	dtAfterCopy, err := updated.ExtractDatetimeClauses()
	require.NoError(t, err)
	copiedDtc := dtAfterCopy[updated.Clauses[0].Id]
	require.NotNil(t, copiedDtc, "copier.Copy: DatetimeClause should be extractable")
	assert.Equal(t, int32(1), copiedDtc.ExecutionCount, "copier.Copy: ExecutionCount should be 1")
	assert.Equal(t, now.Unix(), copiedDtc.LastExecutedAt, "copier.Copy: LastExecutedAt should be set")
	assert.Equal(t, dtc.NextExecutionAt, copiedDtc.NextExecutionAt, "copier.Copy: NextExecutionAt should match")
	assert.NotNil(t, copiedDtc.Recurrence, "copier.Copy: Recurrence should not be nil")
	assert.True(t, updated.Clauses[0].IsRecurring, "copier.Copy: IsRecurring should be true")

	// Step 3: JSON round-trip (what mysql.JSONObject does for storage)
	jsonBytes, err := json.Marshal(updated.Clauses)
	require.NoError(t, err)

	var fromStorage []*autoopsproto.Clause
	err = json.Unmarshal(jsonBytes, &fromStorage)
	require.NoError(t, err)
	require.Len(t, fromStorage, 1)

	storageRule := &AutoOpsRule{AutoOpsRule: &autoopsproto.AutoOpsRule{
		Clauses: fromStorage,
	}}
	dtAfterStorage, err := storageRule.ExtractDatetimeClauses()
	require.NoError(t, err)
	storedDtc := dtAfterStorage[fromStorage[0].Id]
	require.NotNil(t, storedDtc, "storage: DatetimeClause should be extractable")
	assert.Equal(t, int32(1), storedDtc.ExecutionCount, "storage: ExecutionCount should be 1")
	assert.Equal(t, now.Unix(), storedDtc.LastExecutedAt, "storage: LastExecutedAt should be set")
	assert.Equal(t, dtc.NextExecutionAt, storedDtc.NextExecutionAt, "storage: NextExecutionAt should match")
	assert.NotNil(t, storedDtc.Recurrence, "storage: Recurrence should not be nil")
	assert.True(t, fromStorage[0].IsRecurring, "storage: IsRecurring should be true")
}

// TestRecurringLifecycle_AllPatternsAndEndConditions tests every combination of
// frequency (Daily, Weekly, Monthly) × end condition (OnDate, After, Never)
// by simulating init → advance → advance → verify termination. Each case
// executes at least twice before ending.
func TestRecurringLifecycle_AllPatternsAndEndConditions(t *testing.T) {
	t.Parallel()

	type lifecycleCase struct {
		desc           string
		frequency      autoopsproto.RecurrenceRule_Frequency
		daysOfWeek     []int32
		dayOfMonth     int32
		timeOfDay      int64
		startDate      time.Time
		endDate        int64
		maxOccurrences int32
		firstExecAt    time.Time
		secondExecAt   time.Time
		expectThirdRun bool
	}

	cases := []lifecycleCase{
		// ===== DAILY =====
		{
			desc:           "Daily + OnDate: runs 2 days, stops on end date",
			frequency:      autoopsproto.RecurrenceRule_DAILY,
			timeOfDay:      36000, // 10:00
			startDate:      time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2026, 3, 2, 23, 59, 59, 0, time.UTC).Unix(),
			firstExecAt:    time.Date(2026, 3, 1, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 3, 2, 10, 0, 1, 0, time.UTC),
			expectThirdRun: false,
		},
		{
			desc:           "Daily + After 2: runs exactly twice",
			frequency:      autoopsproto.RecurrenceRule_DAILY,
			timeOfDay:      36000,
			startDate:      time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			maxOccurrences: 2,
			firstExecAt:    time.Date(2026, 3, 1, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 3, 2, 10, 0, 1, 0, time.UTC),
			expectThirdRun: false,
		},
		{
			desc:           "Daily + Never: runs 2 times, continues",
			frequency:      autoopsproto.RecurrenceRule_DAILY,
			timeOfDay:      36000,
			startDate:      time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			firstExecAt:    time.Date(2026, 3, 1, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 3, 2, 10, 0, 1, 0, time.UTC),
			expectThirdRun: true,
		},

		// ===== WEEKLY =====
		{
			desc:           "Weekly + OnDate: runs 2 Mondays, stops on end date",
			frequency:      autoopsproto.RecurrenceRule_WEEKLY,
			daysOfWeek:     []int32{1}, // Monday
			timeOfDay:      36000,
			startDate:      time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC), // Monday
			endDate:        time.Date(2026, 3, 15, 23, 59, 59, 0, time.UTC).Unix(),
			firstExecAt:    time.Date(2026, 3, 2, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 3, 9, 10, 0, 1, 0, time.UTC),
			expectThirdRun: false,
		},
		{
			desc:           "Weekly + After 2: runs exactly 2 Wednesdays",
			frequency:      autoopsproto.RecurrenceRule_WEEKLY,
			daysOfWeek:     []int32{3}, // Wednesday
			timeOfDay:      36000,
			startDate:      time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC), // Wednesday
			maxOccurrences: 2,
			firstExecAt:    time.Date(2026, 3, 4, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 3, 11, 10, 0, 1, 0, time.UTC),
			expectThirdRun: false,
		},
		{
			desc:           "Weekly + Never: runs 2 Fridays, continues",
			frequency:      autoopsproto.RecurrenceRule_WEEKLY,
			daysOfWeek:     []int32{5}, // Friday
			timeOfDay:      36000,
			startDate:      time.Date(2026, 3, 6, 0, 0, 0, 0, time.UTC), // Friday
			firstExecAt:    time.Date(2026, 3, 6, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 3, 13, 10, 0, 1, 0, time.UTC),
			expectThirdRun: true,
		},

		// ===== MONTHLY =====
		{
			desc:           "Monthly + OnDate: runs 2 months on the 15th, stops on end date",
			frequency:      autoopsproto.RecurrenceRule_MONTHLY,
			dayOfMonth:     15,
			timeOfDay:      36000,
			startDate:      time.Date(2026, 1, 15, 0, 0, 0, 0, time.UTC),
			endDate:        time.Date(2026, 3, 14, 23, 59, 59, 0, time.UTC).Unix(),
			firstExecAt:    time.Date(2026, 1, 15, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 2, 15, 10, 0, 1, 0, time.UTC),
			expectThirdRun: false,
		},
		{
			desc:           "Monthly + After 2: runs exactly twice on the 1st",
			frequency:      autoopsproto.RecurrenceRule_MONTHLY,
			dayOfMonth:     1,
			timeOfDay:      36000,
			startDate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			maxOccurrences: 2,
			firstExecAt:    time.Date(2026, 1, 1, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 2, 1, 10, 0, 1, 0, time.UTC),
			expectThirdRun: false,
		},
		{
			desc:           "Monthly + Never: runs 2 months on the 10th, continues",
			frequency:      autoopsproto.RecurrenceRule_MONTHLY,
			dayOfMonth:     10,
			timeOfDay:      36000,
			startDate:      time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC),
			firstExecAt:    time.Date(2026, 1, 10, 10, 0, 1, 0, time.UTC),
			secondExecAt:   time.Date(2026, 2, 10, 10, 0, 1, 0, time.UTC),
			expectThirdRun: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			dc := &autoopsproto.DatetimeClause{
				Time:       tc.timeOfDay,
				ActionType: autoopsproto.ActionType_ENABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:      tc.frequency,
					DaysOfWeek:     tc.daysOfWeek,
					DayOfMonth:     tc.dayOfMonth,
					Timezone:       "UTC",
					StartDate:      tc.startDate.Unix(),
					EndDate:        tc.endDate,
					MaxOccurrences: tc.maxOccurrences,
				},
			}

			rule, err := NewAutoOpsRule(
				"feature-1",
				autoopsproto.OpsType_SCHEDULE,
				nil,
				[]*autoopsproto.DatetimeClause{dc},
			)
			require.NoError(t, err)
			require.Len(t, rule.Clauses, 1)
			assert.True(t, rule.Clauses[0].IsRecurring)

			clauseID := rule.Clauses[0].Id

			// Verify initialization
			dtClauses, err := rule.ExtractDatetimeClauses()
			require.NoError(t, err)
			initDtc := dtClauses[clauseID]
			require.NotNil(t, initDtc)
			assert.Equal(t, int32(0), initDtc.ExecutionCount, "init: ExecutionCount should be 0")
			assert.True(t, initDtc.NextExecutionAt > 0, "init: NextExecutionAt should be set")
			assert.False(t, rule.AllClausesFinished(), "init: should not be finished")

			// --- First execution ---
			err = rule.AdvanceRecurringClause(clauseID, tc.firstExecAt)
			require.NoError(t, err)

			dtClauses, err = rule.ExtractDatetimeClauses()
			require.NoError(t, err)
			dtc1 := dtClauses[clauseID]
			assert.Equal(t, int32(1), dtc1.ExecutionCount, "exec1: ExecutionCount should be 1")
			assert.Equal(t, tc.firstExecAt.Unix(), dtc1.LastExecutedAt, "exec1: LastExecutedAt should match")
			assert.True(t, dtc1.NextExecutionAt > 0, "exec1: should have NextExecutionAt for second run")
			assert.False(t, rule.AllClausesFinished(), "exec1: should not be finished yet")

			// --- Second execution ---
			err = rule.AdvanceRecurringClause(clauseID, tc.secondExecAt)
			require.NoError(t, err)

			dtClauses, err = rule.ExtractDatetimeClauses()
			require.NoError(t, err)
			dtc2 := dtClauses[clauseID]
			assert.Equal(t, int32(2), dtc2.ExecutionCount, "exec2: ExecutionCount should be 2")
			assert.Equal(t, tc.secondExecAt.Unix(), dtc2.LastExecutedAt, "exec2: LastExecutedAt should match")

			if tc.expectThirdRun {
				assert.True(t, dtc2.NextExecutionAt > 0, "exec2: should have NextExecutionAt (never-ending)")
				assert.False(t, rule.AllClausesFinished(), "exec2: should NOT be finished (never-ending)")
			} else {
				assert.Equal(t, int64(0), dtc2.NextExecutionAt, "exec2: NextExecutionAt should be 0 (exhausted)")
				assert.True(t, rule.AllClausesFinished(), "exec2: should be finished")
				assert.True(t, rule.Clauses[0].ExecutedAt > 0, "exec2: outer Clause.ExecutedAt should be set")
			}
		})
	}
}

// TestRecurringLifecycle_MultiClause_FinishedStatus verifies that a rule with
// multiple recurring clauses (Enable + Disable) only transitions to FINISHED
// after ALL clauses are exhausted, for each frequency type.
func TestRecurringLifecycle_MultiClause_FinishedStatus(t *testing.T) {
	t.Parallel()

	type execStep struct {
		clauseIdx int
		execTime  time.Time
	}

	type multiClauseCase struct {
		desc           string
		frequency      autoopsproto.RecurrenceRule_Frequency
		daysOfWeek     []int32
		dayOfMonth     int32
		enableTime     int64
		disableTime    int64
		startDate      time.Time
		maxOccurrences int32
		steps          []execStep
	}

	cases := []multiClauseCase{
		{
			desc:           "Daily: Enable 10:00 + Disable 18:00, max 2",
			frequency:      autoopsproto.RecurrenceRule_DAILY,
			enableTime:     36000,
			disableTime:    64800,
			startDate:      time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC),
			maxOccurrences: 2,
			steps: []execStep{
				{0, time.Date(2026, 3, 1, 10, 0, 1, 0, time.UTC)},
				{1, time.Date(2026, 3, 1, 18, 0, 1, 0, time.UTC)},
				{0, time.Date(2026, 3, 2, 10, 0, 1, 0, time.UTC)},
				{1, time.Date(2026, 3, 2, 18, 0, 1, 0, time.UTC)},
			},
		},
		{
			desc:           "Weekly (Mon): Enable 10:00 + Disable 18:00, max 2",
			frequency:      autoopsproto.RecurrenceRule_WEEKLY,
			daysOfWeek:     []int32{1},
			enableTime:     36000,
			disableTime:    64800,
			startDate:      time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC), // Monday
			maxOccurrences: 2,
			steps: []execStep{
				{0, time.Date(2026, 3, 2, 10, 0, 1, 0, time.UTC)},
				{1, time.Date(2026, 3, 2, 18, 0, 1, 0, time.UTC)},
				{0, time.Date(2026, 3, 9, 10, 0, 1, 0, time.UTC)},
				{1, time.Date(2026, 3, 9, 18, 0, 1, 0, time.UTC)},
			},
		},
		{
			desc:           "Monthly (1st): Enable 10:00 + Disable 18:00, max 2",
			frequency:      autoopsproto.RecurrenceRule_MONTHLY,
			dayOfMonth:     1,
			enableTime:     36000,
			disableTime:    64800,
			startDate:      time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
			maxOccurrences: 2,
			steps: []execStep{
				{0, time.Date(2026, 1, 1, 10, 0, 1, 0, time.UTC)},
				{1, time.Date(2026, 1, 1, 18, 0, 1, 0, time.UTC)},
				{0, time.Date(2026, 2, 1, 10, 0, 1, 0, time.UTC)},
				{1, time.Date(2026, 2, 1, 18, 0, 1, 0, time.UTC)},
			},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()

			enableClause := &autoopsproto.DatetimeClause{
				Time:       tc.enableTime,
				ActionType: autoopsproto.ActionType_ENABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:      tc.frequency,
					DaysOfWeek:     tc.daysOfWeek,
					DayOfMonth:     tc.dayOfMonth,
					Timezone:       "UTC",
					StartDate:      tc.startDate.Unix(),
					MaxOccurrences: tc.maxOccurrences,
				},
			}
			disableClause := &autoopsproto.DatetimeClause{
				Time:       tc.disableTime,
				ActionType: autoopsproto.ActionType_DISABLE,
				Recurrence: &autoopsproto.RecurrenceRule{
					Frequency:      tc.frequency,
					DaysOfWeek:     tc.daysOfWeek,
					DayOfMonth:     tc.dayOfMonth,
					Timezone:       "UTC",
					StartDate:      tc.startDate.Unix(),
					MaxOccurrences: tc.maxOccurrences,
				},
			}

			rule, err := NewAutoOpsRule(
				"feature-1",
				autoopsproto.OpsType_SCHEDULE,
				nil,
				[]*autoopsproto.DatetimeClause{enableClause, disableClause},
			)
			require.NoError(t, err)
			require.Len(t, rule.Clauses, 2)

			determineStatus := func() autoopsproto.AutoOpsStatus {
				if rule.AllClausesFinished() {
					return autoopsproto.AutoOpsStatus_FINISHED
				}
				return autoopsproto.AutoOpsStatus_RUNNING
			}

			lastStepIdx := len(tc.steps) - 1
			for i, step := range tc.steps {
				clauseID := rule.Clauses[step.clauseIdx].Id
				err = rule.AdvanceRecurringClause(clauseID, step.execTime)
				require.NoError(t, err, "step %d failed", i)

				if i < lastStepIdx {
					assert.Equal(t, autoopsproto.AutoOpsStatus_RUNNING, determineStatus(),
						"step %d: should be RUNNING (not all clauses exhausted yet)", i)
				} else {
					assert.Equal(t, autoopsproto.AutoOpsStatus_FINISHED, determineStatus(),
						"step %d: should be FINISHED (all clauses exhausted)", i)
				}
			}

			// Verify final state
			dtClauses, err := rule.ExtractDatetimeClauses()
			require.NoError(t, err)
			for i, c := range rule.Clauses {
				dtc := dtClauses[c.Id]
				assert.Equal(t, int64(0), dtc.NextExecutionAt,
					"clause %d: NextExecutionAt should be 0", i)
				assert.Equal(t, tc.maxOccurrences, dtc.ExecutionCount,
					"clause %d: ExecutionCount should be %d", i, tc.maxOccurrences)
				assert.True(t, c.ExecutedAt > 0,
					"clause %d: outer ExecutedAt should be set", i)
			}

			// Verify Update() produces FINISHED status
			status := autoopsproto.AutoOpsStatus_FINISHED
			updated, err := rule.Update(&status, nil, nil)
			require.NoError(t, err)
			assert.Equal(t, autoopsproto.AutoOpsStatus_FINISHED, updated.AutoOpsStatus)
		})
	}
}
