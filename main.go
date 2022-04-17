package main

import (
	"flag"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

func main() {
	var sourceFile string
	var desFile string
	var delaySec int64

	flag.StringVar(&sourceFile, "sourceFile", "", "字幕源文件地址 必填")
	flag.StringVar(&desFile, "desFile", "new.srt", "修改后的文件地址 可不填 不填就是当前目录的new.srt")
	flag.Int64Var(&delaySec, "sec", 0, "延迟的时间单位(秒) 可以为负数")
	flag.Parse()
	bytes, _ := ioutil.ReadFile(sourceFile)

	sourceString := BytesToString(bytes)
	sliceStr := strings.Split(sourceString, "\n")
	for i, v := range sliceStr {
		if strings.Contains(v, "-->") {
			last := modifyTime(v, delaySec)
			sliceStr[i] = last
		}
	}
	var lastStr string
	// 组成字符串
	for _, v := range sliceStr {
		lastStr = lastStr + v + "\n"
	}
 	os.WriteFile(desFile,[]byte(lastStr),0644)

	//print(sourceString)

}

func modifyTime(s string, delay int64) string {
	timeSlice := strings.Split(s, " --> ")
	for i, v := range timeSlice {
		time, append := getTime(v)
		time = time + delay
		lastTime := toTime(time)
		lastV := lastTime + "," + append
		timeSlice[i] = lastV
	}
	return timeSlice[0] + " --> " + timeSlice[1]

}

//  输入 113 输出 00:01:53
func toTime(sec int64) string {
	hour := sec / 3600
	minute := (sec % 3600) / 60
	se := sec - 3600*hour - minute*60

	hours := strconv.Itoa(int64toInt(hour))
	mins := strconv.Itoa(int64toInt(minute))
	secs := strconv.Itoa(int64toInt(se))
	if len(hours) == 1 {
		hours = "0" + hours
	}
	if len(mins) == 1 {
		mins = "0" + mins
	}
	if len(secs) == 1 {
		secs = "0" + secs
	}
	return hours + ":" + mins + ":" + secs
}

func int64toInt(i int64) int {
	return *(*int)(unsafe.Pointer(&i))
}

//00:01:53,340  输出 113,340
func getTime(s string) (int64, string) {
	timeSlice := strings.Split(s, ",")
	time := timeSlice[0]
	timeS := strings.Split(time, ":")
	sec, _ := strconv.ParseInt(timeS[2], 10, 64)
	sec2, _ := strconv.ParseInt(timeS[1], 10, 64)
	sec3, _ := strconv.ParseInt(timeS[0], 10, 64)
	return sec + sec2*60 + sec3*3600, timeSlice[1]
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
