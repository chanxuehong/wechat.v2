// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com), magicshui(shuiyuzhe@gmail.com), Harry Rong(harrykobe@gmail.com)

package material

import (
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chanxuehong/wechat/mp"
)

type ImageInfo struct {
	PicURL string `json:"pic_url"`
}

func Add(clt *mp.Client, imagePath, _type string) (info ImageInfo, err error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return
	}
	defer file.Close()

	return addFromReader(clt, filepath.Base(imagePath), file, _type)
}

func AddFromReader(clt *mp.Client, filename string, reader io.Reader, _type string) (info ImageInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}

	return addFromReader(clt, filename, reader, _type)
}

func addFromReader(clt *mp.Client, filename string, reader io.Reader, _type string) (info ImageInfo, err error) {
	var result struct {
		mp.Error
		ImageInfo `json:"data"`
	}

	var incompleteURL string
	if _type != "" {
		incompleteURL = "https://api.weixin.qq.com/shakearound/material/add?type=" + url.QueryEscape(_type) +
			"&access_token="
	} else {
		incompleteURL = "https://api.weixin.qq.com/shakearound/material/add?access_token="
	}
	fields := []mp.MultipartFormField{{
		ContentType: 0,
		FieldName:   "media",
		FileName:    filename,
		Value:       reader,
	}}
	if err = clt.PostMultipartForm(incompleteURL, fields, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = result.ImageInfo
	return
}
