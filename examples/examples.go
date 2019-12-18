package examples

import (
	"fmt"
	"runtime"
	"sync"
	"time"
	"unsafe"
)

// TODO добавить код ниже в конспект со строками и анально разобрать этот код
//  в т.ч. через  бенчмарки посравнивать с байтес.Буффером и стрингс.Билдером
func concatenate(strs ...string) string {
	var totalLen int
	for i := range strs {
		totalLen += len(strs[i])
	}
	concated := make([]byte, totalLen)
	cursor := 0
	for _, str := range strs {
		cursor += copy(concated[cursor:], *(*[]byte)(unsafe.Pointer(&str)))
	}
	return *(*string)(unsafe.Pointer(&concated))
}

var isTraceEnabled = true

func Run() {
	// Про массивы
	//aboutArrays()

	// Про слайсы
	//aboutSlices()

	// Про мапы
	aboutMaps()

	// Параллельность и всовместно используемые переменные
	//aboutMutex()
}


func aboutArrays() {
	fmt.Println("\n *** \n -------------- arrays ------------------\n *** \n")
	printLineTrace()
	// - массивы это составной типа данных, их ЗНАЧЕНИЯ создаются путем конкатенации в памяти других значений
	// - массивы неизменяемы и гомогенны(могут содержать только значения одинакового данных)
	// - само значение массива можно изменить по индексу на значение такого же типа
	// - по умолчанию элементам массива присваиваются значения их типов данных по умолчанию.
	//      т.е. для var a = [n]int{}; a[n] == 0
	//
	// - частью типа массива является размерность массива, следовательно, мы не можем присвоить
	//      var a [3]int
	//      a = [4]int{} - ошибка компиляции
	// - размер массив должен быть константным значением, которое может быть вычесленно во время компиляции

	// в данном случае создан массив длинной 3 элемента, тип которого [3]int
	var a [3]int = [3]int{0, 1, 2}
	printSliceInt(a[:])
	printLineTrace()
	// [...] - массив инициализировался динамически по кол-ву элементов
	var a1 = [...]int{0, 1, 2, 3, 4, 5}
	printSliceInt(a1[:])

	printLineTrace()
	// элементы можно присваивать по индексу
	var a2 = [...]string{0: "123", 3: "345"}
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(a2), cap(a2), &a2, a2)

	printLineTrace()
	// массивы можно сравнивать, если элементы, из которых он состоит сравниваются между собой
	cA := [3]int{2, 3, 4,}
	cB := [3]int{1, 2}
	cC := [3]int{2, 3, 4,}
	fmt.Println(cA == cB, cB == cC, cA == cC)
}

