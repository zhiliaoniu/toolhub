<?php
error_reporting(0);
define('BASEPATH', dirname(__FILE__) . '/');
require_once("ip_locator.php");
ini_set("memory_limit","-1");

class IpMng {
    function __construct() {
        $this->ip_locator = new IpLocator;
        $this->ip_locator->Initialize(BASEPATH . "final.txt");
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
}
$ip_mng = new IpMng();

$res=array();
$num = 0;
function readipidc(){
    global $ip_mng;
    $file_path = BASEPATH . "ipidc.txt";
    $fd = fopen($file_path, "r");
    while(!feof($fd)) {
        $line = fgets($fd);
        if ($line == false) {
            break;
        }
        $parts = explode("\t", trim($line));
        //print_r($parts);
        if(count($parts) != 2)
          continue;
        $client_ip = $parts[0];
        $ip_info = $ip_mng->GetIpLocationInfo($client_ip);//取出clientip信息 
        if($ip_info == NULL) {
            echo "$ip_info not find";
            continue;
        }

        //print_r($ip_info);
        $idcs = explode(",", trim($parts[1]));
        if(in_array("alisg", $idcs) === FALSE ||
            in_array("alifr", $idcs) === FALSE ||
            in_array("awsuk", $idcs) === FALSE ||
            in_array("awssg", $idcs) === FALSE) {
                continue;
            }

        //print_r($idcs);
        $alisg = array_search("alisg", $idcs);
        $alifr = array_search("alifr", $idcs);
        $awsuk = array_search("awsuk", $idcs);
        $awssg = array_search("awssg", $idcs);
        if(in_array($ip_info[12], array("THAILAND", "INDIA", "SINGAPORE", "MALAYSIA"))){
            $num++;
            if($alisg < $awssg) {
                $res["asia"]["alisg"]+=1; 
            } else {
                $res["asia"]["awssg"]+=1; 
            }
            if($alifr < $awsuk) {
                $res["asia"]["alifr"]+=1; 
            } else {
                $res["asia"]["awsuk"]+=1; 
            }
        } else if(in_array($ip_info[12], array("UNITED KINGDOM", "FRANCE", "RUSSIAN FEDERATION", "GERMANY"))){
            $num++;
            if($alisg < $awssg) {
                $res["europe"]["alisg"]+=1; 
            } else {
                $res["europe"]["awssg"]+=1; 
            }
            if($alifr < $awsuk) {
                $res["europe"]["alifr"]+=1; 
            } else {
                $res["europe"]["awsuk"]+=1; 
            }
        }
        //print_r($res);
        //return;
    }
    fclose($fd);
    echo "num:$num.\n";
    print_r($res);
    //$asia_uk=$res["asia"]["alifr"]+$res["asia"]["awsuk"];
    //echo "alifr:".
}

readipidc();
