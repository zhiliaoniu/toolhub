<?php

class FastDfsTool
{
    public static function makeToken($putPolicy, $secret)
    {
        $encodedPutPolicy = base64_encode((json_encode($putPolicy,JSON_UNESCAPED_SLASHES)));
        $sign = hash_hmac("sha1", $encodedPutPolicy, $secret, true);
        $encodedSign = base64_encode($sign);
        $signPolicy = $encodedSign.':'.$encodedPutPolicy;
        $token = "BIGO"." ".$signPolicy;

        return $token;

    }

}

class Uploader
{
    public static function curl_request($headers, $form, $buffer, $filename, $url, $ip) {

        if (!is_array($headers)) $headers = [];
        $ch = curl_init();

        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_HTTPHEADER, array_merge(['Expect:'], $headers));
        curl_setopt($ch, CURLOPT_TIMEOUT, 60);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);

        if (is_array($form)) { // form方式
            if ($filename) {
                $form = array_merge($form, ["file" => new cURLFile($filename, mime_content_type($filename), basename($filename))]);

            } else { // buffer
                $form = array_merge($form, ["file" => $buffer]);
            }
            curl_setopt($ch, CURLOPT_POSTFIELDS, $form);
        } else { // body方式
            if ($buffer)
                curl_setopt($ch, CURLOPT_POSTFIELDS, $buffer);
            else {

            }
        }

        $content = curl_exec($ch);

        $httpCode = curl_getinfo($ch, CURLINFO_HTTP_CODE);

        curl_close($ch);
        if ($httpCode != 200) {
            echo "http request failed with status: $httpCode\n";

        } else {
            return $content;

        }

    }

    public static function upload_body_buffer($auth, $buffer, $url, $ip = "") {
        $headers = ["Authorization" => $auth];
        return self::curl_request($headers, null, $buffer, null, $url, $ip);

    }

    public static function upload_body_file($auth, $filename, $url, $ip = "") {
        $headers = ["Authorization" => $auth];
        return self::curl_request($headers, null, null, $filename, $url, $ip);

    }

    public static function upload_form_buffer($auth, $buffer, $url, $ip = "") {
        $form = ["Authorization" => $auth];
        return self::curl_request(null, $form, $buffer, null, $url, $ip);

    }

    public static function upload_form_file($auth, $filename, $url, $ip = "") {
        $form = ["Authorization" => $auth];
        return self::curl_request(null, $form, null, $filename, $url, $ip);

    }

}


function getName($n) {
    $characters = '0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ';
    $randomString = '';

    for ($i = 0; $i < $n; $i++) {
        $index = rand(0, strlen($characters) - 1);
        $randomString .= $characters[$index];

    }

    return $randomString;

}

////////// main

$bucket = "bucket_test_1";
$secret = "GnpLqmmhLWzye5bN";
$domain = "199.91.72.74:8079";
$prefix = getName(3);
$test_file = "a.txt";

file_put_contents($test_file, "Hello World");

$uploadNewPolicy = [
    "bucket"=>$bucket,
    "expires"=>15461856000,
    "fsizeLimit"=>[0,5242880],
    "extName" => "txt"
    ];
$auth = FastDfsTool::makeToken($uploadNewPolicy, $secret);
$res = Uploader::upload_form_file($auth, $test_file, "http://$domain/file/new?bucket=$bucket");
echo "upload new (form, file): $res\n";

$json = json_decode($res, true);
$uploadSlavePolicy = [
    "bucket"=>$bucket,
    "expires"=>15461856000,
    "fsizeLimit"=>[0,5242880],
    "masterUrl" => $json['url'],
    "slavePrefix" => $prefix,
    "extName" => "txt"
    ];
$auth = FastDfsTool::makeToken($uploadSlavePolicy, $secret);
$res = Uploader::upload_form_file($auth, $test_file, "http://$domain/file/new?bucket=$bucket");
echo "new file: $res\n";
unlink($test_file);
