package configure

const(
	GLOBAL_ENTITY = "global"
)

type Config struct {
	entities []*Entity
}

type Entity struct {
	name string
	attributes map[string]string
}

func (config *Config)Init() {
	config.entities = make([]*Entity, 0)
	e := &Entity{GLOBAL_ENTITY, make(map[string]string)}
	config.AddEntity(e)
}

func (config *Config)All() []*Entity{
	return config.entities
}

func (config *Config)AddEntity(e *Entity) (err error){
	config.entities = append(config.entities, e)
	return
}

func (config *Config)GetGloablEntity() (entity *Entity){
	for _,e := range config.entities {
		if e.Name() == GLOBAL_ENTITY {
			entity = e
			break 
		}
	}
	return
}

func (config *Config)GetEntity(name string) (entities []*Entity, exist bool){
	for _,e := range config.entities {
		if e.Name() == name {
			entities = append(entities, e)
			exist = true
			break 
		}
	}
	return
}

func (e *Entity)Name() string {
	return e.name
}

func (e *Entity)AddAttr(key string, value string) {
	e.attributes[key] = value
}

func (e *Entity)GetAttr(key string) (value string) {
	return e.attributes[key]
}

func (e *Entity)AllAttrs() map[string]string{
	return e.attributes
}