package imgtp

import (
	"bytes"
	"encoding/json"
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

const user_agent_constant = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36 Edg/107.0.1418.42"

// read user_info.json
const filename_user_info = "user_info.json"

// current path
var current_abs_path string
var file_user_info string
var request_params *userInfoJson

type userInfoJson struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Refresh  int    `json:"refresh,int"` // default 0, not Refresh
	Token    string `json:"token,omitempty"`
}

func init() {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		os.Exit(-1)
	}

	// read user_info.json
	sp := strings.LastIndex(path, "\\")
	path_abs := string([]rune(path)[0 : sp+1])

	current_abs_path = path_abs // end with \
	file_user_info = current_abs_path + filename_user_info

	content := utils.ReadFile(file_user_info)

	err = json.Unmarshal(content, &request_params)
	if err != nil {
		os.Exit(-1)
	}

	// if user_info.json token is empty
	if strings.Compare("", request_params.Token) == 0 {
		token := GetToken()
		request_params.Token = token

		// save
		j, _ := json.Marshal(request_params)
		utils.WriteFile(file_user_info, j)
	}

}

func GetToken() string {

	url := url_get_token + "?email=" + request_params.Email +
		"&password=" + request_params.Password +
		"&refresh=" + strconv.Itoa(request_params.Refresh)
	request, _ := http.NewRequest("POST", url, nil)
	request.Header.Set("User-Agent", user_agent_constant)

	response, err := http.DefaultClient.Do(request)
	if err != nil || response.StatusCode != 200 {
		os.Exit(-response.StatusCode)
	}

	data := utils.ReadResponseBody(response)

	s := string(data)
	start := strings.Index(s, "token")
	end := strings.Index(s, "\"},\"time\"")
	if start+8 >= end-1 {
		os.Exit(-1)
	}

	return string(data[start+8 : end]) // token value
}

func UploadImages(filePath []string) []string {
	// check token
	if strings.Compare("", request_params.Token) == 0 {
		os.Exit(-1)
	}

	success_list := []string{}

	for i, _ := range filePath {
		if !utils.IsFileExist(filePath[i]) {
			os.Exit(-1)
		}
		path := filePath[i]
		url := doUpload(path)
		if strings.Compare("", url) == 0 {
			// failed, return local image path
			success_list = append(success_list, path)
		} else {
			success_list = append(success_list, url)
		}
	}

	return success_list
}

func doUpload(filePath string) string {
	bdbuf := new(bytes.Buffer)
	bodyWriter := multipart.NewWriter(bdbuf)

	bodyWriter.WriteField("token", request_params.Token)

	// 上传文件
	fileWriter, _ := bodyWriter.CreateFormFile("image", filePath)

	f, _ := os.Open(filePath)
	defer f.Close()

	io.Copy(fileWriter, f)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, _ := http.Post(url_image_upload, contentType, bdbuf)
	body := utils.ReadResponseBody(resp)

	if body == nil {
		return ""
	}

	// 解析保存地址
	var find map[string]interface{}
	json.Unmarshal(body, &find)
	//{"code":200,"msg":"success",
	//"data":{"id":"30367","name":"51226312_p0.jpg",
	//"url":"https:\/\/img1.imgtp.com\/2023\/01\/10\/dwWlrSa8.jpg",
	//"size":509716,
	//"mime":"image\/jpeg",
	//"sha1":"b99ca0667495bd342b2cc954f4d36db6ca535b34",
	//"md5":"b6175dadf78dec9d8f2e7acc98a0cb0b",
	//"quota":"32212254720.00",
	//"use_quota":"10085087.00"},
	//"time":1673348172}
	if find["code"].(int) != 200 {
		return ""
	}

	orignal_url := (find["data"].(map[string]interface{}))["url"].(string)

	return strings.Replace(orignal_url, "\\", "", -1)
}
