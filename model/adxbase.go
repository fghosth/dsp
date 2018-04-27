package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/buger/jsonparser"
	"github.com/k0kubun/pp"
	"github.com/valyala/fasttemplate"
	"jvole.com/dsp/config"
	"jvole.com/dsp/util"
)

type ADXbase struct {
	code      uint32    //编号 2的N次方
	data      []byte    //json字符串
	offerInfo OfferInfo //对比数据集合
	num       string    //随机广告位
}

func NewADXbase(data []byte, code uint32, domain string) *ADXbase {
	adx := &ADXbase{
		code,
		data,
		*&OfferInfo{},
		"0",
	}
	adx.SetOfferInfo()
	adx.SetDomain(domain)
	return adx

}

type OfferInfo struct { //优化数据后的集合
	ID             string
	Type           uint32
	Bidfloor       uint32
	Postion        string
	DeviceID       uint32
	DeviceID_STR   string
	UserID         uint32
	UserID_STR     string
	Country        uint32
	Country_STR    string
	City           uint32
	City_STR       string
	Region         uint32
	Region_STR     string
	SourceType     uint32
	ContentCat     []uint32
	ContentCat_STR []string
	ConnType       uint32
	Carriers       uint32
	Carriers_STR   string
	DeviceType     uint32
	OS             uint32
	OS_STR         string
	IdfaGaid       uint32
	IdfaGaid_STR   string
	IP             uint32
	IP_STR         string
	APPSite        uint32
	APPSite_STR    string
	Images         []string
	IMPID_STR      string
	SitePage_STR   string
	SiteIDZone_STR string
	Language_STR   string
	OSstr_STR      string
	OSV_STR        string
	UA_STR         string
	Model_STR      string
	OfferURL       string
	Domain         uint32 //访问的域名,用来筛选成人流量
	Domain_STR     string //访问的域名,用来筛选成人流量

}

//设置Domain
func (st *ADXbase) SetDomain(domain string) {
	st.offerInfo.Domain_STR = domain
	st.offerInfo.Domain = util.Hashcode(domain)
}

//根据条件传入不同模板
func (st *ADXbase) GetBidderResponse(cp Compaign) (response string) {

	return
}

func (st *ADXbase) SetCode(code uint32) {
	st.code = code
}

func (st *ADXbase) GetOfferInfo() OfferInfo {
	return st.offerInfo
}

func (st *ADXbase) GetData() []byte {
	return st.data
}

func (st *ADXbase) SetData(data []byte) {
	st.data = data
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int

	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		// a, _, _, _ := jsonparser.Get(value, "uuid")
		rnd := r.Intn(2)

		if rnd == 1 {
			st.num = fmt.Sprint(i)
			return
		}
		i++
		// fmt.Println(string(a))
	}, "imp")

}

