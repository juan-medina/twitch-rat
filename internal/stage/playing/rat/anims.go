/*
 * Copyright (c) 2023 Juan Antonio Medina Iglesias
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 */

package rat

const (
	IDLE_ANIM  = "idle"
	WALK_ANIM  = "walk"
	RUN_ANIM   = "run"
	FIGHT_ANIM = "fight"
	HURT_ANIM  = "hurt"
	DEAD_ANIM  = "dead"
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
}
