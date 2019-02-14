#!/usr/local/bin/php
<?php
error_reporting(0);
define('BASEPATH', dirname(__FILE__) . '/');
require_once("ip_locator.php");
ini_set("memory_limit","-1");

class IpMng {
    function __construct() {
        $this->ip_locator = new IpLocator;
        $this->ip_locator->Initialize(BASEPATH . "final.txt", BASEPATH . "ip_range.txt");
    }

    //获取ip位置相关信息
    function GetIpLocationInfo($ip) {
        if(trim($ip) == "") {
            return NULL;
        }
        $ip_info = $this->ip_locator->query($ip);

        //if($ip_info[2] == "unknown" && $ip_info[3] == "unknown") return NULL; 
        return $ip_info;
    }

    //获取ip位置相关信息
    function GetIpRangeInfo($ip) { 
        if(trim($ip) == "") {
            return NULL;
        }
        $ip_info = $this->ip_locator->query_city($ip);
        if(count($ip_info) != 3) {
            return NULL;
        }
        return $ip_info;
    }
}

class Location {
    //public $province;
    //public $city;
    public $edu_net;
    public $pre_mid;
    public $pv;
    public $uv;
    
    function __construct($province_tmp, $city_tmp, $edu_net_tmp, $pre_mid_tmp) {
        //$this->province = $province_tmp;
        //$this->city = $city_tmp;
        if($edu_net_tmp == "教育网")$this->edu_net = true;
        else $this->edu_net = false;
        $this->pre_mid=$pre_mid_tmp;
        $this->pv = 1;
        $this->uv = 1;
    }
    function UpdatePvUv($mid) {
        if($this->pre_mid != $mid) {
            $this->uv+=1;
            $this->pre_mid = $mid;
        }
        $this->pv+=1;
    }
}

class LocationOwn {
    public $province;
    public $city;
    public $edu_net;
    public $pre_mid;
    public $pv;
    public $uv;
    
    function __construct($province_tmp, $city_tmp, $edu_net_tmp, $pre_mid_tmp) {
        $this->province = $province_tmp;
        $this->city = $city_tmp;
        if($edu_net_tmp == "教育网")$this->edu_net = true;
        else $this->edu_net = false;
        $this->pre_mid=$pre_mid_tmp;
        $this->pv = 1;
        $this->uv = 1;
    }
    function UpdatePvUv($mid) {
        if($this->pre_mid != $mid) {
            $this->uv+=1;
            $this->pre_mid = $mid;
        }
        $this->pv+=1;
    }
}

class Dns {
    //public $ip;
    public $location;
    public $user_locations;
    public $record_type;//-1 not set; 0 old; 1 new
    public $show;//需要显示的数据

    function __construct($ip_tmp, $timestamp, $time_split) {
        //$this->ip = $ip_tmp;
        $this->user_locations = array();
        $this->record_type = 0;
        $this->show = false;
        if($timestamp >= $time_split) {
            $this->record_type = 1;
            $this->show = true;
        }
    }
    function FillOwnLocation($ip_mng, $mid, $ip_tmp) {
        $ip_info = $ip_mng->GetIpLocationInfo($ip_tmp);//取出clientip信息 
        if($ip_info===NULL) {
            //echo $client_ip." cannot getiprangeinfo"."\n";
            return;
        }
        $this->location = new LocationOwn($ip_info[2], $ip_info[3], $ip_info[4], $mid);
    }
    function UpdataOwnLocation($ip_mng, $mid, $ip_tmp) {
        //location exist
        if($this->location != NULL){
            $this->location->UpdatePvUv($mid);
        }else {//location not exists
            $this->FillOwnLocation($ip_mng, $mid, $ip_tmp);
        }
    }
}

class QdnsRegionInfo {
    //public $code;
    public $pv;
    function __construct($code_tmp) {
        //$this->code = $code_tmp;
        $this->pv = 1;
    }
}

class Url {
    //public $host;
    //public $ip;
    public $location;
    public $user_locations;
    public $qdns_region_infos;
    function __construct($host_tmp, $ip_tmp) {
        //$this->host = $host_tmp;
        //$this->ip = $ip_tmp;
        $this->user_locations = array();
        $this->qdns_region_infos= array();
    }
    function FillOwnLocation($ip_mng, $mid, $ip_tmp) {
        $ip_info = $ip_mng->GetIpLocationInfo($ip_tmp);//取出clientip信息 
        if($ip_info===NULL) {
            //echo $client_ip." cannot getiprangeinfo"."\n";
            return;
        }
        $this->location = new LocationOwn($ip_info[2], $ip_info[3], $ip_info[4], $mid);
    }
    function UpdataOwnLocation($ip_mng, $mid, $ip_tmp) {
        //location exist
        if($this->location != NULL){
            $this->location->UpdatePvUv($mid);
        }else {//location not exists
            $this->FillOwnLocation($ip_mng, $mid, $ip_tmp);
        }
    }
}

class ResultStat {
    public $pre_mid;
    public $pv;
    public $uv;

    function __construct() {
        $this->pre_mid="noexist";
        $this->pv = 0;
        $this->uv = 0;
    }

