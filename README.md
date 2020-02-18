# PJGO

PJSIP implementation on Go using SWIG interface

## Getting started

### Prerequisites

`go`

`swig`

`pjsua` installed in system

### Installing PJSUA in project

Redefine if needed in Makefile:

* `SWIG_DIR` - subdirectory `pjsip-apps/src/swig` in your PJSIP source directory

* `CGO_CXXFLAGS` - directory with include installed pjsua header files - compiler flag

`make pjsua-install`

Then insert into `pjsua.go` compiler flags from `pkgconfig/libpjproject.pc` from installed pjsua lib directory like this:

```go
package pjsua2

/*
#cgo CXXFLAGS: -I/usr/local/include -g -O2 -I/usr/local/include -DPJ_AUTOCONF=1 -O2 -DPJ_IS_BIG_ENDIAN=0 -DPJ_IS_LITTLE_ENDIAN=1
#cgo LDFLAGS: -L/usr/local/lib -lpjsua2-x86_64-unknown-linux-gnu -lstdc++ -lpjsua-x86_64-unknown-linux-gnu -lpjsip-ua-x86_64-unknown-linux-gnu -lpjsip-simple-x86_64-unknown-linux-gnu -lpjsip-x86_64-unknown-linux-gnu -lpjmedia-codec-x86_64-unknown-linux-gnu -lpjmedia-x86_64-unknown-linux-gnu -lpjmedia-videodev-x86_64-unknown-linux-gnu -lpjmedia-audiodev-x86_64-unknown-linux-gnu -lpjmedia-x86_64-unknown-linux-gnu -lpjnath-x86_64-unknown-linux-gnu -lpjlib-util-x86_64-unknown-linux-gnu -lsrtp-x86_64-unknown-linux-gnu -lresample-x86_64-unknown-linux-gnu -lgsmcodec-x86_64-unknown-linux-gnu -lspeex-x86_64-unknown-linux-gnu -lilbccodec-x86_64-unknown-linux-gnu -lg7221codec-x86_64-unknown-linux-gnu -lyuv-x86_64-unknown-linux-gnu -lwebrtc-x86_64-unknown-linux-gnu  -lpj-x86_64-unknown-linux-gnu -lopus -lssl -lcrypto -luuid -lm -lrt -lpthread  -lasound -L/usr/lib -pthread -lSDL2  -lavdevice -lavformat -lavcodec -lswscale -lavutil  -lv4l2 -lopencore-amrnb -lopencore-amrwb


#define intgo swig_intgo
typedef void *swig_voidp;

```

with variables substitution, where:

* `CXXFLAGS` - Cflags
* `LDFLAGS` - Libs