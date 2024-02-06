package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/beevik/etree"
)

var data_list []map[string]string
var data_list1 []map[string]string

var SysTimeLocation, _ = time.LoadLocation("Asia/Chongqing")

func shijian(customTime string) (tstr string) {
	// 使用正则表达式匹配和替换
	re := regexp.MustCompile(`(\d{4})(\d{2})(\d{2})(\d{2})(\d{2})(\d{2})\s(\+\d{4})`)
	iso8601Time := re.ReplaceAllString(customTime, `${1}-${2}-${3}T${4}:${5}:${6}${7}`)
	// 输出转换后的时间
	//fmt.Println(iso8601Time) // 输出: 2024-02-03T05:20:00+0800

	return iso8601Time
}

func RedXml(file string) {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(file); err != nil {
		fmt.Println(err)
	}

	root := doc.SelectElement("tv")
	//fmt.Println("Root element:", root.Tag)

	for _, channels := range root.SelectElements("channel") {
		//fmt.Println("channel element:", channels.Tag)
		iptv_dict1 := make(map[string]string)
		Channel_id := channels.SelectAttrValue("id", "unknown")
		//fmt.Printf("Channel_name:%s,", Channel_id)
		iptv_dict1["Channel_id"] = Channel_id
		if disp_name := channels.SelectElement("display-name"); disp_name != nil {
			//lang := disp_name.SelectAttrValue("lang", "unknown")
			Channel_name := disp_name.Text()
			//fmt.Printf("display_name:%s \n", lang)
			//iptv_dict1["lang"] = lang
			iptv_dict1["Channel_name"] = Channel_name
		}
		data_list1 = append(data_list1, iptv_dict1)
	}

	for _, programmes := range root.SelectElements("programme") {
		iptv_dict2 := make(map[string]string)
		startstr := programmes.SelectAttrValue("start", "none")
		stopstr := programmes.SelectAttrValue("stop", "none")
		chid := programmes.SelectAttrValue("channel", "none")
		//fmt.Printf("start:%s,", startstr)
		//fmt.Printf("stop:%s,", stopstr)
		//fmt.Printf("ch:%s,", chid)
		iptv_dict2["startstr"] = shijian(startstr)
		iptv_dict2["stopstr"] = shijian(stopstr)
		iptv_dict2["chid"] = chid

		if titlestr := programmes.SelectElement("title"); titlestr != nil {
			title := titlestr.Text()
			//lang := titlestr.SelectAttrValue("lang", "unknown")
			//fmt.Printf("title:%s,", titlestr.Text())
			if title != "" {
				iptv_dict2["title"] = title
			} else {
				iptv_dict2["title"] = "none"
			}

		}
		if descstr := programmes.SelectElement("desc"); descstr != nil {
			desc := descstr.Text()
			//fmt.Printf("desc:%s,", desc)
			if desc != "" {
				iptv_dict2["desc"] = desc
			} else {
				iptv_dict2["desc"] = "none"
			}

		}
		data_list = append(data_list, iptv_dict2)

	}
	// for i, v := range data_list {
	// 	fmt.Printf("索引：%d, 值：%s\n", i, v)
	// }

}

// func Make_values() (datalist []map[string]string) {
// 	for i := range data_list {
// 		//fmt.Println(data_list[i])
// 		//fmt.Println(data_list[i]["chid"])
// 		a := data_list[i]["chid"]
// 		for j := range data_list1 {
// 			if data_list1[j]["Channel_id"] == a {
// 				data_list[i]["chid"] = data_list1[j]["Channel_name"]
// 			}
// 		}
// 		//fmt.Printf("索引：%d, 值：%s\n", i, v)
// 	}
// 	return datalist
// }

func Make_values() {
	for i := range data_list {
		//fmt.Println(data_list[i])
		//fmt.Println(data_list[i]["chid"])
		a := data_list[i]["chid"]
		for j := range data_list1 {
			if data_list1[j]["Channel_id"] == a {
				data_list[i]["chid"] = data_list1[j]["Channel_name"]
			}
		}
		//fmt.Printf("索引：%d, 值：%s\n", i, v)
	}
	Liststr = data_list
	// for i, v := range Liststr {
	// 	fmt.Printf("索引：%d, 值：%s\n", i, v)
	// }
}

func downloadxml(urlpath string) {
	resp, err := http.Get(urlpath)
	if err != nil {
		log.Fatalf("无法获取文件： %v", err)
	}
	defer resp.Body.Close()
	// 创建文件用于保存
	filename := path.Base(urlpath)
	flags := os.O_CREATE | os.O_WRONLY
	f, err := os.OpenFile(filename, flags, 0666)
	if err != nil {
		fmt.Println("创建文件失败")
		log.Fatal("err")
	}
	defer f.Close()
	// 将响应流和文件流对接起来
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Fatal("err")
	}
}

func init() {
	do()
	// downloadxml("https://epg.erw.cc/all.xml")
	// RedXml("all.xml")
	// Make_values()

	//定时
	// ticker := time.NewTicker(1 * time.Minute)
	// defer ticker.Stop()

	// for range ticker.C {
	// 	downloadxml("https://epg.erw.cc/all.xml")
	// 	RedXml("all.xml")
	// 	Make_values()
	// }
}

func do() {
	t := time.Now()
	next := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute()+1, 0, 0, SysTimeLocation)
	fmt.Printf("next  type: %T,\t val: %v\n", next, next)
	//获取下次执行时间与当前时间的差
	duration := next.Sub(time.Now())
	fmt.Printf("duration  type: %T,\t val: %v\n", duration, duration)
	/*预约下次执行执行计划，因为在程序初始化的时候已经调用了do()方法，
	*在do()每次执行完，都会再预约下次执行计划，直到主程序die*/
	time.AfterFunc(duration, do)
	fmt.Println("****")
	downloadxml("https://epg.erw.cc/all.xml")
	RedXml("all.xml")
	Make_values()
}