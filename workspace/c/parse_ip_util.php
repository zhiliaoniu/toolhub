#!/usr/local/bin/php
<?php
$json_file="./cloudkill_online_url_ip_info.txt";
$json_string=file_get_contents($json_file);
$json_obj = json_decode($json_string, true);
$new_array=array();
foreach($json_obj as $key => $value_arr) {
    if($key == "offline")$key="global";
    foreach($value_arr as $value) {
        $new_array[$key][$value]="10";
    }
}
//print_r($new_array);
$json_encode = json_encode($new_array);
$file = fopen("./cloudkill_online_url_ip_info.txt.bak2", "w");
fwrite($file, $json_encode);
fclose($file);
