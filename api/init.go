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
	"github.com/cenkalti/backoff/v4"
	"github.com/robfig/cron/v3"
)

var data_list []map[string]string

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
	data_list = nil
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(file); err != nil {
		fmt.Println(err)
	}

	root := doc.SelectElement("tv")
	//fmt.Println("Root element:", root.Tag)
	iptv_dict1 := make(map[string]string)
	for _, channels := range root.SelectElements("channel") {

		Channel_id := channels.SelectAttrValue("id", "unknown")
		//iptv_dict1["Channel_id"] = Channel_id

		if disp_name := channels.SelectElement("display-name"); disp_name != nil {

			Channel_name := disp_name.Text()
			//iptv_dict1["Channel_name"] = Channel_name
			iptv_dict1[Channel_id] = Channel_name
		}

	}
	//fmt.Println(iptv_dict1)

	for _, programmes := range root.SelectElements("programme") {
		iptv_dict2 := make(map[string]string)
		startstr := programmes.SelectAttrValue("start", "none")
		stopstr := programmes.SelectAttrValue("stop", "none")
		chid := programmes.SelectAttrValue("channel", "none")
		//fmt.Println(chid)
		iptv_dict2["startstr"] = shijian(startstr)
		iptv_dict2["stopstr"] = shijian(stopstr)
		iptv_dict2["chid"] = iptv_dict1[chid]
		//fmt.Print(iptv_dict2["ch"])

		if titlestr := programmes.SelectElement("title"); titlestr != nil {
			title := titlestr.Text()
			if title != "" {
				iptv_dict2["title"] = title
			} else {
				iptv_dict2["title"] = "none"
			}

		}
		if descstr := programmes.SelectElement("desc"); descstr != nil {
			desc := descstr.Text()
			if desc != "" {
				iptv_dict2["desc"] = desc
			} else {
				iptv_dict2["desc"] = "none"
			}

		}
		//fmt.Println(chid, startstr, stopstr, iptv_dict2["title"], iptv_dict2["desc"])
		data_list = append(data_list, iptv_dict2)
	}
	//fmt.Println(data_list)
	Liststr = nil
	Liststr = data_list
}

func downloadxml(urlpath string) error {
	operation := func() error {
		resp, err := http.Get(urlpath)
		if err != nil {
			//log.Fatalf("无法获取文件： %v", err)
			fmt.Printf("无法获取文件：%v", err)
			return err
		}
		//defer resp.Body.Close()
		defer func() {
			_ = resp.Body.Close()
		}()

		// 创建文件用于保存
		filename := path.Base(urlpath)
		flags := os.O_CREATE | os.O_WRONLY
		f, err := os.OpenFile(filename, flags, 0666)
		if err != nil {
			fmt.Println("创建文件失败")
			//log.Fatal("err")
			fmt.Printf("创建文件失败: %v", err)
			return err
		}
		//defer f.Close()
		defer func() {
			_ = f.Close()
		}()
		// 将响应流和文件流对接起来
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			log.Fatal("err")
			return err
		}
		return nil
	}
	//重试  https://www.jianshu.com/p/435364fc51ce
	err := backoff.RetryNotify(
		operation,
		backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 2),
		func(err error, duration time.Duration) {
			log.Printf("failed err:%s,and it will be executed again in %v", err.Error(), duration)
		})
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func do(urlpath, filename string) {
	c := cron.New()
	c.AddFunc("3 5,18 * * *", func() {
		downloadxml(urlpath)
		RedXml(filename)
	})
	c.Start()
}

func init() {
	urlpath := "https://e.erw.cc/all.xml"
	//urlpath := "https://epg.pw/xmltv/epg_CN.xml"
	//urlpath := "https://epg.mxdyeah.top/download/all-mxdyeah.xml"
	filename := path.Base(urlpath)
	do(urlpath, filename)
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("File exists\n")
		RedXml(filename)
	} else {
		fmt.Printf("File does not exist\n")
		downloadxml(urlpath)
		RedXml(filename)
	}
}
