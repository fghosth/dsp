package index

import (
	"fmt"
	"io/ioutil"

	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

var CPLNAME = "compList.idx"

//compaignl lists index
type CompaignList struct { //单线程维护不加锁
	ComList []model.Compaign
	// lock    *sync.RWMutex
}

func NewCompainList() *CompaignList {
	return &CompaignList{make([]model.Compaign, 0)}
}

//保存磁盘
func (cl *CompaignList) Save() {
	data, _ := util.EncodeStructToByte(cl)
	data = util.DoZlibCompress(data)
	err := ioutil.WriteFile(CPLNAME, data, 0666) //写入文件(字节数组)
	if err != nil {
		fmt.Println(err)
	}
}

//读取磁盘缓存
func (cl *CompaignList) Load() {
	data := util.ReadFile(CPLNAME)
	data = util.DoZlibUnCompress(data)
	err := util.DecodeByteToStruct(data, cl)
	if err != nil {
		fmt.Println(err)
	}
}

//清空
func (cl *CompaignList) Clear() {
	cl.ComList = nil
}

//根据id获取compaign
func (cl *CompaignList) GetCompain(cid uint32) (cm model.Compaign) {
	length := len(cl.ComList)
	for i := 0; i < length; i++ {
		if cl.ComList[i].ID == cid {
			cm = cl.ComList[i]
		}
	}
	return
}

/*插入compain
 */
func (cl *CompaignList) AddCompain(cp model.Compaign) {
	length := len(cl.ComList)

	if length == 0 {
		cl.ComList = append(cl.ComList, cp)
		return
	}

	for i := length - 1; i >= 0; i-- {
		if cp.Score <= cl.ComList[i].Score {
			// cl.lock.Lock()
			cl.ComList = arrInsertAfter(i+1, cl.ComList, cp)
			// cl.lock.Unlock()
			return
		}
	}
	cl.ComList = arrInsertAfter(0, cl.ComList, cp)
}

//删除compain
func (cl *CompaignList) RemoveCompain(cid uint32) {
	length := len(cl.ComList)
	if length == 0 {
		return
	}
	for i := 0; i < length; i++ {
		if cid == cl.ComList[i].ID {
			cl.ComList = arrDeletePos(i+1, cl.ComList)
			return
		}
	}
}

//修改
func (cl *CompaignList) ModifyCompain(cp model.Compaign) {
	length := len(cl.ComList)
	for i := 0; i < length; i++ {
		if cp.ID == cl.ComList[i].ID {
			cl.ComList[i] = cp
			return
		}
	}
}

//数组删除 pos 位置（下标+1） source原数组
func arrDeletePos(pos int, source []model.Compaign) []model.Compaign {
	length := len(source)
	slice1 := source[0 : pos-1 : pos-1]
	slice2 := source[pos:length:length]
	result := append(slice1, slice2...)
	return result
}

//插入数组
func arrInsertAfter(pos int, source []model.Compaign, data model.Compaign) []model.Compaign {
	length := len(source)
	slice1 := source[0:pos:pos]
	slice1 = append(slice1, data)
	slice2 := source[pos:length:length]
	result := append(slice1, slice2...)
	return result
}
