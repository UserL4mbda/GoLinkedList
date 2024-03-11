
package main

import (
	"fmt"
)

type ALL interface{}
type Nextable  func()         *Node
type Visitor   func(ALL)
type Mappable  func(ALL)      ALL
type Valuable  func()         ALL
type Predicate func(ALL)      bool
type Foldable  func(ALL, ALL) ALL

type Node interface{
	Value() interface{}
	Next()             *Node
	Foreach(Visitor)   *Node
	Map(Mappable)      *Node
	While(Predicate)   *Node
	Filter(Predicate)  *Node
	Cons(ALL)          *Node
	Merge( *Node)      *Node
	Fold(ALL, Foldable) ALL
	Reduce(Foldable)    ALL
}

func Node_Merge(pFirstNode *Node, pSecondNode *Node) *Node {
	tmp := MergedNode {
		firstNode: pFirstNode,
		secondNode: pSecondNode,
	}
	node := Node(tmp)
	return &node
}

func Node_Reduce(pNode *Node, f Foldable) ALL {
	if pNode == nil {
		panic("Node_Reduce pNode ne peut etre null")
	}

	node :=  *pNode
	nextnode := node.Next()
	value := node.Value()

	if nextnode == nil {
		return value
	}

	return node.Fold(value, f)
}

func Node_Fold(pNode *Node, accu ALL, f Foldable) ALL {
	if pNode == nil {
		panic("Node_Fold pNode ne peut etre null")
	}

	node := *pNode
	for{
		value := node.Value()
		accu = f(accu, value)
		pnext := node.Next()

		if pnext == nil {
			return accu
		}

		node = *pnext
	}
}

func Node_Map(ptrnode *Node, f Mappable) *Node {
	tmp := MappedNode {
		supportNode: ptrnode,
		transformer: f,
	}
	node := Node(tmp)
	return &node
}

func Node_Filter(pnode *Node, p Predicate) *Node {
	//ATTENTION: pour les FilterNode, il faudrait implementer
	//un systeme de cache
	if pnode == nil {
		//panic("Error: Node_Filter, pnode can't be nil")
		return nil
	}

	for{
		v := (*pnode).Value()
		if p(v) {
			tmp := FilterNode {
				supportNode: pnode,
				predicate: p,
			}
			node := Node(tmp)
			return &node
		}
		pnode = (*pnode).Next()
		if pnode == nil {
			return nil
		}
	}
}

func Node_While(ptrnode *Node, p Predicate) *Node {
	if ptrnode == nil {
		//panic("Error: Node_While, prtnode can't be nil")
		return nil
	}

	//If the first value doesn't match return nil
	value := (*ptrnode).Value()
	if !p(value) {
		return nil
	}

	tmp := WhileNode {
		supportNode: ptrnode,
		predicate: p,
	}
	node := Node(tmp)
	return &node
}

func Node_Foreach(ptrnode *Node, f Visitor) *Node {
	if ptrnode == nil {
		//panic("Error: Node_While, prtnode can't be nil")
		return nil
	}
	i := ptrnode

	for {
		if i == nil {
			break
		}
		f((*i).Value())
		i = (*i).Next()
	}

	return ptrnode
}

func Node_Cons(ptrnode *Node, value ALL) *Node {
	n := Node(*ptrnode)
	cn := ConcreteNode {
		value: value,
		next:  &n,
	}
	node := Node(cn)
	return &node
}


type SimpleNode struct {
	value Valuable
	next  Nextable
}

func (sn SimpleNode) Value() interface {} {
	return sn.value()
}

func (sn SimpleNode) Next() *Node {
	return sn.next()
}

func (sn SimpleNode) Foreach(f Visitor) *Node {
	node := Node(sn)
	return Node_Foreach(&node, f)
}

func (sn SimpleNode) Map(f Mappable) *Node {
	node := Node(sn)
	return Node_Map(&node, f)
}

func (n SimpleNode) Reduce(f Foldable) ALL {
	node := Node(n)
	return Node_Reduce(&node, f)
}

func (n SimpleNode) Fold(accu ALL, f Foldable) ALL {
	node := Node(n)
	return Node_Fold(&node, accu, f)
}

func (sn SimpleNode) While(p Predicate) *Node {
	node := Node(sn)
	return Node_While(&node, p)
}

func (n SimpleNode) Filter(p Predicate) *Node {
	node := Node(n)
	return Node_Filter(&node, p)
}

func (n SimpleNode) Merge(other *Node) *Node {
	node := Node(n)
	return Node_Merge(&node, other)
}


