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

func TestCollidingHashes(t *testing.T) {
	Convey("Given an empty hash map with colliding hash func", t, func() {
		hm := New()
		hm.hashFunc = func(val Value) (uint32, error) {
			if val == "hello" {
				return 0, nil
			} else {
				return 1024, nil
			}
		}

		Convey("When colliding hashes emerge", func() {
			filled := hm.Assoc("hello", 42)
			filled = filled.Assoc("world", 24)
			Convey("Both values should be stored", func() {
				So(filled.Find("world"), ShouldEqual, 24)
				So(filled.Find("hello"), ShouldEqual, 42)
			})
		})

		Convey("When colliding hash emerge for find", func() {
			filled := hm.Assoc("hello", 42)
			Convey("Should not find the value", func() {
				So(filled.Find("other"), ShouldEqual, nil)
			})
		})

	})
}
