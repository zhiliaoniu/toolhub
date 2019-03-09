package recommendserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/sumaig/glog"
)

type item struct {
	title       string
	source      int
	duration    int
	pv          int
	like        int
	comment     int
	tag         []string
	time        uint64
	topicids    []string
	sourcePv    int //视频在原来渠道的pv
	questionNum int //视频是否有问题
}

const (
	SOURCE_PV     = "source_pv"
	HOT           = "hot"
	NEW           = "new"
	TOPIC         = "topic"
	SOURCE_PREFIX = "source"
	SOURCE1       = "source1"
	SOURCE4       = "source4"
	SOURCE5       = "source5"
	SOURCE6       = "source6"
	QUESTION      = "question"

	maxIndexLength = 10000
)

type index struct {
	id    string
	score uint64
}

type sortedIndex []index

func (a sortedIndex) Len() int {
	return len(a)
}

func (a sortedIndex) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a sortedIndex) Less(i, j int) bool {
	return a[j].score < a[i].score
}

// id -> video struct
var allVideoItems map[string]item

// key -> video index
var allVideoIndex map[string][]index

var itemMutex *sync.RWMutex
var indexMutex *sync.RWMutex

func InitIndex() {
	itemMutex = new(sync.RWMutex)
	indexMutex = new(sync.RWMutex)

	//1.get index file name
	indexPath := GetInitIndexFileName()

	//2.reload index
	DoReloadIndex(indexPath)
}

func DoReloadIndex(indexPath string) (err error) {
	fmt.Println("begin init recommend list")
	if strings.HasPrefix(indexPath, "/") == false {
		//TODO alarm
		glog.Error("init recommend index failed. bad index path:%s", indexPath)
		err = errors.New("bad index path")
		return
	}

	//2.reload index
	start := time.Now()
	if err = LoadIndex(indexPath); err != nil {
		glog.Error("load index failed. err:%v", err)
		return
	}
	glog.Debug("init recommend list cost:%d", time.Since(start))

	glog.Debug("new.len:%d, hot.len:%d, question.len:%d", len(allVideoIndex[NEW]), len(allVideoIndex[HOT]), len(allVideoIndex[QUESTION]))
	fmt.Println("end init recommend list")

	return
}

func GetInitIndexFileName() (filename string) {
	filename = "/data/services/recommend-server/conf/default_index"
	defalutPath := "/data/budao-server/recommend/index/"
	cmd := exec.Command("ls", "-ht", defalutPath)
	out, err := cmd.Output()
	if err != nil {
		glog.Error("get cmd output failed.err:%v", err)
		return
	}
	glog.Debug("cmd output:%s", string(out))
	arr := strings.Split(string(out), "\n")
	glog.Debug("arr:%v. len:%d", arr, len(arr))

	if len(arr) <= 1 {
		return
	}
	glog.Debug("arr[0]:%v", arr[0])
	filename = defalutPath + arr[0]
	return
}

func LoadIndex(filename string) (err error) {
	//1.read index file
	file, err := os.Open(filename)
	if err != nil {
		glog.Error("ERROR: open index file err: %v", err)
		return
	}
	defer file.Close()

	dec := json.NewDecoder(file)

	VideoItemsTmp := make(map[string]item, 0)
	VideoIndexTmp := make(map[string][]index, 0)
	for {
		var jsonObj map[string]interface{}
		if err := dec.Decode(&jsonObj); err != nil {
			if strings.Contains(err.Error(), "EOF") {
				glog.Debug("decode end")
				break
			}
			glog.Error("decode json failed. dec:%v, err:%v", dec, err)
			break
		}

		var oneItem item
		for id, oneObj := range jsonObj {
			glog.Debug("id:%s", id)
			if oneRec, ok := oneObj.([]interface{}); ok {
				glog.Debug("oneRec:%v", oneRec)
				oneItem.title = oneRec[0].(string)
				oneItem.source = int(oneRec[1].(float64))
				oneItem.duration = int(oneRec[2].(float64))
				oneItem.pv = int(oneRec[3].(float64))
				oneItem.like = int(oneRec[4].(float64))
				oneItem.comment = int(oneRec[5].(float64))
				oneItem.tag = strings.Split(oneRec[6].(string), "|")
				oneItem.time = uint64(oneRec[7].(float64))
				if oneRec[8].(string) == "" {
					oneItem.topicids = make([]string, 0)
				} else {
					oneItem.topicids = strings.Split(oneRec[8].(string), ",")
					VideoIndexTmp[TOPIC] = append(VideoIndexTmp[TOPIC], index{id, uint64(oneItem.time)})
				}
				oneItem.sourcePv = int(oneRec[9].(float64))
				oneItem.questionNum = int(oneRec[10].(float64))

				VideoIndexTmp[HOT] = append(VideoIndexTmp[HOT], index{id, uint64(oneItem.pv)})
				VideoIndexTmp[NEW] = append(VideoIndexTmp[NEW], index{id, uint64(oneItem.time)})
				VideoIndexTmp[SOURCE_PV] = append(VideoIndexTmp[SOURCE_PV], index{id, uint64(oneItem.sourcePv)})
				if oneItem.questionNum != 0 {
					VideoIndexTmp[QUESTION] = append(VideoIndexTmp[QUESTION], index{id, uint64(oneItem.sourcePv)})
				}
				sourceKey := fmt.Sprintf("source%d", oneItem.source)
				VideoIndexTmp[sourceKey] = append(VideoIndexTmp[sourceKey], index{id, uint64(oneItem.sourcePv)})
			} else {
				glog.Error("ERROR: json format error. oneObj:%v", oneObj)
			}
			VideoItemsTmp[id] = oneItem
		}
	}

	if len(VideoItemsTmp) < 1000 {
		// TODO alert
		glog.Error("ERROR: indexed video count:%d", len(VideoItemsTmp))
	} else {
		for key, indexArr := range VideoIndexTmp {
			sort.Sort(sortedIndex(indexArr))
			indexArr = formatIndex(indexArr)
			glog.Debug("index key:%s, len:%d", key, len(indexArr))
		}

		itemMutex.Lock()
		allVideoItems = VideoItemsTmp
		itemMutex.Unlock()
		glog.Debug("video item len:%d", len(allVideoItems))

		indexMutex.Lock()
		allVideoIndex = VideoIndexTmp
		indexMutex.Unlock()
	}
	return
}

func getVideoItem(id string) (item, bool) {
	itemMutex.RLock()
	i, ret := allVideoItems[id]
	itemMutex.RUnlock()
	return i, ret
}

func getVideoIndex(key string) ([]index, bool) {
	indexMutex.RLock()
	i, ret := allVideoIndex[key]
	indexMutex.RUnlock()
	return i, ret
}

func formatIndex(i []index) []index {
	if len(i) > maxIndexLength {
		return i[:maxIndexLength]
	} else {
		return i
	}
}
