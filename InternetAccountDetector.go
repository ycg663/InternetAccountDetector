package main

import (
	"flag"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"
	"sourcegraph.com/sourcegraph/go-selenium"
)

const (
	DefaultPassword = "ucas"
	DefaultUrl      = "http://www.baidu.com"
)

var (
	StartUserName string
	MaxTryCnt     int
)

type Detector struct {
	webDriver selenium.WebDriver
}

func NewDetector(browserName, remoteAddr string, implicitWaitTimeout uint) *Detector {
	caps := selenium.Capabilities(map[string]interface{}{"browserName": browserName})
	driver, err := selenium.NewRemote(caps, remoteAddr)
	if err != nil {
		glog.Errorf("Failed to open session: %s\n", err)
		return nil
	}
	driver.SetImplicitWaitTimeout(implicitWaitTimeout)
	driver.SetAsyncScriptTimeout(300)
	return &Detector{webDriver: driver}
}

func (self *Detector) Close() {
	if self.webDriver != nil {
		err := self.webDriver.Close()
		if err != nil {
			glog.Errorf("Close web driver failed. %v", err)
		}
		err = self.webDriver.Quit()
		if err != nil {
			glog.Errorf("Quit web driver failed. %v", err)
		}
		return
	}

}

func (self *Detector) Detect(startUserName string) error {
	err := self.webDriver.Get(DefaultUrl)
	if err != nil {
		glog.Errorf("Get endpoint failed. url: %v, %v", DefaultUrl, err)
		return err
	}
	strli := strings.Split(startUserName, "E")

	prefix := strli[0]
	suffix := ""
	suffixNumber := int64(0)
	if len(strli) > 1 {
		suffix = strli[1]
		suffixNumber, err = strconv.ParseInt(suffix, 10, 64)
		if err != nil {
			glog.Errorf("Invalid start username, %v, %v", startUserName, err)
		}
	} else {
		suffixNumber, err = strconv.ParseInt(prefix, 10, 64)
		if err != nil {
			glog.Errorf("Invalid start username, %v, %v", startUserName, err)
		}
	}

	for i := 0; i < MaxTryCnt; i++ {

		self.webDriver.ExecuteScriptAsync("document.getElementById('username').style.display='inline-block';", []interface{}{})
		self.webDriver.ExecuteScriptAsync("document.getElementById('pwd').style.display='inline-block';", []interface{}{})
		time.Sleep(300 * time.Millisecond)

		userNameEle, err := self.webDriver.FindElement(selenium.ByXPATH, "//input[@name='username']")
		if err != nil {
			glog.Errorf("Find user name input element failed. %v", err)
		}

		passwordEle, err := self.webDriver.FindElement(selenium.ByXPATH, "//input[@name='pwd']")
		if err != nil {
			glog.Errorf("Find user password input element failed. %v", err)
		}

		userNameEle.Clear()
		time.Sleep(time.Second)
		var username = ""
		if suffix != "" {
			curSuffix := strconv.Itoa(int(suffixNumber))
			username = prefix + "E" + curSuffix
		} else {
			username = strconv.Itoa(int(suffixNumber))
		}
		glog.Infof("username: %v\n", username)

		suffixNumber += int64(1)

		userNameEle.Click()
		err = userNameEle.SendKeys(username)
		if err != nil {
			glog.Errorf("Send key to username input failed. %v", err)
		}

		passwordEle.Clear()
		passwordEle.Click()
		err = passwordEle.SendKeys(DefaultPassword)
		if err != nil {
			glog.Errorf("Send key to password input failed. %v", err)
		}

		loginButton, err := self.webDriver.FindElement(selenium.ByXPATH, "//a[@id='loginLink']")
		if err != nil {
			glog.Errorf("Find login button failed. %v", err)
		}
		err = loginButton.Click()
		if err != nil {
			glog.Errorf("Click login button failed. %v", err)
		}
		time.Sleep(500 * time.Millisecond)
		pageSrc, err := self.webDriver.PageSource()
		if strings.Contains(pageSrc, "继续访问") == true {
			glog.Info("Login finish!")
			return nil
		}
	}

	return nil
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	StartUserName = *(flag.String("start_user", "2014E8008744100", ""))
	MaxTryCnt = *(flag.Int("max_cnt", 500, "the number of accounts will try"))

	flag.Parse()

	detector := NewDetector("chrome", "http://127.0.0.1:4444/wd/hub", 1000)
	if detector == nil {
		glog.Fatal("Crawler initialize failed.")
	}
	defer detector.Close()
	err := detector.Detect(StartUserName)
	if err != nil {
		glog.Error(err)
	}
}
