package funcs

/*
 * gomonkey 是 golang 的一款打桩框架，可以认为是monkey的升级版。
 * 功能很全面，支持以下场景：
 *    1. 为一个函数打一个桩
 *    2. 为一个成员方法打一个桩
 *    3. 为一个全局变量打一个桩
 *    4. 为一个函数变量打一个桩
 *    5. 为一个函数打一个特定的桩序列（相当于是为一系列调用返回不同的定制化返回值）
 *    6. 为一个成员方法打一个特定的桩序列
 *    7. 为一个函数变量打一个特定的桩序列
 *
 * Note:
 *    需要关闭内联优化，否则打桩失败：
 *    -gcflags=all=-l
 */
import (
	"github.com/agiledragon/gomonkey"
	"github.com/smartystreets/goconvey/convey"
	"reflect"
	"github.com/hq-cml/go-unittest/common"
	"testing"
)

//为一个函数打一个桩
//定义：
//func ApplyFunc(target, double interface{}) *Patches  //桩源头
//func (this *Patches) ApplyFunc(target, double interface{}) *Patches //桩链
//第一个参数是函数名，第二个参数是桩函数。测试完成后，patches对象通过Reset方法删除所有测试桩。
func TestApplyFunc(t *testing.T) {
	type args struct {
		p string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "case1", args: args{p: "what ever"}, want: []string{"xx", "yy"}},
	}

	patch := gomonkey.ApplyFunc(MyFunc, func(p string) []string {
		return []string{"xx", "yy"}
	})
	defer patch.Reset()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MyFunc(tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MyFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

//为一个成员方法打一个桩
//定义：
//func ApplyMethod(target reflect.Type, methodName string, double interface{}) *Patches
//func (this *Patches) ApplyMethod(target reflect.Type, methodName string, double interface{}) *Patches
//第一个参数是目标类的指针变量的反射类型，第二个参数是字符串形式的方法名，第三个参数是桩函数。
func TestApplyMethod(t *testing.T) {
	var p *common.RealClient
	convey.Convey("TestApplyMethod", t, func() {

		convey.Convey("for succ", func() {

			patches := gomonkey.ApplyMethod(reflect.TypeOf(p), "Get", func(_ *common.RealClient, _ string) (string, bool) {
				return "Hello", true
			})
			defer patches.Reset()

			client := common.NewStorageClient()
			v, ok := client.Get("what ever")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(v, convey.ShouldEqual, "Hello")
		})

		convey.Convey("two methods", func() {

			patches := gomonkey.ApplyMethod(reflect.TypeOf(p), "Get", func(_ *common.RealClient, _ string) (string, bool) {
				return "Hello", true
			})
			patches.ApplyMethod(reflect.TypeOf(p), "Set", func(_ *common.RealClient, _, _ string) error {
				return nil
			})
			defer patches.Reset()

			client := common.NewStorageClient()
			v, ok := client.Get("what ever")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(v, convey.ShouldEqual, "Hello")

			err := client.Set("what ever", "fuck")
			convey.So(err, convey.ShouldEqual, nil)
		})

		convey.Convey("one func and one method", func() {

			patches := gomonkey.ApplyFunc(common.Exec, func(_ string, _ ...string) (string, error) {
				return "World", nil
			})
			patches.ApplyMethod(reflect.TypeOf(p), "Get", func(_ *common.RealClient, _ string) (string, bool) {
				return "Hello", true
			})
			defer patches.Reset()

			client := common.NewStorageClient()
			v, ok := client.Get("what ever")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(v, convey.ShouldEqual, "Hello")

			output, err := common.Exec("", "")
			convey.So(err, convey.ShouldEqual, nil)
			convey.So(output, convey.ShouldEqual, "World")
		})
	})
}

//打桩全局变量
//定义：
//func ApplyGlobalVar(target, double interface{}) *Patches
//func (this *Patches) ApplyGlobalVar(target, double interface{}) *Patches
//ApplyGlobalVar 第一个参数是全局变量的地址，第二个参数是全局变量的桩。
var num = 10
func TestApplyGlobalVar(t *testing.T) {
	convey.Convey("TestApplyGlobalVar", t, func() {

		convey.Convey("change", func() {
			patches := gomonkey.ApplyGlobalVar(&num, 150)
			defer patches.Reset()
			convey.So(num, convey.ShouldEqual, 150)
		})

		convey.Convey("recover", func() {
			convey.So(num, convey.ShouldEqual, 10)
		})
	})
}

//为一个函数变量打一个桩
//定义:
//func ApplyFuncVar(target, double interface{}) *Patches
//func (this *Patches) ApplyFuncVar(target, double interface{}) *Patches
//第一个参数是函数变量的地址，第二个参数是桩函数。
var FuncVar = MyFunc
func TestApplyFuncVar(t *testing.T) {
	convey.Convey("TestApplyFuncVar", t, func() {
		convey.Convey("for succ", func() {
			want := []string{"xx", "yy"}
			patches := gomonkey.ApplyFuncVar(&FuncVar, func(p string) []string {
				return want
			})
			defer patches.Reset()

			s := FuncVar("what ever")
			convey.So(len(s), convey.ShouldEqual, len(want))
			convey.So(s[0], convey.ShouldEqual, want[0])
			convey.So(s[1], convey.ShouldEqual, want[1])
		})
	})
}

//为一个函数打一个特定的桩序列（相当于是为一系列调用返回不同的定制化返回值）
/*
 * type Params []interface{}
 * type OutputCell struct {
 *   Values Params
 *   Times  int
 * }
 */
//func ApplyFuncSeq(target interface{}, outputs []OutputCell) *Patches
//func (this *Patches) ApplyFuncSeq(target interface{}, outputs []OutputCell) *Patches
func TestApplyFuncSeq(t *testing.T) {
	convey.Convey("TestApplyFuncSeq", t, func() {
		convey.Convey("default times is 1", func() {
			ret1 := []string{"hello", "cpp"}
			ret2 := []string{"hello", "golang"}
			ret3 := []string{"hello", "gomonkey"}
			outputs := []gomonkey.OutputCell{
				{Values: gomonkey.Params{ret1}},
				{Values: gomonkey.Params{ret2}},
				{Values: gomonkey.Params{ret3}},
			}
			patches := gomonkey.ApplyFuncSeq(MyFunc, outputs)
			defer patches.Reset()
			output := MyFunc("")
			convey.So(output[1], convey.ShouldEqual, ret1[1])
			output = MyFunc("")
			convey.So(output[1], convey.ShouldEqual, ret2[1])
			output = MyFunc("")
			convey.So(output[1], convey.ShouldEqual, ret3[1])
		})

		convey.Convey("retry succ util the third times", func() {
			ret1 := []string{"hello", "cpp"}
			ret2 := []string{"hello", "golang"}
			outputs := []gomonkey.OutputCell{
				{Values: gomonkey.Params{ret1}, Times: 2}, //次数
				{Values: gomonkey.Params{ret2}},
			}
			patches := gomonkey.ApplyFuncSeq(MyFunc, outputs)
			defer patches.Reset()
			output := MyFunc("")
			convey.So(output[1], convey.ShouldEqual, ret1[1])
			output = MyFunc("")
			convey.So(output[1], convey.ShouldEqual, ret1[1])
			output = MyFunc("")
			convey.So(output[1], convey.ShouldEqual, ret2[1])
		})

	})
}

//为一个成员方法打一个特定的桩序列
//定义：
//func ApplyMethodSeq(target reflect.Type, methodName string, outputs []OutputCell) *Patches
//func (this *Patches) ApplyMethodSeq(target reflect.Type, methodName string, outputs []OutputCell) *Patches
func TestApplyMethodSeq(t *testing.T) {
	e := &common.RealClient{}
	convey.Convey("TestApplyMethodSeq", t, func() {

		convey.Convey("default times is 1", func() {
			ret1 := "hello cpp"
			ret2 := "hello golang"
			ret3 := "hello gomonkey"
			outputs := []gomonkey.OutputCell{
				{Values: gomonkey.Params{ret1, true}},
				{Values: gomonkey.Params{ret2, true}},
				{Values: gomonkey.Params{ret3, true}},
			}
			patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(e), "Get", outputs)
			defer patches.Reset()

			output, ok := e.Get("")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(output, convey.ShouldEqual, ret1)
			output, ok = e.Get("")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(output, convey.ShouldEqual, ret2)
			output, ok = e.Get("")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(output, convey.ShouldEqual, ret3)
		})

		convey.Convey("retry succ util the third times", func() {
			ret1 := "hello cpp"
			ret2 := "hello golang"

			outputs := []gomonkey.OutputCell{
				{Values: gomonkey.Params{ret1, true}, Times: 2},
				{Values: gomonkey.Params{ret2, true}},
			}
			patches := gomonkey.ApplyMethodSeq(reflect.TypeOf(e), "Get", outputs)
			defer patches.Reset()

			output, ok := e.Get("")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(output, convey.ShouldEqual, ret1)
			output, ok = e.Get("")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(output, convey.ShouldEqual, ret1)
			output, ok = e.Get("")
			convey.So(ok, convey.ShouldEqual, true)
			convey.So(output, convey.ShouldEqual, ret2)
		})
	})
}

