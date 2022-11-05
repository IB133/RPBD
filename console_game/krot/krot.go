package krot

import (
	"fmt"
	"math/rand"
)

type Krot struct {
	noraLen int
	hp      int
	rep     int
	weight  float32
}

func New() *Krot {
	return &Krot{
		noraLen: 15,
		hp:      15,
		rep:     5,
		weight:  15,
	}
}

func (k *Krot) Dig(intense bool) {
	if intense {
		k.noraLen += 5
		k.hp -= 30
		return
	}
	k.noraLen += 2
	k.hp -= 10

}

func (k *Krot) Eat(green bool) {
	if green {
		if k.rep < 30 {
			k.hp -= 30
		} else if k.rep >= 30 {
			k.hp += 30
			k.weight += 30
		}
		return
	}
	k.hp += 10
	k.weight += 15
}

func (k *Krot) Fight(enemyWeight int) string {
	var chance = rand.Float32()
	switch enemyWeight {
	case 30:
		if k.weight/30 >= chance {
			k.rep += repChanges(30, k.weight)
			return "Вы победили"
		}
		k.hp -= 15
		return "Вы проиграли"
	case 50:
		if k.weight/50 >= chance {
			k.rep += repChanges(50, k.weight)
			return "Вы победили"
		}
		k.hp -= 30
		return "Вы проиграли"
	case 70:
		if k.weight/70 >= chance {
			k.rep += repChanges(70, k.weight)
			return "Вы победили"
		}
		k.hp -= 45
		return "Вы проиграли"
	}
	return ""
}

func (k *Krot) Sleep() {
	k.noraLen -= 2
	k.hp += 20
	k.rep -= 2
	k.weight -= 5
}

func (k *Krot) Stats() string {
	return fmt.Sprintf("Your stats:\nhp:%v\nrep:%v\nweight:%v\nHole length:%v\n", k.hp, k.rep, k.weight, k.noraLen)
}

func (k *Krot) IsWin() bool {
	return k.rep >= 100
}

func (k *Krot) IsLoose() bool {
	return k.hp <= 0 || k.noraLen <= 0 || k.rep <= 0 || k.weight <= 0
}

func repChanges(eWeight float32, kWeight float32) int {
	switch {
	case eWeight == 30 && kWeight < 50:
		return 10
	case eWeight == 30 && kWeight >= 50:
		return 5
	case eWeight == 50 && kWeight < 70:
		return 20
	case eWeight == 50 && kWeight >= 70:
		return 10
	case eWeight == 70:
		return 40
	}
	return 0
}
