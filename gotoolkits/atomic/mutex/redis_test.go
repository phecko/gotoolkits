package mutex

import (
	"github.com/gomodule/redigo/redis"
	. "github.com/smartystreets/goconvey/convey"
	"sync"
	"testing"
	"time"
)

func TestNewRedisMutex_Lock_Unlock(t *testing.T) {

	Convey("test Lock and unlock", t, func() {
		conn, err := redis.Dial("tcp", "127.0.0.1:6379")
		So(err, ShouldBeNil)
		locker := NewRedisMutex(conn, "TestKey1")

		So(locker.Key(), ShouldEqual, "TestKey1")

		isLock := locker.Lock(10)
		So(isLock, ShouldBeTrue)
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func() {
			Convey("Wait 5 seconds To Unlock", t, func() {
				time.Sleep(time.Second * 5)
				isUnlock, err := locker.Unlock()
				So(err, ShouldBeNil)
				So(isUnlock, ShouldBeTrue)
				wait.Done()
			})
		}()

		Convey("Test Lock again", func() {
			isLock := locker.Lock(10)
			So(isLock, ShouldBeFalse)
		})

		wait.Wait()
	})
}

func TestNewRedisMutex_Expire(t *testing.T) {

	Convey("test Lock and expire", t, func() {
		conn, err := redis.Dial("tcp", "127.0.0.1:6379")
		So(err, ShouldBeNil)
		locker := NewRedisMutex(conn, "TestKey2")
		isLock := locker.Lock(10)
		So(isLock, ShouldBeTrue)
		isLockAgain := locker.Lock(10)
		So(isLockAgain, ShouldBeFalse)

		isLockAgain2 := NewRedisMutex(conn, "TestKey2").Lock(10)
		So(isLockAgain2, ShouldBeFalse)

		Convey("Wait 12 seconds To Wait Lock Expire", func() {
			time.Sleep(time.Second * 12)
			isLock := locker.Lock(10)
			So(isLock, ShouldBeTrue)
		})
	})
}

func TestNewRedisMutex_Replace(t *testing.T) {

	Convey("test Lock be replaced", t, func() {
		conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
		locker := NewRedisMutex(conn, "TestKey3")
		isLock := locker.Lock(3)
		print(isLock)
		So(isLock, ShouldBeTrue)
		newLocker := NewRedisMutex(conn, "TestKey3")

		Convey("Wait 5 seconds Locker expire release", func() {
			time.Sleep(5*time.Second)
			isLock2 := newLocker.Lock(5)
			print(isLock2)
			So(isLock2, ShouldBeTrue)

			Convey("First Locker Unlock", func() {
				ok, err := locker.Unlock()
				So(err, ShouldEqual, ErrLockerBeReplace)
				So(ok, ShouldBeFalse)

				ok, err = newLocker.Unlock()
				So(err, ShouldBeNil)
				So(ok, ShouldBeTrue)
			})
		})
	})
}

func TestRedisMutex_With2(t *testing.T) {

	Convey("Run With Lock Normal", t, func() {

		conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
		locker := NewRedisMutex(conn, "TestKey8")
		getLock := locker.With(50, func() {
			return
		})
		So(getLock, ShouldBeTrue)

		getLock = locker.With(5, func() {
			return
		})
		So(getLock, ShouldBeTrue)
	})
}


func TestRedisMutex_With(t *testing.T) {
	Convey("Run With Lock ", t, func() {
		wait := sync.WaitGroup{}
		wait.Add(1)
		go func(wg *sync.WaitGroup) {
			Convey( "Test 1 can Run", t, func(){
				conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
				locker := NewRedisMutex(conn, "TestKey4")
				getLock := locker.With(5, func() {
					time.Sleep(5 * time.Second)
					return
				})
				wg.Done()
				So(getLock, ShouldBeTrue)
			})
		}(&wait)

		time.Sleep(time.Second * 1)
		wait.Add(1)
		go func(wg *sync.WaitGroup) {
			Convey( "Test 2 Fail ", t, func(){
				conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
				locker := NewRedisMutex(conn, "TestKey4")
				getLock := locker.With(5, func() {
					time.Sleep(5 * time.Second)
				})
				wg.Done()
				So(getLock, ShouldBeFalse)
			})
		}(&wait)

		wait.Add(1)
		go func(wg *sync.WaitGroup) {
			Convey( "Test 3 Fail ", t, func(){
				conn, _ := redis.Dial("tcp", "127.0.0.1:6379")
				locker := NewRedisMutex(conn, "TestKey4")
				getLock := locker.With(5, func() {
					time.Sleep(5 * time.Second)
				})
				wg.Done()
				So(getLock, ShouldBeFalse)
			})
		}(&wait)

		wait.Wait()
	})
}
