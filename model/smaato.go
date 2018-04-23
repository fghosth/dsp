package model

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/valyala/fasttemplate"
	"jvole.com/dsp/config"
	"jvole.com/dsp/util"
)

type Smaato struct {
	ADX ADXbase
}

func NewSmaatoADX(data []byte, code uint32, domain string) *Smaato {
	adx := &Smaato{
		*NewADXbase(data, code, domain),
	}

	return adx
}

//smaato图片竞价
type SMaatoIMGResponseBid struct {
	ID      string           `json:"id"`
	Seatbid []Seatbidstructs `json:"seatbid"`
}
type Seatbidstructs struct {
	Bid []Bidstruct `json:"bid"`
}
type Bidstruct struct {
	Adm     string   `json:"adm"`
	Adomain []string `json:"adomain"`
	Cid     uint32   `json:"cid"`
	Crid    string   `json:"crid"`
	ID      string   `json:"id"`
	Impid   string   `json:"impid"`
	Nurl    string   `json:"nurl"`
	Price   float64  `json:"price"`
}

//获得竞价请求的字符串
func (st *Smaato) GetBidderResponse(cp Compaign) (response string) {
	switch st.GetType() { //广告类型NATIVE：1,BANNER：2,POPUP：3
	case 1:
		response = st.NativeResponse(cp, config.SMAATO_TEMPLATE, "")
	case 2:
		response = st.BannerResponse(cp, config.SMAATO_TEMPLATE, "")
	case 3:
		response = st.PopupResponse(cp, config.SMAATO_TEMPLATE, "")
	}
	return

}

//获得广告位图片尺寸
func (st *Smaato) GetImages() []string {
	// return [][]int{[]int{320, 50}}
	return st.ADX.GetImages()
}

func (st *Smaato) GetOfferInfo() OfferInfo {
	return st.ADX.GetOfferInfo()
}

func (st *Smaato) GetData() []byte {
	return st.ADX.GetData()
}

func (st *Smaato) SetData(data []byte) {
	st.ADX.SetData(data)
	// st.ADX.data = data
}

//设置offerinfo
func (st *Smaato) SetOfferInfo() {
	st.ADX.SetOfferInfo()
}

//获得广告唯一id
func (st *Smaato) GetOfferID() string {
	return st.ADX.GetOfferID()
}

