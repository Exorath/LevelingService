package main

type LevelFunction struct {
	BaseExperience int
	IncrementExperience int
}
func (lf LevelFunction) GetRequiredXp(lvl int) int {
	return lf.BaseExperience + lvl * lf.IncrementExperience
}

func (lf LevelFunction) canLevel(lvl int, exp int) bool {
	return exp >= lf.GetRequiredXp(lvl)
}