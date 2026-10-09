package service

import "fmt"

var spellNames = map[int]string{
	1: "净化", 3: "虚弱", 4: "闪现", 6: "幽灵疾步", 7: "治疗术", 11: "惩戒",
	12: "传送", 13: "清晰术", 14: "引燃", 21: "屏障", 32: "标记", 30: "王者之路", 31: "王者之路",
}

// SpellName returns the Chinese name of a summoner spell ("" for 0 / unset).
func SpellName(id int) string {
	if id <= 0 {
		return ""
	}
	if n, ok := spellNames[id]; ok {
		return n
	}
	return fmt.Sprintf("技能%d", id)
}
