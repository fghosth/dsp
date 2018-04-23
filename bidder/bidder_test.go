package bidder_test

import (
	"fmt"
	"testing"

	"jvole.com/dsp/bidder"
	"jvole.com/dsp/model"
)

func TestGetCompaign(t *testing.T) {
	var cmp []model.Compaign
	cmp = make([]model.Compaign, 5)
	result := make(map[uint32]int, len(cmp))
	for i := 0; i < 5; i++ {
		cp := &model.Compaign{}
		cp.ID = uint32(i)
		switch i {
		case 0:
			cp.Score = 10
		case 1:
			cp.Score = 30
		case 2:
			cp.Score = 60
		case 3:
			cp.Score = 60
		case 4:
			cp.Score = 100
		}
		// cp.Score = rand.Intn(100) + 1
		cmp = append(cmp, *cp)
	}

	for i := 0; i < 500; i++ {
		comp := bidder.GetCompain(cmp)
		result[comp.ID] = result[comp.ID] + 1

	}
	for k, v := range result {
		fmt.Println("id分布:", k, v)
	}
}
