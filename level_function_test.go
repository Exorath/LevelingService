/*
 *    Copyright 2016 Exorath
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */
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