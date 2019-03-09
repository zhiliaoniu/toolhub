package recommendserver

import (
	"math/rand"
	"time"

	"github.com/sumaig/glog"
)

func printIndex(tmp []index) {
	glog.Debug("===== print index =====")
	for _, i := range tmp {
		glog.Debug("item:id:%v, score:%v, source:%v, topicids:%v",
			i.id,
			i.score,
			allVideoItems[i.id].source,
			allVideoItems[i.id].topicids)
	}
	glog.Debug("========= end =========")
}

func shuffle(array []index) {
	rand.Seed(time.Now().UTC().UnixNano())
	for i := len(array) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		array[i], array[j] = array[j], array[i]
	}
}
