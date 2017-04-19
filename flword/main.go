package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"flword/pkg"

	"github.com/astaxie/beego/logs"
)

var (
	firstStart      = true
	syncTime        time.Time
	corpWords       = make(map[string]map[string]string)
	corpIds         = map[string]string{}
	httpJsonRespFmt = `{"err": %v,"msg": %s,"flds": %v,"rows":%v}`
)

var cfg = flag.String("config", "cfg.json", "please use --config=config file path")

func main() {
	pkg.InitConfig(*cfg)
	pkg.InitDB()
	pkg.InitLog()
	go task()
	http.HandleFunc("/flword_check", flwords_filter)
	logs.Info("sync application is started,and listning at :%v", pkg.Config.Addr)
	http.ListenAndServe(pkg.Config.Addr, nil)
}

func flwords_filter(w http.ResponseWriter, r *http.Request) {
	reqdata, err := ioutil.ReadAll(r.Body)
	logs.Info("request data is: %v", string(reqdata))
	if err != nil {
		fmt.Fprintf(w, httpJsonRespFmt, 1, `"`+err.Error()+`"`, `[]`, `[]`)
		return
	}
	contents := strings.Replace(string(reqdata), "&nbsap", "", -1)
	contents = strings.Replace(contents, `\t`, "", -1)
	contents = strings.Replace(contents, `\n`, "", -1)
	contents = strings.Replace(contents, `\r`, "", -1)
	contents = strings.Join(strings.Fields(contents), "")
	logs.Debug("resolved data is: %v", contents)
	req := make(map[string]string)
	err = json.Unmarshal([]byte(contents), &req)
	if err != nil {
		fmt.Fprintf(w, httpJsonRespFmt, 1, `"`+err.Error()+`"`, `[]`, `[]`)
		return
	}
	corp_id := req["corp_id"]
	flds := []string{"word", "cnt"}
	rows := [][]string{}
	words := corpWords[corp_id]
	for _, v := range words {
		num := strings.Count(contents, v)
		if num != 0 {
			word := []string{v, strconv.Itoa(num)}
			rows = append(rows, word)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	logs.Debug("response data is: \n%v", fmt.Sprintf(httpJsonRespFmt, 0, `""`, pkg.ToJson(flds), pkg.ToJson(rows)))
	fmt.Fprintf(w, httpJsonRespFmt, 0, `""`, pkg.ToJson(flds), pkg.ToJson(rows))
}

func query_corpid() error {
	rows, err := pkg.Db.Query(`select id from corps`)
	if err != nil {
		return err
	}
	for rows.Next() {
		corpId := ""
		rows.Scan(&corpId)
		if corpId != "" {
			corpIds[corpId] = corpId
		}
	}
	return nil
}

func query_words(corp_id string) error {
	//首先查询公司所属的敏感词库
	//如果没有在查询公共敏感词库
	sql := `select word ,is_del from flwords where corp_id='` + corp_id + `'`
	if !firstStart {
		if syncTime.Format("2006-01-02 15:04:05") != `0001-01-01 00:00:00` {
			sql += ` and revise_at > '` + syncTime.Format("2006-01-02 15:04:05") + `'`
		}
	}

	rows, err := pkg.Db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	wds := make(map[string]string)
	for rows.Next() {
		word := ""
		is_del := -1
		rows.Scan(&word, &is_del)
		if word != "" {
			if is_del == 1 {
				delete(corpWords[corp_id], word)
			} else {
				wds[word] = word
			}
		}
	}
	corpWords[corp_id] = wds
	if len(wds) == 0 {
		sql := `select word,is_del from flwords where corp_id=''`
		rs, err := pkg.Db.Query(sql)
		if err != nil {
			return err
		}
		defer rs.Close()
		wds := make(map[string]string)
		for rs.Next() {
			word := ""
			is_del := -1
			rs.Scan(&word, &is_del)
			if word != "" {
				if is_del == 1 {
					delete(corpWords[corp_id], word)
				} else {
					wds[word] = word
				}
			}
		}
		corpWords["99999"] = wds
	}
	return nil
}

func task() {
	tick := time.NewTicker(time.Second * time.Duration(pkg.Config.Interval))
	for {
		select {
		case <-tick.C:
			syncTime = time.Now()
			if err := query_corpid(); err != nil {
				logs.Error("query the corpid err:%v", err)
				continue
			}
			for _, corp_id := range corpIds {
				err := query_words(corp_id)
				if err != nil {
					logs.Error("sync the user to redis failed,the err:%v", err)
					continue
				}
			}
			logs.Info("sync the user to memory successed")
		}
	}
}
