#include "netproto/include/test_common.h"
#include "cloudkill_test.h"
#include "request.h"
#include <osl/include/ini_parser.h>
#include <gas/qh_glue_common.h>
#include <glog/logging.h>
//#include "glueasync/src/glueasync_types.h"
#include <boost/property_tree/ptree.hpp>  
#include <boost/property_tree/xml_parser.hpp>  
#include <boost/typeof/typeof.hpp> 
#include <iostream>
#include "../security_tool/security_tool.h"

using namespace std;
evqing::EventLoopThread g_loop_thread;
struct event_base* g_evbase;
gas::ResultPtr g_result;

CPPUNIT_TEST_SUITE_REGISTRATION(CloudKillTest);

static evqing::EventLoopThread loop_thread_;
//static struct event_base* evbase_ = NULL;
static bool g_finished = false;

using gas::CloudKillModule;

void CloudKillTest::setUp() {
    g_evbase = g_loop_thread.event_base();
    osl::INIParser ini_paser;
    CPPUNIT_ASSERT(ini_paser.parse("./mock/qlog.conf"));
	CPPUNIT_ASSERT(qh::config()->Initialize(ini_paser));
    CPPUNIT_ASSERT(cloudkill_.InitInMaster("./mock/cloudkill.json"));
}

void CloudKillTest::tearDown() {
}

void  CloudKillTest::ReadRequestData(const char * file_name, std::string& request_data) {
	ifstream in(file_name, std::ifstream::binary | std::ios::in);
	if (!in.is_open()) {
		std::cout << "open file[" << file_name << "] failed" << std::endl;
		exit(0);
	}
	char buf[512];
	int total = 0;
	int size = 0;
	std::string s;
	while (!in.eof()) {
		in.read(buf, 512);
		size = in.gcount();
		request_data.append(buf, size);
		total += size;
	}
	in.close();
}

std::string CloudKillTest::RunEvent(GasCallBack callback, 
		std::function<void (const gas::ResultPtr &)>  result_cb, gas::ssmap& parms, const std::string& request_data, gas::ContextPtr context) {
	    g_finished = false;
        loop_thread_.Start();
        struct event_base* evbase_ = loop_thread_.event_base();
        xstd::function<void (const gas::ResultPtr& result)> f = xstd::bind(result_cb, xstd::placeholders::_1);
        using gas::CloudKillModule;
        loop_thread_.event_loop()->RunInLoop(xstd::bind(callback, &cloudkill_, &parms, context, request_data, evbase_, f));
        while(!g_finished) {
            usleep(1000);
        }
        loop_thread_.Stop(true);

    return string("");
}


void HttpCloudQueryCallback(const gas::ResultPtr& result) {
	std::string s("../security_tool/sec_data_tool -e deres -i \"");
	s.append(osl::Base64::encode(result->result()));
	s.append("\"");
	FILE * fp = popen(s.c_str(), "r");
	CPPUNIT_ASSERT(fp);
	char buf[5201];
	fgets(buf, sizeof(buf), fp);
	pclose(fp);
	std::string res(osl::Base64::decode(buf));
	std::cout << "http cloudquery result=====>" << res << "=======> http cloudquery end" << std::endl;
	CPPUNIT_ASSERT(res[0] == '0');
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessHttpCloudQuery() {
	std::string data;
	ReadRequestData("cloudquery.txt", data);
    CPPUNIT_ASSERT(data.size() > 0);
    osl::StatTime stat_time;
	gas::ssmap params;
	params["x-360-ver"] = "4";

    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::HTTP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessHttpCloudQuery, &HttpCloudQueryCallback, params, data, ctx); 
}

void HttpFileHealthCallback(const gas::ResultPtr& result) {
	std::string s("retinfo code=\"0\" msg=\"Operation success\"");
	CPPUNIT_ASSERT(result->result().find(s) != std::string::npos);
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessHttpFileHealth() {
	std::string data;
	ReadRequestData("filehealth.txt", data);
    CPPUNIT_ASSERT(data.size() > 0);
    osl::StatTime stat_time;
	gas::ssmap params;
	params["boundary"] = "-----------------------------7d83e2d7a141e";
	//params["x-360-ver"] = "4";

    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::HTTP,
                "/file_health_info.php",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessHttpFileHealth, &HttpFileHealthCallback, params, data, ctx); 
}

