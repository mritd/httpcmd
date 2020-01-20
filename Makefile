all:
	gox -osarch="darwin/amd64 linux/386 linux/amd64" \
        -output="dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
clean:
	rm -rf dist

install:
	go install

.PHONY : all clean install

.EXPORT_ALL_VARIABLES:

GO111MODULE = on
GOPROXY = https://goproxy.io
GOSUMDB = sum.golang.google.cn