//为一个函数变量打一个特定的桩序列
//定义：
//func ApplyFuncVarSeq(target interface{}, outputs []OutputCell) *Patches
//func (this *Patches) ApplyFuncVarSeq(target interface{}, outputs []OutputCell) *Patches
func TestApplyFuncVarSeq(t *testing.T) {
	convey.Convey("TestApplyFuncVarSeq", t, func() {
		convey.Convey("for succ", func() {
			ret1 := []string{"hello", "cpp"}
			ret2 := []string{"hello", "golang"}
			ret3 := []string{"hello", "gomonkey"}
			outputs := []gomonkey.OutputCell{
				{Values: gomonkey.Params{ret1}},
				{Values: gomonkey.Params{ret2}},
				{Values: gomonkey.Params{ret3}},
			}
			patches := gomonkey.ApplyFuncVarSeq(&FuncVar, outputs)
			defer patches.Reset()

			s := FuncVar("what ever")
			convey.So(len(s), convey.ShouldEqual, len(ret1))
			convey.So(s[1], convey.ShouldEqual, ret1[1])

			s = FuncVar("what ever")
			convey.So(len(s), convey.ShouldEqual, len(ret2))
			convey.So(s[1], convey.ShouldEqual, ret2[1])

			s = FuncVar("what ever")
			convey.So(len(s), convey.ShouldEqual, len(ret3))
			convey.So(s[1], convey.ShouldEqual, ret3[1])
		})
	})
}
