// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.14.0
// source: workflow.proto

package massdriver

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type DeploymentStatus int32

const (
	DeploymentStatus_PENDING   DeploymentStatus = 0
	DeploymentStatus_RUNNING   DeploymentStatus = 1
	DeploymentStatus_COMPLETED DeploymentStatus = 2
	DeploymentStatus_FAILED    DeploymentStatus = 3
)

// Enum value maps for DeploymentStatus.
var (
	DeploymentStatus_name = map[int32]string{
		0: "PENDING",
		1: "RUNNING",
		2: "COMPLETED",
		3: "FAILED",
	}
	DeploymentStatus_value = map[string]int32{
		"PENDING":   0,
		"RUNNING":   1,
		"COMPLETED": 2,
		"FAILED":    3,
	}
)

func (x DeploymentStatus) Enum() *DeploymentStatus {
	p := new(DeploymentStatus)
	*p = x
	return p
}

func (x DeploymentStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DeploymentStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_workflow_proto_enumTypes[0].Descriptor()
}

func (DeploymentStatus) Type() protoreflect.EnumType {
	return &file_workflow_proto_enumTypes[0]
}

func (x DeploymentStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DeploymentStatus.Descriptor instead.
func (DeploymentStatus) EnumDescriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{0}
}

type StartDeploymentRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Token string `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
}

func (x *StartDeploymentRequest) Reset() {
	*x = StartDeploymentRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartDeploymentRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartDeploymentRequest) ProtoMessage() {}

func (x *StartDeploymentRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartDeploymentRequest.ProtoReflect.Descriptor instead.
func (*StartDeploymentRequest) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{0}
}

func (x *StartDeploymentRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *StartDeploymentRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ArtifactMetadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProviderResourceId string `protobuf:"bytes,1,opt,name=provider_resource_id,json=providerResourceId,proto3" json:"provider_resource_id,omitempty"`
	Type               string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Name               string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ArtifactMetadata) Reset() {
	*x = ArtifactMetadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArtifactMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArtifactMetadata) ProtoMessage() {}

func (x *ArtifactMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArtifactMetadata.ProtoReflect.Descriptor instead.
func (*ArtifactMetadata) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{1}
}

func (x *ArtifactMetadata) GetProviderResourceId() string {
	if x != nil {
		return x.ProviderResourceId
	}
	return ""
}

func (x *ArtifactMetadata) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *ArtifactMetadata) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Artifact struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Metadata *ArtifactMetadata `protobuf:"bytes,1,opt,name=metadata,proto3" json:"metadata,omitempty"`
	Secrets  *structpb.Struct  `protobuf:"bytes,4,opt,name=secrets,proto3" json:"secrets,omitempty"`
	Specs    *structpb.Struct  `protobuf:"bytes,5,opt,name=specs,proto3" json:"specs,omitempty"`
}

func (x *Artifact) Reset() {
	*x = Artifact{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Artifact) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Artifact) ProtoMessage() {}

func (x *Artifact) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Artifact.ProtoReflect.Descriptor instead.
func (*Artifact) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{2}
}

func (x *Artifact) GetMetadata() *ArtifactMetadata {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Artifact) GetSecrets() *structpb.Struct {
	if x != nil {
		return x.Secrets
	}
	return nil
}

func (x *Artifact) GetSpecs() *structpb.Struct {
	if x != nil {
		return x.Specs
	}
	return nil
}

type UploadArtifactsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DeploymentId    string      `protobuf:"bytes,1,opt,name=deployment_id,json=deploymentId,proto3" json:"deployment_id,omitempty"`
	DeploymentToken string      `protobuf:"bytes,2,opt,name=deployment_token,json=deploymentToken,proto3" json:"deployment_token,omitempty"`
	Artifacts       []*Artifact `protobuf:"bytes,3,rep,name=artifacts,proto3" json:"artifacts,omitempty"`
}

func (x *UploadArtifactsRequest) Reset() {
	*x = UploadArtifactsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UploadArtifactsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadArtifactsRequest) ProtoMessage() {}

func (x *UploadArtifactsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadArtifactsRequest.ProtoReflect.Descriptor instead.
func (*UploadArtifactsRequest) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{3}
}

func (x *UploadArtifactsRequest) GetDeploymentId() string {
	if x != nil {
		return x.DeploymentId
	}
	return ""
}

func (x *UploadArtifactsRequest) GetDeploymentToken() string {
	if x != nil {
		return x.DeploymentToken
	}
	return ""
}

func (x *UploadArtifactsRequest) GetArtifacts() []*Artifact {
	if x != nil {
		return x.Artifacts
	}
	return nil
}

type Deployment struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Status      DeploymentStatus       `protobuf:"varint,2,opt,name=status,proto3,enum=mdtwirp.DeploymentStatus" json:"status,omitempty"`
	Params      *structpb.Struct       `protobuf:"bytes,3,opt,name=params,proto3" json:"params,omitempty"`
	Connections *structpb.Struct       `protobuf:"bytes,4,opt,name=connections,proto3" json:"connections,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
}

