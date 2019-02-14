<?php
//----------------------------
//Read the iplocation file in
class IpLocator { 
    function Initialize($file_path) {
        $this->arr = array();
        $fd = fopen($file_path, "r");
        $this->ip_loc_arr = array();
        while(!feof($fd)) {
            $line = fgets($fd);
            if ($line == false) {
                break;
            }
            $parts = explode("\t", trim($line));
            if(count($parts) < 11)
              continue;
            $this->arr[] = $parts;
        }
        fclose($fd);
        return true;
    }


    //----------------------------
    //ip_locator function
    function query($ip){
        //$ip = ip2long($ip);
        $begin = 0;
        $end = count($this->arr) - 1;
        while($begin <= $end) {
            $mid = intval(($begin + $end)/2);
            if ($this->arr[$mid][0] <= $ip && $this->arr[$mid][1] >= $ip) {
                return $this->arr[$mid];
            }
            if ($this->arr[$mid][0] > $ip) {
                $end = $mid - 1;
            }
            if ($this->arr[$mid][0] < $ip) {
                $begin = $mid + 1;
            }
        }
        return NULL;
    }
}
