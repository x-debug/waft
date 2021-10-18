package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"waft/config"
	"waft/internal"
	"waft/pkg"
	"waft/pkg/array"
)

func setupSignal(pid string) {
	c := make(chan os.Signal)

	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				_ = os.Remove(pid)
				log.Fatalln("Shutdown Waft...")
			default:
			}
		}
	}()
}

var (
	cfg       string
	allowCmds []string
)

func init() {
	allowCmds = []string{"start", "stop", "restart"}
}

func help() string {
	return "waft start|stop|restart"
}

func start(conf *config.ProxyConf) {
	exist := pkg.FileExist(conf.Pid)
	if exist {
		log.Fatalln("pid file exist")
	}

	curPid := os.Getpid()
	err := ioutil.WriteFile(conf.Pid, []byte(strconv.Itoa(curPid)), 0777)
	if err != nil {
		log.Fatalln("write pid error ", err.Error())
	}
	proxy := internal.NewProxy(conf)
	err = proxy.Run()
	if err != nil {
		stop(conf.Pid, false)
	}
}

func stop(pidFile string, sig bool) {
	exist := pkg.FileExist(pidFile)
	if exist {
		if sig {
			fbuf, err := ioutil.ReadFile(pidFile)
			if err != nil {
				log.Fatalln("read pid file error", err.Error())
			}
			pid, _ := strconv.Atoi(string(fbuf))
			err = syscall.Kill(pid, 9)
			if err != nil {
				log.Fatalln("kill proxy error", err.Error())
			}
		}

		_ = os.Remove(pidFile)
	}
}

func startCmd(conf *config.ProxyConf) {
	if len(os.Args) < 2 {
		log.Fatalln(help())
	}

	cmd := os.Args[1]
	if !array.Contains(allowCmds, cmd) {
		log.Fatalln(help())
	}

	if cmd == "start" {
		start(conf)
	} else if cmd == "stop" {
		stop(conf.Pid, true)
	} else {
		stop(conf.Pid, true)
		start(conf)
	}
}

func main() {
	flag.StringVar(&cfg, "config", "./conf_example/waft.yml", "config of server")
	flag.Parse()
	file, err := os.Open(cfg)
	if err != nil {
		log.Fatalln("open config file error: ", err.Error())
	}

	conf, err := config.LoadProxyConf(file)
	if err != nil {
		log.Fatalln("load proxy error: ", err.Error())
	}
	_ = file.Close() //close file immediately
	setupSignal(conf.Pid)
	startCmd(conf)
}
