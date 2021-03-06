/// Author: Bernhard Tittelbach, btittelbach@github  (c) 2014
package bbhw

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func checkSysfsVersusMMapGPIO(gpionum uint, t *testing.T) {
	chipid, gpioid := calcGPIOAddrFromLinuxGPIONum(gpionum)
	t.Logf("Testing sysfs:gpio/gpio%d chip:gpio%d[%d]", gpionum, chipid, gpioid)
	fg := NewMMappedGPIO(gpionum, OUT)
	sg := NewSysfsGPIOOrPanic(gpionum, OUT)
	defer sg.Close()
	defer fg.Close()
	//Step(usr3, 20, time.Duration(200)*time.Millisecond, nil)

	// Test Direction
	d1, err1 := fg.CheckDirection()
	if err1 != nil {
		t.Error(err1.Error())
	}
	if d1 != OUT {
		t.Error("fg.CheckDirection != OUT")
	}
	d2, err2 := sg.CheckDirection()
	if err2 != nil {
		t.Error(err2.Error())
	}
	if d2 != OUT {
		t.Error("sg.CheckDirection != OUT")
	}

	//Test Slow
	sg.SetState(true)
	if GetStateOrPanic(sg) != true {
		t.Error("0: sg.GetState() != sg.SetState()")
	}
	sg.SetState(false)
	if GetStateOrPanic(sg) != false {
		t.Error("1: sg.GetState() != sg.SetState()")
	}

	// Test Fast
	fg.SetState(true)
	if GetStateOrPanic(fg) != true {
		t.Error("0: fg.GetState() != fg.SetState()")
	}
	fg.SetState(false)
	if GetStateOrPanic(fg) != false {
		t.Error("1: fg.GetState() != fg.SetState()")
	}

	// Test SysFS vs MMapped
	fg.SetState(false)
	if GetStateOrPanic(sg) != false {
		t.Error("1: sg.GetState() != fg.SetState()")
	}
	fg.SetState(true)
	if GetStateOrPanic(sg) != true {
		t.Error("2: sg.GetState() != fg.SetState()")
	}
	sg.SetState(false)
	if GetStateOrPanic(fg) != false {
		t.Error("3: sg.GetState() != fg.SetState()")
	}
	sg.SetState(true)
	if GetStateOrPanic(fg) != true {
		t.Error("4: sg.GetState() != fg.SetState()")
	}
	sg.SetState(false)
}

func Test_FastGPIO(t *testing.T) {
	if !verifyAddrIsTIOmap4(omap4_gpio0_offset_) {
		t.Logf("test only works on BeagleBone")
		return
	}
	// fg := NewMMappedGPIO(67, OUT)
	// sg := NewSysfsGPIOOrPanic(67, OUT)
	checkSysfsVersusMMapGPIO(2, t)
	checkSysfsVersusMMapGPIO(3, t)
	checkSysfsVersusMMapGPIO(4, t)
	checkSysfsVersusMMapGPIO(5, t)
	checkSysfsVersusMMapGPIO(50, t)
	checkSysfsVersusMMapGPIO(51, t)
	checkSysfsVersusMMapGPIO(61, t)
	checkSysfsVersusMMapGPIO(67, t)
	checkSysfsVersusMMapGPIO(80, t)
	checkSysfsVersusMMapGPIO(81, t)
	checkSysfsVersusMMapGPIO(88, t)
	checkSysfsVersusMMapGPIO(117, t)
}

func checkSysfsVersusActiveLowMMapGPIO(gpionum uint, t *testing.T) {
	chipid, gpioid := calcGPIOAddrFromLinuxGPIONum(gpionum)
	t.Logf("Testing sysfs:gpio/gpio%d chip:gpio%d[%d]", gpionum, chipid, gpioid)
	fg := NewMMappedGPIO(gpionum, OUT)
	sg := NewSysfsGPIOOrPanic(gpionum, OUT)
	defer sg.Close()
	defer fg.Close()

	// Test SysFS vs MMapped
	fg.SetState(false)
	if GetStateOrPanic(sg) != false {
		t.Error("1: sg.GetState() != fg.SetState()")
	}
	fg.SetActiveLow(true)
	if GetStateOrPanic(sg) != true {
		t.Error("2: sg.GetState() != ! fg.SetState()")
	}
	sg.SetState(false)
	if GetStateOrPanic(fg) != true {
		t.Error("3: sg.GetState() != ! fg.SetState()")
	}
	fg.SetState(true)
	if GetStateOrPanic(sg) != false {
		t.Error("4: sg.GetState() != ! fg.SetState()")
	}
	sg.SetState(false)
}

