package images

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"io"
	"os"
	"strings"
	"unicode/utf16"

	"golang.org/x/image/draw"
)

// Source rules. Only PNG and JPEG sources are accepted: Go's WebP decoders
// (golang.org/x/image/webp and gen2brain/webp) hand libwebp's limited-range
// YCbCr to a full-range conversion, so every WebP source would come out
// with washed-out contrast (golang/go#40173).

// ErrWebP is returned by Inspect for a WebP source, whatever its file name.
var ErrWebP = errors.New("WebP sources are not supported (Go decodes WebP with the wrong colour range): export it as PNG or JPEG")

// Info describes a source raster as it should be displayed.
type Info struct {
	// Width and Height are the upright size, after EXIF orientation.
	Width, Height int
	// Format is "png" or "jpeg".
	Format string
	// Orientation is the EXIF orientation (1-8; 1 when absent). Encoder
	// applies it before resizing.
	Orientation int
	// Profile is the embedded ICC profile's description ("" if none).
	Profile string
	// SRGB is false when the source carries a colour profile that is not
	// sRGB (or grey): its pixels will be read as sRGB and look wrong.
	SRGB bool
}

// Inspect reads a source raster's header and metadata. It fails for
// anything but a PNG or JPEG (ErrWebP for WebP, sniffed from the bytes, not
// the name).
func Inspect(path string) (Info, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Info{}, err
	}
	format := sniff(data)
	switch format {
	case "png", "jpeg":
	case "webp":
		return Info{}, ErrWebP
	default:
		return Info{}, errors.New("not a supported raster image (PNG or JPEG)")
	}
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return Info{}, fmt.Errorf("not a supported raster image (PNG or JPEG): %w", err)
	}
	info := Info{Width: cfg.Width, Height: cfg.Height, Format: format, Orientation: orientationOf(data)}
	info.Profile, info.SRGB = colourProfile(data)
	if info.Orientation >= 5 { // 5-8 swap the axes
		info.Width, info.Height = info.Height, info.Width
	}
	return info, nil
}

// --- EXIF orientation -------------------------------------------------
//
// Go's decoders ignore EXIF (golang/go#4341), so a phone photo stored
// sideways with Orientation 6 or 8 would be resized sideways. The tag is
// read here (pure Go) and applied right after decoding.

// orientationOf is the EXIF Orientation (1-8) of a JPEG (APP1 "Exif") or
// PNG (eXIf chunk), or 1 when there is none or it is malformed.
func orientationOf(data []byte) int {
	var tiff []byte
	switch sniff(data) {
	case "jpeg":
		for _, seg := range jpegSegments(data) {
			if seg.marker == 0xE1 && bytes.HasPrefix(seg.payload, []byte("Exif\x00\x00")) {
				tiff = seg.payload[6:]
				break
			}
		}
	case "png":
		for _, c := range pngChunks(data) {
			if c.typ == "eXIf" {
				tiff = c.data
				break
			}
		}
	}
	if o := tiffOrientation(tiff); o >= 1 && o <= 8 {
		return o
	}
	return 1
}

