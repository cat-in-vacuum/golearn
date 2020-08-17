package examples

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
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
	//// strings
	//aboutStrings()
	//// const
	//constExample()
	//// Про массивы
	//aboutArrays()
	//// Про слайсы
	//aboutSlices()
	//// Про мапы
	//aboutMaps()
	//
	//aboutFunc()
	//
	//aboutInterfaces()
	//
	//// Параллельность и всовместно используемые переменные
	//aboutMutex()

	aboutChannels()
}

func aboutStrings() {
	s := "simple strings"
	subS := s[:5]

	// хранение строк:
	// строки не изменяемы, копирование или получение подстроки Б из строки А
	// очень дешево, т.к скопированная строка или подстрока будут ссылаться на ту же память
	// что и исходная строка(не будет выделяться новая).
	// В примере s:=s := "simple strings", subS := s[:5]
	// fmt.Printf(s, subS,) -> |%s|; |%s|, subS будет ссылаться на первый символ строки s длинной 5

	fmt.Println(s, subS)

	// Детали устройства строк:
	// Сама строка это неизменяемый набор байт
	// Для подсчета кол-ва рун, а не символов байт,
	// лучше использовать utf8.RuneCountInString(s), т.к. некоторые руны могут состоять из нескольких байт.
	// При итерации по строке при помощи range происходит неявное декодирование к UTF-8
	// по этому итарция будет происходить именно по рунам, а на по байтам. Пример итерации по "Hello, 世界"

	for i, r := range "Hello, 世界" {
		fmt.Printf("%d\t%q\t%d\n", i, r, r)
	}

	// еще пример итерации по строке:

	var str = "123abcэюя!?"
	fmt.Println("|type|value")
	for i, char := range str {
		fmt.Printf("str char - byte:%d; type:%T; value:%q\n", char, char, char)
		fmt.Printf("str[i]   - byte:%d; type:%T; value:%q\n", str[i], str[i], str[i])
	}
}

func constExample() {
	const num = 3

	//Константы:
	//Есть два типа констант - типизированные и нетипизированные. Воторые предоставляют большую точность в случае чисел, т.к.
	//могут быть использованны в большем кол-ве вычеслений без  преобразований.
	//В примере ниже
	//const num = 3
	//  можно присвоить всем численным типам:

	var f32 float32
	var f64 float64
	var i int
	var r rune
	f32 = num
	f64 = num
	i = num
	r = num
	// Имеется шесть вариантов таких несвязанных констант,
	// именуемых нетипизированиым булевым значением, нетипизированным целым числом, нетипизированной руной,
	// нетипизированным числом с плавающей точкой, нетипизированным комплексным числом и нетипизированной строкой

	fmt.Printf("var: %f|%T; const: %d||%T\n", f32, f32, num, num)
	fmt.Printf("var: %f|%T; const: %d||%T\n", f64, f64, num, num)
	fmt.Printf("var: %d|%T; const: %d||%T\n", i, i, num, num)
	fmt.Printf("var: %d|%T; const: %d||%T - и даже руне:)\n", r, r, num, num)
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
	cA := [3]int{2, 3, 4}
	cB := [3]int{1, 2}
	cC := [3]int{2, 3, 4}
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
	printSliceInt(s)
	fmt.Printf("a:%#v\n", s)

	// len(s) == 0, s == nil
	s = nil
	printSliceInt(s)
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
	fmt.Println("\n *** \n -------------- maps ------------------\n *** \n")
	printLineTrace()
	// мапа - это ссылка на хеш таблицу
	// ключи, как и значения должны иметь одинаковый тип данных
	// но тип данных  ключа может быть не равен типу значения
	// ключ должен быть сравниваемым значением
	m := map[string]string{
		"cat":   "meow",
		"cat_1": "meow_1",
	}

	printMapStringString(m)

	printLineTrace()
	// это все равно, что
	mm := make(map[string]string)

	mm["cat"] = "meow"
	mm["cat_1"] = "meow_1"

	printMapStringString(mm)

}

type CatInVacuum struct {
	Name  string
	Speed int
}

