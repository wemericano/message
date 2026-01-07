BINARY_NAME=messanger
BUILD_DIR=./bin

# 기본 타겟: 빌드
all:
	@echo Make...
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) .
	@echo Build complete!

# 빌드 후 실행
run: all
	@echo Running...
	@$(BUILD_DIR)/$(BINARY_NAME)

# 빌드 파일 정리
clean:
	@echo Cleaning...
	@powershell -Command "if (Test-Path bin) { Remove-Item -Path bin -Recurse -Force }"
	@go clean
	@echo Clean complete!

# 의존성 설치
install:
	@echo Installing dependencies...
	@go mod download
	@go mod tidy
	@echo Install complete!

.PHONY: all run clean install
