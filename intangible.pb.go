// Code generated by protoc-gen-go.
// source: intangible.proto
// DO NOT EDIT!

/*
Package Intangible is a generated protocol buffer package.

It is generated from these files:
	intangible.proto

It has these top-level messages:
	Vector
	Mesh
	Texture
	Rendering
	API
	Object
*/
package Intangible

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Vector struct {
	X float32 `protobuf:"fixed32,1,opt,name=x" json:"x,omitempty"`
	Y float32 `protobuf:"fixed32,2,opt,name=y" json:"y,omitempty"`
	Z float32 `protobuf:"fixed32,3,opt,name=z" json:"z,omitempty"`
}

func (m *Vector) Reset()                    { *m = Vector{} }
func (m *Vector) String() string            { return proto.CompactTextString(m) }
func (*Vector) ProtoMessage()               {}
func (*Vector) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Mesh struct {
	// source_uri can be any web-retrievable URL, or a primitive under the "platonic" scheme, eg "platonic:cube". "platonic:cube", notably, out of all the possibilities, is the only actual platonic solid that can be described here, but a passing classical reference is always worthwhile, even (or especially) when incorrect.
	SourceUri string  `protobuf:"bytes,1,opt,name=source_uri,json=sourceUri" json:"source_uri,omitempty"`
	Rescale   *Vector `protobuf:"bytes,2,opt,name=rescale" json:"rescale,omitempty"`
	Rotation  *Vector `protobuf:"bytes,3,opt,name=rotation" json:"rotation,omitempty"`
}

func (m *Mesh) Reset()                    { *m = Mesh{} }
func (m *Mesh) String() string            { return proto.CompactTextString(m) }
func (*Mesh) ProtoMessage()               {}
func (*Mesh) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Mesh) GetRescale() *Vector {
	if m != nil {
		return m.Rescale
	}
	return nil
}

func (m *Mesh) GetRotation() *Vector {
	if m != nil {
		return m.Rotation
	}
	return nil
}

type Texture struct {
	SourceUri string `protobuf:"bytes,1,opt,name=source_uri,json=sourceUri" json:"source_uri,omitempty"`
}

func (m *Texture) Reset()                    { *m = Texture{} }
func (m *Texture) String() string            { return proto.CompactTextString(m) }
func (*Texture) ProtoMessage()               {}
func (*Texture) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type Rendering struct {
	Mesh    *Mesh    `protobuf:"bytes,2,opt,name=mesh" json:"mesh,omitempty"`
	Texture *Texture `protobuf:"bytes,3,opt,name=texture" json:"texture,omitempty"`
}

func (m *Rendering) Reset()                    { *m = Rendering{} }
func (m *Rendering) String() string            { return proto.CompactTextString(m) }
func (*Rendering) ProtoMessage()               {}
func (*Rendering) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *Rendering) GetMesh() *Mesh {
	if m != nil {
		return m.Mesh
	}
	return nil
}

func (m *Rendering) GetTexture() *Texture {
	if m != nil {
		return m.Texture
	}
	return nil
}

type API struct {
	// if endpoint is ws://, treat as connection to another server.
	Endpoint string `protobuf:"bytes,1,opt,name=endpoint" json:"endpoint,omitempty"`
	// TODO(dichro): standardized API IDs, to readily identify operations like transfer to another room, moving self/other, pick up/drop objects, etc.
	Name        string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *API) Reset()                    { *m = API{} }
