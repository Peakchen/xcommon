package serverConfig

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/Peakchen/xgameCommon/Config"
	"github.com/Peakchen/xgameCommon/akLog"
)

/*
	export from ExternalGWConfig.json by tool.
*/
type TExternalgwconfigBase struct {
	Id         int32  `json:"id"`
	Listenaddr string `json:"ListenAddr"`
	Pprofaddr  string `json:"PProfAddr"`
	Name       string
}

type TExternalgwconfigConfig struct {
	data []*TExternalgwconfigBase
}

type tArrExternalgwconfig []*TExternalgwconfigBase

var (
	GExternalgwconfigConfig *TExternalgwconfigConfig = &TExternalgwconfigConfig{}
	cstExternalDef                                   = "ExternalGateway"
)

func init() {
	akLog.FmtPrintln("load	ExternalGWConfig.json")
}

func loadExternalgwConfig() {
	var (
		ExternalGWpath string
	)
	if len(SvrPath) == 0 {
		ExternalGWpath = getserverpath()
	}
	ExternalGWpath = filepath.Join(SvrPath, "ExternalGWConfig.json")
	Config.ParseJson2Cache(GExternalgwconfigConfig, &tArrExternalgwconfig{}, ExternalGWpath)
}

func (this *TExternalgwconfigConfig) ComfireAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrExternalgwconfig)
	errlist = []string{}
	for _, item := range *cfg {
		if len(item.Listenaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("ExternalGWConfig listeraddr invalid, id: %v.", item.Id))
		}

		if len(item.Pprofaddr) == 0 {
			errlist = append(errlist, fmt.Sprintf("ExternalGWConfig Pprofaddr invalid, id: %v.", item.Id))
		}
	}
	return
}

func (this *TExternalgwconfigConfig) DataRWAct(data interface{}) (errlist []string) {
	cfg := data.(*tArrExternalgwconfig)
	this.data = []*TExternalgwconfigBase{}
	for _, item := range *cfg {
		item.Name = cstExternalDef + "_" + strconv.Itoa(int(item.Id))
		this.data = append(this.data, item)
	}
	return
}

func (this *TExternalgwconfigConfig) Get(idx int) *TExternalgwconfigBase {
	if idx >= len(this.data) {
		return nil
	}
	return this.data[idx]
}
