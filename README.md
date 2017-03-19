# Tool to detect internet accounts
## Dependencies
- Go enviroment
- selenium server
    <br>
    If you are on MacOS, it's available through `brew install selenium-server-standalone`.
    <br>
    If you are on Ubuntu, please refer to the `https://christopher.su/2015/selenium-chromedriver-ubuntu/`.
- ChromeDriver

## Install
You need to install following go packages 
- go get "sourcegraph.com/sourcegraph/go-selenium"
- go get "github.com/golang/glog"
- Then go build InternetAccountDetector.go

## Usage
```
selenium-server -port 4444
Usage of InternetAccountDetector:
  -alsologtostderr
    	log to standard error as well as files
  -log_backtrace_at value
    	when logging hits line file:N, emit a stack trace
  -log_dir string
    	If non-empty, write log files in this directory
  -logtostderr
    	log to standard error instead of files
  -max_cnt int
    	the number of accounts will try (default 500)
  -start_user string
    	 (default "2014E8008744100")
  -stderrthreshold value
    	logs at or above this threshold go to stderr
  -v value
    	log level for V logs
  -vmodule value
    	comma-separated list of pattern=N settings for file-filtered logging
   ```
    
## TODO 
- store the succeed accounts
- may need to handle the verification code