void HttpGetConfCallback(const gas::ResultPtr& result) {
	std::string s(result->result());
	CPPUNIT_ASSERT(s.find("main") != std::string::npos && s.find("com") != std::string::npos 
			&& s.find("tcp") != std::string::npos && s.find("1100") != std::string::npos && s.find("1101") != std::string::npos);
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessHttpGetConf() {
	std::string data;
    osl::StatTime stat_time;
	gas::ssmap params;

    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::HTTP,
                "/getconf.php",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessHttpGetConf, &HttpGetConfCallback, params, data, ctx); 
}

void HttpStatsCallback(const gas::ResultPtr& result) {
	const std::string &s = result->result();
	std::cout << "stats:" << s << std::endl;  
	CPPUNIT_ASSERT(!s.empty() && s.find("errno=0") != std::string::npos);
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessHttpStats() {
	std::string data;
    osl::StatTime stat_time;
	gas::ssmap params;

    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::HTTP,
                "/stats.php",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessHttpStats, &HttpStatsCallback, params, data, ctx); 
}

void HttpStatusCallback(const gas::ResultPtr& result) {
	std::string s("OK");
	CPPUNIT_ASSERT(s.compare(result->result()) == 0);
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessHttpStatus() {
	std::string data;
    osl::StatTime stat_time;
	gas::ssmap params;

    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::HTTP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessHttpStatus, &HttpStatusCallback, params, data, ctx); 
}

void HttpBoxResultCallback(const gas::ResultPtr& result) {
	std::cout << "result:" << result->result() << std::endl;
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessHttpBoxResult() {
	std::string data;
	//ReadRequestData("box_result.txt", data);
    //CPPUNIT_ASSERT(data.size() > 0);
    osl::StatTime stat_time;
	gas::ssmap params;
	params["md5"] = "fdfsdfsdfdsfs=dfdfdfdfdfsfsdfsdfsdfsdfsfsdfsd"; 

    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::HTTP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessHttpBoxResult, &HttpBoxResultCallback, params, data, ctx); 
}

void UnifyProcess(const gas::ResultPtr& result) {
    g_result = result;
    std::cout << "===========>Result:" << result->result().data() + sizeof(npp::Message::NetHeaderV1) << std::endl;
    g_finished = true;
}

void UnifyProcessByCount(int loop, const gas::ResultPtr& result) {
    g_result = result;
    static int lo = 0;
    std::cout << "===========>Result:" << result->result().data() + sizeof(npp::Message::NetHeaderV1) << std::endl;
    if (++lo >= loop) {
        g_finished = true;
    }
}

void UdpCloudQueryCallback(const gas::ResultPtr& result) {
	const std::string& s = result->result();
	std::cout << "udpcloudquery result====>" << s << "<====udpcloudquery result end" <<s[10] << std::endl;
	CPPUNIT_ASSERT(s[10] == '0');
    g_finished = true;
}