func aboutFunc() {
	printLineTrace()
	// функции ялвяются значениями первого класса:
	// т.е. функция может иметь тип, может быть присвоенна в переменную
	// возвращенна из другой ф-ции и т.д.
	type wrapStringOutFunc func(func(string), string) func()

	var wrapFunc wrapStringOutFunc
	wrapFunc = func(printString func(string), s string) func() {
		return func() {
			fmt.Println("-------")
			printString(s)
			fmt.Println("-------")
		}
	}

	printString := func(s string) {
		fmt.Println(s)
	}

	wrapped := wrapFunc(printString, "i am cat")
	wrapped()

	printLineTrace()
	// функцию можно сравнивать с нил, но нельзя сравнивать между собой
	// по этому нельзя использовать как ключ для мап
	clearF := func() {}
	var nilFunc wrapStringOutFunc

	fmt.Println(wrapped == nil, "ф-ция проинизиализированна")
	fmt.Println(clearF == nil, "ф-ция, которая делает ничего")
	fmt.Println(nilFunc == nil, "ф-ция объявлена, но инициализированна")

	printLineTrace()
	// ф-ция может быть анонимной
	// соответсвенно можно делать замыкания:
	clojure := func() func() int {
		var count int
		return func() int {
			count++
			return count
		}
	}
	// если ты у меня спросишь
	// - братишка, а как так получаестя, что переменная запомниает свое состояние?
	// то я тебе отвечу:
	// - бро, потому что анониманая ф-ция все еще ссылается на перемнную,
	//   соответственно, сборщик мусора, такого, как ты, её не удаляет.
	//   Просто потому что на переменную ссылается всратая анонимная ф-ция.
	counter := clojure()
	fmt.Println(counter())
	fmt.Println(counter())
	fmt.Println(counter())

	// ЗАХВАТ ПЕРЕМЕННЫХ ИТЕРАЦИИ
	printLineTrace()
	var intSlice = []int{1, 2, 3, 4, 5}
	var pointerSlice = make([]*int, 0)

	// вторая переменная, которую возвращает range - является как-бы итеративным указателем
	// на переменную, который переиспользуется на каждой итерации цикла. Т.е добавляя на каждой
	// итерации цикла item в pointerSlice на самом деле добавляется ОДНА И ТАЖЕ переменная, которая является
	// указателем на значение, которое хранится на момент итерации в item.
	// Т.е. получается, что на последней итерации внутри item будет 5, и , т.к. все значения
	// pointerSlice это указатель на одну пеменную item, при выводе значений получается, что
	// все значения равны значению на последней итерации цикла.

	// КОРОТКО: range вторым аргументом возвращает переиспользуемый на каждой итерации указатель,
	// по этому в нашем примере добавляя именно сам указатель, получается что все элементы pointerSlice
	// хранят ссылку на одну переменную.

	for _, item := range intSlice {
		pointerSlice = append(pointerSlice, &item)
	}

	for _, item := range pointerSlice {
		fmt.Println(*item)
	}

	printLineTrace()
	// как этого избегать?
	// очистим слайс
	pointerSlice = make([]*int, 0)
	for _, item := range intSlice {
		// внутри, цикл создает свою область видимости (лексический блок)
		// все значения, созданные циклом захватывают и совместно используют одну перменную
		// адресумое место в памяти

		// просто сохраним значение второй переменной range в новую перемнную
		value := item
		// и уже после будем добавлять ссылку на него в новый слайс
		pointerSlice = append(pointerSlice, &value)
	}
	for _, item := range pointerSlice {
		fmt.Println(*item)
	}
}

type Location struct {
	Vacuum string
}

type Cat struct {
	Legs int
	Name string
	In   Location
}

func aboutStructs() {
	printLineTrace()
	// лучше объявлять эксземпляр структуры, явно указывая поля, т.к. нет необходимости
	// запоминать их порядок, что чревато ошибками
	cat := Cat{
		Legs: 6,
		Name: "Unnamed",
		In: Location{
			Vacuum: "room`s vacuum",
		},
	}
	printStruct(cat)

	printLineTrace()
	// все приведенные ниже записи эквивалентны
	initCat := new(Cat)
	initCat = &Cat{}
	*initCat = Cat{}
	// для модификации, всегда нужно передавать по указателю
	// т.к. при передече ф-ция всегда получает только копию
	modifyCatName("Koluchka", &cat)
	printStruct(cat)

	printLineTrace()
	// сравнить структуры можно только если все её поля сравниваемы
	otherCat := Cat{
		Legs: 8,
		Name: "Cookies",
		In: Location{
			Vacuum: "kitchen vacuum",
		},
	}
	fmt.Println(otherCat == cat)
}

func modifyCatName(name string, cat *Cat) {
	cat.Name = name
}

