ifneq ("$(shell which x86_64-w64-mingw32-gcc)","")
compiler = x86_64-w64-mingw32-gcc
else ifneq ("$(shell which amd64-mingw32msvc-gcc)","")
compiler = amd64-mingw32msvc-gcc
else
ignored = $(error No compatible compiler found on system path)
endif

arch = amd64
mimikatz_version = 2.2.0-20200918-fix
key := $(shell cat /dev/urandom | tr -dc 'a-zA-Z0-9' | fold -w 128 | head -n 1)
inp := $(shell pwd)/MemoryModule/build/MemoryModule.a

all: get encrypt pack mimikatz.exe

mimikatz.exe: MemoryModule mimikatz.go
	CGO_LDFLAGS=$(inp) CGO_LDFLAGS_ALLOW=".*\.a" CC=$(compiler) CGO_ENABLED=1 GOOS=windows GOARCH=$(arch) go build -x -ldflags "-X main.key=$(key)" -o bin/gogokatz.exe ./cmd/gogokatz/

# Dependency build. 
SUBDIRS = MemoryModule
subdirs: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@
# Override default subdir build behavior (make) with cmake. 
MemoryModule:
	[ "`ls -A MemoryModule`" ] || git submodule update --init
	cmake -HMemoryModule -BMemoryModule/build
	cmake --build MemoryModule/build --target MemoryModule

# Create pkged.go with mimikatz.exe.enc
pack:
	pkger -o ./cmd/gogokatz/

encrypt:
	mkdir bin resources || true
	go build --tags pack -ldflags "-X main.key=$(key)" -o bin/packer ./cmd/gogokatz/
	bin/packer -m mimikatz.exe -o resources/mimikatz.exe.enc

get:
	curl -OL https://github.com/gentilkiwi/mimikatz/releases/download/$(mimikatz_version)/mimikatz_trunk.7z
	7z e mimikatz_trunk.7z x64/mimikatz.exe

# Clean target. 
CLEANDIRS = $(SUBDIRS:%=clean-%)
clean: $(CLEANDIRS)
	rm -f mimikatz.exe mimikatz_trunk.7z ./cmd/gogokatz/pkged.go
	rm -rf bin resources/*
$(CLEANDIRS): 
	$(MAKE) -C $(@:clean-%=%) clean
clean-MemoryModule:
	$(MAKE) -C $(@:clean-%=%) clean
	rm -rf MemoryModule/build

test:
	$(MAKE) -C tests test

.PHONY: subdirs $(INSTALLDIRS) $(SUBDIRS) $(CLEANDIRS) clean test pack get check_deps
