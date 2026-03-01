package main

/*
#cgo CFLAGS: -I${SRCDIR}/install/include/freerdp3 -I${SRCDIR}/install/include/winpr3
#cgo LDFLAGS: -L${SRCDIR}/install/lib -lfreerdp3 -lfreerdp-client3 -lwinpr3
#include <freerdp/freerdp.h>
#include <freerdp/codec/color.h>
#include <freerdp/gdi/gdi.h>
#include <freerdp/settings.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>

// Helper function to convert color using FreeRDP 3.x API
static inline UINT32 convertColor(UINT32 color, UINT32 srcBpp, UINT32 dstBpp) {
    UINT32 srcFormat = (srcBpp == 16) ? PIXEL_FORMAT_RGB16 : PIXEL_FORMAT_BGRX32;
    UINT32 dstFormat = (dstBpp == 32) ? PIXEL_FORMAT_BGRX32 : PIXEL_FORMAT_RGB16;
    return FreeRDPConvertColor(color, srcFormat, dstFormat, NULL);
}

// Helper function to flip image data
static inline void flipImageData(BYTE* data, int width, int height, int bpp) {
    int scanline = width * (bpp / 8);
    BYTE* tmpLine = (BYTE*)malloc(scanline);
    if (!tmpLine) return;
    for (int i = 0; i < height / 2; i++) {
        BYTE* line1 = data + (i * scanline);
        BYTE* line2 = data + ((height - 1 - i) * scanline);
        memcpy(tmpLine, line1, scanline);
        memcpy(line1, line2, scanline);
        memcpy(line2, tmpLine, scanline);
    }
    free(tmpLine);
}

// Helper function to get settings from instance
static inline rdpSettings* getSettings(freerdp* instance) {
    return instance->context->settings;
}

extern BOOL preConnect(freerdp* instance);
extern void postConnect(freerdp* instance);
extern void goPrintln(char* text);
extern void goEcho(char* text, rdpContext* context);
extern size_t getPointerSize();
extern BOOL primaryPatBlt(rdpContext* context, PATBLT_ORDER* patblt);
extern BOOL primaryScrBlt(rdpContext* context, SCRBLT_ORDER* scrblt);
extern BOOL primaryOpaqueRect(rdpContext* context, OPAQUE_RECT_ORDER* oro);
extern BOOL primaryMultiOpaqueRect(rdpContext* context, MULTI_OPAQUE_RECT_ORDER* moro);
extern BOOL beginPaint(rdpContext* context);
extern BOOL endPaint(rdpContext* context);
extern BOOL setBounds(rdpContext* context, rdpBounds* bounds);
extern BOOL bitmapUpdate(rdpContext* context, BITMAP_UPDATE* bitmap);

static BOOL cbPrimaryPatBlt(rdpContext* context, PATBLT_ORDER* patblt) {
	return primaryPatBlt(context, patblt);
}

static BOOL cbPrimaryScrBlt(rdpContext* context, const SCRBLT_ORDER* scrblt) {
	return primaryScrBlt(context, (SCRBLT_ORDER*)scrblt);
}

static BOOL cbPrimaryOpaqueRect(rdpContext* context, const OPAQUE_RECT_ORDER* oro) {
	return primaryOpaqueRect(context, (OPAQUE_RECT_ORDER*)oro);
}

static BOOL cbPrimaryMultiOpaqueRect(rdpContext* context, const MULTI_OPAQUE_RECT_ORDER* moro) {
	return primaryMultiOpaqueRect(context, (MULTI_OPAQUE_RECT_ORDER*)moro);
}

static BOOL cbBeginPaint(rdpContext* context) {
	return beginPaint(context);
}
static BOOL cbEndPaint(rdpContext* context) {
	return endPaint(context);
}
static BOOL cbSetBounds(rdpContext* context, const rdpBounds* bounds) {
	return setBounds(context, (rdpBounds*)bounds);
}
static BOOL cbBitmapUpdate(rdpContext* context, const BITMAP_UPDATE* bitmap) {
	return bitmapUpdate(context, (BITMAP_UPDATE*)bitmap);
}

static BOOL cbPreConnect(freerdp* instance) {
	rdpContext* context = instance->context;
	rdpUpdate* update = context->update;
	rdpPrimaryUpdate* primary = update->primary;

	primary->PatBlt = cbPrimaryPatBlt;
	primary->ScrBlt = cbPrimaryScrBlt;
	primary->OpaqueRect = cbPrimaryOpaqueRect;
	primary->MultiOpaqueRect = cbPrimaryMultiOpaqueRect;

	update->BeginPaint = cbBeginPaint;
	update->EndPaint = cbEndPaint;
	update->SetBounds = cbSetBounds;
	update->BitmapUpdate = cbBitmapUpdate;

	return preConnect(instance);
}

static BOOL cbPostConnect(freerdp* instance) {
	postConnect(instance);

	rdpPointer p;
	memset(&p, 0, sizeof(p));

	p.size = getPointerSize();

	//p.New = cbPointer_New;
	//p.Free = cbPointer_Free;
	//p.Set = cbPointer_Set;
	//p.SetNull = cbPointer_SetNull;
	//p.SetDefault = cbPointer_SetDefault;

	graphics_register_pointer(instance->context->graphics, &p);

	return 1;
}

static BITMAP_DATA* nextBitmapRectangle(BITMAP_UPDATE* bitmap, int i) {
	return &bitmap->rectangles[i];
}

static DELTA_RECT* nextMultiOpaqueRectangle(MULTI_OPAQUE_RECT_ORDER* moro, int i) {
	return &moro->rectangles[i];
}

static void bindCallbacks(freerdp* instance) {
	instance->PreConnect = cbPreConnect;
	instance->PostConnect = cbPostConnect;
}
*/
import (
	"C"
)
import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"
)

