package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/ying32/govcl/vcl"
)

// ::private::
type TForm1Fields struct {
}

func (f *TForm1) OnFormCreate(sender vcl.IObject) {
	workpath, _ := GetRunPath()
	f.ShowLog(workpath)
	workpath += "\\"
	bakpath = workpath + ".cache\\"
	f.EdtWorkpath.SetText(workpath)
	f.Tips.SetWordWrap(true)
	f.Tips.SetText("")
	f.SetShowHint(true)
	f.StartListening(workpath + "UFT_equity_Z\\UFTDB_equity\\")
	f.StartListening(workpath + "UFT_equity_A\\UFTDB_equity\\")

	IsMovefromAM4 = false
	IsPullLS = false
	IsPullTrunk = false
	IsDelete = true
	linkA_1 = false
	linkA = false
	linkA_2 = false
	linkZ = false
	linkZ_1 = false
	linkZ_2 = false
	atom_1 = false
	atom_2 = false
	isDownload = false
	isPullFlag = true
	f.Button1.SetAction(f.Action1)
	f.Button1.Hide()

	// 从文件中读取选项并添加到 ComboBox
	f.LoadOptionsFromFile("options.xml")

}

func (f *TForm1) OnEdit1Change(sender vcl.IObject) {

}

func (f *TForm1) OnBtnCodeDownloadClick(sender vcl.IObject) {

	// if f.EditUserID.Text() == "" || f.EditPasswd.Text() == "" {
	// 	vcl.ShowMessage("代码下载必须输入用户名，密码！")
	// 	return
	// }

	root := f.ComboBox1.Text()

	if root == "" && !IsPullLS {
		vcl.ShowMessage("请在下拉框中选择或勾选04LS！")
		return
	}
	if root != "" && IsPullLS {
		vcl.ShowMessage("请不要同时勾选两个版本！")
		return
	} else if root != "" || IsPullLS {
		//do nothing
	} else if !IsPullLS && !IsPullTrunk {
		vcl.ShowMessage("请选择要下载或者更新的版本！")
		return
	} else if IsPullLS && IsPullTrunk {
		vcl.ShowMessage("请不要同时选择两个下载/更新版本！")
		return
	}
	workpath = f.EdtWorkpath.Text()
	if len(workpath) > 0 {
		if workpath[len(workpath)-1] != '\\' {
			workpath += "\\"
		}
	}
	//f.Tips.Lines().Add("开始代码下载")
	f.CodeDownLoad(false)

}

func (f *TForm1) OnBtnCodeUpdateClick(sender vcl.IObject) {
	//f.Tips.Lines().Add("开始代码更新")
	root := f.ComboBox1.Text()
	if root == "" && !IsPullLS {
		vcl.ShowMessage("请在下拉框中选择或勾选04LS！")
		return
	}
	if root != "" && IsPullLS {
		vcl.ShowMessage("请不要同时勾选两个版本！")
		return
	} else if root != "" || IsPullLS {
		//do nothing
	} else if !IsPullLS && !IsPullTrunk {
		vcl.ShowMessage("请选择要下载或者更新的版本！")
		return
	} else if IsPullLS && IsPullTrunk {
		vcl.ShowMessage("请不要同时选择两个下载/更新版本！")
		return
	}
	workpath = f.EdtWorkpath.Text()
	if len(workpath) > 0 {
		if workpath[len(workpath)-1] != '\\' {
			workpath += "\\"
		}
	}
	isPullFlag = true
	f.CodeDownLoad(true)
}

func (f *TForm1) OnTipsChange(sender vcl.IObject) {

}

func (f *TForm1) OnTipsClick(sender vcl.IObject) {

}

func (f *TForm1) OnBtnO45DevelopClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	go f.SymbolicLinkExecute("2")
}

func (f *TForm1) OnBtnAM4DevelopClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	go f.SymbolicLinkExecute("1")
}

func (f *TForm1) OnBtnO45CommitClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	go f.CodeCommit("2")
}

func (f *TForm1) OnBtnAM4CommitClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	go f.CodeCommit("1")
}

