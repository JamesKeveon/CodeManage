unit Form;

{$mode objfpc}{$H+}

interface

uses
  Classes, SysUtils, Forms, Controls, Graphics, Dialogs, StdCtrls, Grids,
  ActnList, ComboEx, ExtCtrls, CheckLst, ColorBox, DBCtrls, Menus, ShellCtrls,
  uRichEdit, ComCtrls;

type

  { TForm1 }

  TForm1 = class(TForm)
    Action1: TAction;
    ActionList1: TActionList;
    BtnCodeDownload: TButton;
    BtnCodeUpdate: TButton;
    BtnFuncCompare: TButton;
    BtnFileAutoMove: TButton;
    BtnO45Develop: TButton;
    BtnAM4Develop: TButton;
    BtnO45Commit: TButton;
    BtnAM4Commit: TButton;
    Button1: TButton;
    BtnClearShow: TButton;
    AutoMoveFromAM4: TCheckBox;
    BtnCreateBrch: TButton;
    BtnSwithTrunk: TButton;
    Button2: TButton;
    Button3: TButton;
    Button4: TButton;
    ChB_Basic: TCheckBox;
    ChB_Businpub: TCheckBox;
    CheB_Equity: TCheckBox;
    ChB_Equity_A: TCheckBox;
    ChB_Equity_Comm: TCheckBox;
    CheckBox1: TCheckBox;
    CheckBox2: TCheckBox;
    CheckBox3: TCheckBox;
    CheckBox4: TCheckBox;
    CheckOpLS: TCheckBox;
    ComboBox1: TComboBox;
    EdtReqNum: TEdit;
    EdtVersion: TEdit;
    FuncCount: TButton;
    EditUserID: TEdit;
    EditPasswd: TEdit;
    EdtFuncName: TEdit;
    Label11: TLabel;
    Label12: TLabel;
    Label2: TLabel;
    Label3: TLabel;
    Label4: TLabel;
    Label5: TLabel;
    Label6: TLabel;
    Label7: TLabel;
    Label8: TLabel;
    Label9: TLabel;
    SelectDirectoryDialog1: TSelectDirectoryDialog;
    Tips: TRichEdit;
    EdtWorkpath: TEdit;
    Label1: TLabel;
    procedure Action1Execute(Sender: TObject);
    procedure Action1Hint(var HintStr: string; var CanShow: Boolean);
    procedure Action1Update(Sender: TObject);
    procedure ActionList1Change(Sender: TObject);
    procedure ActionList1Execute(AAction: TBasicAction; var Handled: Boolean);
    procedure ActionList1Update(AAction: TBasicAction; var Handled: Boolean);
    procedure BtnAM4CommitClick(Sender: TObject);
    procedure BtnAM4DevelopClick(Sender: TObject);
    procedure BtnClearShowClick(Sender: TObject);
    procedure BtnCodeDownloadClick(Sender: TObject);
    procedure BtnCodeUpdateClick(Sender: TObject);
    procedure BtnCreateBrchClick(Sender: TObject);
    procedure BtnFuncCompareClick(Sender: TObject);
    procedure BtnFileAutoMoveClick(Sender: TObject);
    procedure BtnO45CommitClick(Sender: TObject);
    procedure BtnO45DevelopClick(Sender: TObject);
    procedure BtnO45DevelopMouseEnter(Sender: TObject);
    procedure BtnSwithTrunkClick(Sender: TObject);
    procedure Button1Click(Sender: TObject);
    procedure Button2Click(Sender: TObject);
    procedure Button3Click(Sender: TObject);
    procedure Button4Click(Sender: TObject);
    procedure ChB_BasicChange(Sender: TObject);
    procedure ChB_BusinpubChange(Sender: TObject);
    procedure ChB_Equity_AChange(Sender: TObject);
    procedure ChB_Equity_CommChange(Sender: TObject);
    procedure CheB_EquityChange(Sender: TObject);
    procedure CheckBox1Change(Sender: TObject);
    procedure CheckBox2Change(Sender: TObject);
    procedure CheckBox3Change(Sender: TObject);
    procedure CheckBox4Change(Sender: TObject);
    procedure CheckBox5Change(Sender: TObject);
    procedure CheckGroup1Click(Sender: TObject);
    procedure CheckListBox1ItemClick(Sender: TObject; Index: integer);
    procedure CheckOpTrunkChange(Sender: TObject);
    procedure CheckOpLSChange(Sender: TObject);
    procedure ComboBox1Change(Sender: TObject);
    procedure ComboBoxEx1Change(Sender: TObject);
    procedure Label10Click(Sender: TObject);
    procedure ListBox1Click(Sender: TObject);
    procedure ListBox2Click(Sender: TObject);
    procedure RadioGroup1Click(Sender: TObject);
    procedure ScrollBox1Click(Sender: TObject);
    procedure SelectDirectoryDialog1Close(Sender: TObject);
    procedure SelectDirectoryDialog1SelectionChange(Sender: TObject);
    procedure ShellTreeView1Change(Sender: TObject; Node: TTreeNode);
    procedure testChange(Sender: TObject);
    procedure EdtReqNumChange(Sender: TObject);
    procedure EdtVersionChange(Sender: TObject);
    procedure EditUserIDChange(Sender: TObject);
    procedure EdtFuncNameChange(Sender: TObject);
    procedure EdtWorkpathChange(Sender: TObject);
    procedure FormCreate(Sender: TObject);
    procedure FuncCountClick(Sender: TObject);
    procedure Label6Click(Sender: TObject);
    procedure Label9Click(Sender: TObject);
    procedure Memo1Change(Sender: TObject);
    procedure Automovefromam4Change(Sender: TObject);
    procedure RichEdit1Change(Sender: TObject);
    procedure TipsChange(Sender: TObject);
    procedure TipsClick(Sender: TObject);
  private

  public

  end;

