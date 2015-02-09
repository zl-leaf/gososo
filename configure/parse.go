package configure
import(
	"log"
	"os"
	"bufio"
	"io"
	"strings"
)

const(
	DEFAULT_CONFIG_PATH = "./config.ini"
)

func InitConfig(configPath string) (config *Config){
	if strings.TrimSpace(configPath) == "" {
		configPath = DEFAULT_CONFIG_PATH
	}
	str := readConfig(configPath)
	config = parseConfigStr(str)
	return
}

func readConfig(configPath string) (configStr string){
	f, err := os.Open(configPath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}

	buf := bufio.NewReader(f)
	for {
		l, err := buf.ReadString('\n')
		line := strings.TrimSpace(l)
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			} else if len(line) == 0 {
				break
			}
		}
		if len(line) > 0 {
			configStr += line + "\n"
		}
	}
	return
}

func parseConfigStr(configStr string) (config *Config) {
	config = &Config{}
	config.Init()

	lines := strings.Split(configStr, "\n")
	tmpLine := ""
	flag := true
	var e *Entity

	for _, line := range lines {
		if len(line) > 0 {
			switch {
				case line[0] == '[' && line[len(line)-1] == ']':
					tmpLine = line
				case line[len(line)-1] == '[':
					tmpLine = line
					flag = false
				case line[len(line)-1] == ']':
					tmpLine += line
					flag = true
				default:
					if !flag {
						tmpLine += line
					} else {
						tmpLine = line
					}
			}

			if flag {
				if tmpLine[0] == '[' && tmpLine[len(line)-1] == ']' {
					name := parseConfigEName(tmpLine)
					e = &Entity{name, make(map[string][]string)}
					config.AddEntity(e)
				} else {
					key,value := parseConfigLine(tmpLine)
					if e == nil {
						e = config.GetGloablEntity()	
					}
					e.AddAttr(key, value)
					
				}
				
			}
		}
		
	}
	return
}

func parseConfigEName(line string) (name string) {
	name = strings.Trim(line, " []")
	return
}

func parseConfigLine(line string) (key string, value []string) {
	line = strings.Replace(line, " ", "", -1)
	kv := strings.Split(line, "=")
	if len(kv) != 2 {
		log.Fatal("配置文件格式错误")
	}

	key = kv[0]
	v := kv[1]
	if v[0] == '[' && v[len(v)-1] == ']' {
		value = strings.Split(v[1:len(v)-1], ",")
	} else {
		value = []string{v}
	}
	return
}