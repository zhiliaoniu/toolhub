import random
import sys 

sites = ['kuaishou', 'xigua', 'bobo', 'miaopai', 'douyin']
tagids = range(0, 1000)
for vid in range(0, int(sys.argv[1])):
    pv = random.randint(0, 100000)
    duration = random.randint(0, 300)
    like = random.randint(0, 10000)
    comment = random.randint(0, 4000)
    time = random.randint(1500000000, 1530000000)
    source = random.choice(sites)
    tagid = random.sample(tagids, 5)
    tagstr = ""
    for i in tagid:
        tagstr += "TAG" + str(i) + "|"

    tagstr = tagstr.strip("|")

    print "{\"id_%d\":[\"TITLE%d\",\"%s\", %d, %d, %d, %d, \"%s\", %d]}" % (vid, vid, source, duration, pv, like, comment, tagstr, time)
