package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-contrib/sse"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		if r.Header["Authorization"][0] == "token 123" {
			return true
		}
		return false
	},
}


func main() {

	router := gin.Default()

	router.GET("/someGet", getting)
	router.GET("/test", MsgSendHandler)
	router.Run()
}


func CreateFile(fileName string){
	f,err := os.Create( fileName )

	defer f.Close()

	if err !=nil {

		fmt.Println( err.Error() )

	} else {

		for i := 0; i<5000000 ; i++ {
			_,err=f.Write([]byte("这条信息的编号是"+strconv.Itoa(i)+";"))

			if err != nil {
				fmt.Println( err.Error() )
				return
			}
		}

	}
}

func ReadFile(fileName string) *bytes.Buffer {
	f, err := os.OpenFile(fileName, os.O_RDONLY,0600)

	defer f.Close()

	if err == nil {

		contentByte,err :=ioutil.ReadAll(f)
		if err != nil {
			return nil
		}
		return bytes.NewBuffer(contentByte)
	}
	return nil
}

func MsgSendHandler(c *gin.Context) {

	var x = "这条信息的编号是0"
	var y = []byte(x)
	//z,_ := json.Marshal(y)


	c.Render(-1,render.Data{"application/auto",y})
	//c.Writer.Header()["Content-Type"] = []string{"text/event-stream"}
	return
	c.Render(-1, sse.Event{
		Event: "test",
		Data:  y,
	})
}

func getting(c *gin.Context){
	//c.JSON(http.StatusOK,ApprovalCallBackInfo{"c","t","m","t","u",
	//	ApprovalEvent{"1","key","eventtype","instancecode1","APPROVED","approvalcode1"}})
	var err error
	backInfo := &CallBackInfo{} //回调变量绑定

	if err = c.Bind(backInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "请求参数错误",
		})
		return
	}

	//b := ReadFile("dashboard_export_20200815_20200831.xlsx")
	b := ReadFile("测试文件")
	if b == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "文件打开错误",
		})
		return
	}
	fmt.Printf("%v", b.Len())
	if err := SseSliceSend(c, b, 1+b.Len()/10,2, 1); err != nil {
		fmt.Printf("client gone")
		return
	}

}

func SseSliceSend(c *gin.Context, bufferPtr *bytes.Buffer, sliceLen, chanCap int, waitSeconds time.Duration) error {

	if bufferPtr == nil {
		fmt.Printf("bufferPtr is nil")
		return errors.New("bufferPtr is nil")
	}
	if sliceLen <= 0 || chanCap <= 0 {
		fmt.Printf("param error")
		return errors.New("param error")
	}
	chanStream := make(chan *[]byte, chanCap)
	//chanStream := make(chan string, chanCap)
	go func(bPtr *bytes.Buffer, sLen int, seconds time.Duration) {
		defer close(chanStream)
		for bPtr.Len() > 0 {
			var n int
			if bPtr.Len() > sLen {
				n = sLen
			} else {
				n = bPtr.Len()
			}
			buffSlice := make([]byte, n)
			bPtr.Read(buffSlice)
			chanStream <- &buffSlice
			time.Sleep(time.Second*seconds)
		}
		chanStream <- nil
	}(bufferPtr, sliceLen, waitSeconds)
	//c.Writer.Header().Add("Connection","keep-alive")
	//c.Writer.Header().Add("Keep-Alive","timeout=60")
	c.Writer.Header().Set("File-Name","test_file.xlsx")
	if c.Stream(func(w io.Writer) bool {
		if msgPtr, ok := <-chanStream; ok {
			if msgPtr != nil {
				//c.SSEvent("success", msgPtr)
				c.Render(-1,render.Data{})
				c.Render(-1,render.Data{"application/auto",*msgPtr})
				return true
			}
			//c.SSEvent("finish", "")
		}
		return false
	}) {
		fmt.Printf("c.Stream SSEvent client gone")
		return errors.New("c.Stream SSEvent client gone")
	}
	return nil
}