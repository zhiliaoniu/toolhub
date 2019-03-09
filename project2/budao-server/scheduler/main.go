package main

import (
	"master"
	"skel"
)

func main() {
	s := skel.New()
	{
		//regist master module
		m := master.New()
		s.RegisterModule("master", m)
	}
	err := s.Start()

	if err != nil {
		panic(err)
	}
	s.Wait()
}
