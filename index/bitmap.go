package index

import (
	"io/ioutil"
	"path"
	"runtime"

	"github.com/RoaringBitmap/roaring"
	"github.com/go-kit/kit/log/level"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

const (
	CPNAME = "cp.bm"
)

type IndexMap struct {
	Compaign      *roaring.Bitmap //所有待投广告的compaignid
	TypeNative    *roaring.Bitmap //所有广告类别为Native的compaignid
	TypeBanner    *roaring.Bitmap //所有广告类别为Banner的compaignid
	TypePopup     *roaring.Bitmap //所有广告类别为Banner的compaignid
	SourceApp     *roaring.Bitmap //所有资源类别为App的compaignid
	SourceWebsite *roaring.Bitmap //所有资源类别为website的compaignid
	// ConnTypeWifi   *roaring.Bitmap //所有连接类别为wifi的compaignid
	// ConnTypeMobile *roaring.Bitmap //所有连接类别为Mobile的compaignid
	// DeviceTablet  *roaring.Bitmap //所有设备类别为Tablet的compaignid
	// DeviceDeskTop *roaring.Bitmap //所有设备类别为Desktop的compaignid
	// DeviceMobile  *roaring.Bitmap //所有设备类别为Mobile的compaignid
	IsIdfaGaid    *roaring.Bitmap //IsIdfaGaid为true的compaignid
	IsNotIdfaGaid *roaring.Bitmap //IsIdfaGaid为false的compaignid

}

func NewIndexMap() *IndexMap {

	im := &IndexMap{
		roaring.NewBitmap(),
		roaring.NewBitmap(),
		roaring.NewBitmap(),
		roaring.NewBitmap(),
		roaring.NewBitmap(),
		roaring.NewBitmap(),
		roaring.NewBitmap(),
		// roaring.NewBitmap(),
		// roaring.NewBitmap(),
		// roaring.NewBitmap(),
		// roaring.NewBitmap(),
		// roaring.NewBitmap(),
		roaring.NewBitmap(),
	}
	return im
}

type ComCid struct { //添加，修改compaign要保存的结构，如果是0则不添加
	CID              uint32
	TypeNativeCID    uint32
	TypeBannerCID    uint32
	TypePopupCID     uint32
	SourceAppCID     uint32
	SourceWebsiteCID uint32
	// ConnTypeWifiCID   uint32
	// ConnTypeMobileCID uint32
	DeviceTabletCID  uint32
	DeviceDeskTopCID uint32
	DeviceMobileCID  uint32
	IsIdfaGaidCID    uint32
	IsNotIdfaGaidCID uint32
}

//保存磁盘
func (im *IndexMap) Save() {
	data, _ := util.EncodeStructToByte(im)
	data = util.DoZlibCompress(data)
	err := ioutil.WriteFile(CPNAME, data, 0666) //写入文件(字节数组)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "写磁盘错误",
			"err", err,
		)
	}
}

//读取磁盘缓存
func (im *IndexMap) Load() {
	data := util.ReadFile(CPNAME)
	data = util.DoZlibUnCompress(data)
	err := util.DecodeByteToStruct(data, im)
	if err != nil {
		pc, file, line, _ := runtime.Caller(1)
		f := runtime.FuncForPC(pc)
		level.Error(*logger).Log(
			"method", f.Name(),
			"file", path.Base(file),
			"line", line,
			"msg", "byte转结构体错误",
			"err", err,
		)
	}
}

//添加一条记录修改bitmap
func (im *IndexMap) Add(cp model.Compaign) {
	cids := getComCid(cp)
	im.Compaign.Add(cids.CID)
	im.TypeNative.Add(cids.TypeNativeCID)
	im.TypeBanner.Add(cids.TypeBannerCID)
	im.TypePopup.Add(cids.TypePopupCID)
	im.SourceApp.Add(cids.SourceAppCID)
	im.SourceWebsite.Add(cids.SourceWebsiteCID)
	// im.ConnTypeWifi.Add(cids.ConnTypeWifiCID)
	// im.ConnTypeMobile.Add(cids.ConnTypeMobileCID)
	// im.DeviceTablet.Add(cids.DeviceTabletCID)
	// im.DeviceDeskTop.Add(cids.DeviceDeskTopCID)
	// im.DeviceMobile.Add(cids.DeviceMobileCID)
	im.IsIdfaGaid.Add(cids.IsIdfaGaidCID)
	im.IsNotIdfaGaid.Add(cids.IsNotIdfaGaidCID)
}

//删除一条记录修改bitmap
func (im *IndexMap) Remove(cid uint32) {
	im.Compaign.Remove(cid)
	im.TypeNative.Remove(cid)
	im.TypeBanner.Remove(cid)
	im.TypePopup.Remove(cid)
	im.SourceApp.Remove(cid)
	im.SourceWebsite.Remove(cid)
	// im.ConnTypeWifi.Remove(cid)
	// im.ConnTypeMobile.Remove(cid)
	// im.DeviceTablet.Remove(cid)
	// im.DeviceDeskTop.Remove(cid)
	// im.DeviceMobile.Remove(cid)
	im.IsIdfaGaid.Remove(cid)
	im.IsNotIdfaGaid.Remove(cid)
}

//清除某一个bitmap里面的所有内容
func (im *IndexMap) ClearAll() {
	im.Compaign.Clear()
	im.TypeNative.Clear()
	im.TypeBanner.Clear()
	im.TypePopup.Clear()
	im.SourceApp.Clear()
	im.SourceWebsite.Clear()
	// im.ConnTypeWifi.Clear()
	// im.ConnTypeMobile.Clear()
	// im.DeviceTablet.Clear()
	// im.DeviceDeskTop.Clear()
	// im.DeviceMobile.Clear()
	im.IsIdfaGaid.Clear()
	im.IsNotIdfaGaid.Clear()
}

