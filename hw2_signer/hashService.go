package main

import (
	"strings"
	"sync"
)

func SingleHashService(in interface{}, out chan interface{}, wg *sync.WaitGroup, mutex *sync.Mutex) {

	//Создаем срезы строк
	hash1 := make(chan string)
	hash2 := make(chan string)

	//Функция, которая закрывает каналы и завершает работу горутины
	defer func() {
		close(hash1)
		close(hash2)
		wg.Done()
	}()

	//Записываем в переменную dataCrc32, конвертированный интерфейс
	dataCrc32 := ConvToString(in)

	//Горутина, которая записывает в hash1, данные из data шифрованные DataSignerCrc32
	go func(data string, out chan string) {
		out <- DataSignerCrc32(data)
	}(dataCrc32, hash1)

	mutex.Lock()                               //Блокируем для других горутин
	dataMd5 := DataSignerMd5(ConvToString(in)) //Записываем в переменную dataMd5, конвертированный интерфейс и в то же время шифрование DataSignerMd5
	mutex.Unlock()                             //Открываем для других горутин

	//Горутина, которая записывает в hash2, данные из data шифрованные DataSignerCrc32
	go func(data string, out chan string) {
		out <- DataSignerCrc32(data)
	}(dataMd5, hash2)

	//Записываем в out результаты
	out <- (<-hash1) + "~" + (<-hash2)

}

func MultiHashService(in interface{}, out chan interface{}, wg *sync.WaitGroup) {

	//Создаем срез канала строк и даем конкретный размер
	hash := make([]chan string, 6)
	//Данный цикл создает шесть каналов типа string и сохраняет их в слайс hash. Цикл выполняется для инициализации каналов и подготовки их к использованию.
	for i := 0; i < 6; i++ {
		hash[i] = make(chan string)
	}

	//Этот defer блок выполняет закрытие всех каналов в слайсе hash и вызов метода Done() для объекта wg типа sync.WaitGroup.
	defer func() {
		for _, ch := range hash {
			close(ch)
		}
		wg.Done()
	}()

	strBuilder := strings.Builder{} //Создаем новый объект типа strings.Builder{}
	strItem := ConvToString(in)     //Новая переменная присваивает конвертированный интерфейс

	//
	for i, ch := range hash {
		go func(data string, out chan string) {
			out <- DataSignerCrc32(data)
		}(ConvToString(i)+strItem, ch) // В данной строке кода используется функция ConvToString(i) для преобразования целочисленного значения i в строку. Затем к полученной строке добавляется значение переменной strItem.
	}

	//В данном цикле мы проходимся по элементам hash
	for _, ch := range hash {
		//Внутри цикла вызывается оператор <-ch, который блокирует выполнение горутины до тех пор, пока не будет получено значение из канала ch.
		//Полученное значение строки с помощью оператора <- добавляется в strings.Builder с использованием метода WriteString().
		strBuilder.WriteString(<-ch)
	}

	out <- strBuilder.String() //Записываем в out результаты

}
