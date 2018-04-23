package bidder

import (
	"math/rand"
	"time"

	"github.com/go-kit/kit/log"
	"jvole.com/dsp/model"
	"jvole.com/dsp/util"
)

var logger log.Logger

func init() {
	logger = util.KitLogger
	log.With(logger, "component", "Bidder")
}

type Bidder struct{}

/*根据score随机获得一个compaign
@param model.CompainIndex model.Compaign数组
@return model.Compaign
*/
func GetCompain(cp []model.Compaign) model.Compaign {
	var comp model.Compaign
	var sum int
	for _, v := range cp {

		sum = sum + v.Score
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := r.Intn(sum)
	for i := 0; i < len(cp); i++ {
		t = t - cp[i].Score
		if t < 0 {
			comp = cp[i]
			break
		}
	}
	return comp
}
