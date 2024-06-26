package klepman

// Модели данных
/*
relational database management system, RDBMS
Реляционная база данных – это набор данных с предопределенными связями между ними

Рассогласование (impedance mismatch) - ситуация, когда
модели данных ооп не соотвествуют модели данных таблиц
т.е. нужен какой то промежуточный слой для преобразования
данных в RDBMS модель.

Локальность
У JSON-представления — лучшая локальность, чем у многотабличной схемы.
Т.е. все представление модели хранится внутри одной сущности,
в отличие от RDBMS, где данные одной сущности  могут быть разнесены по несокльим
таблицам.

Например, связь один ко многим по сути является древовидной структурой
а json делает это наиболее явно.

https://habr.com/ru/post/254773/
Нормализация базы данных
Способ организации базы данных таким образом, что бы сократить
кол-во дублирования информации
либо сам процесс удаления избыточных данных
при нормализации решаются проблемы производительности, аномалий,
удобства управления данными.
Нормальная форма(НФ) базы данных - это набор правил и критериев, которым должна
отвечать база данных.
Чем выше НФ тем сторже ограничения.
Самой оптимальной считается 3НФ, т.к это компромисс между скоростью и удобством.
Есть и другие (6НФ), но скорее всего они не понадобятся в реальной жизни

Требования 1 НФ
- в таблице нет дублирующих строк
- все значения атомарны (не составные)
- в каждом столбце данные одного типа
- нет массивов и списков в любом виде
по сути, эти требования и являются требованиями реляционной модели

2НФ
- 1НФ
- в таблице должен быть ключ (ключ однозначно идентифицирует строку таблицы)
- если ключ составной, нельзя получить строку по части ключа
Для достижения второй нормальной формы можно использовать декомпозицию таблиц
Т.е. вынесение повторяющихся данных в отдельную таблицу, чтобы затем связать
данные по идентификаторам.

3НФ - отсутсвие транзитивной зависимости
- 2НФ
- Таблица должна содержать правильные неключевые столбцы
(транзитивность - зависимость неключевых столбцов от других неключевых)
по неключевым столбцам нельзя получить данные из других столбцов
т.е. когда каждое поле, которое можно "обобщить" - обобщается в отдельную таблицу.

3НФ+ (форма Бойса-Кода)
таблицы в 3НФ без отсутствия составного ключа автоматом в 3НФ+
- 3НФ
- ключевые атрибуты составного ключа не должны зависеть от неключевых атрибутов.
(непонятно). Скорее всего имеется ввиду декомпозиция составного ключа

ВЫВОД ПО ПОВОДУ НОРМАЛИЗАЦИИ
если есть значения, которые могут храниться в одном месте, но дублируются,
то схема ненормализованна

О ВЫБОРЕ МЕЖДУ ДОКУМЕНТООРИЕНТИРОВАННОЙ БД И РЕЛЯЦИОННОЙ

Если структура данных представляет дерево один ко многим
и все дерево обычно загружается сразу, то лучше брать носкуль.

у документной модели есть свои ограничения - нельзя ссылаться на вложенный
элемент внутри самого документа, вместо этого придется описать путь к нему

MAP REDUCE

модель программирования, для обработки больших объемов данных блоками
на множестве машин.
В урезаном виде есть в nosql (mongo, couch) как механизм, выполняющий чтение
по многим документам.
MapReduce - среднее между декларативным и недекларативным языком.
Логика запроса выражается при помощи фрагментов кода, которые вызывает
фреймворк.
основан на map, reduce (должны быть чистыми, это позволит выполнить их
 где угодно)

map - используется для применения ф-ции к каждому элементу итерируемого
      интерфейса
reduce - используется для последовательной обработки каждого элемента
         массива с сохранением результата


этот подход похож на конвеер элементов, где последовательно к каждому применяется
какая-то ф-ция для выборки.

MAP REDUCE - для распределенных вычислений на кластерах машин.
Языки запросов более высокого уровня можно реализовать в виде конвейера
операций MAP REDUCE.

Графовые модели данных
- используются для моделирования социальных графов
- веб страницы (ноды - веб страницы, ребра ссылки на другие страницы)
- дороги и жд сети

Для работы с графами есть известные механизмы.
В модели графов свойств каждая вершина состоит из:
- уникального идентификатора;
- множества исходящих ребер;
- множества входящих ребер;
- коллекции свойств (пар «ключ — значение»).

ИНДЕКСЫ БД
// https://ru.wikipedia.org/wiki/%D0%98%D0%BD%D0%B4%D0%B5%D0%BA%D1%81_(%D0%B1%D0%B0%D0%B7%D1%8B_%D0%B4%D0%B0%D0%BD%D0%BD%D1%8B%D1%85)
Индекс дополнительная структура БД, создаваемая с целью
повышения производительности поиска данных.
Обычно структура индекса это
(значение одного/нескольких столбцов):(указатель на строку таблицы)
оптимизированная под какой-либо вид поиска (диапазон/вхождение/просто выборка, агрегация)
Аналог в реальном мире - каталог в библиотеке

Индексы бывают двух типов
Кластреные
	строки таблицы упорядочены по значению ключа индекса
    кластерный индекс может быть только один в рамках таблицы

Некластреные
    если нет кластерного индекса - таблицу называют кучей
    содержит только указатели на записи таблицы
    может быть несколько в таблице, каждый определяет собственный порядок следования записей

Индексы создают на тех столбцах, которые часто используются в запросах
увеличение числа индексов замедляет операции добавления, обновления, удаления,
т.к. эти операции требуют обновления индексов.
+ индексы занимают память, по этому нужно убедиться в профите перед созданием индекса

Разреженный индекс
(spare index) каждый ключ ассоциируется с определенным указателем на БЛОК
в сортированном файле данных

Плотный индекс
(dense index) каждый ключ ассоциируется с указателем на запись в сортированном
файле данных

В кластреном индексе с дубированными ключами разреженный индекс
указывает на наименьший ключ в каждом блоке
а плотный указывает на первую запись с указанным ключем

ХЕШ ИНДЕКСЫ - индексы для типа данных ключ-значение.
Задача индекса - быстро находить TID (ссылка на строку таблицы)

ХЕШ таблица
Хороша тем, что скорость доступа к данным (вставка, удалвение, поиск) осуществляется
за О(1) при условии хорошей реализации
но побочным эффектом является скорость доступа к данным и потребление памяти.
Коллизии могут решаться двумя путями - созданием цепочек и двойным хешированием.

Максимальныое значение хеша = максимальному размеру таблицы

т.е. в таблице размером 10 максимальное значение хеша = 9 (0, 1 ..., 9)
Двойное хеширование
- используются две независимые кеш ф-ции к ключам
s = h1(k) - возвращает натурельное число s
t = h2(k) - возвращает шаг, который будет указывать место, куда поставить элемент, если если h1(k) существует
k - ключ
m - размер таблицы
n mod m - остаток от деления n на m

т.е. сначала будет рассматриваться элемент s, если занято, то на s + t, если найдено то на s + 2*t

Метод цепочек (как в го - бакеты)
i - хеш код элемента
h - массив
h[i] - это указатель на начало списка элементов, хеш код которых равен i

Простейшая индексация для случая добавления строки в конец файла
- хранить в ОПЕРАТИВКЕ хеш карту(журнал) ключ:смещение(адрес в файле базы данных на диске)

При добавлении новой пары происходит и обновление хеш карты
Для экономии памяти используют сегменты определенного размера
При достижении размера записыват данные в новый файл
УПОЛОТНЕНИЕ - запись только послдений версии дублирующихся данных для каждого ключа.
              Может применяться к сегментам.
			  Уплотнение можно проводить и в отдельном потоке.

В ПОСТГРЕС:
- хеш ф-ция всегда возвращает integer т.е. до 4 млрд значений
- организованны странично:
     метастраница описывающая то, что именно внутри индекса
     страницы корзин - хранят данные в ввиде хешкод:tid
	 overflow page	 - устроены как страницы корзин, используются, когда одной страницы не хватает
     bitmap          - отмечаются освбодившиеся страницы переполнения, которые можно использовать

недостатки хеш индексов:
- хеш индекс должен полностью помещаться в оперативке
  из-за большого кол-ва операций ввод-выдод с произвольным
  доступом(обращение к любой записи внутри файла)
- запросы по диапазону не эффективны, т.к. нужно искать каждый ключ
  отдельно в хеш картах


недостатки в ПОСТГРЕС:
до 10 версии:
- действия не попадают в журнал упреждющей записи, по этому не могут быть восстановленны после сбоя
  и не участвуют в репликации
- уступает Б-дереву по универсальности

По этому до 10 версии использовать хеш индексы не имеет смысла, но они плотно используются
для внутренней оптимизации запросов

SS-ТАБЛИЦЫ 卍
sorted string table
для ss нужно, что бы ключ встречался только один раз
в каждом объединенном файле сегмента (уплотнение) и был отсортирован по ключу
Преимущество перед хеш
 - объединение сегмнетов происходит быстро, даже если размер сегментов больше
   чем размер оперативки т.к. алгоритм подобен merge sort
   1. Одновременно читаются входные файлы
   2. Просматривается первый ключ в каждом файле
   3. Копируется нижний ключ в выходной файл
   4. Действие повторяется
Если ключи в сегментах дублируются - то берется самый новый.

Алгоритм работы с SS таблицами:
1. При записи добавляем сперва в сбалансированную структуру данных(например чб дерево)
    эта структура лежит в оперативке и называется MemTable
2. Если MemTable (несколько МБ) переполняется, то записываем на диск в виде ss table.
	Это происходит быстро, т.к. данные это отсортированный kv. Эта таблица становится последним
    сегментом БД, пока происходит запись на диск, новые записи летят в новую MemTable
3. Для запроса данных сперва ключ ищут в MemTable, затем в последнем по времени сегменте диска
4. В фоне запускается периодический процесс слияния и уплотнения.

Недостаток схемы в том, что при сбое данные пропадают из оперативки.
Что бы решить это можно держать отдельный журнал в котором есть все записи
из которых можно восстановить мем тейбл. При записи новой мем тейбл на диск журнал можно удалять.

такая индексная структура называется LSM-TREE.
Подсистемы хранения, основанные на слиянии и уплотнении называются
LSM-подсистемами хранения.

в LSM-TREE проблематичен поиск отсутсвтующих значений, т.к. нужно перебрать всю таблицу.
Для оптимизации этого используют фильтры Блума, который позволяет
примерно определить содержимое множества) без многих обращений к диску.

LSM подходит для поиска по диапазону, т.к. отсортирован и поддерживать высокую скорость
записи т.к. записи на диск происходят последовательно.


Теория по деревьям:
Корень  - самый верхний узел дерева
Ребро   - связь
Лист    - узел, который не имеет потомков
Высота  - дерева считается от корня до самого дальнего листа
Глубина - длинна от конкретного узла к корню

Бинарное дерево
Структура данных, в которой каждый узел имеет не более двух потомков
узел называется родительским
а дети - наследниками

B-TREE

Самый универсальный индекс
тоже KV
Эффективны для запросов по ключу и диапазонам
Индекс разбивает бд на блоки(страницы) фиксированного размера, обычно 4кб
       и записывает по одной странице за раз
Все страницы имеют свой адрес и могут ссылаться на другие страницы
как указатели, но только на диске.
Эти ссылки могут огранизовываться в деревья

Одна из страниц назначается корнем дерева, с нее и начинается поиск.
Страница содержит несколько ключей и ссылок на дочерние страницы,
каждая страница содержит непрерывный диапазон ключей
а ключи между ссылками указывют на расположение границ диапазонов.

Кол-во ссылок на дочерние страницы б-дерева называется коэф. ветвления.
*/