var
  Form1: TForm1;

implementation

{$R *.lfm}

{ TForm1 }

procedure TForm1.FormCreate(Sender: TObject);
begin

end;

procedure TForm1.FuncCountClick(Sender: TObject);
begin

end;

procedure TForm1.Label6Click(Sender: TObject);
begin

end;

procedure TForm1.Label9Click(Sender: TObject);
begin

end;

procedure TForm1.Memo1Change(Sender: TObject);
begin

end;

procedure TForm1.Automovefromam4Change(Sender: TObject);
begin

end;

procedure TForm1.BtnCodeDownloadClick(Sender: TObject);
begin

end;

procedure TForm1.BtnAM4CommitClick(Sender: TObject);
begin

end;

procedure TForm1.ActionList1Change(Sender: TObject);
begin

end;

procedure TForm1.ActionList1Execute(AAction: TBasicAction; var Handled: Boolean
  );
begin

end;

procedure TForm1.ActionList1Update(AAction: TBasicAction; var Handled: Boolean);
begin

end;

procedure TForm1.Action1Execute(Sender: TObject);
begin

end;

procedure TForm1.Action1Hint(var HintStr: string; var CanShow: Boolean);
begin

end;

procedure TForm1.Action1Update(Sender: TObject);
begin

end;

procedure TForm1.BtnAM4DevelopClick(Sender: TObject);
begin

end;

procedure TForm1.BtnClearShowClick(Sender: TObject);
begin

end;

procedure TForm1.BtnCodeUpdateClick(Sender: TObject);
begin

end;

procedure TForm1.BtnCreateBrchClick(Sender: TObject);
begin

end;

procedure TForm1.BtnFuncCompareClick(Sender: TObject);
begin

end;

procedure TForm1.BtnFileAutoMoveClick(Sender: TObject);
begin

end;

procedure TForm1.BtnO45CommitClick(Sender: TObject);
begin

end;

procedure TForm1.BtnO45DevelopClick(Sender: TObject);
begin

end;

procedure TForm1.BtnO45DevelopMouseEnter(Sender: TObject);
begin

end;

procedure TForm1.BtnSwithTrunkClick(Sender: TObject);
begin

end;

procedure TForm1.Button1Click(Sender: TObject);
begin

end;

procedure TForm1.Button2Click(Sender: TObject);
begin

end;

procedure TForm1.Button3Click(Sender: TObject);
begin

end;

procedure TForm1.Button4Click(Sender: TObject);
begin

end;

procedure TForm1.ChB_BasicChange(Sender: TObject);
begin

end;

procedure TForm1.ChB_BusinpubChange(Sender: TObject);
begin

end;

procedure TForm1.ChB_Equity_AChange(Sender: TObject);
begin

end;

procedure TForm1.ChB_Equity_CommChange(Sender: TObject);
begin

end;

procedure TForm1.CheB_EquityChange(Sender: TObject);
begin

end;

procedure TForm1.CheckBox1Change(Sender: TObject);
begin

end;

procedure TForm1.CheckBox2Change(Sender: TObject);
begin

end;

procedure TForm1.CheckBox3Change(Sender: TObject);
begin

end;

procedure TForm1.CheckBox4Change(Sender: TObject);
begin

end;

procedure TForm1.CheckBox5Change(Sender: TObject);
begin

end;

procedure TForm1.CheckGroup1Click(Sender: TObject);
begin

end;

procedure TForm1.CheckListBox1ItemClick(Sender: TObject; Index: integer);
begin

end;

procedure TForm1.CheckOpTrunkChange(Sender: TObject);
begin

end;

procedure TForm1.CheckOpLSChange(Sender: TObject);
begin

end;

procedure TForm1.ComboBox1Change(Sender: TObject);
begin

end;

procedure TForm1.ComboBoxEx1Change(Sender: TObject);
begin

end;

procedure TForm1.Label10Click(Sender: TObject);
begin

end;

procedure TForm1.ListBox1Click(Sender: TObject);
begin

end;

procedure TForm1.ListBox2Click(Sender: TObject);
begin

end;

procedure TForm1.RadioGroup1Click(Sender: TObject);
begin

end;

procedure TForm1.ScrollBox1Click(Sender: TObject);
begin

end;

procedure TForm1.SelectDirectoryDialog1Close(Sender: TObject);
begin

end;

procedure TForm1.SelectDirectoryDialog1SelectionChange(Sender: TObject);
begin

end;

procedure TForm1.ShellTreeView1Change(Sender: TObject; Node: TTreeNode);
begin

end;

procedure TForm1.testChange(Sender: TObject);
begin

end;

procedure TForm1.EdtReqNumChange(Sender: TObject);
begin

end;

procedure TForm1.EdtVersionChange(Sender: TObject);
begin

end;

procedure TForm1.AutoMoveFromAM4Change(Sender: TObject);
begin

end;

procedure TForm1.EditUserIDChange(Sender: TObject);
begin

end;

procedure TForm1.EdtFuncNameChange(Sender: TObject);
begin

end;

procedure TForm1.EdtWorkpathChange(Sender: TObject);
begin

end;

procedure TForm1.RichEdit1Change(Sender: TObject);
begin

end;

procedure TForm1.TipsChange(Sender: TObject);
begin

end;

procedure TForm1.TipsClick(Sender: TObject);
begin

end;

end.

