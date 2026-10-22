package buildtest

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"testing"
)

// Hero source fixtures that need more than a plain PNG: WebP sources (which
// the build must reject), JPEGs with an EXIF orientation, and rasters with
// an embedded ICC colour profile. Everything is built by hand so no fixture
// needs an encoder the site itself does not use.

// HeroWebP returns a solid lossless WebP (VP8L) of the given size. With a
// single colour every prefix code has one symbol, so the pixel data takes
// zero bits and the file stays a few dozen bytes at any size.
func HeroWebP(w, h int, c color.RGBA) []byte {
	var bw bitWriter
	bw.write(0x2f, 8) // VP8L signature
	bw.write(uint32(w-1), 14)
	bw.write(uint32(h-1), 14)
	bw.write(0, 1) // alpha_is_used
	bw.write(0, 3) // version
	bw.write(0, 1) // no transform
	bw.write(0, 1) // no colour cache
	bw.write(0, 1) // no meta prefix codes
	for _, sym := range []uint8{c.G, c.R, c.B, c.A} {
		bw.write(1, 1) // simple code
		bw.write(0, 1) // one symbol
		bw.write(1, 1) // 8-bit symbol
		bw.write(uint32(sym), 8)
	}
	bw.write(1, 1) // distance: simple code
	bw.write(0, 1) // one symbol
	bw.write(0, 1) // 1-bit symbol
	bw.write(0, 1) // symbol 0
	data := bw.bytes()

	var chunk bytes.Buffer
	chunk.WriteString("VP8L")
	binary.Write(&chunk, binary.LittleEndian, uint32(len(data)))
	chunk.Write(data)
	if len(data)%2 == 1 {
		chunk.WriteByte(0)
	}
	var out bytes.Buffer
	out.WriteString("RIFF")
	binary.Write(&out, binary.LittleEndian, uint32(4+chunk.Len()))
	out.WriteString("WEBP")
	out.Write(chunk.Bytes())
	return out.Bytes()
}

// bitWriter packs bits least-significant first, as VP8L reads them.
type bitWriter struct {
	buf   []byte
	nbits uint
}

func (b *bitWriter) write(v uint32, n uint) {
	for i := uint(0); i < n; i++ {
		if b.nbits%8 == 0 {
			b.buf = append(b.buf, 0)
		}
		if v&(1<<i) != 0 {
			b.buf[len(b.buf)-1] |= 1 << (b.nbits % 8)
		}
		b.nbits++
	}
}

func (b *bitWriter) bytes() []byte { return b.buf }

