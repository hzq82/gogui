package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"myApp/note"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/text/encoding/simplifiedchinese"
)

func main() {
	os.Setenv("FYNE_FONT", "AlibabaPuHuiTi-3-85-Bold.ttf")
	//fmt.Println("Test Fyne")
	app := app.New()
	//new window and title
	w := app.NewWindow("OPS")
	//set app ICON
	r, _ := fyne.LoadResourceFromPath("f:/go/gogui/sxs.ico")
	w.SetIcon(r)
	//w.SetIcon(theme.FyneLogo())
	//resize window
	w.Resize(fyne.NewSize(500, 500))
	//设置查询ip地址的显示label
	lTitle := widget.NewLabel("您的互联网IP地址如下：")
	//IpconfigTitle := widget.NewLabel("本地ip信息：")
	IpconfigEntry := widget.NewMultiLineEntry()
	lTitle.TextStyle = fyne.TextStyle{Bold: true}
	lIp := widget.NewLabel("")
	lCountry := widget.NewLabel("")
	lCity := widget.NewLabel("")
	logininfo := Logininfo{
		Username: "root",
		Password: "1234qwer",
		SshPort:  22}
	//菜单栏ip查询菜单
	menuItem1 := fyne.NewMenuItem("IP查询", func() {
		lIp.Text = note.MyIP()["query"]
		lCountry.Text = note.MyIP()["country"]
		lCity.Text = note.MyIP()["city"]
		mpuIP := container.NewVBox(
			lTitle,
			lIp,
			lCountry,
			lCity,
			//btnQip,
		)
		mpuQip := widget.NewPopUp(mpuIP, w.Canvas())
		mpuQip.ShowAtPosition(fyne.NewPos(
			w.Canvas().Size().Width/2-mpuIP.MinSize().Width/2,
			w.Canvas().Size().Height/2-mpuIP.MinSize().Height/2,
		))
		//mpuQip.Resize(w.Canvas().Size())
		mpuQip.Show()
	})
	menuItem2 := fyne.NewMenuItem("Ipconfig", func() {
		//fmt.Println("Save pressed")
		cmd := exec.Command("ipconfig", "/all")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalln(err)
		}

		//log.Println("result =>", ConvertByte2String(output, GB18030))
		IpconfigEntry.Text = ConvertByte2String(output, GB18030)

		conIpconfig := container.New(layout.NewGridLayoutWithRows(1), IpconfigEntry)

		mpuIpconfig := widget.NewPopUp(conIpconfig, w.Canvas())
		//mpuIpconfig.Resize(fyne.NewSize(400, 400))
		mpuIpconfig.ShowAtPosition(fyne.NewPos(
			w.Canvas().Size().Width/2-conIpconfig.MinSize().Width/2,
			w.Canvas().Size().Height/2-conIpconfig.MinSize().Height/2,
		))
		//mpuIpconfig.Resize(w.Content().Size())
		mpuIpconfig.Resize(w.Canvas().Size())
		mpuIpconfig.Show()
	})
	menuItem3 := fyne.NewMenuItem("About", func() {
		// w2 := app.NewWindow("About")
		// w2.Resize(fyne.NewSize(200, 200))
		// w2.Show()
		mpuAbout := container.NewVBox(
			widget.NewLabel("软件版本:v0.0.1\n作者:haozq"),
		)
		mpu := widget.NewPopUp(mpuAbout, w.Canvas())
		//mpu.Resize(fyne.NewSize(200, 200))
		mpu.ShowAtPosition(fyne.NewPos(
			w.Canvas().Size().Width/2-mpu.MinSize().Width/2,
			w.Canvas().Size().Height/2-mpu.MinSize().Height/2,
		))
		fmt.Printf("%f", mpu.MinSize().Width)
		mpu.Show()
	})
	menuItem4 := fyne.NewMenuItem("LogClean", func() {
		ipLable := widget.NewLabel("日志清理")
		//logCC := container.New(layout.NewVBoxLayout(), ipLable, layout.NewSpacer())
		//logCC := container.NewHBox(ipLable, layout.NewSpacer(), ipLable)
		//layout.NewSpacer()
		var server string
		var cleanPopup *widget.PopUp
		logCC := container.NewAdaptiveGrid(3, ipLable)
		ipLableR := widget.NewLabel("")
		IPSelect := widget.NewSelect([]string{"192.168.102.50"}, func(s string) {
			server = s
			ipLableR.Text = "选择的服务器：" + s
			ipLableR.Refresh()
		})
		//cleanEntry := widget.NewMultiLineEntry()
		cleanBtn := widget.NewButton("执行", func() {
			openChildWindow(app, server)
		})
		quitBtn := widget.NewButton("退出", func() {
			//app.Quit()
			cleanPopup.Hide()
		})
		cleanCon := container.NewVBox(logCC, IPSelect, ipLableR, cleanBtn, quitBtn, layout.NewSpacer())
		cleanPopup = widget.NewPopUp(cleanCon, w.Canvas())
		cleanPopup.ShowAtPosition(fyne.NewPos(
			w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
			w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
		))
		cleanPopup.Show()
	})
	//192.168.102.40 xxl-job 重启

	menuItem5 := fyne.NewMenuItem("xxl_Restart", func() {
		//openPopupXxljob(w)
		var popup *widget.PopUp
		titleLable := widget.NewLabel("xxl-job重启")
		IpLable := widget.NewLabel("")
		sevSelect := widget.NewSelect([]string{"192.168.102.40"}, func(s string) {
			IpLable.Text = s
			IpLable.Refresh()
		})
		popupContent := container.NewVBox(
			container.NewAdaptiveGrid(2, titleLable),
			sevSelect,
			IpLable,
			widget.NewButton("重启xxl-job", func() {
				xxljobWindow(app, IpLable.Text)
			}),
			widget.NewButton("Close", func() {
				popup.Hide()
			}),
		)
		popup = widget.NewPopUp(popupContent, w.Canvas())
		popup.ShowAtPosition(fyne.NewPos(
			w.Canvas().Size().Height/2-popup.MinSize().Height/2,
			w.Canvas().Size().Width/2-popup.MinSize().Width/2,
		))
		popup.Show()
	})

	//menu := fyne.NewMainMenu(&menuItem)
	newMenu := fyne.NewMenu("Tools", menuItem1, menuItem2)
	workMenu := fyne.NewMenu("Work", menuItem4, menuItem5)
	newMenu1 := fyne.NewMenu("Help", menuItem3)
	menu := fyne.NewMainMenu(newMenu, workMenu, newMenu1)
	//labelTitle := widget.NewLabel("think to do ?")
	//labelTitle.TextStyle = fyne.TextStyle{Bold: true}
	var content1 *fyne.Container
	labelAddr := widget.NewLabel("")
	//labelHead := widget.NewLabel("请选择服务器:")
	labelAddrHear := widget.NewLabel("当前选择的服务器为:")
	labelLayout := container.New(layout.NewHBoxLayout(), labelAddrHear, labelAddr)
	restoreLable := widget.NewLabel("")
	var tomcatFilePath string
	tomcatFilePath = "/usr/local/tomcat/webapps/"
	//labelUser := widget.NewLabel("")
	// sevSelect := widget.NewSelect([]string{"192.168.102.50", "192.168.102.92", "192.168.102.113"}, func(s string) {
	// 	labelAddr.Text = s
	// 	labelAddr.Refresh()
	// })
	// sevSelect.PlaceHolder = "192.168.102.50"

	//gridSelect := container.New(layout.NewVBoxLayout(), sevSelect, form)
	// gridSelect := container.New(layout.NewHBoxLayout(), labelHead, sevSelect)
	// grid := container.New(layout.NewGridLayout(2), gridSelect)

	//服务和ip对应，读取options.json
	resultData, err := ioutil.ReadFile("options.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}
	var server []Server
	err = json.Unmarshal(resultData, &server)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}
	var SevSlice []string
	var SevIplist []string
	var Iplist []string
	// Iplist = []string{"192.168.113.83", "192.168.113.84"}

	radioGroup1 := widget.NewRadioGroup(Iplist, func(selected string) {
		fmt.Println("Selected option from Group 1:", selected)
	})

	for _, servicesIP := range server {
		SevSlice = append(SevSlice, servicesIP.Services)
		SevIplist = append(SevIplist, servicesIP.Services+","+servicesIP.Address)
		//Iplist = append(Iplist, servicesIP.Address)
	}

	serviceIpSelect := widget.NewSelect(SevSlice, func(s string) {
		w.Canvas().Refresh(content1)
		Iplist = []string{}
		for i := 0; i < len(SevSlice); i++ {
			if SevSlice[i] == s {
				fmt.Println(SevIplist[i])
				substrings := strings.Split(SevIplist[i], ",")
				for j := 0; j < len(substrings); j++ {
					if j != 0 {
						Iplist = append(Iplist, substrings[j])
					}
				}
			}
		}
		radioxxx := widget.NewRadioGroup(Iplist, func(s string) {
			fmt.Println("Selected server is:", s)
			labelAddr.Text = s
			labelAddr.Refresh()
		})
		content1.Objects[1] = radioxxx
		w.Canvas().Refresh(content1)
	})

	ipSelectButton := container.NewGridWithColumns(2, serviceIpSelect)

	fileselect := widget.NewLabel("")
	findRestultLabel := widget.NewLabel("")
	//var content1 *fyne.Container
	options := []string{}
	//var findPath string
	radioGroup := widget.NewRadioGroup(options, func(selected string) {
		fmt.Println("Selected option:", selected)
	})
	//pathSlice := []string{""}
	openButton := widget.NewButton("选择本地文件", func() {
		showOpenFileDialog(w, fileselect)
	})
	transButtn := widget.NewButton("替换", func() {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println("远程连接发生错误:", err)
				// 执行适当的错误处理操作
			}
		}()
		if labelAddr.Text == "" {
			reLable := widget.NewLabel("未选择服务器!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable)
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		}
		client, err := CreateSftp(labelAddr.Text, logininfo.Username, logininfo.Password, logininfo.SshPort)
		if err != nil {
			panic(err)
		}
		defer client.Close()
		var localFilePath = fileselect.Text
		//var remoteDir = "/tmp"
		var remoteDir = findRestultLabel.Text
		srcFile, err := os.Open(localFilePath)
		if err != nil {
			panic(err)
		}
		defer srcFile.Close()
		var remoteFileName = path.Base(localFilePath)

		//上传前文件备份
		conn, err := pwdConnect(labelAddr.Text, logininfo.Username, logininfo.Password, logininfo.SshPort)
		if err != nil {
			return
		}
		defer conn.Close()
		//创建ssh session会话
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()
		//cmd := "find /tmp -name Autorun.inf"
		//selectFile := filepath.Base(fileselect.Text)
		//fmt.Println(selectFile)
		currentTime := time.Now()
		formattedTime := currentTime.Format("20060102150405")
		cmd := "cp " + remoteDir + "/" + remoteFileName + " " + remoteDir + "/" + remoteFileName + "." + formattedTime
		fmt.Println(cmd)
		//cmdInfo, err := session.CombinedOutput(cmd)
		cmdInfo, err := session.Output(cmd)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(cmdInfo))
		//

		dstFile, err := client.Create(path.Join(remoteDir, remoteFileName))
		if err != nil {
			panic(err)
		}
		defer dstFile.Close()
		buf := make([]byte, 1024)
		for {
			n, _ := srcFile.Read(buf)
			if n == 0 {
				break
			}
			dstFile.Write(buf)
		}

		var cleanPopup *widget.PopUp
		if len(buf) == 0 {
			reLable := widget.NewLabel("file size is zero")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable, layout.NewSpacer())
			cleanPopup = widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		} else {
			reLable := widget.NewLabel("Replace successfull!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable, layout.NewSpacer())
			cleanPopup = widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		}
	})

	restoreButton := widget.NewButton("还原", func() {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println("远程连接发生错误:", err)
				// 执行适当的错误处理操作
			}
		}()
		if labelAddr.Text == "" {
			reLable := widget.NewLabel("未选择服务器!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable)
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		}
		lastDotIndex := strings.LastIndex(restoreLable.Text, ".")
		oriFileName := restoreLable.Text[:lastDotIndex]
		conn, err := pwdConnect(labelAddr.Text, logininfo.Username, logininfo.Password, logininfo.SshPort)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		//创建ssh session会话
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()
		cmd := "cp " + restoreLable.Text + " " + oriFileName
		cmdInfo, err := session.Output(cmd)
		if err != nil {
			panic(err)
		} else {
			//var cleanPopup *widget.PopUp
			reLable := widget.NewLabel("Restore sucessful!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable, layout.NewSpacer())
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		}
		fmt.Println(string(cmdInfo))
	})

	restartBtuuon := widget.NewButton("重启服务", func() {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println("远程连接发生错误:", err)
				// 执行适当的错误处理操作
			}
		}()

		if labelAddr.Text == "" {
			reLable := widget.NewLabel("未选择服务器!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable, layout.NewSpacer())
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		} else {
			conn, err := pwdConnect(labelAddr.Text, logininfo.Username, logininfo.Password, logininfo.SshPort)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			//创建ssh session会话
			session, err := conn.NewSession()
			if err != nil {
				panic(err)
			}
			defer session.Close()
			cmd := "ps -ef | grep /usr/local/tomcat | grep -v grep | awk '{print $2}'|xargs kill -9 ; sleep 2; /usr/local/tomcat/bin/startup.sh"
			//cmdInfo, err := session.Output(cmd)
			cmdInfo, err := session.CombinedOutput(cmd)
			if err != nil {
				panic(err)
			} else {
				//var cleanPopup *widget.PopUp
				reLable := widget.NewLabel(string(cmdInfo))
				reLable.Refresh()
				cleanCon := container.NewVBox(reLable, layout.NewSpacer())
				cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
				cleanPopup.ShowAtPosition(fyne.NewPos(
					w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
					w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
				))
				cleanPopup.Show()
			}
			fmt.Println(string(cmdInfo))
		}
	})

	transRestore := container.New(layout.NewGridLayoutWithColumns(3), transButtn, restoreButton, restartBtuuon)
	findButton := widget.NewButton("查找文件", func() {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println("远程连接发生错误:", err)
				// 执行适当的错误处理操作
			}
		}()
		if len(labelAddr.Text) == 0 {
			promptLable := widget.NewLabel("请选择服务器和本地文件")
			promptLable.Refresh()
			cleanCon := container.NewVBox(promptLable)
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		}
		conn, err := pwdConnect(labelAddr.Text, "root", "1234qwer", 22)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		//创建ssh session会话
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()
		//cmd := "find /tmp -name Autorun.inf"
		selectFile := filepath.Base(fileselect.Text)
		//fmt.Println(selectFile)
		cmd := "find " + tomcatFilePath + " -name " + "\"" + selectFile + "*" + "\""
		//cmd := "find /usr/local/tomcat/webapps/ -name " + "\"" + selectFile + "*" + "\""
		fmt.Println(cmd)
		//cmdInfo, err := session.CombinedOutput(cmd)
		cmdInfo, err := session.Output(cmd)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(cmdInfo))
		resultSlice := strings.Split(string(cmdInfo), "\n")
		resultSlice = resultSlice[:len(resultSlice)-1]
		if len(resultSlice) == 0 {
			//var cleanPopup *widget.PopUp
			reLable := widget.NewLabel("file not exsit!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable, layout.NewSpacer())
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		} else {
			newRadioGroup := widget.NewRadioGroup(resultSlice, func(selected string) {
				fmt.Println("Selected option:", selected)
				restoreLable.SetText(selected)
				findRestultLabel.SetText(path.Dir(selected))
			})
			content1.Objects[5] = newRadioGroup
			//content1.Add(newRadioGroup)
			w.Canvas().Refresh(content1)
		}

	})
	checkbox := widget.NewCheck("静态化", func(b bool) {
		tomcatFilePath = "/mnt/b2b/"
		fmt.Println(tomcatFilePath)
	})
	fileSelectHead := widget.NewLabel("选择的本地文件:")
	transHbox := container.New(layout.NewHBoxLayout(), fileSelectHead, fileselect, layout.NewSpacer(), checkbox, findButton)
	//stopChan := make(chan int)
	logButton := widget.NewButton("查看tomcat启动日志", func() {
		currentDate := time.Now().Format("2006-01-02")
		filePath := "/usr/local/tomcat/logs/error-debug." + currentDate + ".log"
		stopChan := make(chan int)
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Println("远程连接发生错误:", err)
			}
		}()
		//判断是否选择服务器
		if labelAddr.Text == "" {
			reLable := widget.NewLabel("未选择服务器!")
			reLable.Refresh()
			cleanCon := container.NewVBox(reLable)
			cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
			cleanPopup.ShowAtPosition(fyne.NewPos(
				w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
				w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
			))
			cleanPopup.Show()
		} else {
			//判断文件是否存在
			fileExists, _ := checkFileExists(labelAddr.Text, logininfo.Username, logininfo.Password, filePath)
			// if err != nil {
			// 	panic(err)
			// }
			if fileExists {
				fmt.Println("File exists")
			} else {
				//fmt.Println("File does not exist")
				reLable := widget.NewLabel("文件不存在!")
				reLable.Refresh()
				cleanCon := container.NewVBox(reLable, layout.NewSpacer())
				cleanPopup := widget.NewPopUp(cleanCon, w.Canvas())
				cleanPopup.ShowAtPosition(fyne.NewPos(
					w.Canvas().Size().Height/2-cleanPopup.MinSize().Height/2,
					w.Canvas().Size().Width/2-cleanPopup.MinSize().Width/2,
				))
				cleanPopup.Show()
			}
		}
		//连接服务器执行tail命令
		conn, err := pwdConnect(labelAddr.Text, "root", "1234qwer", 22)
		if err != nil {
			panic(err)
		}
		defer conn.Close()
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()
		//currentDate := time.Now().Format("2006-01-02")
		cmd := "tail -10 /usr/local/tomcat/logs/error-debug." + currentDate + ".log"
		cmdInfo, err := session.Output(cmd)
		if err != nil {
			panic(err)
		}
		//fmt.Println(string(cmdInfo))

		childWindow := app.NewWindow("执行结果")
		logLable := widget.NewLabel("")
		logLable.SetText(string(cmdInfo))
		//childWindow.Resize(fyne.NewSize(300, 200))
		freshButton := widget.NewButton("手动刷新", func() {
			conn, err := pwdConnect(labelAddr.Text, logininfo.Username, logininfo.Password, logininfo.SshPort)
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			session, err := conn.NewSession()
			if err != nil {
				panic(err)
			}
			defer session.Close()
			currentDate := time.Now().Format("2006-01-02")
			cmd := "tail -10 /usr/local/tomcat/logs/error-debug." + currentDate + ".log"
			cmdInfo, err := session.Output(cmd)
			if err != nil {
				panic(err)
			}
			logLable.SetText(string(cmdInfo))
			logLable.Refresh()
		})
		autoFreshButton := widget.NewButton("自动刷新", func() {
			go autoFlushLog(labelAddr.Text, logLable, stopChan)
		})
		stopFreshButton := widget.NewButton("停止刷新", func() {
			stopChan <- 1
			fmt.Println("auto flush stop")
		})
		autoStopContainter := container.New(layout.NewGridLayoutWithColumns(2), autoFreshButton, stopFreshButton)
		logFlushButton := container.New(layout.NewGridLayoutWithColumns(2), freshButton, autoStopContainter)
		childWindowContent := container.NewVBox(
			logFlushButton,
			logLable,
		)
		childWindow.SetContent(childWindowContent)
		childWindow.Show()
		childWindow.SetOnClosed(func() {
			fmt.Println("close the child windows")
			stopChan <- 1
		})
	})
	logButtonContainer := container.New(layout.NewGridLayoutWithColumns(2), logButton, layout.NewSpacer())
	// btn_exit := widget.NewButton("EXIT", func() {
	// 	app.Quit()
	// })

	w.SetMainMenu(menu)
	content1 = container.NewVBox(
		ipSelectButton,
		radioGroup1,
		labelLayout,
		openButton,
		transHbox,
		radioGroup,
		transRestore,
		findRestultLabel,
		logButtonContainer,
		//ipSelectButton,
		//radioGroup1,
	)

	w.SetContent(
		content1,
	)

	w.ShowAndRun()
}

