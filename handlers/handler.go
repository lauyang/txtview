package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
	"txtview/models"
)

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		showError(w, "异常", "非法请求，服务器无法响应")
	} else {
		if r.URL.Path == "/" {
			txtViews, err := models.QueryAll()
			if err != nil {
				fmt.Println(err.Error())
				showError(w, "异常", "查询异常")
				return
			}
			t, err := template.ParseFiles("views/index.html")
			if err != nil {
				showError(w, "异常", "页面渲染异常")
				return
			}
			data := make(map[string][]models.TxtView)
			data["TxtViewList"] = txtViews
			t.Execute(w, data)
		} else {
			// 404页面，路由不到的都会到这里
			showError(w, "404", "页面不存在")
		}
	}
}

func MonitorList(w http.ResponseWriter, r *http.Request) {
	for {
		w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		txtViews, _ := models.QueryAll()
		data := make(map[string][]models.TxtView)
		data["TxtViewList"] = txtViews
		marshal, _ := json.Marshal(data)
		dtstr := "data: " + string(marshal) + "\n\n"
		w.Write([]byte(dtstr))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		} else {
			fmt.Println("no flush")
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func MonitorData(w http.ResponseWriter, r *http.Request) {
	for {
		w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache")
		id := r.FormValue("id")
		intId, _ := strconv.ParseInt(id, 10, 64)
		_, _, unLockTime, err := models.GetTxtView(intId)
		if err != nil {
			showError(w, "异常", "查询txtview内容失败")
			return
		}
		var dtstr string
		if time.Now().Unix() <= unLockTime {
			dtstr = "data: false\n\n"
		} else {
			dtstr = "data: true\n\n"
		}

		w.Write([]byte(dtstr))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		} else {
			fmt.Println("no flush")
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func NewTxtView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 显示new页面
		t, _ := template.ParseFiles("views/new.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		title := r.FormValue("title")
		content := r.FormValue("content")
		fmt.Println("title:", title, "---content:", content)
		isExist, err := models.IsExistTxtViewByTitle(title)
		if err != nil {
			showError(w, "异常", "查询数据异常")
			return
		}

		if isExist {
			t, _ := template.ParseFiles("views/new.html")
			data := make(map[string]string)
			data["IsExist"] = strconv.FormatBool(isExist)
			data["Title"] = title
			data["Content"] = content
			t.Execute(w, data)
			return
		}

		id, err := models.InsertTxtView(title, content)
		if err != nil || id <= 0 {
			showError(w, "异常", "插入数据异常")
			return
		}

		f, err := os.Create("views/file/" + title)
		defer f.Close()
		if err != nil {
			showError(w, "异常", "写入文件异常")
		} else {
			_, err = f.Write([]byte(content))
			if err != nil {
				showError(w, "异常", "写入文件异常")
			}
		}

		// 重定向到主界面
		http.Redirect(w, r, "/", http.StatusSeeOther)
		// 没有return，没有效果，重定向不过去
		return
	}
}

func DelTxtView(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		showError(w, "异常", "非法请求")
	} else {
		// 获取表单参数，也可以这么写
		// r.ParseForm()
		// id := r.Form["id"]
		id := r.FormValue("id")
		del := r.FormValue("del")
		// FormValue取到的数据都为string类型，将id转为int64类型
		// strconv.ParseInt(id, 10, 64) 10意思为10进制，64意思为64位
		intId, _ := strconv.ParseInt(id, 10, 64)
		boolDel, _ := strconv.ParseBool(del)

		_, err := models.DelTxtView(intId, !boolDel)
		if err != nil {
			showError(w, "异常", "完成txtview失败")
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func EditTxtView(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 显示edit页面
		// 本可以将title内容提交至此，但url将会异常难看，还是根据id查询吧
		id := r.FormValue("id")
		intId, _ := strconv.ParseInt(id, 10, 64)
		title, content, unLockTime, err := models.GetTxtView(intId)
		if err != nil {
			showError(w, "异常", "查询txtview内容失败")
			return
		}
		if time.Now().Unix() <= unLockTime {
			showError(w, "异常", "该文本正在被编辑，已被锁定")
			return
		}

		res, err := models.UpdateTxtViewUnLockTime(intId, time.Now().Add(time.Minute).Unix())
		if err != nil || res <= 0 {
			showError(w, "异常", "锁定文本异常")
			return
		}

		t, _ := template.ParseFiles("views/edit.html")
		data := make(map[string]string)
		data["Id"] = id
		data["Title"] = title
		data["Content"] = content
		t.Execute(w, data)

	} else if r.Method == "POST" {
		// edit后的数据post提交至此处
		id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		title := r.FormValue("title")
		content := r.FormValue("content")
		res, err := models.EditTxtview(id, time.Now().Unix(), title, content)
		if err != nil || res <= 0 {
			showError(w, "异常", "修改失败")
			return
		}

		// Write back to file
		err = ioutil.WriteFile("views/file/"+title, []byte(content), 0644)
		if err != nil {
			showError(w, "异常", "查询txtview内容失败")
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return

	}
}

func Download(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()                  //解析url传递的参数，对于POST则解析响应包的主体（request body）
	fileName := r.Form["filename"] //filename  文件名
	path := "views/file/"          //文件存放目录
	file, err := os.Open(path + fileName[0])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	fileNames := url.QueryEscape(fileName[0]) // 防止中文乱码
	w.Header().Add("Content-Type", "application/octet-stream")
	w.Header().Add("Content-Disposition", "attachment; filename=\""+fileNames+"\"")

	if err != nil {
		fmt.Println("Read File Err:", err.Error())
	} else {
		w.Write(content)
	}
}

// 错误处理
func showError(w http.ResponseWriter, title string, message string) {
	t, _ := template.ParseFiles("views/error.html")
	data := make(map[string]string)
	data["title"] = title
	data["message"] = message
	t.Execute(w, data)
}