// tiffOrientation finds tag 0x0112 in IFD0 of a TIFF structure, or 0.
func tiffOrientation(t []byte) int {
	if len(t) < 8 {
		return 0
	}
	var order binary.ByteOrder
	switch string(t[:2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return 0
	}
	if order.Uint16(t[2:]) != 42 {
		return 0
	}
	ifd := int(order.Uint32(t[4:]))
	if ifd < 8 || ifd+2 > len(t) {
		return 0
	}
	n := int(order.Uint16(t[ifd:]))
	for i := range n {
		e := ifd + 2 + 12*i
		if e+12 > len(t) {
			return 0
		}
		if order.Uint16(t[e:]) == 0x0112 && order.Uint16(t[e+2:]) == 3 { // SHORT
			return int(order.Uint16(t[e+8:]))
		}
	}
	return 0
}

// upright applies an EXIF orientation to decoded pixels.
func upright(img image.Image, orientation int) image.Image {
	if orientation <= 1 || orientation > 8 {
		return img
	}
	b := img.Bounds()
	src := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(src, src.Bounds(), img, b.Min, draw.Src)
	sw, sh := b.Dx(), b.Dy()
	dw, dh := sw, sh
	if orientation >= 5 {
		dw, dh = sh, sw
	}
	dst := image.NewRGBA(image.Rect(0, 0, dw, dh))
	// from maps a destination pixel to the stored pixel it shows (EXIF 2.3,
	// Orientation: which visual side row 0 and column 0 of the stored image
	// are).
	from := func(dx, dy int) (int, int) {
		switch orientation {
		case 2: // mirrored horizontally
			return sw - 1 - dx, dy
		case 3: // rotated 180
			return sw - 1 - dx, sh - 1 - dy
		case 4: // mirrored vertically
			return dx, sh - 1 - dy
		case 5: // transposed
			return dy, dx
		case 6: // display rotates 90 clockwise
			return dy, sh - 1 - dx
		case 7: // transversed
			return sw - 1 - dy, sh - 1 - dx
		default: // 8: display rotates 90 counter-clockwise
			return sw - 1 - dy, dx
		}
	}
	for dy := range dh {
		for dx := range dw {
			sx, sy := from(dx, dy)
			copy(dst.Pix[dst.PixOffset(dx, dy):][:4], src.Pix[src.PixOffset(sx, sy):][:4])
		}
	}
	return dst
}

// --- ICC colour profile -------------------------------------------------
//
// Go's decoders drop embedded ICC profiles (golang/go#23936) and every
// output is served untagged, i.e. read as sRGB. A Display P3 or Adobe RGB
// source therefore comes out desaturated. Nothing is converted: the build
// only names the profile so the Author can re-export as sRGB.

// colourProfile reports the embedded ICC profile's description ("" when
// there is none or it has no readable description) and whether the pixels
// can be taken as sRGB: no profile, a PNG sRGB chunk, a profile described
// as sRGB, or a grey profile.
func colourProfile(data []byte) (name string, srgb bool) {
	var icc []byte
	switch sniff(data) {
	case "jpeg":
		icc = jpegICC(data)
	case "png":
		for _, c := range pngChunks(data) {
			switch c.typ {
			case "sRGB":
				return "sRGB", true
			case "iCCP":
				icc = pngICC(c.data)
			}
		}
	}
	if icc == nil {
		return "", true
	}
	name = iccDescription(icc)
	if len(icc) >= 20 && string(icc[16:20]) == "GRAY" {
		return name, true
	}
	lower := strings.ToLower(name)
	return name, strings.Contains(lower, "srgb") || strings.Contains(lower, "61966-2")
}

// jpegICC reassembles an ICC profile split over APP2 "ICC_PROFILE"
// segments (each carries its 1-based sequence number and the total).
func jpegICC(data []byte) []byte {
	parts := map[int][]byte{}
	total := 0
	for _, seg := range jpegSegments(data) {
		const sig = "ICC_PROFILE\x00"
		if seg.marker != 0xE2 || len(seg.payload) < len(sig)+2 || string(seg.payload[:len(sig)]) != sig {
			continue
		}
		parts[int(seg.payload[len(sig)])] = seg.payload[len(sig)+2:]
		total = int(seg.payload[len(sig)+1])
	}
	if len(parts) == 0 {
		return nil
	}
	var icc []byte
	for i := 1; i <= total; i++ {
		icc = append(icc, parts[i]...)
	}
	if icc == nil {
		icc = []byte{} // present but unreadable: still a profile
	}
	return icc
}

// pngICC inflates an iCCP chunk: name, NUL, compression method 0, zlib.
func pngICC(chunk []byte) []byte {
	i := bytes.IndexByte(chunk, 0)
	if i < 0 || i+2 > len(chunk) {
		return []byte{}
	}
	zr, err := zlib.NewReader(bytes.NewReader(chunk[i+2:]))
	if err != nil {
		return []byte{}
	}
	defer zr.Close()
	icc, err := io.ReadAll(io.LimitReader(zr, 16<<20))
	if err != nil {
		return []byte{}
	}
	return icc
}

// iccDescription reads the profile's 'desc' tag: a v2 textDescriptionType
// (ASCII) or a v4 multiLocalizedUnicodeType (UTF-16BE, en-US preferred).
func iccDescription(icc []byte) string {
	if len(icc) < 132 {
		return ""
	}
	be := binary.BigEndian
	n := int(be.Uint32(icc[128:]))
	for i := range n {
		e := 132 + 12*i
		if e+12 > len(icc) {
			return ""
		}
		if string(icc[e:e+4]) != "desc" {
			continue
		}
		off, size := int(be.Uint32(icc[e+4:])), int(be.Uint32(icc[e+8:]))
		if off < 0 || size < 12 || off+size > len(icc) || off+size < off {
			return ""
		}
		tag := icc[off : off+size]
		switch string(tag[:4]) {
		case "desc":
			count := int(be.Uint32(tag[8:]))
			if count < 0 || 12+count > len(tag) {
				return ""
			}
			return strings.TrimRight(string(tag[12:12+count]), "\x00 ")
		case "mluc":
			if len(tag) < 16 {
				return ""
			}
			records, recSize := int(be.Uint32(tag[8:])), int(be.Uint32(tag[12:]))
			best := ""
			for r := range records {
				rec := 16 + r*recSize
				if recSize < 12 || rec+12 > len(tag) {
					break
				}
				length, at := int(be.Uint32(tag[rec+4:])), int(be.Uint32(tag[rec+8:]))
				if at < 0 || length < 0 || at+length > len(tag) {
					continue
				}
				units := make([]uint16, length/2)
				for k := range units {
					units[k] = be.Uint16(tag[at+2*k:])
				}
				s := strings.TrimRight(string(utf16.Decode(units)), "\x00 ")
				if best == "" || string(tag[rec:rec+4]) == "enUS" {
					best = s
				}
			}
			return best
		}
		return ""
	}
	return ""
}

// --- container walking --------------------------------------------------

type jpegSegment struct {
	marker  byte
	payload []byte
}

// jpegSegments lists a JPEG's marker segments up to the first scan.
func jpegSegments(data []byte) []jpegSegment {
	var segs []jpegSegment
	i := 2 // after SOI
	for i+4 <= len(data) {
		if data[i] != 0xFF {
			return segs
		}
		marker := data[i+1]
		switch {
		case marker == 0xFF: // fill byte
			i++
			continue
		case marker == 0x01 || marker >= 0xD0 && marker <= 0xD7: // no length
			i += 2
			continue
		case marker == 0xDA || marker == 0xD9: // start of scan, end of image
			return segs
		}
		n := int(binary.BigEndian.Uint16(data[i+2:]))
		if n < 2 || i+2+n > len(data) {
			return segs
		}
		segs = append(segs, jpegSegment{marker: marker, payload: data[i+4 : i+2+n]})
		i += 2 + n
	}
	return segs
}

type pngChunk struct {
	typ  string
	data []byte
}

// pngChunks lists a PNG's chunks up to the first IDAT (where iCCP and
// eXIf must appear).
func pngChunks(data []byte) []pngChunk {
	var chunks []pngChunk
	i := 8 // after the signature
	for i+8 <= len(data) {
		n := int(binary.BigEndian.Uint32(data[i:]))
		typ := string(data[i+4 : i+8])
		if n < 0 || i+12+n > len(data) || typ == "IDAT" {
			return chunks
		}
		chunks = append(chunks, pngChunk{typ: typ, data: data[i+8 : i+8+n]})
		i += 12 + n
	}
	return chunks
}

// sniff names a raster format from its magic bytes: "png", "jpeg", "webp"
// or "".
func sniff(data []byte) string {
	switch {
	case bytes.HasPrefix(data, []byte("\x89PNG\r\n\x1a\n")):
		return "png"
	case bytes.HasPrefix(data, []byte{0xFF, 0xD8, 0xFF}):
		return "jpeg"
	case len(data) >= 12 && string(data[:4]) == "RIFF" && string(data[8:12]) == "WEBP":
		return "webp"
	}
	return ""
}
