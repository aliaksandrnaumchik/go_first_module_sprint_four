package spentcalories

import (
	"fmt"
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

func parseTraining(data string) (int, string, time.Duration, error) {
	// 1. Разделяем строку по запятой
	parts := strings.Split(data, ",")

	// 2. Проверяем, что получили ровно 3 части
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("неверный формат данных: ожидается 3 части, получено %d", len(parts))
	}

	// 3. Парсим количество шагов
	stepsStr := strings.TrimSpace(parts[0])
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге количества шагов: %w", err)
	}

	// 4. Получаем вид активности
	activity := strings.TrimSpace(parts[1])

	// 5. Парсим продолжительность
	durationStr := strings.TrimSpace(parts[2])
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		return 0, "", 0, fmt.Errorf("ошибка при парсинге длительности: %w", err)
	}

	// 6. Возвращаем успешные значения
	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	// 1. Приводим steps к float64
	stepsFloat := float64(steps)

	// 2. Вычисляем длину шага
	stepLength := height * stepLengthCoefficient

	// 3. Вычисляем общую дистанцию в метрах
	totalDistanceMeters := stepsFloat * stepLength

	// 4. Переводим метры в километры
	distanceKm := totalDistanceMeters / mInKm

	return distanceKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// 1. Проверяем продолжительность
	if duration <= 0 {
		return 0
	}

	// 2. Вычисляем дистанцию
	dist := distance(steps, height)

	// 3. Переводим длительность в часы
	durationHours := duration.Hours()

	// 4. Вычисляем среднюю скорость
	// Скорость = дистанция (км) / время (часы)
	speed := dist / durationHours

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// 1. Парсим входные данные
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		fmt.Println("Ошибка парсинга данных:", err)
		return "", err
	}

	// 2. Определяем тип тренировки и рассчитываем параметры
	var (
		trainingDistance float64
		trainingSpeed    float64
		trainingCalories float64
		errCalories      error
	)

	switch strings.ToLower(activity) {
	case "бег":
		trainingDistance = distance(steps, height)
		trainingSpeed = meanSpeed(steps, height, duration)
		trainingCalories, errCalories = RunningSpentCalories(steps, weight, height, duration)
		if errCalories != nil {
			return "", errCalories
		}
	case "ходьба":
		trainingDistance = distance(steps, height)
		trainingSpeed = meanSpeed(steps, height, duration)
		trainingCalories, errCalories = WalkingSpentCalories(steps, weight, height, duration)
		if errCalories != nil {
			return "", errCalories
		}
	default:
		return "", fmt.Errorf("неизвестный тип тренировки: %s", activity)
	}

	// 3. Формируем результат
	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f",
		activity,
		duration.Hours(),
		trainingDistance,
		trainingSpeed,
		trainingCalories,
	)

	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверка входных параметров на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	// 2. Вычисление средней скорости
	speed := meanSpeed(steps, height, duration)

	// 3. Перевод длительности в минуты
	durationInMinutes := duration.Minutes()

	// 4. Расчет потраченных калорий
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// 1. Проверка входных параметров на корректность
	if steps <= 0 {
		return 0, fmt.Errorf("количество шагов должно быть больше 0")
	}
	if weight <= 0 {
		return 0, fmt.Errorf("вес должен быть больше 0")
	}
	if height <= 0 {
		return 0, fmt.Errorf("рост должен быть больше 0")
	}
	if duration <= 0 {
		return 0, fmt.Errorf("продолжительность должна быть больше 0")
	}

	// 2. Вычисление средней скорости
	speed := meanSpeed(steps, height, duration)

	// 3. Перевод длительности в минуты
	durationInMinutes := duration.Minutes()

	// 4. Базовый расчет калорий
	baseCalories := (weight * speed * durationInMinutes) / minInH

	// 5. Применение корректирующего коэффициента
	calories := baseCalories * walkingCaloriesCoefficient

	return calories, nil
}
