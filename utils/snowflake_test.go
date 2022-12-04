package utils

import (
	"PInfo-server/log"
	"testing"
)

func TestNextVal(t *testing.T) {
	log.InitLogger("./test.log", 1024, 1, 1, -1, 2)

	sn, _ := NewSnowflake(1)
	i := 0
	for i < 65536 {
		i++
		sn.NextVal()
		//fmt.Println(sn.NextVal())
	}
}