const (
	WSOP_SC_BEGINPAINT       uint32 = 0
	WSOP_SC_ENDPAINT         uint32 = 1
	WSOP_SC_BITMAP           uint32 = 2
	WSOP_SC_OPAQUERECT       uint32 = 3
	WSOP_SC_SETBOUNDS        uint32 = 4
	WSOP_SC_PATBLT           uint32 = 5
	WSOP_SC_MULTI_OPAQUERECT uint32 = 6
	WSOP_SC_SCRBLT           uint32 = 7
	WSOP_SC_PTR_NEW          uint32 = 8
	WSOP_SC_PTR_FREE         uint32 = 9
	WSOP_SC_PTR_SET          uint32 = 10
	WSOP_SC_PTR_SETNULL      uint32 = 11
	WSOP_SC_PTR_SETDEFAULT   uint32 = 12
)

type bitmapUpdateMeta struct {
	op  uint32
	x   uint32
	y   uint32
	w   uint32
	h   uint32
	dw  uint32
	dh  uint32
	bpp uint32
	cf  uint32
	sz  uint32
}

type primaryPatBltMeta struct {
	op  uint32
	x   int32
	y   int32
	w   int32
	h   int32
	fg  uint32
	rop uint32
}

type primaryScrBltMeta struct {
	op  uint32
	rop uint32
	x   int32
	y   int32
	w   int32
	h   int32
	sx  int32
	sy  int32
}

type rdpConnectionSettings struct {
	hostname *string
	username *string
	password *string
	width    int
	height   int
	port     int
}

type rdpContext struct {
	_p       C.rdpContext
	recvq    chan []byte
	sendq    chan []byte
	settings *rdpConnectionSettings
}

type rdpPointer struct {
	pointer *C.rdpPointer
	id      int
}

func rdpconnect(sendq chan []byte, recvq chan []byte, settings *rdpConnectionSettings) {
	var instance *C.freerdp

	fmt.Println("RDP Connecting...")

	instance = C.freerdp_new()

	C.bindCallbacks(instance)
	instance.ContextSize = C.size_t(unsafe.Sizeof(rdpContext{}))
	C.freerdp_context_new(instance)

	var context *rdpContext
	context = (*rdpContext)(unsafe.Pointer(instance.context))
	context.sendq = sendq
	context.recvq = recvq
	context.settings = settings

	C.freerdp_connect(instance)

	mainEventLoop := true

	for mainEventLoop {
		select {
		case <-recvq:
			fmt.Println("Disconnecting (websocket error)")
			mainEventLoop = false
		default:
			e := int(C.freerdp_error_info(instance))
			if e != 0 {
				switch e {
				case 1:
				case 2:
				case 7:
				case 9:
					// Manual disconnections and such
					fmt.Println("Disconnecting (manual)")
					mainEventLoop = false
					break
				case 5:
					// Another user connected
					break
				default:
					// Unknown error?
					break
				}
			}
			if int(C.freerdp_shall_disconnect(instance)) != 0 {
				fmt.Println("Disconnecting (RDC said so)")
				mainEventLoop = false
			}
			if mainEventLoop {
				C.freerdp_check_fds(instance)
			}
			C.usleep(1000)
		}
	}
	C.freerdp_free(instance)
}