//设置offerinfo
func (st *ADXbase) SetOfferInfo() {
	st.offerInfo.ID = st.GetOfferID()
	st.offerInfo.Type = st.GetType()
	st.offerInfo.Bidfloor = st.GetBidFloor()
	st.offerInfo.Postion = st.GetPostion()
	st.offerInfo.DeviceID = util.Hashcode(st.GetDeviceID())
	st.offerInfo.DeviceID_STR = st.GetDeviceID()
	st.offerInfo.UserID = util.Hashcode(st.GetUserID())
	st.offerInfo.UserID_STR = st.GetUserID()
	st.offerInfo.Country = util.Hashcode(st.GetCountry())
	st.offerInfo.Country_STR = st.GetCountry()
	st.offerInfo.City = util.Hashcode(st.GetCity())
	st.offerInfo.City_STR = st.GetCity()
	st.offerInfo.Region = util.Hashcode(st.GetRegion())
	st.offerInfo.Region_STR = st.GetRegion()
	st.offerInfo.SourceType = st.GetSourceType()
	ccc := st.GetContentCat()
	st.offerInfo.ContentCat = make([]uint32, len(ccc))
	st.offerInfo.ContentCat_STR = make([]string, len(ccc))
	for k, v := range ccc {
		st.offerInfo.ContentCat[k] = util.Hashcode(v)
		st.offerInfo.ContentCat_STR[k] = v
	}
	st.offerInfo.ConnType = st.GetConnType()
	st.offerInfo.Carriers = util.Hashcode(st.GetCarriers())
	st.offerInfo.Carriers_STR = st.GetCarriers()
	st.offerInfo.DeviceType = st.GetDeviceType()
	st.offerInfo.OS = util.Hashcode(st.GetOS())
	st.offerInfo.OS_STR = st.GetOS()
	st.offerInfo.IdfaGaid = util.Hashcode(st.GetIdfaGaid())
	st.offerInfo.IdfaGaid_STR = st.GetIdfaGaid()
	st.offerInfo.IP = util.Hashcode(st.GetIP())
	st.offerInfo.IP_STR = st.GetIP()
	st.offerInfo.APPSite = util.Hashcode(st.GetAppSite())
	st.offerInfo.APPSite_STR = st.GetAppSite()
	st.offerInfo.Images = st.GetImages()
	st.offerInfo.IMPID_STR = st.GetIMPID()
	st.offerInfo.UA_STR = st.GetUA()
	st.offerInfo.Language_STR = st.GetLanguage()
	st.offerInfo.SiteIDZone_STR = strconv.Itoa(st.GetSiteIDzone())
	st.offerInfo.SitePage_STR = st.GetSitePage()
	st.offerInfo.OSV_STR = st.GetOSVStr()
	st.offerInfo.OSstr_STR = st.GetOSStr()
	st.offerInfo.Model_STR = st.GetModel()

	imUrl := config.IMPURL_SiteORAPPID + "=" + st.offerInfo.Postion
	imUrl = imUrl + "&" + config.IMPURL_Domain + "=" + st.offerInfo.APPSite_STR
	imUrl = imUrl + "&" + config.IMPURL_Cat + "=" + strings.Join(st.offerInfo.ContentCat_STR, ",")
	imUrl = imUrl + "&" + config.IMPURL_Page + "=" + st.offerInfo.SitePage_STR
	imUrl = imUrl + "&" + config.IMPURL_IDZone + "=" + st.offerInfo.SiteIDZone_STR
	imUrl = imUrl + "&" + config.IMPURL_IP + "=" + st.offerInfo.IP_STR
	imUrl = imUrl + "&" + config.IMPURL_Country + "=" + st.offerInfo.Country_STR
	imUrl = imUrl + "&" + config.IMPURL_City + "=" + st.offerInfo.City_STR
	imUrl = imUrl + "&" + config.IMPURL_Region + "=" + st.offerInfo.Region_STR
	imUrl = imUrl + "&" + config.IMPURL_OS + "=" + st.offerInfo.OSstr_STR
	imUrl = imUrl + "&" + config.IMPURL_Language + "=" + st.offerInfo.Language_STR
	imUrl = imUrl + "&" + config.IMPURL_Model + "=" + st.offerInfo.Model_STR
	imUrl = imUrl + "&" + config.IMPURL_OSV + "=" + st.offerInfo.OSV_STR
	// imUrl = imUrl + "&" + config.IMPURL_UA + "=" + st.offerInfo.UA_STR
	imUrl = imUrl + "&" + config.IMPURL_DeviceType + "=" + fmt.Sprint(st.offerInfo.DeviceType)
	imUrl = imUrl + "&" + config.IMPURL_DeviceID + "=" + st.offerInfo.DeviceID_STR
	imUrl = imUrl + "&" + config.IMPURL_ADX + "=" + fmt.Sprint(config.ADXCode[config.ADX])
	st.offerInfo.OfferURL = imUrl
}

// //获得native 广告的title长度
// func (st *ADXbase) GetTitleLen() int {
//
// }
//
// //获得native 广告的Description信息
// func (st *ADXbase) GetDescription() []string {
//
// }

//获得广告位图片尺寸
func (st *ADXbase) GetImages() []string {
	images := make([]string, 0)
	switch st.offerInfo.Type { //广告类型NATIVE：1,BANNER：2,POPUP：3
	case 1:
		res, _ := jsonparser.GetString(st.data, "imp", "["+st.num+"]", "native", "request")
		res = strings.Replace(res, "\\", "", -1)
		w, _ := jsonparser.GetInt([]byte(res), "native", "assets", "[0]", "img", "wmin")
		h, _ := jsonparser.GetInt([]byte(res), "native", "assets", "[0]", "img", "hmin")
		images = append(images, strconv.Itoa(int(w))+"x"+strconv.Itoa(int(h)))
	case 2:
		w, _ := jsonparser.GetInt(st.data, "imp", "["+st.num+"]", "banner", "w")
		h, _ := jsonparser.GetInt(st.data, "imp", "["+st.num+"]", "banner", "h")
		images = append(images, strconv.Itoa(int(w))+"x"+strconv.Itoa(int(h)))
	case 3:
		jsonparser.ArrayEach(st.data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			w, _ := jsonparser.GetInt(value, "w")
			h, _ := jsonparser.GetInt(value, "h")
			images = append(images, strconv.Itoa(int(w))+"x"+strconv.Itoa(int(h)))
		}, "imp", "["+st.num+"]", "popup", "format")
	default:
		images = nil
	}

	return images
}