func (f *TForm1) OnBtnFileAutoMoveClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	FileCheName := f.EdtFuncName.Text()
	if strings.Contains(FileCheName, ".txt") {
		FromMovePath := workpath + "UFT_equity_Z"
		if IsMovefromAM4 {
			FromMovePath = workpath + "UFT_equity_A"
		}
		go f.AutoMoveFromFile(FileCheName, FromMovePath, "", true)
	} else {
		DistPath := ""
		if strings.Contains(FileCheName, "AS") || strings.Contains(FileCheName, "AF") {
			DistPath = workpath + "UFT-Common/UFT-Atom"
		} else if strings.Contains(FileCheName, "LS") || strings.Contains(FileCheName, "LF") {
			DistPath = workpath + "UFT-Common/UFT-Business"
		} else if strings.Contains(FileCheName, "RS") || strings.Contains(FileCheName, "RF") {
			DistPath = workpath + "UFT-Common/UFT-Factor"
		} else {
			DistPath = workpath + "UFT-Common/UFT-Structure"
		}
		//迁移资管银信从此处修改
		if IsMovefromAM4 {
			go f.AutoMoveFile(FileCheName, workpath+"UFT_equity_A", DistPath, true)
		} else {
			go f.AutoMoveFile(FileCheName, workpath+"UFT_equity_Z", DistPath, true)
		}
	}

}

func (f *TForm1) OnEdtFuncNameChange(sender vcl.IObject) {

}

func (f *TForm1) OnButton1Click(sender vcl.IObject) {
	f.testCMD()
}

func (f *TForm1) OnBtnO45DevelopMouseEnter(sender vcl.IObject) {
	f.BtnO45Develop.SetHint("开发模式是将UFT-Common/UFT-Metadata、UFT-Common/UFT-Structure、UFT-Common/UFT-Atom、UFT-Common/UFT-Factor、UFT-Common/UFT-Business拷贝到原始工程下.")
	//f.Tips.Lines().Add("开发模式是将UFT-Metadata、UFT-Structure、UFT-Atom、UFT-Factor、UFT-Business拷贝到原始工程下.")
}

func (f *TForm1) OnAction1Execute(sender vcl.IObject) {
	vcl.ShowMessage("点击了测试按钮")
}

func (f *TForm1) OnAction1Update(sender vcl.IObject) {
	//vcl.ShowMessage("OnAction1Update")
}

func (f *TForm1) OnEditUserIDChange(sender vcl.IObject) {

}

func (f *TForm1) OnEdtWorkpathChange(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	workpath += "/"
	strWriteMsg := fmt.Sprint("工作目录修改成：", workpath)
	f.Tips.Lines().Add(strWriteMsg)
}

func (f *TForm1) OnBtnClearShowClick(sender vcl.IObject) {
	f.Tips.SetText("")
}

func (f *TForm1) OnAutoMoveFromAM4Change(sender vcl.IObject) {
	IsMovefromAM4 = !IsMovefromAM4
	//f.Automovefromam4.SetChecked(IsMovefromAM4)\
	if IsMovefromAM4 {
		f.Tips.Lines().Add("从银信版迁移")
	} else {
		f.Tips.Lines().Add("从资管版迁移")
	}

}

func (f *TForm1) OnCheckOpLSChange(sender vcl.IObject) {
	IsPullLS = !IsPullLS
	//f.Automovefromam4.SetChecked(IsMovefromAM4)\
	if IsPullLS {
		f.Tips.Lines().Add("下载/更新04补丁版本代码")
	} else {
		f.Tips.Lines().Add("取消下载/更新04补丁版本代码")
	}
}

func (f *TForm1) OnCheckOpTrunkChange(sender vcl.IObject) {
	IsPullTrunk = !IsPullTrunk
	if IsPullTrunk {
		f.Tips.Lines().Add("下载/更新最新版本代码")
	} else {
		f.Tips.Lines().Add("取消下载/更新最新版本代码")
	}
}

func (f *TForm1) OnBtnCreateBrchClick(sender vcl.IObject) {
	go f.osExecCreateFeatures(true)
}

func (f *TForm1) OnBtnSwithTrunkClick(sender vcl.IObject) {
	go f.osExecCreateFeatures(false)
}

func (f *TForm1) OnChB_BasicChange(sender vcl.IObject) {
	Featrure_basic = !Featrure_basic
	if Featrure_basic {
		f.Tips.Lines().Add("拉取（切换）basic分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）basic分支")
	}
}

func (f *TForm1) OnEdit2Change(sender vcl.IObject) {

}

func (f *TForm1) OnChB_BusinpubChange(sender vcl.IObject) {
	Featrure_businpub = !Featrure_businpub
	if Featrure_businpub {
		f.Tips.Lines().Add("拉取（切换）businpub分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）businpub分支")
	}
}

