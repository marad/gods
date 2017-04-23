package hashmap

import (
	. "github.com/smartystreets/goconvey/convey"
	. "gods"
	"testing"
)

func TestPopCountFunc(t *testing.T) {
	var key uint32
	Convey("Given a key", t, func() {
		key = 314456588
		Convey("When counting population", func() {
			So(popcnt(key), ShouldEqual, 14)
		})
	})
}

func TestGettingHashParts(t *testing.T) {
	var key uint32
	var part uint32
	Convey("Given a hash key", t, func() {
		key = 314456588
		Convey("Expect to have proper key parts extracted", func() {
			Convey("First key part", func() {
				part = hashPart(key, 0)
				So(part, ShouldEqual, 12)
			})
			Convey("Second key part", func() {
				part = hashPart(key, 1)
				So(part, ShouldEqual, 16)
			})
			Convey("Third key part", func() {
				part = hashPart(key, 2)
				So(part, ShouldEqual, 14)
			})
		})
	})
}

func TestHashing(t *testing.T) {
	Convey("Hashing string works as expected", t, func() {
		hashed, _ := Hash("hello")
		So(hashed, ShouldEqual, 314456588)
	})

	Convey("Trying to hash unknown type returns an error", t, func() {
		_, err := Hash(42)
		So(err, ShouldNotBeNil)
	})
}

func TestValueNodeImplementation(t *testing.T) {
	var vn ValueNode
	var copied ValueNode
	Convey("Given a ValueNode", t, func() {
		vn = ValueNode{123, "Hello"}
		Convey("When assoc-ing new value to the node", func() {
			copied = vn.Assoc(123, "World", 0).(ValueNode)

			Convey("Copied node should have the new value", func() {
				So(copied.Key, ShouldEqual, 123)
				So(copied.BaseValue, ShouldEqual, "World")
			})

			Convey("Original node should not have changed", func() {
				So(vn.Key, ShouldEqual, 123)
				So(vn.BaseValue, ShouldEqual, "Hello")
			})
		})
	})
}

func TestBasicHashMapFunctionality(t *testing.T) {
	var emptyMap *HashMap
	Convey("Given an empty map", t, func() {
		emptyMap = New()
		Convey("When adding new values", func() {
			modified := emptyMap.Assoc("myKey", 42)
			modified2 := modified.Assoc("anotherKey", 24)

			Convey("The keys should be findable in the map", func() {
				So(modified.Find("myKey"), ShouldEqual, 42)
				So(modified2.Find("myKey"), ShouldEqual, 42)
				So(modified2.Find("anotherKey"), ShouldEqual, 24)
			})

			Convey("But the value of original maps should not be modified", func() {
				So(emptyMap.Find("myKey"), ShouldEqual, nil)
				So(emptyMap.Find("anotherKey"), ShouldEqual, nil)
				So(modified.Find("anotherKey"), ShouldEqual, nil)
			})
		})
	})
}

func TestSearchingInEmptyMap(t *testing.T) {
	var emptyMap *HashMap
	var value Value
	Convey("Given an empty map", t, func() {
		emptyMap = New()
		Convey("When searching for nonexistent key in the map", func() {
			value = emptyMap.Find("not-existing")
			Convey("Returned value should be nil", func() {
				So(value, ShouldEqual, nil)
			})
		})
	})
}

func TestOverridingAValue(t *testing.T) {
	var first *HashMap
	var second *HashMap

	Convey("Given a map with a value", t, func() {
		first = New().Assoc("key", 42)
		Convey("When value is overriden", func() {
			second = first.Assoc("key", 24)
			Convey("The new value should be returned", func() {
				So(second.Find("key"), ShouldEqual, 24)
			})
			Convey("But the original map should remain unchanged", func() {
				So(first.Find("key"), ShouldEqual, 42)
			})
		})
	})

}

//func TestCollidingHashes(t *testing.T) {
//	hash := func(val Value) (uint32, error) {
//	}
//}