func sendBinary(sendq chan []byte, data *bytes.Buffer) {
	sendq <- data.Bytes()
}

//export getPointerSize
func getPointerSize() C.size_t {
	return C.size_t(unsafe.Sizeof(rdpPointer{}))
}

//export primaryPatBlt
func primaryPatBlt(rawContext *C.rdpContext, patblt *C.PATBLT_ORDER) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))

	if C.GDI_BS_SOLID == patblt.brush.style {
		// Convert color from 16-bit to 32-bit using FreeRDP 3.x API
		color := uint32(C.convertColor(patblt.foreColor, 16, 32))

		meta := primaryPatBltMeta{
			WSOP_SC_PATBLT,
			int32(patblt.nLeftRect),
			int32(patblt.nTopRect),
			int32(patblt.nWidth),
			int32(patblt.nHeight),
			color,
			uint32(patblt.bRop),
		}

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.LittleEndian, meta)
		sendBinary(context.sendq, buf)
	}
	return C.TRUE
}

//export primaryScrBlt
func primaryScrBlt(rawContext *C.rdpContext, scrblt *C.SCRBLT_ORDER) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))

	meta := primaryScrBltMeta{
		WSOP_SC_SCRBLT,
		uint32(scrblt.bRop),
		int32(scrblt.nLeftRect),
		int32(scrblt.nTopRect),
		int32(scrblt.nWidth),
		int32(scrblt.nHeight),
		int32(scrblt.nXSrc),
		int32(scrblt.nYSrc),
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, meta)
	sendBinary(context.sendq, buf)
	return C.TRUE
}

//export primaryOpaqueRect
func primaryOpaqueRect(rawContext *C.rdpContext, oro *C.OPAQUE_RECT_ORDER) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))

	// Note: oro is const in C, so we can't modify it directly
	// Convert color from 16-bit to 32-bit
	color := C.convertColor(oro.color, 16, 32)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, WSOP_SC_OPAQUERECT)

	// Create a copy with converted color
	type opaqueRectOrder struct {
		nLeftRect  int32
		nTopRect   int32
		nWidth     int32
		nHeight    int32
		color      uint32
	}

	order := opaqueRectOrder{
		nLeftRect: int32(oro.nLeftRect),
		nTopRect:  int32(oro.nTopRect),
		nWidth:    int32(oro.nWidth),
		nHeight:   int32(oro.nHeight),
		color:     uint32(color),
	}

	binary.Write(buf, binary.LittleEndian, order)
	sendBinary(context.sendq, buf)
	return C.TRUE
}

//export primaryMultiOpaqueRect
func primaryMultiOpaqueRect(rawContext *C.rdpContext, moro *C.MULTI_OPAQUE_RECT_ORDER) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))

	// Convert color from 16-bit to 32-bit using FreeRDP 3.x API
	color := C.convertColor(moro.color, 16, 32)

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, WSOP_SC_MULTI_OPAQUERECT)
	binary.Write(buf, binary.LittleEndian, int32(color))
	binary.Write(buf, binary.LittleEndian, int32(moro.numRectangles))

	var r *C.DELTA_RECT
	var i int
	for i = 1; i <= int(moro.numRectangles); i++ {
		r = C.nextMultiOpaqueRectangle(moro, C.int(i))
		binary.Write(buf, binary.LittleEndian, r)
	}

	sendBinary(context.sendq, buf)
	return C.TRUE
}

//export beginPaint
func beginPaint(rawContext *C.rdpContext) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, WSOP_SC_BEGINPAINT)
	sendBinary(context.sendq, buf)
	return C.TRUE
}

//export endPaint
func endPaint(rawContext *C.rdpContext) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, WSOP_SC_ENDPAINT)
	sendBinary(context.sendq, buf)
	return C.TRUE
}

//export setBounds
func setBounds(rawContext *C.rdpContext, bounds *C.rdpBounds) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))
	buf := new(bytes.Buffer)

	if bounds != nil {
		binary.Write(buf, binary.LittleEndian, WSOP_SC_SETBOUNDS)
		binary.Write(buf, binary.LittleEndian, bounds)
		sendBinary(context.sendq, buf)
	}
	return C.TRUE
}

