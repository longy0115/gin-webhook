package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/gitpush", gitpush)
	r.GET("/hello", hello)
	r.GET("/sl", sl)
	r.Run(":8091")
}

func hello(c *gin.Context) {
	c.JSON(200, gin.H{
		"mess": "helloWordLong",
	})
}

//shell 命令
func sl(c *gin.Context) {
	out, errout, err := outShell("ls -ltr")
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)
}

func outShell(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

// Github Webhooks Post请求处理函数
func gitpush(c *gin.Context) {
	// 验证签名
	if matched, _ := verifySignature(c); !matched {
		err := "Signatures didn't match!"
		c.String(http.StatusForbidden, err)
		fmt.Println(err)
		return
	}

	fmt.Println("Signatures is match! go!")

	// 你自己的业务逻辑......
	cmd := "cd /home/dnmp/www/doris-ucat && git pull"
	out, errout, err := outShell(cmd)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	fmt.Println("--- stdout ---")
	fmt.Println(out)
	fmt.Println("--- stderr ---")
	fmt.Println(errout)

	c.String(http.StatusOK, "OK")
}

// 验证签名
func verifySignature(c *gin.Context) (bool, error) {
	payloadBody, err := c.GetRawData()
	if err != nil {
		return false, err
	}

	// 获取请求头中的签名信息
	hSignature := c.GetHeader("X-Hub-Signature")

	// 计算Payload签名
	signature := hmacSha1(payloadBody)
	fmt.Println(signature)
	return (hSignature == signature), nil

}

// hmac-sha1
func hmacSha1(payloadBody []byte) string {
	h := hmac.New(sha1.New, []byte("webhoods-ucat"))
	h.Write(payloadBody)
	return "sha1=" + hex.EncodeToString(h.Sum(nil))
}
