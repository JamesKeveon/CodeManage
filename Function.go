package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/axgle/mahonia"
	"github.com/beevik/etree"
	"github.com/fsnotify/fsnotify"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/gogs/chardet"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type strCodePath struct {
	Name       string
	LocalPath  string
	url        string
	branchname string
}

type strGitOpResult struct {
	index    int
	bResult  bool
	bSuccess bool
	url      string
	errmsg   string
}

type Watch struct {
	watch *fsnotify.Watcher
}

var version, fistlocalpath, atomfistlocalpath, atomfistlocalpath_2, workpath, bakpath string
var lock sync.Mutex
var RecordFlag bool
var IsMovefromAM4 bool
var IsBatchCompare bool
var IsPullLS bool
var IsPullTrunk bool
var IsDelete bool
var GitOpResult []strGitOpResult
var LSGitOpResult [8]strGitOpResult

var Featrure_basic, Featrure_businpub, Featrure_equity, Featrure_equity_A, Featrure_equity_Common, Feature_equity_server, Feature_equity_server_common, Feature_basic_server, Feature_businpub_server bool
var RequirementNum, Version, baisc_Version, businpub_Version, equity_Verison, equity_A_Version, equity_Comm_Version, equity_server_Version, equity_server_common_Version, basic_server_Version, businpub_server_Version string
var linkA_1, linkA_2, linkA bool
var linkZ_1, linkZ_2, linkZ bool
var atom_1, atom_2 bool
var isDownload bool
var isPullFlag bool

// func main() {
// 	workpath, _ = GetRunPath()
// 	workpath += "\\"
// 	bakpath = workpath + ".cache\\"
// 	for {
// 		fmt.Printf("\n\n\n")
// 		fmt.Println("工作路径：", workpath)
// 		fmt.Println("输入要执行的操作:")
// 		fmt.Println("代码下载:-----------------输入1;")
// 		fmt.Println("开发模式:-----------------输入2;")
// 		fmt.Println("代码递交:-----------------输入3;")
// 		fmt.Println("文件差异比较:-------------输入4;")
// 		fmt.Println("迁移已合并文件:-----------输入5;")
// 		fmt.Println("代码更新:-----------------输入6;")
// 		fmt.Println("退出:---------------------输入exit;")
// 		Mode := ""
// 		fmt.Scanln(&Mode)
// 		if Mode == "1" {
// 			CodeDownLoad(false)
// 		} else if Mode == "2" {
// 			DevelopMode()
// 		} else if Mode == "3" {
// 			CodeCommit()
// 		} else if Mode == "4" {
// 			for {
// 				fmt.Println("输入要比较的的函数中文名，输入exit-退出......")
// 				CompareCheName := ""
// 				fmt.Scanln(&CompareCheName)
// 				if CompareCheName == "exit" {
// 					break
// 				}

// 				CompareFile(CompareCheName)
// 			}

// 		} else if Mode == "5" {
// 			fmt.Println("迁移已合并文件默认从资管版的路径迁移，所以先将合并完成的文件放在资管版的路径！")
// 			for {
// 				fmt.Println("输入要迁移的函数中文名，输入exit-退出......")
// 				FileCheName := ""
// 				fmt.Scanln(&FileCheName)
// 				if FileCheName == "exit" {
// 					break
// 				}
// 				DistPath := ""
// 				if strings.Contains(FileCheName, "AS") || strings.Contains(FileCheName, "AF") {
// 					DistPath = workpath + "UFT-Atom"
// 				} else if strings.Contains(FileCheName, "LS") || strings.Contains(FileCheName, "LF") {
// 					DistPath = workpath + "UFT-Business"
// 				}
// 				AutoMoveFile(FileCheName, workpath+"UFT_equity_Z", DistPath, true)
// 			}
// 		} else if Mode == "6" {
// 			CodeDownLoad(true)
// 		} else if Mode == "exit" {
// 			break
// 		} else {

// 			fmt.Println("\n\n\n付费定制！")
// 			fmt.Printf("\n\n\n")
// 		}
// 	}

// 	// fmt.Println("输入任意键退出！")
// 	// exitflag := ""
// 	// fmt.Scanln(&exitflag)

// }
func (f *TForm1) GetVariables(nodeName string) ([]string, []string, []string, []string, []string, error) {
	//open the file
	xmlFile, err := os.Open("options.xml")
	if err != nil {
		return nil, nil, nil, nil, nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var name1 []string
	var name2 []string
	var url []string
	var url2 []string
	var brunchName []string
	var inElement string

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == nodeName {
				inElement = nodeName
			}
			if inElement == nodeName && se.Name.Local == "Variable" {
				// for _, attr := range se.Attr {
				// 	if attr.Name.Local == "name1" {
				// 		name1 = append(name1, attr.Value)
				// 	}
				// }
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						name2 = append(name2, attr.Value)
					}
				}
				for _, attr := range se.Attr {
					if attr.Name.Local == "url" {
						url = append(url, attr.Value)
					}
				}
				// for _, attr := range se.Attr {
				// 	if attr.Name.Local == "url2" {
				// 		url2 = append(url2, attr.Value)
				// 	}
				// }
				for _, attr := range se.Attr {
					if attr.Name.Local == "brunchName" {
						brunchName = append(brunchName, attr.Value)
					}
				}
			}
		case xml.EndElement:
			if se.Name.Local == nodeName {
				inElement = ""
			}
		}
	}
	name1 = name2
	url2 = name2
	return name1, name2, url, url2, brunchName, nil
}

