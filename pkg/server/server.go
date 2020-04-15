package server

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/kvaps/dnsmasq-controller/pkg/conf"
	ctrl "sigs.k8s.io/controller-runtime"
)

var (
	serverLog = ctrl.Log.WithName("server")
)

func Start() error {
	config := conf.GetConfig()
	oldGen := 0

	dnsmasqBinary, err := exec.LookPath("dnsmasq")
	if err != nil {
		panic("dnsmasq binary is not found!")
	}

	err = setupDir(config.DnsmasqConfDir, config.CleanupDir)
	if err != nil {
		panic(err)
	}

	args := []string{
		"dnsmasq",
		"--no-daemon",
		"--no-hosts",
		"--conf-dir=" + config.DnsmasqConfDir,
	}

	if config.EnableDNS {
		args = append(args,
			"--addn-hosts="+config.DnsmasqConfDir+"/hosts",
		)
	} else {
		args = append(args, "--port=0")
	}
	if config.EnableDHCP {
		args = append(args,
			"--dhcp-hostsfile="+config.DnsmasqConfDir+"/dhcp-hosts",
			"--dhcp-optsfile="+config.DnsmasqConfDir+"/dhcp-opts",
		)
	}

	args = append(args, config.DnsmasqOptions...)

	serverLog.Info("Starting dnsmasq: " + strings.Join(args, " "))
	cmd := serverStart(dnsmasqBinary, args)
	go func() {
		for {
			newGen := config.Generation
			time.Sleep(time.Duration(config.SyncDelay) * time.Second)
			if newGen != oldGen {
				serverLog.Info("Configuration changed, restarting dnsmasq.")
				serverStop(cmd)
				cmd = serverStart(dnsmasqBinary, args)
				serverLog.Info("Configuration reloaded.")
				oldGen = newGen
			}
		}
	}()

	return nil
}

func serverStart(dnsmasqBinary string, args []string) *exec.Cmd {
	cmd := &exec.Cmd{
		Path: dnsmasqBinary,
		Args: args,
	}
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return cmd
}

func serverStop(cmd *exec.Cmd) {
	timer := time.AfterFunc(1*time.Second, func() {
		err := cmd.Process.Kill()
		if err != nil {
			panic(err)
		}
	})
	cmd.Wait()
	timer.Stop()
}

func setupDir(p string, cleanup bool) error {
	dir, err := ioutil.ReadDir(p)
	if cleanup {
		for _, d := range dir {
			err = os.RemoveAll(path.Join([]string{p, d.Name()}...))
			if err != nil {
				return err
			}
		}
	}

	dirs := []string{"/", "/hosts", "/dhcp-hosts", "/dhcp-opts"}
	for _, dir := range dirs {
		if _, err := os.Stat(p + dir); os.IsNotExist(err) {
			err := os.MkdirAll(p+dir, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
