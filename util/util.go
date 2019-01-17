package util

import (
	"time"
	"github.com/sony/sonyflake"
	"net"
	"errors"
	"math/rand"
)

func GetUniqueID()(uint64,error) {


	st := sonyflake.Settings{}
	st.StartTime = getTime()
	st.MachineID = machineID
	st.CheckMachineID = checkMachineID
	sf := sonyflake.NewSonyflake(st)
	if sf == nil {
		return 0,errors.New("sonyflake is nil")
	}
	if id,err:=sf.NextID();err!=nil{
		return 0, err
	}else{
		return id,nil
	}


}
func getTime() time.Time {

	return time.Now()
}
func machineID() (uint16, error) {
	//ip,_:=getGlobalUnicastIP()
	////ipStr := os.Getenv("MY_IP")
	////if len(ipStr) == 0 {
	////	return 0, errors.New("'MY_IP' environment variable not set")
	////}
	//ip1 := net.ParseIP(ip.String())
	//if len(ip) < 4 {
	//	return 0, errors.New("invalid IP")
	//}
	//return uint16(ip1[15])<<8 + uint16(ip1[16]), nil
	r:=RandomInt(1,9999999999)

	return uint16(r),nil
}
func checkMachineID(id uint16) (bool) {
	return id != 0
}
func getGlobalUnicastIP()( net.IP,error){
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
					return ip,nil
					//t.Logf("ip0 IsGlobalUnicast %s\n", ip.String())
				}
				//if ip.IsInterfaceLocalMulticast() {
				//	t.Logf("ip0 IsInterfaceLocalMulticast %s\n", ip.String())
				//}
				//if ip.IsLinkLocalMulticast() {
				//	t.Logf("ip0 IsLinkLocalMulticast %s\n", ip.String())
				//}
				//if ip.IsLinkLocalUnicast() {
				//	t.Logf("ip0 IsLinkLocalUnicast %s\n", ip.String())
				//}
				//if ip.IsLoopback() {
				//	t.Logf("ip0 IsLoopback %s\n", ip.String())
				//}
				//if ip.IsMulticast() {
				//	t.Logf("ip0 IsMulticast %s\n", ip.String())
				//}
			}
		}
	}
	return nil,errors.New("not found")
}
func RandomInt(min, max int) int {
	if min >= max {
		return min
	}

	rand.Seed(time.Now().UTC().UnixNano())
	return  min+rand.Intn(max-min+1)
}
