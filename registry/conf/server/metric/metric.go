package metric

import (
	"fmt"

	"github.com/asaskevich/govalidator"
	"github.com/micro-plat/hydra/registry/conf"
)

type IMetric interface {
	GetConf() (*Metric, bool)
}

type Metric struct {
	Host     string `json:"host" valid:"requrl,required"`
	DataBase string `json:"dataBase" valid:"ascii,required"`
	Cron     string `json:"cron" valid:"ascii,required"`
	UserName string `json:"userName,omitempty" valid:"ascii"`
	Password string `json:"password,omitempty" valid:"ascii"`
	Disable  bool   `json:"disable,omitempty"`
	*option
}

//NewMetric 构建api server配置信息
func NewMetric(host string, db string, cron string, opts ...Option) *Metric {
	m := &Metric{
		Host:     host,
		DataBase: db,
		Cron:     cron,
		option:   &option{},
	}
	for _, opt := range opts {
		opt(m.option)
	}
	return m
}

//GetConf 设置metric
func GetConf(cnf conf.IMainConf) (metric *Metric) {
	_, err := cnf.GetSubObject("metric", &metric)
	if err != nil && err != conf.ErrNoSetting {
		panic(err)
	}
	if err == conf.ErrNoSetting {
		metric.Disable = true
		return
	}
	if b, err := govalidator.ValidateStruct(&metric); !b {
		panic(fmt.Errorf("metric配置有误:%v", err))
	}
	return
}