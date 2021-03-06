package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zgs225/alfred-youdao/alfred"
	"github.com/zgs225/youdao"
)

const (
	APPID     = "2f871f8481e49b4c"
	APPSECRET = "CQFItxl9hPXuQuVcQa5F2iPmZSbN0hYS"
	MAX_LEN   = 255
)

func init() {
	log.SetPrefix("[i] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	log.Println(os.Args)

	client := &youdao.Client{
		AppID:     APPID,
		AppSecret: APPSECRET,
	}
	agent := newAgent(client)
	q := strings.TrimSpace(strings.Join(os.Args[1:], " "))
	items := alfred.NewResult()

	if len(q) > 255 {
		items.Append(&alfred.ResultElement{
			Valid:    false,
			Title:    "错误: 最大查询字符数为255",
			Subtitle: q,
		})
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(items)
		fmt.Print(b.String())
		os.Exit(1)
	}

	r, err := agent.Query(q)
	if err != nil {
		panic(err)
	}

	if r.Basic != nil {
		item := alfred.ResultElement{
			Valid:    true,
			Title:    r.Basic.Explains[0],
			Subtitle: r.Basic.Phonetic,
			Arg:      r.Basic.Explains[0],
			Mods: map[string]*alfred.ModElement{
				alfred.Mods_Shift: &alfred.ModElement{
					Valid:    true,
					Arg:      toYoudaoDictUrl(q),
					Subtitle: "回车键打开词典网页",
				},
			},
		}
		items.Append(&item)
	}

	if r.Translation != nil {
		item := alfred.ResultElement{
			Valid:    true,
			Title:    (*r.Translation)[0],
			Subtitle: "翻译结果",
			Arg:      (*r.Translation)[0],
			Mods: map[string]*alfred.ModElement{
				alfred.Mods_Shift: &alfred.ModElement{
					Valid:    true,
					Arg:      toYoudaoDictUrl(q),
					Subtitle: "回车键打开词典网页",
				},
			},
		}
		items.Append(&item)
	}

	if r.Web != nil {
		items.Append(&alfred.ResultElement{
			Valid:    true,
			Title:    "网络释义",
			Subtitle: "有道词典",
		})

		for _, elem := range *r.Web {
			items.Append(&alfred.ResultElement{
				Valid:    true,
				Title:    elem.Key,
				Subtitle: elem.Value[0],
				Arg:      elem.Key,
				Mods: map[string]*alfred.ModElement{
					alfred.Mods_Shift: &alfred.ModElement{
						Valid:    true,
						Arg:      toYoudaoDictUrl(q),
						Subtitle: "回车键打开词典网页",
					},
				},
			})
		}
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(items)
	fmt.Print(b.String())

	if agent.Dirty {
		agent.Cache.SaveFile(CACHE_FILE)
	}
}
