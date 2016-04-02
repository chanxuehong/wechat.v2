// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package suite

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("item not found")

type TicketGetter interface {
	// 根据 suiteId 获取套件当前的 suiteTicket, 如果没有找到返回 ErrNotFound
	GetSuiteTicket(suiteId string) (ticket string, err error)

	// 没有实际意义, 接口标识
	TagEF8503CCFE9811E4959AA4DB30FED8E1()
}

var _ TicketGetter = (*TicketCache)(nil)

type TicketCache struct {
	rwmutex sync.RWMutex
	ticket  string
}

func (cache *TicketCache) TagEF8503CCFE9811E4959AA4DB30FED8E1() {}

func (cache *TicketCache) SetSuiteTicket(suiteId string, ticket string) (err error) {
	//if suiteId == "" {
	//	return errors.New("empty suiteId")
	//}
	if ticket == "" {
		return errors.New("empty ticket")
	}

	cache.rwmutex.Lock()
	cache.ticket = ticket
	cache.rwmutex.Unlock()
	return
}

func (cache *TicketCache) GetSuiteTicket(suiteId string) (ticket string, err error) {
	cache.rwmutex.RLock()
	ticket = cache.ticket
	if ticket == "" {
		err = ErrNotFound
	}
	cache.rwmutex.RUnlock()
	return
}

var _ TicketGetter = (*TicketCache2)(nil)

type TicketCache2 struct {
	rwmutex sync.RWMutex
	m       map[string]string
}

func NewTicketCache2() *TicketCache2 {
	return &TicketCache2{
		m: make(map[string]string),
	}
}

func (cache *TicketCache2) TagEF8503CCFE9811E4959AA4DB30FED8E1() {}

func (cache *TicketCache2) SetSuiteTicket(suiteId string, ticket string) (err error) {
	if suiteId == "" {
		return errors.New("empty suiteId")
	}
	if ticket == "" {
		return errors.New("empty ticket")
	}

	cache.rwmutex.Lock()
	cache.m[suiteId] = ticket
	cache.rwmutex.Unlock()
	return
}

func (cache *TicketCache2) GetSuiteTicket(suiteId string) (ticket string, err error) {
	cache.rwmutex.RLock()
	ticket = cache.m[suiteId]
	if ticket == "" {
		err = ErrNotFound
	}
	cache.rwmutex.RUnlock()
	return
}
