package os

import (
	"fmt"
	"sync"
	"time"
)

func RunLock() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
		}
	}()
	work()
}

// блокировки

// атомарные операции
//   - самая маленькая операция(та, которую нельзя прервать)
//   - неделима, не может выполнится частично

// критическая секция
//   - часть программы, где происходит работа с ресурсом
//   - чтение или запись

// взаимная блокировка (deadlock)
//   - когда процессы ждут один другого


// livelock
//  когда процессы работают, но ничего не происходит

// взаимное исключение
//     - свойство, которое зависит от порядка исполнения


// Дедлок возникает
// - 1 взаимного исключения
// - 2 hold and wait (процесс блокируется и ждет освобождения другого ресурса)
// - 3 Система не может отбирать ресурсы у процесса
// - 4 Циклическое ожидание

type ResourceA struct {
	mu sync.Mutex
	data string
}

type ResourceB struct {
	mu sync.Mutex
	data string
}

func work() {
	resa := ResourceA{}
	resb := ResourceB{}
	wg := new(sync.WaitGroup)


	wg.Add(1)
	// Процесс 1
	go func(wg *sync.WaitGroup) {

		resa.mu.Lock()
		// Ждет разблокировки ResourceB{}
		resb.mu.Lock()
		wg.Done()
	}(wg)

	// Процесс 2
	wg.Add(1)
	go func(wg *sync.WaitGroup) {

		resb.mu.Lock()
		// Ждет разблокировки ResourceA{}
		resa.mu.Lock()

		wg.Done()
	}(wg)

	// разблокировка никогда не случится

	// но мы узнаем про это только тогда, когда
	// родительская горутина полностью завершит работу

	// работа ниже будет продолжаться
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(time.Millisecond * 500)
		fmt.Println("ща будет дедлок")
	}(wg)

	// пока не произойдет выход
	wg.Wait()
}