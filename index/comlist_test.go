package index_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"jvole.com/dsp/index"
	"jvole.com/dsp/model"
)

var CMPL index.CompaignList

func TestAddCompain(t *testing.T) {
	cm := &model.Compaign{}
	count := 20
	cl := index.NewCompainList()
	t1 := time.Now() // get current time
	for i := 0; i < count; i++ {
		cm.ID = uint32(i + 1)
		cm.Score = rand.Intn(1000) + 1
		// fmt.Print(cm.Score, ",")
		cl.AddCompain(*cm)
	}

	cm.ID = uint32(90000)
	cm.Score = 342
	cl.AddCompain(*cm)

	cm.ID = uint32(90001)
	cm.Score = 343
	cl.AddCompain(*cm)
	elapsed := time.Since(t1)
	fmt.Println("App elapsed: ", elapsed)
	CMPL = *cl
	// fmt.Println("rb1.And(rb2)", c.String())

	// for i := 0; i < len(cl.ComList); i++ {
	// 	fmt.Println(cl.ComList[i].Score, cl.ComList[i].ID)
	// }
}

func TestRemoveCompain(t *testing.T) {
	CMPL.RemoveCompain(90000)

	for i := 0; i < len(CMPL.ComList); i++ {
		if CMPL.ComList[i].ID == 90000 {
			fmt.Println("找到了:", CMPL.ComList[i].Score, CMPL.ComList[i].ID)
		}
		// fmt.Println(CMPL.ComList[i].Score, CMPL.ComList[i].ID)
	}
}

func TestModifyCompain(t *testing.T) {
	cm := &model.Compaign{}
	cm.ID = uint32(90001)
	cm.Score = 34200
	cm.Countries.Data = []string{"shanghai", "bengjing"}
	CMPL.ModifyCompain(*cm)

	for i := 0; i < len(CMPL.ComList); i++ {
		if CMPL.ComList[i].ID == 90001 {
			fmt.Println("修改:", CMPL.ComList[i].Score, CMPL.ComList[i].ID, CMPL.ComList[i].Countries.Data)
		}
		// fmt.Println(CMPL.ComList[i].Score, CMPL.ComList[i].ID)
	}
}

func TestSave(t *testing.T) {
	CMPL.Save()
}

func TestLoad(t *testing.T) {
	cl := index.NewCompainList()
	cl.Load()
	for i := 0; i < len(cl.ComList); i++ {
		if cl.ComList[i].ID == 90001 {
			fmt.Println("load修改:", cl.ComList[i].Score, cl.ComList[i].ID, cl.ComList[i].Countries)
		}
		// fmt.Println(CMPL.ComList[i].Score, CMPL.ComList[i].ID)
	}
}