func (x *Deployment) Reset() {
	*x = Deployment{}
	if protoimpl.UnsafeEnabled {
		mi := &file_workflow_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Deployment) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Deployment) ProtoMessage() {}

func (x *Deployment) ProtoReflect() protoreflect.Message {
	mi := &file_workflow_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Deployment.ProtoReflect.Descriptor instead.
func (*Deployment) Descriptor() ([]byte, []int) {
	return file_workflow_proto_rawDescGZIP(), []int{4}
}

func (x *Deployment) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Deployment) GetStatus() DeploymentStatus {
	if x != nil {
		return x.Status
	}
	return DeploymentStatus_PENDING
}

func (x *Deployment) GetParams() *structpb.Struct {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *Deployment) GetConnections() *structpb.Struct {
	if x != nil {
		return x.Connections
	}
	return nil
}

func (x *Deployment) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Deployment) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

var File_workflow_proto protoreflect.FileDescriptor

var file_workflow_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x77, 0x6f, 0x72, 0x6b, 0x66, 0x6c, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72, 0x70, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x16, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x6c, 0x0a, 0x10, 0x41, 0x72, 0x74, 0x69,
	0x66, 0x61, 0x63, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x30, 0x0a, 0x14,
	0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x49, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0xa3, 0x01, 0x0a, 0x08, 0x41, 0x72, 0x74, 0x69, 0x66,
	0x61, 0x63, 0x74, 0x12, 0x35, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72, 0x70, 0x2e,
	0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0x52, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x31, 0x0a, 0x07, 0x73, 0x65,
	0x63, 0x72, 0x65, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x75, 0x63, 0x74, 0x52, 0x07, 0x73, 0x65, 0x63, 0x72, 0x65, 0x74, 0x73, 0x12, 0x2d, 0x0a,
	0x05, 0x73, 0x70, 0x65, 0x63, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x05, 0x73, 0x70, 0x65, 0x63, 0x73, 0x22, 0x99, 0x01, 0x0a,
	0x16, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x10,
	0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x2f, 0x0a, 0x09, 0x61, 0x72, 0x74, 0x69, 0x66,
	0x61, 0x63, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x64, 0x74,
	0x77, 0x69, 0x72, 0x70, 0x2e, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x52, 0x09, 0x61,
	0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73, 0x22, 0xb1, 0x02, 0x0a, 0x0a, 0x44, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72,
	0x70, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2f, 0x0a, 0x06, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x39, 0x0a, 0x0b, 0x63,
	0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x5f, 0x61, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x39, 0x0a, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x52, 0x09, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x2a, 0x47, 0x0a, 0x10,
	0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45, 0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x00, 0x12, 0x0b, 0x0a,
	0x07, 0x52, 0x55, 0x4e, 0x4e, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x43, 0x4f,
	0x4d, 0x50, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x03, 0x32, 0xa0, 0x01, 0x0a, 0x08, 0x57, 0x6f, 0x72, 0x6b, 0x66, 0x6c,
	0x6f, 0x77, 0x12, 0x49, 0x0a, 0x0f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x44, 0x65, 0x70, 0x6c, 0x6f,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x1f, 0x2e, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72, 0x70, 0x2e,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13, 0x2e, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72, 0x70,
	0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x12, 0x49, 0x0a,
	0x0f, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73,
	0x12, 0x1f, 0x2e, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72, 0x70, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x41, 0x72, 0x74, 0x69, 0x66, 0x61, 0x63, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x13, 0x2e, 0x6d, 0x64, 0x74, 0x77, 0x69, 0x72, 0x70, 0x2e, 0x44, 0x65, 0x70, 0x6c,
	0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x22, 0x00, 0x42, 0x0e, 0x5a, 0x0c, 0x2e, 0x2f, 0x6d, 0x61,
	0x73, 0x73, 0x64, 0x72, 0x69, 0x76, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_workflow_proto_rawDescOnce sync.Once
	file_workflow_proto_rawDescData = file_workflow_proto_rawDesc
)