func (n SimpleNode) Cons(value ALL) *Node {
	node := Node(n)
	return Node_Cons(&node, value)
}

func Constant(k int) Node {
	var tmp SimpleNode
	tmp = SimpleNode {
		value: func () ALL {
			return k
		},
		next: func () *Node {
			var node Node
			node = Node(tmp)
			return &node
		},
	}
	return Node(tmp)
}

func IntegersNodeFrom(from int) Node {
	var tmp SimpleNode
	//init := 0
	//var getter func (int) (func () ALL, func () *Node)
	var gValue func (int) (func () ALL)
	var gNext  func (int) (func () *Node)

	gValue = func (i int) (func () ALL) {
		return func () ALL {
			return i
		}
	}

	gNext  = func (i int) (func () *Node) {
		return func () *Node {
			sn := SimpleNode {
				value: gValue(i),
				next:  gNext(i+1),
			}
			n := Node(sn)
			return &n
		}
	}

	tmp = SimpleNode {
		value: gValue(from),
		next:  gNext(from + 1),
	}

	return Node(tmp)
}

func IntegersNode() Node {
	return IntegersNodeFrom(0)
}

type MappedNode struct {
	supportNode *Node
	transformer Mappable
}

func (mn MappedNode) Value() interface{} {
	return mn.transformer((*(mn.supportNode)).Value())
}

func (mn MappedNode) Next() *Node {
	if mn.supportNode == nil {
		return nil
	}
	nxt := (*(mn.supportNode)).Next()
	if nxt == nil {
		return nil
	}
	nod := Node(MappedNode {
		supportNode: (*(mn.supportNode)).Next(),
		transformer: mn.transformer,
	})
	return &nod
}


func (mn MappedNode) Map(f Mappable) *Node {
	node := Node(mn)
	return Node_Map(&node, f)
}

func (mn MappedNode) Reduce(f Foldable) ALL {
	node := Node(mn)
	return Node_Reduce(&node, f)
}

func (n MappedNode) Fold(accu ALL, f Foldable) ALL {
	node := Node(n)
	return Node_Fold(&node, accu, f)
}

func (mn MappedNode) While(p Predicate) *Node {
	node := Node(mn)
	return Node_While(&node, p)
}

func (n MappedNode) Filter(p Predicate) *Node {
	node := Node(n)
	return Node_Filter(&node, p)
}

func (mn MappedNode) Foreach(f Visitor) *Node {
	node := Node(mn)
	return Node_Foreach(&node, f)
}

func (n MappedNode) Merge(other *Node) *Node {
	node := Node(n)
	return Node_Merge(&node, other)
}

func (n MappedNode) Cons(value ALL) *Node {
	node := Node(n)
	return Node_Cons(&node, value)
}

type FilterNode struct {
	supportNode *Node
	predicate Predicate
}

func (n FilterNode) Value() interface{} {
	return (*(n.supportNode)).Value()
}

func (n FilterNode) Next() *Node {
	if n.supportNode == nil {
		return nil
	}

	predicate := n.predicate
	//WARNING: En doublon dans Node_Filter
	//TODO: Factoriser cette partie
	node := *(n.supportNode)
	pnext := node.Next()
	for{
		if pnext == nil {
			return nil
		}

		nextNode  := *pnext
		nextValue := nextNode.Value()

		if predicate(nextValue) {
			tmp := FilterNode {
				supportNode: pnext,
				predicate: predicate,
			}
			rNode := Node(tmp)
			return &rNode
		}
		pnext = nextNode.Next()
	}
}

func (n FilterNode) Foreach(f Visitor) *Node {
	node := Node(n)
	return Node_Foreach(&node, f)
}

func (n FilterNode) Map(f Mappable) *Node {
	node := Node(n)
	return Node_Map(&node, f)
}

func (n FilterNode) Reduce(f Foldable) ALL {
	node := Node(n)
	return Node_Reduce(&node, f)
}

func (n FilterNode) Fold(accu ALL, f Foldable) ALL {
	node := Node(n)
	return Node_Fold(&node, accu, f)
}

func (n FilterNode) Merge(other *Node) *Node {
	node := Node(n)
	return Node_Merge(&node, other)
}

func (n FilterNode) While(p Predicate) *Node {
	node := Node(n)
	return Node_While(&node, p)
}

func (n FilterNode) Filter(p Predicate) *Node {
	node := Node(n)
	return Node_Filter(&node, p)
}

func (n FilterNode) Cons(value ALL) *Node {
	node := Node(n)
	return Node_Cons(&node, value)
}

type WhileNode struct {
	supportNode *Node
	predicate Predicate
}

