package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/ProtossGenius/SureMoonNet/basis/smn_data"
	"github.com/ProtossGenius/SureMoonNet/basis/smn_file"
	"github.com/ProtossGenius/pglang/analysis/lex_pgl"
	"github.com/ProtossGenius/pglang/snreader"
	"golang.org/x/crypto/ssh/terminal"
)

func parseGrammer(c <-chan snreader.ProductItf) {
	for {
		it := (<-c).(*lex_pgl.LexProduct)
		typ, value := it.ProductType(), it.Value
		if value == "\t" {
			fmt.Println("haha")
		}
		if typ == int(lex_pgl.PGLA_PRODUCT_HAN) {
			fmt.Println("Warn")
		}

	}
}

func read() {
	var char rune
	var err error
	sm := lex_pgl.NewLexAnalysiser()
	go parseGrammer(sm.GetResultChan())
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(fd, state)
	reader := bufio.NewReader(os.Stdin)
	for {
		char, _, err = reader.ReadRune()
		if char != '\t' {
			fmt.Printf("%c", char)
		}
		if err = sm.Read(&lex_pgl.PglaInput{Char: char}); err != nil {
			log.Fatal("when parse lex error : ", err)
		}
	}
}

// get user home.
func getUserHome() string {
	user, err := user.Current()
	if err == nil {
		return user.HomeDir
	}
	return ""

}

func getPath(f ...string) string {
	return strings.Join(f, smn_file.PathSep)

}

var (
	configPath = getPath(getUserHome(), ".config", "mws")
)

func getConfigPath() string {
	return configPath
}

type WorkSpace struct {
	Projects []Proj `json:"projects"`
}

// Proj comments.
type Proj struct {
	Path    string
	SubPath []string
	Tasks   []string
}

type Task struct {
	Branch string
	Desc   string
	Status string
}

type CmdList struct {
	Cmds []string
}

func loadConfig(cfgPath string, obj interface{}) {
	wsCfg := "{}"
	if smn_file.IsFileExist(cfgPath) {
		data, err := smn_file.FileReadAll(cfgPath)
		if err != nil {
			log.Fatal(err)
		}

		wsCfg = string(data)
	}
	smn_data.GetDataFromStr(wsCfg, obj)
}

// Manager manage current status.
type Manager struct {
	Commands map[string]interface{}
	Ws       *WorkSpace
}

func (m *Manager) Init(ws *WorkSpace) {
	m.Ws = ws

}

func (m *Manager) OnCmd() {
}

func main() {
	if !smn_file.IsFileExist(getConfigPath()) {
		os.MkdirAll(getConfigPath(), os.ModeDir)
	}

	ws := &WorkSpace{}
	mgr := &Manager{}
	loadConfig(getPath(getConfigPath(), "ws.json"), ws)
	loadConfig(getPath(getConfigPath(), "cmd.json"), &mgr.Commands)
	mgr.Init(ws)

	mgr.OnCmd()
}