//所有集合交集
func (im *IndexMap) AndAll(bm ...*roaring.Bitmap) *roaring.Bitmap {
	return roaring.FastAnd(bm...)
}

//清除某一个bitmap里面的所有内容
func (im *IndexMap) ClearOne(bm roaring.Bitmap) {
	bm.Clear()
}
func (im *IndexMap) AddIsNotIdfaGaid(cid []uint32) {
	im.IsNotIdfaGaid.AddMany(cid)
}

func (im *IndexMap) ReomveIsNotIdfaGaid(cid uint32) {
	im.IsNotIdfaGaid.Remove(cid)
}
func (im *IndexMap) AddIsIdfaGaid(cid []uint32) {
	im.IsIdfaGaid.AddMany(cid)
}

func (im *IndexMap) ReomveIsIdfaGaid(cid uint32) {
	im.IsIdfaGaid.Remove(cid)
}

// func (im *IndexMap) AddDeviceMobile(cid []uint32) {
// 	im.DeviceMobile.AddMany(cid)
// }
//
// func (im *IndexMap) ReomveDeviceMobile(cid uint32) {
// 	im.DeviceMobile.Remove(cid)
// }
//
// func (im *IndexMap) AddDeviceDeskTop(cid []uint32) {
// 	im.DeviceDeskTop.AddMany(cid)
// }
//
// func (im *IndexMap) ReomveDeviceDeskTop(cid uint32) {
// 	im.DeviceDeskTop.Remove(cid)
// }
//
// func (im *IndexMap) AddDeviceTablet(cid []uint32) {
// 	im.DeviceTablet.AddMany(cid)
// }
//
// func (im *IndexMap) ReomveDeviceTablet(cid uint32) {
// 	im.DeviceTablet.Remove(cid)
// }

// func (im *IndexMap) AddConnTypeMobile(cid []uint32) {
// 	im.ConnTypeMobile.AddMany(cid)
// }
//
// func (im *IndexMap) ReomveConnTypeMobile(cid uint32) {
// 	im.ConnTypeMobile.Remove(cid)
// }
//
// func (im *IndexMap) AddConnTypeWifi(cid []uint32) {
// 	im.ConnTypeWifi.AddMany(cid)
// }
//
// func (im *IndexMap) ReomveConnTypeWifi(cid uint32) {
// 	im.ConnTypeWifi.Remove(cid)
// }
func (im *IndexMap) AddSourceWebsite(cid []uint32) {
	im.SourceWebsite.AddMany(cid)
}

func (im *IndexMap) ReomveSourceWebsite(cid uint32) {
	im.SourceWebsite.Remove(cid)
}
func (im *IndexMap) AddSourceApp(cid []uint32) {
	im.SourceApp.AddMany(cid)
}

func (im *IndexMap) ReomveSourceApp(cid uint32) {
	im.SourceApp.Remove(cid)
}
func (im *IndexMap) AddTypePopup(cid []uint32) {
	im.TypePopup.AddMany(cid)
}

func (im *IndexMap) ReomveTypePopup(cid uint32) {
	im.TypePopup.Remove(cid)
}

func (im *IndexMap) AddTypeBanner(cid []uint32) {
	im.TypeBanner.AddMany(cid)
}

func (im *IndexMap) ReomveTypeBanner(cid uint32) {
	im.TypeBanner.Remove(cid)
}
func (im *IndexMap) AddTypeNative(cid []uint32) {
	im.TypeNative.AddMany(cid)
}

func (im *IndexMap) ReomveTypeNative(cid uint32) {
	im.TypeNative.Remove(cid)
}
func (im *IndexMap) AddCompaign(cid []uint32) {
	im.Compaign.AddMany(cid)
}

func (im *IndexMap) ReomveCompaign(cid uint32) {
	im.Compaign.Remove(cid)
}

//根据compaign获得comcid
func getComCid(cp model.Compaign) (cids ComCid) {
	cids.CID = cp.ID
	switch cp.Type {
	case 1:
		cids.TypeNativeCID = cp.ID
	case 2:
		cids.TypeBannerCID = cp.ID
	case 3:
		cids.TypePopupCID = cp.ID
	}
	switch cp.SourceType {
	case 0:
		cids.SourceAppCID = cp.ID
		cids.SourceWebsiteCID = cp.ID
	case 1:
		cids.SourceAppCID = cp.ID
	case 2:
		cids.SourceWebsiteCID = cp.ID
	}
	// switch cp.ConnectionType {
	// case 0:
	// 	cids.ConnTypeWifiCID = cp.ID
	// 	cids.ConnTypeMobileCID = cp.ID
	// case 1:
	// 	cids.ConnTypeWifiCID = cp.ID
	// case 2:
	// 	cids.ConnTypeMobileCID = cp.ID
	// }
	// for _, v := range cp.DeviceTypes.Data {
	// 	switch v { //TABLET,MOBILE,DESKTOP
	// 	case "5":
	// 		cids.DeviceTabletCID = cp.ID
	// 	case "4":
	// 		cids.DeviceMobileCID = cp.ID
	// 	case "2":
	// 		cids.DeviceDeskTopCID = cp.ID
	// 	}
	// }
	if cp.IsIdfaGaid {
		cids.IsIdfaGaidCID = cp.ID
	} else {
		cids.IsNotIdfaGaidCID = cp.ID
	}
	return
}
