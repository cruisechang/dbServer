package util

import (
	"errors"
	"math/rand"
	"time"

	"github.com/sony/sonyflake"
)

func GetUniqueID() (uint64, error) {

	st := sonyflake.Settings{}
	st.StartTime = getTime()
	st.MachineID = machineID
	//st.CheckMachineID = checkMachineID   //不需要
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return 0, errors.New("sonyflake is nil")
	}
	if id, err := sf.NextID(); err != nil {
		return 0, err
	} else {
		return id, nil
	}

}
func getTime() time.Time {

	const shortForm = "2007-Jan-02"
	t, _ := time.Parse(shortForm, "2013-Feb-03")
	return t
}
func machineID() (uint16, error) {
	r := RandomInt(1, 9999999999)
	return uint16(r), nil
}
func checkMachineID(id uint16) bool {
	return id != 0
}

func RandomInt(min, max int) int {
	if min >= max {
		return min
	}

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min+1)
}
