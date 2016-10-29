package main
import(
	"testing"
)
func TestGetRequiredXp(t *testing.T){
	lf := LevelFunction{IncrementExperience: 2, BaseExperience: 5}
	if tRequired := lf.GetRequiredXp(3); tRequired != (5 + 2*3){
		t.Error("expected", (5 + 2*3), "got", tRequired)
	}
}
func TestCanLevel(t *testing.T){
	lf := LevelFunction{IncrementExperience: 2, BaseExperience: 5}
	if tCanLevel := lf.canLevel(3, 11); tCanLevel != true{
		t.Error("expected", true, "got", tCanLevel)
	}
	if tCanLevel := lf.canLevel(3, 12); tCanLevel != true{
		t.Error("expected", true, "got", tCanLevel)
	}
	if tCanLevel := lf.canLevel(3, 10); tCanLevel != false{
		t.Error("expected", true, "got", tCanLevel)
	}
	if tCanLevel := lf.canLevel(0, 5); tCanLevel != true{
		t.Error("expected", true, "got", tCanLevel)
	}
	if tCanLevel := lf.canLevel(0, 4); tCanLevel != false{
		t.Error("expected", false, "got", tCanLevel)
	}
}