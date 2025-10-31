package lib

import "fmt"

type Stateway struct{}

func (s *Stateway) Start() {
	fmt.Println("Stateway started")
}