func (wn WhileNode) Value() interface{} {
	return (*(wn.supportNode)).Value()
}

func (wn WhileNode) Next() *Node {
	if wn.supportNode == nil {
		return nil
	}

	ptrnext := (*(wn.supportNode)).Next()
	if ptrnext == nil {
		return nil
	}

	nextnode  := *ptrnext
	value     := nextnode.Value()

	if !(wn.predicate(value)) {
		return nil
	}

	rNode := Node(WhileNode {
		supportNode: ptrnext, //nextnode,
		predicate:   wn.predicate,
	})

	return &rNode
}

func (wn WhileNode) Foreach(f Visitor) *Node {
	node := Node(wn)
	return Node_Foreach(&node, f)
}

func (wn WhileNode) Map(f Mappable) *Node {
	node := Node(wn)
	return Node_Map(&node, f)
}

func (n WhileNode) Reduce(f Foldable) ALL {
	node := Node(n)
	return Node_Reduce(&node, f)
}

func (n WhileNode) Fold(accu ALL, f Foldable) ALL {
	node := Node(n)
	return Node_Fold(&node, accu, f)
}

func (wn WhileNode) While(p Predicate) *Node {
	node := Node(wn)
	return Node_While(&node, p)
}

func (n WhileNode) Filter(p Predicate) *Node {
	node := Node(n)
	return Node_Filter(&node, p)
}

func (n WhileNode) Merge(other *Node) *Node {
	node := Node(n)
	return Node_Merge(&node, other)
}

func (n WhileNode) Cons(value ALL) *Node {
	node := Node(n)
	return Node_Cons(&node, value)
}

type MergedNode struct {
	firstNode  *Node
	secondNode *Node
}

func (n MergedNode) Value() interface{} {
	if n.firstNode != nil {
		return (*(n.firstNode)).Value()
	}

	if n.secondNode != nil {
		return (*(n.secondNode)).Value()
	}

	//We should panic insteed!
	return nil
}

func (n MergedNode) Next() *Node {
	if n.firstNode != nil {
		pnode := (*(n.firstNode)).Next()
		if pnode != nil {
			mn:= MergedNode {
				firstNode: pnode,
				secondNode: n.secondNode,
			}
			node := Node(mn)
			return &node
		}
	}

	if n.secondNode != nil {
		pnode := (*(n.secondNode)).Next()
		if pnode != nil {
			mn:= MergedNode {
				firstNode: nil,
				secondNode: pnode,
			}
			node := Node(mn)
			return &node
		}
	}

	return nil
}

func (n MergedNode) Foreach(f Visitor) *Node {
	node := Node(n)
	return Node_Foreach(&node, f)
}

func (n MergedNode) Map(f Mappable) *Node {
	node := Node(n)
	return Node_Map(&node, f)
}

func (n MergedNode) Reduce(f Foldable) ALL {
	node := Node(n)
	return Node_Reduce(&node, f)
}

func (n MergedNode) Fold(accu ALL, f Foldable) ALL {
	node := Node(n)
	return Node_Fold(&node, accu, f)
}

func (n MergedNode) While(p Predicate) *Node {
	node := Node(n)
	return Node_While(&node, p)
}

func (n MergedNode) Filter(p Predicate) *Node {
	node := Node(n)
	return Node_Filter(&node, p)
}

func (n MergedNode) Cons(value ALL) *Node {
	node := Node(n)
	return Node_Cons(&node, value)
}

func (n MergedNode) Merge(other *Node) *Node {
	node := Node(n)
	return Node_Merge(&node, other)
}

type ConcreteNode struct{
	value interface{}
	next *Node
	//next *ConcreteNode
}

func (n ConcreteNode) Value() interface{} {
	return n.value
}

func (n ConcreteNode) Next() *Node {
	if n.next == nil {
		return nil
	}
	m := Node(*(n.next))
	return &m
}

func (n ConcreteNode) Cons(value ALL) *Node {
	node := Node(n)
	return Node_Cons(&node, value)
}

func (n ConcreteNode) Foreach(f Visitor) *Node {
	node := Node(n)
	return Node_Foreach(&node, f)
}

func (n ConcreteNode) Map(f Mappable) *Node {
	node := Node(n)
	return Node_Map(&node, f)
}

func (n ConcreteNode) Reduce(f Foldable) ALL {
	node := Node(n)
	return Node_Reduce(&node, f)
}

func (n ConcreteNode) Fold(accu ALL, f Foldable) ALL {
	node := Node(n)
	return Node_Fold(&node, accu, f)
}

func (n ConcreteNode) While(p Predicate) *Node {
	node := Node(n)
	return Node_While(&node, p)
}