func (f *TForm1) GetCodePath(nodeName string) (string, error) {
	//open the file
	xmlFile, err := os.Open("options.xml")
	if err != nil {
		return "", fmt.Errorf("无法打开文件: %w", err)
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var codepath string
	var find bool = false

	for {
		t, err := decoder.Token()
		if err != nil || find {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == nodeName {
				for _, attr := range se.Attr {
					if attr.Name.Local == "codepath" {
						codepath = attr.Value
					}
				}
				fmt.Println("nodeName:" + nodeName + ",se:" + se.Name.Local + ",codepath:" + codepath)
				find = true
				break
			}
			fmt.Println("nodeName:" + nodeName + ",se:" + se.Name.Local + ",codepath:" + codepath)
		}
	}
	return codepath, nil
}

func (f *TForm1) CodeDownLoad(isPull bool) {
	root := f.ComboBox1.Text()

	name1, name2, url, url2, brunchName, err := f.GetVariables(root)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(name1) != len(brunchName) || len(name1) != len(name2) || len(name2) != len(url) || len(url) != len(brunchName) {
		fmt.Println("xml文件中" + root + "数据数量不一致")
		return
	}

	var CodePath []strCodePath
	GitOpResult = GitOpResult[:0]
	for i := 0; i < len(name1); i++ {
		// 创建一个新的strCodePath实例并添加到切片中
		CodePath = append(CodePath, strCodePath{name1[i], workpath + name2[i], url[i], brunchName[i]})
		GitOpResult = append(GitOpResult, strGitOpResult{i, false, false, url2[i], ""})
	}

	var LSCodePath [10]strCodePath
	LSCodePath[0] = strCodePath{"basic", workpath + "UFT_basic", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-basic/UFT.git", "dev-IPS1.0-basicV202301.04.000.LS"}
	LSCodePath[1] = strCodePath{"businpub", workpath + "UFT_businpub", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-businpub/UFT.git", "dev-IPS1.0-businpubV202301.04.000.LS"}
	LSCodePath[2] = strCodePath{"equity_A", workpath + "UFT_equity_A", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT.git", "dev-IPS1.0-equityV202301A.04.000.LS"}
	LSCodePath[3] = strCodePath{"UFT-Metadata", workpath + "UFT-Common/UFT-Metadata", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/uft-metadata.git", "dev-IPS1.0-equityV202301A.04.000.LS"}
	LSCodePath[4] = strCodePath{"UFT-Structure", workpath + "UFT-Common/UFT-Structure", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT-Structure.git", "dev-IPS1.0-equityV202301A.04.000.LS"}
	LSCodePath[5] = strCodePath{"UFT-Atom", workpath + "UFT-Common/UFT-Atom", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT-Atom.git", "dev-IPS1.0-equityV202301A.04.000.LS"}
	LSCodePath[6] = strCodePath{"UFT-Factor", workpath + "UFT-Common/UFT-Factor", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT-Factor.git", "dev-IPS1.0-equityV202301A.04.000.LS"}
	LSCodePath[7] = strCodePath{"UFT-Business", workpath + "UFT-Common/UFT-Business", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT-Business.git", "dev-IPS1.0-equityV202301A.04.000.LS"}

	LSGitOpResult[0] = strGitOpResult{0, false, false, "IPS1.0-basic", ""}
	LSGitOpResult[1] = strGitOpResult{1, false, false, "IPS1.0-businpub", ""}
	LSGitOpResult[2] = strGitOpResult{3, false, false, "IPS1.0-equity-A", ""}
	LSGitOpResult[3] = strGitOpResult{4, false, false, "IPS1.0-equity/uft-metadata", ""}
	LSGitOpResult[4] = strGitOpResult{5, false, false, "IPS1.0-equity/UFT-Structure", ""}
	LSGitOpResult[5] = strGitOpResult{6, false, false, "IPS1.0-equity/UFT-Atom", ""}
	LSGitOpResult[6] = strGitOpResult{7, false, false, "IPS1.0-equity/UFT-Factor", ""}
	LSGitOpResult[7] = strGitOpResult{8, false, false, "IPS1.0-equity/UFT-Business", ""}

	//url := "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/uft-metadata.git"
	//username := f.EditUserID.Text()
	//password := f.EditPasswd.Text()
	//username := ""
	//password := ""
	strWriteMsg := fmt.Sprint("拉取代码在[", workpath, "]目录下")
	f.ShowLog(strWriteMsg)
	if IsPullLS {
		for index := range LSCodePath {
			if LSCodePath[index].Name != "" {
				if isPull {
					go f.osExecPull(index, LSCodePath[index].LocalPath, LSCodePath[index].url, LSCodePath[index].branchname, "branch")
				} else {
					_, err := os.Stat(workpath + "UFT-Common")
					if os.IsNotExist(err) {
						err := os.Mkdir(workpath+"UFT-Common", 0755)
						if err != nil {
							strWriteMsg := fmt.Sprint("创建目录失败", workpath+"UFT-Common", err)
							f.ShowLog(strWriteMsg)
							return
						}
					}
					//go f.gogitclone(index, username, password, CodePath[index])
					go f.osExecGitClone(index, LSCodePath[index].LocalPath, LSCodePath[index].url, LSCodePath[index].branchname, "branch")
				}

			}
		}
	} else {
		for index := range CodePath {
			if CodePath[index].Name != "" {
				if isPull {
					go f.osExecPull(index, CodePath[index].LocalPath, CodePath[index].url, CodePath[index].branchname, "branch")
				} else {
					//go f.gogitclone(index, username, password, CodePath[index])
					go f.osExecGitClone(index, CodePath[index].LocalPath, CodePath[index].url, CodePath[index].branchname, "branch")
				}

			}
		}
	}

	go f.ShowGitResult()

}

func (f *TForm1) gogitclone(index int, username string, password string, CodePath strCodePath) bool {
	//strWriteMsg := fmt.Sprint("拉取:   ", CodePath.Name)
	//f.ShowLog(strWriteMsg)
	strWriteMsg := fmt.Sprint("本地路径:", CodePath.LocalPath, ", Git路径:", CodePath.url, ", 分支名称:", CodePath.branchname)
	f.ShowLog(strWriteMsg)

	r, err := git.PlainClone(CodePath.LocalPath, false, &git.CloneOptions{
		URL:           CodePath.url,
		RemoteName:    "origin",
		ReferenceName: plumbing.ReferenceName(CodePath.branchname),
		Auth:          &http.BasicAuth{Username: username, Password: password},
		//Bare:          true,
		//SingleBranch: true,
		//Progress: os.Stdout,
	})

	if err != nil {
		strWriteMsg = fmt.Sprint("拉取:   ", CodePath.Name, "  失败", err)
		//f.ShowLog(strWriteMsg)
		if IsPullLS {
			LSGitOpResult[index].errmsg = err.Error()
			LSGitOpResult[index].bSuccess = false
			LSGitOpResult[index].bResult = true
		} else {
			GitOpResult[index].errmsg = err.Error()
			GitOpResult[index].bSuccess = false
			GitOpResult[index].bResult = true
		}
		return false
	} else {
		ref, err := r.Head()
		if err != nil {
			fmt.Println("Error reading HEAD: ", err)
		}

		strWriteMsg = fmt.Sprint("拉取:   ", CodePath.Name, "  成功!   最后一次递交SHA-1是:", ref)
		//f.ShowLog(strWriteMsg)
		if IsPullLS {
			LSGitOpResult[index].bSuccess = true
			LSGitOpResult[index].bResult = true
		} else {
			GitOpResult[index].bSuccess = true
			GitOpResult[index].bResult = true
		}

	}

	return true
}

func (f *TForm1) osExecGitClone(index int, workspace, url, referenceName, refType string) error {

	strWriteMsg := fmt.Sprint("git client clone :", url, " by  branch: ", referenceName)
	f.ShowLog(strWriteMsg)

	isEmpty, errorinfo, error := CheckFolderISEmpty(workspace)
	if !isEmpty {
		if IsPullLS {
			LSGitOpResult[index].bSuccess = false
			LSGitOpResult[index].errmsg = errorinfo
			LSGitOpResult[index].bResult = true
		} else {
			GitOpResult[index].bSuccess = false
			GitOpResult[index].errmsg = errorinfo
			GitOpResult[index].bResult = true
		}
		return error
	}
	// git clone -- branch dev-trunk https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-basic/UFT.git E:\VS_PRO\GO_PRO\GO_VCL_PRO\CodeManage\workspace
	cmd := exec.Command("git", "clone", "--branch", referenceName, url, workspace)

	cmd.Dir = workspace
	// 设置工作环境变量，防止cmd弹窗
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // 创建一个新的会话组
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		f.ShowLog(err.Error())
	}

	if strings.Contains(string(out), "error:") || strings.Contains(string(out), "failed") || strings.Contains(string(out), "fatal") {
		if IsPullLS {
			LSGitOpResult[index].bSuccess = false
			LSGitOpResult[index].errmsg = string(out)
		} else {
			GitOpResult[index].bSuccess = false
			GitOpResult[index].errmsg = string(out)
		}
	} else {
		if IsPullLS {
			LSGitOpResult[index].bSuccess = true
		} else {
			GitOpResult[index].bSuccess = true
		}

	}
	if IsPullLS {
		LSGitOpResult[index].bResult = true
	} else {
		GitOpResult[index].bResult = true
	}

	//strWriteMsg := fmt.Sprint("git client pull " + string(out))
	//f.ShowLog(strWriteMsg)

	return err
}

func CheckFolderISEmpty(dirName string) (bool, string, error) {
	dirPath, err := filepath.Abs(dirName)
	errorinfo := ""
	if err != nil {
		fmt.Println("无法获取目录路径:", err)
		errorinfo = fmt.Sprint("无法获取目录路径:", err)
		return false, errorinfo, err
	}
	// 检查目录是否存在
	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			errorinfo = fmt.Sprint("创建目录失败:", err)
			return false, errorinfo, err
		}

		return true, errorinfo, err
	}

	// 检查目录是否为空
	isEmpty, err := isDirectoryEmpty(dirPath)
	if err != nil {
		errorinfo = fmt.Sprint("检查目录是否为空时出错:", err)
		return false, errorinfo, err
	}
	if !isEmpty {
		errorinfo = fmt.Sprint("目录不为空，不能执行下载。")
	}
	return isEmpty, errorinfo, err
}

func isDirectoryEmpty(dir string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}
	// 如果没有文件或子目录，则目录为空
	return len(files) == 0, nil
}

// git client pull
func (f *TForm1) osExecPull(index int, workspace, url, referenceName, refType string) error {
	if refType == "branch" {
		strWriteMsg := fmt.Sprint("git client pull :", url, " by  branch: ", referenceName)
		f.ShowLog(strWriteMsg)
	} else if refType == "tag" {
		strWriteMsg := fmt.Sprint("git client pull :", url, " by  tag: ", referenceName)
		f.ShowLog(strWriteMsg)
	} else if refType == "commit" {
		strWriteMsg := fmt.Sprint("git client pull :", url, " by  commit id: ", referenceName)
		f.ShowLog(strWriteMsg)
	}
	cmd := exec.Command("git", "pull", url, referenceName)
	cmd.Dir = workspace
	// 设置工作环境变量，防止cmd弹窗
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // 创建一个新的会话组
	}
	out, err := cmd.CombinedOutput()
	if strings.Contains(string(out), "error:") || strings.Contains(string(out), "failed") || strings.Contains(string(out), "fatal") {
		if IsPullLS {
			LSGitOpResult[index].bSuccess = false
			LSGitOpResult[index].errmsg = string(out)
		} else {
			GitOpResult[index].bSuccess = false
			GitOpResult[index].errmsg = string(out)
		}

	} else {
		if IsPullLS {
			LSGitOpResult[index].bSuccess = true
		} else {
			GitOpResult[index].bSuccess = true
		}

	}
	if IsPullLS {
		LSGitOpResult[index].bResult = true
	} else {
		GitOpResult[index].bResult = true
	}

	//strWriteMsg := fmt.Sprint("git client pull " + string(out))
	//f.ShowLog(strWriteMsg)

	return err
}

func (f *TForm1) osExecCreateFeatures(isCreateNewbrach bool) {
	var CodePath [9]strCodePath
	CodePath[0] = strCodePath{"basic", workpath + "UFT_basic", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-basic/UFT.git", "dev-trunk"}
	CodePath[1] = strCodePath{"businpub", workpath + "UFT_businpub", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-businpub/UFT.git", "dev-trunk"}
	CodePath[2] = strCodePath{"equity_Z", workpath + "UFT_equity_Z", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT.git", "dev-trunk"}
	CodePath[3] = strCodePath{"equity_A", workpath + "UFT_equity_A", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT.git", "dev-trunk-A"}
	CodePath[4] = strCodePath{"UFT-Common", workpath + "UFT-Common", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/UFT-Common.git", "dev-trunk"}

	CodePath[5] = strCodePath{"equity-server", workpath + "Server_equity", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/Server.git", "dev-trunk"}
	CodePath[6] = strCodePath{"equity-server-common", workpath + "Server_Common", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-equity/Server_common.git", "dev-trunk"}
	CodePath[7] = strCodePath{"basic-server", workpath + "Server_basic", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-basic/Server.git", "dev-trunk"}
	CodePath[8] = strCodePath{"businpub-server", workpath + "Server_businpub", "https://hsgit.hundsun.com/O45/IPS1.0/IPS1.0-businpub/Server.git", "dev-trunk"}

	GitOpResult = append(GitOpResult, strGitOpResult{0, false, false, "IPS1.0-basic", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{1, false, false, "IPS1.0-businpub", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{2, false, false, "IPS1.0-equity-Z", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{3, false, false, "IPS1.0-equity-A", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{4, false, false, "UFT-Common", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{5, false, false, "IPS1.0-equity-server", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{6, false, false, "IPS1.0-equity-server-common", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{7, false, false, "IPS1.0-basic-server", ""})
	GitOpResult = append(GitOpResult, strGitOpResult{8, false, false, "IPS1.0-businpub-server", ""})

	// GitOpResult[0] = strGitOpResult{0, false, false, "IPS1.0-basic", ""}
	// GitOpResult[1] = strGitOpResult{1, false, false, "IPS1.0-businpub", ""}
	// GitOpResult[2] = strGitOpResult{2, false, false, "IPS1.0-equity-Z", ""}
	// GitOpResult[3] = strGitOpResult{3, false, false, "IPS1.0-equity-A", ""}
	// GitOpResult[4] = strGitOpResult{4, false, false, "UFT-Common", ""}

	for index := range CodePath {
		if CodePath[index].Name != "" {

			isCreatebrach, Featrure_Name := GetFeatureName(CodePath[index].Name)
			if !isCreatebrach {
				strWriteMsg := fmt.Sprint(CodePath[index].Name, "  未勾选不创建（切换）分支。")
				f.ShowLog(strWriteMsg)
				continue
			}

			if isCreateNewbrach {
				cmd := exec.Command("git", "pull", CodePath[index].url, CodePath[index].branchname)
				cmd.Dir = CodePath[index].LocalPath
				// 设置工作环境变量，防止cmd弹窗
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true, // 创建一个新的会话组
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					strWriteMsg := fmt.Sprint(CodePath[index].Name + "更新失败，不创建分支。\n" + string(output))
					f.ShowLog(strWriteMsg)
					continue
				}

				cmd = exec.Command("git", "checkout", "-b", Featrure_Name)
				cmd.Dir = CodePath[index].LocalPath
				// 设置工作环境变量，防止cmd弹窗
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true, // 创建一个新的会话组
				}

				output, err = cmd.CombinedOutput()
				if err != nil {
					//log.Fatalf("Failed to create remote branch: %s", err)
					strWriteMsg := fmt.Sprint(CodePath[index].Name + "创建本地分支失败。\n" + string(output))

					if strings.Contains(strWriteMsg, "already exists") {
						strWriteMsg := fmt.Sprint("\n" + CodePath[index].Name + " 本地已存在,自动切到feature分支。")
						f.ShowLog(strWriteMsg)
						f.SwitchGitBrach(Featrure_Name, CodePath[index].LocalPath, CodePath[index].Name)
					} else {
						f.ShowLog(strWriteMsg)
					}
					continue
				} else {
					//log.Printf("Created remote branch: %s", output)
					strWriteMsg := fmt.Sprint("Created branch:" + string(output))
					f.ShowLog(strWriteMsg)
				}

				cmd = exec.Command("git", "push", "origin", "--no-verify", "-u", Featrure_Name)
				cmd.Dir = CodePath[index].LocalPath
				// 设置工作环境变量，防止cmd弹窗
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true, // 创建一个新的会话组
				}
				output, err = cmd.CombinedOutput()
				if err != nil {
					//log.Fatalf("Failed to create remote branch: %s", err)
					strWriteMsg := fmt.Sprint(CodePath[index].Name + "推送远程分支失败。\n" + string(output))
					f.ShowLog(strWriteMsg)
				} else {

					//log.Printf("Created remote branch: %s", output)
					strWriteMsg := fmt.Sprint("Created remote branch:" + string(output))
					f.ShowLog(strWriteMsg)
				}
			} else {
				// 切换分支
				cmd := exec.Command("git", "checkout", CodePath[index].branchname)
				cmd.Dir = CodePath[index].LocalPath
				// 设置工作环境变量，防止cmd弹窗
				cmd.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true, // 创建一个新的会话组
				}
				output, err := cmd.CombinedOutput()
				if err != nil {
					strWriteMsg := fmt.Sprint(CodePath[index].Name + "切换到  [" + CodePath[index].branchname + "]  失败。" + string(output))
					f.ShowLog(strWriteMsg)
				} else {
					strWriteMsg := fmt.Sprint(CodePath[index].Name + "切换到  [" + CodePath[index].branchname + "]  成功。")
					f.ShowLog(strWriteMsg)
				}
			}
		}
	}
}
func (f *TForm1) SwitchGitBrach(Brachname string, LocalPath string, Name string) {
	// 切换分支
	cmd := exec.Command("git", "checkout", Brachname)
	cmd.Dir = LocalPath
	// 设置工作环境变量，防止cmd弹窗
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // 创建一个新的会话组
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		strWriteMsg := fmt.Sprint(Name + "切换到  [" + Brachname + "]  失败。" + string(output))
		f.ShowLog(strWriteMsg)
	} else {
		strWriteMsg := fmt.Sprint(Name + "切换到  [" + Brachname + "]  成功。")
		f.ShowLog(strWriteMsg)
	}
}
func GetFeatureName(name string) (bool, string) {

	Featrure_Name := ""
	isCreate := false
	switch name {
	case "basic":
		if Featrure_basic {
			isCreate = true
			Featrure_Name = "feature-" + baisc_Version + "-" + RequirementNum
		}
	case "businpub":
		if Featrure_businpub {
			isCreate = true
			Featrure_Name = "feature-" + businpub_Version + "-" + RequirementNum
		}
	case "equity_Z":
		if Featrure_equity {
			isCreate = true
			Featrure_Name = "feature-" + equity_Verison + "-" + RequirementNum
		}
	case "equity_A":
		if Featrure_equity_A {
			isCreate = true
			Featrure_Name = "feature-" + equity_A_Version + "-" + RequirementNum
		}
	case "UFT-Common":
		if Featrure_equity_Common {
			isCreate = true
			Featrure_Name = "feature-" + equity_Comm_Version + "-" + RequirementNum
		}
	case "equity-server":
		if Feature_equity_server {
			isCreate = true
			Featrure_Name = "feature-" + equity_server_Version + "-" + RequirementNum
		}
	case "equity-server-common":
		if Feature_equity_server_common {
			isCreate = true
			Featrure_Name = "feature-" + equity_server_common_Version + "-" + RequirementNum
		}
	case "basic-server":
		if Feature_basic_server {
			isCreate = true
			Featrure_Name = "feature-" + basic_server_Version + "-" + RequirementNum
		}
	case "businpub-server":
		if Feature_businpub_server {
			isCreate = true
			Featrure_Name = "feature-" + businpub_server_Version + "-" + RequirementNum
		}
	default:
		return false, ""
	}
	return isCreate, Featrure_Name
}

// 包含了代码下载完成后的自动递交功能
func (f *TForm1) ShowGitResult() {

	time.Sleep(2 * time.Second)
	//FinishNum := 0
	for {

		if IsPullLS {
			bFinishNum := true
			for _, tmpGitOpResult := range LSGitOpResult {
				if !tmpGitOpResult.bResult {
					strWriteMsg := fmt.Sprint("正在下载/更新" + tmpGitOpResult.url)
					f.ShowLog(strWriteMsg)
					bFinishNum = false
				}
			}

			if bFinishNum {

				f.ShowLog("\n所有分支下载/更新完成...")
				for _, tmpGitOpResult := range LSGitOpResult {
					strWriteMsg := fmt.Sprint("下载/更新结果:" + tmpGitOpResult.url + ",成功.")
					if !tmpGitOpResult.bSuccess {
						strWriteMsg = fmt.Sprint("下载/更新结果:" + tmpGitOpResult.url + ",失败.")
					}
					f.ShowLog(strWriteMsg)
					isDownload = true
				}
				for _, tmpGitOpResult := range LSGitOpResult {
					if !tmpGitOpResult.bSuccess {
						strWriteMsg := fmt.Sprint(tmpGitOpResult.url + ",下载/更新失败原因:" + tmpGitOpResult.errmsg)
						f.ShowLog(strWriteMsg)
					}

				}
				break
			} else {
				f.ShowLog("---------------------------------------------------------------------------")
				time.Sleep(2 * time.Second)
			}
		} else {
			bFinishNum := true
			for _, tmpGitOpResult := range GitOpResult {
				if !tmpGitOpResult.bResult {

					strWriteMsg := fmt.Sprint("正在下载/更新" + tmpGitOpResult.url)
					f.ShowLog(strWriteMsg)
					bFinishNum = false
				}
			}

			if bFinishNum {
				f.ShowLog("\n所有分支下载/更新完成...")
				for _, tmpGitOpResult := range GitOpResult {
					strWriteMsg := fmt.Sprint("下载/更新结果:" + tmpGitOpResult.url + ",成功.")
					if !tmpGitOpResult.bSuccess {
						strWriteMsg = fmt.Sprint("下载/更新结果:" + tmpGitOpResult.url + ",失败.")
					}
					f.ShowLog(strWriteMsg)
					isDownload = true
				}
				for _, tmpGitOpResult := range GitOpResult {
					if !tmpGitOpResult.bSuccess {
						strWriteMsg := fmt.Sprint(tmpGitOpResult.url + ",下载/更新失败原因:" + tmpGitOpResult.errmsg)
						f.ShowLog(strWriteMsg)
					}

				}
				break
			} else {
				f.ShowLog("---------------------------------------------------------------------------")
				time.Sleep(3 * time.Second)
			}
		}
	}
	//执行更改所属用户的问题
	workpath = f.EdtWorkpath.Text()
	if len(workpath) > 0 {
		if workpath[len(workpath)-1] != '\\' {
			workpath += "\\"
		}
	}
	currentUser, err := user.Current()
	if err != nil {
		f.ShowLog(err.Error())
	}
	f.setOwner(workpath, currentUser.Username)

	//开始创建软链接

	if !isPullFlag {
		for {
			if isDownload {
				f.ShowLog("--------------------开始创建软链接-------------------------")
				go f.SymbolicLinkExecute("1")
				go f.SymbolicLinkExecute("2")
				isDownload = false
				break
			}
		}
	}

}

func (f *TForm1) CopyAtomToA_1(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_core"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom       拷贝到UFT_equity_A对应目录下")
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_core       拷贝到UFT_equity_A对应目录下---完成")
	} else {
		BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		err := Copy(SrcFilePath, BakFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
	}
}
func (f *TForm1) CopyAtomToA_2(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_manage"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_manage       拷贝到UFT_equity_A对应目录下")
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_manage       拷贝到UFT_equity_A对应目录下---完成")
	} else {
		BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		err := Copy(SrcFilePath, BakFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
	}
}

func (f *TForm1) CopyBusinessToA(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Business"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Business   拷贝到UFT_equity_A对应目录下")
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/uftbusiness/equity/"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		f.ShowLog("将  UFT-Common/UFT-Business   拷贝到UFT_equity_A对应目录下---完成")
	} else {
		BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/uftbusiness/equity/"
		err := Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}
	}
}

func (f *TForm1) CopyAtomToZ_1(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_core"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_core       拷贝到UFT_equity_Z对应目录下")
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_core       拷贝到UFT_equity_Z对应目录下---完成")
	} else {
		BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		err := Copy(SrcFilePath, BakFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
	}

}

func (f *TForm1) CopyAtomToZ_2(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_manage"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_manage       拷贝到UFT_equity_Z对应目录下")
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_manage       拷贝到UFT_equity_Z对应目录下---完成")
	} else {
		BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		err := Copy(SrcFilePath, BakFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
	}

}

func (f *TForm1) CopyBusinessToZ(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Business"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Business   拷贝到UFT_equity_Z对应目录下")
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/uftbusiness/equity/"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		f.ShowLog("将  UFT-Common/UFT-Business   拷贝到UFT_equity_Z对应目录下---完成")
	} else {
		BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/uftbusiness/equity/"
		err := Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}
	}
}

func (f *TForm1) DevelopMode(version string) {
	//f.Tips.Lines().Add("开发模式是将UFT-Metadata、UFT-Structure、UFT-Atom、UFT-Factor、UFT-Business拷贝到原始工程下：")
	//fmt.Println("只拷贝银信版:-----------------输入1;")
	//fmt.Println("只拷贝资管版:-----------------输入2;")
	//fmt.Println("两个版本都拷贝:---------------输入3;")
	//fmt.Scanln(&version)
	f.ShowLog("--------------------开始执行开发模式--------------------")
	if version == "1" || version == "3" {

		go f.CopyAtomToA_1(false)
		go f.CopyAtomToA_2(false)
		go f.CopyAtomToA_1(true)
		go f.CopyAtomToA_2(true)
		go f.CopyBusinessToA(false)
		go f.CopyBusinessToA(true)

		f.ShowLog("将  UFT-Common/UFT-Metadata   拷贝到UFT_equity_A对应目录下")
		//bT := time.Now()
		SrcFilePath := workpath + "UFT-Common/UFT-Metadata"
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/metadata"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/metadata"
		err = Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}
		f.ShowLog("将  UFT-Common/UFT-Metadata   拷贝到UFT_equity_A对应目录下---完成")

		f.ShowLog("将  UFT-Common/UFT-Structure  拷贝到UFT_equity_A对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Structure"
		DistFilePath = workpath + "UFT_equity_A/UFTDB_equity/uftstructure/equity/"
		err = Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		BakFilePath = bakpath + "UFT_equity_A/UFTDB_equity/uftstructure/equity/"
		err = Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}
		f.ShowLog("将  UFT-Common/UFT-Structure  拷贝到UFT_equity_A对应目录下---完成")

		f.ShowLog("将  UFT-Common/UFT-Factor     拷贝到UFT_equity_A对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Factor"
		DistFilePath = workpath + "UFT_equity_A/UFTDB_equity/uftfactor/equity/"
		err = Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		BakFilePath = bakpath + "UFT_equity_A/UFTDB_equity/uftfactor/equity/"
		err = Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}
		f.ShowLog("将  UFT-Common/UFT-Factor     拷贝到UFT_equity_A对应目录下---完成")
	}

	if version == "2" || version == "3" {

		go f.CopyAtomToZ_1(false)
		go f.CopyAtomToZ_2(false)
		go f.CopyAtomToZ_1(true)
		go f.CopyAtomToZ_2(true)
		go f.CopyBusinessToZ(false)
		go f.CopyBusinessToZ(true)

		f.ShowLog("将  UFT-Common/UFT-Metadata   拷贝到UFT_equity_Z对应目录下")
		SrcFilePath := workpath + "UFT-Common/UFT-Metadata"
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/metadata"
		err := Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/metadata"
		err = Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
		}
		f.ShowLog("将  UFT-Common/UFT-Metadata   拷贝到UFT_equity_Z对应目录下---完成")

		f.ShowLog("将  UFT-Common/UFT-Structure  拷贝到UFT_equity_Z对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Structure"
		DistFilePath = workpath + "UFT_equity_Z/UFTDB_equity/uftstructure/equity/"
		err = Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		BakFilePath = bakpath + "UFT_equity_Z/UFTDB_equity/uftstructure/equity/"
		err = Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
		}
		f.ShowLog("将  UFT-Common/UFT-Structure  拷贝到UFT_equity_Z对应目录下---完成")

		f.ShowLog("将  UFT-Common/UFT-Factor     拷贝到UFT_equity_Z对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Factor"
		DistFilePath = workpath + "UFT_equity_Z/UFTDB_equity/uftfactor/equity/"
		err = Copy(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		BakFilePath = bakpath + "UFT_equity_Z/UFTDB_equity/uftfactor/equity/"
		err = Copy(SrcFilePath, BakFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, BakFilePath)
		}
		f.ShowLog("将  UFT-Common/UFT-Factor     拷贝到UFT_equity_Z对应目录下---完成")
	}
	//f.ShowLog("==资源拷贝完成==")
	//fmt.Println("拷贝Factor")
}

func (f *TForm1) SymbolicLinkExecute(version string) {

	if version == "1" || version == "3" {
		go f.LinkAtomToA_1(false)
		go f.LinkAtomToA_2(false)
		go f.LinkAtomToA_1(true)
		go f.LinkAtomToA_2(true)
		go f.LinkBusinessToA(false)
		go f.LinkBusinessToA(true)

		f.ShowLog("将  UFT-Common/UFT-Metadata   链接到UFT_equity_A对应目录下")
		//bT := time.Now()
		SrcFilePath := workpath + "UFT-Common/UFT-Metadata"
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/metadata"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		// BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/metadata"
		// err = Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// 	// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
		// 	// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		// }

		f.ShowLog("将  UFT-Common/UFT-Structure  链接到UFT_equity_A对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Structure"
		DistFilePath = workpath + "UFT_equity_A/UFTDB_equity/uftstructure/equity/"
		err = Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		// BakFilePath = bakpath + "UFT_equity_A/UFTDB_equity/uftstructure/equity/"
		// err = Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// 	// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
		// 	// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		// }

		f.ShowLog("将  UFT-Common/UFT-Factor     链接到UFT_equity_A对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Factor"
		DistFilePath = workpath + "UFT_equity_A/UFTDB_equity/uftfactor/equity/"
		err = Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		// BakFilePath = bakpath + "UFT_equity_A/UFTDB_equity/uftfactor/equity/"
		// err = Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// 	// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
		// 	// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		// }

		for {
			time.Sleep(7 * time.Second)
			if linkA && linkA_1 && linkA_2 {
				f.ShowLog("--------------------银信软链接创建完成-------------------------")
				break
			} else {
				f.ShowLog("----银信软链接尚未完成----")
			}
		}
		linkA = false
		linkA_1 = false
		linkA_2 = false
	}

	if version == "2" || version == "3" {

		go f.LinkAtomToZ_1(false)
		go f.LinkAtomToZ_2(false)
		go f.LinkAtomToZ_1(true)
		go f.LinkAtomToZ_2(true)
		go f.LinkBusinessToZ(false)
		go f.LinkBusinessToZ(true)

		f.ShowLog("将  UFT-Common/UFT-Metadata   链接到UFT_equity_Z对应目录下")
		SrcFilePath := workpath + "UFT-Common/UFT-Metadata"
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/metadata"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		// BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/metadata"
		// err = Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// }

		f.ShowLog("将  UFT-Common/UFT-Structure  链接到UFT_equity_Z对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Structure"
		DistFilePath = workpath + "UFT_equity_Z/UFTDB_equity/uftstructure/equity/"
		err = Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		// BakFilePath = bakpath + "UFT_equity_Z/UFTDB_equity/uftstructure/equity/"
		// err = Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// }

		f.ShowLog("将  UFT-Common/UFT-Factor     链接到UFT_equity_Z对应目录下")
		SrcFilePath = workpath + "UFT-Common/UFT-Factor"
		DistFilePath = workpath + "UFT_equity_Z/UFTDB_equity/uftfactor/equity/"
		err = Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		}

		// BakFilePath = bakpath + "UFT_equity_Z/UFTDB_equity/uftfactor/equity/"
		// err = Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// }
		for {
			time.Sleep(7 * time.Second)
			if linkZ && linkZ_1 && linkZ_2 {
				f.ShowLog("--------------------资管软链接创建完成-------------------------")
				break
			} else {
				f.ShowLog("----资管软链接尚未完成----")
			}
		}
		linkZ = false
		linkZ_1 = false
		linkZ_2 = false
	}

}

func (f *TForm1) LinkAtomToA_1(isbak bool) {
	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_core"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_core       链接到UFT_equity_A对应目录下")
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}

		linkA_1 = true
	} else {
		// BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		// err := Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
		// 	f.ShowLog(strWriteMsg)
		// }
	}

}

func (f *TForm1) LinkAtomToA_2(isbak bool) {
	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_manage"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_manage       链接到UFT_equity_A对应目录下")
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}

		linkA_2 = true

	} else {
		// BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		// err := Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
		// 	f.ShowLog(strWriteMsg)
		// }
	}
}

