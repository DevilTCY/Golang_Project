package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

/*请求url*/
var player_host = "127.0.0.1:9999"   /*用户管理平台*/
var lottery_host = "127.0.0.1:12220" /*游戏引擎*/

var url_map = map[string]string{
	"player_login":   "http://" + player_host + "/user/un/login",
	"get_game_issue": "http://" + lottery_host + "/engine/T300111",
	"wager":          "http://" + lottery_host + "/engine/T301200",
}

/*用户登陆信息*/
type player struct {
	playercode string /*用户编号*/
	playerID   string /*用户登陆ID*/
	passwd     string /*用户密码*/
	token      string
}

var player_set = []player{
	{"100001222", "13800000002", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001226", "13800000006", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001240", "13800000020", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001262", "13800000042", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001266", "13800000046", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001280", "13800000060", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001284", "13800000064", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001288", "13800000068", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001302", "13800000082", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001306", "13800000086", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001354", "13800000134", "14e1b600b1fd579f47433b88e8d85291", ""},
	{"100001256", "13800000036", "14e1b600b1fd579f47433b88e8d85291", ""},
}

/*投注总票数*/
var (
	total_ticket int
	mutex        sync.Mutex
)

/*票流水号*/
var tickSN int64

//用户登陆，获取token
/*用户登陆返回数据结构
type tmp struct{
	code int
	data string		//token
	error bool
	msg string
	success bool
}
*/
func player_login(count int) {
	url := url_map["player_login"]
	loginData := map[string]interface{}{
		"loginType":         0,
		"deviceCode":        "d520c7a8-421b-4563-b955-f5abc56b97ec",
		"deviceSoftVersion": "V1.0.0",
		"platform":          4,
		"platformModel":     "1.0",
		"playerAccount":     "",
		"password":          "",
		"timestamp":         "",
	}
	for k, val := range player_set {
		//达到请求数量，退出
		if k == count {
			break
		}

		//初始化客户端
		client := &http.Client{}

		//填充请求参数
		loginData["playerAccount"] = val.playerID
		loginData["password"] = val.passwd
		loginData["timestamp"] = time.Now().Unix()

		//请求map转换为json
		data, _ := json.Marshal(loginData)
		//fmt.Println(string(data))

		//发起http请求
		requestData := strings.NewReader(string(data))
		req, _ := http.NewRequest("POST", url, requestData)
		req.Header.Add("Content-Type", "application/json")
		response, err := client.Do(req)
		if err != nil {
			fmt.Println("PlayerCode: ", val.playerID, "login error, ", err)
			return
		}
		body, _ := ioutil.ReadAll(response.Body)
		//fmt.Println(string(body))

		//取出结果中的token，填充用户结构体
		result := make(map[string]interface{})
		json.Unmarshal(body, &result)
		//fmt.Println(result)
		player_set[k].token = fmt.Sprintf("%v", result["data"])

		//关闭连接
		response.Body.Close()
	}
}

//获取游戏信息
func get_game_issue(game_id int, wager_issue *int) {
	url := url_map["get_game_issue"]

	//填充请求参数
	request := map[string]interface{}{
		"GameID":     game_id,
		"Issue":      -1,
		"QueryCount": 1,
		"RecBegin":   0,
	}
	//fmt.Println(request)

	//初始化客户端
	client := &http.Client{}

	//请求map转换为json
	data, _ := json.Marshal(request)
	//fmt.Println(string(data))

	//发起http请求
	requestData := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", url, requestData)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("PlayerCode: ", game_id, "login error, ", err)
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))

	//取出结果中的期号
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	//判断返回码
	retcode, _ := strconv.Atoi(result["RetCode"].(string))
	if retcode != 0 {
		fmt.Printf("Get Game %d wager_issue error, retcode=%d\n", game_id, retcode)
		*wager_issue = -1
		//关闭连接
		response.Body.Close()
		return
	}
	str := ""
	for _, item := range result["Format02"].([]interface{}) {
		str = fmt.Sprintf("%s", item.(map[string]interface{})["Issue"])
	}
	//fmt.Println(str)
	*wager_issue, _ = strconv.Atoi(str)

	//关闭连接
	response.Body.Close()
}

var wg sync.WaitGroup