func (c CatInVacuum) Meow() {
	printLineTrace()
	fmt.Printf("I am %s, meow!\n", c.Name)
}

func (c CatInVacuum) Walk() {
	printLineTrace()
	fmt.Printf("My speed %d!\n", c.Speed)
}

func aboutInterfaces() {
	fmt.Println("\n *** \n -------------- interfaces ------------------\n *** \n")
	printLineTrace()
	// интерфейс - абстрактный тип, который определяет поведение
	// в отличие от конкретного типа, который определяет данные
	// интерфейсы могут встраиваться друг в друга при этом одновремнно удовлетворять
	// тем реализациям, поведение которых полностью соответсвует методам ожидаемой реализации

	// интерфейс, который описывает что-то мяукающее
	type Meower interface {
		Meow()
	}

	// описывает нечто ходящее
	type Walker interface {
		Walk()
	}

	// интерфейс, который состоит из двух абстракций:
	// - нечто мяукающее
	// - нечто ходящее
	type Cat interface {
		Meower
		Walker
	}

	letWalk := func(c Walker) {
		c.Walk()
	}

	letMeow := func(c Meower) {
		c.Meow()
	}

	// опишем тип, который удовлетворит двум интерфейсам (см реализацию CatInVacuum)
	var cInV = CatInVacuum{
		Name:  "CatInVacuum",
		Speed: 55,
	}

	// т.к. тип имеет поведение Meow(), мы можем его передать в те ф-ции, которые
	// имеют в своей сигнатуре Meower
	letMeow(cInV)
	// тоже самое касается Walk(). Мы можем передать переменную в ф-цию, которая ожидает
	// Walker
	letWalk(cInV)

	printLineTrace()

	// объявим нулевой интерфейс типа io.Writter
	// нулевое значение интерфейса - nil
	var w io.Writer
	fmt.Println(w)

	printLineTrace()
	// в своей концепции интерфейс состоит из дескриптора типа
	// и самого значения этого типа
	// называются они динамический тип  и динамическое значение типа
	// разберем это:

	//  на этом этапе тип и значение будут равны nil
	// Вызов метода нулевого интерфейса приводит к панике
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				fmt.Println("recovered from the panic that was caused by calling the nil-interface method")
			}
		}()

		w.Write([]byte("hello"))
	}()

	// именно в этом месте присовим уже проинциализированную реализацию io.Writer
	w = os.Stdout
	// это преобразование включает в себя неявное преобразование из конкретного типа
	// в тип интерфейса, как будто это выглядит io.Writer(os.Stdout)
	// в данном случае, динмический тип переменной становится
	// равным дескриптору типа указателя *os.File, а динамесческое значение хранит копию os.Stdout
	// которое является уже самим указателем на некую переменную типа os.File
	// т.е. еще раз : динамический тип хранит сам ДЕСКРИПТОР ТИПА указателя на os.File - т.е. некую абстрактную дичь, которая определяет сам тип
	// а динамическое значение хранит уже саму ПЕРЕМЕННУЮ указателя на os.File

	printLineTrace()
	// В общем случае, во время компиляции мы не знаем, каким будет динамический тип значния, которое запакованно в интерфейс
	// по этому в этом случае используется динамическая диспетчеризация.
	// Вместо непосредственного вызова компилятор генерирует код для получения адреса вызываемого метода, в данном случае Write
	// После чего осуществляет косвенный вызов этого метода.
	// Аргументом для вызова становится копия
	w.Write([]byte("Hello, cat!"))

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