func Test_FastActiveLowGPIO(t *testing.T) {
	if !verifyAddrIsTIOmap4(omap4_gpio0_offset_) {
		t.Logf("test only works on BeagleBone")
		return
	}
	checkSysfsVersusActiveLowMMapGPIO(3, t)
	checkSysfsVersusActiveLowMMapGPIO(4, t)
	checkSysfsVersusActiveLowMMapGPIO(5, t)
	checkSysfsVersusActiveLowMMapGPIO(50, t)
	checkSysfsVersusActiveLowMMapGPIO(51, t)
	checkSysfsVersusActiveLowMMapGPIO(61, t)
	checkSysfsVersusActiveLowMMapGPIO(67, t)
	checkSysfsVersusActiveLowMMapGPIO(80, t)
	checkSysfsVersusActiveLowMMapGPIO(81, t)
	checkSysfsVersusActiveLowMMapGPIO(88, t)
	checkSysfsVersusActiveLowMMapGPIO(117, t)
}

func Test_SysfsGPIOwCable(t *testing.T) {
	if !verifyAddrIsTIOmap4(omap4_gpio0_offset_) {
		t.Logf("test only works on BeagleBone")
		return
	}
	outg := NewSysfsGPIOOrPanic(67, OUT) //P8_8
	ing := NewSysfsGPIOOrPanic(66, IN)   //P8_7
	defer outg.Close()
	defer ing.Close()
	outg.SetState(false)
	if GetStateOrPanic(ing) != false {
		fmt.Println("For this test, please connect Pin P8_7 to P8_8")
		t.Error("1: ing.GetState() != outg.SetState()")
	}
	outg.SetState(true)
	time.Sleep(time.Duration(10) * time.Millisecond)
	if GetStateOrPanic(ing) != true {
		t.Error("2: ing.GetState() != outg.SetState()")
	}
	outg.SetState(false)
	//Step(outg, 20, time.Duration(200)*time.Millisecond, nil)
}

func Test_MmappedGPIOwCable(t *testing.T) {
	if !verifyAddrIsTIOmap4(omap4_gpio0_offset_) {
		t.Logf("test only works on BeagleBone")
		return
	}
	outg := NewMMappedGPIO(67, OUT) //P8_8
	ing := NewMMappedGPIO(66, IN)   //P8_7

	// Test Direction
	d1, err1 := outg.CheckDirection()
	if err1 != nil {
		t.Error(err1.Error())
	}
	if d1 != OUT {
		t.Error("outg.CheckDirection != OUT")
	}
	d2, err2 := ing.CheckDirection()
	if err2 != nil {
		t.Error(err2.Error())
	}
	if d2 != IN {
		t.Error("ing.CheckDirection != IN")
	}

	outg.SetState(false)
	if GetStateOrPanic(ing) != false {
		fmt.Println("For this test, please connect Pin P8_7 to P8_8")
		t.Error("1: ing.GetState() != outg.SetState()")
	}
	outg.SetState(true)
	if GetStateOrPanic(ing) != true {
		t.Error("2: ing.GetState() != outg.SetState()")
	}
	outg.SetState(false)
	//Step(outg, 20, time.Duration(200)*time.Millisecond, nil)
}

func Test_MmappedGPIOwCableInGoroutines(t *testing.T) {
	if !verifyAddrIsTIOmap4(omap4_gpio0_offset_) {
		t.Logf("test only works on BeagleBone")
		return
	}
	outg := NewMMappedGPIO(67, OUT) //P8_8
	outslow := NewSysfsGPIOOrPanic(67, OUT)
	ing := NewMMappedGPIO(66, IN) //P8_7

	go outg.SetState(false)
	time.Sleep(10 * time.Millisecond)
	if GetStateOrPanic(outslow) != false {
		fmt.Println("For this test, please connect Pin P8_7 to P8_8")
		t.Error("1: outslow.GetState() != outg.SetState()")
	}
	if GetStateOrPanic(ing) != false {
		fmt.Println("For this test, please connect Pin P8_7 to P8_8")
		t.Error("1: ing.GetState() != outg.SetState()")
	}
	go outg.SetState(true)
	time.Sleep(10 * time.Millisecond)
	if GetStateOrPanic(outslow) != true {
		fmt.Println("For this test, please connect Pin P8_7 to P8_8")
		t.Error("2: outslow.GetState() != outg.SetState()")
	}
	if GetStateOrPanic(ing) != true {
		t.Error("2: ing.GetState() != outg.SetState()")
	}
	go outg.SetState(false)
	time.Sleep(10 * time.Millisecond)
	if GetStateOrPanic(ing) != false {
		fmt.Println("For this test, please connect Pin P8_7 to P8_8")
		t.Error("3: ing.GetState() != outg.SetState()")
	}
	//Step(outg, 20, time.Duration(200)*time.Millisecond, nil)
}

func Test_FakeGPIO(t *testing.T) {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	f1 := NewFakeNamedGPIO("fake 1", OUT, logger)
	f2 := NewFakeGPIO(2, IN)
	//next line should not generate output
	f2.FakeInput(false)
	FakeGPIODefaultLogTarget_ = logger
	//now this should write output
	f2.FakeInput(true)
	if GetStateOrPanic(f2) != true {
		t.Error("f2 fake input did not work")
	}
	f1.SetState(true)
	f1.ConnectTo(f2)
	f1.SetState(false)
	if GetStateOrPanic(f1) != false {
		t.Error("f1 SetState did not work")
	}
	if GetStateOrPanic(f2) != false {
		t.Error("Fake connection to f2 did not work")
	}
}
