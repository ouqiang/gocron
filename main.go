package main

import (
    "context"
    "fmt"
	"path/filepath"
    "os"
    "os/exec"
    "os/signal"
    "time"
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
    if err := run(); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }
}