// @description wechat 是腾讯微信公众平台 api 的 golang 语言封装
// @link        https://github.com/chanxuehong/wechat for the canonical source repository
// @license     https://github.com/chanxuehong/wechat/blob/master/LICENSE
// @authors     chanxuehong(chanxuehong@gmail.com)

package media

import (
	"errors"
	"io"
	"net/url"
	"os"
	"path/filepath"

	"github.com/chanxuehong/wechat/mp"
)

// 上传多媒体图片
func (clt *Client) UploadImage(filepath string) (info *MediaInfo, err error) {
	return clt.uploadMedia(MediaTypeImage, filepath)
}

// 上传多媒体语音
func (clt *Client) UploadVoice(filepath string) (info *MediaInfo, err error) {
	return clt.uploadMedia(MediaTypeVoice, filepath)
}

// 上传多媒体视频
func (clt *Client) UploadVideo(filepath string) (info *MediaInfo, err error) {
	return clt.uploadMedia(MediaTypeVideo, filepath)
}

// 上传多媒体
func (clt *Client) uploadMedia(mediaType, _filepath string) (info *MediaInfo, err error) {
	file, err := os.Open(_filepath)
	if err != nil {
		return
	}
	defer file.Close()

	return clt.uploadMediaFromReader(mediaType, filepath.Base(_filepath), file)
}

// 上传多媒体图片
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadImageFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMediaFromReader(MediaTypeImage, filename, reader)
}

// 上传多媒体语音
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadVoiceFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMediaFromReader(MediaTypeVoice, filename, reader)
}

// 上传多媒体视频
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadVideoFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadMediaFromReader(MediaTypeVideo, filename, reader)
}

func (clt *Client) uploadMediaFromReader(mediaType, filename string, reader io.Reader) (info *MediaInfo, err error) {
	var result struct {
		mp.Error
		MediaInfo
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/upload?type=" +
		url.QueryEscape(mediaType) + "&access_token="
	fields := []mp.MultipartFormField{{
		ContentType: 0,
		FieldName:   "media",
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
	info = &result.MediaInfo
	return
}

// =============================================================================

// 上传多媒体缩略图
func (clt *Client) UploadThumb(_filepath string) (info *MediaInfo, err error) {
	file, err := os.Open(_filepath)
	if err != nil {
		return
	}
	defer file.Close()

	return clt.uploadThumbFromReader(filepath.Base(_filepath), file)
}

// 上传多媒体缩略图
//  NOTE: 参数 filename 不是文件路径, 是指定 multipart/form-data 里面文件名称
func (clt *Client) UploadThumbFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	if filename == "" {
		err = errors.New("empty filename")
		return
	}
	if reader == nil {
		err = errors.New("nil reader")
		return
	}
	return clt.uploadThumbFromReader(filename, reader)
}

func (clt *Client) uploadThumbFromReader(filename string, reader io.Reader) (info *MediaInfo, err error) {
	var result struct {
		mp.Error
		MediaType string `json:"type"`
		MediaId   string `json:"thumb_media_id"`
		CreatedAt int64  `json:"created_at"`
	}

	incompleteURL := "https://api.weixin.qq.com/cgi-bin/media/upload?type=thumb&access_token="
	fields := []mp.MultipartFormField{{
		ContentType: 0,
		FieldName:   "media",
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
	info = &MediaInfo{
		MediaType: result.MediaType,
		MediaId:   result.MediaId,
		CreatedAt: result.CreatedAt,
	}
	return
}
