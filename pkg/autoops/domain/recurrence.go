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
	"errors"
	"sort"
	"time"

	proto "github.com/bucketeer-io/bucketeer/v2/proto/autoops"
)

// CalculateNextExecution calculates the next execution time for a recurring clause.
// Returns (nextExecutionUnix, shouldContinue).
// nextExecutionUnix is 0 when no more executions remain.
func CalculateNextExecution(
	clause *proto.DatetimeClause,
	currentExecutionTime time.Time,
) (int64, bool) {
	if clause == nil {
		return 0, false
	}

	recurrence := clause.Recurrence
	if recurrence == nil || recurrence.Frequency == proto.RecurrenceRule_ONCE ||
		recurrence.Frequency == proto.RecurrenceRule_FREQUENCY_UNSPECIFIED {
		return 0, false
	}

	if clause.Time < 0 || clause.Time >= 24*60*60 {
		return 0, false
	}

	if recurrence.MaxOccurrences > 0 && clause.ExecutionCount >= recurrence.MaxOccurrences {
		return 0, false
	}

	if recurrence.EndDate > 0 && currentExecutionTime.Unix() >= recurrence.EndDate {
		return 0, false
	}

	loc := loadTimezone(recurrence.Timezone)
	currentTime := currentExecutionTime.In(loc)

	var nextExec time.Time
	var found bool

	switch recurrence.Frequency {
	case proto.RecurrenceRule_DAILY:
		nextExec, found = calculateNextDaily(clause, currentTime, loc)
	case proto.RecurrenceRule_WEEKLY:
		nextExec, found = calculateNextWeekly(clause, currentTime, loc)
	case proto.RecurrenceRule_MONTHLY:
		nextExec, found = calculateNextMonthly(clause, currentTime, loc)
	default:
		return 0, false
	}

	if !found {
		return 0, false
	}

	if recurrence.EndDate > 0 && nextExec.Unix() >= recurrence.EndDate {
		return 0, false
	}

	return nextExec.Unix(), true
}

// InitializeRecurringClause sets up initial tracking fields for a new recurring clause.
func InitializeRecurringClause(clause *proto.DatetimeClause) error {
	if clause.Recurrence == nil || clause.Recurrence.Frequency == proto.RecurrenceRule_ONCE ||
		clause.Recurrence.Frequency == proto.RecurrenceRule_FREQUENCY_UNSPECIFIED {
		return nil
	}

	if clause.Time < 0 || clause.Time >= 24*60*60 {
		return errors.New("time must be seconds since midnight in range [0, 86399]")
	}

	if clause.Recurrence.StartDate == 0 {
		return errors.New("start_date is required for recurring schedules")
	}

	switch clause.Recurrence.Frequency {
	case proto.RecurrenceRule_WEEKLY:
		if _, ok := validateDaysOfWeek(clause.Recurrence.DaysOfWeek); !ok {
			return errors.New("days_of_week must be non-empty with values 0-6 for weekly recurrence")
		}
	case proto.RecurrenceRule_MONTHLY:
		if clause.Recurrence.DayOfMonth < 1 || clause.Recurrence.DayOfMonth > 31 {
			return errors.New("day_of_month must be 1-31 for monthly recurrence")
		}
	}

	loc := loadTimezone(clause.Recurrence.Timezone)
	startTime := time.Unix(clause.Recurrence.StartDate, 0).In(loc)

	nextExec := computeFirstExecution(clause, startTime, loc)

	if clause.Recurrence.EndDate > 0 && nextExec.Unix() >= clause.Recurrence.EndDate {
		clause.NextExecutionAt = 0
	} else {
		clause.NextExecutionAt = nextExec.Unix()
	}
	clause.ExecutionCount = 0
	clause.LastExecutedAt = 0

	return nil
}

// IsRecurring returns true if the DatetimeClause has a recurring schedule.
func IsRecurring(clause *proto.DatetimeClause) bool {
	if clause == nil {
		return false
	}
	return clause.Recurrence != nil &&
		clause.Recurrence.Frequency != proto.RecurrenceRule_ONCE &&
		clause.Recurrence.Frequency != proto.RecurrenceRule_FREQUENCY_UNSPECIFIED
}

func loadTimezone(tz string) *time.Location {
	if tz == "" {
		return time.UTC
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.UTC
	}
	return loc
}

func timeOfDayComponents(clause *proto.DatetimeClause) (int, int, int) {
	tod := clause.Time
	hours := int(tod / 3600)
	minutes := int((tod % 3600) / 60)
	seconds := int(tod % 60)
	return hours, minutes, seconds
}

func calculateNextDaily(
	clause *proto.DatetimeClause,
	currentTime time.Time,
	loc *time.Location,
) (time.Time, bool) {
	hours, minutes, seconds := timeOfDayComponents(clause)

	nextExec := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day()+1,
		hours, minutes, seconds,
		0, loc,
	)

	return nextExec, true
}

