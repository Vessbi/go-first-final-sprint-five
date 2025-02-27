package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
)

var (
	StepLength = 0.65  // длина шага в метрах
	distance   float64 //пройденная дистанция

)

func parsePackage(data string) (int, time.Duration, error) {
	// Строка в слайс строк
	word := strings.Split(data, ",")
	if len(word) < 2 {
		err := errors.New("len(word)!=2")
		return 0, 0, fmt.Errorf("slice length error: %w", err)
	}
	//преобразование строки в целое число
	step, err := strconv.Atoi(word[0])
	if err != nil {
		return 0, 0, fmt.Errorf("conversion error to integer: %w", err)
	}
	//Парсим время
	duration, err := time.ParseDuration(word[1])
	if err != nil {
		return 0, 0, fmt.Errorf("conversion to time format error: %w", err)
	}
	return step, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	//Парсим data через parsePackage
	step, t, err := parsePackage(data)
	if err != nil {
		return ""
	}
	if step <= 0 {
		return ""
	}
	distance = float64(step) * StepLength / 1000
	killCal := spentcalories.WalkingSpentCalories(step, weight, height, t)
	p := fmt.Sprintf("Количество шагов: %d.\nДистанция составила: %.2f км.\nВы сожгли: %.2f\n", step, distance, killCal)
	return p
}
