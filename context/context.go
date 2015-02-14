package context
import(
	"errors"
)

type Context struct {
	components map[string]interface{}
}

func New() *Context {
	context := &Context{}
	context.components = make(map[string]interface{})
	return context
}

func (context *Context) AddComponent(name string, component interface{}) error {
	if _,exist := context.components[name];exist {
		return errors.New("已经存在该组件")
	}
	context.components[name] = component
	return nil
}

func (context *Context) GetComponent(name string) (component interface{}, exist bool) {
	component,exist = context.components[name]
	return
}