func (f *TForm1) LinkBusinessToA(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Business"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Business   链接到UFT_equity_A对应目录下")
		DistFilePath := workpath + "UFT_equity_A/UFTDB_equity/uftbusiness/equity/"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		linkA = true
	} else {
		// BakFilePath := bakpath + "UFT_equity_A/UFTDB_equity/uftbusiness/equity/"
		// err := Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// 	// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
		// 	// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		// }
	}
}

func (f *TForm1) LinkAtomToZ_1(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_core"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_core       链接到UFT_equity_Z对应目录下")
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}
		linkZ_1 = true
	} else {
		// BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_core/"
		// err := Copy(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
		// 	f.ShowLog(strWriteMsg)
		// }
	}

}

func (f *TForm1) LinkAtomToZ_2(isbak bool) {

	SrcFilePath := workpath + "UFT-Common/UFT-Atom/atom_equity_manage"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Atom/atom_equity_manage       链接到UFT_equity_Z对应目录下")
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			fmt.Printf("%s:%s", err, DistFilePath)
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}

		linkZ_2 = true
	} else {
		// BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/"
		// err := Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	strWriteMsg := fmt.Sprint(err, BakFilePath, "\n")
		// 	f.ShowLog(strWriteMsg)
		// }
	}

}

func (f *TForm1) LinkBusinessToZ(isbak bool) {
	SrcFilePath := workpath + "UFT-Common/UFT-Business"
	if !isbak {
		f.ShowLog("将  UFT-Common/UFT-Business   链接到UFT_equity_Z对应目录下")
		DistFilePath := workpath + "UFT_equity_Z/UFTDB_equity/uftbusiness/equity/"
		err := Link(SrcFilePath, DistFilePath)
		if err != nil {
			strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
		}

		linkZ = true
	} else {
		// BakFilePath := bakpath + "UFT_equity_Z/UFTDB_equity/uftbusiness/equity/"
		// err := Link(SrcFilePath, BakFilePath)
		// if err != nil {
		// 	fmt.Printf("%s:%s", err, BakFilePath)
		// 	// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
		// 	// WriteStingToFile(strWriteMsg, "./执行信息.txt")
		// }
	}
}

