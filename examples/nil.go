package examples

import (
	"fmt"
)

// https://speakerdeck.com/campoy/understanding-nil?
func aboutNil() {
	// nil имеет много смыслов, это может быть 0, ничего, никогда и ничто
	// nil в го имеет ввиду отсутствие типа у значения
	// nil это предопроеделнное (нулевое) значение для составных типов в go
	// pointer, channel, func, interface, map, slice
	// все эти типы имеют встроенным какой-то базовый тип
	printLineTrace()
	var nilChan chan struct{}
	var nilSlice []string
	fmt.Println(nilSlice)
	fmt.Println(nilChan)

	// ??
	// var nil = 123
	// fmt.Println(nil)
	// fmt.Println(nilChan==nil)


	printLineTrace()
	// nil for pointers
	// - указыват на nil, т.е. ни на что не ссылается
	// - нулевое значение для указателя

	// для слайсов
	// - имеет ввиду, что у слайса отсутсвует указатель на сам низлежащий массив

	// для функций, каналов и функций
	// просто отсутсвие занчения низлежащего типа
	// для функций занчением является реализация функции

	// для интерфейсов нулевое значение означает отсутствие реализации (динамического типа)
	// при этом само динамическое значения типа может быть nil, в том время как динамический тип
	// уже присутствует
	// интерфейс пустой, будет nil
	var doer Doer
	fmt.Println(doer==nil)
	var human *Human
	// не инициализированный human == nil
	fmt.Println(doer==nil)
	// как только интерфейс получил динамический тип, он больше не равен nil
	doer = human
	fmt.Println(doer==nil)


	printLineTrace()
	// из этого следует один крутой пример:
	// что не следует объявлять реализации для кастомных error заранее, т.к. они не будут nil, даже если
	// не будет присвоено значение реализации!!!
	err := processDoError()
	// !!FALSE!!!
	fmt.Println(err==nil)

	// В целом, все приводит к выводу о том, что ниловое значение нужно делать полезным
	// т.е. там где может возникает паника из-за нилового указателя
	// делат проверку на нил и возвращать значение которое будет инициализированно
	// особенно в случае методов с получателями через указатель
	printLineTrace()
	var pointerReceiver *SomeType
	// тут паники не будет, но без проверки на нил - паниковали бы
	fmt.Println(pointerReceiver.MethodAtPointerReceiver())

	// про нилвое значение в слайсах - часто использование  нил слайса
	// через аппенд быстрее выделения через мейк
	printLineTrace()
	var s []int
	for i:=0; i <= 10; i ++ {
		s = append(s, i )
	}
	fmt.Println(s)

	printLineTrace()
	// ниловые мапы паникуют при попытке записи
	// но при попытке чтения все норм
	var m map[string]int
	fmt.Println(m["key"])
	// паника
	// m["key"] = 0
	// но паники не будет вот так
	for range m {
		m["key"] = 0
	}
	// это можно использовать для защиты от случая записи в ниловую мапу
	// т.к. итерация не произойдет
	// (аля ниловая мапа только для чтения)))
	//
	// паники не будет
	writeToMap(nil)


	printLineTrace()
	// каналы же паникуют, если закрыть канал с нил и при отправке в закрытый канал
	// (закрытие закрытого канала тоже вызывает панику)
	// при отправке в нил  канал - вечное ожидание
	// расскоментить для примера
	// nilChanDeadlock()
}

type Doer interface {
	Do()
}

type Human struct {}; func(h Human) Do() { fmt.Println("human Do()")}


type doError struct {}; func (d doError) Error() string {return "doError"}

func processDoError() error {
	return &doError{}
}

type SomeType struct {
	name string
}

func (s *SomeType) MethodAtPointerReceiver() string {
	if s == nil {
		return ""
	}
	 return s.name
}

func writeToMap(m map[string]int)  {
	for range m {
		m["key"]  = 1
	}
}

func nilChanDeadlock() {
	wait := make(chan struct{})
	var nilCh chan int

	go func() {
		nilCh <- 0

	}()

	go func() {
		// deadlock
		fmt.Println("nil chan out", <-nilCh)
		wait <- struct{}{}
	}()

	<-wait
}