package main

import (
	"bytes"
	"fmt"
    "time"
    "encoding/json"

	. "github.com/boke0/att/pkg/entities"
	. "github.com/boke0/att/pkg/messages"
	"github.com/boke0/att/pkg/primitives"
)

func main() {
    //r := 100
    r := 20
    n := 50

    entities := []AttAlice{}
    for i := 0; i<n; i++ {
        entities = append(entities, NewAttAlice())
    }

    var (
        message AttMessage
        sent []byte
    )

    fmt.Printf("sent_data_size, max_receive_time\n",)
    var (
        sent_data_size int
        max_receive_time int64
    )
    for i := 0; i<r; i++ {
        /** # 送信 **/
        /** ## 準備 **/
        // i番目のEntity視点で、AliceとBobに分ける
        alice := entities[i%n]
        bobs := map[string]AttBob{}
        for _, entity := range entities {
            if entity.Id != alice.Id {
                bobs[entity.Id] = entity.Bob()
            }
        }

        if(i == 0) {
            // i番目のEntityが初期化する
            message = alice.Initialize(bobs)
        } else {
            // i番目のEntityがメッセージを送る
            sent = primitives.RandomByte()
            message = alice.Send(sent)
        }
        data, _ := json.Marshal(message)
        sent_data_size = len(data)
        entities[i%n] = alice

        /** # 受信 **/
        for j := 0; j<n; j++ {
            if(j == i%n){
                continue
            }
            /** ## 準備 **/
            // i番目のEntity視点で、AliceとBobに分ける
            alice := entities[j]
            bobs := map[string]AttBob{}
            for _, entity := range entities {
                if entity.Id != alice.Id {
                    bobs[entity.Id] = entity.Bob()
                }
            }
            t := time.Now()
            //time.Sleep(time.Millisecond * 200)
            if message.InitializeMessage != nil {
                alice.Receive(message, bobs)
            }else{
                received := alice.Receive(message, bobs)
                if !bytes.Equal(received, sent) {
                    panic("invalid message")
                }
            }
            if max_receive_time < time.Since(t).Milliseconds() {
                max_receive_time = time.Since(t).Milliseconds()
            }
            entities[j] = alice
        }
        fmt.Printf("%d, %d\n", sent_data_size, max_receive_time)
    }
}
