package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"lol-teammate-helper/internal/reqLib"
	"net/http"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
var port int = 56737
var token string = "JaffTt5fcKeFztB5KHO4Hw"
var baseDataUrl string = "/lol-summoner/v1/current-summoner"

func (a *App) Greet(name string) string {
	fmt.Println("app.go Greet")
	url := fmt.Sprintf("https://127.0.0.1:%d%s", port, baseDataUrl)

	// 创建HTTP客户端，跳过TLS证书验证（因为这是本地自签名证书）
	//client := &http.Client{
	//	Timeout: 10 * time.Second,
	//	Transport: &http.Transport{
	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	},
	//}
	client := reqLib.Instance().Client()

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Sprintf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 创建Authorization头：Basic + base64("riot:token")
	authString := "riot:" + token
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Set("Authorization", "Basic "+encodedAuth)
	fmt.Println("Auth String:", encodedAuth)
	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("读取响应失败: %v", err)
	}

	// 返回响应结果
	result := fmt.Sprintf("状态码: %d\n响应体: %s", resp.StatusCode, string(body))
	fmt.Println(result)

	return string(body)
}

func (a *App) GetImgSrc() string {
	url := "https://127.0.0.1:56737/lol-game-data/assets/v1/profile-icons/4804.jpg"

	// 创建HTTP客户端，跳过TLS证书验证（因为这是本地自签名证书）
	//client := &http.Client{
	//	Timeout: 10 * time.Second,
	//	Transport: &http.Transport{
	//		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//	},
	//}
	client := reqLib.Instance().Client()

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Sprintf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 创建Authorization头：Basic + base64("riot:token")
	authString := "riot:" + token
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP状态码:cite[1]
	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("请求失败，状态码: %d", resp.StatusCode)
	}

	// 读取图片数据:cite[1]:cite[3]
	imgData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("读取图片数据失败: %v", err)
	}

	// 将图片数据转换为Base64字符串
	base64Str := base64.StdEncoding.EncodeToString(imgData)
	//fmt.Println("头像结果" + base64Str)
	// 返回Base64编码的图片，可直接用于<img>标签的src属性

	res := "data:image/jpeg;base64," + base64Str
	//fmt.Println("头像地址：" + res)
	return res
}