func Link(from, to string) error {
	// cmd := exec.Command("cmd", "/C", "mklink", "/D", from, to)
	// err := cmd.Run()
	// return err
	f, e := os.Stat(from)
	if e != nil {
		return e
	}

	if f.IsDir() {
		// Create the target directory if it doesn't exist
		if _, e := os.Stat(to); os.IsNotExist(e) {
			if e := os.MkdirAll(to, 0777); e != nil {
				return e
			}
		}

		// Read the source directory
		list, e := os.ReadDir(from)
		if e != nil {
			return e
		}

		for _, item := range list {
			// Skip directories containing ".git"
			if strings.Contains(item.Name(), ".git") {
				continue
			}
			// Recursively create symlinks for the contents
			srcPath := filepath.Join(from, item.Name())
			dstPath := filepath.Join(to, item.Name())
			if e := Link(srcPath, dstPath); e != nil {
				return e
			}
		}
	} else {
		// Create a symbolic link for the file
		if e := os.Symlink(from, to); e != nil {
			return e
		}
	}
	return nil
}

func (f *TForm1) CodeCommitAtom_1(version string) {
	if version == "1" {
		f.ShowLog("将 UFT-Common/UFT-Atom     中的文件从UFT_equity_A还原")
		listFilesAtom(workpath+"UFT-Common/UFT-Atom", workpath+"UFT_equity_A/UFTDB_equity/uftatom/atom_equity/", true)

		atom_1 = true
	} else if version == "2" {
		f.ShowLog("将 UFT-Common/UFT-Atom      中的文件从UFT_equity_Z还原")
		listFilesAtom(workpath+"UFT-Common/UFT-Atom", workpath+"UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/", true)

		atom_2 = true
	} else {
		f.ShowLog("对不起，没有找到您要提交的版本！！！")
	}

}

func (f *TForm1) CodeCommitAtom_2(version string) {
	if version == "1" {
		f.ShowLog("将 UFT-Common/UFT-Atom/atom_equity_manage     中的文件从UFT_equity_A还原")
		listFilesAtom_2(workpath+"UFT-Common/UFT-Atom/atom_equity_manage", workpath+"UFT_equity_A/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/", true)
		f.ShowLog("将 UFT-Common/UFT-Atom/atom_equity_manage     中的文件从UFT_equity_A还原---完成")
	} else if version == "2" {
		f.ShowLog("将 UFT-Common/UFT-Atom/atom_equity_manage      中的文件从UFT_equity_Z还原")
		listFilesAtom_2(workpath+"UFT-Common/UFT-Atom/atom_equity_manage", workpath+"UFT_equity_Z/UFTDB_equity/uftatom/atom_equity/atom_equity_manage/", true)

	} else {
		f.ShowLog("对不起，没有找到您要提交的版本！！！")
	}

}

func (f *TForm1) CodeCommitBusiness(version string) {
	if version == "1" {
		f.ShowLog("将 UFT-Common/UFT-Business 中的文件从UFT_equity_A还原")
		listFiles(workpath+"UFT-Common/UFT-Business", workpath+"UFT_equity_A/UFTDB_equity/uftbusiness/equity/", true)

	} else if version == "2" {
		f.ShowLog("将 UFT-Common/UFT-Business 中的文件从UFT_equity_Z还原")
		listFiles(workpath+"UFT-Common/UFT-Business", workpath+"UFT_equity_Z/UFTDB_equity/uftbusiness/equity/", true)

	} else {
		f.ShowLog("对不起，没有找到您要提交的版本！！！")
	}
}

