<?php
//Read the iplocation file in

class Vector{
    public $index;

}
class IpLocator { 
    function Initialize($file_path, $ip_range_path) {
        $this->arr = array();
        $fd = fopen($file_path, "r");
        $count  = 0;
        while(!feof($fd)) {
            ++$count;
           // if ($count % 10000 == 0) {
           //     $buf = "memory_size:\t" . memory_get_usage() . "\n";
           //     fwrite(STDERR, $buf);
           // }

            $line = fgets($fd);
            if ($line == false) {
                break;
            }
            $parts = explode("\t", trim($line), 6);
            if(count($parts) < 6)
                continue;
            unset($parts[5]);
            $this->arr[] = $parts;
        }
        fclose($fd);

        //加载ip_range_region
        $this->ip_range_arr = array();
        $fd2 = fopen($ip_range_path, "r");
        while(!feof($fd2)) {
            $line = fgets($fd2);
            if ($line == false) {
                break;
            }
            $parts = explode(" ", trim($line));
            if(count($parts) < 3)
                continue;
            $this->ip_range_arr[] = $parts;
        }
        //print_r($this->ip_range_arr);
        //echo count($this->ip_range_arr) . "\n";
        fclose($fd2);
        return true;
    }

    //ip_locator function
    function query($ip){
        $ip = ip2long($ip);
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

    //ip_locator function
    function query_city($ip){
        $ip = ip2long($ip);
        $begin = 0;
        $end = count($this->ip_range_arr) - 1;
        while($begin <= $end) {
            $mid = intval(($begin + $end)/2);
            if ($this->ip_range_arr[$mid][0] <= $ip && $this->ip_range_arr[$mid][1] >= $ip) {
                return $this->ip_range_arr[$mid];
            }
            if ($this->ip_range_arr[$mid][0] > $ip) {
                $end = $mid - 1;
            }
            if ($this->ip_range_arr[$mid][0] < $ip) {
                $begin = $mid + 1;
            }
        }
        return NULL;
    }
}