func (f *TForm1) OnCheB_EquityChange(sender vcl.IObject) {
	Featrure_equity = !Featrure_equity
	if Featrure_equity {
		f.Tips.Lines().Add("拉取（切换）equity分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）equity分支")
	}
}

func (f *TForm1) OnChB_Equity_AChange(sender vcl.IObject) {
	Featrure_equity_A = !Featrure_equity_A
	if Featrure_equity_A {
		f.Tips.Lines().Add("拉取（切换）equity_A分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）equity_A分支")
	}
}

func (f *TForm1) OnChB_Equity_CommChange(sender vcl.IObject) {
	Featrure_equity_Common = !Featrure_equity_Common
	if Featrure_equity_Common {
		f.Tips.Lines().Add("拉取（切换）equity_Common分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）equity_Common分支")
	}
}

func (f *TForm1) OnLabel6Click(sender vcl.IObject) {

}

//	func (f *TForm1) OnEdtReqNumChange(sender vcl.IObject) {
//		RequirementNum = f.EdtReqNum.Text()
//		strWriteMsg := fmt.Sprint("需求编号是：", RequirementNum)
//		f.Tips.Lines().Add(strWriteMsg)
//	}
func (f *TForm1) OnEdtReqNumExit(sender vcl.IObject) {
	RequirementNum = f.EdtReqNum.Text()
	strWriteMsg := fmt.Sprint("需求编号是：", RequirementNum)
	f.Tips.Lines().Add(strWriteMsg)
}

func (f *TForm1) OnEdtVersionExit(sender vcl.IObject) {
	Version = f.EdtVersion.Text()
	cutVersion := Version
	if strings.Contains(Version, ".") {
		innerText := strings.Split(Version, ".")
		cutVersion = innerText[1]
	}
	baisc_Version = "IPS1.0-basicV202401." + cutVersion + ".000"
	businpub_Version = "IPS1.0-businpubV202401." + cutVersion + ".000"
	equity_Verison = "IPS1.0-equityV202401Z." + cutVersion + ".000"
	equity_A_Version = "IPS1.0-equityV202401A." + cutVersion + ".000"
	equity_Comm_Version = "IPS1.0-equityV202401" + Version + ".000"

	equity_server_Version = "IPS1.0-equityV202401Z." + cutVersion + ".000"
	equity_server_common_Version = "IPS1.0-equityV202401" + Version + ".000"
	basic_server_Version = "IPS1.0-basicV202401." + cutVersion + ".000"
	businpub_server_Version = "IPS1.0-businpubV202401." + cutVersion + ".000"
	strWriteMsg := fmt.Sprint("版本是：\n", "basic:         "+baisc_Version+"\n"+
		"businpub:      "+businpub_Version+"\n"+
		"equity:        "+equity_Verison+"\n"+
		"equity_A:     "+equity_A_Version+"\n"+
		"equity_Common: "+equity_Comm_Version+"\n"+
		"equity_server:"+equity_server_Version+"\n"+
		"equity_server_common"+equity_server_common_Version+"\n"+
		"basic_server"+basic_server_Version+"\n"+
		"businpub_server"+businpub_server_Version+"\n")
	f.Tips.Lines().Add(strWriteMsg)
}

func (f *TForm1) OnEdtVersionChange(sender vcl.IObject) {

}

func (f *TForm1) OnLabel9Click(sender vcl.IObject) {

}

func (f *TForm1) LoadOptionsFromFile(filename string) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		vcl.ShowMessage(fmt.Sprintf("无法打开文件: %s", err.Error()))
		return
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var rootLevel bool

	for {
		t, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == "root" {
				rootLevel = true
			} else if rootLevel && se.Name.Local != "Variable" {
				// 只添加与根元素同一级的节点名称到 ComboBox
				f.ComboBox1.Items().Add(se.Name.Local)
			}
		case xml.EndElement:
			if se.Name.Local == "root" {
				rootLevel = false
			}
		}
	}
	f.ComboBox1.Items().Add("")
}

