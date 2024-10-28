OUTPUT_DIR = build

ifeq ($(OS),Windows_NT)
    RM = if exist $(OUTPUT_DIR) (rmdir /s /q $(OUTPUT_DIR))
else
    RM = rm -rf $(OUTPUT_DIR)
endif

build:
	@go build -o $(OUTPUT_DIR)/

clean:
	@$(RM)

.PHONY: build clean
