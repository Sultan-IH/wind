package datahandler

import (
	"strconv"
	"time"

	"bitbucket.org/emotech/common/golang/logs"

	"github.com/Sultan-IH/wind/plug"
)

const threshold = 50

// goroutine to control ventilator
func controlVentilatorGrad(dataChannel DataChannel, plug plug.Plug) {
	// calculate gradient of channel in windows
	logs.Printf("[ID:%s] control ventilator started", plug.Alias)
	previous := 0
	for data := range dataChannel {
		current, err := strconv.Atoi(string(data))
		if err != nil {
			logs.Printf("error in control ventilator: atoi threw and error: %v", err)
			continue
		}
		delta := current - previous
		logs.Printf("[%s] delta is %d", plug.ID, delta)
		if delta > threshold {
			logs.Printf("[ID:%s] turning plug on", plug.Alias)
			plug.TurnON()
		} else {
			logs.Printf("[ID:%s] turning plug off", plug.Alias)
			plug.TurnOFF()
		}
	}
	logs.Printf("[ID:%s] control ventilator finished", plug.Alias)

}

var (
	samplePeriod       = time.Second * 3 / 2
	maxWindSensorValue = 1024
	powerUpRate        = samplePeriod
	powerDownRate      = samplePeriod // could be different to samplePeriod
)

func controlVentilatorPWM(dataChannel DataChannel, vplug plug.Plug) {
	// every this
	buffer := []int{}
	lastSample := time.Now()
	logs.Printf("[CONTROL] control the ventilator goroutine started")
	for v := range dataChannel {
		i, err := strconv.Atoi(string(v))
		if err != nil {
			logs.Printf("error in control ventilator: atoi threw and error: %v", err)
			continue
		}
		buffer = append(buffer, i)

		if time.Now().After(lastSample.Add(samplePeriod)) {
			logs.Printf("[CONTROL] starting PWM")

			// calculate the average
			avg := average(buffer)
			logs.Printf("[CONTROL] average is: %d", avg)
			// pwm this shit
			percentage := float32(avg) / float32(maxWindSensorValue)
			logs.Printf("[CONTROL] percentage is :%f; avg: %f; maxval: %f;", percentage, float32(avg), float32(maxWindSensorValue))

			remainder := percentage - vplug.VentilatorState
			logs.Printf("[CONTROL] remainder is: %f", remainder)

			timeON := time.Duration(remainder) * powerUpRate
			logs.Printf("[CONTROL] time on is: %v", timeON)

			go pwmVentilator(timeON, vplug)

			timeOFF := samplePeriod - timeON
			endState := 1 - (timeOFF / powerDownRate)
			vplug.VentilatorState = float32(endState)
			lastSample = time.Now()
		}

	}
}

func pwmVentilator(timeON time.Duration, vplug plug.Plug) {
	// for the next samplePeriod seconds
	if timeON == time.Duration(0) {
		return
	}
	vplug.TurnON()
	timer := time.NewTimer(timeON)
	<-timer.C
	vplug.TurnOFF()
	// change plug state

}
func average(ints []int) int {
	s := 0
	for _, v := range ints {
		s += v
	}
	return s / len(ints)
}
