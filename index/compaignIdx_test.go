package index

import (
	"math/rand"
	"testing"

	"jvole.com/dsp/model"
)

var compidx *CompaignIdx

func TestGetcompaign(t *testing.T) {
	// pp.Println("========", index.CPINDEX.GetCompaign(87).ID)s
}

// func TestRedisDB_tmp(t *testing.T) {
// 	index.CPINDEX.SaveRedis()
// }
func TestAdd(t *testing.T) {
	compidx = NewCompaignIdx()
	cp := &model.Compaign{}
	count := 5
	for i := 0; i < count; i++ {
		cp.ID = uint32(i + 1)
		cp.Score = rand.Intn(1000) + 1
		compidx.Add(*cp)
	}
	cp.ID = uint32(555)
	cp.Score = 16
	compidx.Add(*cp)
	// fmt.Println(compidx.CompaignOrderByScore.ComList[0].Score)
	// for i := 0; i < len(compidx.CompaignOrderByScore.ComList); i++ {
	// 	fmt.Println(compidx.CompaignOrderByScore.ComList[i].Score)
	// 	// fmt.Println(compidx.CompaignOrderByScore)
	// }
	if compidx.GetCPLen() != 6 {
		t.Errorf("错误:长度应该是6.结果是：%d", compidx.GetCPLen())
	}
}

func TestRemove(t *testing.T) {
	compidx.Remove(555)
	key, ok := compidx.compaign.Load(555)
	if ok {
		t.Error("REMOVE错误： 应该被删除id555，实际：", key)
	}
}

func TestSaveDiskCMIDX(t *testing.T) {
	compidx.SaveDisk()
}

func TestLoadDiskCMIDX(t *testing.T) {
	cmidx := NewCompaignIdx()
	cmidx.LoadDisk()
	if compidx.GetCPLen() != 5 {
		t.Errorf("LOAD错误长度应该是5.结果是：%d", compidx.GetCPLen())
	}
}

func aTestSaveRedisCMIDX(t *testing.T) {
	compidx.SaveRedis()
}

func aTestLoadRedisCMIDX(t *testing.T) {
	cmidx := NewCompaignIdx()
	cmidx.LoadRedis()
	if compidx.GetCPLen() != 5 {
		t.Errorf("LOAD错误长度应该是5.结果是：%d", compidx.GetCPLen())
	}
}

func aTestClear(t *testing.T) {
	compidx.Clear()
	if compidx.GetCPLen() > 0 {
		t.Errorf("CLEAR错误长度应该是0.结果是：%d", compidx.GetCPLen())
	}
}
