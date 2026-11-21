package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type TemperatureController struct {
	minTemp int
	maxTemp int
}

func NewTemperatureController() *TemperatureController {
	return &TemperatureController{
		minTemp: 15,
		maxTemp: 30,
	}
}

func (tc *TemperatureController) Update(op string, temp int) error {
	switch op {
	case ">=":
		if temp > tc.minTemp {
			tc.minTemp = temp
		}
	case "<=":
		if temp < tc.maxTemp {
			tc.maxTemp = temp
		}
	default:
		return fmt.Errorf("неизвестный оператор %q", op)
	}
	return nil
}

func (tc *TemperatureController) Current() int {
	if tc.minTemp <= tc.maxTemp {
		return tc.minTemp
	}
	return -1
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		log.Println(fmt.Errorf("ошибка чтения параметров n и k: %w", scanner.Err()))
		return
	}

	headerFields := strings.Fields(scanner.Text())
	if len(headerFields) < 2 {
		log.Println(fmt.Errorf("недостаточно параметров в первой строке: получено %d", len(headerFields)))
		return
	}

	n, err := strconv.Atoi(headerFields[0])
	if err != nil {
		log.Println(fmt.Errorf("ошибка парсинга n (%q): %w", headerFields[0], err))
		return
	}

	k, err := strconv.Atoi(headerFields[1])
	if err != nil {
		log.Println(fmt.Errorf("ошибка парсинга k (%q): %w", headerFields[1], err))
		return
	}

	controllers := make([]*TemperatureController, n)
	for i := 0; i < n; i++ {
		controllers[i] = NewTemperatureController()
	}

	total := n * k
	for i := 0; i < total; i++ {
		if !scanner.Scan() {
			log.Println(fmt.Errorf("ошибка чтения запроса %d: %w", i+1, scanner.Err()))
			break
		}

		reqFields := strings.Fields(scanner.Text())
		if len(reqFields) < 2 {
			log.Println(fmt.Errorf("некорректный формат запроса %q", scanner.Text()))
			continue
		}

		operatorName := reqFields[0]
		tempVal, parseErr := strconv.Atoi(reqFields[1])
		if parseErr != nil {
			wrappedErr := fmt.Errorf("ошибка парсинга температуры из %q: %w", reqFields[1], parseErr)
			log.Println(wrappedErr)
			continue
		}

		dept := i / k
		if err := controllers[dept].Update(operatorName, tempVal); err != nil {
			log.Println(fmt.Errorf("ошибка обновления температуры для отдела %d: %w", dept, err))
			continue
		}

		fmt.Println(controllers[dept].Current())
	}
}
