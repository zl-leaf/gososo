package context
import(
	"errors"
)

type Context struct {
	services map[string]Service
	components map[string]Component
}

type Service interface {
	Init() error
	Start() error
	Stop() error
}

type Component interface{

}

func New() *Context {
	context := &Context{}
	context.services = make(map[string]Service)
	context.components = make(map[string]Component)
	return context
}

func (context *Context) AddService(name string, service Service) error {
	if _,exist := context.services[name];exist {
		return errors.New("已经存在该组件")
	}
	context.services[name] = service
	return nil
}

func (context *Context) GetService(name string) (service Service, exist bool) {
	service,exist = context.services[name]
	return
}

func (context *Context) AddComponent(name string, component Component) error {
	if _,exist := context.components[name];exist {
		return errors.New("已经存在该组件")
	}
	context.components[name] = component
	return nil
}

func (context *Context) GetComponent(name string) (component Component, exist bool) {
	component,exist = context.components[name]
	return
}