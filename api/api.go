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
	Datestr  string     `json:"date"`
	Ch       string     `json:"channel_name"`
	Url      string     `json:"url"`
	Epg_data []Datalist `json:"epg_data"`
	Success  string     `json:"success"`
}

func ss(iso8601Time string) (s string) {
	ss := iso8601Time
	return ss[20:]
}

func iso8601_uni(iso8601Time string) (t string) {
	sf := "2006-01-02T15:04:05+" + ss(iso8601Time)
	//fmt.Println(sf)
	//result, err := time.ParseInLocation("2006-01-02T15:04:05+0000", iso8601Time, time.Local)
	result, err := time.ParseInLocation(sf, iso8601Time, time.Local)
	//如果错误则退出
	if err != nil {
		fmt.Println(err)
	}
	timestr := fmt.Sprintf("%02d", result.Hour()) + ":" + fmt.Sprintf("%02d", result.Minute())
	return timestr
}

func teshu(iso8601Time string) (t string) {
	sf := "2006-01-02T15:04:05+" + ss(iso8601Time)
	result, err := time.ParseInLocation(sf, iso8601Time, time.Local)
	//result, err := time.ParseInLocation("2006-01-02T15:04:05+0000", iso8601Time, time.Local)
	//如果错误则退出
	if err != nil {
		fmt.Println(err)
	}
	timestr := strconv.Itoa(result.Year()) + "-" + fmt.Sprintf("%02d", result.Month()) + "-" + fmt.Sprintf("%02d", result.Day())
	//fmt.Println(timestr)
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

			//if strings.ToLower(v["ch"]) == strings.ToLower(ch) && teshu(v["startstr"]) == datestr {
			if strings.EqualFold(v["chid"], ch) && teshu(v["startstr"]) == datestr {
				//fmt.Println(ch + v["title"])
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
	response, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