func file_workflow_proto_rawDescGZIP() []byte {
	file_workflow_proto_rawDescOnce.Do(func() {
		file_workflow_proto_rawDescData = protoimpl.X.CompressGZIP(file_workflow_proto_rawDescData)
	})
	return file_workflow_proto_rawDescData
}

var file_workflow_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_workflow_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_workflow_proto_goTypes = []interface{}{
	(DeploymentStatus)(0),          // 0: mdtwirp.DeploymentStatus
	(*StartDeploymentRequest)(nil), // 1: mdtwirp.StartDeploymentRequest
	(*ArtifactMetadata)(nil),       // 2: mdtwirp.ArtifactMetadata
	(*Artifact)(nil),               // 3: mdtwirp.Artifact
	(*UploadArtifactsRequest)(nil), // 4: mdtwirp.UploadArtifactsRequest
	(*Deployment)(nil),             // 5: mdtwirp.Deployment
	(*structpb.Struct)(nil),        // 6: google.protobuf.Struct
	(*timestamppb.Timestamp)(nil),  // 7: google.protobuf.Timestamp
}
var file_workflow_proto_depIdxs = []int32{
	2,  // 0: mdtwirp.Artifact.metadata:type_name -> mdtwirp.ArtifactMetadata
	6,  // 1: mdtwirp.Artifact.secrets:type_name -> google.protobuf.Struct
	6,  // 2: mdtwirp.Artifact.specs:type_name -> google.protobuf.Struct
	3,  // 3: mdtwirp.UploadArtifactsRequest.artifacts:type_name -> mdtwirp.Artifact
	0,  // 4: mdtwirp.Deployment.status:type_name -> mdtwirp.DeploymentStatus
	6,  // 5: mdtwirp.Deployment.params:type_name -> google.protobuf.Struct
	6,  // 6: mdtwirp.Deployment.connections:type_name -> google.protobuf.Struct
	7,  // 7: mdtwirp.Deployment.created_at:type_name -> google.protobuf.Timestamp
	7,  // 8: mdtwirp.Deployment.updated_at:type_name -> google.protobuf.Timestamp
	1,  // 9: mdtwirp.Workflow.StartDeployment:input_type -> mdtwirp.StartDeploymentRequest
	4,  // 10: mdtwirp.Workflow.UploadArtifacts:input_type -> mdtwirp.UploadArtifactsRequest
	5,  // 11: mdtwirp.Workflow.StartDeployment:output_type -> mdtwirp.Deployment
	5,  // 12: mdtwirp.Workflow.UploadArtifacts:output_type -> mdtwirp.Deployment
	11, // [11:13] is the sub-list for method output_type
	9,  // [9:11] is the sub-list for method input_type
	9,  // [9:9] is the sub-list for extension type_name
	9,  // [9:9] is the sub-list for extension extendee
	0,  // [0:9] is the sub-list for field type_name
}

func init() { file_workflow_proto_init() }
func file_workflow_proto_init() {
	if File_workflow_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_workflow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartDeploymentRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArtifactMetadata); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflow_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Artifact); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflow_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UploadArtifactsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_workflow_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Deployment); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_workflow_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_workflow_proto_goTypes,
		DependencyIndexes: file_workflow_proto_depIdxs,
		EnumInfos:         file_workflow_proto_enumTypes,
		MessageInfos:      file_workflow_proto_msgTypes,
	}.Build()
	File_workflow_proto = out.File
	file_workflow_proto_rawDesc = nil
	file_workflow_proto_goTypes = nil
	file_workflow_proto_depIdxs = nil
}
