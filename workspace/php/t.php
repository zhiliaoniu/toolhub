#!/usr/local/bin/php
<?php

$nextWeek = time() + (7 * 24 * 60 * 60);
echo 'nextWeek'.$nextWeek."\n";
                   // 7 days; 24 hours; 60 mins; 60 secs
echo 'Now:       '. date('Y-m-d H:i:s', time()) ."\n";
echo 'Next Week: '. date('Y-m-d', $nextWeek) ."\n";
// or using strtotime():
echo 'Next Week: '. date('Y-m-d', strtotime('+1 week')) ."\n";

echo strtotime("Sat, 06 May 2017 05:36:39 GMT");

$f="/home/s/apps/qlogd/log/qlogd_misslog/CloudSafeLine_OutChainNewsUrl_output-0000000005_missing_20170531.log";
echo basename($f);

$arr=array();
if ($arr["aa"] == NULL){
    echo "aa == NULL\n";
}
