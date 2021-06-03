// description: sync_eth 
// 
// @author: xwc1125
// @date: 2020/10/05
package es

import (
	"fmt"
	"testing"
)

func TestTweet_Mapping(t *testing.T) {
	mapping := Mapping(Tweet{})
	fmt.Println(mapping)
}