func aboutSlices() {
	fmt.Println("\n *** \n -------------- slices ------------------\n *** \n")
	printLineTrace()
	// слайсы представляют собой легковесную структуру, которая предоставляет доступ к подпоследовательности элементов массива
	// и, возможно даже ко всем элементам БАЗОВАОГО МАССИВА
	// срез состоит из длинны, размера и указателя (ЩА ВАЖНО) на:
	//     ***
	//     первый элемент массива доступный через срез. Т.е. не обязательно на первый, потому что
	//     мы можем создать слайс, например таким обращом:
	//              массив arr := [3]int{...}
	//              а слайс создать на его основе со второго по 3ий эелемент s := arr[2:3]
	//     ***
	// пример:
	arrBase := [6]string{"Cat", "in", "a", "vacuum", "is", "anything"}

	sliceFromArrBase := arrBase[3:6]
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(arrBase), cap(arrBase), &arrBase, arrBase)
	printSliceString(&sliceFromArrBase)

	// указатель на 2ой элемент массива равен указателю на слайс
	printLineTrace()
	fmt.Printf("base array elem ptr: %p; slice elem ptr: %p;\n&arrBase[2] == &sliceFromArrBase[0]: %t\n", &arrBase[3], sliceFromArrBase, &arrBase[3] == &sliceFromArrBase[0])

	printLineTrace()
	// несколько слайсов могут ссылаться на один и тот же базовый массив, в т.ч. когда один слайс создается на основе другого слайса
	otherSlice := sliceFromArrBase[1:3]
	fmt.Printf("base slice elem ptr: %p; slice elem ptr: %p;\n&otherSlice[3] == &sliceFromArrBase[0]: %t\n", &arrBase[4], otherSlice, &arrBase[4] == &otherSlice[0])

	printLineTrace()
	// если создавать слайс на основе другого слайса, который будет выходить за пределы cap, то будет паника
	func() {
		// востановим прогу
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered:", r)
			}

		}()
		// будет паника
		_ = sliceFromArrBase[2:20]
	}()

	printLineTrace()
	// СРАЗУ РЕМАРОЧКА НАХОЙ: У ДОНОВАНА НИ СЛОВА
	// если создавать слайс срезанием слева - размерность сократится и будет равна длинне
	otherArr := [7]string{"Cat", "in", "the", "fog", "in", "a", "vacuum"}
	shortenedCapSlice := otherArr[3:7]
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(otherArr), cap(otherArr), &otherArr, otherArr)
	printSliceString(&shortenedCapSlice)
	// если справа - то сократится длинна, но вместимость будет как и прежде
	savedCapSlice := otherArr[0:4]
	printSliceString(&savedCapSlice)

	printLineTrace()
	// Если увеличивать слайс в пределах его cap, но за пределами его len, то слайс просто расширится
	// как мы видим, значения нижележащего массива никуда не девались
	extendedSlice := savedCapSlice[0:7]
	printSliceString(&extendedSlice)

	printLineTrace()
	// при передаче слайса в фцию можно изменять значения внутри
	// т.к. при копировании слайса создаётся псевдоним указателя
	// p.s. про расширение через append будет ниже
	mutateSlice := func(s []string) bool {
		if len(s) == 0 {
			return false
		}
		s[0] = "mutated"
		return true
	}
	// изменяем слайс
	ok := mutateSlice(extendedSlice)
	if ok {
		// слайс изменился
		printSliceString(&extendedSlice)

		// его базовый массив тоже
		fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(otherArr), cap(otherArr), &otherArr, otherArr)
	}

	printLineTrace()
	// для сравнения двух слайсов нет встроенных инструментов в го
	// сравнение можно написать самоятоятельно
	equal := func(x, y []string) bool {
		if len(x) != len(y) {
			return false
		}
		for i := range x {
			if x[i] != y[i] {
				return false
			}
		}
		return true
	}
	fmt.Println(equal(extendedSlice, savedCapSlice))
	printSliceString(&extendedSlice)
	printSliceString(&savedCapSlice)

	printLineTrace()
	// нулевым значением для слайса является нил
	// лучше проверять нуденвой слайс через len(s)
	var s []int
	printSliceInt(s);
	fmt.Printf("a:%#v\n", s)

	// len(s) == 0, s == nil
	s = nil
	printSliceInt(s);
	fmt.Printf("a:%#v\n", s)

	// len(s) == 0, s == nil
	s = []int(nil)
	printSliceInt(s)
	fmt.Printf("a:%#v\n", s)

	// len(s) == 0, s != nil
	s = []int{}
	printSliceInt(s)
	fmt.Printf("a:%#v\n", s)

	// Создание слайсов
	printSliceInt(s)
	// make создает неименнованную переменную массива и возвращает его срез
	// сам массив при этом доступен только через срез
	//  - в этом случае слайс будет представлять весь массив
	msFull := make([]string, 5) // cap == len

	//  - в этом случае только первые пять элементов базового массива, остальные доступны для расширения
	msSliced := make([]string, 5, 10) // len == 5; cap == 10
	printSliceString(&msFull)
	printSliceString(&msSliced)

	printLineTrace()
	// append - добавляет новые элементы в слайс
	// рассмотрим ф-ци. для понимания работы append
	appendCat := func(cats []string, cat string) []string {
		var newCats []string
		newCatsLen := len(cats) + 1

		if newCatsLen <= cap(cats) {
			// значит, есть место для роста (хоть 1 свободный элемент)
			newCats = cats[:newCatsLen] // расширили длинну на один элемент, но в пределах cap
		} else {
			// в другом случае, получается, что места для роста нет
			// значит нужно выделить новый массив. При этом увеличив вместимость
			// в два раза для линейной амортизированной сложности(что это я пока хз)
			newCatsCap := newCatsLen
			if newCatsCap < 2*len(cats) {
				newCatsCap = 2 * len(cats)
			}

			newCats = make([]string, newCatsLen, newCatsCap)
			// copy может возвращать кол-во фактически скопированных элементов
			copy(newCats, cats)
		}
		newCats[len(cats)] = cat
		return newCats
	}

	cats := []string{"cat1", "cat2", "cat3"}

	cats = appendCat(cats, "newCat")

	// как видно, был возвращен новый слайс (pointer не равны)
	printSliceString(&cats)
	printSliceString(&cats)

	printLineTrace()
	// встроенная фция append обладает более сложными мехнизмами стратегии роста
	var x = make([]string, 0, 10)
	printSliceString(&x)
	x = append(x, "val1")
	// слайс будет ссылаться на тот же базовый массив
	z := append(x, "val1")
	printSliceString(&x)
	printSliceString(&z)
	printSliceString(&x)
}

