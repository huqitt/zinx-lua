package LuaState

import (
	"github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
	"os"
	"bufio"
	"strings"
	"io/ioutil"
	"fmt"
	"path/filepath"
)

var pathList = []string{
	"Script",
}

var FunctionList = []*lua.FunctionProto{}

func init(){
	// 获取当前路径（返回值：路径，错误信息）
	dir,_ := filepath.Abs(`.`)
	dir += "\\src\\"
	for _,v := range pathList {
		files,_,err := getFilesAndDirs(dir + v)
		if err != nil {
			panic(err)
		}
		for _,n := range files {
			function,err := CompileLuaFile(n)
			if err == nil {
				FunctionList = append(FunctionList, function)
				fmt.Println("CompileLuaFile succ, path:", n)
			} else {
				fmt.Println("CompileLuaFile err:", err, "\nPath:", n)
			}
		}
	}
}

// DoCompiledFile takes a FunctionProto, as returned by CompileLua, and runs it in the LState. It is equivalent
// to calling DoFile on the LState with the original source file.
func DoCompiledFile(L *lua.LState) error {
	for _,v := range FunctionList {
		lfunc := L.NewFunctionFromProto(v)
		L.Push(lfunc)
		err := L.PCall(0, lua.MultRet, nil)
		if err != nil {
			return err
		}
	}
	return nil
}


//获取指定目录下的所有文件和目录
func getFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	PthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+PthSep+fi.Name())
			getFilesAndDirs(dirPth + PthSep + fi.Name())
		} else {
			// 过滤指定格式
			ok := strings.HasSuffix(fi.Name(), ".lua")
			if ok {
				files = append(files, dirPth+PthSep+fi.Name())
			}
		}
	}

	return files, dirs, nil
}

// 编译 lua 代码文件
func CompileLuaFile(filePath string) (*lua.FunctionProto, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	chunk, err := parse.Parse(reader, filePath)
	if err != nil {
		return nil, err
	}
	proto, err := lua.Compile(chunk, filePath)
	if err != nil {
		return nil, err
	}
	return proto, nil
}