//获得广告唯一id
func (st *ADXbase) GetOfferID() string {
	res, _, _, _ := jsonparser.Get(st.data, "id")
	return string(res)
}

//获得impid
func (st *ADXbase) GetIMPID() string {
	res, _, _, _ := jsonparser.Get(st.data, "imp", "["+st.num+"]", "id")
	return string(res)
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
func (st *ADXbase) GetType() uint32 {
	var otype uint32
	// arrType := []string{"banner", "native", "popup"}
	res, _, _, _ := jsonparser.Get(st.data, "imp", "["+st.num+"]", "banner", "w")
	if string(res) != "" {
		otype = 2
	}
	res, _, _, _ = jsonparser.Get(st.data, "imp", "["+st.num+"]", "native", "request")
	if string(res) != "" {
		otype = 1
	}

	if otype == 0 {
		otype = 3
	}

	return otype
}

//获得低价 *1000000
func (st *ADXbase) GetBidFloor() uint32 {
	res, _, _, _ := jsonparser.Get(st.data, "imp", "["+st.num+"]", "bidfloor")
	tmpbit, _ := strconv.ParseFloat(string(res), 5)
	bits := uint32(tmpbit * float64(config.Cashfix))
	return bits
}

//获得广告位标识
func (st *ADXbase) GetPostion() string {
	res, _, _, _ := jsonparser.Get(st.data, "app", "id")
	if len(res) > 0 {
		return string(res)
	}
	res, _, _, _ = jsonparser.Get(st.data, "site", "id")
	return string(res)
}

//获得设备号
func (st *ADXbase) GetDeviceID() string {
	res, _, _, _ := jsonparser.Get(st.data, "device", "ifa")
	return string(res)
}

//获得用户标识
func (st *ADXbase) GetUserID() string {
	res, _, _, _ := jsonparser.Get(st.data, "user", "id")
	return string(res)
}

//获得country
func (st *ADXbase) GetCountry() string {
	res, _, _, _ := jsonparser.Get(st.data, "device", "geo", "country")
	return string(res)
}

//获得city
func (st *ADXbase) GetCity() string {
	res, _, _, _ := jsonparser.Get(st.data, "device", "geo", "city")
	return string(res)
}

//获得region/state
func (st *ADXbase) GetRegion() string {
	res, _, _, _ := jsonparser.Get(st.data, "device", "geo", "region")
	return string(res)
}

//获得ADX
func (st *ADXbase) GetADX() uint32 {
	return st.code
}

//获得App或site domain
func (st *ADXbase) GetAppSite() string {
	res, _, _, _ := jsonparser.Get(st.data, "app", "domain")
	if len(res) > 0 {
		return string(res)
	}
	res, _, _, _ = jsonparser.Get(st.data, "site", "domain")
	return string(res)
}

//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
func (st *ADXbase) GetSourceType() uint32 {
	var otype uint32
	arrType := []string{"app", "site"}
	for _, v := range arrType {
		res, _, _, _ := jsonparser.Get(st.data, v)
		if string(res) != "" {
			switch v {
			case "app":
				otype = 2
			case "site":
				otype = 1
			}

		}
	}
	return otype
}

//获得广告类型
func (st *ADXbase) GetContentCat() []string {
	var result []string
	res, _, _, _ := jsonparser.Get(st.data, "app", "cat")
	if len(res) > 0 {
		json.Unmarshal(res, &result)
		return result
	}
	res, _, _, _ = jsonparser.Get(st.data, "site", "cat")
	json.Unmarshal(res, &result)
	return result
}

//获得site的page
func (st *ADXbase) GetSitePage() string {
	var res string
	res, _ = jsonparser.GetString(st.data, "site", "page")
	return res
}

//获得site的IDzone
func (st *ADXbase) GetSiteIDzone() int {
	var res int64
	res, _ = jsonparser.GetInt(st.data, "site", "ext", "idzone")
	return int(res)
}

//获得device/ua
func (st *ADXbase) GetUA() string {
	var res string
	res, _ = jsonparser.GetString(st.data, "device", "ua")
	return res
}

//获得language
func (st *ADXbase) GetLanguage() string {
	var res string
	res, _ = jsonparser.GetString(st.data, "device", "language")
	return res
}

//获得osversion
func (st *ADXbase) GetOSVStr() string {
	var res string
	res, _ = jsonparser.GetString(st.data, "device", "osv")
	return res
}

//获得os
func (st *ADXbase) GetOSStr() string {
	var res string
	res, _ = jsonparser.GetString(st.data, "device", "os")
	return res
}

//获得model
func (st *ADXbase) GetModel() string {
	var res string
	res, _ = jsonparser.GetString(st.data, "device", "model")
	return res
}

//获得连接类型
func (st *ADXbase) GetConnType() uint32 {
	res, _, _, _ := jsonparser.Get(st.data, "device", "connectiontype")
	tmpbit, _ := strconv.ParseUint(string(res), 10, 32)
	bits := uint32(tmpbit)
	return bits
}

//获得运营商
func (st *ADXbase) GetCarriers() string {
	res, _, _, _ := jsonparser.Get(st.data, "ext", "carriername")
	return string(res)
}

//获得设备类型
func (st *ADXbase) GetDeviceType() uint32 {
	res, _, _, _ := jsonparser.Get(st.data, "device", "devicetype")
	tmpbit, _ := strconv.ParseUint(string(res), 10, 32)
	bits := uint32(tmpbit)
	return bits
}

//获得os系统
func (st *ADXbase) GetOS() string {
	os, _, _, _ := jsonparser.Get(st.data, "device", "os")
	version, _, _, _ := jsonparser.Get(st.data, "device", "osv")
	if string(version) != "" {
		return strings.ToLower(string(os)) + " " + strings.ToLower(string(version))
	} else {
		return strings.ToLower(string(os))
	}
}

//获得IdfaGaid
func (st *ADXbase) GetIdfaGaid() string {
	res, _, _, _ := jsonparser.Get(st.data, "device", "ifa")
	return string(res)
}

//获得IP
func (st *ADXbase) GetIP() string {
	res, _, _, _ := jsonparser.Get(st.data, "device", "ip")
	return string(res)
}

//获得连接Nurl
func (st *ADXbase) GetNURL(price string, cp Compaign, uuid string) string {
	//TODO 时间待验证 应该是utc时间
	nurl := config.NURL + "?" + config.NURLParam["OID"] + "=" + st.offerInfo.ID + "&" + config.NURLParam["CID"] + "=" + fmt.Sprint(cp.ID) + "&" + config.NURLParam["UID"] + "=" + fmt.Sprint(cp.UID) + "&" + config.NURLParam["Postion"] + "=" + st.offerInfo.Postion + "&" + config.NURLParam["T"] + "=" + strconv.Itoa(int(time.Now().Unix())) + "&" + config.NURLParam["User"] + "=" + st.offerInfo.UserID_STR + "&" + config.NURLParam["Device"] + "=" + st.offerInfo.DeviceID_STR + "&" + config.NURLParam["Price"] + "=" + price //"${AUCTION_PRICE}"
	nurl = nurl + "&" + config.IMPURL_ClickID + "=" + uuid
	// u, _ := url.Parse(nurl)
	// q := u.Query()
	//
	// u.RawQuery = q.Encode() //urlencode
	// nurl = u.String()
	return nurl
}

//处理impressionurl
func (st *ADXbase) GetImpressionUrl(urls, imgurl string, uuid string) (imUrl string) {
	imUrl = urls + "?" + config.IMPURL_ImageURLName + "=" + imgurl
	imUrl = imUrl + "&" + config.IMPURL_ClickID + "=" + uuid
	imUrl = imUrl + "&" + st.offerInfo.OfferURL
	u, err := url.Parse(imUrl)
	if err != nil {
		pp.Println(err)
	}
	q := u.Query()
	q.Add(config.IMPURL_UA, st.offerInfo.UA_STR) //ua直接拼接 encode会出错
	u.RawQuery = q.Encode()                      //urlencode
	imUrl = u.String()
	// pp.Println("====", imUrl)
	return
}

//获得native的response字符串 NativeResponse
func (st *ADXbase) NativeResponse(cp Compaign, template, admTemplate string) string {
	var res string
	var pos int
	uuid := util.NewUUID(time.Now().UTC().String())
	admt := fasttemplate.New(admTemplate, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	t := fasttemplate.New(template, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(cp.Creatives) == 0 { // 只支持图片广告
		return ""
	}
	carr := make([]Creative, 0)
	for _, v := range cp.Creatives { //获得合法规格的图片

		imgsize := strings.Split(st.offerInfo.Images[0], "x")
		w, _ := strconv.Atoi(imgsize[0])
		h, _ := strconv.Atoi(imgsize[1])
		if v.Width >= w && v.Height >= h {
			carr = append(carr, v)
		}
	}
	if len(carr) > 1 {
		pos = r.Intn(len(carr) - 1) //随机数
	} else if len(carr) == 0 {
		return ""
	} else {
		pos = 0
	}
	//win 连接
	nurl := st.GetNURL(config.AUCTION_PRICE, cp, uuid)
	clickurl := cp.RedirectURL
	clickurl = st.GetImpressionUrl(clickurl, "", uuid)
	// impressionURL := cp.ImpressionURL + "&" + config.IMPURL_ImageURLName + "=" + carr[pos].ImgUrl
	impressionURL := st.GetImpressionUrl(cp.ImpressionURL, carr[pos].ImgUrl, uuid)
	adid := carr[pos].CreativeID
	title := carr[pos].Headline
	description := carr[pos].Description
	admstr := admt.ExecuteString(map[string]interface{}{
		"CLICKURL":    util.DealXMLURL(clickurl),
		"IMGURL":      util.DealXMLURL(impressionURL),
		"HEIGHT":      strings.Split(st.offerInfo.Images[0], "x")[1], //TODO 修改 选出所有符合规格的图片，随机选一张
		"WIDTH":       strings.Split(st.offerInfo.Images[0], "x")[0],
		"TITLE":       title,
		"DESCRIPTION": description,
	})
	// fmt.Println(nurl)
	res = t.ExecuteString(map[string]interface{}{
		"CLICKURL": clickurl,
		"NURL":     nurl,
		"ADM":      admstr,
		"PRICE":    fmt.Sprint(float32(cp.Maxbid/config.CPMBidFix) / float32(config.Cashfix)),
		"ID":       st.offerInfo.ID,
		"IMPID":    st.offerInfo.IMPID_STR,
		"CID":      fmt.Sprint(cp.ID),
		"ADID":     adid,
	})
	return res
}

//获得banner的response字符串
func (st *ADXbase) BannerResponse(cp Compaign, template, admTemplate string) string {
	var res string
	var pos int
	uuid := util.NewUUID(time.Now().UTC().String())
	admt := fasttemplate.New(admTemplate, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	t := fasttemplate.New(template, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(cp.Creatives) == 0 { // 只支持图片广告
		return ""
	}
	carr := make([]Creative, 0)
	for _, v := range cp.Creatives { //获得合法规格的图片
		if v.MainImgFormat == st.offerInfo.Images[0] {
			carr = append(carr, v)
		}
	}

	if len(carr) > 1 {
		pos = r.Intn(len(carr) - 1) //随机数
	} else if len(carr) == 0 {
		return ""
	} else {
		pos = 0
	}

	//win 连接
	nurl := st.GetNURL(config.AUCTION_PRICE, cp, uuid)
	clickurl := cp.RedirectURL
	clickurl = st.GetImpressionUrl(clickurl, "", uuid)
	// impressionURL := cp.ImpressionURL + "&" + config.IMPURL_ImageURLName + "=" + carr[pos].ImgUrl
	impressionURL := st.GetImpressionUrl(cp.ImpressionURL, carr[pos].ImgUrl, uuid)
	adid := carr[pos].CreativeID
	admstr := admt.ExecuteString(map[string]interface{}{
		"CLICKURL": clickurl,
		"IMGURL":   impressionURL,
	})
	res = t.ExecuteString(map[string]interface{}{
		"NURL":  nurl,
		"ADM":   admstr,
		"PRICE": fmt.Sprint(float32(cp.Maxbid/config.CPMBidFix) / float32(config.Cashfix)),
		"ID":    st.offerInfo.ID,
		"IMPID": st.offerInfo.IMPID_STR,
		"CID":   fmt.Sprint(cp.ID),
		"ADID":  adid,
	})

	return res
}

//获得Popup的response字符串
func (st *ADXbase) PopupResponse(cp Compaign, template, admTemplate string) string {
	var res string
	uuid := util.NewUUID(time.Now().UTC().String())
	admt := fasttemplate.New(admTemplate, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	t := fasttemplate.New(template, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)

	//win 连接
	nurl := st.GetNURL(config.AUCTION_PRICE, cp, uuid)

	clickurl := cp.RedirectURL
	clickurl = st.GetImpressionUrl(clickurl, "", uuid)
	adid := ""
	admstr := admt.ExecuteString(map[string]interface{}{
		"CLICKURL": clickurl,
	})

	res = t.ExecuteString(map[string]interface{}{
		"CLICKURL": clickurl,
		"NURL":     nurl,
		"ADM":      admstr,
		"PRICE":    fmt.Sprint(float32(cp.Maxbid/config.CPMBidFix) / float32(config.Cashfix)),
		"ID":       st.offerInfo.ID,
		"IMPID":    st.offerInfo.IMPID_STR,
		"CID":      fmt.Sprint(cp.ID),
		"ADID":     adid,
	})

	return res
}
