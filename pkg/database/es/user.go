// description: sync_eth
// 
// @author: xwc1125
// @date: 2020/10/05
package es

import (
	"github.com/olivere/elastic/v7"
	"time"
)

type Tweet struct {
	User     string                `json:"user" es:"text"`
	Message  string                `json:"message" es:"text"`
	Retweets int                   `json:"retweets"  es:"text"`
	Image    string                `json:"image,omitempty"  es:"text"`
	Created  time.Time             `json:"created,omitempty"  es:"store"`
	Tags     []string              `json:"tags,omitempty"  es:"text"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"  es:"text"`
	A        float32
}
