package main

import (
	"fmt"
	"time"
	"net"
	"message"
)

func main() {
	fmt.Println("test client...")
	time.Sleep(3 * time.Second)

	con,err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println(err)
	}

	for {
		testMsg := message.NewLoginMsg(1, "test:::")
		propertyList := testMsg.GetPropertyList()
		fmt.Println(propertyList)
		bts, _ := message.WriteByte(propertyList)
		fmt.Println(bts)
		_,err := con.Write(bts)
		if err != nil {
			fmt.Println(err)
			break
		}
		buf := make([]byte, 512)

		count,err := con.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		fmt.Printf("\nserver call back: %s, count = %d\n", buf, count)

		time.Sleep(1 * time.Second)
	}

	fmt.Print("客户端正常退出。")
}