void CloudKillTest::test_UcsProcessUdpCloudQuery() {
    g_finished = false;
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));

    npp::Message::NetHeaderV1 header;
    std::string d = "md5s=d98dd92c867582f16bd28238788095da\t1100\tzhichao\t\tbbb\r\n"
        "product=360safe\r\ncombo=ipc\r\n"
        "mid=42305285be077a7b54dc016a9885411a\r\n"
        "ext=ext:qvm||dlog:4,adsf||dl.ppath:aa.ace||wzc:cbcbb"
                       "type:0\tp1:NtShutdownSystem\tp2:c:\\windows\\system32\\drivers\\xxabcdr.sys\n"
                       "type:0\tp1:NtShutdownSystem\tp2:c:\\windows\\system32\\drivers\\xxabcdr.sys\r\n"
        "osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string request_data;
    request_data.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryCallback, params, request_data, ctx); 
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_1() {
    g_finished = false;
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));

    npp::Message::NetHeaderV1 header;
    std::string d = "bWQ1cz0xYWViYzk1NWZiNWQ2YmVkNjkwZDY1NjU5Yzc1MzdkNis3ODI1MThhY2QwYzJiMGUxOTFmZTIxYmQzNjcwYmZkYWE3MGVjNmUzCTQ5MTg2NAljOlxwcm9ncmFtIGZpbGVzXDM2MFwzNjBzYWZlXGRlZXBzY2FuXGhlYXZ5Z2F0ZS5kbGwJMQkKDQpwcm9kdWN0PWRlZXBzY2FuDQpjb21ibz1wcmltYXJ5DQpzcHJvZHVjdD0zNjBzYWZlDQpzY29tYm89bmV3c3B5X2tpbGxlcg0Kdms9YTQ5ZWU4NTMNCm1pZD0xNmZiOGRhZGE1MDc0OGJhOWM2NDY1OTdlZjA0ZDhhNA0KbGFuZ2lkPTIwNTINCm9zdmVyPTYuMS43NjAwLjI1Ni4xLjANCnB2ZXI9My4zLjguMzAwMA0KdXY9NA0KZXQ9MQ0KaWV2ZXI9OC4wLjc2MDAuMTYzODUNCndpbjY0PTENCmNvbnN1bWU9MywwLDMsMw0K";
    d = osl::Base64::decode(d.data(), d.length());
    header.Init();
    header.set_version(1);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
    std::string request_data;
    request_data.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryCallback, params, request_data, ctx); 
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_2() {
    g_finished = false;
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    {
       //level=50
       std::string d = "md5s=c16aa8b006580a54861b6382cfc74c27+2cee871d0cb5afe9920561a56eee992a16188dc6\t1100\tzhichao\t\tbbb\r\n"
           "product=360safe\r\ncombo=ipc\r\n"
           "mid=42305285be077a7b54dc016a9885411a\r\n"
           "osver=5.1.2600.256.1.3\r\n"
           "vk=vkdebug\r\n";
       header.Init();
       header.set_version(2);
       header.set_reserve(0);
       std::string message;
       message.append((char*)&header, sizeof(header)).append(d);
	   RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryCallback, params, message, ctx); 
    }

    {
       //level=30
       std::string d = "md5s=c16aa8b006580a54861b6382cfc74c27+f1047ebb5e37f70bb07599ab1674972d19a8ee55\t1100\tzhichao\t\tbbb\r\n"
           "product=360safe\r\ncombo=ipc\r\n"
           "mid=42305285be077a7b54dc016a9885411a\r\n"
           "osver=5.1.2600.256.1.3\r\n"
           "vk=vkdebug\r\n";
       header.Init();
       header.set_version(2);
       header.set_reserve(0);
       std::string message;
       message.append((char*)&header, sizeof(header)).append(d);
	   RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryCallback, params, message, ctx); 
    }

    {
        //level=30,50
        std::string d = 
            "md5s=c16aa8b006580a54861b6382cfc74c27+f1047ebb5e37f70bb07599ab1674972d19a8ee55\t1100\tzhichao\t\tbbb\nc16aa8b006580a54861b6382cfc74c27+2cee871d0cb5afe9920561a56eee992a16188dc6\t1100\tzhichao\t\tbbb\r\n"
            "product=360safe\r\ncombo=ipc\r\n"
            "mid=42305285be077a7b54dc016a9885411a\r\n"
            "osver=5.1.2600.256.1.3\r\n"
            "vk=vkdebug\r\n";
        header.Init();
        header.set_version(2);
        header.set_reserve(0);
        std::string message;
        message.append((char*)&header, sizeof(header)).append(d);
	   RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryCallback, params, message, ctx); 
	}
}

void CloudKillTest::test_ExtractLocalDesc() {
    std::string desc = "%&{\"0\":\"Windowsxxxx\",\"1\":\"Windows trandition\"}";
    std::string localdesc = cloudkill::MCFileInfo::ExtractLocalDesc(desc, "0", NULL, cloudkill_.config()->log_name());
    H_TEST_ASSERT(localdesc == "Windowsxxxx");
    localdesc = cloudkill::MCFileInfo::ExtractLocalDesc(desc, "1", NULL, cloudkill_.config()->log_name());
    H_TEST_ASSERT(localdesc == "Windows trandition");
}

