package gomock

/*
 * GoMock是由Golang官方开发维护的测试框架，主要用于interface的Mock功能，实现依赖注入之类
 * 与Golang内置的testing包良好集成，含了GoMock包和mockgen工具两部分，其中
 *    GoMock包完成对桩对象生命周期的管理
 *    mockgen工具用来生成interface对应的Mock类源文件。
 *
 * mockgen有两种操作模式：源文件和反射。
 *    源文件模式通过一个包含interface定义的文件生成mock类文件，它通过 -source 标识生效，举例：
 *       mockgen -source=foo.go [other options]
 *    反射模式通过构建一个程序用反射理解接口生成一个mock类文件，它通过两个非标志参数生效：导入路径和用逗号分隔的符号列表（多个interface）。举例：
 *       mockgen database/sql/driver Conn,Driver
 *       注意：第一个参数是基于GOPATH的相对路径，第二个参数可以为多个interface，并且多个interface之间只能用逗号分隔，不能有空格。
 *
 * mockgen支持的选项如下：
 *   -source: 一个文件包含打算mock的接口列表，源文件模式下使用
 *   -destination: 存放生成mock类代码的文件。如果你没有设置这个选项，代码将被打印到标准输出
 *   -package: 用于指定mock类源文件的包名。如果你没有设置这个选项，则包名由mock_和输入文件的包名级联而成
 *
 * 本例：
 *   mockgen -destination=./mocks/mock.go -package=mocks github.com/hq-cml/go-unittest/common StorageClient,Decoder
 *
 * 测试套路：
 *	 1. mock控制器生成
 *      mock控制器通过NewController接口生成，是mock生态系统的顶层控制。
 *      它定义了mock对象的作用域和生命周期，以及它们的期望。多个协程同时调用控制器的方法是安全的。
 *   2. 创建mock对象，并将mock对象注入控制器
 *      （如果有多个mock对象则注入同一个控制器 )
 *   3. mock对象的行为注入
 *       对于mock对象的行为注入，控制器是通过map+slice来维护的，一个方法对应map的一项。
 *       因为一个方法在一个用例中可能调用多次，所以map的值类型是数组切片。
 *       当mock对象进行行为注入时，控制器会将行为Add。当该方法被调用时，控制器会将该行为Remove。
 *       所以通常有几个case，就会有几个EXPECT()，或者也可以通过AnyTimes设置无限次数
 *       如果需要mock的方法特别复杂，仅指定返回值不能够满足，可以试用DoAndReturn来解决
 *   4. 依赖注入或者打桩
 *       如果实现是作为参数传入，则作为参数传入（依赖注入）
 *       如果是内部生成的实现，则需要打桩替换
 *       （理论上可以用gostub，但是有侵入性，需要用函数变量替换原函数，所以用gomonkey更好[go test参数：-gcflags=all=-l，取消内联]）
 *
 * 行为调用的保序
 *  默认情况下，行为调用顺序按每个方法注册EXPECT方式一致。如果要严格保序，有两种方法：
 *    通过After关键字来实现保序
 *    通过InOrder关键字来实现保序，见例子
 *  Note：当mock对象行为的注入保序后，如果行为调用的顺序和其不一致，就会触发测试失败。
 *
 */
