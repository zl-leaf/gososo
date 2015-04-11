package robots

import(
    "net/http"
    "net/url"
    "io/ioutil"
    "errors"
    "strings"
    "regexp"
)

type Robots struct {
    agent string
    robots map[string]Robot
}

type Robot struct {
    KeyValues []KeyValue
}

type KeyValue struct {
    Key string
    Val string
}

func New(agent string) *Robots {
    robots := &Robots{agent:agent, robots:make(map[string]Robot)}
    return robots
}

/**
 * 获取单个robot
 */
func (robots *Robots) GetRobot(host string) (robot Robot) {
    if host[len(host)-1:] != "/" {
        host = host + "/"
    }
    if r,exist := robots.robots[host];exist {
        robot = r
        return
    }
    robot,err := NewRobot(host, robots.agent)
    if err == nil {
        robots.robots[host] = robot
    }
    return
}

/**
 * 新建一个robot
 */
func NewRobot(host string, userAgent string) (robot Robot, err error) {
    content,err := GetRobotsContent(host)
    if err != nil {
        err = errors.New("获取robots内容失败")
        return
    }

    robot,err = parse(content, userAgent)
    return
}

/**
 * 获取指定UserAgent的Robot
 */
func parse(content string, userAgent string) (robot Robot, err error) {
    robot = Robot{KeyValues:make([]KeyValue, 0)}

    lines := strings.Split(content, "\n")
    lineCount := len(lines)
    ua := ""
    for i:=0; i<lineCount; i++ {
        line := lines[i]
        if strings.TrimSpace(line) == "" {
            continue
        }
        if line[0:1] == "#" {
            continue
        }
        kv := parseKeyValue(line)
        if kv.Key == "User-agent" {
            ua = kv.Val
        } else {
            if ua == userAgent || ua == "*" {
                robot.KeyValues = append(robot.KeyValues, kv)
            }
        }
    }

    return
}

func parseKeyValue(str string) (kv KeyValue) {
    tmpKV := strings.Split(str, ":")
    if len(tmpKV) < 2 {
        return
    }

    key := strings.TrimSpace(tmpKV[0])
    val := strings.TrimSpace(tmpKV[1])

    kv = KeyValue{Key:key, Val:val}
    return
}

/**
 * 获取Robots内容
 */
func GetRobotsContent(host string) (content string, err error) {
    if strings.TrimSpace(host) == "" {
        err = errors.New("域名为空")
        return
    }
    if host[len(host)-1:] != "/" {
        host = host + "/"
    }
    robotUrl := host + "robots.txt"
    resp, err := http.Get(robotUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()
    if resp.StatusCode == 404 {
        return
    }

    b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

    content = string(b)
    return
}

/**
 * 检查URL是否允许抓取
 */
func (robot *Robot)  IsAllow(e string) bool {
    u,err := url.Parse(e)
	if err != nil {
		return false
	}
    for _,kv := range robot.KeyValues {
		k := kv.Key
		v := kv.Val
		if k == "Disallow" || k == "Allow" {
            r := regexp.MustCompile(v)
            match := r.Match([]byte(u.Path))
            if match {
                if k == "Disallow" {
                    return false
                } else {
                    return true
                }
            }
        }
	}
    return true
}