void CloudKillTest::test_decrypt() {
    std::string d = "CgMBR0MUAAEAABEBOhx6HLvQlz5ixQyofpkw764ABevBEbpKiRnkvVCqwfmlJ7TQcpvvwmo1iUWaT4bEraL8abdxoKeVbssGnzzwAasF9FWiVaoA81OiBPNVpwHxU6oPpgKjBqUH9FOqB5sHohiiP6I/oj+iP6MDpAeiP6U/mz+bP5gPplKjA6YPoFD0AqRUp1OkUqoE8A/zU6VXplemUKcFoz+jBrwGmwabB5sHmwemBKsPmwGbP5s/30+ycv1V51v3WOZFsnD9WvZT4BbHf5s8pg6jB/cD9lTwD6QC91DzBPQFogD2U6tUo1WmB/cO9AWbB6IYoj+iP6I/oj+jA6QHoj+lP5s/mz+YVfEC9FCgUqAHowWlDqsD9gH3AfEHoQKqD6EB9gLxAqM/owa8BpsGmwabBpsHpwOqBZsBmz+bP6EAooYg/jn4UvwtP5g7mA==";
    d = osl::Base64::decode(d.data(), d.length());
    size_t skip_len = d[12] +10 +3;
    char key1[2] = {(char)0x36, (char)0x92};
    cloudkill::Security security( cloudkill_.processor().get(), cloudkill_.processor()->security_conf_udp());
    security.Xor(&d[skip_len], d.size() - skip_len, key1);
    std::cout << "XXXXXXXXXXXXXXX:" << d.data() + skip_len << std::endl;
    H_TEST_ASSERT(true);
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_midtag() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t3145728\taaa\\bbb\nd6a6e60eeb3bbe4dd455d0991891c840\t562568\tccc\\aaa\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,test:1234\r\n"//||dlog:4,adsf||dl.ppath:aa.ace||wzc:bbb\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        //"osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryCallback, params, message, ctx); 
}

void UdpCloudQueryTimestampCallback(const gas::ResultPtr& result) {
	const std::string& s = result->result();
	int n = s.find("\r\nmd5=");
	std::cout << "udpcloudquery result====>" << s << "<====udpcloudquery result end" << s[10] << std::endl;
	CPPUNIT_ASSERT(s[10] == '0');
	std::string res(s.data() + n + sizeof("\r\nmd5=") - 1);
    std::vector<std::string> md5_vec;
    osl::StringUtil::split(md5_vec, res, "\n");
    std::cout << "md5_vec.size():" << md5_vec.size() << std::endl;
    H_TEST_ASSERT(md5_vec.size() >= 2);
    
    std::vector<std::string> row_vec;
    osl::StringUtil::split(row_vec, md5_vec[0], "\t");
    std::cout << "row_vector:" << osl::ext::DumpVector(row_vec, false) << std::endl;
    H_TEST_ASSERT(row_vec[0] == "f73a83fea9ea0ea702f6b36203c8fa9f");
    //H_TEST_ASSERT(row_vec[1] == "10.0");
    std::string extinfo = osl::Base64::decode(row_vec[9]);
    std::cout << "extinfo:" << extinfo << "\n";
    H_TEST_ASSERT(extinfo.find("<time>1") != std::string::npos);

    row_vec.clear();
    osl::StringUtil::split(row_vec, md5_vec[1], "\t");
    std::cout << "row_vector:" << osl::ext::DumpVector(row_vec, false) << std::endl;
    H_TEST_ASSERT(row_vec[0] == "d6a6e60eeb3bbe4dd455d0991891c840");
    //H_TEST_ASSERT(row_vec[1] == "10.0");
    extinfo = osl::Base64::decode(row_vec[9]);
    std::cout << "extinfo:" << extinfo << "\n";
    H_TEST_ASSERT(extinfo.find("<time>1") != std::string::npos);
    g_finished = true;
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_timestamp() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t3145728\taaa\\bbb\nd6a6e60eeb3bbe4dd455d0991891c840\t562568\tccc\\aaa\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,test:1234||time:1,1\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        //"osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryTimestampCallback, params, message, ctx); 
}