func (f *TForm1) OnComboBox1Change(sender vcl.IObject) {
	comboBox := vcl.AsComboBox(sender)
	if comboBox != nil {
		selectedText := comboBox.Text()
		codepath, err := f.GetCodePath(selectedText)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(codepath) != 0 {
			workpath = codepath
			if len(workpath) > 0 {
				if workpath[len(workpath)-1] != '\\' {
					workpath += "\\"
				}
			}
			fmt.Println("workpath:" + workpath)
			f.EdtWorkpath.SetText(workpath)
			fmt.Println("EdtWorkpath:" + f.EdtWorkpath.Text())
		}

		message := "选择了一个分支: " + selectedText + "，工作路径：" + workpath
		f.Tips.Lines().Add(message)
	} else {
		fmt.Println("无法将 sender 转换为 TComboBox")
	}
}

func (f *TForm1) OnCheckBox1Change(sender vcl.IObject) {
	Feature_equity_server = !Feature_equity_server
	if Feature_equity_server {
		f.Tips.Lines().Add("拉取（切换）equity分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）equity分支")
	}
}

func (f *TForm1) OnCheckBox2Change(sender vcl.IObject) {
	Feature_equity_server_common = !Feature_equity_server_common
	if Feature_equity_server_common {
		f.Tips.Lines().Add("拉取（切换）common分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）common分支")
	}
}

func (f *TForm1) OnCheckBox3Change(sender vcl.IObject) {
	Feature_basic_server = !Feature_basic_server
	if Feature_basic_server {
		f.Tips.Lines().Add("拉取（切换）basic分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）basic分支")
	}
}

func (f *TForm1) OnCheckBox4Change(sender vcl.IObject) {
	Feature_businpub_server = !Feature_businpub_server
	if Feature_businpub_server {
		f.Tips.Lines().Add("拉取（切换）businpub分支")
	} else {
		f.Tips.Lines().Add("取消拉取（切换）businpub分支")
	}
}

func (f *TForm1) OnCheckBox5Change(sender vcl.IObject) {
	IsDelete = !IsDelete
	if IsDelete {
		f.Tips.Lines().Add("提交时删除equity-AZ中部分代码")
	} else {
		f.Tips.Lines().Add("提交时保留equity-AZ中部分代码")
	}
}

func (f *TForm1) OnFuncCountClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	CountFileName := f.EdtFuncName.Text()
	if len(workpath) > 0 {
		if workpath[len(workpath)-1] != '\\' {
			workpath += "\\"
		}
	}
	f.CountMerged(CountFileName)
}

func (f *TForm1) OnBtnFuncCompareClick(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	CompareCheName := f.EdtFuncName.Text()
	if len(workpath) > 0 {
		if workpath[len(workpath)-1] != '\\' {
			workpath += "\\"
		}
	}
	if strings.Contains(CompareCheName, ".txt") {
		IsBatchCompare = true
		go f.CompareFromFile(CompareCheName)
	} else {
		IsBatchCompare = false
		go f.CompareFile(CompareCheName)
	}
}

func (f *TForm1) OnButton2Click(sender vcl.IObject) {
	workpath = f.EdtWorkpath.Text()
	CountFileName := f.EdtFuncName.Text()
	// if CountFileName == "" || CountFileName == " " {
	// 	f.ShowLog("请输入文件名")
	// 	return
	// }
	if len(workpath) > 0 {
		if workpath[len(workpath)-1] != '\\' {
			workpath += "\\"
		}
	}
	go f.logAnalysis(CountFileName)
}

func (f *TForm1) OnScrollBox1Click(sender vcl.IObject) {

}

func (f *TForm1) OnListBox1Click(sender vcl.IObject) {

}

func (f *TForm1) OnCheckListBox1ItemClick(sender vcl.IObject, index int32) {

}

func (f *TForm1) OnSelectDirectoryDialog1Close(sender vcl.IObject) {

}

func (f *TForm1) OnButton3Click(sender vcl.IObject) {
	f.ShowLog(workpath)
	f.SelectDirectoryDialog1.SetFileName(workpath[:len(workpath)-1])
	if f.SelectDirectoryDialog1.Execute() {

		path := f.SelectDirectoryDialog1.FileName()
		message := "选择了路径： " + path
		f.ShowLog(message)
		f.EdtWorkpath.SetText(path)
	}
}

func (f *TForm1) OnSelectDirectoryDialog1SelectionChange(sender vcl.IObject) {

}

func (f *TForm1) OnDirectoryButtonClick(sender vcl.IObject) {

}

func (f *TForm1) OnButton4Click(sender vcl.IObject) {

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
	f.ShowLog(currentUser.Username)
	go f.setOwner(workpath, currentUser.Username)

}
