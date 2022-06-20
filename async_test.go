package generics

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var ErrFatal = errors.New("fatal")

func calculate(ctx context.Context, i int) (int, error) {
	return 20 * i, nil
}

func atoi(ctx context.Context, s string) (int, error) {
	return strconv.Atoi(s)
}

func itoa(ctx context.Context, i int) (string, error) {
	return strconv.Itoa(i), nil
}

func sleep(ctx context.Context, d time.Duration) (int, error) {
	time.Sleep(d)
	return int(d), nil
}

func fatal(ctx context.Context, i int) (int, error) {
	return 0, ErrFatal
}

func TestRun(t *testing.T) {
	ctx := context.Background()

	promise1 := Go(ctx, 10, calculate)
	result1, err := promise1.Resolve()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result1 != 200 {
		t.Error("expecting 200, but got ", result1)
	}

	promise2 := Go(ctx, time.Second, sleep)
	result2, err := promise2.Resolve()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result2 != 1000000000 {
		t.Error("expecting 1000000000, got ", result2)
	}

	promise3 := Go(ctx, 1, fatal)
	result3, err := promise3.Resolve()
	if !errors.Is(err, ErrFatal) {
		t.Error("expecting ErrFatal, got ", err)
	}
	if result3 != 0 {
		t.Error("expecting 0, got ", result3)
	}
}

func TestTry(t *testing.T) {
	ctx := context.Background()
	promise1 := Go(ctx, time.Second, sleep)
	result1, err := promise1.Try()
	if result1 != 0 {
		t.Error("expecting 0, got ", result1)
	}
	if !errors.Is(err, ErrNotDone) {
		t.Error("expecting ErrNotDone, got ", err)
	}

	result1, err = promise1.Resolve()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result1 != 1000000000 {
		t.Error("expecting 1000000000, got ", result1)
	}

	result1, err = promise1.Try()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result1 != 1000000000 {
		t.Error("expecting 1000000000, got ", result1)
	}
}

func TestWait(t *testing.T) {
	ctx := context.Background()

	promise1 := Go(ctx, "100", atoi)
	promise2 := Go(ctx, 200, calculate)

	err := Wait(promise1, promise2)
	if err != nil {
		t.Error("expecting no error, got ", err)
	}

	result1, err := promise1.Resolve()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result1 != 100 {
		t.Error("expecting 100, got ", result1)
	}

	result2, err := promise2.Resolve()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result2 != 4000 {
		t.Error("expecting 4000, got ", result2)
	}

	promise3 := Go(ctx, 2*time.Second, sleep)
	promise4 := Go(ctx, 0, fatal)

	err = Wait(promise3, promise4)
	if !errors.Is(err, ErrFatal) {
		t.Error("expecting ErrFatal, got ", err)
	}

	result3, err := promise3.Try()
	if result3 != 0 {
		t.Error("expecting 0, got ", result3)
	}
	if !errors.Is(err, ErrNotDone) {
		t.Error("expecting ErrNotDone, got ", err)
	}

	result4, err := promise4.Resolve()
	if result4 != 0 {
		t.Error("expecting 0, got ", result4)
	}
	if !errors.Is(err, ErrFatal) {
		t.Error("expecting ErrFatal, got ", err)
	}
}

func TestWithCancel(t *testing.T) {
	withCancel := WithCancel(sleep)

	ctx := context.Background()

	promise1 := Go(ctx, time.Second, withCancel)
	result1, err := promise1.Resolve()
	if err != nil {
		t.Error("expecting no error, got ", err)
	}
	if result1 != 1000000000 {
		t.Error("expecting 1000000000, got ", result1)
	}

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	promise2 := Go(ctx, time.Second, withCancel)
	result2, err := promise2.Resolve()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Error("expecting DeadlineExceeded, got ", err)
	}

	if result2 != 0 {
		t.Error("expecting 0, got ", result2)
	}
}

