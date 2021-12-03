package klepman

import (
	"container/list"
	"fmt"
	"runtime"
)

const (
	isTraceEnabled = true
)

// https://medium.com/nuances-of-programming/%D0%B2%D1%81%D0%B5-%D1%87%D1%82%D0%BE-%D0%BD%D1%83%D0%B6%D0%BD%D0%BE-%D0%B7%D0%BD%D0%B0%D1%82%D1%8C-%D0%BE-%D0%B4%D1%80%D0%B5%D0%B2%D0%BE%D0%B2%D0%B8%D0%B4%D0%BD%D1%8B%D1%85-%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D0%B0%D1%85-%D0%B4%D0%B0%D0%BD%D0%BD%D1%8B%D1%85-d750444a77ec

func Run() {
	// инит дерева вида
/*
	       1
	    /     \
       2       5
      / \     / \
	 3   4   6   7
*/
	printLineTrace()

	root := NewRoot("1")
	n2 := root.root.
		InsertLeft("2")

	n2.InsertLeft("3")
	n2.InsertRight("4")

	n5 := root.root.
		InsertRight("5")

	n5.InsertLeft("6")
	n5.InsertRight("7")

	// DFS обход (поиск в глубину) :
	// - тип PreOrder - предварительный поиск
	printLineTrace()
	root.root.DFS(PreOrderDFS)

	// - тип InOrder  - симметричный поиск
	printLineTrace()
	root.root.DFS(InOrderDFS)

	// - тип PostOrder  - как InOrder, но в обратном порядке
	printLineTrace()
	root.root.DFS(PostOrderDFS)

	// BFS()
	printLineTrace()
	root.root.BFS()
}

type BinaryTree struct {
	root *BinaryNode
}

type BinaryNode struct {
	Left *BinaryNode
	Right *BinaryNode
	Value string
}

func NewRoot(v string) *BinaryTree {
	return &BinaryTree{root: &BinaryNode{Value: v}}
}

func (n *BinaryNode) InsertLeft(v string) *BinaryNode{
	if n.Left == nil {
		n.Left = &BinaryNode{Value: v}
	} else {
		newNode := &BinaryNode{Value: v}
		newNode.Left = n.Left
		n.Left = newNode
	}
	return n.Left
}

func (n *BinaryNode) InsertRight(v string) *BinaryNode{
	if n.Right == nil {
		n.Right = &BinaryNode{Value: v}
	} else {
		newNode := &BinaryNode{Value: v}
		newNode.Right = n.Right
		n.Right = newNode
	}

	return n.Right
}

// DFS Depth-first search ПОИСК В ГЛУБИНУ
// Проход в глубь дерева, а затем возврат к исходной точке называется алгоритмом DFS
// метод обхода дерева
// - рекурсивно перебираются все исходящие из рассматриваемой вершины ребра
// - если ребро ведет в нерассмотренную вершину, запускаем алгоритм от этой вершины
// - возврат в том случае, если в не осталось ребер в текущей вершине
// - если после завершения остались нерассмотренные вершины, то запускаем алгоритм на нерассмотренной

// DfsOrder Типы DFS поиска
type DfsOrder int
const (
	// PreOrderDFS - предварительный обход
	// 1. Записать значение узла
	// 2. Если есть левый потомок перейти и записать значение
	// 3. п2 для правого
	PreOrderDFS DfsOrder = iota
	// InOrderDFS Симметричный обход
	// 1. Перейти к левому, записать
	// 2. Записать
	// 3. Перейти к правому, записать
	InOrderDFS
	// PostOrderDFS Обход в обратном порядке
	// 1. Перейти к левому потомку и записать
	// 2. Перейти к правому и записать
	// 3. Записать
	PostOrderDFS
)

func(n *BinaryNode) DFS(order DfsOrder) {
	switch order {
	case PreOrderDFS:
		dfsPreOrder(n)
	case InOrderDFS:
		dfsInOrder(n)
	case PostOrderDFS:
		dfsPostOrder(n)
	}
}

func dfsPreOrder(n *BinaryNode) {
	if n == nil {
		panic("node is nil")
	}
	// сразу выполнит действие с нодой
	fmt.Println(n.Value)

	// будет лезть влево, и делать действие
	// пока не наткнется на самый глубокий узел у которого нет левого ребенка
	if n.Left != nil {
		dfsPreOrder(n.Left)
	}

	// когда добрались в левый самый глубокий узел пошли по его правой части
	if n.Right != nil {
		dfsPreOrder(n.Right)
	}

	// понятно, что все рекурсивные вызовы накапливаются в стеке
	// и раскручиваются после выхода из самого последнего вызова
	return
}

func dfsInOrder(n *BinaryNode) {
	if n == nil {
		panic("node is nil")
	}

	// сперва ищем самый глубокий левый узел
	if n.Left != nil {
		dfsInOrder(n.Left)
	}

	// ток нашли - напечатали
	fmt.Println(n.Value)

	// потом пошли вправо
	if n.Right != nil {
		dfsInOrder(n.Right)
	}

	return
}

func dfsPostOrder(n *BinaryNode) {
	if n == nil {
		panic("node is nil")
	}

	// пошли в левый самый глубокий
	// и нашли его
	if n.Left != nil {
		dfsPostOrder(n.Left)
	}

	// пошли в самый глубокий правый самого глубокого левого
	if n.Right != nil {
		dfsPostOrder(n.Right)
	}

	// и напечатли его
	fmt.Println(n.Value)
}

// BFS Breadth-first search
// обходит по уровням
// пример:
/*
	       1        - lvl 0
	    /     \
       2       5    - lvl 1
      / \     / \
	 3   4   6   7  - lvl 2
*/
// сам алгоритм изпользует очередь
// 1. Добавить рут в очередь (put)
// 2. Повторять, пока очередь не пуста
// 3. Получить первый узел в очереди и записать его значение
// 4. Добавить левый и правый потомок в очередь
func(n *BinaryNode) BFS() {
 	q := list.New()
 	q.PushBack(n)

 	for q.Len() > 0 {
		currEl := q.Front()
		curr, ok := currEl.Value.(*BinaryNode)
		if !ok {
			panic("error assert BFS.List.Front.Value to *BinaryNode")
		}
		q.Remove(currEl)

		fmt.Println(curr.Value)

		if curr.Left != nil {
			q.PushBack(curr.Left)
		}

		if curr.Right != nil {
			q.PushBack(curr.Right)
		}
	}
}


// Бинарное дерево поиска
// упорядоченное бинарное дерево
// хранит значения таким образом, что может быть пригодно
// для принципов бинарного поиска
// важное свойство это то, что величина узла дерева бинарного поиска больше,
// чем кол-во его потомков левого элемента потомка, но меньше, чем кол-во его потомков
// правого элемента потомка.


func printLineTrace() {
	if !isTraceEnabled {
		return
	}
	_, fn, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d\n", fn, line)
}