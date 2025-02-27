package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep   = 0.65  // средняя длина шага.
	mInKm     = 1000  // количество метров в километре.
	minInH    = 60    // количество минут в часе.
	kmhInMsec = 0.278 // коэффициент для преобразования км/ч в м/с.
	cmInM     = 100   // количество сантиметров в метре.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// Парсинг строки
	words := strings.Split(data, ",")
	if len(words) != 3 {
		err := errors.New("len(word)!=3")
		return 0, "", 0, fmt.Errorf("slice length error: %w", err)
	}
	step, err := strconv.Atoi(words[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("conversion error to integer: %w", err)
	}
	//вид активности activity
	activity := words[1]

	t, err := time.ParseDuration(words[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("conversion to time format error: %w", err)
	}
	return step, activity, t, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	dist := float64(steps) * lenStep / mInKm
	return dist
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	timeTraining := time.Duration.Hours(duration)
	dist := distance(steps)
	speedAverage := dist / timeTraining
	return speedAverage
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	step, activity, t, err := parseTraining(data)
	if err != nil {
		return ""
	}
	timeTraining := time.Duration.Hours(t)
	dist := distance(step)
	averageSpeed := meanSpeed(step, t)
	killCallSprint := RunningSpentCalories(step, weight, t)
	killCallWalk := WalkingSpentCalories(step, weight, height, t)
	switch activity {
	case "Ходьба":
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий %.2f\n", activity, timeTraining, dist, averageSpeed, killCallWalk)
	case "Бег":
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий %.2f\n", activity, timeTraining, dist, averageSpeed, killCallSprint)
	}
	return fmt.Sprintf("Неизвестный тип тренировки")
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	speedAverage := meanSpeed(steps, duration)
	kcal := ((float64(runningCaloriesMeanSpeedMultiplier) * speedAverage) - runningCaloriesMeanSpeedShift) * weight
	return kcal
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	speedAverage := meanSpeed(steps, duration)
	durHor := time.Duration.Hours(duration)
	kCal := ((float64(walkingCaloriesWeightMultiplier) * weight) + (speedAverage*2/height)*float64(walkingSpeedHeightMultiplier)*durHor*float64(minInH))
	return kCal
}
