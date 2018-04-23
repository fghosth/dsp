package model

import (
	"jvole.com/dsp/config"
)

type Popundertotal struct {
	ADX ADXbase
}

func NewPopundertotalADX(data []byte, code uint32, domain string) *Popundertotal {
	adx := &Popundertotal{
		*NewADXbase(data, code, domain),
	}
	return adx
}

//获得竞价请求的字符串
func (st *Popundertotal) GetBidderResponse(cp Compaign) (response string) {
	switch st.GetType() { //广告类型NATIVE：1,BANNER：2,POPUP：3
	case 1:
		response = st.NativeResponse(cp, config.POPUNDERTOTAL_POPUP_XML_TEMPLATE, "")
	case 2:
		// response = st.BannerResponse(cp, config.BASE_BID_TEMPLATE, config.EXADS_BANNER_XML_TEMPLATE)
	case 3:
		response = st.PopupResponse(cp, config.POPUNDERTOTAL_POPUP_XML_TEMPLATE, "")
	}
	return
}

//GetUA()
func (st *Popundertotal) GetUA() string {
	return st.ADX.GetUA()
}

//获得site的page
func (st *Popundertotal) GetSitePage() string {
	return st.ADX.GetSitePage()
}

//获得site的IDzone
func (st *Popundertotal) GetSiteIDzone() int {
	return st.ADX.GetSiteIDzone()
}

//获得language
func (st *Popundertotal) GetLanguage() string {
	return st.ADX.GetLanguage()
}

//获得osversion
func (st *Popundertotal) GetOSVStr() string {
	return st.ADX.GetOSVStr()
}

//获得os
func (st *Popundertotal) GetOSStr() string {
	return st.ADX.GetOSStr()
}

//获得model
func (st *Popundertotal) GetModel() string {
	return st.ADX.GetModel()
}

//获得广告位图片尺寸
func (st *Popundertotal) GetImages() []string {
	return st.ADX.GetImages()
}

func (st *Popundertotal) GetOfferInfo() OfferInfo {
	return st.ADX.GetOfferInfo()
}

func (st *Popundertotal) GetData() []byte {
	return st.ADX.GetData()
}

func (st *Popundertotal) SetData(data []byte) {
	st.ADX.data = data
}

//设置offerinfo
func (st *Popundertotal) SetOfferInfo() {
	st.ADX.SetOfferInfo()
	st.ADX.offerInfo.Images = st.GetImages()
}

//获得广告唯一id
func (st *Popundertotal) GetOfferID() string {
	return st.ADX.GetOfferID()
}

//获得impid
func (st *Popundertotal) GetIMPID() string {
	return st.ADX.GetIMPID()
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
func (st *Popundertotal) GetType() uint32 {
	return st.ADX.GetType()
}

//获得低价 *1000000
func (st *Popundertotal) GetBidFloor() uint32 {
	return st.ADX.GetBidFloor()
}

//获得广告位标识
func (st *Popundertotal) GetPostion() string {

	return st.ADX.GetPostion()
}

//获得设备号
func (st *Popundertotal) GetDeviceID() string {
	return st.ADX.GetDeviceID()
}

//获得用户标识
func (st *Popundertotal) GetUserID() string {
	return st.ADX.GetUserID()
}

//获得country
func (st *Popundertotal) GetCountry() string {
	return st.ADX.GetCountry()
}

//获得city
func (st *Popundertotal) GetCity() string {
	return st.ADX.GetCity()
}

//获得region/state
func (st *Popundertotal) GetRegion() string {
	return st.ADX.GetRegion()
}

//获得ADX
func (st *Popundertotal) GetADX() uint32 {
	return st.ADX.GetADX()
}

//获得ADX
func (st *Popundertotal) SetCode(code uint32) {
	st.ADX.SetCode(code)
}

//获得App或site标识
func (st *Popundertotal) GetAppSite() string {
	return st.ADX.GetAppSite()
}

//获得sourceType //投放app还是网站 网站：1 app：2 全部：0
func (st *Popundertotal) GetSourceType() uint32 {

	return st.ADX.GetSourceType()
}

//获得广告类型
func (st *Popundertotal) GetContentCat() []string {

	return st.ADX.GetContentCat()
}

//获得连接类型
func (st *Popundertotal) GetConnType() uint32 {
	return st.ADX.GetConnType()
}

//获得运营商
func (st *Popundertotal) GetCarriers() string {
	return st.ADX.GetCarriers()
}

//获得设备类型
func (st *Popundertotal) GetDeviceType() uint32 {

	return st.ADX.GetDeviceType()
}

//获得os系统
func (st *Popundertotal) GetOS() string {
	return st.ADX.GetOS()
}

//获得IdfaGaid
func (st *Popundertotal) GetIdfaGaid() string {
	return st.ADX.GetIdfaGaid()
}

//获得IP
func (st *Popundertotal) GetIP() string {
	return st.ADX.GetIP()
}

//获得native的response字符串 NativeResponse
func (st *Popundertotal) NativeResponse(cp Compaign, template, admTemplate string) string {
	return st.ADX.NativeResponse(cp, template, admTemplate)
}

//获得banner的response字符串
func (st *Popundertotal) BannerResponse(cp Compaign, template, admTemplate string) string {

	return ""
}

//获得Popup的response字符串
func (st *Popundertotal) PopupResponse(cp Compaign, template, admTemplate string) string {

	return st.ADX.PopupResponse(cp, template, admTemplate)
}
