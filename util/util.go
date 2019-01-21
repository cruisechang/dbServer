package util

import (
	"errors"
	"math/rand"
	"time"

	"github.com/sony/sonyflake"
)

//RandomInt return ran int
//include max
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//return min + r.Intn(max-min+1)

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min+1)

	// return `0 <= n < 100`  not include 100
	//fmt.Print(rand.Intn(100), ",")
	//fmt.Print(rand.Intn(100))
}

func CreateUniqueIDProvider() (*sonyflake.Sonyflake, error) {

	const shortForm = "2007-Jan-02"
	t, _ := time.Parse(shortForm, "2013-Feb-03")

	st := sonyflake.Settings{}
	st.StartTime = t
	st.MachineID = machineID
	//st.CheckMachineID = checkMachineID   //不需要
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return nil, errors.New("sonyFlake is nil")
	}
	return sf, nil
}
func machineID() (uint16, error) {

	r:= RandomInt(1,9999999)
	return uint16(r), nil
}