func TestThen(t *testing.T) {
	ctx := context.Background()

	promise1 := Go(ctx, 10, calculate)
	promise2 := Then(ctx, promise1, itoa)

	result2, err := promise2.Resolve()
	if result2 != "200" {
		t.Error(`expecting 200, but got`, result2)
	}
	if err != nil {
		t.Error("expecting no error, got ", err)
	}

	promise3 := Go(ctx, 0, fatal)
	promise4 := Then(ctx, promise3, itoa)
	result4, err := promise4.Resolve()
	if len(result4) > 0 {
		t.Error("expecting empty string, got ", result4)
	}
	if !errors.Is(err, ErrFatal) {
		t.Error("expecting ErrFatal, got", err)
	}
}

func TestInvoke(t *testing.T) {
	is := assert.New(t)

	err := fmt.Errorf("failed")

	t1, err1 := Invoke(12, func(i int) error {
		return nil
	})
	is.Equal(t1, 1)
	is.Equal(err1, nil)

	t2, err2 := Invoke(12, func(i int) error {
		if i == 11 {
			return nil
		}
		return err
	})
	is.Equal(t2, 12)
	is.Equal(err2, nil)

	t3, err3 := Invoke(2, func(i int) error {
		if i == 11 {
			return nil
		}
		return err
	})
	is.Equal(t3, 2)
	is.Equal(err3, err)

	t4, err4 := Invoke(0, func(i int) error {
		if i < 100 {
			return err
		}

		return nil
	})
	is.Equal(t4, 101)
	is.Equal(err4, nil)
}

func TestDelayedInvoke(t *testing.T) {
	is := assert.New(t)

	err := fmt.Errorf("failed")

	t1, time1, err1 := DelayedInvoke(42, 10*time.Millisecond, func(i int, d time.Duration) error {
		return nil
	})
	is.Equal(t1, 1)
	is.Greater(time1, 0*time.Millisecond)
	is.Less(time1, 1*time.Millisecond)
	is.Equal(err1, nil)

	t2, time2, err2 := DelayedInvoke(42, 10*time.Millisecond, func(i int, d time.Duration) error {
		if i == 5 {
			return nil
		}

		return err
	})
	is.Equal(t2, 6)
	is.Greater(time2, 50*time.Millisecond)
	is.Less(time2, 60*time.Millisecond)
	is.Equal(err2, nil)

	t3, time3, err3 := DelayedInvoke(2, 10*time.Millisecond, func(i int, d time.Duration) error {
		if i == 5 {
			return nil
		}

		return err
	})
	is.Equal(t3, 2)
	is.Greater(time3, 10*time.Millisecond)
	is.Less(time3, 20*time.Millisecond)
	is.Equal(err3, err)

	t4, time4, err4 := DelayedInvoke(0, 10*time.Millisecond, func(i int, d time.Duration) error {
		if i < 10 {
			return err
		}

		return nil
	})
	is.Equal(t4, 11)
	is.Greater(time4, 100*time.Millisecond)
	is.Less(time4, 115*time.Millisecond)
	is.Equal(err4, nil)
}

func TestDebounce(t *testing.T) {
	fn1 := func() {
		t.Log("step 1 : called once after 10ms when func stopped")
	}
	fn2 := func() {
		t.Log("step 2 : called once after 10ms when func stopped")
	}
	fn3 := func() {
		t.Log("step 3 : called once after 10ms when func stopped")
	}

	t1, _ := NewDebounce(10*time.Millisecond, fn1)
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			t1()
		}
		time.Sleep(20 * time.Millisecond)
	}

	t2, _ := NewDebounce(10*time.Millisecond, fn2)
	for i := 0; i < 3; i++ {
		for j := 0; j < 5; j++ {
			t2()
		}
		time.Sleep(5 * time.Millisecond)
	}

	time.Sleep(10 * time.Millisecond)

	t3, cancel := NewDebounce(10*time.Millisecond, fn3)
	for i := 0; i < 3; i++ {
		for j := 0; j < 10; j++ {
			t3()
		}
		time.Sleep(20 * time.Millisecond)
		if i == 0 {
			cancel()
		}
	}
}
