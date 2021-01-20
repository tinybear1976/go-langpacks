package langpacks

import (
	"bufio"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tinybear1976/go-langpacks/redisdb"
)

type LoadResult struct {
	LangTag  string
	FileName string
	Estimate int
	Reality  int
}

type LoadMode uint8

const (
	InMemory LoadMode = 0
	InRedis  LoadMode = 1
)

var (
	version       string
	langs         map[string]map[int]string
	loadmode      LoadMode
	redis_ip_port string
	redis_pwd     string
	redis_db      int
	lps_path      string
	lps_suffix    string
	lps_separator string
	is_loaded     bool
)

func init() {
	loadmode = InMemory
	lps_path = "./"
	lps_suffix = ".lps"
	lps_separator = "~"
	is_loaded = false
}

func InitLangPacksDefault() {
	InitLangPacks("", "", "", "", "", 0)
}

func InitLangPacksDefaultRedis(ipWithPort, pwd string, db int) {
	InitLangPacks("", "", "", ipWithPort, pwd, db)
	SetLoadMode(InRedis)
}

func InitLangPacks(lpsPath string, lpsSuffix string, separator string, ipWithPort, pwd string, db int) {
	redis_ip_port = ipWithPort
	redis_pwd = pwd
	redis_db = db
	if lpsPath != "" {
		lps_path = lpsPath
	}
	if lpsSuffix != "" {
		lps_suffix = strings.ToLower(lpsSuffix)
	}
	if separator != "" {
		lps_separator = separator
	}
}

func SetLoadMode(mode LoadMode) {
	loadmode = mode
	is_loaded = false
}

func Query(langTag string, textId int) (str string) {
	if !is_loaded {
		return
	}
	switch loadmode {
	case InRedis:
		key := "lang::" + langTag + "::" + strconv.Itoa(textId)
		str, _ = redisdb.GET("lang", key)
	case InMemory:
		m, ok := langs[langTag]
		if ok {
			str = m[textId]
		}
	}
	return
}

func Load() (rst []LoadResult, err error) {
	dir, err := ioutil.ReadDir(lps_path)
	if err != nil {
		return nil, err
	}
	switch loadmode {
	case InRedis:
		redisdb.Destroy()
		redisdb.New("go-langpacks", redis_ip_port, redis_pwd, redis_db)
	case InMemory:
		langs = make(map[string]map[int]string)
	}

	for _, fi := range dir {
		if fi.IsDir() {
			continue
		}
		if strings.HasSuffix(strings.ToLower(fi.Name()), lps_suffix) {
			file := path.Join(lps_path, fi.Name())
			switch loadmode {
			case InMemory:
				langpack := make(map[int]string)
				tag, es, rely := loadLangPacksbyMemory(file, lps_separator, langpack)
				if tag != "" {
					langs[tag] = langpack
					rst = append(rst, LoadResult{
						LangTag:  tag,
						FileName: file,
						Estimate: es,
						Reality:  rely,
					})
				}
			case InRedis:
				tag, es, rely := loadLangPacks(file, lps_separator)
				if tag != "" {
					rst = append(rst, LoadResult{
						LangTag:  tag,
						FileName: file,
						Estimate: es,
						Reality:  rely,
					})
				}
			}
		}
	}
	is_loaded = true
	return
}

func loadLangPacks(filenameWithPath string, separator string) (langTag string, estimate int, reality int) {
	file, err := os.Open(filenameWithPath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	firtRow := true
	for scanner.Scan() {
		line := scanner.Text()
		if firtRow {
			firtRow = false
			langTag = strings.Trim(line, " ")
			if langTag == "" {
				return
			}
			continue
		}
		segs := strings.Split(line, separator)
		if len(segs) != 2 {
			continue
		}
		estimate++
		no, ok := strconv.Atoi(strings.Trim(segs[0], " "))
		if ok == nil {
			key := "lang::" + langTag + "::" + strconv.Itoa(no)
			err := redisdb.SET("go-langpacks", key, segs[1])
			if err == nil {
				reality++
			}
		}
	}
	return
}

func loadLangPacksbyMemory(filenameWithPath string, separator string, m map[int]string) (langTag string, estimate int, reality int) {
	file, err := os.Open(filenameWithPath)
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	firtRow := true
	for scanner.Scan() {
		line := scanner.Text()
		if firtRow {
			firtRow = false
			langTag = strings.Trim(line, " ")
			if langTag == "" {
				return
			}
			continue
		}
		segs := strings.Split(line, separator)
		if len(segs) != 2 {
			continue
		}
		estimate++
		no, ok := strconv.Atoi(strings.Trim(segs[0], " "))
		if ok == nil {
			m[no] = segs[1]
			reality++
		}
	}
	return
}
