all: macos linux windows 

windows:
	GOOARCH=amd64 GOOS=windows go build -ldflags="-s -w"; \
	tar -czvf windows-amd64-rlatool.tar.gz rlatool.exe README.md
	rm rlatool.exe

macos: macos-intel macos-arm

macos-intel:
	GOOARCH=amd64 GOOS=darwin go build -ldflags="-s -w"; \
	chmod +x rlatool; \
	tar -czvf macos-amd64-rlatool.gz rlatool README.md
	rm rlatool

macos-arm:
	GOOARCH=arm64 GOOS=darwin go build -ldflags="-s -w"; \
	chmod +x rlatool; \
	tar -czvf macos-arm64-rlatool.gz rlatool README.md
	rm rlatool

linux:
	GOOARCH=amd64 GOOS=linux go build -ldflags="-s -w"; \
	chmod +x rlatool; \
	tar -czvf linux-amd64-rlatool.gz rlatool README.md
	rm rlatool
