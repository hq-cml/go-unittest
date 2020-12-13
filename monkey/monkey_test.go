package monkey
/*
 * Monkey是Golang的一个猴子补丁（monkeypatching）框架
 * 原理是在运行时通过汇编语句重写可执行文件，将待打桩函数或方法的实现跳转到桩实现，原理和热补丁类似。
 * Monkey框架的使用场景很多，依次为：
 *    场景1：为一个函数打桩
 *    场景2：为一个过程打桩
 *    场景3：为一个方法打桩
 *    场景4：桩中桩
 *
 * Monkey的问题：
 *    不支持多次调用桩函数（方法）而呈现不同行为的复杂情况。
 *
 * Note:
 *    Monkey其实完全可以用gomonkey进行替代，所以意义不大了
 *
 * Ps:
 *    mockgen -destination=./mocks/mock.go -package=mocks github.com/hq-cml/go-unittest/common Decoder
 */
import (
	mk "bou.ke/monkey"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/smartystreets/goconvey/convey"
	"reflect"
	"github.com/hq-cml/go-unittest/common"
	"github.com/hq-cml/go-unittest/monkey/mocks"
	"testing"
)

//为函数打桩
func TestMonkeyExec1(t *testing.T) {
	patch := mk.Patch(common.Exec, func(_ string, _ ...string) (string, error) {
		return "xxx-vethName100-yyy", nil
	})
	defer patch.Unpatch()

	ret, err := common.Exec("what ever", "arg1", "arg2")
	if err != nil {
		t.Errorf("Should not error: %v", err.Error())
	}
	if ret != "xxx-vethName100-yyy" {
		t.Errorf("The ret should = xxx-vethName100-yyy")
	}
}

//为方法打桩
func TestMonkeyGet(t *testing.T) {
	var client *common.RealClient
	patch := mk.PatchInstanceMethod(reflect.TypeOf(client), "Get",
		func( _ *common.RealClient, _ string) (string,bool) {
			return "Hello", true
	})
	defer patch.Unpatch()

	ret, ok := client.Get("what ever")
	if ok != true {
		t.Errorf("Should return true")
	}
	if ret != "Hello" {
		t.Errorf("The ret should = Hello")
	}
}

//桩中桩
//感觉这个桩中装的写法，没有太大必要，可以直接用gomock的DoAndReturn来代替
func TestStubInStub(t *testing.T) {
	convey.Convey("TestStubInStub", t, func() {
		convey.Convey("case1", func() {
			//mock控制器
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			//mock对象注入控制器
			mockDecoder := mocks.NewMockDecoder(ctrl)
			//mockDecoder.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).Return(nil)

			//利用monkey打桩NewDecoder函数
			mk.Patch(common.NewDecoder, func() common.Decoder {
				return mockDecoder
			})
			defer mk.UnpatchAll()
			//利用monkey对mock对象进行二次打桩
			mk.PatchInstanceMethod(reflect.TypeOf(mockDecoder), "Unmarshal", func(_ *mocks.MockDecoder, data []byte, movie interface{}) error {
				mv, ok := movie.(*common.Movie)
				if !ok {
					return errors.New("Wrong Type")
				}
				mv.Name = "Go"
				mv.Type = "Love"
				mv.Score = 95
				return nil
			})
			decoder := common.NewDecoder()
			var movie = &common.Movie{}
			err := decoder.Unmarshal([]byte("Titanic"), movie)
			convey.So(err, convey.ShouldBeNil)
			convey.So(movie.Name, convey.ShouldEqual, "Go")
			convey.So(movie.Type, convey.ShouldEqual, "Love")
			convey.So(movie.Score, convey.ShouldEqual, 95)
		})
	})
}

