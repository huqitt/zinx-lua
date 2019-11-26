package LuaState

import (
	"sync"
	"github.com/yuin/gopher-lua"
	"fmt"
	"path/filepath"
	"strings"
	"utils"
)

type LuaStateData struct {
	L	*lua.LState
	IsUsed bool
	Id	int
}

var m     sync.Mutex
var statePool []*LuaStateData

func init() {
	statePool = make([]*LuaStateData, 0, utils.GlobalServer.WorkerPoolSize)
}

func GetLState() (*lua.LState, int) {
	m.Lock()
	defer m.Unlock()

	var res *LuaStateData = nil
	for i := 0; i < len(statePool); i++ {
		if !statePool[i].IsUsed {
			res = statePool[i]
			break
		}
	}

	if res == nil {
		res = &LuaStateData{
			L:new(),
			Id:len(statePool),
		}
		statePool = append(statePool, res)
	}

	res.IsUsed = true

	return res.L, res.Id
}

func AddServer(ser_id uint32, name string, desc string){
	m.Lock()
	defer m.Unlock()

	for i := 0; i < len(statePool); i++ {
		if statePool[i].IsUsed {
			lState := statePool[i].L
			err := lState.CallByParam(
				lua.P{
					Fn: lState.GetGlobal("Main").(*lua.LTable).RawGetString("AddServer").(*lua.LFunction),
					NRet: 0,
					Protect: true,
				},
				lua.LNumber(ser_id),
				lua.LString(name),
				lua.LString(desc),
			)
			if err != nil {
				panic(err)
			}
		}
	}
}

func CloseLState(id int){
	m.Lock()
	defer m.Unlock()

	if id < len(statePool) {
		statePool[id].IsUsed = false
	} else {
		fmt.Println("CloseLState err:\nId = ", id, ", statePool.len = ", len(statePool))
	}
}

func new() *lua.LState {
	L := lua.NewState()
	// 获取当前路径（返回值：路径，错误信息）
	dir, _ := filepath.Abs(`.`)
	luapath := dir + "\\src\\Script\\"
	//fmt.Println("dir:", dir)
	source := "package.path = '" + luapath + "'..'?.lua;'.. package.path"
	source = strings.Replace(source, "\\", "\\\\", -1)
	err0 := L.DoString(source)
	if err0 != nil {
		panic(err0)
	}
	err1 := DoCompiledFile(L)
	if err1 != nil {
		panic(err1)
	}

	// setting the L up here.
	// load scripts, set global variables, share channels, etc...
	return L
}

func ClearLStateAll() {
	m.Lock()
	defer m.Unlock()
	for i := 0; i < len(statePool); i++ {
		statePool[i].L.Close()
		if statePool[i].IsUsed {
			panic("关闭了一个正在使用的虚拟机!")
		}
	}
	statePool = nil
}
