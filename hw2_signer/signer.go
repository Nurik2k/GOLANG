package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	// Создаем канал для входных данных
	in := make(chan interface{})
	out := make(chan interface{})
	wg := &sync.WaitGroup{} // Создаем объект sync.WaitGroup{} для ожидания завершения всех горутин

	for _, value := range jobs { // Итерируемся по значениям в слайсе jobs
		in = out
		out = make(chan interface{}) // Создаем новый канал выходных данных для текущего job

		wg.Add(1)

		go func(job2 job, input, output chan interface{}) {
			defer wg.Done()     // Уменьшаем счетчик WaitGroup на 1 при завершении горутины
			defer close(output) // Закрываем канал output при завершении горутины

			job2(input, output) // Выполняем текущий job, передавая входной и выходной каналы
		}(value, in, out) // Запускаем горутину для выполнения функции value с переданными аргументами value, in и out
	}

	wg.Wait() // Ожидаем завершения всех горутин, ожидающих в WaitGroup
}

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{} // Создаем объект sync.WaitGroup{} для ожидания завершения всех горутин
	mutex := &sync.Mutex{}  // Создаем объект sync.Mutex{} для синхронизации доступа к общим данным

	for i := range in { // Итерируемся по значениям в канале in до его закрытия
		wg.Add(1)                               // Увеличиваем счетчик WaitGroup на 1 перед запуском новой горутины
		go SingleHashService(i, out, wg, mutex) // Запускаем горутину для выполнения SingleHashService с переданными аргументами i, out, wg и mutex
	}

	wg.Wait() // Ожидаем завершения всех горутин, ожидающих в WaitGroup
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{} // Создаем объект sync.WaitGroup{} для ожидания завершения всех горутин

	for value := range in { // Итерируемся по значениям в канале in до его закрытия
		wg.Add(1)                           // Увеличиваем счетчик WaitGroup на 1 перед запуском новой горутины
		go MultiHashService(value, out, wg) // Запускаем горутину для выполнения MultiHashService с переданными аргументами value, out и wg
	}

	wg.Wait() // Ожидаем завершения всех горутин, ожидающих в WaitGroup
}

func CombineResults(in, out chan interface{}) {
	inputValue := make([]string, 0, MaxInputDataLen) // Создаем пустой слайс типа string с начальной емкостью MaxInputDataLen = 100

	for value := range in { // Итерируемся по значениям в канале in до его закрытия
		inputValue = append(inputValue, ConvToString(value)) // Преобразуем значение в строку с помощью ConvToString() и добавляем его в слайс inputValue
	}

	sort.Slice(inputValue, func(i, j int) bool { // Сортируем слайс inputValue с использованием анонимной функции сравнения
		return inputValue[i] < inputValue[j] // Возвращаем true, если элемент с индексом i меньше элемента с индексом j
	})

	out <- strings.Join(inputValue, "_") // Преобразуем слайс inputValue в строку, объединяя элементы с помощью разделителя "_" и отправляем результат в канал out
}

// Функция конвертации интерфейса в строку
func ConvToString(inter interface{}) string {
	str := fmt.Sprintf("%v", inter)
	return str
}