func (f *TForm1) CodeCommit(version string) {
	//f.ShowLog("代码提交是将目录UFT-Common/UFT-Metadata、UFT-Common/UFT-Structure、UFT-Common/UFT-Atom、UFT-Common/UFT-Factor、UFT-Common/UFT-Business 中的文件从原始工程还原出来:")
	f.ShowLog("--------------------开始执行递交模式--------------------")
	//fmt.Println("提交银信版:-----------------输入1;")
	//fmt.Println("提交资管版:-----------------输入2;")
	timestr := time.Now().Format("2006-01-02 15:04:05")
	strWriteMsg := fmt.Sprint(timestr, "发生变化的文件有：", "\n")
	WriteStingToFile(strWriteMsg, "./执行信息.txt")
	//fmt.Scanln(&version)
	if version == "1" {

		go f.CodeCommitAtom_1(version)
		//go f.CodeCommitAtom_2(version)

		f.ShowLog("将 UFT-Common/UFT-Metadata 中的文件从UFT_equity_A还原")
		listFiles(workpath+"UFT-Common/UFT-Metadata", workpath+"UFT_equity_A/UFTDB_equity/metadata/", true)

		f.ShowLog("将 UFT-Common/UFT-Structure中的文件从UFT_equity_A还原")
		listFiles(workpath+"UFT-Common/UFT-Structure", workpath+"UFT_equity_A/UFTDB_equity/uftstructure/equity/", true)

		f.ShowLog("将 UFT-Common/UFT-Factor   中的文件从UFT_equity_A还原")
		listFiles(workpath+"UFT-Common/UFT-Factor", workpath+"UFT_equity_A/UFTDB_equity/uftfactor/equity/", true)

		f.CodeCommitBusiness(version)
		for {
			time.Sleep(3 * time.Second)
			if atom_1 {
				f.ShowLog("=====银信资源还原完成，可继续进行代码提交=====")
				break
			} else {
				f.ShowLog("===银信资源还原尚未完成===")
			}
		}
		atom_1 = false
	} else if version == "2" {
		go f.CodeCommitAtom_1(version)
		//go f.CodeCommitAtom_2(version)

		f.ShowLog("将 UFT-Common/UFT-Metadata  中的文件从UFT_equity_Z还原")
		listFiles(workpath+"UFT-Common/UFT-Metadata", workpath+"UFT_equity_Z/UFTDB_equity/metadata/", true)

		f.ShowLog("将 UFT-Common/UFT-Structure 中的文件从UFT_equity_Z还原")
		listFiles(workpath+"UFT-Common/UFT-Structure", workpath+"UFT_equity_Z/UFTDB_equity/uftstructure/equity/", true)

		f.ShowLog("将 UFT-Common/UFT-Factor   中的文件从UFT_equity_Z还原")
		listFiles(workpath+"UFT-Common/UFT-Factor", workpath+"UFT_equity_Z/UFTDB_equity/uftfactor/equity/", true)

		f.CodeCommitBusiness(version)

		for {
			time.Sleep(3 * time.Second)
			if atom_2 {
				f.ShowLog("=====资管资源还原完成，可继续进行代码提交=====")
				break
			} else {
				f.ShowLog("===资管资源还原尚未完成===")
			}
		}
		atom_2 = false
	} else {
		f.ShowLog("对不起，没有找到您要提交的版本！！！")
	}

	//fmt.Println("还原Factor")
}
func (f *TForm1) AutoMoveFile(SeachFilechineseName string, dirname string, DistPath string, bfirst bool) {

	if bfirst {
		if strings.Contains(SeachFilechineseName, "AS") || strings.Contains(SeachFilechineseName, "AF") {
			fistlocalpath = dirname + "/UFTDB_equity/uftatom/atom_equity"
		} else if strings.Contains(SeachFilechineseName, "LS") || strings.Contains(SeachFilechineseName, "LF") {
			fistlocalpath = dirname + "/UFTDB_equity/uftbusiness/equity"
		} else if strings.Contains(SeachFilechineseName, "RS") || strings.Contains(SeachFilechineseName, "RF") {
			fistlocalpath = dirname + "/UFTDB_equity/uftfactor/equity"
		} else {
			fistlocalpath = dirname + "/UFTDB_equity/uftstructure/equity"
		}
	}

	fileInfos, err := os.ReadDir(dirname)
	if err != nil {
		//log.Fatal(err)
		return
	}
	for _, fi := range fileInfos {
		if strings.EqualFold(".git", fi.Name()) {
			continue
		}
		filename := dirname + "/" + fi.Name() //绝对路径带后缀文件名
		Suffix := path.Ext(filename)          //文件后缀名
		//fmt.Printf("Suffix = %v, len(Suffix)=%v\n", Suffix, len(Suffix))
		if len(Suffix) > 0 && SupportCompareFileType(Suffix) {
			var chinesename, strObjectId string
			readXml(filename, &chinesename, &strObjectId)
			if SeachFilechineseName == chinesename {
				//filenameall := path.Base(filename)
				srclocalpath := filename
				RealDistPath := strings.Replace(srclocalpath, fistlocalpath, DistPath, 1)
				err = Copy(filename, RealDistPath)
				if err != nil {
					fmt.Printf("迁移失败：%s:%s", err, filename)
					// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
					// WriteStingToFile(strWriteMsg, "./执行信息.txt")
				}
				// fmt.Println("迁移完成：从\n", filename, "\n迁移到\n", RealDistPath)
				// fmt.Printf("\n")

				strWriteMsg := fmt.Sprint("---------------------------------------------------\n迁移完成：从\n", filename, "\n迁移到\n", RealDistPath, "\n")
				f.ShowLog(strWriteMsg)
				WriteStingToFile(strWriteMsg, "./迁移文件汇总.txt")
				strWriteMsg = fmt.Sprint(SeachFilechineseName, "\n")
				WriteStingToFile(strWriteMsg, "./迁移完成函数汇总.txt")
				//DeleteFile(DistFilePath)
			} else {
				continue
			}
		}

		if fi.IsDir() {
			//继续遍历fi这个目录
			f.AutoMoveFile(SeachFilechineseName, filename, DistPath, false)
		}
	}
}

func (f *TForm1) CompareFile(chineseName string) {

	strWriteMsg := fmt.Sprint("正在查找资管版【", chineseName, "】函数")
	f.ShowLog(strWriteMsg)
	IsO45Find, O45DistFilePath, err := FindFile(workpath+"UFT_equity_Z/UFTDB_equity", "", chineseName, "")
	if err != nil {
		f.ShowLog(err.Error())
		return
	}
	if !IsO45Find {
		strWriteMsg := fmt.Sprint("资管版不存在该函数:", chineseName)
		f.ShowLog(strWriteMsg)
		return
	}

	strWriteMsg = fmt.Sprint("正在查找银信版【", chineseName, "】函数")
	f.ShowLog(strWriteMsg)
	IsAM4Find, AM4DistFilePath, err := FindFile(workpath+"UFT_equity_A/UFTDB_equity", "", chineseName, "")
	if err != nil {
		f.ShowLog(err.Error())
		return
	}
	if !IsAM4Find {
		strWriteMsg := fmt.Sprint("银信版不存在该函数:", chineseName)
		f.ShowLog(strWriteMsg)
		return
	}

	if IsO45Find && IsAM4Find {
		todaystr := time.Now().Format("2006-01-02") //当前日期
		var logfile string
		var resultfile string
		logfile = "./" + "/log/" + chineseName + ".log"                         //比较日志
		resultfile = "./" + "/result/" + chineseName + "_" + todaystr + ".html" //比较报告

		bcscriptfile := "@" + "./BeyondCompare" + "/bcscript.txt" //BeyondCompare比较配置文件
		// "/silent"去除弹窗提示
		command := exec.Command("./BeyondCompare/BCompare.exe", "/silent", bcscriptfile, logfile, resultfile, O45DistFilePath, AM4DistFilePath)
		//command := exec.Command("./BeyondCompare/BCompare.exe", "/", bcscriptfile, O45DistFilePath, O45DistFilePath)
		// 设置工作环境变量，防止cmd弹窗
		command.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true, // 创建一个新的会话组
		}
		output, err := command.CombinedOutput()
		if err != nil {
			strWriteMsg := fmt.Sprint("\n比较出错:", AM4DistFilePath, "\n", AM4DistFilePath, "\n")
			f.ShowLog(strWriteMsg)
			strWriteMsg = fmt.Sprint(fmt.Sprint(err) + ": " + string(output))
			f.ShowLog(strWriteMsg)
		} else {
			//fmt.Println(string(output))
			//fmt.Printf("比较完成:%s,%s\n", filename, newdirpath)
			//调用浏览器打开比较报告
			if !IsBatchCompare {
				command = exec.Command("cmd", "/C", "start", resultfile)
				command.SysProcAttr = &syscall.SysProcAttr{
					HideWindow: true, // 创建一个新的会话组
				}
				command.Run()
			}

		}
	}
}

// 复制整个文件夹或单个文件
func Copy(from, to string) error {
	f, e := os.Stat(from)
	if e != nil {
		return e
	}

	if f.IsDir() {
		//from是文件夹，那么定义to也是文件夹
		if list, e := os.ReadDir(from); e == nil {
			for _, item := range list {
				if strings.Contains(from, ".git") {
					continue
				}
				if e = Copy(filepath.Join(from, item.Name()), filepath.Join(to, item.Name())); e != nil {
					return e
				}
			}
		}
	} else {
		//from是文件，那么创建to的文件夹
		p := filepath.Dir(to)
		if _, e = os.Stat(p); e != nil {
			if e = os.MkdirAll(p, 0777); e != nil {
				return e
			}
		}

		//读取源文件
		file, e := os.Open(from)
		if e != nil {
			return e
		}

		defer file.Close()
		bufReader := bufio.NewReader(file)
		// 创建一个文件用于保存
		out, e := os.Create(to)
		if e != nil {
			return e
		}

		defer out.Close()
		// 然后将文件流和文件流对接起来
		_, e = io.Copy(out, bufReader)
	}
	return e
}