    function PrintResult() {
        echo "pv\t" . $this->pv . "\n";
        echo "uv\t" . $this->uv . "\n";
    }

    function UpdatePvUv($mid) {
        if($this->pre_mid != $mid) {
            $this->uv+=1;
            $this->pre_mid = $mid;
        }
        $this->pv+=1;
    }
}

class Clear {
    public $dns_arr;
    public $url_arr;
    public $result_stat;

    function __construct(&$dns_arr, &$url_arr, &$result_stat) {
        $this->dns_arr = &$dns_arr;
        //print_r($dns_arr);
        $this->url_arr = &$url_arr;
        $this->result_stat   = &$result_stat;
    }

    function ClearArrayItem(&$arr, $key){
        $keys = array_keys($arr);  
        $index = array_search($key, $keys); 
        array_splice($arr, $index, 1);
    }

    function ClearPremid__(&$arrs, $key) {
        foreach($arrs as $ip => &$arr){
            $arr->location = (array)$arr->location;
            unset($arr->location[$key]);
            $user_locations = &$arr->user_locations;
            foreach($user_locations as $client_id => &$location){
                $location = (array)$location;
                unset($location[$key]);
            }
        }
    }

    function ClearPremidAndEchoResult() {
        $this->ClearPremid__(&$this->dns_arr, "pre_mid");
        $this->ClearPremid__(&$this->url_arr, "pre_mid");
        //echo "dns_arr_result:\n";
        //print_r($dns_arr);
        $dns_result_arr["dnss"]=$this->dns_arr;
        unset($this->dns_arr);
        echo json_encode($dns_result_arr) . "\n";
        unset($dns_result_arr);
        //echo "UUUUUUUUUUUUUUUUUUU\n";
        //print_r($url_arr);
        $url_result_arr["urls"]=$this->url_arr;
        unset($this->url_arr);
        echo json_encode($url_result_arr) . "\n";
        unset($url_result_arr);

        //echo "dns_arr.count:".count($dns_arr)."\t"."url_arr.count:".count($url_arr)."\n";
        //$result_stat->PrintResult();
        $pv_uv["pv_uv"]=$this->result_stat;
        echo json_encode($pv_uv) . "\n";
    }
}

class WhiteHost {
    public $hosts;

    function __construct() {
        $this->hosts = array();
    }

    function Init($url){
        $ch = curl_init($url);
        curl_setopt($ch , CURLOPT_RETURNTRANSFER, true);
        curl_setopt($ch , CURLOPT_HTTP_VERSION, CURL_HTTP_VERSION_1_0);
        curl_setopt($ch , CURLOPT_TIMEOUT, 10);

        $ret = curl_exec($ch);
        $http_code = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $http_err = curl_error($ch);
        //echo $http_code,":", $http_err,"\n";
        $ret = json_decode($ret, true);
        $host_arr = explode(",", $ret["hosts"]);
        //print_r($ret);
        foreach($ret as $host => $value){
            if($value == 1 || $host=="kjjs.360.cn") continue;
            $this->hosts[$host]=$value;
        }
        //print_r($this->hosts);
    }
    function IsWhiteHost($host) {
        if(array_key_exists($host, $this->hosts))return true;
        return false;
    }
}


function IsLocalIp($ip) {
    //if ( $ip == "127.0.0.1" ) return true;

    $ipNum = ip2long($ip);

    //0.0.0.0/8
    if ( ($ipNum & 0xFF000000) == 0x00000000 ) return true;

    //127.0.0.0/8
    if ( ($ipNum & 0xFF000000) == 0x7F000000 ) return true;

    //10.0.0.0/8
    if ( ($ipNum & 0xFF000000) == 0x0A000000 ) return true;

    //172.16.0.0/12
    if ( ($ipNum & 0xFFF00000) == 0xAC100000 ) return true;

    //192.168.0.0/16
    if ( ($ipNum & 0xFFFF0000) == 0xC0A80000 ) return true;

    return false;
}

function CheckIp($ip){
    $arr=explode(".", $ip);
    if(count($arr)!=4)return false;
    foreach($arr as $data) {
        if($data<0 || $data>255)return false;
    }
    return true;
}

