package test

import (
	"errors"
	"github.com/sony/sonyflake"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestSonyflake(t *testing.T) {

	tt := time.Now()
	t.Logf("time %s\n", tt.String())

	st := sonyflake.Settings{}
	st.StartTime = getTime()
	st.MachineID = machineID
	//st.CheckMachineID = checkMachineID   //不需要
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		t.Fatalf("failed to initialize sonyflake")
	}

	for i := 0; i < 100; i++ {
		id, _ := sf.NextID()
		t.Logf("TestSonyflake id =%d\n", id)

	}

	//if id, err := sf.NextID(); err != nil {
	//	//t.Fatalf("TestSonyflake error=%s/n",err.Error())
	//} else {
	//	t.Logf("TestSonyflake id =%d\n", id)
	//}

}
func getTime() time.Time {
	const shortForm = "2007-Jan-02"
	t, _ := time.Parse(shortForm, "2013-Feb-03")

	return t
}
func machineID() (uint16, error) {
	//ip, _ := getGlobalUnicastIP()
	//ip1 := net.ParseIP(ip.String())
	//if len(ip) < 4 {
	//	return 0, errors.New("invalid IP")
	//}
	//return uint16(ip1[15])<<8 + uint16(ip1[16]), nil

	r := randomInt(1, 9999999999)
	return uint16(r), nil
}
func checkMachineID(id uint16) (bool) {
	return id != 0
}

func getGlobalUnicastIP() (net.IP, error) {
	ifaces, _ := net.Interfaces()
	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:

				ip = v.IP
				if ip.IsGlobalUnicast() {
					return ip, nil
					//t.Logf("ip0 IsGlobalUnicast %s\n", ip.String())
				}
				if ip.IsInterfaceLocalMulticast() {
					return ip, nil
					//t.Logf("ip0 IsInterfaceLocalMulticast %s\n", ip.String())
				}
				if ip.IsLinkLocalMulticast() {
					return ip, nil
					//t.Logf("ip0 IsLinkLocalMulticast %s\n", ip.String())
				}
				if ip.IsLinkLocalUnicast() {
					return ip, nil
					//t.Logf("ip0 IsLinkLocalUnicast %s\n", ip.String())
				}
				if ip.IsLoopback() {
					return ip, nil
					//t.Logf("ip0 IsLoopback %s\n", ip.String())
				}
				if ip.IsMulticast() {
					return ip, nil
					//t.Logf("ip0 IsMulticast %s\n", ip.String())
				}
			}
		}
	}
	return nil, errors.New("not found")
}

func randomInt(min, max int) int {
	if min >= max {
		return min
	}

	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min+1)
}
