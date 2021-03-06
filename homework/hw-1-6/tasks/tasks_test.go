package tasks

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	// Make randomization
	rand.Seed(time.Now().UnixNano())
}

func emulateActivity(taskNumber, delaySec int, error bool) {
	time.Sleep(time.Duration(delaySec) * time.Second)
	if !error {
		log.Printf("Task %d fired in %d seconds", taskNumber, delaySec)
	}
}

var tasks1 = []func() error{
	func() error {
		emulateActivity(1, rand.Intn(5), false)
		return nil
	},
	func() error {
		emulateActivity(2, rand.Intn(5), false)
		return nil
	},
	func() error {
		emulateActivity(3, rand.Intn(5), false)
		return nil
	},
	func() error {
		emulateActivity(4, rand.Intn(5), false)
		return nil
	},
}

var tasks2 = []func() error{
	func() error {
		emulateActivity(1, 1, true)
		return errors.New("error in task 1")
	},
	func() error {
		emulateActivity(2, 1, false)
		return nil
	},
	func() error {
		emulateActivity(3, 1, true)
		return errors.New("error in task 3")
	},
	func() error {
		emulateActivity(4, 2, true)
		return errors.New("error in task 4")
	},
	func() error {
		emulateActivity(5, 3, false)
		return nil
	},
}

var tasks3 []func() error

func tasks3Builder(number int) {
	for i := 0; i < number; i++ {
		i := i
		tasks3 = append(tasks3, func() error {
			errFlag := rand.Intn(2) == 0
			delaySec := rand.Intn(5)
			emulateActivity(i, delaySec, errFlag)

			// Emulate result || error
			var err error = nil
			if errFlag {
				err = errors.New(fmt.Sprintf("error in task %d", i))
			}
			return err
		})
	}
}

const (
	concurrentTaskNumber = 4
	errorsLimit          = 2
)

func TestRun(t *testing.T) {
	var err error
	log.Println("************ Test 1 ************")
	if err = Run(tasks1, concurrentTaskNumber, errorsLimit); err != nil {
		log.Printf("An error occurred: %s", err.Error())
	}
	assert.Nil(t, err, "There should be no errors in this test")

	log.Println()
	log.Println("************ Test 2 ************")
	if err = Run(tasks2, concurrentTaskNumber, errorsLimit); err != nil {
		log.Printf("An error occurred: %s", err.Error())
	}
	assert.NotNil(t, err, "There should be error in this test")

	log.Println()
	log.Println("************ Test 3 ************")
	tasks3Builder(20)
	if err = Run(tasks3, 10, 4); err != nil {
		log.Printf("An error occurred: %s", err.Error())
	}
	assert.NotNil(t, err, "There should be error in this test")
}