// func myIP() map[string]string {
// 	m := make(map[string]string)
// 	req, err := http.Get("http://ip-api.com/json/")
// 	if err != nil {
// 		//return err.Error()
// 		fmt.Sprintln(err)
// 	}
// 	defer req.Body.Close()
// 	body, err := ioutil.ReadAll(req.Body)
// 	if err != nil {
// 		//return err.Error()
// 		fmt.Sprintln("error")
// 	}
// 	//var ip IP
// 	json.Unmarshal(body, &m)
// 	return m
// }

// windows command 中文支持
type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

type Server struct {
	Services string `json:"label"`
	Address  string `json:"value"`
}

type Logininfo struct {
	Username string
	Password string
	SshPort  int
}

func toggleBool(b *bool) {
	*b = !*b
}

func ConvertByte2String(byte []byte, charset Charset) string {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes, _ = simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}
	return str
}
func openChildWindow(app fyne.App, sevIP string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			log.Println("远程连接错误:", err)
			// 执行适当的错误处理操作
		}
	}()
	childWindow := app.NewWindow("执行结果")
	childWindow.Resize(fyne.NewSize(200, 200))

	conn, err := pwdConnect(sevIP, "root", "1234qwer", 22)
	if err != nil {
		return
	}
	defer conn.Close()
	//创建ssh session会话
	session, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	cmd := "find /usr/local/tomcat/logs/* -type f -ctime +30 |xargs -i rm -rf {}"
	cmdInfo, err := session.CombinedOutput(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cmdInfo))
	if cmdInfo != nil {
		childWindowContent := container.NewVBox(
			widget.NewLabel(string(cmdInfo)),
		)
		childWindow.SetContent(childWindowContent)
		childWindow.Show()
	} else {
		childWindowContent := container.NewVBox(
			//widget.NewLabel("This is a child window"),
			widget.NewLabel("exec sucessful"),
		)
		childWindow.SetContent(childWindowContent)
		childWindow.Show()
	}

}

