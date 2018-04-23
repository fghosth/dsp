package model

import (
	"jvole.com/dsp/config"
)

type Exads struct {
	ADX ADXbase
}

func NewExadsADX(data []byte, code uint32, domain string) *Exads {
	adx := &Exads{
		*NewADXbase(data, code, domain),
	}
	return adx
}

//获得竞价请求的字符串
func (st *Exads) GetBidderResponse(cp Compaign) (response string) {
	switch st.GetType() { //广告类型NATIVE：1,BANNER：2,POPUP：3
	case 1:
		response = st.NativeResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_NATIVE_JSON_TEMPLATE)
	case 2:
		response = st.BannerResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_BANNER_XML_TEMPLATE)
	case 3:
		response = st.PopupResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_POPUP_XML_TEMPLATE)
	}
	return
}

//GetUA()
func (st *Exads) GetUA() string {
	return st.ADX.GetUA()
}

//获得site的page
func (st *Exads) GetSitePage() string {
	return st.ADX.GetSitePage()
}

//获得site的IDzone
func (st *Exads) GetSiteIDzone() int {
	return st.ADX.GetSiteIDzone()
}

//获得language
func (st *Exads) GetLanguage() string {
	return st.ADX.GetLanguage()
}

//获得osversion
func (st *Exads) GetOSVStr() string {
	return st.ADX.GetOSVStr()
}

//获得os
func (st *Exads) GetOSStr() string {
	return st.ADX.GetOSStr()
}

//获得model
func (st *Exads) GetModel() string {
	return st.ADX.GetModel()
}

//获得广告位图片尺寸
func (st *Exads) GetImages() []string {
	return st.ADX.GetImages()
}

func (st *Exads) GetOfferInfo() OfferInfo {
	return st.ADX.GetOfferInfo()
}

func (st *Exads) GetData() []byte {
	return st.ADX.GetData()
}

func (st *Exads) SetData(data []byte) {
	st.ADX.SetData(data)
}

//设置offerinfo
func (st *Exads) SetOfferInfo() {
	st.ADX.SetOfferInfo()
	st.ADX.offerInfo.Images = st.GetImages()
}

//获得广告唯一id
func (st *Exads) GetOfferID() string {
	return st.ADX.GetOfferID()
}

//获得impid
func (st *Exads) GetIMPID() string {
	return st.ADX.GetIMPID()
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
func (st *Exads) GetType() uint32 {
	return st.ADX.GetType()
}

//获得低价 *1000000
func (st *Exads) GetBidFloor() uint32 {
	return st.ADX.GetBidFloor()
}

//获得广告位标识
func (st *Exads) GetPostion() string {

	return st.ADX.GetPostion()
}

//获得设备号
func (st *Exads) GetDeviceID() string {
	return st.ADX.GetDeviceID()
}

//获得用户标识
func (st *Exads) GetUserID() string {
	return st.ADX.GetUserID()
}

//获得country
func (st *Exads) GetCountry() string {
	return st.ADX.GetCountry()
}

//获得city
func (st *Exads) GetCity() string {
	return st.ADX.GetCity()
}

//获得region/state
func (st *Exads) GetRegion() string {
	return st.ADX.GetRegion()
}

//获得ADX
func (st *Exads) GetADX() uint32 {
	return st.ADX.GetADX()
}

//获得ADX
func (st *Exads) SetCode(code uint32) {
	st.ADX.SetCode(code)
}

//获得App或site标识
func (st *Exads) GetAppSite() string {
	return st.ADX.GetAppSite()
}

//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
func (st *Exads) GetSourceType() uint32 {

	return st.ADX.GetSourceType()
}

//获得广告类型
func (st *Exads) GetContentCat() []string {

	return st.ADX.GetContentCat()
}

//获得连接类型
func (st *Exads) GetConnType() uint32 {
	return st.ADX.GetConnType()
}

//获得运营商
func (st *Exads) GetCarriers() string {
	return st.ADX.GetCarriers()
}

//获得设备类型
func (st *Exads) GetDeviceType() uint32 {

	return st.ADX.GetDeviceType()
}

//获得os系统
func (st *Exads) GetOS() string {
	return st.ADX.GetOS()
}

//获得IdfaGaid
func (st *Exads) GetIdfaGaid() string {
	return st.ADX.GetIdfaGaid()
}

//获得IP
func (st *Exads) GetIP() string {
	return st.ADX.GetIP()
}

//获得native的response字符串 NativeResponse
func (st *Exads) NativeResponse(cp Compaign, template, admTemplate string) string {
	return st.ADX.NativeResponse(cp, template, admTemplate)
}

//获得banner的response字符串
func (st *Exads) BannerResponse(cp Compaign, template, admTemplate string) string {
	return st.ADX.BannerResponse(cp, template, admTemplate)
}

//获得Popup的response字符串
func (st *Exads) PopupResponse(cp Compaign, template, admTemplate string) string {
	return st.ADX.PopupResponse(cp, template, admTemplate)
}
