import simplejson
import demjson
import json
import time

class JsonTest():
    def test(self):
        d = {
    "out_put_pos": 1,
        "picker": {
            "class": "DictPicker",
            "manager_name":"proc_behavior"
        },
        "header": {
            "class": "TimeStampHeader",
            "key_pos": 1,
            "header_creator": {
                "module": "ycs_proc_chain",
                "class": "proc_behavior_header_creator"
            }
        },
        "events": {
            "type": {
                "class": "UniqEvent",
                "seg_creator": {
                    "class": "key_type_seg_creator",
                    "module": "proc_behavior"
                }
            },
            "version": {
                "class": "UniqEvent",
                "pos": "hi.tfv"
            },
            "size": {
                "class": "UniqEvent",
                "pos": "hi.dsi"
            },
            "dna": {
                "class": "UniqEvent",
                "seg_creator": {
                    "class": "dna_seg_creator",
                    "module": "proc_behavior"
                }
            },
            "icon_dna": {
                "class": "UniqEvent",
                "seg_creator": {
                    "class": "icon_dna_seg_creator",
                    "module": "proc_behavior"
                }
            },
            "ssdeep": {
                "class": "UniqEvent",
                "seg_creator": {
                    "class": "ssdeep_seg_creator",
                    "module": "proc_behavior"
                }
            },
            "icon_ssdeep": {
                "class": "UniqEvent",
                "seg_creator": {
                    "class": "icon_ssdeep_seg_creator",
                    "module": "proc_behavior"
                }
            },

            "pmd5": {
                "class": "UniqEvent",
                "pos": "hi.pmd5"
            },
            "related": {
                "class": "UniqListHolderEvent",
                "max_len": 200,
                "need_update_functor":{
                    "class": "need_related_md5_sha1",
                    "module": "proc_behavior"
                },
                "primary_key":["md5", "sha1"],
                "primary_key_creator": {
                    "class": "default_primary_key_creator",
                    "module": "container"
                },
                "childs": {
                    "md5": {
                        "class": "UniqEvent",
                        "seg_creator": {
                            "class": "related_md5_creator",
                            "module": "proc_behavior"
                        }
                    },
                    "sha1": {
                        "class": "UniqEvent",
                        "seg_creator": {
                            "class": "related_sha1_creator",
                            "module": "proc_behavior"
                        }
                    }
                }
            },
            "signature": {
                "class": "DictHolderEvent",
                "childs": {
                    "time": {
                        "class": "UniqEvent",
                        "pos": "hi.signtime"
                    },
                    "fingerprint": {
                        "class": "UniqEvent",
                        "seg_creator": {
                            "class": "fingerprint_creator",
                            "module": "proc_behavior"
                        }
                    },
                    "internal_name": {
                        "class": "UniqEvent",
                        "pos": "hi.itn"
                    },
                    "original_name": {
                        "class": "UniqEvent",
                        "pos": "hi.orn"
                    },
                    "product": {
                        "class": "UniqEvent",
                        "pos": "hi.gen"
                    },
                    "signers": {
                        "class": "UniqEvent",
                        "pos": "hi.sig"
                    }
                }
            },
            "parent": {
                "class": "UniqListHolderEvent",
                "max_len": 200,
                "primary_key":["md5"],
                "primary_key_creator": {
                    "class": "default_primary_key_creator",
                    "module": "container"
                },
                "childs": {
                    "path": {
                        "class": "UniqEvent",
                        "pos": "hi.src"
                    },
                    "wmi_path": {
                        "class": "UniqEvent",
                        "pos": "hi.wsrc"
                    },
                    "thread_mod": {
                        "class": "UniqEvent",
                        "pos": "hi.tmod"
                    },
                    "md5": {
                        "class": "UniqEvent",
                        "pos": "hi.md5par"
                    },
                    "parent_path": {
                        "class": "UniqEvent",
                        "pos": "hi.psrc"
                    },
                    "file_defender": {
                        "class": "DictHolderEvent",
                        "childs": {
                            "link_path": {
                                "class": "UniqEvent",
                                "pos": "hi.lnk"
                            },
                            "link_dest_path": {
                                "class": "UniqEvent",
                                "pos": "hi.lnkdst"
                            }
                        }
                    }
                }
            },
            "traceable": {
                "class": "DictHolderEvent",
                "childs": {
                    "dl_info": {
                        "class": "UniqListHolderEvent",
                        "max_len": 200,
                        "primary_key":["downlad_url", "parent_url"],
                        "primary_key_creator": {
                            "class": "default_primary_key_creator",
                            "module": "container"
                        },
                        "childs": {
                            "downlad_url": {
                                "class": "UniqEvent",
                                "pos": "hi.durl"
                            },
                            "parent_url": {
                                "class": "UniqEvent",
                                "pos": "hi.purl"
                            },
                            "from_zip": {
                                "class": "UniqEvent",
                                "pos": "hi.uzp"
                            }
                        }
                    },
                    "download_tool": {
                        "class": "DictHolderEvent",
                        "childs": {
                        }
                    }
                }
            },
            "actions": {
                "class": "DictHolderEvent",
                "childs": {
                    "inject": {
                        "class": "ListEvent",
                        "max_len": 200,
                        "pos":"hi.injdst"
                    },
                    "defended": {
                        "class": "ListEvent",
                        "max_len": 200,
                        "pos":"hi.fdt"
                    }
                }
            },
            "cmd": {
                "class": "UniqListHolderEvent",
                "max_len": 200,
                "primary_key":["command_line", "path"],
                "primary_key_creator": {
                    "class": "default_primary_key_creator",
                    "module": "container"
                },
                "childs": {
                    "command_line": {
                        "class": "UniqEvent",
                        "pos": "hi.cle"
                    },
                    "path": {
                        "class": "UniqEvent",
                        "pos": "hi.dst"
                    }
                }
            },
            "registry": {
                "class": "DictHolderEvent",
                "childs": {
                    "set": {
                        "class": "UniqListHolderEvent",
                        "max_len": 200,
                        "primary_key":["data", "path", "name"],
                        "primary_key_creator": {
                            "class": "default_primary_key_creator",
                            "module": "container"
                        },
                        "childs": {
                            "data": {
                                "class": "UniqEvent",
                                "pos": "hi.regd"
                            },
                            "path": {
                                "class": "UniqEvent",
                                "pos": "hi.regk"
                            },
                            "name": {
                                "class": "UniqEvent",
                                "pos": "hi.regv"
                            }
                        }
                    }
                }
            },
            "file_system": {
                "class": "DictHolderEvent",
                "childs": {
                    "mbr": {
                        "class": "UniqListHolderEvent",
                        "max_len": 200,
                        "primary_key":["data", "index", "size"],
                        "primary_key_creator": {
                            "class": "default_primary_key_creator",
                            "module": "container"
                        },
                        "childs": {
                            "data": {
                                "class": "UniqEvent",
                                "pos": "hi.mbr"
                            },
                            "index": {
                                "class": "UniqEvent",
                                "pos": "hi.mbrindex"
                            },
                            "size": {
                                "class": "UniqEvent",
                                "pos": "hi.mbrsize"
                            }
                        }
                    }
                }
            }
        }
}
        num = 100000
        print("num:%d" % num)

        #dump
        t = time.clock()
        for i in xrange(num):
            json.dumps(d)
        t1 = time.clock() - t
        print('json dump:%s' % (t1))

        t = time.clock()
        for i in xrange(num):
            simplejson.dumps(d)
        t1 = time.clock() - t
        print('simplejson dump:%s' % (t1))

        #t = time.clock()
        #for i in xrange(num):
        #    demjson.decode(d)
        #t1 = time.clock() - t
        #print('demjson dump:%s' % (t1))

        #load
        d = json.dumps(d)
        t = time.clock()
        for i in xrange(num):
            json.loads(d)
        t1 = time.clock() - t
        print('\njson load:%s' % (t1))

        t = time.clock()
        for i in xrange(num):
            simplejson.loads(d)
        t1 = time.clock() - t
        print('simplejson load:%s' % (t1))

        #t = time.clock()
        #for i in xrange(num):
        #    demjson.encode(d)
        #t1 = time.clock() - t
        #print('demjson load:%s' % (t1))
j = JsonTest()
j.test()