func pwdConnect(sshHost, sshUser, sshPassword string, sshPort int) (*ssh.Client, error) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			log.Println("远程连接错误:", err)
			// 执行适当的错误处理操作
		}
	}()
	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	Results, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		//log.Fatal("connect is failed!", err)
		panic(err)
	}
	return Results, nil
}

// sftp
func CreateSftp(sshHost, sshUser, sshPassword string, sshPort int) (*sftp.Client, error) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			log.Println("远程连接错误:", err)
			// 执行适当的错误处理操作
		}
	}()
	conn, err := pwdConnect(sshHost, sshUser, sshPassword, sshPort)
	if err != nil {
		fmt.Println("sftp connect is failed!")
		panic(err)
	}
	Results, err := sftp.NewClient(conn)
	if err != nil {
		return nil, err
	}
	return Results, err
}

// 判断文件是否存在
func checkFileExists(sshHost, sshUser, sshPassword string, filePath string) (bool, error) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			log.Println("远程连接错误:", err)
			// 执行适当的错误处理操作
		}
	}()
	config := &ssh.ClientConfig{
		Timeout:         5 * time.Second,
		User:            sshUser,
		Auth:            []ssh.AuthMethod{ssh.Password(sshPassword)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	addr := fmt.Sprintf("%s:%d", sshHost, 22)
	Client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		//log.Fatal("connect is failed!", err)
		panic(err)
	}
	defer Client.Close()

	// 执行命令检查文件是否存在
	session, err := Client.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()
	output, err := session.Output("test -e " + filePath)
	if err != nil {
		return false, err
	}
	return !strings.Contains(string(output), "No such file"), nil
}

