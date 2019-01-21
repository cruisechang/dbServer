package handler

import (
	"errors"
	"math/rand"
	"time"

	"github.com/sony/sonyflake"
)

func createUniqueIDProvider() (*sonyflake.Sonyflake, error) {

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

	rand.Seed(time.Now().UTC().UnixNano())
	r := 1 + rand.Intn(99999999-1+1)

	return uint16(r), nil
}
