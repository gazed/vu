package gltf

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"reflect"
	"testing"
	"testing/fstest"
)

func readFile(path string) []byte {
	r, _ := ioutil.ReadFile(path)
	return r
}

func TestDecoder_decodeBuffer(t *testing.T) {
	type args struct {
		buffer *Buffer
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		want    []byte
		wantErr bool
	}{
		{"byteLength_0", &Decoder{}, args{&Buffer{ByteLength: 0, URI: "a.bin"}}, nil, true},
		{"noURI", &Decoder{}, args{&Buffer{ByteLength: 1, URI: ""}}, nil, true},
		{"invalidURI", &Decoder{}, args{&Buffer{ByteLength: 1, URI: "../a.bin"}}, nil, true},
		{"noSchemeErr", NewDecoder(nil), args{&Buffer{ByteLength: 3, URI: "ftp://a.bin"}}, nil, false},
		{"base", NewDecoderFS(nil, fstest.MapFS{"a.bin": &fstest.MapFile{Data: []byte("abcdfg")}}), args{&Buffer{ByteLength: 6, URI: "a.bin"}}, []byte("abcdfg"), false},
		{"dotdot", NewDecoderFS(nil, fstest.MapFS{"a..b.bin": &fstest.MapFile{Data: []byte("abcdfg")}}), args{&Buffer{ByteLength: 6, URI: "a..b.bin"}}, []byte("abcdfg"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.decodeBuffer(tt.args.buffer); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.buffer.Data, tt.want) {
				t.Errorf("Decoder.decodeBuffer() buffer = %v, want %v", tt.args.buffer.Data, tt.want)
			}
		})
	}
}

func TestDecoder_decodeBinaryBuffer(t *testing.T) {
	type args struct {
		buffer *Buffer
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		want    []byte
		wantErr bool
	}{
		{"base", NewDecoder(bytes.NewBuffer([]byte{0x06, 0x00, 0x00, 0x00, 0x42, 0x49, 0x4e, 0x00, 1, 2, 3, 4, 5, 6})),
			args{&Buffer{ByteLength: 6}}, []byte{1, 2, 3, 4, 5, 6}, false},
		{"smallbuffer", NewDecoder(bytes.NewBuffer([]byte{0x6, 0x00, 0x00, 0x00, 0x42, 0x49, 0x4e, 0x00, 1, 2, 3, 4, 5, 6})),
			args{&Buffer{ByteLength: 5}}, []byte{1, 2, 3, 4, 5}, false},
		{"bigbuffer", NewDecoder(bytes.NewBuffer([]byte{0x6, 0x00, 0x00, 0x00, 0x42, 0x49, 0x4e, 0x00, 1, 2, 3, 4, 5, 6})),
			args{&Buffer{ByteLength: 7}}, nil, true},
		{"invalidBuffer", new(Decoder), args{&Buffer{ByteLength: 0}}, nil, true},
		{"readErr", NewDecoder(bytes.NewBufferString("")), args{&Buffer{ByteLength: 1}}, nil, true},
		{"invalidHeader", NewDecoder(bytes.NewBufferString("aaaaaaaa")), args{&Buffer{ByteLength: 1}}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.decodeBinaryBuffer(tt.args.buffer); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.decodeBinaryBuffer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(tt.args.buffer.Data, tt.want) {
				t.Errorf("Decoder.decodeBinaryBuffer() buffer = %v, want %v", tt.args.buffer.Data, tt.want)
			}
		})
	}
}

func TestDecoder_Decode(t *testing.T) {
	type args struct {
		doc *Document
	}
	tests := []struct {
		name    string
		d       *Decoder
		args    args
		wantErr bool
	}{
		{"baseJSON", NewDecoderFS(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 1, \"URI\": \"a.bin\"}]}"), fstest.MapFS{"a.bin": &fstest.MapFile{Data: []byte("abcdfg")}}), args{new(Document)}, false},
		{"onlyGLBHeader", NewDecoderFS(bytes.NewBuffer([]byte{0x67, 0x6c, 0x54, 0x46, 0x02, 0x00, 0x00, 0x00, 0x40, 0x0b, 0x00, 0x00, 0x5c, 0x06, 0x00, 0x00, 0x4a, 0x53, 0x4f, 0x4e}), fstest.MapFS{"a.bin": &fstest.MapFile{Data: []byte("abcdfg")}}), args{new(Document)}, true},
		{"glbNoJSONChunk", NewDecoderFS(bytes.NewBuffer([]byte{0x67, 0x6c, 0x54, 0x46, 0x02, 0x00, 0x00, 0x00, 0x40, 0x0b, 0x00, 0x00, 0x5c, 0x06, 0x00, 0x00, 0x4a, 0x52, 0x4f, 0x4e}), fstest.MapFS{"a.bin": &fstest.MapFile{Data: []byte("abcdfg")}}), args{new(Document)}, true},
		{"empty", NewDecoder(bytes.NewBufferString("")), args{new(Document)}, true},
		{"invalidJSON", NewDecoder(bytes.NewBufferString("{asset: {}}")), args{new(Document)}, true},
		{"invalidBuffer", NewDecoder(bytes.NewBufferString("{\"buffers\": [{\"byteLength\": 0}]}")), args{new(Document)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.Decode(tt.args.doc); (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSampler_Decode(t *testing.T) {

	tests := []struct {
		name    string
		s       []byte
		want    *Sampler
		wantErr bool
	}{
		{"empty", []byte(`{}`), &Sampler{}, false},
		{"nondefault",
			[]byte(`{"minFilter":9728,"wrapT":33071}`),
			&Sampler{MagFilter: MagUndefined, MinFilter: MinNearest, WrapS: WrapRepeat, WrapT: WrapClampToEdge},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Sampler
			err := json.Unmarshal(tt.s, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshaling Sampler error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&got, tt.want) {
				t.Errorf("Unmarshaling Sampler = %v, want %v", string(tt.s), tt.want)
			}
		})
	}
}
