package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// 1. Разделяем входную строку по запятой
	parts := strings.Split(data, ",")

	// 2. Проверяем, что получили ровно 2 части
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат данных: ожидается 2 части, разделённые запятой")
	}

	// 3. Парсим количество шагов
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при преобразовании количества шагов: " + err.Error())
	}

	// 4. Проверяем, что количество шагов больше нуля
	if steps <= 0 {
		return 0, 0, fmt.Errorf("количество шагов должно быть положительным числом")
	}

	// 5. Парсим длительность прогулки
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка при парсинге длительности: " + err.Error())
	}

	// 6. Возвращаем результаты при успешном выполнении
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// 1. Парсим входные данные
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка при парсинге данных:", err)
		return ""
	}

	// 2. Проверяем количество шагов
	if steps <= 0 {
		return ""
	}

	// 3. Вычисляем дистанцию в метрах
	distanceMeters := float64(steps) * stepLength

	// 4. Переводим дистанцию в километры
	distanceKm := distanceMeters / mInKm

	// 5. Вычисляем потраченные калории
	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		fmt.Println("Ошибка при рассчёте каллорий:", err)
		return ""
	}

	// 6. Формируем результирующую строку
	result := fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.",
		steps,
		distanceKm,
		calories,
	)

	return result
}
