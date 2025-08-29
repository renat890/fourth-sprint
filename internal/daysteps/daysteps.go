package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

var (
	ErrInvalidLen = errors.New("некорректная длина слайса, длина должна быть равна 2")
	ErrStepsZero  = errors.New("количество шагов должно быть больше 0")
	ErrDurZero    = errors.New("продолжительность должна быть больше 0")
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	stepsDur := strings.Split(data, ",")

	if len(stepsDur) != 2 {
		return 0, 0, ErrInvalidLen
	}

	steps, err := strconv.Atoi(stepsDur[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, ErrStepsZero
	}

	dur, err := time.ParseDuration(stepsDur[1])
	if err != nil {
		return 0, 0, err
	}
	if dur <= 0 {
		return 0, 0, ErrDurZero
	}

	return steps, dur, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, dur, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceInM := stepLength * float64(steps)
	distanceInKm := distanceInM / mInKm
	calories, err := spentcalories.WalkingSpentCalories(steps, weight,height, dur)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceInKm, calories)
}
