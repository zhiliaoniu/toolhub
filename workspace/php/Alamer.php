<?php
class Alamer {
    static private function config() {
        $config = array();
        $config['uname'] = php_uname('n');
        $uname = explode('.', $config['uname']);
        if (count($uname) > 2) {
            if (substr($uname[0], 0, 2) == 'w-') {
                $uname[0] = substr($uname[0], 2);
            }
            $config['sname'] = "{$uname[0]}.{$uname[2]}"; 
        }
        $config['use_qalarm'] = false;
        $qalarm_path = '/home/q/php/Qalarm/Qalarm.php';
        if (is_readable($qalarm_path)) {
            require_once $qalarm_path;
            if (!method_exists('Qalarm', 'send')) {
                $config['use_qalarm'] = true;
            }
        }
        $config['qalarm_pid'] = "72";
        $config['qalarm_mid'] = "14";
        $config['qalarm_code'] = "1001";
        $config['alarm_id'] = "cloud_rule";
        return $config;

    }

    static private function lvsIp($retry){
        switch($retry)
        {
            case 1:
                $lvsIp = '106.38.184.154';
                break;
            case 2:
                $lvsIp = '111.206.79.29';
                break;
            default:
                break;
        }
        return $lvsIp;
    }

    public static function Alarm($subject, $content) {
        $config = self::config();
        $pid = isset($config['qalarm_pid']) ? $config['qalarm_pid'] : '';
        $mid = isset($config['qalarm_mid']) ? $config['qalarm_mid'] : '';
        $code = isset($config['qalarm_code']) ? $config['qalarm_code'] : '';
        $host = urlencode($config['uname']);
        $message = urlencode(iconv('', 'UTF-8', "{$subject}\t{$content}"));
        $subject = "{$config['sname']}:{$subject}";
        $subject = urlencode(iconv('UTF-8', 'GBK', $subject));
        $content = urlencode(iconv('UTF-8', 'GBK', $content));

        $retry = 2;
        while($retry >= 0)
        {
            if($retry == 0 ){
                $alarm_url = 'http://alarms.ops.qihoo.net:8360/intfs/alarm_intf?'
                    . "group_name={$config['alarm_id']}&subject={$subject}&content={$content}";
            }else{
                $lvsIp = self::lvsIp($retry);
                $alarm_url = "http://$lvsIp/qalarm.php?"
                    ."pid={$pid}&mid={$mid}&code={$code}&host={$host}&content={$message}";
            }

            $alarminfo = file_get_contents($alarm_url);
            if ($alarminfo == 'ok') {
                return true;
            }

            $retry -= 1;
        }
        return false;
    }
}