// HeroJPEG encodes img as a quality-90 baseline JPEG (no metadata).
func HeroJPEG(img image.Image) []byte {
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

// WithEXIFOrientation inserts an APP1 Exif segment carrying only the
// Orientation tag (0x0112) right after a JPEG's SOI marker. TIFF byte order
// is big-endian ("MM") when bigEndian is set, else little-endian ("II").
func WithEXIFOrientation(jpg []byte, orientation uint16, bigEndian bool) []byte {
	var order binary.AppendByteOrder = binary.LittleEndian
	tiff := []byte("II*\x00")
	if bigEndian {
		order, tiff = binary.BigEndian, []byte("MM\x00*")
	}
	tiff = order.AppendUint32(tiff, 8)      // IFD0 offset
	tiff = order.AppendUint16(tiff, 1)      // one entry
	tiff = order.AppendUint16(tiff, 0x0112) // Orientation
	tiff = order.AppendUint16(tiff, 3)      // SHORT
	tiff = order.AppendUint32(tiff, 1)      // count
	tiff = order.AppendUint16(tiff, orientation)
	tiff = order.AppendUint16(tiff, 0) // value padding
	tiff = order.AppendUint32(tiff, 0) // no next IFD
	return withJPEGSegment(jpg, 0xE1, append([]byte("Exif\x00\x00"), tiff...))
}

// withJPEGSegment inserts one APPn segment after the SOI marker.
func withJPEGSegment(jpg []byte, marker byte, payload []byte) []byte {
	seg := []byte{0xFF, marker}
	seg = binary.BigEndian.AppendUint16(seg, uint16(len(payload)+2))
	seg = append(seg, payload...)
	out := append([]byte{}, jpg[:2]...)
	out = append(out, seg...)
	return append(out, jpg[2:]...)
}

// ICCProfile returns a minimal RGB display ICC profile whose only tag is its
// description: a v2 'desc' (textDescriptionType) or, when v4 is set, a v4
// 'mluc' (multiLocalizedUnicodeType, en-US), as Apple's Display P3 uses.
// The build reads only the description, so no colorants are needed.
func ICCProfile(description string, v4 bool) []byte {
	var tag []byte
	if v4 {
		var utf16 []byte
		for _, r := range description {
			utf16 = binary.BigEndian.AppendUint16(utf16, uint16(r))
		}
		tag = append([]byte("mluc"), 0, 0, 0, 0)
		tag = binary.BigEndian.AppendUint32(tag, 1)  // records
		tag = binary.BigEndian.AppendUint32(tag, 12) // record size
		tag = append(tag, "enUS"...)
		tag = binary.BigEndian.AppendUint32(tag, uint32(len(utf16)))
		tag = binary.BigEndian.AppendUint32(tag, 28) // string offset in the tag
		tag = append(tag, utf16...)
	} else {
		tag = append([]byte("desc"), 0, 0, 0, 0)
		tag = binary.BigEndian.AppendUint32(tag, uint32(len(description)+1))
		tag = append(tag, description...)
		tag = append(tag, 0)
		tag = append(tag, make([]byte, 4+4+2+1+67)...) // empty Unicode and ScriptCode
	}
	for len(tag)%4 != 0 {
		tag = append(tag, 0)
	}
	const tagOffset = 128 + 4 + 12
	size := tagOffset + len(tag)
	header := make([]byte, 128)
	binary.BigEndian.PutUint32(header[0:], uint32(size))
	version := uint32(0x02100000)
	if v4 {
		version = 0x04300000
	}
	binary.BigEndian.PutUint32(header[8:], version)
	copy(header[12:], "mntr")
	copy(header[16:], "RGB ")
	copy(header[20:], "XYZ ")
	copy(header[36:], "acsp")
	p := append(header, 0, 0, 0, 1) // tag count
	p = append(p, "desc"...)
	p = binary.BigEndian.AppendUint32(p, tagOffset)
	p = binary.BigEndian.AppendUint32(p, uint32(len(tag)))
	return append(p, tag...)
}

// WithJPEGICCProfile embeds profile in a JPEG as one APP2 ICC_PROFILE
// segment (sequence 1 of 1).
func WithJPEGICCProfile(jpg, profile []byte) []byte {
	payload := append([]byte("ICC_PROFILE\x00\x01\x01"), profile...)
	return withJPEGSegment(jpg, 0xE2, payload)
}

// WithPNGICCProfile embeds profile in a PNG as an iCCP chunk right after
// IHDR.
func WithPNGICCProfile(pngData, profile []byte) []byte {
	var z bytes.Buffer
	zw := zlib.NewWriter(&z)
	zw.Write(profile)
	zw.Close()
	body := append([]byte("ICC Profile\x00\x00"), z.Bytes()...) // name, compression 0
	chunk := binary.BigEndian.AppendUint32(nil, uint32(len(body)))
	chunk = append(chunk, "iCCP"...)
	chunk = append(chunk, body...)
	chunk = binary.BigEndian.AppendUint32(chunk, crc32.ChecksumIEEE(chunk[4:]))
	const afterIHDR = 8 + 4 + 4 + 13 + 4
	out := append([]byte{}, pngData[:afterIHDR]...)
	out = append(out, chunk...)
	return append(out, pngData[afterIHDR:]...)
}

// DecodePixels fully decodes the image written for urlPath. The test fails
// if the file is missing or does not decode.
func DecodePixels(tb testing.TB, outDir, urlPath string) image.Image {
	tb.Helper()
	data, err := os.ReadFile(StaticPath(outDir, urlPath))
	if err != nil {
		tb.Fatalf("image %s: %v", urlPath, err)
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		tb.Fatalf("image %s does not decode: %v", urlPath, err)
	}
	return img
}

// Quadrants is a w x h image whose top-left quadrant is red and the rest
// blue: rotate or flip it and the red corner moves, so a test can see which
// way up a decoded variant is.
func Quadrants(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	red := color.RGBA{R: 0xDC, G: 0x26, B: 0x26, A: 0xFF}
	blue := color.RGBA{R: 0x25, G: 0x63, B: 0xEB, A: 0xFF}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			if x < w/2 && y < h/2 {
				img.SetRGBA(x, y, red)
			} else {
				img.SetRGBA(x, y, blue)
			}
		}
	}
	return img
}