func (m *API) String() string            { return proto.CompactTextString(m) }
func (*API) ProtoMessage()               {}
func (*API) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type Object struct {
	Id          string     `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	Position    *Vector    `protobuf:"bytes,2,opt,name=position" json:"position,omitempty"`
	Rendering   *Rendering `protobuf:"bytes,3,opt,name=rendering" json:"rendering,omitempty"`
	Rotation    *Vector    `protobuf:"bytes,4,opt,name=rotation" json:"rotation,omitempty"`
	BoundingBox *Vector    `protobuf:"bytes,5,opt,name=bounding_box,json=boundingBox" json:"bounding_box,omitempty"`
	Removed     bool       `protobuf:"varint,6,opt,name=removed" json:"removed,omitempty"`
	Api         []*API     `protobuf:"bytes,7,rep,name=api" json:"api,omitempty"`
}

func (m *Object) Reset()                    { *m = Object{} }
func (m *Object) String() string            { return proto.CompactTextString(m) }
func (*Object) ProtoMessage()               {}
func (*Object) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *Object) GetPosition() *Vector {
	if m != nil {
		return m.Position
	}
	return nil
}

func (m *Object) GetRendering() *Rendering {
	if m != nil {
		return m.Rendering
	}
	return nil
}

func (m *Object) GetRotation() *Vector {
	if m != nil {
		return m.Rotation
	}
	return nil
}

func (m *Object) GetBoundingBox() *Vector {
	if m != nil {
		return m.BoundingBox
	}
	return nil
}

func (m *Object) GetApi() []*API {
	if m != nil {
		return m.Api
	}
	return nil
}

func init() {
	proto.RegisterType((*Vector)(nil), "Intangible.Vector")
	proto.RegisterType((*Mesh)(nil), "Intangible.Mesh")
	proto.RegisterType((*Texture)(nil), "Intangible.Texture")
	proto.RegisterType((*Rendering)(nil), "Intangible.Rendering")
	proto.RegisterType((*API)(nil), "Intangible.API")
	proto.RegisterType((*Object)(nil), "Intangible.Object")
}

func init() { proto.RegisterFile("intangible.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 376 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x52, 0x4d, 0x4b, 0xc3, 0x40,
	0x10, 0x25, 0x1f, 0x26, 0xcd, 0xb4, 0x68, 0x19, 0x11, 0x82, 0x20, 0xd4, 0xe0, 0xa1, 0x07, 0xcd,
	0xa1, 0xc5, 0x1f, 0xa0, 0xb7, 0x1e, 0x44, 0x59, 0xfc, 0x38, 0xd6, 0xa4, 0x59, 0xea, 0x4a, 0x9b,
	0x0d, 0x9b, 0x44, 0xd2, 0x5e, 0xfd, 0x25, 0xfe, 0x53, 0x37, 0x9b, 0x8f, 0x46, 0xb0, 0xf4, 0xb6,
	0x33, 0xfb, 0xde, 0xbc, 0x37, 0x6f, 0x17, 0x86, 0x2c, 0xce, 0x82, 0x78, 0xc9, 0xc2, 0x15, 0xf5,
	0x13, 0xc1, 0x33, 0x8e, 0x30, 0x6b, 0x3b, 0xde, 0x04, 0xac, 0x57, 0xba, 0xc8, 0xb8, 0xc0, 0x01,
	0x68, 0x85, 0xab, 0x8d, 0xb4, 0xb1, 0x4e, 0xb4, 0xa2, 0xac, 0x36, 0xae, 0x5e, 0x55, 0x9b, 0xb2,
	0xda, 0xba, 0x46, 0x55, 0x6d, 0xbd, 0x6f, 0x0d, 0xcc, 0x07, 0x9a, 0x7e, 0xe0, 0x05, 0x40, 0xca,
	0x73, 0xb1, 0xa0, 0xf3, 0x5c, 0x30, 0xc5, 0x75, 0x88, 0x53, 0x75, 0x5e, 0x04, 0xc3, 0x6b, 0xb0,
	0x05, 0x4d, 0x17, 0xc1, 0x8a, 0xaa, 0x49, 0xfd, 0x09, 0xfa, 0x3b, 0x65, 0xbf, 0x92, 0x25, 0x0d,
	0x04, 0x7d, 0xe8, 0x49, 0x77, 0x41, 0xc6, 0x78, 0xac, 0xa4, 0xfe, 0x87, 0xb7, 0x18, 0x6f, 0x0c,
	0xf6, 0x33, 0x2d, 0xb2, 0x5c, 0xd0, 0x03, 0x3e, 0xbc, 0x77, 0x70, 0x08, 0x8d, 0x23, 0x2a, 0x58,
	0xbc, 0xc4, 0x2b, 0x30, 0xd7, 0xd2, 0x7b, 0xed, 0x68, 0xd8, 0x95, 0x28, 0x77, 0x22, 0xea, 0x16,
	0x6f, 0xc0, 0xce, 0xaa, 0xe1, 0xb5, 0x97, 0xd3, 0x2e, 0xb0, 0xd6, 0x25, 0x0d, 0xc6, 0x7b, 0x03,
	0xe3, 0xee, 0x69, 0x86, 0xe7, 0xd0, 0x93, 0x3a, 0x09, 0x97, 0x89, 0xd7, 0x2e, 0xda, 0x1a, 0x11,
	0xcc, 0x38, 0x58, 0x57, 0x49, 0x38, 0x44, 0x9d, 0x71, 0x04, 0xfd, 0x48, 0x6e, 0x2f, 0x58, 0xd2,
	0x6e, 0xed, 0x90, 0x6e, 0xcb, 0xfb, 0xd1, 0xc1, 0x7a, 0x0c, 0x3f, 0xe5, 0xee, 0x78, 0x0c, 0x3a,
	0x8b, 0xea, 0xb1, 0xf2, 0x54, 0xe6, 0x95, 0xf0, 0x94, 0x29, 0xe6, 0xfe, 0x78, 0x5b, 0x0c, 0x4e,
	0xc1, 0x11, 0x4d, 0x0a, 0xf5, 0x52, 0x67, 0x5d, 0x42, 0x1b, 0x11, 0xd9, 0xe1, 0xfe, 0x3c, 0x8a,
	0x79, 0xf8, 0x51, 0xf0, 0x16, 0x06, 0x21, 0xcf, 0xe3, 0x48, 0x72, 0xe7, 0x21, 0x2f, 0xdc, 0xa3,
	0xbd, 0x9c, 0x7e, 0x83, 0xbb, 0xe7, 0x05, 0xba, 0xe5, 0x4f, 0x59, 0xf3, 0x2f, 0x1a, 0xb9, 0x96,
	0x64, 0xf4, 0x48, 0x53, 0xe2, 0x25, 0x18, 0x41, 0xc2, 0x5c, 0x7b, 0x64, 0xc8, 0x39, 0x27, 0xdd,
	0x39, 0x32, 0x70, 0x52, 0xde, 0x85, 0x96, 0xfa, 0xd5, 0xd3, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x98, 0xfe, 0x2c, 0xb0, 0xe9, 0x02, 0x00, 0x00,
}
