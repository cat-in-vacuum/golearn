package os

/*
процесс - это исполняемая программа
Структура называется process control block
состоит из самой программы
данных (указатель на данные)
контекста выполнения

ПРОЦЕСС

Идентификатор
     уникальное имя процесса в ос что бы к нему обратиться
Состояние
     описание состояние процесса(запущен, остановлен и т.д.)
Приоритет
     каждый процесс можно оценить (как более или менее важный)
		обычно это число
Счетчик команд
     показывает на текущий или следующий процесс
Указатель на память
     указатель на адрес памяти
Контекст
     информация которая позволяет вернуть процесс в работу
     после паузы
Информация о IO
Другая информация

МОДЕЛЬ СОСТОЯНИЙ ПРОЦЕССА

--- процесс не запущен --- процесс запущен ---- процесс не запущен
вход  -------> запуск --->  работа -----------> пауза или выход --

При запуске процесс становится в очередь (с приоритетами)

СОЗДАНИЕ ПРОЦЕССА
-- создание структуры процесса
-- обычно процессы создает ОС, но может созадаваться и процессами
  process spawning:
  	родительский --> дочерний

УНИЧТОЖЕНИЕ ПРОЦЕССА
-- нужно понимание, когда процесс завршен
-- сигнал HALT
-- действие пользователя
-- Ошибка
-- Заврешние родительского процесса

*/