import (
	"errors"
	"github.com/agiledragon/gomonkey"
	"github.com/golang/mock/gomock"
	"github.com/hq-cml/go-unittest/common"
	"github.com/hq-cml/go-unittest/gomock/mocks"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCheckItemKey1(t *testing.T) {
	type args struct {
		client common.StorageClient
		key    string
	}

	//mock控制器生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//创建mock对象，并将mock对象注入控制器（如果有多个mock对象则注入同一个控制器 )
	mockCli := mocks.NewMockStorageClient(mockCtrl)

	//mock对象的行为注入
	//这里是3个case，所以是3个EXPECT()
	mockCli.EXPECT().Get(gomock.Any()).Return("fuck", false)
	mockCli.EXPECT().Get(gomock.Any()).Return("fuck", true)
	mockCli.EXPECT().Get(gomock.Any()).Return("Hello world", true)

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "case1",
			args: args {
				client: mockCli,  //将mock对象从外部注入
				key: "what ever",
			},
			want: false,
			wantErr: true,
		},
		{
			name: "case2",
			args: args {
				client: mockCli,
				key: "what ever",
			},
			want: false,
			wantErr: true,
		},
		{
			name: "case3",
			args: args {
				client: mockCli,
				key: "what ever",
			},
			want: true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckItemKey1(tt.args.client, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckItemKey1() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckItemKey1() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckItemKey2(t *testing.T) {
	type args struct {
		key string
	}

	//mock控制器生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//创建mock对象，并将mock对象注入控制器（如果有多个mock对象则注入同一个控制器 )
	mockCli := mocks.NewMockStorageClient(mockCtrl)

	//mock对象的行为注入
	//这里是3个case，所以是3个EXPECT()
	mockCli.EXPECT().Get(gomock.Any()).Return("fuck", false)
	mockCli.EXPECT().Get(gomock.Any()).Return("fuck", true)
	mockCli.EXPECT().Get(gomock.Any()).Return("Hello world", true)

	//打桩注册
	//gostub需要借助函数变量，有侵入性
	//stubs := gostub.Stub(common.NewStorageClient, func() common.StorageClient {
	//	return mockCli
	//})
	//defer stubs.Reset()

	//gomonkey
	//go test参数：-gcflags=-l，取消内联
	patch := gomonkey.ApplyFunc(common.NewStorageClient, func() common.StorageClient {
		return mockCli
	})
	defer patch.Reset()

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "case1",
			args: args {
				key: "what ever",
			},
			want: false,
			wantErr: true,
		},
		{
			name: "case2",
			args: args {
				key: "what ever",
			},
			want: false,
			wantErr: true,
		},
		{
			name: "case3",
			args: args {
				key: "what ever",
			},
			want: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckItemKey2(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckItemKey2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckItemKey2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//测试保序
func TestReplace(t *testing.T) {
	type args struct {
		client common.StorageClient
		key    string
		def    string
	}

	//mock控制器生成
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	//创建mock对象，并将mock对象注入控制器（如果有多个mock对象则注入同一个控制器 )
	mockCli := mocks.NewMockStorageClient(mockCtrl)

	//mock对象的行为注入
	//mockCli.EXPECT().Get(gomock.Any()).Return("fuck", false)
	//mockCli.EXPECT().Get(gomock.Any()).Return("Hello", true)
	//mockCli.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)

	//InOrder 保序，顺序必须严格匹配，否则出凑
	gomock.InOrder(
		mockCli.EXPECT().Get(gomock.Any()).Return("fuck", false),
		mockCli.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil),
		mockCli.EXPECT().Get(gomock.Any()).Return("Hello", true),
	)

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "case1",
			args: args {
				client: mockCli,  //将mock对象从外部注入
				key: "key1",
				def: "Hello",
			},
			want: "Hello",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Replace(tt.args.client, tt.args.key, tt.args.def)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Replace() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//利用DoAndReturn来mock更加复杂的方法逻辑
func TestDoAndReturn(t *testing.T) {
	convey.Convey("TestDoAndReturn", t, func() {
		convey.Convey("case1", func() {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockDecoder := mocks.NewMockDecoder(ctrl)
			mockDecoder.EXPECT().Unmarshal(gomock.Any(), gomock.Any()).DoAndReturn(
				func(data []byte, movie interface{}) (error) {
					mv, ok := movie.(*common.Movie)
					if !ok {
						return errors.New("Wrong Type")
					}
					mv.Name = "Go"
					mv.Type = "Love"
					mv.Score = 95
					return nil
				}).AnyTimes()

			patch := gomonkey.ApplyFunc(common.NewDecoder, func() common.Decoder {
				return mockDecoder
			})
			defer patch.Reset()

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