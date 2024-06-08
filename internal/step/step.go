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

package step

type Step struct {
	initialValue float32
	endValue     float32
	timeToChange int
}

type LoopValue struct {
	steps        []Step
	currentStep  int
	currentValue float32
	elapsedTime  int
}

func (l *LoopValue) AddStep(initialValue float32, endValue float32, timeToChange int) {
	l.steps = append(l.steps, Step{
		initialValue: initialValue,
		endValue:     endValue,
		timeToChange: timeToChange,
	})
}

func (l *LoopValue) Update(elapsedTime int) bool {
	currentValue := l.currentValue
	current := l.steps[l.currentStep]
	l.elapsedTime += elapsedTime

	if l.elapsedTime >= current.timeToChange {
		l.currentValue = current.endValue
		if l.currentStep == len(l.steps)-1 {
			l.currentStep = 0
		} else {
			l.currentStep++
		}
		l.elapsedTime = 0
	} else {
		diff := current.endValue - current.initialValue
		l.currentValue = current.initialValue + diff*float32(l.elapsedTime)/float32(current.timeToChange)
	}

	return currentValue != l.currentValue
}

func (l LoopValue) GetValue() float32 {
	return l.currentValue
}

func NewLoopValue() LoopValue {
	return LoopValue{
		steps:       []Step{},
		currentStep: 0,
		elapsedTime: 0,
	}
}

func NewPingPongValue(initial, end float32, timeToChange int, middleTime int) LoopValue {
	value := NewLoopValue()
	value.AddStep(initial, end, timeToChange)
	value.AddStep(end, end, middleTime)
	value.AddStep(end, initial, timeToChange)
	return value
}
