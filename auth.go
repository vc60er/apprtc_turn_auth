package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"net/http"
	"time"
)

//curl 'https://networktraversal.googleapis.com/v1alpha/iceconfig?key=AIzaSyAJdh2HkajseEIltlZ3SIXO02Tze9sO3NY' -X POST -H 'origin: https://appr.tc' -H 'accept-encoding: gzip, deflate, br' -H 'accept-language: zh-CN,zh;q=0.8,en;q=0.6' -H 'user-agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36' -H 'accept: */*' -H 'referer: https://appr.tc/' -H 'authority: networktraversal.googleapis.com' -H 'content-length: 0' --compressed
/**
{
  "lifetimeDuration": "86400.000s",
  "iceServers": [
    {
      "urls": [
        "turn:216.58.221.30:19305?transport=udp",
        "turn:[2404:6800:4008:C01::7F]:19305?transport=udp",
        "turn:216.58.221.30:443?transport=tcp"
      ],
      "username": "CPv22b0FEgaR/mw7LOIYzc/s6OMT",
      "credential": "lNbJSs6grl/rFPiapr5ke6GSwt8="
    },
    {
      "urls": [
        "stun:stun.l.google.com:19302"
      ]
    }
  ]
}
*/

type iceServer struct {
	Urls       []string `json:"urls"`
	Username   string   `json:"username"`
	Credential string   `json:"credential"`
}

type turnAuth struct {
	LifetimeDuration string      `json:"lifetimeDuration"`
	IceServers       []iceServer `json:"iceServers"`
}

func test_json() {
	ta := turnAuth{
		LifetimeDuration: "2345",
		IceServers: []iceServer{
			iceServer{
				Username:   "cn",
				Credential: "cn",
				Urls: []string{
					"turn:216.58.221.30:19305?transport=udp",
					"turn:[2404:6800:4008:C01::7F]:19305?transport=udp",
					"turn:216.58.221.30:443?transport=tcp",
				},
			},
		},
	}

	b, _ := json.Marshal(ta)
	fmt.Println(string(b))
}

// curl https://computeengineondemand.appspot.com/turn\?username\=iapprtc\&key\=4080218913
/**
{
    username: "1472037322:iapprtc",
    password: "w8+/YssgjqhFXxaxq/AF5K4SkNo=",
    uris: [
        "turn:107.167.189.134:3478?transport=udp",
        "turn:107.167.189.134:3478?transport=tcp",
        "turn:107.167.189.134:3479?transport=udp",
        "turn:107.167.189.134:3479?transport=tcp"
    ]
}
*/

type turnAuth2 struct {
	Username string   `json:"username"`
	Password string   `json:"password"`
	Uris     []string `json:"uris"`
}

func test_json2() {
	ta := turnAuth2{
		Username: "cn",
		Password: "cn",
		Uris: []string{
			"turn:107.167.189.134:3478?transport=udp",
			"turn:107.167.189.134:3478?transport=tcp",
		},
	}

	b, _ := json.Marshal(ta)
	fmt.Println(string(b))
}

var key = "4080218913" // 这里的 key 是事先设置好的, 我们把他当成一个常亮来看, 所以就不从HTTP请求参数里读取了

func hmac_func(key string, content string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(content))

	input := mac.Sum(nil)

	str := base64.StdEncoding.EncodeToString(input)
	fmt.Println(str)

	return str
}

func test_hmac_func() {
	h := hmac_func(key, "helloworld")
	fmt.Println(h)
	h = hmac_func(key, "hello")
	fmt.Println(h)
}

func turnHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	//  w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE")

	req.ParseForm()

	var time_to_live int = 60 * 60 * 24
	var timestamp float64 = math.Floor(float64(time.Now().UnixNano()/100000000)) + float64(time_to_live)
	var turn_username string = fmt.Sprintf("%d", int(timestamp))
	var password string = hmac_func(key, turn_username)

	var is iceServer
	is.Username = turn_username
	is.Credential = password
	is.Urls = append(is.Urls, "turn:10.58.60.236:9000?transport=udp")
	is.Urls = append(is.Urls, "turn:10.58.60.236:9000?transport=tcp")

	var ta turnAuth
	ta.LifetimeDuration = "86400.000s"
	ta.IceServers = append(ta.IceServers, is)

	b, _ := json.Marshal(ta)
	fmt.Println(string(b))

	fmt.Fprintf(w, string(b))
}

func turnHandler2(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	//  w.Header().Add("Access-Control-Allow-Methods", "POST, DELETE")

	req.ParseForm()

	var time_to_live int = 60 * 60 * 24
	var timestamp float64 = math.Floor(float64(time.Now().UnixNano()/100000000)) + float64(time_to_live)
	var turn_username string = fmt.Sprintf("%d", int(timestamp))
	var password string = hmac_func(key, turn_username)

	var ta turnAuth2
	ta.Username = turn_username
	ta.Password = password
	ta.Uris = append(ta.Uris, "turn:10.58.60.236:9000?transport=udp")
	ta.Uris = append(ta.Uris, "turn:10.58.60.236:9000?transport=tcp")

	b, _ := json.Marshal(ta)
	fmt.Println(string(b))

	fmt.Fprintf(w, string(b))
}

var port = flag.Int("port", 8081, "The TCP port that the server listens on")

func main() {
	flag.Parse()
	fmt.Printf("Starting auth: port = %d\n", *port)

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	http.HandleFunc("/turn", turnHandler)
	http.HandleFunc("/v2/turn", turnHandler2)

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	fmt.Println(err)
}
