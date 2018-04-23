package model

import (
	"jvole.com/dsp/config"
)

type MGID struct {
	ADX ADXbase
}

func NewMGIDADX(data []byte, code uint32, domain string) *MGID {
	adx := &MGID{
		*NewADXbase(data, code, domain),
	}
	return adx
}

//获得竞价请求的字符串
func (st *MGID) GetBidderResponse(cp Compaign) (response string) {
	switch st.GetType() { //广告类型NATIVE：1,BANNER：2,POPUP：3
	case 1:
		response = st.NativeResponse(cp, config.BASE_BID_TEMPLATE, config.MGID_NATIVE_JSON_TEMPLATE)
	case 2:
		// response = st.BannerResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_BANNER_XML_TEMPLATE)
	case 3:
		// response = st.PopupResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_POPUP_XML_TEMPLATE)
	}
	return
}

//GetUA()
func (st *MGID) GetUA() string {
	return st.ADX.GetUA()
}

//获得site的page
func (st *MGID) GetSitePage() string {
	return st.ADX.GetSitePage()
}

//获得site的IDzone
func (st *MGID) GetSiteIDzone() int {
	return st.ADX.GetSiteIDzone()
}

//获得language
func (st *MGID) GetLanguage() string {
	return st.ADX.GetLanguage()
}

//获得osversion
func (st *MGID) GetOSVStr() string {
	return st.ADX.GetOSVStr()
}

//获得os
func (st *MGID) GetOSStr() string {
	return st.ADX.GetOSStr()
}

//获得model
func (st *MGID) GetModel() string {
	return st.ADX.GetModel()
}

//获得广告位图片尺寸
func (st *MGID) GetImages() []string {
	return st.ADX.GetImages()
}

func (st *MGID) GetOfferInfo() OfferInfo {
	return st.ADX.GetOfferInfo()
}

func (st *MGID) GetData() []byte {
	return st.ADX.GetData()
}

func (st *MGID) SetData(data []byte) {
	st.ADX.data = data
}

//设置offerinfo
func (st *MGID) SetOfferInfo() {
	st.ADX.SetOfferInfo()
	st.ADX.offerInfo.Images = st.GetImages()
}

//获得广告唯一id
func (st *MGID) GetOfferID() string {
	return st.ADX.GetOfferID()
}

//获得impid
func (st *MGID) GetIMPID() string {
	return st.ADX.GetIMPID()
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
func (st *MGID) GetType() uint32 {
	return st.ADX.GetType()
}

//获得低价 *1000000
func (st *MGID) GetBidFloor() uint32 {
	return st.ADX.GetBidFloor()
}

//获得广告位标识
func (st *MGID) GetPostion() string {

	return st.ADX.GetPostion()
}

//获得设备号
func (st *MGID) GetDeviceID() string {
	return st.ADX.GetDeviceID()
}

//获得用户标识
func (st *MGID) GetUserID() string {
	return st.ADX.GetUserID()
}

//获得country
func (st *MGID) GetCountry() string {
	return st.ADX.GetCountry()
}

//获得city
func (st *MGID) GetCity() string {
	return st.ADX.GetCity()
}

//获得region/state
func (st *MGID) GetRegion() string {
	return st.ADX.GetRegion()
}

//获得ADX
func (st *MGID) GetADX() uint32 {
	return st.ADX.GetADX()
}

//获得ADX
func (st *MGID) SetCode(code uint32) {
	st.ADX.SetCode(code)
}

//获得App或site标识
func (st *MGID) GetAppSite() string {
	return st.ADX.GetAppSite()
}

//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
func (st *MGID) GetSourceType() uint32 {

	return st.ADX.GetSourceType()
}

//获得广告类型
func (st *MGID) GetContentCat() []string {

	return st.ADX.GetContentCat()
}

//获得连接类型
func (st *MGID) GetConnType() uint32 {
	return st.ADX.GetConnType()
}

//获得运营商
func (st *MGID) GetCarriers() string {
	return st.ADX.GetCarriers()
}

//获得设备类型
func (st *MGID) GetDeviceType() uint32 {

	return st.ADX.GetDeviceType()
}

//获得os系统
func (st *MGID) GetOS() string {
	return st.ADX.GetOS()
}

//获得IdfaGaid
func (st *MGID) GetIdfaGaid() string {
	return st.ADX.GetIdfaGaid()
}

//获得IP
func (st *MGID) GetIP() string {
	return st.ADX.GetIP()
}

//获得native的response字符串 NativeResponse
func (st *MGID) NativeResponse(cp Compaign, template, admTemplate string) string {
	// var res string
	// var pos int
	// admt := fasttemplate.New(admTemplate, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	// t := fasttemplate.New(template, config.TEMPLATE_LEFT_TAG, config.TEMPLATE_RIGHT_TAG)
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// if len(cp.Creatives) == 0 { // 只支持图片广告
	// 	return ""
	// }
	// carr := make([]Creative, 0)
	// for _, v := range cp.Creatives { //获得合法规格的图片
	//
	// 	imgsize := strings.Split(st.ADX.offerInfo.Images[0], "x")
	// 	w, _ := strconv.Atoi(imgsize[0])
	// 	h, _ := strconv.Atoi(imgsize[1])
	// 	if v.Width >= w && v.Height >= h {
	// 		carr = append(carr, v)
	// 	}
	// }
	// if len(carr) > 1 {
	// 	pos = r.Intn(len(carr) - 1) //随机数
	// } else if len(carr) == 0 {
	// 	return ""
	// } else {
	// 	pos = 0
	// }
	// //win 连接
	// nurl := st.ADX.GetNURL(config.AUCTION_PRICE, cp)
	// clickurl := cp.RedirectURL
	// impressionURL := cp.ImpressionURL + "&" + config.ImageURLName + "=" + carr[pos].ImgUrl
	// title := carr[pos].Headline
	// description := carr[pos].Description
	// admstr := admt.ExecuteString(map[string]interface{}{
	// 	"CLICKURL":    clickurl,
	// 	"IMGURL":      impressionURL,
	// 	"HEIGHT":      strings.Split(st.ADX.offerInfo.Images[0], "x")[1], //TODO 修改 选出所有符合规格的图片，随机选一张
	// 	"WIDTH":       strings.Split(st.ADX.offerInfo.Images[0], "x")[0],
	// 	"TITLE":       title,
	// 	"DESCRIPTION": description,
	// })
	// // fmt.Println(nurl)
	// res = t.ExecuteString(map[string]interface{}{
	// 	"NURL":  nurl,
	// 	"ADM":   admstr,
	// 	"PRICE": fmt.Sprint(float32(cp.Maxbid/config.CPMBidFix) / float32(config.Cashfix)),
	// 	"ID":    st.ADX.offerInfo.ID,
	// 	"IMPID": st.ADX.offerInfo.IMPID_STR,
	// })
	// return res
	return st.ADX.NativeResponse(cp, template, admTemplate)
}

//获得banner的response字符串
func (st *MGID) BannerResponse(cp Compaign, template, admTemplate string) string {

	return ""
}

//获得Popup的response字符串
func (st *MGID) PopupResponse(cp Compaign, template, admTemplate string) string {

	return ""
}
