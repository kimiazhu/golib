// Description: download 包提供http相关工具。downloader.go提供下载工具。
// 目前只支持单线程下载，但提供下载各个阶段的控制，以及可以中断下载过程。
// 下载服务是同步的。
// TODO: 是否需要支持断点续传？
// Author: ZHU HAIHUA
// Since: 2016-03-09 16:32
package download

import (
	"fmt"
	log "github.com/kimiazhu/log4go"
	"github.com/xgsdk2/betatest/tako.lib/util"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	DOWNLOAD_FAILED    = -1
	CREATE_FILE_FAILED = -2
	SAVE_FAILED        = -3
	USER_CANCELED      = -4
	SERVER_ERROR       = -5
)

type DownloadError struct {
	code int
	msg  string
}

type DownloadProgress struct {
	finished int64
	total    int64
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("error %d: %s", e.code, e.msg)
}

type Downloader struct {
	Url        string
	SaveDir    string
	Override   bool
	onFinish   func(filepath string)
	onCancel   func()
	onError    func(error)
	onProgress func(finished int64, total int64)
	abort      bool
}

func NewDownloader(url, dir string, override bool) (downloader *Downloader, err error) {
	p, _ := filepath.Abs(dir)
	log.Debug("download dir is: %v", p)
	var fi os.FileInfo
	fi, err = os.Lstat(dir)
	if err != nil {
		p, _ := filepath.Abs(dir)
		log.Info("dir [%v] is not exist, try to create", p)
		if err = os.MkdirAll(dir, os.ModeDir|os.ModePerm); err != nil {
			log.Error("Cannot create dir [%v], error is: %v", dir, err)
			return nil, err
		}
		fi, err = os.Lstat(dir)
	}

	if !fi.IsDir() {
		err = fmt.Errorf("the path [%v] is not a dir", dir)
		log.Error("the path [%v] is not a dir", dir)
		return
	}

	downloader = &Downloader{
		Url:      url,
		SaveDir:  dir,
		Override: override,
	}

	return
}

func (d *Downloader) genFilename() (filename string) {
	defer log.Recover(func(err interface{}) string {
		// 如果生产文件中间产生任何异常，返回一个32位随机字符串。
		filename = util.RandStrN(32)
		return fmt.Sprintf("generate filename panic! error: %v", err)
	})
	if strings.Index(d.Url, "?") != -1 {
		url := d.Url[:strings.Index(d.Url, "?")]
		s := strings.Split(url, "/")
		filename = s[len(s)-1]
	} else {
		s := strings.Split(d.Url, "/")
		filename = s[len(s)-1]
	}

	if filename == "" {
		filename = util.RandStrN(32)
	}

	var basename, suffix string
	index := strings.LastIndex(filename, ".")
	if index < 0 {
		// dot not exists
		basename = filename
	} else if index == 0 {
		// remove the dot
		basename = string([]rune(filename)[1:])
		filename = basename
	} else if index > 0 {
		basename = string([]rune(filename)[:index])
		// suffix contains the dot(.)
		suffix = string([]rune(filename)[index:])
	}

	num := 1
	for {
		fullpath := filepath.Join(d.SaveDir, filename)
		if _, err := os.Stat(fullpath); err == nil {
			// file exists
			if d.Override {
				// remove exists
				os.Remove(fullpath)
				break
			} else {
				// find another name by append suffix
				filename = fmt.Sprintf("%s(%d)%s", basename, num, suffix)
				num++
			}
		} else {
			break
		}
	}
	return
}

// Start 会开启下载服务。
// callbacks可以不给任何参数，也可以接收最多三个参数。
// 其中第一个是onFinish func(string)会在下载完成时回调。
// 第二个是onError func(int, error)会在下载出错时回调
// 第三个是onProgress func(int64, int64)用于在下载过程中通知实时进度
// 这三个参数会覆盖Downloader对应的三个方法。
// 分别是OnFinish(f func())以及OnError(f func(int, error))
// 以及OnProgress(f func(int64, int64))。
func (d *Downloader) Start(callbacks ...interface{}) {
	switch len(callbacks) {
	case 1:
		d.onFinish = callbacks[0].(func(string))
	case 2:
		d.onFinish = callbacks[0].(func(string))
		d.onError = callbacks[1].(func(error))
	case 3:
		d.onFinish = callbacks[0].(func(string))
		d.onError = callbacks[1].(func(error))
		d.onProgress = callbacks[2].(func(int64, int64))
	}

	fn := d.genFilename()
	fullpath := filepath.Join(d.SaveDir, fn)

	if _, err := os.Stat(fullpath); err == nil && d.Override {
		if d.Override {
			// remove exists
			os.Remove(fullpath)
		} else {
			fullpath = filepath.Join()
		}
	}

	err := d.download(fullpath)
	if err != nil {
		re := os.Remove(fullpath)
		if re != nil {
			log.Error("remove file failed: %v", re)
		}
		msg := fmt.Errorf("download [%v] failed. error is: %v", d.Url, err)
		log.Error(msg)
		call(d.onError, msg)
	} else {
		log.Info("download url [%v] success", d.Url)
		call(d.onFinish, fullpath)
	}
}

func (d *Downloader) download(localpath string) error {
	file, err := os.Create(localpath)
	defer file.Close()

	if err != nil {
		e := &DownloadError{CREATE_FILE_FAILED, err.Error()}
		return e
	}

	req, err := http.NewRequest("GET", d.Url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		e := &DownloadError{DOWNLOAD_FAILED, err.Error()}
		return e
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e := &DownloadError{SERVER_ERROR, fmt.Sprintf("remote error: %v", resp.Status)}
		return e
	}

	var buf = make([]byte, 1024*32)
	var written int64
	for {
		if d.abort {
			log.Warn("user canceled download [%v]", d.Url)
			call(d.onCancel)
			return &DownloadError{USER_CANCELED, "user canceled"}
		}

		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			call(d.onProgress, written, resp.ContentLength)
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = &DownloadError{SAVE_FAILED, fmt.Sprintf("save [%v] failed", localpath)}
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}

	return err
}

func call(funcname interface{}, args ...interface{}) {
	switch f := funcname.(type) {
	case func():
		if f != nil {
			f()
		}
	case func(error):
		if f != nil {
			f(args[0].(error))
		}
	case func(int64, int64):
		if f != nil {
			f(args[0].(int64), args[1].(int64))
		}
	case func(string):
		if f != nil {
			f(args[0].(string))
		}
	default:
	}
}

// Cancel 取消当前下载
// onCancel可以指定一个参数：func() 用于在取消成功后进行回调。
// 此回调方法会覆盖OnCancel(f func())进行的设置
func (d *Downloader) Cancel(onCancel ...func()) {
	if len(onCancel) > 0 {
		d.onCancel = onCancel[0]
	}
	d.abort = true
}

func (d *Downloader) OnFinish(f func(string)) {
	d.onFinish = f
}

func (d *Downloader) OnCancel(f func()) {
	d.onCancel = f
}

func (d *Downloader) OnError(f func(error)) {
	d.onError = f
}

func (d *Downloader) OnProgress(f func(int64, int64)) {
	d.onProgress = f
}
