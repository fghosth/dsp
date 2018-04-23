package index_test

import (
	"testing"

	"github.com/RoaringBitmap/roaring"
	"jvole.com/dsp/index"
)

func TestBitmap(t *testing.T) {

	im := index.NewIndexMap()
	var da []uint32
	da = make([]uint32, 100)
	for i := 0; i < 100; i++ {
		da[i] = uint32(i)
	}
	im.AddCompaign(da)
	da = make([]uint32, 100)
	for i := 51; i < 100; i++ {
		da[i] = uint32(i)
	}
	im.AddTypePopup(da)
	x := im.AndAll([]*roaring.Bitmap{im.Compaign, im.TypePopup}...)
	if x.GetCardinality() != 50 {
		t.Error("错误： 期待50，实际：", x.GetCardinality())
	}
	// fmt.Println(x.GetCardinality())
}

func TestBitmapSave(t *testing.T) {

	im := index.NewIndexMap()
	var da []uint32
	da = make([]uint32, 100)
	for i := 0; i < 100; i++ {
		da[i] = uint32(i)
	}
	im.AddCompaign(da)
	da = make([]uint32, 100)
	for i := 51; i < 100; i++ {
		da[i] = uint32(i)
	}
	im.AddTypePopup(da)
	im.Save()
	// fmt.Println(x.GetCardinality())
}

func TestBitmapLoad(t *testing.T) {

	im := index.NewIndexMap()
	im.Load()
	if im.Compaign.GetCardinality() != 100 {
		t.Error("错误： 期待100，实际：", im.Compaign.GetCardinality())
	}
}
