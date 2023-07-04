package main

import (
    . "github.com/boke0/att/pkg/entities"
    "github.com/boke0/att/pkg/state"
    "github.com/boke0/att/pkg/tree"
    "github.com/boke0/att/pkg/builder"
    "github.com/boke0/att/pkg/primitives"
)

func main() {
    r := 100
    n := 5

    entities := []AttAlice{}
    for i := 0; i<n; i++ {
        entities = append(entities, NewAttAlice())
    }

    var message AttMessage

    for i := 0; i<r; i++ {
        /** # 送信 **/
        /** ## 準備 **/
        // i番目のEntity視点で、AliceとBobに分ける
        alice := entities[i]
        bobs := map[string]AttBob{}
        for _, entity := range entities {
            if entity.Id != alice.Id {
                bobs[entity.Id] = entity.ToBob()
            }
        }
        // i番目のEntity視点で、Treeを作る

        /** # 受信 **/
        for j := 0; j<n; i++ {
            /** ## 準備 **/
            // i番目のEntity視点で、AliceとBobに分ける
            alice := entities[i]
            bobs := map[string]AttBob{}
            for _, entity := range entities {
                if entity.Id != alice.Id {
                    bobs[entity.Id] = entity.ToBob()
                }
            }
        }
    }
}