// 重启xxljob
func xxljobWindow(app fyne.App, sevIP string) {
	defer func() {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
			log.Println("远程连接错误:", err)
			// 执行适当的错误处理操作
		}
	}()
	childWindow := app.NewWindow("执行结果")
	childWindow.Resize(fyne.NewSize(400, 300))

	conn, err := pwdConnect(sevIP, "root", "1Qaz2Wsx", 22)
	if err != nil {
		return
	}
	defer conn.Close()
	//创建ssh session会话
	session, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	cmd := "ps -ef | grep xxl-job-admin-2.0.2.jar | grep -v grep | awk '{print $2}'|xargs kill -9 ; sleep 2; nohup /usr/java/jdk1.8.0_144/bin/java -jar /usr/local/xxjob/xxl-job-admin-2.0.2.jar>/dev/null &"
	cmdInfo, err := session.CombinedOutput(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cmdInfo))
	if cmdInfo != nil {
		childWindowContent := container.NewVBox(
			//widget.NewLabel("This is a child window"),
			widget.NewLabel(string(cmdInfo)),
		)
		childWindow.SetContent(childWindowContent)
		childWindow.Show()
	} else {
		childWindowContent := container.NewVBox(
			//widget.NewLabel("This is a child window"),
			widget.NewLabel("exec sucessful"),
		)
		childWindow.SetContent(childWindowContent)
		childWindow.Show()
	}

}

