package serverConfig

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Peakchen/xgameCommon/Config"
	"github.com/Peakchen/xgameCommon/akLog"
)

/*
	export from mgoconfig.json by tool.
*/
type TMgoconfigBase struct {
	Id            int32  `json:"id"`
	Passwd        string `json:"Passwd"`
	Username      string `json:"UserName"`
	Shareusername string `json:"ShareUserName"`
	Host          string `json:"Host"`
	Sharehost     string `json:"ShareHost"`
	Sharepasswd   string `json:"SharePasswd"`
	Pprofaddr     string `json:"PProfAddr"`
	Name          string
}

type TMgoconfigConfig struct {
	data []*TMgoconfigBase
}

type tArrMgoconfig []*TMgoconfigBase

var (
	GMgoconfigConfig *TMgoconfigConfig = &TMgoconfigConfig{}
	cstMgoDef                          = "mongo"
)

func init() {
	akLog.FmtPrintln("load	mgoconfig.json")
}

func loadMgoConfig() {
	var (
		mgopath string
	)
	if len(SvrPath) == 0 {
		mgopath = getserverpath()
	}
	mgopath = filepath.Join(SvrPath, "mgoconfig.json")
	Config.ParseJson2Cache(GMgoconfigConfig, &tArrMgoconfig{}, mgopath)
}

func (this *TMgoconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrMgoconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Host) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Host invalid, id: %v.", item.Id))
		}

		if len(item.Sharehost) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Sharehost invalid, id: %v.", item.Id))
		}

		if len(item.Sharehost) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Sharehost invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("mgoconfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TMgoconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrMgoconfig)
	errlist = []string{}
	this.data = []*TMgoconfigBase{}
	for _, item := range *cfg {
		item.Name = cstMgoDef + "_" + strconv.Itoa(int(item.Id))
		this.data = append(this.data, item)
	}
	return
}

func (this *TMgoconfigConfig) Get(idx int) *TMgoconfigBase {
	if idx >= len(this.data) {
		return nil
	}
	return this.data[idx]
}