func listFiles(dirname string, FindPath string, bfirst bool) {

	if bfirst {
		fistlocalpath = dirname
	}

	fileInfos, err := os.ReadDir(dirname)
	if err != nil {
		//log.Fatal(err)

		return
	}
	for _, fi := range fileInfos {
		if strings.EqualFold(".git", fi.Name()) {
			continue
		}
		filename := dirname + "/" + fi.Name() //绝对路径带后缀文件名
		Suffix := path.Ext(filename)          //文件后缀名
		//fmt.Printf("Suffix = %v, len(Suffix)=%v\n", Suffix, len(Suffix))
		if len(Suffix) > 0 {
			// MergeFileCnt++
			// fmt.Printf("-")
			// if MergeFileCnt%100 == 0 {
			// 	//fmt.Println("\n进度:", float32(MergeFileCnt/TotleFileNum))
			// 	fmt.Printf("%v%%\n", float32((MergeFileCnt*100.00)/(TotleFileNum*1.0000)))
			// }
			var chinesename, strObjectId string
			chinesename = ""
			if false {
				readXml(filename, &chinesename, &strObjectId)
			}

			filenameall := path.Base(filename)

			// if strings.EqualFold(".uftatomfunction", Suffix) || strings.EqualFold(".uftatomservice", Suffix) {
			// 	FindPath = FindPath + "/uftatom"
			// } else if strings.EqualFold(".uftfunction", Suffix) || strings.EqualFold(".uftservice", Suffix) {
			// 	FindPath = FindPath + "/uftbusiness"
			// }
			srclocalpath := path.Dir(filename)

			RealFindPath := strings.Replace(srclocalpath, fistlocalpath, FindPath, 1)
			//fmt.Printf("%s:%s:%s:%s", filenameall, srclocalpath, RealFindPath, fistlocalpath)
			//fmt.Printf("\n")
			IsFind, DistFilePath, err := FindFile(RealFindPath, filenameall, chinesename, Suffix)
			if err != nil {
				ShowMessage(err.Error())
				break
			}
			if IsFind {
				BakFilePath := strings.Replace(DistFilePath, workpath, bakpath, 1)
				bakisExists, _ := PathExists(BakFilePath)
				// FileisSame := false
				if bakisExists {
					// FileisSame = fileToMd5(DistFilePath) == fileToMd5(BakFilePath)
					DeleteFile(BakFilePath)
					// strWriteMsg := fmt.Sprint("删除备份文件:", BakFilePath, "\n")
					// WriteStingToFile(strWriteMsg, "./删除备份文件汇总.txt")
				}
				//将更新过的文件复制到common中
				// if !bakisExists || !FileisSame {
				// 	err = Copy(DistFilePath, filename)
				// 	if err != nil {
				// 		fmt.Printf("%s:%s", err, DistFilePath)
				// 		// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
				// 		// WriteStingToFile(strWriteMsg, "./执行信息.txt")
				// 	}
				// 	strWriteMsg := fmt.Sprint(DistFilePath, "\n")
				// 	WriteStingToFile(strWriteMsg, "./执行信息.txt")
				// }
				if IsDelete {
					DeleteFile(DistFilePath)
				}

				// strWriteMsg := fmt.Sprint("删除文件:", DistFilePath, "\n")
				// WriteStingToFile(strWriteMsg, "./删除文件汇总.txt")
			} else {
				continue
			}
		}

		if fi.IsDir() {
			//继续遍历fi这个目录
			listFiles(filename, FindPath, false)
		}
	}
}

func listFilesAtom(dirname string, FindPath string, bfirst bool) {

	if bfirst {
		atomfistlocalpath = dirname
	}

	fileInfos, err := os.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fileInfos {
		if strings.EqualFold(".git", fi.Name()) {
			continue
		}
		filename := dirname + "/" + fi.Name() //绝对路径带后缀文件名
		Suffix := path.Ext(filename)          //文件后缀名
		//fmt.Printf("Suffix = %v, len(Suffix)=%v\n", Suffix, len(Suffix))
		if len(Suffix) > 0 {
			// MergeFileCnt++
			// fmt.Printf("-")
			// if MergeFileCnt%100 == 0 {
			// 	//fmt.Println("\n进度:", float32(MergeFileCnt/TotleFileNum))
			// 	fmt.Printf("%v%%\n", float32((MergeFileCnt*100.00)/(TotleFileNum*1.0000)))
			// }
			var chinesename, strObjectId string
			chinesename = ""
			if false {
				readXml(filename, &chinesename, &strObjectId)
			}

			filenameall := path.Base(filename)

			// if strings.EqualFold(".uftatomfunction", Suffix) || strings.EqualFold(".uftatomservice", Suffix) {
			// 	FindPath = FindPath + "/uftatom"
			// } else if strings.EqualFold(".uftfunction", Suffix) || strings.EqualFold(".uftservice", Suffix) {
			// 	FindPath = FindPath + "/uftbusiness"
			// }
			srclocalpath := path.Dir(filename)

			RealFindPath := strings.Replace(srclocalpath, atomfistlocalpath, FindPath, 1)
			//fmt.Printf("%s:%s:%s:%s", filenameall, srclocalpath, RealFindPath, fistlocalpath)
			//fmt.Printf("\n")
			IsFind, DistFilePath, err := FindFile(RealFindPath, filenameall, chinesename, Suffix)
			if err != nil {
				ShowMessage(err.Error())

			}
			if IsFind {
				BakFilePath := strings.Replace(DistFilePath, workpath, bakpath, 1)
				bakisExists, _ := PathExists(BakFilePath)
				//FileisSame := false
				if bakisExists {
					//FileisSame = fileToMd5(DistFilePath) == fileToMd5(BakFilePath)
					DeleteFile(BakFilePath)
					// strWriteMsg := fmt.Sprint("删除备份文件:", BakFilePath, "\n")
					// WriteStingToFile(strWriteMsg, "./删除备份文件汇总.txt")
				}

				// if !bakisExists || !FileisSame {
				// 	err = Copy(DistFilePath, filename)
				// 	if err != nil {
				// 		fmt.Printf("%s:%s", err, DistFilePath)
				// 		// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
				// 		// WriteStingToFile(strWriteMsg, "./执行信息.txt")
				// 	}
				// 	strWriteMsg := fmt.Sprint(DistFilePath, "\n")
				// 	WriteStingToFile(strWriteMsg, "./执行信息.txt")
				// }
				if IsDelete {
					DeleteFile(DistFilePath)
				}
				// strWriteMsg := fmt.Sprint("删除文件:", DistFilePath, "\n")
				// WriteStingToFile(strWriteMsg, "./删除文件汇总.txt")
			} else {
				continue
			}
		}

		if fi.IsDir() {
			//继续遍历fi这个目录
			listFilesAtom(filename, FindPath, false)
		}
	}
}

func listFilesAtom_2(dirname string, FindPath string, bfirst bool) {

	if bfirst {
		atomfistlocalpath_2 = dirname
	}

	fileInfos, err := os.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range fileInfos {
		if strings.EqualFold(".git", fi.Name()) {
			continue
		}
		filename := dirname + "/" + fi.Name() //绝对路径带后缀文件名
		Suffix := path.Ext(filename)          //文件后缀名
		//fmt.Printf("Suffix = %v, len(Suffix)=%v\n", Suffix, len(Suffix))
		if len(Suffix) > 0 {
			// MergeFileCnt++
			// fmt.Printf("-")
			// if MergeFileCnt%100 == 0 {
			// 	//fmt.Println("\n进度:", float32(MergeFileCnt/TotleFileNum))
			// 	fmt.Printf("%v%%\n", float32((MergeFileCnt*100.00)/(TotleFileNum*1.0000)))
			// }
			var chinesename, strObjectId string
			chinesename = ""
			if false {
				readXml(filename, &chinesename, &strObjectId)
			}

			filenameall := path.Base(filename)

			// if strings.EqualFold(".uftatomfunction", Suffix) || strings.EqualFold(".uftatomservice", Suffix) {
			// 	FindPath = FindPath + "/uftatom"
			// } else if strings.EqualFold(".uftfunction", Suffix) || strings.EqualFold(".uftservice", Suffix) {
			// 	FindPath = FindPath + "/uftbusiness"
			// }
			srclocalpath := path.Dir(filename)

			RealFindPath := strings.Replace(srclocalpath, atomfistlocalpath_2, FindPath, 1)
			//fmt.Printf("%s:%s:%s:%s", filenameall, srclocalpath, RealFindPath, fistlocalpath)
			//fmt.Printf("\n")
			IsFind, DistFilePath, err := FindFile(RealFindPath, filenameall, chinesename, Suffix)
			if err != nil {
				ShowMessage(err.Error())
				break
			}
			if IsFind {
				BakFilePath := strings.Replace(DistFilePath, workpath, bakpath, 1)
				bakisExists, _ := PathExists(BakFilePath)
				//FileisSame := false
				if bakisExists {
					//FileisSame = fileToMd5(DistFilePath) == fileToMd5(BakFilePath)
					DeleteFile(BakFilePath)
					// strWriteMsg := fmt.Sprint("删除备份文件:", BakFilePath, "\n")
					// WriteStingToFile(strWriteMsg, "./删除备份文件汇总.txt")
				}

				// if !bakisExists || !FileisSame {
				// 	err = Copy(DistFilePath, filename)
				// 	if err != nil {
				// 		fmt.Printf("%s:%s", err, DistFilePath)
				// 		// strWriteMsg := fmt.Sprint(err, DistFilePath, "\n")
				// 		// WriteStingToFile(strWriteMsg, "./执行信息.txt")
				// 	}
				// 	strWriteMsg := fmt.Sprint(DistFilePath, "\n")
				// 	WriteStingToFile(strWriteMsg, "./执行信息.txt")
				// }
				DeleteFile(DistFilePath)
				// strWriteMsg := fmt.Sprint("删除文件:", DistFilePath, "\n")
				// WriteStingToFile(strWriteMsg, "./删除文件汇总.txt")
			} else {
				continue
			}
		}

		if fi.IsDir() {
			//继续遍历fi这个目录
			listFilesAtom(filename, FindPath, false)
		}
	}
}

// 在指定路径中查找文件名
func FindFile(FindPath string, filenameall string, chinesename string, Suffix string) (bool, string, error) {
	//fmt.Printf("在路径%s中查询文件%s\n", path, filenameall)
	fileInfos, err := os.ReadDir(FindPath)
	if err != nil {
		//log.Fatal(err)
		return false, "", err
	}

	for _, fi := range fileInfos {
		filename := FindPath + "/" + fi.Name() //绝对路径带后缀文件名
		distFileSuffix := path.Ext(filename)
		if Suffix == "" || strings.EqualFold(Suffix, distFileSuffix) {
			fmt.Printf("找到文件：filenameall %s 在路径下： %s\n", filenameall, filename)
			if chinesename != "" {
				var dist_name, DistObjectId string
				readXml(filename, &dist_name, &DistObjectId)
				if strings.EqualFold(chinesename, dist_name) {
					return true, filename, nil
				}
			} else {
				findfilename := path.Base(filename)
				if strings.EqualFold(filenameall, findfilename) {
					return true, filename, nil
				}
			}
		}

		if fi.IsDir() {
			//继续遍历fi这个目录
			bIsFind, findpath, err := FindFile(filename, filenameall, chinesename, Suffix)
			if err != nil {
				return bIsFind, findpath, err
			}
			if bIsFind {
				return bIsFind, findpath, nil
			}
		}
	}
	return false, "", nil
}

