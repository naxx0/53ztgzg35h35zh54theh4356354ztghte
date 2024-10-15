
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net"
    "os"
    "sync"
    "time"

    "gopkg.in/gomail.v2"
)

type Config struct {
    SystemName string    `json:"system_name"`
    SMTP       SMTPConfig `json:"smtp"`
    Ports      []string   `json:"ports"`
}

type SMTPConfig struct {
    Server    string `json:"server"`
    Port      int    `json:"port"`
    User      string `json:"user"`
    Pass      string `json:"pass"`
    Recipient string `json:"recipient"`
}

var config Config
var wg sync.WaitGroup

func main() {
    // Load configuration from config.json
    err := loadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    ipAddress, err := getLocalIP()
    if err != nil {
        log.Fatalf("Failed to get local IP address: %v", err)
    }

    // Send startup email with IP address and ports
    startupMessage := fmt.Sprintf("Protos Listener started on %s.\nIP Address: %s\nListening on ports: %v", config.SystemName, ipAddress, config.Ports)
    if err := sendEmail(fmt.Sprintf("Listener on %s started successfully", config.SystemName), startupMessage); err != nil {
        log.Fatalf("Failed to send startup email: %v", err)
    }

    // Launch listeners for each port
    for _, port := range config.Ports {
        wg.Add(1)
        go func(port string) {
            defer wg.Done()
            listener(port)
        }(port)
    }

    wg.Wait()
}

func listener(port string) {
    ln, err := net.Listen("tcp", ":"+port)
    if err != nil {
        log.Printf("Failed to listen on port %s: %v", port, err)
        return
    }
    defer ln.Close()
    log.Printf("Listening on port %s", port)

    for {
        conn, err := ln.Accept()
        if err != nil {
            log.Printf("Failed to accept connection on port %s: %v", port, err)
            continue
        }

        go handleConnection(conn, port)
    }
}

func handleConnection(conn net.Conn, port string) {
    defer conn.Close()

    remoteAddr := conn.RemoteAddr().String()
    log.Printf("Connection established from %s on port %s", remoteAddr, port)

    // Read the first 1024 bytes from the connection
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil && err != io.EOF {
        log.Printf("Failed to read from connection: %v", err)
    }
    dataSnippet := string(buffer[:n])

    // Send email when a connection is detected
    accessMessage := fmt.Sprintf("Connection detected from %s on port %s\n\nFirst 1024 bytes of data:\n%s", remoteAddr, port, dataSnippet)
    if err := sendEmail(fmt.Sprintf("Access detected on %s on %s", port, config.SystemName), accessMessage); err != nil {
        log.Printf("Failed to send access email: %v", err)
    }

    // Terminate the program after access is detected
    log.Printf("Terminating program after access detected on port %s", port)
    go func() {
        time.Sleep(2 * time.Second)  // Wait 2 seconds to ensure logs and email are sent before terminating
        os.Exit(0)  // Terminate the program
    }()
}

func sendEmail(subject, body string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", config.SMTP.User)
    m.SetHeader("To", config.SMTP.Recipient)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    d := gomail.NewDialer(config.SMTP.Server, config.SMTP.Port, config.SMTP.User, config.SMTP.Pass)

    return d.DialAndSend(m)
}

func getLocalIP() (string, error) {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        return "", err
    }

    for _, addr := range addrs {
        if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
            return ipNet.IP.String(), nil
        }
    }
    return "", fmt.Errorf("no external IP address found")
}

func loadConfig(filename string) error {
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return err
    }

    return json.Unmarshal(data, &config)
}
