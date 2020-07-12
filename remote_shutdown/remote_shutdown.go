package remote_shutdown

import (
	"errors"
	"remote_shutdown/utils"
	//"os/exec"
	"runtime"
	"strconv"
	"github.com/codeskyblue/go-sh"
	"github.com/go-resty/resty/v2"
)

type RShutdown struct {
	Protocol      string `json:"protocol"`
	Host          string `json:"hostname"`
	Port          int    `json:"port"`
	SecondsToWait int    `json:"wait"`
	SecurityCode  string `json:"token"`
	Path          string `json:"path"`
}

func NewRShutdown(proto, h, s, path string, p, wait int) *RShutdown {
	return  &RShutdown{
		Protocol: proto,
		Host:         h,
		Port: p,
		Path: path,
		SecondsToWait: wait,
		SecurityCode: s,
	}
}

func (rs *RShutdown) URL() string {
	url := ""
	switch rs.Protocol {
	case "http":
		url += "http://"
	case "https":
		url += "https://"
	default:
		url += "http://"
	}

	url += rs.Host

	switch rs.Port {
	case 0:
		url += "/"
	default:
		url += ":" + strconv.Itoa(rs.Port) + "/"
	}

	url += rs.Path
	url += "?s=" + rs.SecurityCode
	return url
}

type Shutdown struct {
	SecondsToWait int
}

// See: https://stackoverflow.com/questions/42660690/is-it-possible-to-shut-down-the-host-machine-by-executing-a-command-on-one-of-it
func (s *Shutdown) Execute() error {
	u := utils.New()
	u.Log.Infof("Reached Shutdown state")
	return nil
	switch runtime.GOOS {
	case "illumos":
		// https://illumos.org/man/1M/shutdown
		//return exec.Command("shutdown", "-y", "-g" + strconv.Itoa(s.SecondsToWait), "-i 5").Run()
		return sh.Command("shutdown", "-y", "-g" + strconv.Itoa(s.SecondsToWait), "-i 5").Run()
	case "dragonfly":
		// https://man.dragonflybsd.org/?command=shutdown&section=8
		//return exec.Command("shutdown", "-p", strconv.Itoa(s.SecondsToWait), "System Shutdown initiated via remote_shutdown").Run()
		return sh.Command("shutdown", "-p", strconv.Itoa(s.SecondsToWait), "System Shutdown initiated via remote_shutdown").Run()
	case "freebsd":
	case "netbsd":
	case "openbsd":
	case "plan9":
	case "solaris":
	case "darwin":
		// https://www.howtogeek.com/512304/how-to-shut-down-your-mac-using-terminal/#:~:text=When%20Terminal%20opens%2C%20type%20sudo,want%20your%20Mac%20to%20restart.
		//return exec.Command("shutdown", "-h", strconv.Itoa(s.SecondsToWait), "System Shutdown initiated via remote_shutdown").Run()
		return sh.Command("shutdown", "-h", strconv.Itoa(s.SecondsToWait), "System Shutdown initiated via remote_shutdown").Run()
	case "windows":
		// https://www.easeus.com/computer-instruction/shutdown-power-off-windows-10-using-cmd.html
		//return exec.Command("shutdown", "/s", "/t", strconv.Itoa(s.SecondsToWait), "/c \"System Shutdown initiated via remote_shutdown\"", "/d 0:0", "/f").Run()
		return sh.Command("shutdown", "/s", "/t", strconv.Itoa(s.SecondsToWait), "/c \"System Shutdown initiated via remote_shutdown\"", "/d 0:0", "/f").Run()
	case "linux":
		//return exec.Command("shutdown", "-P", strconv.Itoa(s.SecondsToWait), "System Shutdown initiated via remote_shutdown").Run()
		return sh.Command("shutdown", "-P", strconv.Itoa(s.SecondsToWait), "System Shutdown initiated via remote_shutdown").Run()
	default:
		return errors.New("Unknown OS " + runtime.GOOS)
	}
	return nil
}

// aix
// android
// js/wasm
// windows

func (rs *RShutdown) Execute(key string) error {
	u := utils.New()
	u.Log.Infof("rs = %#v", rs)
	u.Log.Infof("key = %s", key)
	if rs.SecurityCode == key {
		client := resty.New()
		resp, err := client.R().EnableTrace().SetHeader("Content-Type", "application/json").
			SetBody(rs).Post(rs.URL())
		// Explore response object
		u.Log.Infoln("Response Info:")
		u.Log.Infoln("Error      :", err)
		u.Log.Infoln("Status Code:", resp.StatusCode())
		u.Log.Infoln("Status     :", resp.Status())
		u.Log.Infoln("Proto      :", resp.Proto())
		u.Log.Infoln("Time       :", resp.Time())
		u.Log.Infoln("Received At:", resp.ReceivedAt())
		u.Log.Infoln("Body       :\n", resp)
		u.Log.Infoln()

		// Explore trace info
		u.Log.Infoln("Request Trace Info:")
		ti := resp.Request.TraceInfo()
		u.Log.Infoln("DNSLookup    :", ti.DNSLookup)
		u.Log.Infoln("ConnTime     :", ti.ConnTime)
		u.Log.Infoln("TCPConnTime  :", ti.TCPConnTime)
		u.Log.Infoln("TLSHandshake :", ti.TLSHandshake)
		u.Log.Infoln("ServerTime   :", ti.ServerTime)
		u.Log.Infoln("ResponseTime :", ti.ResponseTime)
		u.Log.Infoln("TotalTime    :", ti.TotalTime)
		u.Log.Infoln("IsConnReused :", ti.IsConnReused)
		u.Log.Infoln("IsConnWasIdle:", ti.IsConnWasIdle)
		u.Log.Infoln("ConnIdleTime :", ti.ConnIdleTime)
		s := &Shutdown{SecondsToWait: rs.SecondsToWait}
		return s.Execute()
	} else {
		return errors.New("Invalid secret key for host: " + rs.Host)
	}
}