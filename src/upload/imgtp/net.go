package imgtp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"typora-image-upload/src/utils"
)

// post
const url_get_token string = "https://www.imgtp.com/api/token"

// post 注意：请求时header如果有参数 token，接口则认证该token，上传的图片也是在该token用户下，否则为游客上传。
const url_image_upload string = "https://www.imgtp.com/api/upload"

// read user_info.json
const filename_user_info = "user_info.json"

type getTokenRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Refresh  int    `json:"refresh,int"` // default 0, not Refresh
}

var request_params getTokenRequestParams

//type getTokenResponse struct {
//	code int    `json:"code"`
//	msg  string `json:"msg"`
//	time int64  `json:"time"`
//	data string `json:"data"`
//}

func init() {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		os.Exit(-1)
	}

	// read user_info.json
	sp := strings.LastIndex(path, "\\")
	path_abs := string([]rune(path)[0 : sp+1])
	content := utils.ReadFile(path_abs + filename_user_info)

	err = json.Unmarshal(content, &request_params)
	if err != nil {
		os.Exit(-1)
	}
}

func GetToken() string {
	//obuf := make([]byte, 512)
	args := fasthttp.AcquireArgs()
	args.Set("email", request_params.Email)
	args.Set("password", request_params.Password)
	args.Set("refresh", strconv.Itoa(request_params.Refresh))

	statusCode, body, err := fasthttp.Post(nil, url_get_token, args)

	if err != nil || statusCode != 200 {
		os.Exit(-1)
	}

	//jm := make(map[string]interface{})
	//json.Unmarshal(body, &jm)
	//return (jm["data"].(map[string]interface{}))["token"].(string)

	s := string(body)
	start := strings.Index(s, "token")
	end := strings.Index(s, "\"},\"time\"")
	if start+8 >= end-1 {
		os.Exit(-1)
	}

	return string(body[start+8 : end])
	// ====
	//url := strings.Builder{}
	//url.WriteString(url_get_token)
	//url.WriteString("?email=")
	//url.WriteString(request_params.Email)
	//url.WriteString("&password=")
	//url.WriteString(request_params.Password)
	//url.WriteString("&refresh=")
	//url.WriteString(strconv.Itoa(request_params.Refresh))
	//
	//request, err := http.NewRequest("POST", url.String(), nil)
	//client := &http.Client{}
	//resp, err := client.Do(request)
	//
	//if err != nil || resp.StatusCode != 200 {
	//	os.Exit(-1)
	//}
	//
	//
	//
	//var dst []byte
	//reader := resp.Body
	//defer reader.Close()
	//
	//for i, _ := reader.Read(obuf); i > 0; i, _ = reader.Read(obuf) {
	//	dst = append(dst, obuf...)
	//}
	//m := make(map[string]interface{}, 8)
	//fmt.Println(m["token"])
}

func UploadImages(filePath []string) {
	// check exist

	for i, _ := range filePath {
		if !utils.IsFileExist(filePath[i]) {
			os.Exit(-1)
		}

		file, err := os.Open(filePath[i])
		if err != nil {
			os.Exit(-1)
		}

		bufReader := bufio.NewReader(file)
		bufReader.Read([]byte{})
	}
}

const image_path = "D:\\Image\\Pictures\\imagename.jpg"

func Try() {
	token := GetToken()

	bdbuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bdbuf)

	bodyWriter.WriteField("token", token)

	// 上传文件
	fileWriter, _ := bodyWriter.CreateFormFile("image", image_path)

	f, _ := os.Open(image_path)
	defer f.Close()

	_, _ = io.Copy(fileWriter, f)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	request := fasthttp.AcquireRequest()
	request.Header.SetMethod("POST")
	request.Header.SetContentType(contentType)
	request.SetRequestURI(url_image_upload)
	request.SetBodyRaw(bdbuf.Bytes())

	response := fasthttp.AcquireResponse()

	fasthttp.Do(request, response)

	var resp_bd map[string]interface{}
	json.Unmarshal(response.Body(), &resp_bd)

	fmt.Println(resp_bd["data"].(map[string]interface{})["url"])
}

func Try2() {
	token := GetToken()

	bdbuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bdbuf)

	bodyWriter.WriteField("token", token)

	// 上传文件
	fileWriter, _ := bodyWriter.CreateFormFile("image", image_path)

	f, _ := os.Open(image_path)
	defer f.Close()

	_, _ = io.Copy(fileWriter, f)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	//req := fasthttp.AcquireRequest()
	//req.Header.Set("token", token)
	//req.Header.Add(fasthttp.HeaderContentType, contentType)
	//req.SetBody(bdbuf.Bytes())
	//req.SetRequestURI("https://www.imgtp.com/api/upload")
	//
	//resp := fasthttp.AcquireResponse()
	//
	//fasthttp.Do(req, resp)
	//
	//fmt.Println(string(resp.Body()))

	// {"code":200,"msg":"success",
	//"data":{"id":"29986","name":"gamersky_02origin_03_201833184750D.jpg",
	//"url":"https:\/\/img1.imgtp.com\/2023\/01\/10\/RxOhrxCg.jpg",
	//"size":835827,"mime":"image\/jpeg",
	//"sha1":"2b13e2a1c09f7320671c12f310b87669f9948760",
	//"md5":"f43242a9a25cccf6ac0d361a4b7125bd",
	//"quota":"32212254720.00","use_quota":"10411198.00"},
	//"time":1673287586}
	resp, _ := http.Post(url_image_upload, contentType, bdbuf)
	defer resp.Body.Close()
	all, _ := io.ReadAll(resp.Body)
	fmt.Println(string(all))

}
