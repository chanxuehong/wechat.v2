// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package media

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/chanxuehong/wechat/mp"
)

// 上传到微信服务器的图片信息.
type ImageInfo struct {
	URL string `json:"url"`
}

// 上传图片到微信服务器, 给其他场景使用, 比如卡卷, POI.
func (clt *Client) UploadImagePermanent(imgPath string) (info ImageInfo, err error) {
	file, err := os.Open(imgPath)
	if err != nil {
		return
	}
	defer file.Close()

	return clt.uploadImagePermanentFromReader(filepath.Base(imgPath), file)
}

// 上传图片到微信服务器, 给其他场景使用, 比如卡卷, POI.
func (clt *Client) UploadImagePermanentFromReader(filename string, reader io.Reader) (info ImageInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}

	return clt.uploadImagePermanentFromReader(filename, reader)
}

func (clt *Client) uploadImagePermanentFromReader(filename string, reader io.Reader) (info ImageInfo, err error) {
	var result struct {
		mp.Error
		ImageInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/uploadimg?access_token="
	fields := []mp.MultipartFormField{{
		ContentType: 0,
		FieldName:   "buffer",
		FileName:    filename,
		Value:       reader,
	}}
	if err = ((*mp.Client)(clt)).PostMultipartForm(incompleteURL, fields, &result); err != nil {
		return
	}

	if result.ErrCode != mp.ErrCodeOK {
		err = &result.Error
		return
	}
	info = result.ImageInfo
	return
}