func validateDaysOfWeek(daysOfWeek []int32) ([]int, bool) {
	if len(daysOfWeek) == 0 {
		return nil, false
	}
	sortedDays := make([]int, len(daysOfWeek))
	for i, d := range daysOfWeek {
		day := int(d)
		if day < 0 || day > 6 {
			return nil, false
		}
		sortedDays[i] = day
	}
	sort.Ints(sortedDays)
	return sortedDays, true
}

func calculateNextWeekly(
	clause *proto.DatetimeClause,
	currentTime time.Time,
	loc *time.Location,
) (time.Time, bool) {
	sortedDays, ok := validateDaysOfWeek(clause.Recurrence.DaysOfWeek)
	if !ok {
		return time.Time{}, false
	}

	hours, minutes, seconds := timeOfDayComponents(clause)
	currentWeekday := int(currentTime.Weekday())

	// Check remaining days this week (strictly after today)
	for _, day := range sortedDays {
		if day > currentWeekday {
			daysUntil := day - currentWeekday
			nextExec := time.Date(
				currentTime.Year(),
				currentTime.Month(),
				currentTime.Day()+daysUntil,
				hours, minutes, seconds,
				0, loc,
			)
			return nextExec, true
		}
	}

	// Wrap to the first scheduled day of next week
	firstDay := sortedDays[0]
	daysUntil := (7 - currentWeekday) + firstDay
	nextExec := time.Date(
		currentTime.Year(),
		currentTime.Month(),
		currentTime.Day()+daysUntil,
		hours, minutes, seconds,
		0, loc,
	)

	return nextExec, true
}

func calculateNextMonthly(
	clause *proto.DatetimeClause,
	currentTime time.Time,
	loc *time.Location,
) (time.Time, bool) {
	dayOfMonth := clause.Recurrence.DayOfMonth
	if dayOfMonth < 1 || dayOfMonth > 31 {
		return time.Time{}, false
	}

	hours, minutes, seconds := timeOfDayComponents(clause)

	nextMonth := currentTime.AddDate(0, 1, 0)
	nextExec := time.Date(
		nextMonth.Year(),
		nextMonth.Month(),
		int(dayOfMonth),
		hours, minutes, seconds,
		0, loc,
	)

	// time.Date auto-normalizes overflows (e.g., Feb 30 -> Mar 2).
	// If the month shifted, skip to the following month.
	if nextExec.Month() != nextMonth.Month() {
		nextMonth = nextMonth.AddDate(0, 1, 0)
		nextExec = time.Date(
			nextMonth.Year(),
			nextMonth.Month(),
			int(dayOfMonth),
			hours, minutes, seconds,
			0, loc,
		)
	}

	return nextExec, true
}

// computeFirstExecution determines the first execution timestamp for a newly
// created recurring clause based on its start_date and recurrence pattern.
func computeFirstExecution(
	clause *proto.DatetimeClause,
	startTime time.Time,
	loc *time.Location,
) time.Time {
	hours, minutes, seconds := timeOfDayComponents(clause)

	switch clause.Recurrence.Frequency {
	case proto.RecurrenceRule_DAILY:
		candidate := time.Date(
			startTime.Year(), startTime.Month(), startTime.Day(),
			hours, minutes, seconds, 0, loc,
		)
		if candidate.Before(startTime) {
			candidate = candidate.AddDate(0, 0, 1)
		}
		return candidate

	case proto.RecurrenceRule_WEEKLY:
		sortedDays, ok := validateDaysOfWeek(clause.Recurrence.DaysOfWeek)
		if !ok {
			return startTime
		}

		startWeekday := int(startTime.Weekday())
		for _, day := range sortedDays {
			if day < startWeekday {
				continue
			}
			daysUntil := day - startWeekday
			candidate := time.Date(
				startTime.Year(), startTime.Month(), startTime.Day()+daysUntil,
				hours, minutes, seconds, 0, loc,
			)
			if !candidate.Before(startTime) {
				return candidate
			}
		}
		// Wrap to next week
		firstDay := sortedDays[0]
		daysUntil := (7 - startWeekday) + firstDay
		return time.Date(
			startTime.Year(), startTime.Month(), startTime.Day()+daysUntil,
			hours, minutes, seconds, 0, loc,
		)

	case proto.RecurrenceRule_MONTHLY:
		dom := int(clause.Recurrence.DayOfMonth)
		if dom < 1 || dom > 31 {
			return startTime
		}
		candidate := time.Date(
			startTime.Year(), startTime.Month(), dom,
			hours, minutes, seconds, 0, loc,
		)
		if !candidate.Before(startTime) && candidate.Month() == startTime.Month() {
			return candidate
		}
		// Move to next month
		nextMonth := startTime.AddDate(0, 1, 0)
		candidate = time.Date(
			nextMonth.Year(), nextMonth.Month(), dom,
			hours, minutes, seconds, 0, loc,
		)
		if candidate.Month() != nextMonth.Month() {
			nextMonth = nextMonth.AddDate(0, 1, 0)
			candidate = time.Date(
				nextMonth.Year(), nextMonth.Month(), dom,
				hours, minutes, seconds, 0, loc,
			)
		}
		return candidate

	default:
		return startTime
	}
}
