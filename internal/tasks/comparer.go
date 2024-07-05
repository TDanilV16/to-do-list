package tasks

import (
	"time"
)

func OrderByDeadlineDescending(task1, task2 *Task) int {
	return orderBy(task1.Deadline, task2.Deadline, false)
}

func OrderByDeadlineAscending(task1, task2 *Task) int {
	return orderBy(task1.Deadline, task2.Deadline, true)
}

func orderBy(time1, time2 time.Time, ascending bool) int {
	mul := 1
	if ascending {
		mul *= -1
	}
	if time1.Before(time2) {
		return mul
	} else if time1.After(time2) {
		return -mul
	}

	return 0
}