func (n ConcreteNode) Filter(p Predicate) *Node {
	node := Node(n)
	return Node_Filter(&node, p)
}

func (n ConcreteNode) Merge(other *Node) *Node {
	node := Node(n)
	return Node_Merge(&node, other)
}

type NullableNode struct{
	node *Node
}

type NullableValue struct {
	value *interface {}
	undef func () bool
	//cache *interface {}
}

func (nv NullableValue) Value() interface{} {
	return *(nv.value)
}

func (nv NullableValue) Undef() bool {
	return nv.undef()
}

func (nn NullableNode) Next() NullableNode {
	pnode := nn.node
	if pnode == nil {
		return nn
	}
	pnextnode := Node(*pnode).Next()
	return NullableNode{node: pnextnode}
}

func (nn NullableNode) Value() NullableValue {
	pnode := nn.node
	if pnode == nil {
		return NullableValue{ value: nil, }
	}
	node := Node(*pnode)
	val  := node.Value()
	return NullableValue { value: &val, }
}

func (nn NullableNode) Reduce(f Foldable) NullableValue {
	pnode := nn.node
	if pnode == nil {
		return NullableValue { value: nil, undef: func () bool{return true}, }
	}

	node := *pnode
	value := node.Reduce(f)
	v := interface{}(value)
	return NullableValue {
		value: &v,
		undef: func () bool {return false},
	}
}

func (nn NullableNode) Fold(accu ALL, f Foldable) NullableValue {
	pnode := nn.node
	if pnode == nil {
		return NullableValue { value: nil, undef: func () bool{return true}, }
	}

	node := *pnode
	v := node.Fold(accu, f)
	vw := interface{}(v)

	return NullableValue {
		value: &vw,
		undef: func () bool {return false},
	}
}

func (nn NullableNode) Map(f Mappable) NullableNode {
	pnode := nn.node
	if pnode == nil {
		return nn
	}
	return NullableNode{ node: Node(*pnode).Map(f), }
}

func (nn NullableNode) While(p Predicate) NullableNode {
	pnode := nn.node
	if pnode == nil {
		return nn
	}
	return NullableNode{ node: Node(*pnode).While(p), }
}

func (nn NullableNode) Filter(p Predicate) NullableNode {
	pnode := nn.node
	if pnode == nil {
		return nn
	}
	return NullableNode{ node: Node(*pnode).Filter(p), }
}

func (nn NullableNode) Foreach(f Visitor) NullableNode {
	pnode := nn.node
	if pnode == nil {
		return nn
	}
	inner := Node(*pnode)
	inner.Foreach(f)
	return nn
}

func (nn NullableNode) Cons(value ALL) NullableNode {
	pnode := nn.node
	innerNode := ConcreteNode {
		value: value,
		next: pnode,
	}
	node := Node(innerNode)
	return NullableNode {
		node: &node,
	}
}

func (n NullableNode) Merge(n2 NullableNode) NullableNode {
	pnode := n.node
	node := *pnode

	pnode2 := n2.node
	inner := node.Merge(pnode2)
	return NullableNode {
		node: inner,
	}
}

func NullableNode_New(n *Node) NullableNode {
	return NullableNode {
		node: n,
	}
}

func NodeValue(value ALL) NullableNode {
	n := ConcreteNode {
		value: value,
		next: nil,
	}

	node := Node(n)

	return NullableNode {
		node: &node,
	}
}

func (n NullableNode) Return() *Node {
	return n.node
}

func Integers() NullableNode {
	integers := IntegersNode()
	return NullableNode_New(&integers)
}

