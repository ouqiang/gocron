package main

import (
    "context"
    "fmt"
	"path/filepath"
    "os"
    "os/exec"
    "os/signal"
    "time"
    "github.com/robfig/cron/v3"
	"log"
)

func run() error {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()
    cmd := exec.Command("cmd", "/C", "python C:\\github\\kabu-backend\\job\\jpx-index.py")
	fmt.Println(filepath.Dir("C:\\github\\kabu-backend\\job\\jpx-index.py"))
	cmd.Dir = filepath.Dir("C:\\github\\kabu-backend\\job\\jpx-index.py")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Start(); err != nil {
        return err
    }
    errCh := make(chan error, 1)
    go func() {
        defer close(errCh)
        errCh <- cmd.Wait()
    }()
    for {
        select {
        case exitErr := <-errCh:
            return exitErr
        case <-ctx.Done():
            fmt.Println("Send SIGINT")
            cmd.Process.Signal(os.Interrupt)
            select {
            case exitErr := <-errCh:
                return exitErr
            case <-time.After(5 * time.Second):
                fmt.Println("Send SIGKILL")
                cmd.Process.Kill()
                return <-errCh
            }
        }
    }
}

func main() {
    p := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	s, err := p.Parse("35 16 ? * 1-5")
	if err != nil {
		log.Fatal(err)
	}
	t := time.Now()
	for i := 0; i < 5; i++ {
		t = s.Next(t)
		fmt.Println(t)
	}
}