//export bitmapUpdate
func bitmapUpdate(rawContext *C.rdpContext, bitmap *C.BITMAP_UPDATE) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(rawContext))

	var bmd *C.BITMAP_DATA
	var i int

	for i = 0; i < int(bitmap.number); i++ {
		bmd = C.nextBitmapRectangle(bitmap, C.int(i))

		buf := new(bytes.Buffer)

		meta := bitmapUpdateMeta{
			WSOP_SC_BITMAP,                           // op
			uint32(bmd.destLeft),                     // x
			uint32(bmd.destTop),                      // y
			uint32(bmd.width),                        // w
			uint32(bmd.height),                       // h
			uint32(bmd.destRight - bmd.destLeft + 1), // dw
			uint32(bmd.destBottom - bmd.destTop + 1), // dh
			uint32(bmd.bitsPerPixel),                 // bpp
			uint32(bmd.compressed),                   // cf
			uint32(bmd.bitmapLength),                 // sz
		}
		if int(bmd.compressed) == 0 {
			// Use FreeRDP 3.x helper function to flip image
			C.flipImageData(bmd.bitmapDataStream, C.int(bmd.width), C.int(bmd.height), C.int(bmd.bitsPerPixel))
		}

		binary.Write(buf, binary.LittleEndian, meta)

		// Unsafe copy bmd.bitmapLength bytes out of bmd.bitmapDataStream
		var bitmapDataStream []byte
		clen := int(bmd.bitmapLength)
		bitmapDataStream = (*[1 << 30]byte)(unsafe.Pointer(bmd.bitmapDataStream))[:clen]
		(*reflect.SliceHeader)(unsafe.Pointer(&bitmapDataStream)).Cap = clen
		binary.Write(buf, binary.LittleEndian, bitmapDataStream)

		sendBinary(context.sendq, buf)
	}
	return C.TRUE
}

//export postConnect
func postConnect(instance *C.freerdp) {
	fmt.Println("Connected.")
}

//export preConnect
func preConnect(instance *C.freerdp) C.BOOL {
	context := (*rdpContext)(unsafe.Pointer(instance.context))
	settings := C.getSettings(instance)

	// 设置连接参数 - 使用 FreeRDP 3.x settings API
	hostname := C.CString(*context.settings.hostname)
	username := C.CString(*context.settings.username)
	password := C.CString(*context.settings.password)
	defer C.free(unsafe.Pointer(hostname))
	defer C.free(unsafe.Pointer(username))
	defer C.free(unsafe.Pointer(password))

	C.freerdp_settings_set_string(settings, C.FreeRDP_ServerHostname, hostname)
	C.freerdp_settings_set_string(settings, C.FreeRDP_Username, username)
	C.freerdp_settings_set_string(settings, C.FreeRDP_Password, password)
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_DesktopWidth, C.UINT32(context.settings.width))
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_DesktopHeight, C.UINT32(context.settings.height))
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_ServerPort, C.UINT32(context.settings.port))
	C.freerdp_settings_set_bool(settings, C.FreeRDP_IgnoreCertificate, C.TRUE)
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_ColorDepth, 16)

	// Performance flags
	perfFlags := C.PERF_DISABLE_WALLPAPER | C.PERF_DISABLE_THEMING |
		C.PERF_DISABLE_MENUANIMATIONS | C.PERF_DISABLE_FULLWINDOWDRAG
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_PerformanceFlags, C.UINT32(perfFlags))

	// Connection type
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_ConnectionType, C.CONNECTION_TYPE_BROADBAND_HIGH)

	// Other settings
	C.freerdp_settings_set_bool(settings, C.FreeRDP_RemoteFxCodec, C.FALSE)
	C.freerdp_settings_set_bool(settings, C.FreeRDP_FastPathOutput, C.TRUE)
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_FrameAcknowledge, 1)
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_LargePointerFlag, 1)
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_GlyphSupportLevel, C.GLYPH_SUPPORT_FULL)
	C.freerdp_settings_set_bool(settings, C.FreeRDP_BitmapCacheV3Enabled, C.FALSE)
	C.freerdp_settings_set_uint32(settings, C.FreeRDP_OffscreenSupportLevel, 0)

	// Order support settings are typically handled internally by FreeRDP 3.x

	return C.TRUE
}