void UdpCloudQueryTimestampFalseCallback(const gas::ResultPtr& result) {
	const std::string& s = result->result();
	int n = s.find("\r\nmd5=");
	std::cout << "udpcloudquery result====>" << s << "<====udpcloudquery result end" << s[10] << std::endl;
	CPPUNIT_ASSERT(s[10] == '0');
	std::string res(s.data() + n + sizeof("\r\nmd5=") - 1);
    std::vector<std::string> md5_vec;
    osl::StringUtil::split(md5_vec, res, "\n");
    std::cout << "md5_vec.size():" << md5_vec.size() << std::endl;
    H_TEST_ASSERT(md5_vec.size() >= 2);
    
    std::vector<std::string> row_vec;
    osl::StringUtil::split(row_vec, md5_vec[0], "\t");
    std::cout << "row_vector:" << osl::ext::DumpVector(row_vec, false) << std::endl;
    H_TEST_ASSERT(row_vec[0] == "f73a83fea9ea0ea702f6b36203c8fa9f");
    //H_TEST_ASSERT(row_vec[1] == "10.0");
	std::string extinfo = osl::Base64::decode(row_vec[9]);
    std::cout << "extinfo:" << extinfo << "\n";
    H_TEST_ASSERT(extinfo.find("<time>1") == std::string::npos);

    row_vec.clear();
    osl::StringUtil::split(row_vec, md5_vec[1], "\t");
    std::cout << "row_vector:" << osl::ext::DumpVector(row_vec, false) << std::endl;
    H_TEST_ASSERT(row_vec[0] == "d6a6e60eeb3bbe4dd455d0991891c840");
    //H_TEST_ASSERT(row_vec[1] == "10.0");
    extinfo = osl::Base64::decode(row_vec[9]);
    std::cout << "extinfo:" << extinfo << "\n";
    H_TEST_ASSERT(extinfo.find("<time>1") == std::string::npos);
    g_finished = true;
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_timestamp1() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t3145728\taaa\\bbb\nd6a6e60eeb3bbe4dd455d0991891c840\t562568\tccc\\aaa\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,adid:10001||time:1\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        //"osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryTimestampFalseCallback, params, message, ctx); 
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_timestamp2() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t3145728\taaa\\bbb\nd6a6e60eeb3bbe4dd455d0991891c840\t562568\tccc\\aaa\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,test:1234||time:1,eJwzBAAAMgAy\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        //"osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryTimestampFalseCallback, params, message, ctx); 
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_timestamp3() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t3145728\taaa\\bbb\nd6a6e60eeb3bbe4dd455d0991891c840\t562568\tccc\\aaa\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,adid:10001||time:3,eJwzBAAAMgAy||\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        //"osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryTimestampCallback, params, message, ctx); 
}

void CloudKillTest::test_UcsProcessUdpCloudQuery_timestamp4() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t3145728\taaa\\bbb\nd6a6e60eeb3bbe4dd455d0991891c840\t562568\tccc\\aaa\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,adid:10001||time:2,MQ==||\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        //"osver=5.1.2600.256.1.3\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryTimestampCallback, params, message, ctx); 
}


void UdpCloudQueryDNSCallback(const gas::ResultPtr& result) {
	const std::string& s = result->result();
	int n = s.find("\r\nmd5=");
	std::cout << "udpcloudquery result====>" << s << "<====udpcloudquery result end" << s[10] << std::endl;
	CPPUNIT_ASSERT(s[10] == '0');
	std::string res(s.data() + n + sizeof("\r\nmd5=") - 1);
    std::vector<std::string> md5_vec;
    osl::StringUtil::split(md5_vec, res, "\n");
    std::cout << "md5_vec.size():" << md5_vec.size() << std::endl;
    H_TEST_ASSERT(md5_vec.size() >= 2);
    
    std::vector<std::string> row_vec;
    osl::StringUtil::split(row_vec, md5_vec[0], "\t");
    std::cout << "row_vector:" << osl::ext::DumpVector(row_vec, false) << std::endl;
    H_TEST_ASSERT(row_vec[0] == "f73a83fea9ea0ea702f6b36203c8fa9f");
    //H_TEST_ASSERT(row_vec[1] == "10.0");
	std::string extinfo = osl::Base64::decode(row_vec[9]);
    std::cout << "extinfo:" << extinfo << "\n";
    H_TEST_ASSERT(extinfo.find("<dns_ip>") != std::string::npos);
    g_finished = true;
}


void CloudKillTest::test_UcsProcessUdpCloudQuery_dns() {
    osl::StatTime stat_time;
	gas::ssmap params;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t562568\tccc\t0\tdns:103.24.2.169||ret_dns_ip:1||dns_test_client_ip:10.19.1.142\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,adid:10001||time:1\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(0);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpCloudQueryDNSCallback, params, message, ctx); 
}

void UdpGetConfCallback(const gas::ResultPtr& result) {
	const std::string & s = result->result();
	std::cout << "udp get conf result======>" << result->result() << "======>udp get conf result end" << std::endl;
	CPPUNIT_ASSERT(s.find("main") != std::string::npos && s.find("com") != std::string::npos 
			&& s.find("tcp") != std::string::npos && s.find("1100") != std::string::npos && s.find("1101") != std::string::npos);
    g_finished = true;
}