func SupportCompareFileType(Suffix string) bool {
	return strings.EqualFold(".uftatomfunction", Suffix) ||
		strings.EqualFold(".uftatomservice", Suffix) ||
		strings.EqualFold(".uftfunction", Suffix) ||
		strings.EqualFold(".uftservice", Suffix) ||
		strings.EqualFold(".uftfactorfunction", Suffix) ||
		strings.EqualFold(".uftfactorservice", Suffix) ||
		strings.EqualFold(".uftstructure", Suffix)
	//strings.EqualFold(".txt", Suffix)
}

func readXml(Strpath string, chinesename *string, strObjectId *string) {

	*chinesename = ""
	*strObjectId = ""

	Suffix := path.Ext(Strpath)
	var rootstr string
	if strings.EqualFold(".uftatomfunction", Suffix) {
		rootstr = "business:Function"
	} else if strings.EqualFold(".uftatomservice", Suffix) {
		rootstr = "business:Service"
	} else if strings.EqualFold(".uftfunction", Suffix) {
		rootstr = "business:Function"
	} else if strings.EqualFold(".uftservice", Suffix) {
		rootstr = "business:Service"
	} else if strings.EqualFold(".uftfactorfunction", Suffix) {
		rootstr = "business:FactorFunction"
	} else if strings.EqualFold(".uftfactorservice", Suffix) {
		rootstr = "business:FactorService"
	} else if strings.EqualFold(".uftstructure", Suffix) {
		rootstr = "structure:Structure"
	} else {
		return
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(Strpath); err != nil {
		panic(err)
	}

	//解析XML根节点
	root := doc.SelectElement(rootstr)
	//fmt.Println("FilePath:", Strpath)
	if nil != root {
		for _, attr := range root.Attr {
			//fmt.Printf("  ATTR: %s=%s\n", attr.Key, attr.Value)
			if strings.EqualFold("chineseName", attr.Key) {
				*chinesename = attr.Value
			} else if strings.EqualFold("objectId", attr.Key) {
				*strObjectId = attr.Value
			}
		}
	}
}

func DeleteFile(FilePath string) {
	err := os.Remove(FilePath)
	if err != nil {
		fmt.Printf("%s:%s", err, FilePath)
	}
}

func WriteStingToFile(writeMsg string, filePath string) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()

	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(writeMsg)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
}

func GetRunPath() (string, error) {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return path, err
}

// 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 判断文件是否相同
func fileToMd5(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		//log.Fatal(err)
		return err.Error()
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		//log.Fatal(err)
		return err.Error()
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (f *TForm1) ShowLog(showstr string) {
	//lock.Lock()
	f.Tips.Lines().Add(showstr)
	//lock.Unlock()
}

func (f *TForm1) StartListening(dir string) {
	//f.ShowLog(dir)

	idExist, _ := PathExists(dir)
	if idExist {
		watch, _ := fsnotify.NewWatcher()
		w := Watch{
			watch: watch}
		w.watchDir(dir)
		//select {}
	}

}

// 监控目录
func (w *Watch) watchDir(dir string) {

	LogFileName := ""
	if strings.Contains(dir, "UFT_equity_Z") {
		LogFileName = fmt.Sprint(workpath + "监控UFT_equity_Z.txt")
	} else if strings.Contains(dir, "UFT_equity_A") {
		LogFileName = fmt.Sprint(workpath + "监控UFT_equity_A.txt")
	}
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		//目录下的文件也在监控范围内，不需要我们一个一个加
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
			//strWriteMsg := fmt.Sprint("监控 : ", path, "\n")
			//WriteStingToFile(strWriteMsg, LogFileName)
		}
		return nil
	})
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {

						if RecordFlag {
							strWriteMsg := fmt.Sprint("创建文件 : ", ev.Name, "\n")
							WriteStingToFile(strWriteMsg, LogFileName)
						}

						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name)
							//strWriteMsg := fmt.Sprint("添加监控 : ", ev.Name, "\n")
							//WriteStingToFile(strWriteMsg, LogFileName)
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						//strWriteMsg := fmt.Sprint("写入文件 : ", ev.Name, "\n")
						//WriteStingToFile(strWriteMsg, LogFileName)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						//strWriteMsg := fmt.Sprint("删除文件 : ", ev.Name, "\n")
						//WriteStingToFile(strWriteMsg, LogFileName)
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Remove(ev.Name)
							//strWriteMsg := fmt.Sprint("删除监控 : ", ev.Name, "\n")
							//WriteStingToFile(strWriteMsg, LogFileName)
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						//strWriteMsg := fmt.Sprint("重命名文件 : ", ev.Name, "\n")
						//WriteStingToFile(strWriteMsg, LogFileName)
						//如果重命名文件是目录，则移除监控
						//注意这里无法使用os.Stat来判断是否是目录了
						//因为重命名后，go已经无法找到原文件来获取信息了
						//所以这里就简单粗爆的直接remove好了
						w.watch.Remove(ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						//strWriteMsg := fmt.Sprint("修改权限 : ", ev.Name, "\n")
						//WriteStingToFile(strWriteMsg, LogFileName)
					}
				}
			case err := <-w.watch.Errors:
				{
					strWriteMsg := fmt.Sprint("error : ", err, "\n")
					WriteStingToFile(strWriteMsg, LogFileName)
					return
				}
			}
		}
	}()
}

func (f *TForm1) testCMD() {

	cmd1 := exec.Command("ping", "www.baidu.com")
	//ppReader, err := cmd1.StdoutPipe()
	// defer ppReader.Close()
	// var bufReader = bufio.NewReader(ppReader)
	// if err != nil {
	// 	fmt.Printf("create cmd stdoutpipe failed,error:%s\n", err)
	// 	os.Exit(1)
	// }
	// err = cmd1.Start()
	// if err != nil {
	// 	fmt.Printf("cannot start cmd1,error:%s\n", err)
	// 	os.Exit(1)
	// }
	out, _ := cmd1.CombinedOutput()
	strWriteMsg := fmt.Sprint("ping: " + string(out))
	f.ShowLog(strWriteMsg)
	// go func() {
	// 	var buffer []byte = make([]byte, 4096)
	// 	for {
	// 		n, err := bufReader.Read(buffer)
	// 		f.Tips.Lines().Add(string(buffer[:n]))
	// 		if err != nil {
	// 			if err == io.EOF {
	// 				f.Tips.Lines().Add("pipi has Closed\n")
	// 				break
	// 			} else {
	// 				f.Tips.Lines().Add("Read content failed")
	// 			}
	// 		}
	// 		f.Tips.Lines().Add(string(buffer[:n]))
	// 	}
	// }()
	// time.Sleep(10 * time.Second)
	// err = stopProcess(cmd1)
	// if err != nil {
	// 	fmt.Printf("stop child process failed,error:%s", err)
	// 	os.Exit(1)
	// }
	// cmd1.Wait()
	// time.Sleep(1 * time.Second)
}

func stopProcess(cmd *exec.Cmd) error {
	pro, err := os.FindProcess(cmd.Process.Pid)
	if err != nil {
		return err
	}
	err = pro.Signal(syscall.SIGINT)
	if err != nil {
		return err
	}
	fmt.Printf("结束子进程%s成功\n", cmd.Path)
	return nil
}

func (f *TForm1) CountMerged(Filename string) {
	go f.CountMergedFunc(Filename)
}
func (f *TForm1) CountMergedFunc(Filename string) {

	FilePath := fmt.Sprint(workpath + Filename)
	err := f.convertFileFormat(FilePath)
	if err != nil {
		f.ShowLog(err.Error())
		return
	}
	time.Sleep(2 * time.Second)
	f.ShowLog(FilePath)
	f.ShowLog("处理开始")

	content := make(map[int]string)
	file, err := os.Open(FilePath)
	if err != nil {
		f.ShowLog("文件打开失败！")
		return
	}
	defer file.Close()

	conut := 0
	//fileAbspath, _ := filepath.Abs(path)
	buf := make([]byte, 1024*1024)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		content[conut] = string(scanner.Text())
		conut++
	}

	mergedCount := 0
	unmergedCount := 0
	for i := 0; i < conut; i++ {
		funcName := content[i]
		if strings.Contains(content[i], "(") {
			innerText := strings.Split(content[i], "(")
			innerText_1 := strings.Split(innerText[1], ")")
			funcName = innerText_1[0]
		}

		if strings.Contains(funcName, "基础公共") || strings.Contains(funcName, "业务公共") {
			continue
		}

		Findpath := ""
		isFind := false
		if strings.Contains(funcName, "AS") || strings.Contains(funcName, "AF") {
			Findpath := workpath + "UFT-Common/UFT-Atom"

			isFind, _, err = FindFile(Findpath, "", funcName, "")
			if err != nil {
				f.ShowLog(err.Error())
				break
			}

		} else {
			Findpath = workpath + "UFT-Common/UFT-Business"

			isFind, _, err = FindFile(Findpath, "", funcName, "")
			if err != nil {
				f.ShowLog(err.Error())
				break
			}
		}

		if isFind {
			mergedCount++
		} else {
			unmergedCount++
		}
		strWriteMsg := fmt.Sprintf("%d-%d", conut, i+1)
		strWriteMsg = fmt.Sprint("(" + strWriteMsg + ") : " + funcName + " : " + strconv.FormatBool(isFind))
		f.ShowLog(strWriteMsg)
	}

	strWriteMsg := fmt.Sprintf("处理完成，总共(%d个),已合并（%d）,未合并（%d）个。", conut, mergedCount, unmergedCount)
	f.ShowLog(strWriteMsg)
}

func (f *TForm1) convertFileFormat(filename string) error {
	filePath := filename // 要检测的文件路径

	// 打开文件并读取所有内容
	file, err := os.Open(filePath)
	if err != nil {
		f.ShowLog("Error opening file:" + string(err.Error()))
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		f.ShowLog("Error reading file:" + string(err.Error()))
		return err
	}

	// 使用 chardet 检测文件编码
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(content)
	if err != nil {
		f.ShowLog("Error detecting encoding:" + string(err.Error()))
		return err
	}
	file.Close()
	if result.Charset != "UTF-8" {
		f.ShowLog("----------文件格式不是UTF-8----------")
		f.ShowLog("----------开始转换文件格式为UTF-8----------")
		originalFile := filename
		tempFile := "temp.txt"
		// 打开原始文件（GB2312 编码）
		file, err := os.Open(originalFile)
		if err != nil {
			fmt.Println("Error opening original file:", err)
			return err
		}
		defer file.Close()

		// 创建临时文件（用于存储 UTF-8 编码内容）
		outFile, err := os.Create(tempFile)
		if err != nil {
			fmt.Println("Error creating temp file:", err)
			return err
		}
		defer outFile.Close()

		// 创建 GBK 解码器（GB2312 是 GBK 的子集）
		reader := transform.NewReader(file, simplifiedchinese.GBK.NewDecoder())
		writer := bufio.NewWriter(outFile)

		// 将内容从 GBK 编码转换为 UTF-8 编码并写入临时文件
		_, err = io.Copy(writer, reader)
		if err != nil {
			fmt.Println("Error converting file:", err)
			return err
		}

		// 确保所有内容写入文件
		writer.Flush()

		// 关闭文件，以确保所有数据都已写入
		outFile.Close()
		file.Close()

		err = os.Remove(originalFile)
		if err != nil {
			fmt.Println("Error deleting temp file:", err)
			return err
		}
		time.Sleep(3 * time.Second)
		err = os.Rename("temp.txt", originalFile)
		if err != nil {
			fmt.Println("Error in renaming")
			return err
		}
		f.ShowLog("----------文件格式转换结束-----------")
		return nil
	} else {
		return nil
	}

}