//获得impid
func (st *Smaato) GetIMPID() string {
	return st.ADX.GetIMPID()
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
func (st *Smaato) GetType() uint32 {
	return st.ADX.GetType()
}

//获得低价 *1000000
func (st *Smaato) GetBidFloor() uint32 {
	return st.ADX.GetBidFloor()
}

//获得广告位标识
func (st *Smaato) GetPostion() string {

	return st.ADX.GetPostion()
}

//获得设备号
func (st *Smaato) GetDeviceID() string {
	return st.ADX.GetDeviceID()
}

//获得用户标识
func (st *Smaato) GetUserID() string {
	return st.ADX.GetUserID()
}

//获得country
func (st *Smaato) GetCountry() string {
	return st.ADX.GetCountry()
}

//获得city
func (st *Smaato) GetCity() string {
	return st.ADX.GetCity()
}

//获得region/state
func (st *Smaato) GetRegion() string {
	return st.ADX.GetRegion()
}

//获得ADX
func (st *Smaato) GetADX() uint32 {
	return st.ADX.GetADX()
}

//获得ADX
func (st *Smaato) SetCode(code uint32) {
	st.ADX.SetCode(code)
}

//获得App或site标识
func (st *Smaato) GetAppSite() string {
	return st.ADX.GetAppSite()
}

//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
func (st *Smaato) GetSourceType() uint32 {

	return st.ADX.GetSourceType()
}

//获得广告类型
func (st *Smaato) GetContentCat() []string {

	return st.ADX.GetContentCat()
}

//获得连接类型
func (st *Smaato) GetConnType() uint32 {
	return st.ADX.GetConnType()
}

//获得运营商
func (st *Smaato) GetCarriers() string {
	return st.ADX.GetCarriers()
}

//获得设备类型
func (st *Smaato) GetDeviceType() uint32 {

	return st.ADX.GetDeviceType()
}

//获得os系统
func (st *Smaato) GetOS() string {
	return st.ADX.GetOS()
}

//获得IdfaGaid
func (st *Smaato) GetIdfaGaid() string {
	return st.ADX.GetIdfaGaid()
}

//获得IP
func (st *Smaato) GetIP() string {
	return st.ADX.GetIP()
}

//GetUA()
func (st *Smaato) GetUA() string {
	return st.ADX.GetUA()
}

//获得site的page
func (st *Smaato) GetSitePage() string {
	return st.ADX.GetSitePage()
}

//获得site的IDzone
func (st *Smaato) GetSiteIDzone() int {
	return st.ADX.GetSiteIDzone()
}

//获得language
func (st *Smaato) GetLanguage() string {
	return st.ADX.GetLanguage()
}

//获得osversion
func (st *Smaato) GetOSVStr() string {
	return st.ADX.GetOSVStr()
}

//获得os
func (st *Smaato) GetOSStr() string {
	return st.ADX.GetOSStr()
}

//获得model
func (st *Smaato) GetModel() string {
	return st.ADX.GetModel()
}

//获得banner的response字符串
func (st *Smaato) BannerResponse(cp Compaign, template, admTemplate string) string {
	var res string
	var pos int
	uuid := util.NewUUID(time.Now().UTC().String())
	t := fasttemplate.New(template, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(cp.Creatives) == 0 { // 只支持图片广告
		return ""
	}
	carr := make([]Creative, 0)
	for _, v := range cp.Creatives { //获得合法规格的图片
		if v.MainImgFormat == st.ADX.offerInfo.Images[0] {
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
	beaconSlice := [][]byte{
	// []byte(`<beacon><![CDATA[` + cp.RedirectURL + `]]><\/beacon>`),
	// []byte(`<beacon><![CDATA[` + cp.ImpressionURL + `]]><\/beacon>`),
	}
	tagfun := fasttemplate.TagFunc(func(w io.Writer, tag string) (int, error) {
		var nn int
		for _, x := range beaconSlice {
			n, err := w.Write(x)
			if err != nil {
				return nn, err
			}
			nn += n
		}
		return nn, nil
	})
	//win 连接
	nurl := st.ADX.GetNURL(config.AUCTION_PRICE, cp, uuid)
	clickurl := cp.RedirectURL
	// impressionURL := cp.ImpressionURL + "&" + config.IMPURL_ImageURLName + "=" + carr[pos].ImgUrl
	impressionURL := st.ADX.GetImpressionUrl(cp.ImpressionURL, carr[pos].ImgUrl, uuid)
	// fmt.Println(nurl)
	res = t.ExecuteString(map[string]interface{}{
		"NURL":      nurl,
		"CRID":      cp.Creatives[pos].CreativeID,
		"ADOMAIN":   cp.Adomain,
		"PRICE":     fmt.Sprint(float32(cp.Maxbid/config.CPMBidFix) / float32(config.Cashfix)),
		"ID":        st.ADX.offerInfo.ID,
		"CLICKURL":  clickurl,
		"IMGURL":    impressionURL,
		"HEIGHT":    strings.Split(st.ADX.offerInfo.Images[0], "x")[1], //TODO 修改 选出所有符合规格的图片，随机选一张
		"WIDTH":     strings.Split(st.ADX.offerInfo.Images[0], "x")[0],
		"IMPID":     st.ADX.offerInfo.IMPID_STR,
		"CID":       fmt.Sprint(cp.ID),
		"BEACONURL": tagfun,
	})

	return res
}

//获得Popup的response字符串
func (st *Smaato) PopupResponse(cp Compaign, template, admTemplate string) string {
	return ""
}

//获得native的response字符串 NativeResponse
func (st *Smaato) NativeResponse(cp Compaign, template, admTemplate string) string {
	return ""
}