func aboutChannels() {
	// каналы нужны для общения между горутинами
	// канал является ссылкой на структуру данных
	// т.е. вызывающая и вызываемая ф-ция ссылаются на одну и ту-же структу данных
	// нулевое значение канала - nil
	var nilChan chan int
	printLineTrace()
	fmt.Println(nilChan == nil)
	// после инициализации канал уже не считается нулевым
	ch := make(chan int)
	printLineTrace()
	fmt.Println(ch == nil)
	// каналы можно сравнивать между собой
	printLineTrace()
	// каналы равны, если они являются ссылкой на одну и ту же структуру данных
	copyCh := ch
	fmt.Println(ch == copyCh)
	// каналы могут отправлять и получать значения. Это называется communications

	printLineTrace()
	// закрытый канал close(chan)
	//    - устанавалиется флаг, указывающий что по каналу больше не будут передаваться данные
	//    - при отправке в закрытый канал - паника
	//    - при получении из закрытого канала сперва возвращаются раннее записанные значения
	//          после того, как они закончатся - вернется тип данных канала по умолчанию
	close(ch)
	// проверить состояние канала(открыт или закрыт) можно при помощи второй переменной bool
	v, ok := <-ch
	fmt.Println(v, ok)
	// закрытие закрытого канала вызовет панику
	func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recovered from: ", r)
			}
		}()
		close(ch)
	}()

	// канал может быть буфферизированным
	//  buffCh := make(chan int, 10)
	// небуфферизированный
	// ch := make(chan int)
	// ch := make(chan int, 0)
	// при отправке данных  в небуферизированный канал горутина, внутри которой
	// отправляются данные будет заблокированна до вычитки другой горутиной и наоборот:
	// если данные не переданы данные в вычитывающую горутину, она заблокируется до отправки данных в канал

	// связь по небуферизированному каналу приводит к "синхронизации" операций отправления и получения.
	// по этому связь по такому каналу называют синхронной

	// из каналов можно строить конвееры, где значения передаются через канал из стартовой ф-ции в последующую
	// startFunc(ch) -> processFunc(ch) -> finalFunc(ch)
	// конвееры часто используются в "долгоиграющих" серверах, для организации связи внутри горутин,
	// которые имеют бесконечные циклы

	// вычитывать из  каналов можно при помощи range

	// закрывать каналы всегда по завершению работы с ним не нужно
	// закрывать канал стоит тогда, когда важно сообщить принимающей горутине, что все данные уже отправлены

	// каналы могут иметь направление
	// chSend := make(<-chan int) // на получение
	// канал для получения закрыть может только отправляющая горутина
	// close(chSend) будет ошибка рантайма

	printLineTrace()
	// каналы могут быть буфферизированными
	// при записи в буф канал значение добавляется в конце очереди
	// при вычитке - удаляется из очереди
	buffCh := make(chan int, 3)
	buffCh <- 1
	buffCh <- 2
	// len() - сколько элементов сейчас внутри
	// cap() - вместимость канала
	fmt.Println(len(buffCh), cap(buffCh))
	// буфф каналы позволяют горутинам не блокироваться при записи значения
	// пока буфер не будет заполнен
	// как физический аналог работы буфферизированного канала можно представить
	// как кондитеров, которые пекут торт
	// один делает заготовку, другой мажет крем, третий выкладывает фрукты
	// если между каждым из них есть место что бы ставить торт, выполнив свою часть работы и браться за следующий торт
	// то они будут работать быстрее

	// буферизированные каналы особо эффективны если в каждой ноде конвеера, через который проходит передача данных
	// по таким каналам работа выполняется примерно с одинаковой скоростью то работа будет эфективна, однако, если
	// в начале или конце работа будет происходить быстрее, то смысла работы не будет, т.к. буфер окажется забит и
	// связанные по каналу горутины будет простаивать
	}

func printChanInt(c chan int) {
	fmt.Println(c)
}

func printStruct(s interface{}) {
	fmt.Printf("%+v\n", s)
}

func printSliceInt(s []int) {
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(s), cap(s), &s, s)
}

func printSliceString(s *[]string) {
	fmt.Printf("Len: %d | cap: %d | pointer: %p | elements: %+v\n", len(*s), cap(*s), s, s)
}

func printMapStringString(m map[string]string) {
	fmt.Printf("len: %d; value: %+v\n", len(m), m)
}

// TODO разобраться, как печатать комменты выше вызываемой ф-ции
func FuncPathAndName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func FuncName(f interface{}) string {
	splitFuncName := strings.Split(FuncPathAndName(f), ".")
	return splitFuncName[len(splitFuncName)-1]
}

// Get description of a func
func FuncDescription(f interface{}) string {
	fileName, _ := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).FileLine(0)
	funcName := FuncName(f)
	fset := token.NewFileSet()

	// Parse src
	parsedAst, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
		return ""
	}

	pkg := &ast.Package{
		Name:  "Any",
		Files: make(map[string]*ast.File),
	}
	pkg.Files[fileName] = parsedAst

	importPath, _ := filepath.Abs("/")
	myDoc := doc.New(pkg, importPath, doc.AllDecls)
	for _, theFunc := range myDoc.Funcs {
		if theFunc.Name == funcName {
			return theFunc.Doc
		}
	}
	return ""
}

func printLineTrace() {
	if !isTraceEnabled {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d\n", fn, line)
}