//投注
func wager(game_id, wager_issue, index int) {
	url := url_map["wager"]
	//填充请求信息
	request := map[string]interface{}{
		"Format01": map[string]interface{}{
			"UserID":       "",
			"TimeStamp":    "",
			"Access-Token": "",
			"Lang":         "en_US",
			"GameID":       game_id,
			"WagerIssue":   wager_issue,
			"TickSN":       "",
			"WagerType":    "4",
			"WagerMoney":   "6700",
			"MultiIssue":   "1",
			"PayMode":      "0",
			"couponHid":    "",
			"couponCode":   "",
		},
		"Format02": []map[string]interface{}{
			{
				"Num":           "6|-9|-0|",
				"BetSN":         1,
				"PlayTypeID":    "2000",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 0,
			},
			{
				"Num":           "6|-9|-0|",
				"BetSN":         2,
				"PlayTypeID":    "2272",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 0,
			},
			{
				"Num":           "6|-255|-255|",
				"BetSN":         3,
				"PlayTypeID":    "2213",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 0,
			},
			{
				"Num":           "255|-1|2|3|4|-255|",
				"BetSN":         4,
				"PlayTypeID":    "2214",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 1,
			},
			{
				"Num":           "1|2|3|-4|5|6|-7|8|9|",
				"BetSN":         5,
				"PlayTypeID":    "2020",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 1,
			},
			{
				"Num":           "1|2|3|-4|5|6|-7|8|9|",
				"BetSN":         6,
				"PlayTypeID":    "2273",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 1,
			},
			{
				"Num":           "10|",
				"BetSN":         7,
				"PlayTypeID":    "2617",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 0,
			},
			{
				"Num":           "10|-11|",
				"BetSN":         8,
				"PlayTypeID":    "2618",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 1,
			},
			{
				"Num":           "12|",
				"BetSN":         9,
				"PlayTypeID":    "2619",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 0,
			},
			{
				"Num":           "12|-13|",
				"BetSN":         10,
				"PlayTypeID":    "2620",
				"WagerMode":     1,
				"MultiTimes":    1,
				"SubPlayTypeID": 1,
			},
		},
	}
	//fmt.Println(reflect.TypeOf(request["Format01"]), request["Format01"])

	val := player_set[index]

	//初始化客户端
	client := &http.Client{}

	//填充请求参数
	//data["Format01"]["UserID"] type interface {} does not support indexing，需要使用类型断言，将data["Format01"]转换为map[string]interface{}类型
	request["Format01"].(map[string]interface{})["UserID"] = val.playercode
	request["Format01"].(map[string]interface{})["TimeStamp"] = string(time.Now().Unix())
	request["Format01"].(map[string]interface{})["Access-Token"] = val.token
	request["Format01"].(map[string]interface{})["TickSN"] = tickSN
	//fmt.Println(request)

	//请求map转换为json
	data, _ := json.Marshal(request)
	//fmt.Println(string(data))

	//发起http请求
	requestData := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", url, requestData)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("http post error, %v\n", err)

		//关闭连接
		response.Body.Close()
		wg.Done()
		return
	}
	body, _ := ioutil.ReadAll(response.Body)
	//fmt.Println(string(body))

	//判断返回码是否等于0
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	retcode, _ := strconv.Atoi(result["RetCode"].(string))
	if retcode == 0 {
		mutex.Lock()
		total_ticket++
		fmt.Println("total wager ticket count: ", total_ticket)
		mutex.Unlock()
	} else {
		fmt.Printf("GameID:%d, WagerIssue:%d, UserID:%s, RetCode:%d\n", game_id, wager_issue, val.playerID, retcode)
	}
	//关闭连接
	response.Body.Close()
	wg.Done()
}

func main() {
	count := 1
	game_id := 11001
	wager_issue := 0
	total_ticket = 0
	//tickSN = time.Now().Unix()
	tickSN = 1
	fmt.Println(tickSN)

	//用户登陆，更新token
	player_login(count)
	//fmt.Println(player_set)

	//获取游戏期信息
	get_game_issue(game_id, &wager_issue)
	fmt.Println(wager_issue)
	if wager_issue <= 0 {
		return
	}

	//循环投注
	for index := 0; index <= 12; index++ {
		if index == 12 {
			index = 0
		}
		wg.Add(1)
		go wager(game_id, wager_issue, index)
		tickSN++
		time.Sleep(10 * time.Millisecond)
	}
	wg.Wait()
}
