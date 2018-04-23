package filter

import (
	"reflect"
	"testing"

	"github.com/RoaringBitmap/roaring"
	"jvole.com/dsp/index"
	"jvole.com/dsp/model"
)

func TestLevelFilterComp(t *testing.T) {
	type args struct {
		compaign index.IndexMap
		offer    model.Offer
	}
	tests := []struct {
		name string
		args args
		want *roaring.Bitmap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		if got := LevelFilterComp(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. LevelFilterComp() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestTypeFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign index.IndexMap
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *roaring.Bitmap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		tf := &TypeFilter{
			Code: tt.fields.Code,
		}
		if got := tf.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. TypeFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestSourceFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign index.IndexMap
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *roaring.Bitmap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		sf := &SourceFilter{
			Code: tt.fields.Code,
		}
		if got := sf.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. SourceFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

// func TestConnTypeFilter_Filter(t *testing.T) {
// 	type fields struct {
// 		Code uint32
// 	}
// 	type args struct {
// 		compaign index.IndexMap
// 		offer    model.Offer
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *roaring.Bitmap
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		ctf := &ConnTypeFilter{
// 			Code: tt.fields.Code,
// 		}
// 		if got := ctf.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
// 			t.Errorf("%q. ConnTypeFilter.Filter() = %v, want %v", tt.name, got, tt.want)
// 		}
// 	}
// // }
//
// func TestDeviceTypeFilter_Filter(t *testing.T) {
// 	type fields struct {
// 		Code uint32
// 	}
// 	type args struct {
// 		compaign index.IndexMap
// 		offer    model.Offer
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 		want   *roaring.Bitmap
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		dtf := &DeviceTypeFilter{
// 			Code: tt.fields.Code,
// 		}
// 		if got := dtf.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
// 			t.Errorf("%q. DeviceTypeFilter.Filter() = %v, want %v", tt.name, got, tt.want)
// 		}
// 	}
// }

func TestIdfaGaidFilter_Filter(t *testing.T) {
	type fields struct {
		Code uint32
	}
	type args struct {
		compaign index.IndexMap
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *roaring.Bitmap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		igf := &IdfaGaidFilter{
			Code: tt.fields.Code,
		}
		if got := igf.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. IdfaGaidFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAndLevelFilter_Filter(t *testing.T) {
	type fields struct {
		filter []LevelFilter
	}
	type args struct {
		compaign index.IndexMap
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *roaring.Bitmap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &AndLevelFilter{
			filter: tt.fields.filter,
		}
		if got := af.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. AndLevelFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestAndLevelFilter_AndFilter(t *testing.T) {
	type fields struct {
		filter []LevelFilter
	}
	type args struct {
		filter []LevelFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &AndLevelFilter{
			filter: tt.fields.filter,
		}
		af.AndFilter(tt.args.filter...)
	}
}

func TestOrLevelFilter_Filter(t *testing.T) {
	type fields struct {
		filter []LevelFilter
	}
	type args struct {
		compaign index.IndexMap
		offer    model.Offer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *roaring.Bitmap
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &OrLevelFilter{
			filter: tt.fields.filter,
		}
		if got := af.Filter(tt.args.compaign, tt.args.offer); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. OrLevelFilter.Filter() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestOrLevelFilter_OrFilter(t *testing.T) {
	type fields struct {
		filter []LevelFilter
	}
	type args struct {
		filter []LevelFilter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		af := &OrLevelFilter{
			filter: tt.fields.filter,
		}
		af.OrFilter(tt.args.filter...)
	}
}
