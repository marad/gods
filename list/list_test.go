package list

import (
	. "github.com/smartystreets/goconvey/convey"
	. "gods"
	"testing"
)

func TestCreatingEmptyList(t *testing.T) {
	var l *List
	Convey("Given an empty list", t, func() {
		l = Empty()

		Convey("Expect to have the list empty", func() {
			So(l.IsEmpty(), ShouldEqual, true)
		})

		Convey("First and Rest on empty list should return nil", func() {
			So(l.First(), ShouldEqual, nil)
			So(l.Rest(), ShouldEqual, nil)
		})
	})
}

func TestCeratingListFromArray(t *testing.T) {
	var array []Value
	var list *List

	Convey("Given an array", t, func() {
		array = []Value{1, 2, 3}

		Convey("When creating list from array", func() {
			list = FromArray(array)

			Convey("List should not be empty", func() {
				So(list.IsEmpty(), ShouldEqual, false)
			})

			Convey("Elements should match", func() {
				So(list.First(), ShouldEqual, 1)
				So(list.Rest().First(), ShouldEqual, 2)
				So(list.Rest().Rest().First(), ShouldEqual, 3)
			})

			Convey("Last element should have an empty list as Rest", func() {
				So(list.Rest().Rest().Rest(), ShouldEqual, EMPTY_LIST)
			})
		})
	})
}

func TestInsertingAndGettingFirstElement(t *testing.T) {
	var empty *List
	var one *List
	var two *List

	Convey("Given an empty list", t, func() {
		empty = Empty()

		Convey("When adding new element", func() {
			one = empty.Cons(1)

			Convey("Result should be different than empty list", func() {
				So(one, ShouldNotEqual, empty)
			})

			Convey("First should return added element", func() {
				So(one.First(), ShouldEqual, 1)
			})

			Convey("Empty should not be modified", func() {
				So(empty.IsEmpty(), ShouldEqual, true)
			})
		})

		Convey("And when adding another element", func() {
			two = one.Cons(2)

			Convey("Result should be different", func() {
				So(two, ShouldNotEqual, empty)
				So(two, ShouldNotEqual, one)
			})
		})
	})
}

func TestCopyingList(t *testing.T) {
	var list *List
	var copy *List

	Convey("Given a list", t, func() {
		list = FromArray([]Value{1, 2, 3})

		Convey("When copying the list", func() {
			copy = list.Copy()

			Convey("List should not be empty", func() {
				So(copy.IsEmpty(), ShouldEqual, false)
			})

			Convey("Copy should be different list", func() {
				So(copy, ShouldNotEqual, list)
			})

			Convey("Elements should match", func() {
				So(copy.First(), ShouldEqual, 1)
				So(copy.Rest().First(), ShouldEqual, 2)
				So(copy.Rest().Rest().First(), ShouldEqual, 3)
			})

			Convey("Last element should have an empty list as Rest", func() {
				So(copy.Rest().Rest().Rest(), ShouldEqual, EMPTY_LIST)
			})
		})
	})

	Convey("Given an empty list", t, func() {
		list = Empty()
		Convey("When copied", func() {
			copy = list.Copy()
			Convey("Result should be empty list", func() {
				So(copy, ShouldEqual, Empty())
			})
		})
	})
}

func TestInsertingElement(t *testing.T) {
	var list *List
	var ins *List

	Convey("Given a list", t, func() {
		list = FromArray([]Value{1, 3})

		Convey("After inserting an element", func() {
			ins = list.Insert(2, 1)

			Convey("Should not modify given list", func() {
				So(list.First(), ShouldEqual, 1)
				So(list.Rest().First(), ShouldEqual, 3)
			})

			Convey("Inserted elements should be in good order", func() {
				So(ins.First(), ShouldEqual, 1)
				So(ins.Rest().First(), ShouldEqual, 2)
				So(ins.Rest().Rest().First(), ShouldEqual, 3)
			})
		})

		Convey("Inserting an element beyond the length", func() {
			ins = list.Insert(10, 100)
			Convey("Should append the element to the list", func() {
				So(ins.First(), ShouldEqual, 1)
				So(ins.Rest().First(), ShouldEqual, 3)
				So(ins.Rest().Rest().First(), ShouldEqual, 10)
			})
		})
	})

	Convey("Given longer list", t, func() {
		list = FromArray([]Value{1, 2, 3, 5})

		Convey("After inserting an element", func() {
			ins = list.Insert(4, 3)
			Convey("Elements should be in good order", func() {
				So(ins.First(), ShouldEqual, 1)
				So(ins.Rest().First(), ShouldEqual, 2)
				So(ins.Rest().Rest().First(), ShouldEqual, 3)
				So(ins.Rest().Rest().Rest().First(), ShouldEqual, 4)
				So(ins.Rest().Rest().Rest().Rest().First(), ShouldEqual, 5)
			})
		})
	})

	Convey("Given an empty list", t, func() {
		list = Empty()

		Convey("When inserting emelent at index 0", func() {
			ins = list.Insert(1, 0)

			Convey("First element should be set", func() {
				So(ins.First(), ShouldEqual, 1)
			})
			Convey("Tail should be empty list", func() {
				So(ins.Rest(), ShouldEqual, Empty())
			})
		})
	})
}

func TestHeterogenousList(t *testing.T) {
	var list *List

	Convey("Given an empty list", t, func() {
		list = Empty()

		Convey("When adding elements of different types", func() {
			list = list.Cons(1)
			list = list.Cons("Hello")

			Convey("List should not complain", func() {
				So(list.First(), ShouldEqual, "Hello")
				So(list.Rest().First(), ShouldEqual, 1)
			})
		})
	})
}
