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

package step

type Action int

const (
	NoAction Action = iota
	Next
	Loop
	End
)

type Value interface {
	GetValue() float32
	Update(elapsedTime int) bool
	Reset()
	IsAtEnd() bool
}

type step struct {
	initialValue float32
	endValue     float32
	timeToChange int
	action       Action
}

type stepImpl struct {
	steps        []step
	currentStep  int
	currentValue float32
	elapsedTime  int
	atEnd        bool
}

func (l *stepImpl) AddStep(initialValue float32, endValue float32, timeToChange int, action Action) {
	l.steps = append(l.steps, step{
		initialValue: initialValue,
		endValue:     endValue,
		timeToChange: timeToChange,
		action:       action,
	})
}

func (l *stepImpl) Update(elapsedTime int) bool {
	currentValue := l.currentValue
	current := l.steps[l.currentStep]
	l.elapsedTime += elapsedTime

	if l.elapsedTime >= current.timeToChange {
		l.currentValue = current.endValue
		switch current.action {
		case Next:
			if l.currentStep == len(l.steps)-1 {
				l.currentStep = 0
			} else {
				l.currentStep++
			}
		case Loop:
			l.currentStep = 0
		case End:
			l.atEnd = true
			return true
		}
		l.elapsedTime = 0
	} else {
		diff := current.endValue - current.initialValue
		l.currentValue = current.initialValue + diff*float32(l.elapsedTime)/float32(current.timeToChange)
	}

	return currentValue != l.currentValue
}

func (l stepImpl) GetValue() float32 {
	return l.currentValue
}
func (l *stepImpl) Reset() {
	l.currentStep = 0
	l.currentValue = l.steps[0].initialValue
	l.elapsedTime = 0
	l.atEnd = false
}

func (l stepImpl) IsAtEnd() bool {
	return l.atEnd
}

func newStep() stepImpl {
	return stepImpl{
		steps:       []step{},
		currentStep: 0,
		elapsedTime: 0,
	}
}

func NewPingPongValue(initial, end float32, timeToChange int, middleTime int) Value {
	value := newStep()
	value.currentValue = initial
	value.AddStep(initial, end, timeToChange, Next)
	value.AddStep(end, end, middleTime, Next)
	value.AddStep(end, initial, timeToChange, Loop)
	return &value
}

func NewLoopValue(initial, end float32, timeToChange int) Value {
	value := newStep()
	value.currentValue = initial
	value.AddStep(initial, end, timeToChange, Loop)
	return &value
}

func NewFromToPauseValue(initial, end float32, timeToChange int, timeAtEnd int) Value {
	value := newStep()
	value.AddStep(initial, end, timeToChange, Next)
	value.AddStep(end, end, timeAtEnd, End)
	return &value
}

func NewFromMiddleToPauseValue(initial, middle float32, end float32, timeToMiddle int, timeToEnd int, timeAtEnd int) Value {
	value := newStep()
	value.AddStep(initial, middle, timeToMiddle, Next)
	value.AddStep(middle, end, timeToEnd, Next)
	value.AddStep(end, end, timeAtEnd, End)
	return &value
}
