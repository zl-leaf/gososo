package robots

import(
    "testing"
    "log"
)
func Test_Parse(t *testing.T) {
    rs := New("*")
    r := rs.GetRobot("http://v.baidu.com/")
    log.Println(r.IsAllow("http://v.baidu.com/test"))
    log.Println(rs)

}