void CloudKillTest::test_UcsProcessUdpGetConf()
{/*{{{*/
    PrintLog log("CloudKillTest::test_UcsProcessHttpGetConf");
    g_finished = false;
    osl::StatTime st;
    osl::KVStatLog kvlog;
    gas::ContextPtr ctx( new gas::Context(
                qs::RequestContext::UDP,
                "1",
                "223.255.242.3", &st, xstd::shared_ptr<osl::KVStatLog>()));

    npp::Message::NetHeaderV1 header;
    std::string d = "product=test\r\ncombo=test\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(1);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
	gas::ssmap params;
	RunEvent(&CloudKillModule::UcsProcessUdpCloudQuery, &UdpGetConfCallback, params, message, ctx); 
}/*}}}*/

void UdpStatsCallback(const gas::ResultPtr& result) {
	std::string s(result->result());
	//CPPUNIT_ASSERT(!s.empty() && s.find("errno=0") != std::string::npos);
    g_finished = true;
	
}

void  CloudKillTest::test_UcsProcessUdpStats() {
    osl::StatTime stat_time;
	gas::ssmap params;

    npp::Message::NetHeaderV1 header;
    header.Init();
    header.set_version(2);
    header.set_reserve(2);
    std::string message;
    message.append((char*)&header, sizeof(header));
    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::UDP,
                "/stats.php",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessUdpStats, &UdpStatsCallback, params, message, ctx); 
}

void UdpStatusCallback(const gas::ResultPtr& result) {
	std::cout << result->result() << std::endl;
	CPPUNIT_ASSERT(result->result().find("OK") != std::string::npos);
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessUdpStatus() {
    osl::StatTime stat_time;
	gas::ssmap params;

    npp::Message::NetHeaderV1 header;
    header.Init();
    header.set_version(2);
    header.set_reserve(3);
    std::string message;
    message.append((char*)&header, sizeof(header));
    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::UDP,
                "3",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessUdpStatus, &UdpStatusCallback, params, message, ctx); 
}

void UdpBoxResultCallback(const gas::ResultPtr& result) {
	CPPUNIT_ASSERT(result->result()[10] == '0');
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessUdpBoxResult() {
    osl::StatTime stat_time;
	gas::ssmap params;
    npp::Message::NetHeaderV1 header;
    std::string d = 
        "md5s=f73a83fea9ea0ea702f6b36203c8fa9f+\t562568\tccc\t0\tdns:103.24.2.169||ret_dns_ip:1||dns_test_client_ip:10.19.1.142\r\n"
        "product=360safe\r\n"
        "combo=ipc\r\n"
        "ext=ext:1,adid:10001||time:1,1\r\n"
        "mid=5239ebd29e8621475f38fae77cb761e0\r\n"
        "vk=vkdebug\r\n"
        "nl=2\r\n"
        "uv=4\r\n";
    header.Init();
    header.set_version(2);
    header.set_reserve(6);
    std::string message;
    message.append((char*)&header, sizeof(header)).append(d);
    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::UDP,
                "0",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent(&CloudKillModule::UcsProcessUdpBoxResult, &UdpBoxResultCallback, params, message, ctx); 
}

std::string CloudKillTest::RunEvent1(GasCallBack1 callback, 
		std::function<void (const std::string &)>  result_cb, gas::ssmap& parms, const std::string& request_data, gas::ContextPtr context) {
	    g_finished = false;
        loop_thread_.Start();
        struct event_base* evbase_ = loop_thread_.event_base();
        gas::ssmap params;
        xstd::function<void (const std::string & result)> f = xstd::bind(result_cb, xstd::placeholders::_1);
        using gas::CloudKillModule;
        loop_thread_.event_loop()->RunInLoop(xstd::bind(callback, &cloudkill_, params, context, request_data, evbase_, f));
        while(!g_finished) {
            usleep(1000);
        }
        loop_thread_.Stop(true);

    return string("");
}

void StatusMonitorCallback(const std::string& result) {
	std::string s("OK");
	CPPUNIT_ASSERT(s.compare(result) == 0);
    g_finished = true;
}

void  CloudKillTest::test_UcsProcessStatusMonitor() {
    osl::StatTime stat_time;
	gas::ssmap params;
    npp::Message::NetHeaderV1 header;
    header.Init();
    header.set_version(4);
    header.set_reserve(3);
    std::string message;
    message.append((char*)&header, sizeof(header));
    gas::ContextPtr ctx(new gas::Context(
                qs::RequestContext::UDP,
                "3",
                "202.99.8.1",
                &stat_time, xstd::shared_ptr<osl::KVStatLog>()));
	RunEvent1(&CloudKillModule::UcsProcessStatusMonitor, &StatusMonitorCallback, params, message, ctx); 
}

