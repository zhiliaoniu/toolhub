package recommendserver

import (
	"github.com/sumaig/glog"
)

const maxTopicCount = 1

func isIndexEmpty(candidate map[string][]index, indexCursor map[string]int) (ret bool) {
	ret = true
	for key, i := range candidate {
		if len(i) > indexCursor[key]+1 {
			ret = false
			break
		}
	}
	return
}

func rerank(candidate map[string][]index, count int) (ret []index) {

	maxTypeCountMap := map[int]int{
		1: 3, //type, num
		4: 3,
	}

	typeWeight := map[string]int{
		"source1": 1, //type, num
		"source4": 1,
		"topic":   1,
	}

	maxTypeCount := 1

	typeCount := make(map[int]int)
	topicCount := make(map[string]int)

	retMap := make(map[string]index)

	// record index cursor
	indexCursor := make(map[string]int)
	// make random weight
	var randomChoices []Choice
	// index not selected (for iteration select)
	missIndex := make(map[string][]index, 0)

	// init
	for key, i := range candidate {
		glog.Debug("index rerank init(key, count): %s, %d", key, len(i))
		indexCursor[key] = 0
		weight := 1

		if _, ok := typeWeight[key]; ok {
			weight = typeWeight[key]
		}

		c := Choice{
			weight,
			key,
		}

		randomChoices = append(randomChoices, c)
	}
	glog.Debug("random choices init: %v", randomChoices)

	// TODO: better judge for exit
	for !isIndexEmpty(candidate, indexCursor) {

		// weighted random
		c, err := WeightedRandom(randomChoices)
		if err != nil {
			glog.Debug("random choices err: %v", randomChoices)
			break
		}
		key := c.Item.(string)
		glog.Debug("rand source type:(type, count, cursor) %s, %d, %d", key, len(candidate[key]), indexCursor[key])

		indexCursor[key]++
		if indexCursor[key] >= len(candidate[key]) {
			for i, c := range randomChoices {
				if c.Item == key {
					randomChoices[i].Weight = 0
				}
			}
			glog.Debug("random choices %v", randomChoices)
			continue
		}

		i := candidate[key][indexCursor[key]]
		glog.Debug("select i: %v, %d", i, indexCursor[key])

		if _, ok := retMap[i.id]; ok {
			missIndex[key] = append(missIndex[key], i)
			continue
		}
		glog.Debug("ran item: %v", i)

		source := allVideoItems[i.id].source
		topicIds := allVideoItems[i.id].topicids

		glog.Debug("item start: %v", allVideoItems[i.id])
		if _, ok := maxTypeCountMap[source]; ok {
			maxTypeCount = maxTypeCountMap[source]
		} else {
			maxTypeCount = 1
		}

		if _, ok := typeCount[source]; !ok {
			typeCount[source] = 0
		}

		glog.Debug("type: %v, %v, %v", source, typeCount[source], maxTypeCount)
		if typeCount[source] >= maxTypeCount {
			missIndex[key] = append(missIndex[key], i)
			continue
		} else {
			typeCount[source]++
		}
		glog.Debug("item type ok: %v", allVideoItems[i.id])

		isTopicConflict := false
		for _, topicId := range topicIds {
			if _, ok := topicCount[topicId]; !ok {
				topicCount[topicId] = 0
			}
			if topicCount[topicId] >= maxTopicCount {
				glog.Debug("topic: %v, %v", topicId, topicCount)
				isTopicConflict = true
				break
			} else {
				topicCount[topicId]++
			}
		}

		if isTopicConflict {
			missIndex[key] = append(missIndex[key], i)
			continue
		}
		glog.Debug("item topic ok: %v", allVideoItems[i.id])

		glog.Debug("item ok: %v, %s", allVideoItems[i.id], key)
		retMap[i.id] = i
		glog.Debug("ret item: %v", retMap)

		if len(retMap) >= count {
			break
		}
	}

	for _, i := range retMap {
		ret = append(ret, i)
	}

	if len(ret) == 0 {
		glog.Debug("no candidate selected")
		return
	}

	//
	if len(ret) < count {
		glog.Debug("less")
		for _, i := range rerank(missIndex, count-len(ret)) {
			ret = append(ret, i)
			if len(ret) >= count {
				break
			}
		}
	}
	printIndex(ret)
	return
}