func aboutMaps() {
	// мапа - это ссылка на хеш таблицу
	// ключи, как и значения должны иметь одинаковый тип данных
	// но тип данных  ключа может быть не равен типу значения
	// ключ должен быть сравниваемым значением
	m := map[string]string{
		"key1": "value1",
	}

	printMapStringString(m)
}

func aboutMutex() {
	fmt.Println("\n *** \n -------------- mutex ------------------ \n *** \n")
	printLineTrace()
	// Состояние гонки - когда есть несколько одновремнных обращений к переменной
	// их нескольких горутин
	// и хоть одно из них является записью
	var balance int

	Deposit := func(amount int) {
		balance = balance + amount
	}

	Balance := func() int {
		return balance
	}

	for i := 0; i <= 800; i++ {
		go Deposit(100)
	}

	time.Sleep(time.Millisecond * 100)
	fmt.Println(Balance())

	printLineTrace()
	// Все оказывается еще хуже, если гонка данных включает переменную типа, большего,
	// чем одно машинное слово, такого как интерфейс, строка или срез.
	// например
	var x []int
	func() {
		// востановим прогу, т.к. она очень вероятно,
		// что будет падать с паникой при x[999999] = 1
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered:", r)
			}
			// слайс состоит из 3х элементов
			// длинна, указатель на массив и емкость
			fmt.Println(len(x), cap(x))
		}()
		// так вот, может получится такое, что у слайса длинной 1000000
		// под капотом окажется массив с длинной 10, если выделять память под него конкурентно
		go func() { x = make([]int, 10) }()
		go func() { x = make([]int, 1000000) }()
		// соотвественно тут вылезет паника.
		x[999999] = 1
	}()

	printLineTrace()
	// первый способ избежать гонки - не записывать переменную
	// т.е. если только чтение из нескольких горутин - это всегда будет потоко-безопасно

	// второй способ - это ограничить работу с данными одной горутиной.
	// т.е. в случае параллельного выполнения одной горутине отдаётся только одна переменная для выполнения конкретной цели

	// Если все таки нужно работать с переменной в несколько потоков, то можно использовать каналы
	// перепишем симулятор банка
	var deposits = make(chan int)
	var balances = make(chan int)

	Deposit = func(amount int) { deposits <- amount }
	Balance = func() int { return <-balances }

	teller := func() {
		var balance int
		for {
			select {
			case amount := <-deposits:
				{
					balance += amount
				}
			case balances <- balance:

			}
		}
	}
	// лучше вынести в инит
	// сейчас это тут только для непрерывности повествования
	go teller()

	for i := 0; i <= 99; i++ {
		go Deposit(10)
	}

	time.Sleep(time.Millisecond * 500)

	fmt.Println(Balance())

	printLineTrace()
	// для синхронизации данных можно использовать также мьютекс
	// sync.Mutex
	// Мьютексы призваны обеспечить сохранение определенных инвариантов совместно используемых
	// переменных в критических точках во время выполнения программы.
	// Т.е. имеется ввиду, что операции блокировок мьютексами созданы
	// для определения условий целостности данных при конкурентном использовании.
	type sMu struct {
		// принято мьютекс объявлять сверху данных, которые мьютекс призван защищать
		sync.Mutex
		i int
	}

	smu := sMu{sync.Mutex{}, 0}

	// проверим блокировку Lock() когда
	// работа с переменной ограничивается для чтения и записи во время блокировки
	for i := 0; i <= 5; i++ {
		go func() {
			smu.Lock()
			// Unlock() принято вызывать через дефер, т.к. это гарантирует, что мы закроем все исключения
			// в т.ч. паники, которые могут возникать во время работы
			defer smu.Unlock()
			fmt.Println("smu lock()")
			smu.i++
			// видим, что пока работает блокировка
			// горутины ждут своей очереди для записи
			time.Sleep(time.Second * 1)
			fmt.Println("smu.i: ", smu.i)

			fmt.Println("smu unlock()")
		}()
	}

	// эмулируем ожидание работы горутин
	time.Sleep(time.Second * 6)

	fmt.Println(smu.i)
	// Мьютексы нереентерабельны(нельза заблокировать уже заблокированный мьютекс)
	// При работе мьютексами всегда нужно гарантировать, что ф-ция вызывающая блокировку, верент все
	// на свои места. Что бы если вдруг будет в тот же момент вызываться другая ф-ция, которая
	// также вызывает блокировку не случилась взаимоблокировка, в результате чего  программа остановится.
	// распространенное решение для этого - делают ф-цию, которая не экспортируема и просто делает работу
	// предполагая, что блокировку делает уже потребитель и делают для нее экспортируемую обертку, которая выполняет блокировку
	// ГОВЕЙ в данном случае -
	//             --------------------------------------------------------------------
	//             При использовании мьютекса убедитесь, что и он, и переменные,
	//             которые он защищает, не экспортируются, независимо от того, являются
	//             ли они переменными уровня пакета или полями структуры.
	//             --------------------------------------------------------------------

	printLineTrace()
	// sync.RwMutex
	// Блокировка типа "несколько читателей - один писатель"
	// позволяет использовать блокироваться только при записи, при это оставляя возможость
	// для чтения.
	// Объяснение от Артема Кузнецова - https://gist.github.com/blank-teer/2b446ae0c3df62518ef06aa99f4423af
	//
	// Пример: все горутины, которые сейчас хотят прочитать данные будут блокироваться
	// пока происходит запись, после чего получат конкурентный доступ к данным и наоборот - операция записи
	// будет ожидать, пока все читатели завершат работу с переменной и освободят свои блокировки.
	type sRWmu struct {
		sync.RWMutex
		i int
	}

	srwmu := sRWmu{sync.RWMutex{}, 10}

	// тут будем  конкурентно читать, блокируясь только на запись
	for i := 0; i <= 5; i++ {
		go func() {
			// блокируем операции записи
			srwmu.RLock()
			// снимаем блокировку
			defer srwmu.RUnlock()

			// видим, что пока работает блокировка
			// читать мы все равно можем конкурентно. В отличие от предыдущего примера,
			// не смотря на то, что спим 1 сек, 10 операций чтения выполнятся за ~ 1 сек
			time.Sleep(time.Second * 1)
			fmt.Println("srwmu.i: ", srwmu.i, "Read")
			fmt.Println("srwmu RUnlock()")
		}()
	}

	// а тут будем писать
	// операции записи будут ждать, пока закончат читатели и наоборот
	for i := 0; i <= 2; i++ {
		go func() {
			srwmu.Lock()
			srwmu.i++
			time.Sleep(time.Second * 1)
			fmt.Println("srwmu.i: ", srwmu.i, "Write")
			srwmu.Unlock()
			fmt.Println("srwmu unlock()")
		}()
	}

	// эмулируем ожидание работы горутин
	time.Sleep(time.Second * 6)
}

// todo
// о стеке вызовов
// и о функциях

func printSliceInt(s []int) {
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(s), cap(s), &s, s)
}

func printSliceString(s *[]string) {
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(*s), cap(*s), s, s)
}

func printMapStringString(m map[string]string){
	fmt.Printf("len: %d; value: %+v;", len(m))
}

func printLineTrace() {
	if !isTraceEnabled {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d\n", fn, line)
}
