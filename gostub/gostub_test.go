package gostub

/*
 * GoStub框架的使用场景主要是为了打桩：
 *     场景1：为一个全局变量打桩
 *     场景2：为一个函数打桩
 *     场景3：为一个过程(无返回值的函数)打桩
 *     场景4：由任意相同或不同的基本场景组合而成
 *
 * 注意：
 *    不论是调用Stub函数还是StubFunc函数，都会生成一个stubs对象
 *    	 该对象有Reset操作，即将全局变量的值恢复为原值。
 *    	 该对象仍然有Stub方法和StubFunc方法，所以在一个测试用例中可以同时对多个全局变量、函数或过程打桩。
 *       (这些全局变量、函数或过程会将初始值存在一个map中，并在延迟语句中通过Reset方法统一做回滚处理)
 *
 * 一些问题：
 *    对于函数的stub，必须是函数变量，所以存在一定的侵入性
 *    没法对方法（method，即结构体的方法）进行打桩
 *
 * Note：
 *    GoStub的功能完全可以被gomonkey取代
 */
import (
	"testing"
	"github.com/smartystreets/goconvey/convey"
	"github.com/prashantv/gostub"
)

//为全局变量打桩
func TestCheckGlobalNum(t *testing.T) {
	stubs := gostub.Stub(&num, 150)
	defer stubs.Reset()

	if num != 150 {
		t.Errorf("The num should = 150")
	}
}

//为函数打桩
//利用Stub
func TestStubExec1(t *testing.T) {
	stubs := gostub.Stub(&MyExec, func(cmd string, args ...string) (string, error) {
		return "xxx-vethName100-yyy", nil
	})
	defer stubs.Reset()

	ret, err := MyExec("what ever", "arg1", "arg2")
	if err != nil {
		t.Errorf("Should not error: %v", err.Error())
	}
	if ret != "xxx-vethName100-yyy" {
		t.Errorf("The ret should = xxx-vethName100-yyy")
	}
}

//为函数打桩
//直接利用StubFunc
func TestStubExec2(t *testing.T) {
	stubs := gostub.StubFunc(&MyExec, "xxx-vethName100-yyy", nil)
	defer stubs.Reset()

	ret, err := MyExec("what ever", "arg1", "arg2")
	if err != nil {
		t.Errorf("Should not error: %v", err.Error())
	}
	if ret != "xxx-vethName100-yyy" {
		t.Errorf("The ret should = xxx-vethName100-yyy")
	}
}

//结合Convey框架使用
func TestFuncDemo(t *testing.T) {
	convey.Convey("TestFuncDemo", t, func() {
		convey.Convey("Global variety", func() {
			stubs := gostub.Stub(&num, 50)
			defer stubs.Reset()

			convey.So(num == 50, convey.ShouldBeTrue)
		})

		convey.Convey("For Stub, StubFunc chains", func() {
			stubs := gostub.Stub(&num, 150)
			defer stubs.Reset()

			//继续挂链
			stubs.StubFunc(&MyExec,"xxx-vethName100-yyy", nil)

			convey.So(num == 150, convey.ShouldBeTrue)
			convey.So(func() bool{
				ret, err := MyExec("what ever", "arg1")
				if err != nil {
					return false
				}
				if ret != "xxx-vethName100-yyy" {
					return false
				}
				return true
			}(), convey.ShouldBeTrue)
		})
	})
}