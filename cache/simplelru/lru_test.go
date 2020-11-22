package simplelru

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestLRU(t *testing.T) {
	Convey("Test LRU", t, func() {
		evictCount := 0
		l,_ := NewLRU(128, func(key interface{}, value interface{}) {
			evictCount ++
		})
		for i := 1; i <= 300; i++ {
			l.Add(i, i)
		}

		So(evictCount, ShouldEqual, 300-128)

		_,_ = Println("Len should be 128")
		So(l.Len(), ShouldEqual, 128)

		expectKeys := []interface{}{}
		for i:=300-128+1;;i++ {
			expectKeys = append(expectKeys, i)
			if len(expectKeys) >= 128{
				break
			}
		}

		keys := l.Keys()
		So(keys, ShouldResemble, expectKeys)


		v, ok := l.Get(128)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeNil)

		v, ok = l.Get(299)
		So(ok, ShouldBeTrue)
		So(v.(int), ShouldEqual, 299)

		//"The Oldest One Should Be 173, After get should be 174"
		_, v, ok = l.GetOldest()
		So(v.(int), ShouldEqual, 173)

		v, ok = l.Peek(173)
		So(v.(int), ShouldEqual, 173)
		_, v, ok = l.GetOldest()
		So(v.(int), ShouldEqual, 173)

		v, ok = l.Get(173)
		_, v, ok = l.GetOldest()
		So(v.(int), ShouldEqual, 174)


		ok = l.Contains(299)
		So(ok, ShouldBeTrue)

		ok = l.Remove(299)
		So(ok, ShouldBeTrue)
		ok = l.Remove(110)
		So(ok, ShouldBeFalse)

		ok = l.Contains(299)
		So(ok, ShouldBeFalse)




	})
}


func TestLRU_Add(t *testing.T) {
	Convey("Test LRU add", t, func() {
		l, _ := NewLRU(2, nil)
		for i:=0;i<3 ;i++  {
			l.Add(i, i)
		}
		v, ok := l.Get(0)
		So(ok, ShouldBeFalse)
		So(v, ShouldBeNil)

		v, ok = l.Get(1)
		So(ok, ShouldBeTrue)
		So(v.(int), ShouldEqual, 1)
		So(l.Len(), ShouldEqual, 2)
	})
}

func TestLRU_Purge(t *testing.T) {
	Convey("Test purge", t, func() {
		evictCount := 0
		l, _ := NewLRU(128, func(key interface{}, value interface{}) {
			evictCount++
		})
		for i := 1; i <= 300; i++ {
			l.Add(i, i)
		}

		Convey("Test purge", func() {
			l.Purge()
			So(l.Len(), ShouldEqual, 0)
		})

		Convey("Test resize", func() {
			for i := 1; i <= 300; i++ {
				l.Add(i, i)
			}
			So(l.Len(), ShouldEqual, 128)

			d := l.Resize(200)
			So(d, ShouldEqual, 128-200)

			for i := 1; i <= 300; i++ {
				l.Add(i, i)
			}
			So(l.Len(), ShouldEqual, 128)

		})
	})

}


