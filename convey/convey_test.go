package convey

/*
 * 参考：https://www.jianshu.com/p/e3b2b1194830
 * GoConvey是一个单元测试框架，可以将单测case进行整合，并
 * GoConvey框架兼容Golang原生的单元测试，所以可以使用go test -v来运行测试。
 *
 * 每个测试用例必须使用Convey函数包裹起来:
 *   第一个参数为string类型的测试描述
 *   第二个参数为测试函数的入参（类型为*testing.T）
 *   第三个参数为不接收任何参数也不返回任何值的函数（习惯使用闭包）
 *
 * Convey函数的第三个参数闭包的实现中通过So函数完成断言判断
 *   第一个参数为实际值
 *   第二个参数为断言函数变量
 *   第三个参数或者没有（当第二个参数为类ShouldBeTrue形式的函数变量）或者有（当第二个函数为类ShouldEqual形式的函数变量）
 *
 * 定制化断言函数
 * 	 type assertion func(actual interface{}, expected ...interface{}) string
 *   当assertion的返回值为""时表示断言成功，否则表示失败
 */
import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

//最常规的例子，只有一个case，并且测试通过
func TestStringSliceEqual1(t *testing.T) {
	convey.Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
		a := []string{"hello", "goconvey"}
		b := []string{"hello", "goconvey"}
		convey.So(StringSliceEqual(a, b), convey.ShouldBeTrue)
	})
}

//测试不通过的情况，将上面例子强行改成不通过
//PS：这里用的是SkipSo跳过！
func TestStringSliceEqual2(t *testing.T) {
	convey.Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
		a := []string{"hello", "goconvey"}
		b := []string{"hello", "goconvey"}
		// SkipSo !!
		convey.SkipSo(StringSliceEqual(a, b), convey.ShouldBeFalse)
	})
}

//多个case的情况下，每一个Convey语句对应一个测试用例
//那么一个函数的多个测试用例可以通过一个测试函数的多个Convey语句来呈现。
func TestStringSliceEqual3(t *testing.T) {
	convey.Convey("TestStringSliceEqual should return true when a != nil  && b != nil", t, func() {
		a := []string{"hello", "goconvey"}
		b := []string{"hello", "goconvey"}
		convey.So(StringSliceEqual(a, b), convey.ShouldBeTrue)
	})

	convey.Convey("TestStringSliceEqual should return true when a ＝= nil  && b ＝= nil", t, func() {
		convey.So(StringSliceEqual(nil, nil), convey.ShouldBeTrue)
	})

	convey.Convey("TestStringSliceEqual should return false when a ＝= nil  && b != nil", t, func() {
		a := []string(nil)
		b := []string{}
		convey.So(StringSliceEqual(a, b), convey.ShouldBeFalse)
	})

	convey.Convey("TestStringSliceEqual should return false when a != nil  && b != nil", t, func() {
		a := []string{"hello", "world"}
		b := []string{"hello", "goconvey"}
		convey.So(StringSliceEqual(a, b), convey.ShouldBeFalse)
	})
}

//更地道的用法：Convey语句无限嵌套（无限嵌套可以体现测试用例之间的关系）
//注意：无限嵌套中，只有最外层的Convey需要传入*testing.T类型的变量t。
func TestStringSliceEqual4(t *testing.T) {
	convey.Convey("TestStringSliceEqual", t, func() {
		convey.Convey("should return true when a != nil  && b != nil", func() {
			a := []string{"hello", "goconvey"}
			b := []string{"hello", "goconvey"}
			convey.So(StringSliceEqual(a, b), convey.ShouldBeTrue)
		})

		convey.Convey("should return true when a ＝= nil  && b ＝= nil", func() {
			convey.So(StringSliceEqual(nil, nil), convey.ShouldBeTrue)
		})

		convey.Convey("should return false when a ＝= nil  && b != nil", func() {
			a := []string(nil)
			b := []string{}
			convey.So(StringSliceEqual(a, b), convey.ShouldBeFalse)
		})

		convey.Convey("should return false when a != nil  && b != nil", func() {
			a := []string{"hello", "world"}
			b := []string{"hello", "goconvey"}
			convey.So(StringSliceEqual(a, b), convey.ShouldBeFalse)
		})
	})
}

//定制化断言函数
//相当于自己定制的deepEqual
func ShouldSummerBeComming(actual interface{}, expected ...interface{}) string {
	if actual == "summer" && expected[0] == "comming" {
		return ""
	} else {
		return "summer is not comming!"
	}
}

func TestSummer(t *testing.T) {
	convey.Convey("TestSummer", t, func() {
		convey.So("summer", ShouldSummerBeComming, "comming") //通过
		convey.SkipSo("winter", ShouldSummerBeComming, "comming")       //不通过，用SkipSo跳过
	})
}