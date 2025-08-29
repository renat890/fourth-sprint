package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

var (
	ErrInvalidLen          = errors.New("некорректная длина слайса, длина должна быть равна 3")
	ErrStepsZero           = errors.New("количество шагов должно быть больше 0")
	ErrDurationZero        = errors.New("продолжительность должна быть больше 0")
	ErrUnknownTypeTraining = errors.New("неизвестный тип тренировки")
	ErrHeightZero          = errors.New("рост должен быть больше 0")
	ErrWeightZero          = errors.New("вес должен быть больше 0")
)

func parseTraining(data string) (int, string, time.Duration, error) {
	training := strings.Split(data, ",")
	if len(training) != 3 {
		return 0, "", 0, ErrInvalidLen
	}

	steps, err := strconv.Atoi(training[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, ErrStepsZero
	}

	dur, err := time.ParseDuration(training[2])
	if err != nil {
		return 0, "", 0, err
	}
	if dur <= 0 {
		return 0, "", 0, ErrDurationZero
	}

	return steps, training[1], dur, nil
}

func distance(steps int, height float64) float64 {
	stepLen := height * stepLengthCoefficient
	distanceInM := stepLen * float64(steps)
	distanceInKm := distanceInM / mInKm
	
	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	
	wayInKm := distance(steps, height)
	avgSpeed := wayInKm / duration.Hours()

	return avgSpeed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// steps, type, dur, err  := 
	steps, typeTraining, dur, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	switch typeTraining {
		case "Ходьба":
			wayInKm := distance(steps, height)
			avgSpeed := meanSpeed(steps, height, dur)
			calories, err := WalkingSpentCalories(steps, weight, height, dur)
			if err != nil {
				log.Println(err)
				return "", err
			}
			return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeTraining, dur.Hours(), wayInKm, avgSpeed, calories), nil
		case "Бег":
			wayInKm := distance(steps, height)
			avgSpeed := meanSpeed(steps, height, dur)
			calories, err := RunningSpentCalories(steps, weight, height, dur)
			if err != nil {
				log.Println(err)
				return "", err
			}
			return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeTraining, dur.Hours(), wayInKm, avgSpeed, calories), nil
		default:
			log.Println(err)
			return "", ErrUnknownTypeTraining
	}
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrStepsZero
	}
	if duration <= 0 {
		return 0, ErrDurationZero
	}
	if height <= 0 {
		return 0, ErrHeightZero
	}
	if weight <= 0 {
		return 0, ErrWeightZero
	}
	
	avgSpeed := meanSpeed(steps, height, duration)
	calories := (weight * avgSpeed * duration.Minutes()) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, ErrStepsZero
	}
	if duration <= 0 {
		return 0, ErrDurationZero
	}
	if height <= 0 {
		return 0, ErrHeightZero
	}
	if weight <= 0 {
		return 0, ErrWeightZero
	}

	avgSpeed := meanSpeed(steps, height, duration)
	calories := (weight * avgSpeed * duration.Minutes()) / minInH
	calories *= walkingCaloriesCoefficient

	return calories, nil
}
