package git

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"testing"

	"github.com/opensourceways/xihe-server/config"
)

func TestGitlab(t *testing.T) {
	cfg, err := config.LoadConfig("../../../conf/app.conf.yaml")
	if err != nil {
		t.Fatalf("LoadConfig error :%v", err)
	}
	userGitlabClient, err := NewGitlabClient(cfg)
	if err != nil {
		t.Fatal("NewUserGitlabClient error", err)
	}

	// userClient := NewGitUserClient(userGitlabClient)
	// _, err = userClient.CreateUser("ceshiyonghu5", "ceshiyonghu5", "23523243233@qq.com", "ceshiyonghu", "a good boy", true)
	// if err != nil {
	// 	t.Fatal("CreateUser error", err)
	// }

	projectClient := NewGitProjectClient(userGitlabClient)
	// err = projectClient.CreateProject("name", "desc", "private", true, true)
	// if err != nil {
	// 	t.Fatal("CreateProject error", err)
	// }
	desc := "改一个新的描述"
	_, err = projectClient.UpdateProject("root/demo_project", nil, &desc, nil, nil, nil)
	if err != nil {
		t.Fatal("DeleteProject error", err)
	}
	// err = projectClient.DeleteProject("gitlab-instance-6bb21082/testproject")
	// if err != nil {
	// 	t.Fatal("DeleteProject error", err)
	// }
}

func TestWriteData(t *testing.T) {

	b := new(bytes.Buffer)
	w := multipart.NewWriter(b)
	fw, err := os.Create("filename.bin")
	// fw, err := w.CreateFormFile("file", "C:/workspace/src/xihe-server/infrastructure/git/filename.bin")
	if err != nil {
		t.Fatal(err)
	}
	// for i := 0; i < 10; i++ {

	// 	fw.Write([]byte("abc"))
	// }
	content := new(bytes.Buffer)
	content.WriteString("abcdefg123456789ABCDEFG1234")
	temp := make([]byte, 5)
	place := 0
	for i := 0; place < content.Len(); i++ {
		fmt.Printf("第%d轮,值是%d \n", i, place)
		io.ReadFull(content, temp)
		place, err = fw.Write(temp)
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = io.Copy(fw, content)
	if err != nil {
		t.Fatal(err)
	}
	if err = w.Close(); err != nil {
		t.Fatal(err)
		return
	}
	t.Log("ok")
}