func main () {
	fmt.Println("Bonjour le monde !")

	myNode := ConcreteNode{
		value : 3,
		next : nil,
	}


	var printage Visitor = func (v ALL) {
		fmt.Println("PRINTAGE:")
		fmt.Printf("%v\n", v)
	}

//	var double Mappable = func (v ALL) ALL {
//		return v.(int) * 2
//	}
//
//	var sup14 Predicate = func (v ALL) bool {
//		return v.(int) > 14
//	}

//	var inf102 Predicate = func (v ALL) bool {
//		return v.(int) < 102
//	}

	inf := func (i int) Predicate {
		return func(v ALL) bool {
			return v.(int) <= i
		}
	}

	inc := func (i int) Mappable {
		return func(v ALL) ALL {
			return v.(int) + i
		}
	}

	fmt.Println(myNode)

	newNode := myNode.Cons(4)
	fmt.Println(newNode)


	nn := NodeValue(3).Cons(4).Cons(5).Cons(6).Cons(7).Cons(8).Cons(9).Cons(10).Cons(11).Cons(12)

	nm := NodeValue(100).Cons(101).Cons(102).Cons(103).Cons(104).Cons(105).Cons(106).Cons(107)
	fmt.Println(nn)

	nn.Merge(nm).Filter(inf(102)).Foreach(printage)
	k := nn.Merge(nm).Filter(inf(102)).Fold(0, func(accu ALL, value ALL) ALL { return accu.(int) + value.(int) })
	fmt.Println("k =",*(k.value))

	w := NodeValue(5).Cons(4).Cons(3).Cons(2).Cons(1)
	add := func(a ALL, b ALL) ALL { return a.(int) + b.(int) }
	mult:= func(a ALL, b ALL) ALL { return a.(int) * b.(int) }

	//val1 := w.Fold(0, add)
	fmt.Println("5+4+3+2+1 =", w.Reduce(add).Value())

	//val2 := w.Fold(1, mult)
	fmt.Println("5*4*3*2*1 =", w.Reduce(mult).Value())


	fmt.Println(inf(3)(2))

	Integers().While(inf(6)).Foreach(printage)

	fmt.Println("\nDe 1 a 6 potentiellement !")
	Integers().Map(inc(1)).While(inf(6)).Foreach(printage)

	fmt.Println( "Fact(6) =", Integers().Map( inc(1) ).While( inf(6) ).Reduce( mult ).Value() )


	fmt.Println("Test Entiers:")
	nod := IntegersNode()
	fmt.Println("0 ->", nod.Value())

	un := nod.Next()
	uno := *un
	fmt.Println("1 ->", uno.Value())

	deux := uno.Next()
	due := *deux
	fmt.Println("2 ->", due.Value())

	//With nn
//	nn.
//		Foreach(printage).
//		Map(double).
//		Foreach(printage).
//		While(sup14).
//		Foreach(printage).
//		Merge(nm).
//		Foreach(printage)




//	fmt.Println(newNode)
//
//	fmt.Println(newNode.Next())
//
//	newNode.Foreach(func(v ALL){fmt.Printf("Node: %v\n",v)})
//
//	newNode.Cons(5).Foreach(func(v ALL){fmt.Printf("new: %v\n",v)}) //.Foreach(func(w ALL){fmt.Printf("Toto: %v\n",w)})
//
//	nod3 := newNode.Cons(5).Cons(6).Cons(7).Cons(8).Cons(9).Cons(10)
//	fmt.Println("nod3:", nod3)
//
//	pn4 := nod3.Map(func (x ALL) ALL {return x.(int) * 2})
//	fmt.Println("nod3:", nod3)
//
//	n4 := *pn4
//	var printage Visitor = func (v ALL) {
//		fmt.Println("PRINTAGE:")
//		fmt.Printf("%v\n", v)
//	}
//	fmt.Printf("N4: %v\n", n4.Value())
//	fmt.Printf("N4: %v\n", (*(*n4.Next()).Next()).Value())
//	fmt.Printf("N4: %v\n", (*(*(*n4.Next()).Next()).Next()).Value())
//	n4.Foreach(printage)
//
//	cinq := Constant(5)
//	fmt.Printf("Cinq: %v\n", cinq.Value())
//	fmt.Printf("Cinq4: %v\n", (*(*cinq.Next()).Next()).Value())
//
//	entiers := Integers()
//	fmt.Printf("Entiers: %v\n", entiers.Value())
//	fmt.Printf("Entiers4: %v\n", (*(*entiers.Next()).Next()).Value())
//
//	pentiers2 := Integers().Map(func (n ALL) ALL {return n.(int) * 2})
//	entiers2 := *pentiers2
//	fmt.Printf("Entiers2: %v\n", entiers2.Value())
//	fmt.Printf("Entiers8: %v\n", (*(*entiers2.Next()).Next()).Value())
//
//	pwEntier := Integers().While(func (n ALL) bool {if n.(int) < 12 {return true}; return false})
//	wEntier := *pwEntier
//	wEntier.Foreach(printage)
//
//	fmt.Println("nod3:", nod3)
//	//fmt.Println("NOD3")
//	nod3.Foreach(printage)
//
//	ptst := nod3.While(func(n ALL)bool{if n.(int)>4 {return true};return false})
//	if ptst != nil {
//		fmt.Println("node3 > 4 printage")
//		(*ptst).Foreach(printage)
//	}
//
//	//fmt.Println("NOD3")
//	//nod3.While(func (n ALL) bool {if n.(int) < 8 {return true}; return false}).Foreach(printage)
//
}
