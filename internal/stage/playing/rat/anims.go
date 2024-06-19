/*
 *  Copyright (c) 2024 Juan Medina
 *
 *   All rights reserved. This software and related documentation are proprietary to Juan Medina.
 *
 *   This source code is for internal use only and may not be copied, modified, or distributed
 *   without the express written permission of Juan Medina. Any use of this software for any
 *   purpose other than its intended use is strictly prohibited and may result in severe civil
 *   and criminal penalties.
 *
 *   THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
 *   INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
 *   PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE
 *   FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
 *   OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
 *   DEALINGS IN THE SOFTWARE.
 */

package rat

const (
	IDLE_ANIM  = "idle"
	WALK_ANIM  = "walk"
	RUN_ANIM   = "run"
	FIGHT_ANIM = "fight"
	HURT_ANIM  = "hurt"
	DEAD_ANIM  = "dead"
	JUMP_ANIM  = "jump"
	HEAL_ANIM  = "heal"
)

var animationMap = map[string]animation{
	IDLE_ANIM: {
		pattern:       "rat_normal_idle_%02d",
		startFrame:    1,
		endFrame:      6,
		frameDuration: 100,
		loop:          true,
	},
	WALK_ANIM: {
		pattern:       "rat_normal_walk_%02d",
		startFrame:    1,
		endFrame:      8,
		frameDuration: 100,
		loop:          true,
	},
	RUN_ANIM: {
		pattern:       "rat_normal_run_%02d",
		startFrame:    1,
		endFrame:      8,
		frameDuration: 50,
		loop:          true,
	},
	FIGHT_ANIM: {
		pattern:       "rat_normal_fight_%02d",
		startFrame:    1,
		endFrame:      6,
		frameDuration: 100,
		loop:          false,
	},
	HURT_ANIM: {
		pattern:       "rat_normal_hurt_%02d",
		startFrame:    1,
		endFrame:      3,
		frameDuration: 100,
		loop:          false,
	},
	DEAD_ANIM: {
		pattern:       "rat_normal_dead_%02d",
		startFrame:    1,
		endFrame:      9,
		frameDuration: 100,
		loop:          false,
	},
	JUMP_ANIM: {
		pattern:       "rat_normal_jump_%02d",
		startFrame:    1,
		endFrame:      5,
		frameDuration: 100,
		loop:          false,
	},
	HEAL_ANIM: {
		pattern:       "rat_normal_heal_%02d",
		startFrame:    1,
		endFrame:      5,
		frameDuration: 150,
		loop:          false,
	},
}
