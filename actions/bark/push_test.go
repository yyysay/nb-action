package bark

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPush(t *testing.T) {

	// 创建一个假的 Bark Server

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				writer http.ResponseWriter,
				request *http.Request,
			) {

				// 检查请求方法

				if request.Method != http.MethodGet {

					t.Errorf(
						"expected GET, got %s",
						request.Method,
					)

				}

				// 检查请求路径

				expectedPath :=
					"/device123/testciphertext"

				if request.URL.Path != expectedPath {

					t.Errorf(
						"expected path %s, got %s",
						expectedPath,
						request.URL.Path,
					)

				}

				// 返回成功状态

				writer.WriteHeader(
					http.StatusOK,
				)

			},
		),
	)

	// 测试结束关闭服务器

	defer server.Close()

	// 创建测试配置

	config := Config{

		Server: server.URL,

		DeviceKey: "device123",
	}

	// 测试推送

	err := push(
		config,
		"testciphertext",
	)

	if err != nil {

		t.Fatal(err)

	}

}
