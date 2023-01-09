package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var wggg sync.WaitGroup
var room int64 = 24326168 //输入房间号  这也找不到建议速速remake

//cookie1=  后面请直接放内容  不要放入cookie:  直接从 buvid3= 或 _uuid=开始复制完成就行
//csrf1=  后面请直接放内容  在cookie的bili_jct里面  不要放入 bili_jct=
//如果第一条的返回结果不是 举办成功  那应该是数据放错了  

//建议用小号挂着

var cookie1 string = "" //

var csrf string = "" //     cookie里的bili_jct

func main() {
	wggg.Add(2)
	
	go danmu()   
	
	go zbj()  //如果不想举报直播间 注释这一行即可
	
	wggg.Wait()
}

func danmu() {

	for {
		time.Sleep(time.Second * 5) //几秒一次？  底边调10也行 反正没弹幕   弹幕很快的就调0吧
		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("https://api.live.bilibili.com/ajax/msg?roomid=%v", room), nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("authority", "api.live.bilibili.com")
		req.Header.Set("pragma", "no-cache")
		req.Header.Set("cache-control", "no-cache")
		req.Header.Set("upgrade-insecure-requests", "1")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
		req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
		req.Header.Set("sec-fetch-site", "none")
		req.Header.Set("sec-fetch-mode", "navigate")
		req.Header.Set("sec-fetch-user", "?1")
		req.Header.Set("sec-fetch-dest", "document")
		req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%s\n", bodyText)
		//miao, err := jsonparser.GetString(bodyText, "data")
		for i := 0; i < 5; i++ {
			text := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.text", i)).String()
			uid := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.uid", i)).Int()
			ts := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.check_info.ts", i)).Int()
			nickname := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.nickname", i)).String()

			seqing(text, uid, ts, nickname)

		}
		for i := 5; i < 9; i++ {
			text := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.text", i)).String()
			uid := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.uid", i)).Int()
			ts := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.check_info.ts", i)).Int()
			nickname := gjson.Parse(string(bodyText)).Get(fmt.Sprintf("data.room.%v.nickname", i)).String()
			ruma(text, uid, ts, nickname)
		}
		//i := 0
	}
}
func seqing(text string, uid int64, ts int64, nickname string) {
	strData := url.QueryEscape(text)
	//str := url.QueryEscape(nickname)
	seqin := url.QueryEscape("低俗色情")

	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf("id=0&roomid=%v&tuid=%v&msg=%v&reason=%v&ts=%v&sign=&reason_id=2&token=&dm_type=0&csrf_token=%v&csrf=%v&visit_id=", room, uid, strData, seqin, ts, csrf, csrf))
	req, err := http.NewRequest("POST", "https://api.live.bilibili.com/xlive/web-ucenter/v1/dMReport/Report", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "api.live.bilibili.com")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("cookie", fmt.Sprintf("%v", cookie1))
	req.Header.Set("origin", "https://live.bilibili.com")
	req.Header.Set("referer", fmt.Sprintf("https://live.bilibili.com/%v?spm_id_from=444.41.live_users.item.click", room))
	req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Microsoft Edge";v="108"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.76")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)
	hlen := len(bodyText)
	if hlen < 53 {
		fmt.Printf("已举办---%v---%v---低俗色情\n", nickname, text)
	} else {
		fmt.Printf("举办成功---%v---%v---低俗色情\n", nickname, text)
	}
}

func ruma(text string, uid int64, ts int64, nickname string) {
	strData := url.QueryEscape(text)
	//str := url.QueryEscape(nickname)
	seqin := url.QueryEscape("辱骂引战")

	client := &http.Client{}
	var data = strings.NewReader(fmt.Sprintf("id=0&roomid=%v&tuid=%v&msg=%v&reason=%v&ts=%v&sign=&reason_id=4&token=&dm_type=0&csrf_token=%v&csrf=%v&visit_id=", room, uid, strData, seqin, ts, csrf, csrf))
	req, err := http.NewRequest("POST", "https://api.live.bilibili.com/xlive/web-ucenter/v1/dMReport/Report", data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "api.live.bilibili.com")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36")
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("origin", "https://live.bilibili.com")
	req.Header.Set("sec-fetch-site", "same-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("referer", fmt.Sprintf("https://live.bilibili.com/%v?spm_id_from=444.41.live_users.item.click", room))
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cookie", fmt.Sprintf("%v", cookie1))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("%s\n", bodyText)
	hlen := len(bodyText)
	if hlen < 53 {
		fmt.Printf("已举办---%v---%v---辱骂引战\n", nickname, text)
	} else {
		fmt.Printf("举办成功---%v---%v---辱骂引战\n", nickname, text)
	}
}

// 举报直播间
func zbj() {
	for {

		time.Sleep(time.Second * 100) //100秒举报一次感觉差不多了  太频繁似乎有问题？

		neirong := "%E6%81%B6%E6%84%8F%E6%8C%82%E6%9C%BA%E8%BE%B1%E9%AA%82%E8%89%B2%E6%83%85%E6%94%BF%E6%B2%BB%E6%93%A6%E8%BE%B9"

		client := &http.Client{}
		var data = strings.NewReader(fmt.Sprintf("room_id=%v&picUrl=&reason=%v&csrf_token=%v&csrf=%v&visit_id=", room, neirong, csrf, csrf))
		req, err := http.NewRequest("POST", "https://api.live.bilibili.com/liveact/report_room", data)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("authority", "api.live.bilibili.com")
		req.Header.Set("accept", "application/json, text/plain, */*")
		req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
		req.Header.Set("cookie", fmt.Sprintf("%v"))
		req.Header.Set("origin", "https://live.bilibili.com")
		req.Header.Set("referer", fmt.Sprintf("https://live.bilibili.com/%v?spm_id_from=333.1007.0.0", room))
		req.Header.Set("sec-ch-ua", `"Not?A_Brand";v="8", "Chromium";v="108", "Microsoft Edge";v="108"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-site")
		req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36 Edg/108.0.1462.76")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		text := gjson.Parse(string(bodyText)).Get("data").Int()
		if text == 0 {
			fmt.Printf("成功举报一次直播间\n")
		} else {
			fmt.Printf("%v\n", bodyText)
		}
	}
}
