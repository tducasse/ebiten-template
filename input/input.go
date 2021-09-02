package input

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var Input map[string][]ebiten.Key

func IsPressed(action string) bool {
	keys, ok := Input[action]
	if !ok {
		log.Println("Action " + action + " is not mapped")
		return false
	}
	for _, k := range keys {
		if ebiten.IsKeyPressed(k) {
			return true
		}
	}
	return false
}

func Init(keys map[string][]ebiten.Key) {
	Input = keys
}

func AddAction(action string, keys []ebiten.Key) {
	Input[action] = keys
}
