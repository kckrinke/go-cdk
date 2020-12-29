package cdk

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSignalingBasics(t *testing.T) {
	Convey("Basic Signaling Features", t, func() {
		s := new(CSignaling)
		So(s.Init(), ShouldEqual, false)
		So(s.Init(), ShouldEqual, true)
		someData := "some data"
		signalCaught := false
		s.Connect(
			SignalEventError,
			"basic-test",
			func(data []interface{}, argv ...interface{}) EventFlag {
				So(data, ShouldHaveLength, 1)
				So(data[0], ShouldEqual, "some data")
				So(argv, ShouldHaveLength, 1)
				err := fmt.Errorf("an error")
				So(argv[0], ShouldHaveSameTypeAs, err)
				So(argv[0].(error).Error(), ShouldEqual, err.Error())
				signalCaught = true
				return EVENT_PASS
			},
			someData,
		)
		So(signalCaught, ShouldEqual, false)
		So(
			s.Emit(SignalEventError, fmt.Errorf("an error")),
			ShouldEqual,
			EVENT_PASS,
		)
		So(signalCaught, ShouldEqual, true)
		// So(s.Disconnect(SignalEventError, "basic-test"), ShouldBeNil)
	})
}

func TestSignalingPatterns(t *testing.T) {
	var hit0, hit1, hit2 bool
	reset := func() { hit0, hit1, hit2 = false, false, false }
	reset()
	hit0fn := func(data []interface{}, argv ...interface{}) EventFlag {
		hit0 = true
		return EVENT_PASS
	}
	hit1fn := func(data []interface{}, argv ...interface{}) EventFlag {
		hit1 = true
		return EVENT_STOP
	}
	hit2fn := func(data []interface{}, argv ...interface{}) EventFlag {
		hit2 = true
		return EVENT_PASS
	}
	s := new(CSignaling)
	Convey("Signaling Init", t, func() {
		So(s, ShouldNotBeNil)
		So(s.Init(), ShouldEqual, false)
		So(s.Init(), ShouldEqual, true)
	})
	Convey("Signaling Listeners", t, func() {
		s.Connect(SignalEventError, "many-errors-0", hit0fn)
		s.Connect(SignalEventError, "many-errors-1", hit1fn)
		s.Connect(SignalEventError, "many-errors-2", hit2fn)
		reset()
		So(
			s.Emit(SignalEventError, fmt.Errorf("an error")),
			ShouldEqual,
			EVENT_STOP,
		)
		So(hit0, ShouldEqual, true)
		So(hit1, ShouldEqual, true)
		So(hit2, ShouldEqual, false)
		So(s.Disconnect(SignalEventError, "many-errors-1"), ShouldBeNil)
		reset()
		So(
			s.Emit(SignalEventError, fmt.Errorf("an error")),
			ShouldEqual,
			EVENT_PASS,
		)
		So(hit0, ShouldEqual, true)
		So(hit1, ShouldEqual, false)
		So(hit2, ShouldEqual, true)
		err := s.Disconnect(SignalEventError, "many-errors-1")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "unknown signal handle: many-errors-1")
	})
	Convey("Signaling Regulating", t, func() {
		So(s.Disconnect(SignalEventError, "many-errors-0"), ShouldBeNil)
		So(s.Disconnect(SignalEventError, "many-errors-1"), ShouldNotBeNil)
		So(s.Disconnect(SignalEventError, "many-errors-2"), ShouldBeNil)
		s.StopSignal(SignalEventError)
		So(s.IsSignalStopped(SignalEventError), ShouldEqual, true)
		So(s.IsSignalPassed(SignalEventError), ShouldEqual, false)
		So(
			s.Emit(SignalEventError, fmt.Errorf("an error dropped")),
			ShouldEqual,
			EVENT_STOP,
		)
		s.ResumeSignal(SignalEventError)
		s.ResumeSignal(SignalEventError)
		s.ResumeSignal("this is not really signal")
		So(s.IsSignalStopped(SignalEventError), ShouldEqual, false)
		So(s.IsSignalPassed(SignalEventError), ShouldEqual, false)
		So(
			s.Emit(SignalEventError, fmt.Errorf("an error stopped")),
			ShouldEqual,
			EVENT_PASS,
		)
		s.Connect(SignalEventError, "many-errors-1", hit1fn)
		s.PassSignal(SignalEventError)
		So(s.IsSignalStopped(SignalEventError), ShouldEqual, false)
		So(s.IsSignalPassed(SignalEventError), ShouldEqual, true)
		So(
			s.Emit(SignalEventError, fmt.Errorf("an error passed")),
			ShouldEqual,
			EVENT_PASS,
		)
		s.ResumeSignal(SignalEventError)
	})
}
