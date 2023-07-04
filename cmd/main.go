package main

import (
    . "github.com/boke0/att/pkg/entities"
    "github.com/boke0/att/pkg/state"
    "github.com/boke0/att/pkg/tree"
    "github.com/boke0/att/pkg/builder"
    "github.com/boke0/att/pkg/primitives"
)

func main() {
    entities := []AttAlice{}
    for i := 0; i<5; i++ {
        entities = append(entities, NewAttAlice())
    }
}