function reduce() {
    //$buf = "memory_size:\t" . memory_get_usage() . "\n";
    //fwrite(STDERR, $buf);
    $ip_mng = new IpMng();
    $result_stat = new ResultStat();
    $white_host=new WhiteHost();
    //$white_host->Init("http://180.163.251.205/data/special_kill/white_dnshost");
    $white_host->Init("http://dp.safe.qihoo.net:8360/index.php/api/host");
    $time_split = mktime(0, 0, 0, date("m")  , date("d")-1, date("Y"));//昨天凌晨 
    //echo "time_split:". $time_split . "\n";
    //$buf = "memory_size:\t" . memory_get_usage() . "\n";
    //fwrite(STDERR, $buf);

    $dns_arr = array();
    $url_arr = array();

    $num = 0;
    while (!feof(STDIN))
    {
        $line = trim(fgets(STDIN));
        if(!$line) {
            continue;
        }

        $cols = explode("\t",$line);
        //print_r($cols);
        //echo count($cols)."\n";
        if(count($cols)!=7) {
            continue;
        }

        //查询clientip所属城市编号
        //解析请求数据
        $mid       = $cols[0];
        $client_ip = $cols[2];
        $host      = $cols[3];
        $hostips   = array();
        if($cols[4] != "-") {
            if(strpos($cols[4], ",") === false){
                $hostips = explode("|", trim($cols[4]));
            }else {
                //echo $cols[4]."\n";
                continue;
            }
        }
        $dnsips = array();
        if($cols[5] != "-") {
            if(strpos($cols[5], ",") === false){
                $dnsips= explode("|", trim($cols[5]));
            }else {
                //echo $cols[5]."\n";
                continue;
            }
        }
        $timestamp = $cols[6];

        /*
        echo $mid . "\t" . $client_ip . "\t" . $host . "\t" . $timestamp . "\n";
        print_r($hostips);
        print_r($dnsips);
        */

        //统计昨天pv、uv
        if($timestamp > $time_split) {
            $result_stat->UpdatePvUv($mid);
        } 
        $num+=1;
        if($num % 5000000 == 0){
            $num=0;
            $buf = "memory_size:\t" . memory_get_usage() . "\n";
            fwrite(STDERR, $buf);
        }
        //echo "line_num:".$num."\t"."pv:".$result_stat->pv."\t"."uv:".$result_stat->uv."\n";

        $ip_info = $ip_mng->GetIpLocationInfo($client_ip);//取出clientip信息 
        if($ip_info===NULL) {
            //echo $client_ip." cannot getiprangeinfo"."\n";
            continue;
        }
        //print_r($ip_info);

        $client_id = $ip_info[2] . ":" . $ip_info[3];
        //分析所有DNS IP
        if(count($dnsips)>1){
            $dnsips=array_slice($dnsips, 0, 1);
        }

        foreach($dnsips as $dnsip) {
            //dnsip existed
            if(CheckIp($dnsip)===false || IsLocalIp($dnsip)===true)continue;
            if(array_key_exists($dnsip, $dns_arr)){
                $dns = $dns_arr[$dnsip];
                $dns->UpdataOwnLocation($ip_mng, $mid, $dnsip);
                //province city exists
                if(array_key_exists($client_id, $dns->user_locations)){
                    $dns->user_locations[$client_id]->UpdatePvUv($mid);
                }else {//province city not exists
                    $location = new Location($ip_info[2], $ip_info[3], $ip_info[4], $mid);
                    $dns->user_locations[$client_id]=$location;
                }
                //判断dns新旧，-1 not set; 0 old; 1 new
                if($dns->record_type == 1 && $timestamp < $time_split) $dns->record_type = 0;
                if($dns->show == false && $timestamp > $time_split)$dns->show = true;
                //print_r($dns);
            } else {//dnsip not existed
                $dns = new Dns($dnsip, $timestamp, $time_split);
                $dns->UpdataOwnLocation($ip_mng, $mid, $dnsip);
                $location = new Location($ip_info[2], $ip_info[3], $ip_info[4], $mid);
                $dns->user_locations[$client_id]=$location;
                $dns_arr[$dnsip]=$dns;
            }
            //print_r($dns_arr);
        }

        //统计URL相关数据
        if(!$white_host->IsWhiteHost($host))continue;
        foreach($hostips as $ip){
            if(CheckIp($ip)===false || IsLocalIp($ip)===true)continue;
            $hostip=$host . ":" . $ip;
            $ip_range = $ip_mng->GetIpRangeInfo($client_ip);
            //print_r($ip_range);
            if($ip_range == NULL) continue;
            $code = $ip_range[2];
            // hostip existed
            if(array_key_exists($hostip, $url_arr)){
                $url = $url_arr[$hostip];
                $url->UpdataOwnLocation($ip_mng, $mid, $ip);
                //province city exists
                if(array_key_exists($client_id, $url->user_locations)){
                    $url->user_locations[$client_id]->UpdatePvUv($mid);
                }else {//province city not exists
                    $location = new Location($ip_info[2], $ip_info[3], $ip_info[4], $mid);
                    $url->user_locations[$client_id]=$location;
                }
                if(array_key_exists($code, $url->qdns_region_infos)){
                    $dns_region_info = $url->qdns_region_infos[$code];
                    $dns_region_info->pv+=1;
                }else {
                    $dns_region_info = new QdnsRegionInfo($code);
                    $url->qdns_region_infos[$code] = $dns_region_info;
                }
            } else {// hostip not existed
                $url = new Url($host, $ip);
                $url->UpdataOwnLocation($ip_mng, $mid, $ip);

                $location = new Location($ip_info[2], $ip_info[3], $ip_info[4], $mid);
                $url->user_locations[$client_id]=$location;

                $dns_region_info = new QdnsRegionInfo($code);
                $url->qdns_region_infos[$code] = $dns_region_info;
                $url_arr[$hostip]=$url;
            }
        }
    }
    $clear_manager = new Clear(&$dns_arr, &$url_arr, &$result_stat);
    $clear_manager->ClearPremidAndEchoResult();
}

reduce(); 