// 查看文件内容
func contentOpenFileDialog(win fyne.Window) {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader != nil {
			defer reader.Close()

			// 读取文件内容
			content, err := ioutil.ReadAll(reader)
			if err != nil {
				fmt.Println("Failed to read file:", err)
				return
			}

			// 将文件内容打印到控制台
			fmt.Println(string(content))
		}
	}, win)
}

// 获取文件路径
func showOpenFileDialog(win fyne.Window, label *widget.Label) {
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err == nil && reader != nil {
			// 获取文件路径
			filePath := reader.URI().Path()
			fmt.Println("Selected file:", filePath)
			// 关闭文件读取器
			label.SetText(filePath)
			label.Refresh()
			_ = reader.Close()
		}
	}, win)

	fileDialog.Show()
}

func autoFlushLog(labelAddr string, logLable *widget.Label, stopChan chan int) {
	conn, err := pwdConnect(labelAddr, "root", "1234qwer", 22)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for {
		session, err := conn.NewSession()
		if err != nil {
			panic(err)
		}
		defer session.Close()
		currentDate := time.Now().Format("2006-01-02")
		cmd := "tail -10 /usr/local/tomcat/logs/error-debug." + currentDate + ".log"
		cmdInfo, err := session.Output(cmd)
		if err != nil {
			panic(err)
		}
		logLable.SetText(string(cmdInfo))
		logLable.Refresh()
		//time.Sleep(5 * time.Second)
		fmt.Println("log flush test")
		select {
		case <-stopChan:
			return
		case <-time.After(5 * time.Second):
			// 等待五秒钟后继续执行下一次
		}
	}
}
