package configuration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"sync"
	"testing"
)

func TestParseKeyOrder(t *testing.T) {

	wg := &sync.WaitGroup{}

	fn := func() {

		defer func() {
			wg.Done()
		}()

		for i := 0; i < 100000; i++ {
			conf := LoadConfig("tests/configs.conf")
			for g := 1; g < 3; g++ {
				for i := 1; i < 4; i++ {
					key := fmt.Sprintf("test.out.a.b.c.d.groups.g%d.o%d.order", g, i)
					order := conf.GetInt32(key, -1)

					if order != int32(i) {
						fmt.Println(conf)
						t.Fatalf("order not match,group %d, except: %d, real order: %d", g, i, order)
						return
					}
				}
			}
			conf = nil
			runtime.Gosched()
		}
	}

	wg.Add(2)

	go fn()
	go fn()

	wg.Wait()
}

func TestConfig_GetBoolean(t *testing.T) {
	conf := ParseString("{k1:TRUE, k2:faLSE, k3: yes,k4:no, k5 : on , k6:oFf}")
	assert.Equal(t, true, conf.GetBoolean("k1"))
	assert.Equal(t, false, conf.GetBoolean("k2"))
	assert.Equal(t, true, conf.GetBoolean("k3"))
	assert.Equal(t, false, conf.GetBoolean("k4"))
	assert.Equal(t, true, conf.GetBoolean("k5"))
	assert.Equal(t, false, conf.GetBoolean("k6"))
}

func TestConfig_GetInt32(t *testing.T) {
	conf := ParseString("{k1:2147483647, k2:-2147483648,}")
	assert.Equal(t, int32(2147483647), conf.GetInt32("k1"))
	assert.Equal(t, int32(-2147483648), conf.GetInt32("k2"))
}

func TestConfig_GetInt64(t *testing.T) {
	conf := ParseString("{k1:9223372036854775807, k2:-9223372036854775808,}")
	assert.Equal(t, int64(9223372036854775807), conf.GetInt64("k1"))
	assert.Equal(t, int64(-9223372036854775808), conf.GetInt64("k2"))
}

func TestConfig_GetFloat32(t *testing.T) {
	conf := ParseString("{k1:1e3, k2:1e-3,}")
	assert.Equal(t, float32(1000), conf.GetFloat32("k1"))
	assert.Equal(t, float32(0.001), conf.GetFloat32("k2"))
}

func TestConfig_GetFloat64(t *testing.T) {
	conf := ParseString("{k1:1e3, k2:1e-3,}")
	assert.Equal(t, 1000., conf.GetFloat64("k1"))
	assert.Equal(t, 0.001, conf.GetFloat64("k2"))
}

func TestConfig_GetBooleanSafely(t *testing.T) {
	conf := ParseString("{k1:TRUE, k2:faLSE, k3: yes,k4:no, k5 : on , k6:oFf}")
	var v bool
	var err error

	v, err = conf.GetBooleanSafely("k1")
	if assert.Nil(t, err) {
		assert.Equal(t, true, v)
	}

	v, err = conf.GetBooleanSafely("k2")
	if assert.Nil(t, err) {
		assert.Equal(t, false, v)
	}

	v, err = conf.GetBooleanSafely("k3")
	if assert.Nil(t, err) {
		assert.Equal(t, true, v)
	}

	v, err = conf.GetBooleanSafely("k4")
	if assert.Nil(t, err) {
		assert.Equal(t, false, v)
	}

	v, err = conf.GetBooleanSafely("k5")
	if assert.Nil(t, err) {
		assert.Equal(t, true, v)
	}

	v, err = conf.GetBooleanSafely("k6")
	if assert.Nil(t, err) {
		assert.Equal(t, false, v)
	}
}

func TestConfig_GetBooleanSafelyError(t *testing.T) {
	conf := ParseString("{k1:qwerty,}")

	_, err := conf.GetBooleanSafely("k1")
	assert.Error(t, err)
}

func TestConfig_GetInt32Safely(t *testing.T) {
	conf := ParseString("{k1:2147483647, k2:-2147483648,}")
	var v int32
	var err error

	v, err = conf.GetInt32Safely("k1")
	if assert.Nil(t, err) {
		assert.Equal(t, int32(2147483647), v)
	}

	v, err = conf.GetInt32Safely("k2")
	if assert.Nil(t, err) {
		assert.Equal(t, int32(-2147483648), v)
	}
}

func TestConfig_GetInt32SafelyError(t *testing.T) {
	conf := ParseString("{k1:2147483648, k2:-2147483649, k3: qwerty}")
	var err error

	_, err = conf.GetInt32Safely("k1")
	assert.Error(t, err)

	_, err = conf.GetInt32Safely("k2")
	assert.Error(t, err)

	_, err = conf.GetInt32Safely("k3")
	assert.Error(t, err)
}

func TestConfig_GetInt64Safely(t *testing.T) {
	conf := ParseString("{k1:9223372036854775807, k2:-9223372036854775808,}")
	var v int64
	var err error

	v, err = conf.GetInt64Safely("k1")
	if assert.Nil(t, err) {
		assert.Equal(t, int64(9223372036854775807), v)
	}

	v, err = conf.GetInt64Safely("k2")
	if assert.Nil(t, err) {
		assert.Equal(t, int64(-9223372036854775808), v)
	}
}

func TestConfig_GetInt64SafelyError(t *testing.T) {
	conf := ParseString("{k1:9223372036854775808, k2:-9223372036854775809, k3: qwerty}")
	var err error

	_, err = conf.GetInt64Safely("k1")
	assert.Error(t, err)

	_, err = conf.GetInt64Safely("k2")
	assert.Error(t, err)

	_, err = conf.GetInt64Safely("k3")
	assert.Error(t, err)
}

func TestConfig_GetFloat32Safely(t *testing.T) {
	conf := ParseString("{k1:1e3, k2:1e-3,}")
	var v float32
	var err error

	v, err = conf.GetFloat32Safely("k1")
	if assert.Nil(t, err) {
		assert.Equal(t, float32(1000), v)
	}

	v, err = conf.GetFloat32Safely("k2")
	if assert.Nil(t, err) {
		assert.Equal(t, float32(0.001), v)
	}
}

func TestConfig_GetFloat32SafelyError(t *testing.T) {
	conf := ParseString("{k1:qwerty,}")

	_, err := conf.GetFloat32Safely("k1")
	assert.Error(t, err)
}

func TestConfig_GetFloat64Safely(t *testing.T) {
	conf := ParseString("{k1:1e3, k2:1e-3,}")
	var v float64
	var err error

	v, err = conf.GetFloat64Safely("k1")
	if assert.Nil(t, err) {
		assert.Equal(t, 1000., v)
	}

	v, err = conf.GetFloat64Safely("k2")
	if assert.Nil(t, err) {
		assert.Equal(t, 0.001, v)
	}
}

func TestConfig_GetFloat64SafelyError(t *testing.T) {
	conf := ParseString("{k1:qwerty,}")

	_, err := conf.GetFloat64Safely("k1")
	assert.Error(t, err)
}
