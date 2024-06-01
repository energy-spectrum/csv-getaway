package helper

import "time"

type CampaignSchedule struct {
	StartTime time.Time
	EndTime   time.Time
}

func IsScheduleValid(schedule []CampaignSchedule) bool {
	for i, s1 := range schedule {
		for j, s2 := range schedule {
			if i != j && isTimeOverlap(s1.StartTime, s1.EndTime, s2.StartTime, s2.EndTime) {
				// найдено пересечение времени в расписании.
				return false
			}
		}
		// время начала и конца интервала равны
		if s1.StartTime.Equal(s1.EndTime) {
			return false
		}
	}
	// расписание валидно, нет пересечений времени.
	return true
}

func isTimeOverlap(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && end1.After(start2) || (start1.Equal(start2) && end1.Equal(end2))
}

func ValidateIntervals(intervals []CampaignSchedule) bool {
	for _, interval := range intervals {
		// TODO: добавить и переработать:  interval.StartTime.Before(currentTime)
		// interval.EndTime.Before(currentTime) ||
		if interval.EndTime.Before(interval.StartTime) {
			return false
		}
	}
	return true
}
