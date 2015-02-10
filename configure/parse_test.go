package configure
import(
	"testing"
	"fmt"
)

func Test_InitConfig(t *testing.T) {
	config := InitConfig("./config.example.ini")
	es := config.All() 
	for _,e := range es {
		fmt.Println(e.Name())
		for key,_ := range e.AllAttrs() {
			fmt.Print(key+" ")
			fmt.Println(e.GetAttr(key))
		}
	}
}