func (f *TForm1) AutoMoveFromFile(Filename string, dirname string, DistPath string, bfirst bool) {

	FilePath := fmt.Sprint(workpath + Filename)
	f.ShowLog(FilePath)
	fmt.Println("处理开始")
	content := make(map[int]string)
	file, err := os.Open(FilePath)
	if err != nil {
		f.ShowLog("文件打开失败！")
		return
	}
	defer file.Close()

	conut := 0
	//fileAbspath, _ := filepath.Abs(path)
	buf := make([]byte, 1024*1024)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		content[conut] = string(scanner.Text())
		conut++
	}

	for i := 0; i < conut; i++ {
		FileCheName := content[i]
		f.ShowLog(FileCheName)

		realDistPath := ""
		if strings.Contains(FileCheName, "AS") || strings.Contains(FileCheName, "AF") {
			realDistPath = workpath + "UFT-Common/UFT-Atom"
		} else if strings.Contains(FileCheName, "LS") || strings.Contains(FileCheName, "LF") {
			realDistPath = workpath + "UFT-Common/UFT-Business"
		} else if strings.Contains(FileCheName, "RS") || strings.Contains(FileCheName, "RF") {
			realDistPath = workpath + "UFT-Common/UFT-Factor"
		} else {
			realDistPath = workpath + "UFT-Common/UFT-Structure"
		}
		f.AutoMoveFile(FileCheName, dirname, realDistPath, bfirst)
	}

	f.ShowLog("迁移结束了。")
}

func (f *TForm1) CompareFromFile(Filename string) {

	FilePath := fmt.Sprint(workpath + Filename)
	err := f.convertFileFormat(FilePath)
	time.Sleep(2 * time.Second)
	if err != nil {
		f.ShowLog(err.Error())
		return
	}
	f.ShowLog(FilePath)
	fmt.Println("处理开始")
	content := make(map[int]string)
	file, err := os.Open(FilePath)
	if err != nil {
		f.ShowLog("文件打开失败！")
		return
	}
	defer file.Close()

	conut := 0
	//fileAbspath, _ := filepath.Abs(path)
	buf := make([]byte, 1024*1024)
	scanner := bufio.NewScanner(file)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		content[conut] = string(scanner.Text())
		conut++
	}

	for i := 0; i < conut; i++ {
		FileCheName := content[i]
		if strings.Contains(content[i], "(") {
			innerText := strings.Split(content[i], "(")
			innerText_1 := strings.Split(innerText[1], ")")
			FileCheName = innerText_1[0]
		}
		f.ShowLog(FileCheName)
		go f.CompareFile(FileCheName)
	}

	f.ShowLog("比较完成了。")
}

func (f *TForm1) logAnalysis(filename string) {
	f.ShowLog("------开始日志分析------")
	filePath := workpath + filename

	fileBytes, _ := os.ReadFile(filePath)
	Codefmt := f.detectEncoding(fileBytes)
	fmt.Println(Codefmt)

	if Codefmt != "UFT-8" {
		// 创建GBK到UTF-8的编码转换器
		decoder := mahonia.NewDecoder(Codefmt)
		// 进行编码转换
		_, utf8Bytes, err := decoder.Translate(fileBytes, true)
		if err != nil {
			panic(err)
		}
		// 写入UTF-8文件
		err = os.WriteFile("./utf-8.txt", utf8Bytes, 0644)
		if err != nil {
			panic(err)
		}
		filePath = "./utf-8.txt"

	}

	f.ShowLog("处理开始")
	content := make(map[int]string)
	f__, err := os.Open(filePath)
	if err != nil {
		f.ShowLog(err.Error())
		return
	}
	defer f__.Close()
	conut := 0
	//fileAbspath, _ := filepath.Abs(path)
	buf := make([]byte, 10*1024*1024)
	scanner := bufio.NewScanner(f__)
	scanner.Buffer(buf, 10*1024*1024)
	for scanner.Scan() {
		content[conut] = string(scanner.Text())
		conut++
	}

	writeMsg := "功能号,函数名字,调用开始时间,调用结束时间,总耗时(ms)\n"
	//WriteStingToFile(writeMsg, "./执行信息.txt")
	var stringBuilder strings.Builder

	stringBuilder.WriteString("\xEF\xBB\xBF") //防止cvsz中文乱码
	stringBuilder.WriteString(writeMsg)

	totleCount := 0
	search := 0
	f_1, _ := os.Open(filePath)
	scanner_1 := bufio.NewScanner(f_1)
	scanner_1.Buffer(buf, 10*1024*1024)
	for scanner_1.Scan() {
		search++
		if search == conut*20/100 {
			f.ShowLog("----分析进度20%----")
		} else if search == conut*40/100 {
			f.ShowLog("----分析进度40%----")
		} else if search == conut*60/100 {
			f.ShowLog("----分析进度60%----")
		} else if search == conut*80/100 {
			f.ShowLog("----分析进度80%----")
		} else if search == conut {
			f.ShowLog("----分析进度100%----")
		}

		//fmt.Println("文件内容：", scanner_1.Text())
		if strings.Contains(scanner_1.Text(), ",开始") {
			var funcId, FuncName, CallTime, threadId string
			f.GetFuncMsg(scanner_1.Text(), &funcId, &FuncName, &CallTime, &threadId)
			for i := search; i < conut; i++ {
				if strings.Contains(content[i], FuncName+",结束") {
					var funcId_1, FuncName_1, endTime, threadId_1 string
					f.GetFuncMsg(content[i], &funcId_1, &FuncName_1, &endTime, &threadId_1)
					if threadId == threadId_1 {
						DurationTime := f.DurationTimeCalc(CallTime, endTime)
						strDurationTime := fmt.Sprintf("%v", ((float64)(DurationTime.Microseconds()) / 1000.0))
						writeMsg = funcId + "," + FuncName + "," + CallTime + "," + endTime + "," + strDurationTime + "\n"
						//WriteStingToFile(writeMsg, "./执行信息.txt")
						stringBuilder.WriteString(writeMsg)
						totleCount++
						break
					}
				}
			}
		}

	}
	message := "\n文件行数：" + strconv.Itoa(conut)
	f.ShowLog(message)
	message = "统计到服务调用次数：" + strconv.Itoa(totleCount)
	f.ShowLog(message)

	new_filename := "./" + filename + ".csv"
	file, _ := os.OpenFile(new_filename, os.O_RDWR|os.O_CREATE, os.ModeAppend|os.ModePerm)
	defer file.Close()
	dataString := stringBuilder.String()
	file.WriteString(dataString)
	file.Close()

	f__.Close()
	f_1.Close()
	err = os.Remove("utf-8.txt")
	if err != nil {
		f.ShowLog(string(err.Error()))
	}
	f.ShowLog("处理完成")
	f.ShowLog("------日志分析结束------")

}

func (f *TForm1) GetFuncMsg(CallContent string, funcId, FuncName, CallTime, threadId *string) {

	string_Slice := strings.Split(CallContent, "]")
	//fmt.Println("第一行内容分割：", strings.Split(content[0], "]"))
	for i := 0; i < len(string_Slice); i++ {
		//fmt.Println("内容分割：", string_Slice[i])
	}

	if len(string_Slice) > 2 {
		string_Slice_0 := strings.Split(string_Slice[0], " ")
		if len(string_Slice_0) > 2 {
			for i := 0; i < len(string_Slice_0); i++ {
				//fmt.Println("切片0分割：", string_Slice_0[i])
			}

			string_Slice_0_0 := strings.Split(string_Slice_0[0], "-")
			*CallTime = string_Slice_0_0[1]
		}

		//获取funcId
		string_Slice_1 := strings.Split(string_Slice[1], "|")
		if len(string_Slice_1) > 6 {
			for i := 0; i < len(string_Slice_1); i++ {
				//fmt.Println("切片1分割：", string_Slice_1[i])
			}
			*funcId = string_Slice_1[4]
			*threadId = string_Slice_1[6]
		}

		string_Slice_2 := strings.Split(string_Slice[2], ",")
		if len(string_Slice_1) > 6 {
			for i := 0; i < len(string_Slice_1); i++ {
				//fmt.Println("切片2分割：", string_Slice_1[i])
			}
			*FuncName = string_Slice_2[0]
		}

	} else {
		*funcId = ""
		*FuncName = ""
		*CallTime = ""
		*threadId = ""
	}

}

func (f *TForm1) DurationTimeCalc(beginTime, endTime string) time.Duration {

	fbeginTime, _ := strconv.ParseFloat(beginTime, 64)
	fendTime, _ := strconv.ParseFloat(endTime, 64)

	_, beginMis := math.Modf(fbeginTime)
	_, endMis := math.Modf(fendTime)

	time_Slice := strings.Split(beginTime, ".")
	tbeginTime, _ := time.Parse("150405", time_Slice[0])

	time_Slice = strings.Split(endTime, ".")
	tendTime, _ := time.Parse("150405", time_Slice[0])
	subM := tendTime.Sub(tbeginTime)
	//fmt.Println(subM.Minutes(), "分钟")

	Microseconds := (endMis - beginMis) * 1000000
	return subM + time.Duration(Microseconds)*time.Microsecond
}

func (f *TForm1) detectEncoding(file []byte) string {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(file)
	if err != nil {
		panic(err)
	}
	return result.Charset
}

func debugLog(message string) {
	message += "\n"

	// 文件路径
	filePath := "debugLog.txt"

	// 使用os.OpenFile打开文件，如果不存在则创建
	// os.O_APPEND | os.O_CREATE | os.O_WRONLY 用于追加内容，如果文件不存在则创建
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 写入字符串
	_, err = file.WriteString(message)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

}

func IntPtr(n int) uintptr {
	return uintptr(n)
}
func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}

// ShowMessage windows下的另一种DLL方法调用
func ShowMessage(text string) {
	user32dll, _ := syscall.LoadLibrary("user32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")
	MessageBoxW := user32.NewProc("MessageBoxW")
	MessageBoxW.Call(IntPtr(0), StrPtr(text), IntPtr(0))
	defer syscall.FreeLibrary(user32dll)
}
func (f *TForm1) setOwner(path string, owner string) error {
	f.ShowLog("-------开始更改所有者----------")

	// 更改文件所有者
	icaclsCmd := exec.Command("icacls", path, "/setowner", owner, "/t", "/c", "/l", "/q")
	icaclsCmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true, // 创建一个新的会话组
	}
	output, err := icaclsCmd.CombinedOutput()
	if err != nil {
		f.ShowLog(fmt.Sprintf("Error setting owner for %s: %v\nOutput: %s", path, err, string(output)))
		return err
	}

	f.ShowLog("-------更改所有者结束----------")

	return nil
}
