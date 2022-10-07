package krot

import (
	"fmt"
	"math/rand"
)

type Krot struct {
	Nora_len int
	Hp       int
	Rep      int
	Weight   float32
}

func (k *Krot) Dig(intense bool) {
	if intense {
		k.Nora_len += 5
		k.Hp -= 30
	} else {
		k.Nora_len += 2
		k.Hp -= 10
	}

}

func (k *Krot) Eat(green bool) {
	if green {
		if k.Rep < 30 {
			k.Hp -= 30
		} else if k.Rep >= 30 {
			k.Hp += 30
			k.Weight += 30
		}
	} else {
		k.Hp += 10
		k.Weight += 15
	}
}

func (k *Krot) Fight(enemy_Weight int) string {
	var chance = rand.Float32()
	switch enemy_Weight {
	case 30:
		if k.Weight/30 >= chance {
			k.Rep += Rep_changes(30, k.Weight)
			return "Вы победили"
		} else {
			k.Hp -= 15
			return "Вы проиграли"
		}
	case 50:
		if k.Weight/50 >= chance {
			k.Rep += Rep_changes(50, k.Weight)
			return "Вы победили"
		} else {
			k.Hp -= 30
			return "Вы проиграли"
		}
	case 70:
		if k.Weight/70 >= chance {
			k.Rep += Rep_changes(70, k.Weight)
			return "Вы победили"
		} else {
			k.Hp -= 45
			return "Вы проиграли"
		}
	}
	return ""
}

func (k *Krot) Sleep() {
	k.Nora_len -= 2
	k.Hp += 20
	k.Rep -= 2
	k.Weight -= 5
}

func (k *Krot) Stats() string {
	return fmt.Sprintf("Your stats:\nHp:%v\nRep:%v\nWeight:%v\nHole length:%v", k.Hp, k.Rep, k.Weight, k.Nora_len)
}

func Rep_changes(e_Weight float32, k_Weight float32) int {
	switch {
	case e_Weight == 30 && k_Weight < 50:
		return 10
	case e_Weight == 30 && k_Weight >= 50:
		return 5
	case e_Weight == 50 && k_Weight < 70:
		return 20
	case e_Weight == 50 && k_Weight >= 70:
		return 10
	case e_Weight == 70:
		return 40

	}
	return 0
}
