package app

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"modernc.org/libc/limits"
)

// При наличии правила повторения, возвращает следующую дату текущей задачи.
// `now` - время от которого ищется ближайшая дата;
// `date` - текущее время выполнения задачи;
// `repeat` - правило повторения в спецаильном формате
func NextDate(now time.Time, date string, repeat string) (string, error) {
	beginDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("nextDate: invalid date format: <%s>, %w", date, err)
	}

	repeatSlice := strings.Split(repeat, " ")
	if len(repeatSlice) < 1 {
		return "", fmt.Errorf("nextDate: invalid repeat format: <%s>, repeat is empty", repeat)
	}

	// Разбиваем на модификатор и значение
	modif := repeatSlice[0]

	switch modif {

	case "y":
		if len(repeatSlice) != 1 {
			return "", fmt.Errorf("nextDate: invalid repeat format: [%s], a year cant have additional values", repeat)
		}

		next := beginDate.AddDate(1, 0, 0)
		for ; next.Before(now); next = next.AddDate(1, 0, 0) {

		}
		return next.Format("20060102"), nil

	case "d":
		if len(repeatSlice) != 2 {
			return "", fmt.Errorf("nextDate: invalid repeat format: [%s], days must have only one additional value", repeat)
		}

		days, err := strconv.ParseInt(repeatSlice[1], 10, 32)
		if err != nil {
			return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
		}

		if days > 400 {
			return "", fmt.Errorf("nextDate: invalid repeat format: [%s], days must be less than 400", repeat)
		}

		next := beginDate.AddDate(0, 0, int(days))
		for ; next.Before(now); next = next.AddDate(0, 0, int(days)) {

		}
		return next.Format("20060102"), nil

	case "w":
		if len(repeatSlice) < 2 {
			return "", fmt.Errorf("nextDate: invalid repeat format: [%s], a weekday must have one or more additional value", repeat)
		}
		weekDaysStringList := strings.Split(repeatSlice[1], ",")

		// Ищем ближайший день недели
		minDif := int64(limits.LONG_MAX)
		closestDate := time.Time{}
		for _, ds := range weekDaysStringList {
			weekDay, err := strconv.ParseInt(ds, 10, 32)
			if err != nil {
				return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
			}

			dt := now
			if beginDate.After(now) {
				dt = beginDate
			}
			d, err := nextWeekDay(dt, int(weekDay))
			if err != nil {
				return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
			}

			dif := d.Sub(now).Milliseconds()
			if dif < minDif {
				minDif = dif
				closestDate = d
			}
		}

		return closestDate.Format("20060102"), nil

	case "m":
		if len(repeatSlice) < 2 {
			return "", fmt.Errorf("nextDate: invalid repeat format: [%s], a monthday must have one or more additional value", repeat)
		}

		dt := now
		if beginDate.After(now) {
			dt = beginDate
		}

		minDif := int64(limits.LONG_MAX)
		closestDate := time.Time{}
		daysStringList := strings.Split(repeatSlice[1], ",")

		if len(repeatSlice) == 2 {
			// Если месяц не указан (повторение каждый месяц)

			for _, ds := range daysStringList {
				day, err := strconv.ParseInt(ds, 10, 32)
				if err != nil {
					return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
				}

				if day < -31 || day > 31 {
					return "", fmt.Errorf("nextDate: invalid repeat format: [%s], monthDays must be between -31 and 31", repeat)
				}

				d, err := nextMonthDay(dt, int(day))
				if err != nil {
					return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
				}

				dif := d.Sub(dt).Milliseconds()
				if dif < minDif {
					minDif = dif
					closestDate = d
				}
			}
			return closestDate.Format("20060102"), nil
		} else if len(repeatSlice) == 3 {
			// Если указаны конкретные месяцы

			monthStringList := strings.Split(repeatSlice[2], ",")

			for _, ms := range monthStringList {
				month, err := strconv.ParseInt(ms, 10, 32)
				if err != nil {
					return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
				}
				for _, ds := range daysStringList {
					day, err := strconv.ParseInt(ds, 10, 32)
					if err != nil {
						return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
					}
					d, err := nextSpecifiedDay(dt, int(day), int(month))
					if err != nil {
						return "", fmt.Errorf("nextDate: invalid repeat format: [%s], %w", repeat, err)
					}
					dif := d.Sub(dt).Milliseconds()
					if dif < minDif {
						minDif = dif
						closestDate = d
					}
				}
			}
			return closestDate.Format("20060102"), nil
		}
	}
	return "", fmt.Errorf("nextDate: invalid repeat format: [%s], a modificator must be y,d,w, m", repeat)
}

func monthLength(m time.Month) int {
	return time.Date(2000, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// Возвращает ближайшую дату указанного дня месяца `monthDay`
func nextMonthDay(current time.Time, monthDay int) (time.Time, error) {
	if monthDay > 31 || monthDay < -2 || monthDay == 0 {
		return time.Time{}, fmt.Errorf("monthDay couldn't be equal to %d", monthDay)
	}
	month := current.Month()
	if monthDay < 0 {
		monthDay += 1
		month++
	}

	for i := 0; i < 4; i++ {
		date := time.Date(current.Year(), month+time.Month(i), monthDay, 0, 0, 0, 0, current.Location())

		// Проверка чтобы искомый день остался внутри выбранного месяца, если нет, то может быть следующий?
		if date.After(current) && ((monthDay > 0 && month+time.Month(i) == date.Month()) || (monthDay <= 0 && month+time.Month(i-1) == date.Month())) {
			return date, nil
		}
	}
	return time.Time{}, fmt.Errorf("nextMonthDay not found")
}

// Возвращает ближайшую дату указанного дня недели `weekDay`
func nextWeekDay(current time.Time, weekDay int) (time.Time, error) {
	if weekDay < 1 || weekDay > 7 {
		return time.Time{}, fmt.Errorf("weekDay couldn't be equal to %d", weekDay)
	}

	dif := (time.Weekday(weekDay) - current.Weekday())
	if dif <= 0 {
		dif += 7
	}
	return current.AddDate(0, 0, int(dif)), nil
}

// Возвращает ближайшую дату указанного месяца `month` и дня месяца `monthDay`
func nextSpecifiedDay(current time.Time, monthDay, month int) (time.Time, error) {
	if month > 12 || month < 1 {
		return time.Time{}, fmt.Errorf("month couldn't be equal to %d", month)
	}
	if monthDay > monthLength(time.Month(month)) || monthDay > monthLength(time.Month(month)) {
		return time.Time{}, fmt.Errorf("monthDay couldn't be equal to %d in %s", monthDay, time.Month(month).String())
	}

	if monthDay < 0 {
		monthDay += 1
		month++
	}

	date := time.Date(current.Year(), time.Month(month), monthDay, 0, 0, 0, 0, current.Location())
	if date.Before(current) {
		return date.AddDate(1, 0, 0), nil
	}
	return date, nil
}
