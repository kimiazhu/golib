// Description: download
// Author: ZHU HAIHUA
// Since: 2016-03-10 11:35
package download

import (
	log "github.com/kimiazhu/log4go"
	"github.com/xgsdk2/betatest/tako.lib/mgox"
	"testing"
	"time"
	"fmt"
)

func init() {
	cfg := `
	<logging>
	    <filter enabled="true">
		    <tag>stdout</tag>
		    <type>console</type>
		    <level>DEBUG</level>
		    <exclude>github.com/xgsdk2/betatest/tako.lib/mgox</exclude>
	    </filter>
	</logging>
	`
	log.Setup([]byte(cfg))

	mgox.Config("42.62.96.68", "tako", "tako", "tako")
}

func TestDownload(t *testing.T) {
	//url := "http://mirrors.ustc.edu.cn/opensuse/distribution/12.3/iso/openSUSE-12.3-GNOME-Live-i686.iso"
	url := "http://localhost:8000/openSUSE.iso"

	d, e := NewDownloader(url, "C:/Users/KC/Desktop", false)
	if e != nil {
		log.Error("download failed: %v", e)
	} else {
		d.OnFinish(func(fullpath string){
			fmt.Println("success download file: " + fullpath)
		})
		d.OnError(func(err error){
			fmt.Printf("Error occur: %v\n", err)
		})
		d.OnProgress(func(finish, total int64){
			fmt.Printf("finished: %.03f%%\n", float64(finish)/float64(total) * 100)
		})
		go d.Start()
		time.Sleep(10 * time.Second)
		d.Cancel(func(){
			fmt.Println("Canceled by user after 10 second")
		})
	}

	time.Sleep(15 * time.Second)
}
