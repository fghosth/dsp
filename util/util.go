package util

import (
	"bytes"
	"compress/zlib"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	uuid "github.com/satori/go.uuid"
	"jvole.com/dsp/config"
)

var (
	KitLogger kitlog.Logger
	IPSegment = `^(?:(?:1[0-9][0-9]\.)|(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))/[123][12]?$`
	IP        = `^(?:(?:1[0-9][0-9]\.)|(?:2[0-4][0-9]\.)|(?:25[0-5]\.)|(?:[1-9][0-9]\.)|(?:[0-9]\.)){3}(?:(?:1[0-9][0-9])|(?:2[0-4][0-9])|(?:25[0-5])|(?:[1-9][0-9])|(?:[0-9]))$`
)

func init() {
	KitLogger = kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	switch config.Loglevel {
	case "all":
		KitLogger = level.NewFilter(KitLogger, level.AllowAll())
	case "debug":
		KitLogger = level.NewFilter(KitLogger, level.AllowDebug())
	case "info":
		KitLogger = level.NewFilter(KitLogger, level.AllowInfo())
	case "error":
		KitLogger = level.NewFilter(KitLogger, level.AllowError())
	}
	KitLogger = kitlog.With(KitLogger, "ts", kitlog.DefaultTimestampUTC)
}

//列出目录下所有文件
func ListFiles(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}

//处理clickurl
//http://iytg3a.nbtrk6.com/22b90d96-8adf-4978-9229-4cc643a0e12f  to  http://iytg3a.nbtrk6.com/impressionClick/22b90d96-8adf-4978-9229-4cc643a0e12f
func DealWithClickURL(url string) string {
	urlarr := strings.Split(url, "/")
	var result string
	for i := 0; i < len(urlarr); i++ {
		if result == "" {
			result = urlarr[i]
		} else {
			if i == 2 {
				result = result + "/" + urlarr[i] + "/impressionClick"
			} else {
				result = result + "/" + urlarr[i]
			}
		}

	}
	Tracefile(result+"------"+url, "url.log")
	return result
}

//生成uuid
func NewUUID(mark string) string {
	uuidV4, _ := uuid.NewV4()
	return uuid.NewV5(uuidV4, mark).String()
}

//判断一个数是否是2的N次方
func Is2N(n int) bool {
	return n > 0 && ((n & (n - 1)) == 0)
}

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//数组插入 pos 位置（下标+1） data数据，source原数组
func ArrInsertAfter(pos int, data interface{}, source []interface{}) []interface{} {
	length := len(source)
	slice1 := source[0:pos:pos]
	slice1 = append(slice1, data)
	slice2 := source[pos:length:length]
	result := append(slice1, slice2...)
	return result
}

//数组删除 pos 位置（下标+1） source原数组
func ArrDeletePos(pos int, source []interface{}) []interface{} {
	length := len(source)
	slice1 := source[0 : pos-1 : pos-1]
	slice2 := source[pos:length:length]
	result := append(slice1, slice2...)
	return result
}

//生成Guid字串
func GetGuid() string {
	b := make([]byte, 48)

	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

//进行zlib压缩
func DoZlibCompress(src []byte) []byte {
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(src)
	w.Close()
	return in.Bytes()
}

//进行zlib解压缩
func DoZlibUnCompress(compressSrc []byte) []byte {
	b := bytes.NewReader(compressSrc)
	var out bytes.Buffer
	r, _ := zlib.NewReader(b)
	io.Copy(&out, r)
	return out.Bytes()
}

//struct转byte
func EncodeStructToByte(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

//byte转struct
func DecodeByteToStruct(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

//读取文件
func ReadFile(path string) []byte {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	return fd
}

//判断文件或文件夹是否存在 error为nil 则存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// offset:+08:00/+05:45/-08:00
// return:seconds east of UTC
func ZoneOffset(offset string) int {
	if offset == "" {
		return 0
	}
	ts := strings.Split(offset[1:], ":")
	var h, m int64
	var err error
	if h, err = strconv.ParseInt(ts[0], 10, 64); err != nil {

		fmt.Printf("zoneOffset(%s) hour error:%s\n", offset, err.Error())
		return 0
	}
	if m, err = strconv.ParseInt(ts[1], 10, 64); err != nil {
		fmt.Printf("zoneOffset(%s) minute error:%s\n", offset, err.Error())
		return 0
	}
	os := (h*60 + m) * 60
	if strings.HasPrefix(offset, "-") {
		return int(os * -1)
	}
	return int(os)
}

//判断是否是ip
func IsIP(ip string) bool {
	match, _ := regexp.MatchString(IP, ip)
	return match
}

//判断是否是ip网段
func IsIPSegment(ips string) bool {
	match, _ := regexp.MatchString(IPSegment, ips)
	return match
}

//判断ip地址是否属于某一网段 ip地址，ip网段 返回bool 2个ip与子网掩码 and 相同则在同个网段 --子网掩码方便求得技巧，补码+1
func IsInsegment(ip, ipS string) (isIn bool) {
	var ipsInt, ipInt int
	arrips := strings.Split(ipS, "/")
	netmask, _ := strconv.Atoi(arrips[1])
	netmask = 0xFFFFFFFF << uint(32-netmask)
	ips := arrips[0]
	ipsarr := strings.Split(ips, ".")

	tmpint, _ := strconv.Atoi(ipsarr[0])
	tmpint2, _ := strconv.Atoi(ipsarr[1])
	tmpint3, _ := strconv.Atoi(ipsarr[2])
	tmpint4, _ := strconv.Atoi(ipsarr[3])
	ipsInt = tmpint<<24 | tmpint2<<16 | tmpint3<<8 | tmpint4

	iparr := strings.Split(ip, ".")
	tmpint, _ = strconv.Atoi(iparr[0])
	tmpint2, _ = strconv.Atoi(iparr[1])
	tmpint3, _ = strconv.Atoi(iparr[2])
	tmpint4, _ = strconv.Atoi(iparr[3])
	ipInt = tmpint<<24 | tmpint2<<16 | tmpint3<<8 | tmpint4

	// fmt.Println(ipsInt, ipInt, ipsInt-ipInt)
	if ipInt&netmask == ipsInt&netmask {
		isIn = true
	}
	return
}

func Int64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func BytesToInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

//处理连接 用于竞价成功后返回的xml 『/』替换为【\/】
func DealXMLURL(url string) string {
	url = strings.Replace(url, "/", "\\/", -1)
	return url
}

//追踪内容文件，追加模式添加
func Tracefile(str_content string, path string) {
	fd, _ := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	fd_time := time.Now().Format("2006-01-02 15:04:05")
	fd_content := strings.Join([]string{"======", fd_time, "=====", str_content, "\n"}, "")
	buf := []byte(fd_content)
	fd.Write(buf)
	fd.Close()
}

//用户http reponse zip内容
func GetResponseZip(w io.Writer, data []byte) error {
	zipW := zlib.NewWriter(w)
	defer zipW.Close()
	_, err := zipW.Write(data)
	if err != nil {
		return err
	}
	return nil
}
