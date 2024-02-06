package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Liststr []map[string]string

type Datalist struct {
	Start string `json:"start"`
	Desc  string `json:"desc"`
	End   string `json:"end"`
	Title string `json:"title"`
}

type Epg_data struct {
	Datestr  string `json:"date"`
	Ch       string `json:"channel_name"`
	Url      string `json:"url"`
	Epg_data []Datalist
	Success  string `json:"success"`
}

func iso8601_uni(iso8601Time string) (t string) {
	result, err := time.ParseInLocation("2006-01-02T15:04:05+0800", iso8601Time, time.Local)
	//如果错误则退出
	if err != nil {
		fmt.Println(err)
		//return -1
	}
	timestr := fmt.Sprintf("%02d", result.Hour()) + ":" + fmt.Sprintf("%02d", result.Minute())
	return timestr
}

func teshu(iso8601Time string) (t string) {
	result, err := time.ParseInLocation("2006-01-02T15:04:05+0800", iso8601Time, time.Local)
	//如果错误则退出
	if err != nil {
		fmt.Println(err)
		//return -1
	}
	timestr := strconv.Itoa(result.Year()) + "-" + fmt.Sprintf("%02d", result.Month()) + "-" + fmt.Sprintf("%02d", result.Day())
	return timestr
}

func Api_Handler(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	ch := vars.Get("ch")
	dates := vars.Get("date")
	datestr := ""
	if dates != "" {
		datestr = dates
	} else {
		datestr = time.Now().Format("2006-01-02")
	}

	data := new(Epg_data)
	data.Datestr = datestr
	data.Url = r.Host
	if ch != "" {
		data.Ch = strings.ToUpper(ch)
		data.Success = "Success"
		for _, v := range Liststr {

			if strings.ToLower(v["chid"]) == strings.ToLower(ch) && teshu(v["startstr"]) == datestr {
				//fmt.Printf("索引：%d, 值: %s\n", i, v)
				datas := Datalist{
					Start: iso8601_uni(v["startstr"]),
					Desc:  v["desc"],
					End:   iso8601_uni(v["stopstr"]),
					Title: v["title"],
				}
				data.Epg_data = append(data.Epg_data, datas)
			}
		}
	} else {
		data.Ch = "none"
		data.Success = "Failed"
		data.Epg_data = nil
	}

	// datas := Datalist{
	// 	Start: "20240209170000 +0800",
	// 	Desc:  "desc:none",
	// 	End:   "20240209185959 +0800",
	// 	Title: "title:none",
	// }
	// data.Epg_data = append(data.Epg_data, datas)

	//fmt.Println(Liststr)
	// for _, v := range Liststr {

	// 	if strings.ToLower(v["chid"]) == strings.ToLower(ch) {
	// 		//fmt.Printf("索引：%d, 值: %s\n", i, v)
	// 		datas := Datalist{
	// 			Start: v["startstr"],
	// 			Desc:  v["desc"],
	// 			End:   v["stopstr"],
	// 			Title: v["title"],
	// 		}
	// 		data.Epg_data = append(data.Epg_data, datas)
	// 	}
	// }
	response, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
