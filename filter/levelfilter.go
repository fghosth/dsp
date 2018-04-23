package filter

import (
	"github.com/RoaringBitmap/roaring"
	"github.com/go-kit/kit/log"

	"jvole.com/dsp/config"
	"jvole.com/dsp/index"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

var logger *log.Logger

func init() {
	logger = &util.KitLogger
	log.With(*logger, "component", "filter")
}

type LevelFilter interface { //过滤器接口 大范围过滤
	/*
		过滤器
		@param index.IndexMap
		@param offer model.Offer adx的offer
		@return Compaignid 的集合  过滤后剩下的compaign
	*/
	Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap
}

//执行所有levelfilter 返回过滤后的compid集合
func LevelFilterComp(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
	tf := &TypeFilter{1}
	sf := &SourceFilter{2}
	// ctf := &ConnTypeFilter{4}
	// dtf := &DeviceTypeFilter{8}
	igf := &IdfaGaidFilter{16}
	levelf := []LevelFilter{igf, tf, sf}
	andf := &AndLevelFilter{}
	andf.AndFilter(levelf...)

	return andf.Filter(compaign, offer)
}

//广告类型NATIVE：1,BANNER：2,POPUP：3
type TypeFilter struct { //编号1
	Code uint32 //过滤器编号1,2,4
}

func (tf *TypeFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
	offerinfo := offer.GetOfferInfo()
	bitm := roaring.NewBitmap()
	bitm.Or(compaign.Compaign)
	switch offerinfo.Type {
	case 1:
		bitm.And(compaign.TypeNative)
	case 2:
		bitm.And(compaign.TypeBanner)
	case 3:
		bitm.And(compaign.TypePopup)
	}
	return bitm
}

//投放app还是网站 网站：1 app：2 全部：0
type SourceFilter struct { //编号2
	Code uint32 //过滤器编号1,2,4
}

func (sf *SourceFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
	offerinfo := offer.GetOfferInfo()
	bitm := roaring.NewBitmap()
	bitm.Or(compaign.Compaign)
	switch offerinfo.SourceType {
	case 1:
		bitm.And(compaign.SourceApp)
	case 2:
		bitm.And(compaign.SourceWebsite)
	}
	return bitm
}

// // 连接类型：wifi：1，mobile：2，any：0
// type ConnTypeFilter struct { //编号4
// 	Code uint32 //过滤器编号1,2,4
// }
//
// func (ctf *ConnTypeFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
// 	offerinfo := offer.GetOfferInfo()
// 	bitm := compaign.Compaign
// 	switch offerinfo.ConnType {
// 	case 1:
// 		bitm.And(compaign.ConnTypeWifi)
// 	case 2:
// 		bitm.And(compaign.ConnTypeMobile)
// 	}
// 	return bitm
// }

//设备类型TABLET,MOBILE,1:手机/平板电脑  2:个人电脑 3:	联网电视 4:	手机 5:	平板电脑 6:	联网设备 7:	数字电视机顶盒
// type DeviceTypeFilter struct { //编号8
// 	Code uint32 //过滤器编号1,2,4
// }
//
// func (dtf *DeviceTypeFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
// 	offerinfo := offer.GetOfferInfo()
// 	bitm := compaign.Compaign
// 	tmpbitm := roaring.NewBitmap()
// 	if offerinfo.DeviceType == 1 || offerinfo.DeviceType == 4 {
// 		tmpbitm.Or(compaign.DeviceMobile)
// 	}
// 	if offerinfo.DeviceType == 5 || offerinfo.DeviceType == 4 {
// 		tmpbitm.Or(compaign.DeviceTablet)
// 	}
// 	if offerinfo.DeviceType == 2 {
// 		tmpbitm.Or(compaign.DeviceDeskTop)
// 	}
// 	bitm.And(tmpbitm)
// 	return bitm
// }

//安卓或苹果的主动接受广告选项是否打开 是：true  否：false
type IdfaGaidFilter struct { //编号16
	Code uint32 //过滤器编号1,2,4
}

func (igf *IdfaGaidFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
	offerinfo := offer.GetOfferInfo()
	bitm := roaring.NewBitmap()
	bitm.Or(compaign.Compaign)
	if offerinfo.IdfaGaid == 0 {
		bitm.AndNot(compaign.IsIdfaGaid)
	}
	return bitm
}

//===============================================
type AndLevelFilter struct {
	filter []LevelFilter //过滤器集合
}

func (af *AndLevelFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
	com := roaring.NewBitmap()
	com.Or(compaign.Compaign)
	for _, v := range af.filter { //遍历过滤器
		com.And(v.Filter(compaign, offer))
		if config.TrackCompaignISON { //追踪代码
			for _, cid := range config.DspTrackCIDs {
				cmp := index.CPINDEX.GetCompaign(cid)
				if !com.Contains(cid) && cmp.ID != 0 {
					trackingLevelFilter(cid, v, offer)
				}
			}
		}
		// pp.Println(com.GetCardinality(), v)
	}
	return com
}

/*
多个过滤器and
*/
func (af *AndLevelFilter) AndFilter(filter ...LevelFilter) {
	af.filter = filter
}

type OrLevelFilter struct {
	filter []LevelFilter //过滤器集合
}

func (af *OrLevelFilter) Filter(compaign index.IndexMap, offer model.Offer) *roaring.Bitmap {
	com := roaring.NewBitmap()
	for _, v := range af.filter { //遍历过滤器
		com.Or(v.Filter(compaign, offer))
	}
	return com
}

/*
多个过滤器or
*/
func (af *OrLevelFilter) OrFilter(filter ...LevelFilter) {
	af.filter = filter
}
