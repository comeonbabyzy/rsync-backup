package commands

//func shellCmd(command string, arg ...string) (string, error) {
/*
func RunCmd(env string, command string, arg ...string) (string, error) {
	cmd := exec.Command(command, arg...)

	cmd.Env = append(os.Environ(),
		//"RSYNC_PASSWORD="+rsyncPassword,
		env,
	)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return "", errors.Wrap(err, stderr.String())
	}

	return "Result: " + out.String(), nil
}
*/

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func RunCmd(env string, cmd string, arg ...string) error {
	ctx, cancel := context.WithCancel(context.Background())

	go func(cancelFunc context.CancelFunc) {
		time.Sleep(3 * time.Second)
		//cancelFunc()
	}(cancel)

	err := command(ctx, env, cmd, arg...)

	return err
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

func command(ctx context.Context, env string, cmd string, arg ...string) error {
	//var c *exec.Cmd

	c := exec.CommandContext(ctx, cmd, arg...)
	c.Env = append(c.Env, env)

	log.WithFields(log.Fields{
		"env": c.Env,
		"cmd": cmd,
		"arg": arg,
	}).Debug("command")

	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := c.StderrPipe()
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	// 因为有2个任务, 一个需要读取stderr 另一个需要读取stdout
	wg.Add(2)
	go read(ctx, &wg, stderr)
	go read(ctx, &wg, stdout)
	// 这里一定要用start,而不是run 详情请看下面的图
	err = c.Start()
	// 等待任务结束
	wg.Wait()
	return err
}

func read(ctx context.Context, wg *sync.WaitGroup, std io.ReadCloser) {
	reader := bufio.NewReader(std)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}

			decodedString := ConvertByte2String([]byte(readString), GB18030)

			fmt.Print(decodedString)
		}
	}
}
