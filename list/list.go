package list

/**
 *   Copyright (c) Marcin Radoszewski. All rights reserved.
 *   The use and distribution terms for this software are covered by the
 *   Eclipse Public License 1.0 (http://opensource.org/licenses/eclipse-1.0.php)
 *   which can be found in the file epl-v10.html at the root of this distribution.
 *   By using this software in any fashion, you are agreeing to be bound by
 * 	 the terms of this license.
 *   You must not remove this notice, or any other, from this software.
 **/

import (
	. "github.com/marad/gods"
)

/**
 * Implementation of persistent linked list.
 **/

type List struct {
	head Value
	tail *List
}

var EMPTY_LIST = &List{nil, nil}

func Empty() *List {
	return EMPTY_LIST
}

func FromArray(arr []Value) *List {
	list := EMPTY_LIST
	for i := len(arr) - 1; i >= 0; i-- {
		list = list.Cons(arr[i])
	}
	return list
}

func (l *List) Cons(val Value) *List {
	return &List{val, l}
}

func (l *List) First() Value {
	return l.head
}

func (l *List) Rest() *List {
	return l.tail
}

func (l *List) IsEmpty() bool {
	return l.head == nil && l.tail == nil
}

func (l *List) Copy() *List {
	var current *List
	var newHead *List
	var newCurrent *List

	if l.IsEmpty() {
		return EMPTY_LIST
	} else {
		newHead = EMPTY_LIST.Cons(l.First())
		newCurrent = newHead
		current = l.Rest()
		for !current.IsEmpty() {
			newCurrent.tail = EMPTY_LIST.Cons(current.First())
			newCurrent = newCurrent.Rest()
			current = current.Rest()
		}
		return newHead
	}
}

func (l *List) Insert(element Value, position int) *List {
	var current *List
	var newHead *List
	var newCurrent *List

	if l.IsEmpty() {
		return EMPTY_LIST.Cons(element)
	} else {
		newHead = EMPTY_LIST.Cons(l.First())
		newCurrent = newHead
		current = l.Rest()
		for copied := 1; copied < position && !current.IsEmpty(); copied++ {
			newCurrent.tail = EMPTY_LIST.Cons(current.First())
			newCurrent = newCurrent.Rest()
			current = current.Rest()
		}

		newCurrent.tail = current.Cons(element)
		return newHead
	}
}
