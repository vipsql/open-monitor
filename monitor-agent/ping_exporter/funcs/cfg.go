package funcs

import (
	"encoding/json"
	"log"
	"sync"
	"os/exec"
	"os"
	"bufio"
	"strings"
	"io"
	"io/ioutil"
)

type TransferConfig struct {
	Addrs    []string `json:"addrs"`
	Interval int      `json:"interval"`
	Timeout  int      `json:"timeout"`
	Sn       int      `json:"sn"`
}

type OpenFalconConfig struct {
	Enabled  bool     `json:"enabled"`
	Transfer      *TransferConfig   `json:"transfer"`
}

type PrometheusCOnfig struct {
	Enabled  bool     `json:"enabled"`
	Port  string  `json:"port"`
	Path  string  `json:"path"`
}

type IpSourceConfig struct {
	Const  SourceConstConfig  `json:"const"`
	File  SourceFileConfig  `json:"file"`
	Remote  SourceRemoteConfig  `json:"remote"`
	Listen  SourceListenConfig  `json:"listen"`
}

type SourceConstConfig struct {
	Enabled  bool     `json:"enabled"`
	Ips  []string  `json:"ips"`
	Weight  int  `json:"weight"`
}

type SourceFileConfig struct {
	Enabled  bool     `json:"enabled"`
	Path  string  `json:"path"`
	Weight  int  `json:"weight"`
}

type SourceRemoteConfig  struct {
	Enabled  bool     `json:"enabled"`
	GroupTag  string  `json:"group_tag"`
	Url  string  `json:"url"`
	Interval  int  `json:"interval"`
	Weight  int  `json:"weight"`
}

type SourceListenConfig  struct {
	Enabled  bool     `json:"enabled"`
	Port  string  `json:"port"`
	Path  string  `json:"path"`
	Weight  int  `json:"weight"`
}

type MetricConfig  struct {
	Default  string  `json:"default"`
	CountNum  string  `json:"count_num"`
	CountSuccess  string  `json:"count_success"`
	CountFail  string  `json:"count_fail"`
}

type GlobalConfig struct {
	Debug         bool              `json:"debug"`
	Interval      int               `json:"interval"`
	OpenFalcon    OpenFalconConfig  `json:"open-falcon"`
	Prometheus    PrometheusCOnfig  `json:"prometheus"`
	IpSource      IpSourceConfig    `json:"ip_source"`
	Metrics       MetricConfig      `json:"metrics"`
}

var (
	ConfigFile string
	config     *GlobalConfig
	lock       = new(sync.RWMutex)
)

func Config() *GlobalConfig {
	lock.RLock()
	defer lock.RUnlock()
	return config
}

func ParseConfig(cfg string) error {
	if cfg == "" {
		log.Fatalln("use -c to specify configuration file")
	}
	_, err := os.Stat(cfg)
	if os.IsExist(err) {
		log.Fatalln("config file not found")
		return err
	}
	b,err := ioutil.ReadFile(cfg)
	if err != nil {
		log.Printf("read file %s error %v \n", cfg, err)
		return err
	}
	configContent := strings.TrimSpace(string(b))
	var c GlobalConfig
	err = json.Unmarshal([]byte(configContent), &c)
	if err != nil {
		log.Fatalln("parse config file:", cfg, "fail:", err)
		return err
	}
	lock.Lock()
	defer lock.Unlock()
	config = &c
	log.Println("read config file:", cfg, "successfully")
	return nil
}

func Uuid() (string) {
	commandName := "/usr/sbin/dmidecode"
	params := []string{"|", "grep UUID"}
	cmd := exec.Command(commandName, params...)
	//显示运行的命令
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalln(os.Stderr, "error=>", err.Error())
		return ""
	}
	cmd.Start() // Start开始执行c包含的命令，但并不会等待该命令完成即返回。Wait方法会返回命令的返回状态码并在命令返回后释放相关的资源。

	reader := bufio.NewReader(stdout)

	var index int
	var uuid string
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		s := strings.Split(line, ":")
		t := strings.TrimSpace(s[0])
		if t == "UUID" {
			uuid = strings.TrimSpace(s[1])
		}
		index++
	}
	cmd.Wait()
	return uuid
}