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
	// TODO(dichro): schema for request+reply. Avro? There appears to be no easily accessible protobuf, thrift, or json-schema IDL parser that'll work with Unity.
	AvroRequestSchema string `protobuf:"bytes,4,opt,name=avro_request_schema,json=avroRequestSchema" json:"avro_request_schema,omitempty"`
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
	// 408 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x84, 0x53, 0x4b, 0xcf, 0xd3, 0x30,
	0x10, 0x54, 0x1e, 0x24, 0xcd, 0xf6, 0x13, 0x94, 0xad, 0x90, 0x2c, 0x24, 0xa4, 0x12, 0x71, 0xe8,
	0x01, 0x72, 0x68, 0xc5, 0x0f, 0x80, 0x5b, 0x0f, 0x08, 0x64, 0x1e, 0xd7, 0x90, 0x87, 0xd5, 0x1a,
	0x35, 0x71, 0x70, 0x92, 0x2a, 0xed, 0x95, 0x03, 0xbf, 0x83, 0x7f, 0x8a, 0xe3, 0x3c, 0x1a, 0x24,
	0xaa, 0xde, 0x3c, 0xde, 0xd9, 0x9d, 0xf1, 0x6c, 0x02, 0x0b, 0x9e, 0x57, 0x51, 0xbe, 0xe7, 0xf1,
	0x91, 0x05, 0x85, 0x14, 0x95, 0x40, 0xd8, 0x8d, 0x37, 0xfe, 0x06, 0x9c, 0x6f, 0x2c, 0xa9, 0x84,
	0xc4, 0x07, 0x30, 0x1a, 0x62, 0xac, 0x8c, 0xb5, 0x49, 0x8d, 0xa6, 0x45, 0x67, 0x62, 0x76, 0xe8,
	0xdc, 0xa2, 0x0b, 0xb1, 0x3a, 0x74, 0xf1, 0x7f, 0x19, 0x60, 0x7f, 0x60, 0xe5, 0x01, 0x5f, 0x00,
	0x94, 0xa2, 0x96, 0x09, 0x0b, 0x6b, 0xc9, 0x75, 0xaf, 0x47, 0xbd, 0xee, 0xe6, 0xab, 0xe4, 0xf8,
	0x1a, 0x5c, 0xc9, 0xca, 0x24, 0x3a, 0x32, 0x3d, 0x69, 0xbe, 0xc1, 0xe0, 0xaa, 0x1c, 0x74, 0xb2,
	0x74, 0xa0, 0x60, 0x00, 0x33, 0xe5, 0x2e, 0xaa, 0xb8, 0xc8, 0xb5, 0xd4, 0xff, 0xe9, 0x23, 0xc7,
	0x5f, 0x83, 0xfb, 0x85, 0x35, 0x55, 0x2d, 0xd9, 0x1d, 0x1f, 0xfe, 0x77, 0xf0, 0x28, 0xcb, 0x53,
	0x26, 0x79, 0xbe, 0xc7, 0x57, 0x60, 0x67, 0xca, 0x7b, 0xef, 0x68, 0x31, 0x95, 0x68, 0xdf, 0x44,
	0x75, 0x15, 0xdf, 0x80, 0x5b, 0x75, 0xc3, 0x7b, 0x2f, 0xcb, 0x29, 0xb1, 0xd7, 0xa5, 0x03, 0xc7,
	0xff, 0x6d, 0x80, 0xf5, 0xee, 0xd3, 0x0e, 0x9f, 0xc3, 0x4c, 0x09, 0x15, 0x42, 0x45, 0xde, 0xdb,
	0x18, 0x31, 0x22, 0xd8, 0x79, 0x94, 0x75, 0x51, 0x78, 0x54, 0x9f, 0x71, 0x05, 0xf3, 0x54, 0x3d,
	0x5f, 0xf2, 0x62, 0x7c, 0xb6, 0x47, 0xa7, 0x57, 0x2a, 0x95, 0x65, 0x74, 0x92, 0x22, 0x94, 0xec,
	0x67, 0xcd, 0xca, 0x2a, 0x2c, 0x93, 0x03, 0xcb, 0x22, 0x62, 0x6b, 0xe6, 0xd3, 0xb6, 0x44, 0xbb,
	0xca, 0x67, 0x5d, 0xf0, 0xff, 0x98, 0xe0, 0x7c, 0x8c, 0x7f, 0xa8, 0xb0, 0xf0, 0x31, 0x98, 0x3c,
	0xed, 0x6d, 0xa8, 0x53, 0x1b, 0x70, 0x21, 0x4a, 0xae, 0x95, 0x6e, 0xef, 0x63, 0xe4, 0xe0, 0x16,
	0x3c, 0x39, 0xc4, 0xd6, 0xa7, 0xf0, 0x6c, 0xda, 0x30, 0x66, 0x4a, 0xaf, 0xbc, 0x7f, 0xb6, 0x68,
	0xdf, 0xdf, 0x22, 0xbe, 0x85, 0x87, 0x58, 0xd4, 0x79, 0xaa, 0x7a, 0xc3, 0x58, 0x34, 0xe4, 0xd1,
	0xcd, 0x9e, 0xf9, 0xc0, 0x7b, 0x2f, 0x1a, 0x24, 0xed, 0xa7, 0x95, 0x89, 0x13, 0x4b, 0x89, 0xa3,
	0x3a, 0x66, 0x74, 0x80, 0xf8, 0x12, 0xac, 0xa8, 0xe0, 0xc4, 0x5d, 0x59, 0x6a, 0xce, 0x93, 0xe9,
	0x1c, 0xb5, 0x20, 0xda, 0xd6, 0x62, 0x47, 0xff, 0x06, 0xdb, 0xbf, 0x01, 0x00, 0x00, 0xff, 0xff,
	0xa4, 0x5f, 0xf9, 0x30, 0x1a, 0x03, 0x00, 0x00,
}
