package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"messanger/api"
	"messanger/config"
	"messanger/dbcall"
)

const PID_FILE = "server.pid"

var server *http.Server

func main() {
	// Kill existing process
	killExistingProcess()

	// Config Init
	initConfig()

	// Database Init
	initDB()

	// DB 연결을 dbcall에 전달
	dbcall.InitDB(db)

	// WebSocket Init
	initWebSocket()

	// Server Init
	initServer()

	// Save PID
	savePID()

	// Web Server start
	startWebServer()

	// ...
}

// Config Init
func initConfig() {
	if err := config.Init(); err != nil {
		log.Fatalf("### Failed to load config: %v", err)
	}
	log.Println("### Config Init Complete")
}

// Server Init
func initServer() {
	port := config.GetPort()
	mux := http.NewServeMux()

	// 정적 파일 서빙
	mux.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("public/html"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("public/js"))))

	// API 라우트 설정
	api.SetupRoutes(mux)

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/ws", handleWebSocket)

	server = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Println("### Server Init Complete")
}

// Web Server Start
func startWebServer() {
	port := config.GetPort()
	go func() {
		log.Printf("### Server starting on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("### Server failed to start: %v", err)
		}
	}()

	log.Println("### Web Server Start Complete")

	waitForShutdown()
}

// rootHandler Root Path Handler
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "public/html/index.html")
}

// waitForShutdown Wait for Shutdown
func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("### Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("### Server forced to shutdown: %v", err)
	}

	removePIDFile()
	log.Println("### Server exited")
}

// killExistingProcess 기존 프로세스 확인 및 종료
func killExistingProcess() {
	if _, err := os.Stat(PID_FILE); os.IsNotExist(err) {
		return
	}

	data, err := os.ReadFile(PID_FILE)
	if err != nil {
		os.Remove(PID_FILE)
		return
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		os.Remove(PID_FILE)
		return
	}

	// 프로세스가 실행 중인지 확인
	process, err := os.FindProcess(pid)
	if err != nil {
		os.Remove(PID_FILE)
		return
	}

	// Windows와 Unix 계열 다른 처리
	if runtime.GOOS == "windows" {
		// Windows: taskkill 사용
		cmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(pid))
		cmd.Run()
	} else {
		// Unix: SIGTERM 전송
		process.Signal(syscall.SIGTERM)
	}

	log.Printf("### Killed existing process (PID: %d)", pid)
	os.Remove(PID_FILE)
}

// savePID 현재 PID를 파일에 저장
func savePID() {
	pid := os.Getpid()
	err := os.WriteFile(PID_FILE, []byte(strconv.Itoa(pid)), 0644)
	if err != nil {
		log.Fatalf("### Failed to save PID file: %v", err)
	}
	log.Printf("### Server PID: %d (saved to %s)", pid, PID_FILE)
}

// removePIDFile PID 파일 삭제
func removePIDFile() {
	if _, err := os.Stat(PID_FILE); err == nil {
		os.Remove(PID_FILE)
		log.Println("### PID file removed")
	}
}
