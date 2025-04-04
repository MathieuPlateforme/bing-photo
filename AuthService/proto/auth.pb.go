// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.30.2
// source: proto/auth.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request and response messages
type LoginRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginRequest) Reset() {
	*x = LoginRequest{}
	mi := &file_proto_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRequest) ProtoMessage() {}

func (x *LoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRequest.ProtoReflect.Descriptor instead.
func (*LoginRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{0}
}

func (x *LoginRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *LoginRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type LoginResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LoginResponse) Reset() {
	*x = LoginResponse{}
	mi := &file_proto_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LoginResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginResponse) ProtoMessage() {}

func (x *LoginResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginResponse.ProtoReflect.Descriptor instead.
func (*LoginResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{1}
}

func (x *LoginResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *LoginResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type RegisterRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	Username      string                 `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	FirstName     string                 `protobuf:"bytes,4,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName      string                 `protobuf:"bytes,5,opt,name=lastName,proto3" json:"lastName,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterRequest) Reset() {
	*x = RegisterRequest{}
	mi := &file_proto_auth_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRequest) ProtoMessage() {}

func (x *RegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRequest.ProtoReflect.Descriptor instead.
func (*RegisterRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{2}
}

func (x *RegisterRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *RegisterRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *RegisterRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *RegisterRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

type RegisterResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RegisterResponse) Reset() {
	*x = RegisterResponse{}
	mi := &file_proto_auth_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterResponse) ProtoMessage() {}

func (x *RegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterResponse.ProtoReflect.Descriptor instead.
func (*RegisterResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ValidateTokenRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidateTokenRequest) Reset() {
	*x = ValidateTokenRequest{}
	mi := &file_proto_auth_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenRequest) ProtoMessage() {}

func (x *ValidateTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenRequest.ProtoReflect.Descriptor instead.
func (*ValidateTokenRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{4}
}

func (x *ValidateTokenRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type ValidateTokenResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Valid         bool                   `protobuf:"varint,1,opt,name=valid,proto3" json:"valid,omitempty"`
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ValidateTokenResponse) Reset() {
	*x = ValidateTokenResponse{}
	mi := &file_proto_auth_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ValidateTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ValidateTokenResponse) ProtoMessage() {}

func (x *ValidateTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ValidateTokenResponse.ProtoReflect.Descriptor instead.
func (*ValidateTokenResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{5}
}

func (x *ValidateTokenResponse) GetValid() bool {
	if x != nil {
		return x.Valid
	}
	return false
}

func (x *ValidateTokenResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ForgotPasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ForgotPasswordRequest) Reset() {
	*x = ForgotPasswordRequest{}
	mi := &file_proto_auth_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ForgotPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForgotPasswordRequest) ProtoMessage() {}

func (x *ForgotPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForgotPasswordRequest.ProtoReflect.Descriptor instead.
func (*ForgotPasswordRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{6}
}

func (x *ForgotPasswordRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type ForgotPasswordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ForgotPasswordResponse) Reset() {
	*x = ForgotPasswordResponse{}
	mi := &file_proto_auth_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ForgotPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ForgotPasswordResponse) ProtoMessage() {}

func (x *ForgotPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ForgotPasswordResponse.ProtoReflect.Descriptor instead.
func (*ForgotPasswordResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{7}
}

func (x *ForgotPasswordResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type ResetPasswordRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Token         string                 `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	NewPassword   string                 `protobuf:"bytes,3,opt,name=new_password,json=newPassword,proto3" json:"new_password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResetPasswordRequest) Reset() {
	*x = ResetPasswordRequest{}
	mi := &file_proto_auth_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetPasswordRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetPasswordRequest) ProtoMessage() {}

func (x *ResetPasswordRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetPasswordRequest.ProtoReflect.Descriptor instead.
func (*ResetPasswordRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{8}
}

func (x *ResetPasswordRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *ResetPasswordRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

func (x *ResetPasswordRequest) GetNewPassword() string {
	if x != nil {
		return x.NewPassword
	}
	return ""
}

type ResetPasswordResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ResetPasswordResponse) Reset() {
	*x = ResetPasswordResponse{}
	mi := &file_proto_auth_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ResetPasswordResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResetPasswordResponse) ProtoMessage() {}

func (x *ResetPasswordResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResetPasswordResponse.ProtoReflect.Descriptor instead.
func (*ResetPasswordResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{9}
}

func (x *ResetPasswordResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GoogleAuthRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoogleAuthRequest) Reset() {
	*x = GoogleAuthRequest{}
	mi := &file_proto_auth_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoogleAuthRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoogleAuthRequest) ProtoMessage() {}

func (x *GoogleAuthRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoogleAuthRequest.ProtoReflect.Descriptor instead.
func (*GoogleAuthRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{10}
}

type GoogleAuthResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	AuthUrl       string                 `protobuf:"bytes,1,opt,name=authUrl,proto3" json:"authUrl,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoogleAuthResponse) Reset() {
	*x = GoogleAuthResponse{}
	mi := &file_proto_auth_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoogleAuthResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoogleAuthResponse) ProtoMessage() {}

func (x *GoogleAuthResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoogleAuthResponse.ProtoReflect.Descriptor instead.
func (*GoogleAuthResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{11}
}

func (x *GoogleAuthResponse) GetAuthUrl() string {
	if x != nil {
		return x.AuthUrl
	}
	return ""
}

type GoogleAuthCallbackRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoogleAuthCallbackRequest) Reset() {
	*x = GoogleAuthCallbackRequest{}
	mi := &file_proto_auth_proto_msgTypes[12]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoogleAuthCallbackRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoogleAuthCallbackRequest) ProtoMessage() {}

func (x *GoogleAuthCallbackRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[12]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoogleAuthCallbackRequest.ProtoReflect.Descriptor instead.
func (*GoogleAuthCallbackRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{12}
}

func (x *GoogleAuthCallbackRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

type GoogleAuthCallbackResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	UserInfo      string                 `protobuf:"bytes,2,opt,name=userInfo,proto3" json:"userInfo,omitempty"`
	Token         string                 `protobuf:"bytes,3,opt,name=Token,proto3" json:"Token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoogleAuthCallbackResponse) Reset() {
	*x = GoogleAuthCallbackResponse{}
	mi := &file_proto_auth_proto_msgTypes[13]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoogleAuthCallbackResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoogleAuthCallbackResponse) ProtoMessage() {}

func (x *GoogleAuthCallbackResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[13]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoogleAuthCallbackResponse.ProtoReflect.Descriptor instead.
func (*GoogleAuthCallbackResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{13}
}

func (x *GoogleAuthCallbackResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *GoogleAuthCallbackResponse) GetUserInfo() string {
	if x != nil {
		return x.UserInfo
	}
	return ""
}

func (x *GoogleAuthCallbackResponse) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type LogoutRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogoutRequest) Reset() {
	*x = LogoutRequest{}
	mi := &file_proto_auth_proto_msgTypes[14]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogoutRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogoutRequest) ProtoMessage() {}

func (x *LogoutRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[14]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogoutRequest.ProtoReflect.Descriptor instead.
func (*LogoutRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{14}
}

func (x *LogoutRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type LogoutResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogoutResponse) Reset() {
	*x = LogoutResponse{}
	mi := &file_proto_auth_proto_msgTypes[15]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogoutResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogoutResponse) ProtoMessage() {}

func (x *LogoutResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[15]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogoutResponse.ProtoReflect.Descriptor instead.
func (*LogoutResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{15}
}

func (x *LogoutResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type UpdateUserRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	FirstName     string                 `protobuf:"bytes,3,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName      string                 `protobuf:"bytes,4,opt,name=lastName,proto3" json:"lastName,omitempty"`
	Picture       string                 `protobuf:"bytes,5,opt,name=picture,proto3" json:"picture,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserRequest) Reset() {
	*x = UpdateUserRequest{}
	mi := &file_proto_auth_proto_msgTypes[16]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRequest) ProtoMessage() {}

func (x *UpdateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[16]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{16}
}

func (x *UpdateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UpdateUserRequest) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *UpdateUserRequest) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *UpdateUserRequest) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *UpdateUserRequest) GetPicture() string {
	if x != nil {
		return x.Picture
	}
	return ""
}

type UpdateUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Message       string                 `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserResponse) Reset() {
	*x = UpdateUserResponse{}
	mi := &file_proto_auth_proto_msgTypes[17]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserResponse) ProtoMessage() {}

func (x *UpdateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[17]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserResponse.ProtoReflect.Descriptor instead.
func (*UpdateUserResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{17}
}

func (x *UpdateUserResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetMeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMeRequest) Reset() {
	*x = GetMeRequest{}
	mi := &file_proto_auth_proto_msgTypes[18]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeRequest) ProtoMessage() {}

func (x *GetMeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[18]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeRequest.ProtoReflect.Descriptor instead.
func (*GetMeRequest) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{18}
}

func (x *GetMeRequest) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type GetMeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Username      string                 `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
	FirstName     string                 `protobuf:"bytes,3,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName      string                 `protobuf:"bytes,4,opt,name=lastName,proto3" json:"lastName,omitempty"`
	Picture       string                 `protobuf:"bytes,5,opt,name=picture,proto3" json:"picture,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetMeResponse) Reset() {
	*x = GetMeResponse{}
	mi := &file_proto_auth_proto_msgTypes[19]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetMeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMeResponse) ProtoMessage() {}

func (x *GetMeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_auth_proto_msgTypes[19]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMeResponse.ProtoReflect.Descriptor instead.
func (*GetMeResponse) Descriptor() ([]byte, []int) {
	return file_proto_auth_proto_rawDescGZIP(), []int{19}
}

func (x *GetMeResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *GetMeResponse) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *GetMeResponse) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *GetMeResponse) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *GetMeResponse) GetPicture() string {
	if x != nil {
		return x.Picture
	}
	return ""
}

var File_proto_auth_proto protoreflect.FileDescriptor

const file_proto_auth_proto_rawDesc = "" +
	"\n" +
	"\x10proto/auth.proto\x12\x05proto\"@\n" +
	"\fLoginRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\"?\n" +
	"\rLoginResponse\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"\x99\x01\n" +
	"\x0fRegisterRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\bpassword\x18\x02 \x01(\tR\bpassword\x12\x1a\n" +
	"\busername\x18\x03 \x01(\tR\busername\x12\x1c\n" +
	"\tfirstName\x18\x04 \x01(\tR\tfirstName\x12\x1a\n" +
	"\blastName\x18\x05 \x01(\tR\blastName\",\n" +
	"\x10RegisterResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\",\n" +
	"\x14ValidateTokenRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\"G\n" +
	"\x15ValidateTokenResponse\x12\x14\n" +
	"\x05valid\x18\x01 \x01(\bR\x05valid\x12\x18\n" +
	"\amessage\x18\x02 \x01(\tR\amessage\"-\n" +
	"\x15ForgotPasswordRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\"2\n" +
	"\x16ForgotPasswordResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"e\n" +
	"\x14ResetPasswordRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x14\n" +
	"\x05token\x18\x02 \x01(\tR\x05token\x12!\n" +
	"\fnew_password\x18\x03 \x01(\tR\vnewPassword\"1\n" +
	"\x15ResetPasswordResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"\x13\n" +
	"\x11GoogleAuthRequest\".\n" +
	"\x12GoogleAuthResponse\x12\x18\n" +
	"\aauthUrl\x18\x01 \x01(\tR\aauthUrl\"/\n" +
	"\x19GoogleAuthCallbackRequest\x12\x12\n" +
	"\x04code\x18\x01 \x01(\tR\x04code\"h\n" +
	"\x1aGoogleAuthCallbackResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\x12\x1a\n" +
	"\buserInfo\x18\x02 \x01(\tR\buserInfo\x12\x14\n" +
	"\x05Token\x18\x03 \x01(\tR\x05Token\"%\n" +
	"\rLogoutRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\"*\n" +
	"\x0eLogoutResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"\x99\x01\n" +
	"\x11UpdateUserRequest\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\busername\x18\x02 \x01(\tR\busername\x12\x1c\n" +
	"\tfirstName\x18\x03 \x01(\tR\tfirstName\x12\x1a\n" +
	"\blastName\x18\x04 \x01(\tR\blastName\x12\x18\n" +
	"\apicture\x18\x05 \x01(\tR\apicture\".\n" +
	"\x12UpdateUserResponse\x12\x18\n" +
	"\amessage\x18\x01 \x01(\tR\amessage\"$\n" +
	"\fGetMeRequest\x12\x14\n" +
	"\x05token\x18\x01 \x01(\tR\x05token\"\x95\x01\n" +
	"\rGetMeResponse\x12\x14\n" +
	"\x05email\x18\x01 \x01(\tR\x05email\x12\x1a\n" +
	"\busername\x18\x02 \x01(\tR\busername\x12\x1c\n" +
	"\tfirstName\x18\x03 \x01(\tR\tfirstName\x12\x1a\n" +
	"\blastName\x18\x04 \x01(\tR\blastName\x12\x18\n" +
	"\apicture\x18\x05 \x01(\tR\apicture2\xb6\x05\n" +
	"\vAuthService\x122\n" +
	"\x05Login\x12\x13.proto.LoginRequest\x1a\x14.proto.LoginResponse\x12;\n" +
	"\bRegister\x12\x16.proto.RegisterRequest\x1a\x17.proto.RegisterResponse\x12J\n" +
	"\rValidateToken\x12\x1b.proto.ValidateTokenRequest\x1a\x1c.proto.ValidateTokenResponse\x12M\n" +
	"\x0eForgotPassword\x12\x1c.proto.ForgotPasswordRequest\x1a\x1d.proto.ForgotPasswordResponse\x12J\n" +
	"\rResetPassword\x12\x1b.proto.ResetPasswordRequest\x1a\x1c.proto.ResetPasswordResponse\x12F\n" +
	"\x0fLoginWithGoogle\x12\x18.proto.GoogleAuthRequest\x1a\x19.proto.GoogleAuthResponse\x12Y\n" +
	"\x12GoogleAuthCallback\x12 .proto.GoogleAuthCallbackRequest\x1a!.proto.GoogleAuthCallbackResponse\x125\n" +
	"\x06Logout\x12\x14.proto.LogoutRequest\x1a\x15.proto.LogoutResponse\x12A\n" +
	"\n" +
	"UpdateUser\x12\x18.proto.UpdateUserRequest\x1a\x19.proto.UpdateUserResponse\x122\n" +
	"\x05GetMe\x12\x13.proto.GetMeRequest\x1a\x14.proto.GetMeResponseB\x0eZ\f/proto;protob\x06proto3"

var (
	file_proto_auth_proto_rawDescOnce sync.Once
	file_proto_auth_proto_rawDescData []byte
)

func file_proto_auth_proto_rawDescGZIP() []byte {
	file_proto_auth_proto_rawDescOnce.Do(func() {
		file_proto_auth_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_auth_proto_rawDesc), len(file_proto_auth_proto_rawDesc)))
	})
	return file_proto_auth_proto_rawDescData
}

var file_proto_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 20)
var file_proto_auth_proto_goTypes = []any{
	(*LoginRequest)(nil),               // 0: proto.LoginRequest
	(*LoginResponse)(nil),              // 1: proto.LoginResponse
	(*RegisterRequest)(nil),            // 2: proto.RegisterRequest
	(*RegisterResponse)(nil),           // 3: proto.RegisterResponse
	(*ValidateTokenRequest)(nil),       // 4: proto.ValidateTokenRequest
	(*ValidateTokenResponse)(nil),      // 5: proto.ValidateTokenResponse
	(*ForgotPasswordRequest)(nil),      // 6: proto.ForgotPasswordRequest
	(*ForgotPasswordResponse)(nil),     // 7: proto.ForgotPasswordResponse
	(*ResetPasswordRequest)(nil),       // 8: proto.ResetPasswordRequest
	(*ResetPasswordResponse)(nil),      // 9: proto.ResetPasswordResponse
	(*GoogleAuthRequest)(nil),          // 10: proto.GoogleAuthRequest
	(*GoogleAuthResponse)(nil),         // 11: proto.GoogleAuthResponse
	(*GoogleAuthCallbackRequest)(nil),  // 12: proto.GoogleAuthCallbackRequest
	(*GoogleAuthCallbackResponse)(nil), // 13: proto.GoogleAuthCallbackResponse
	(*LogoutRequest)(nil),              // 14: proto.LogoutRequest
	(*LogoutResponse)(nil),             // 15: proto.LogoutResponse
	(*UpdateUserRequest)(nil),          // 16: proto.UpdateUserRequest
	(*UpdateUserResponse)(nil),         // 17: proto.UpdateUserResponse
	(*GetMeRequest)(nil),               // 18: proto.GetMeRequest
	(*GetMeResponse)(nil),              // 19: proto.GetMeResponse
}
var file_proto_auth_proto_depIdxs = []int32{
	0,  // 0: proto.AuthService.Login:input_type -> proto.LoginRequest
	2,  // 1: proto.AuthService.Register:input_type -> proto.RegisterRequest
	4,  // 2: proto.AuthService.ValidateToken:input_type -> proto.ValidateTokenRequest
	6,  // 3: proto.AuthService.ForgotPassword:input_type -> proto.ForgotPasswordRequest
	8,  // 4: proto.AuthService.ResetPassword:input_type -> proto.ResetPasswordRequest
	10, // 5: proto.AuthService.LoginWithGoogle:input_type -> proto.GoogleAuthRequest
	12, // 6: proto.AuthService.GoogleAuthCallback:input_type -> proto.GoogleAuthCallbackRequest
	14, // 7: proto.AuthService.Logout:input_type -> proto.LogoutRequest
	16, // 8: proto.AuthService.UpdateUser:input_type -> proto.UpdateUserRequest
	18, // 9: proto.AuthService.GetMe:input_type -> proto.GetMeRequest
	1,  // 10: proto.AuthService.Login:output_type -> proto.LoginResponse
	3,  // 11: proto.AuthService.Register:output_type -> proto.RegisterResponse
	5,  // 12: proto.AuthService.ValidateToken:output_type -> proto.ValidateTokenResponse
	7,  // 13: proto.AuthService.ForgotPassword:output_type -> proto.ForgotPasswordResponse
	9,  // 14: proto.AuthService.ResetPassword:output_type -> proto.ResetPasswordResponse
	11, // 15: proto.AuthService.LoginWithGoogle:output_type -> proto.GoogleAuthResponse
	13, // 16: proto.AuthService.GoogleAuthCallback:output_type -> proto.GoogleAuthCallbackResponse
	15, // 17: proto.AuthService.Logout:output_type -> proto.LogoutResponse
	17, // 18: proto.AuthService.UpdateUser:output_type -> proto.UpdateUserResponse
	19, // 19: proto.AuthService.GetMe:output_type -> proto.GetMeResponse
	10, // [10:20] is the sub-list for method output_type
	0,  // [0:10] is the sub-list for method input_type
	0,  // [0:0] is the sub-list for extension type_name
	0,  // [0:0] is the sub-list for extension extendee
	0,  // [0:0] is the sub-list for field type_name
}

func init() { file_proto_auth_proto_init() }
func file_proto_auth_proto_init() {
	if File_proto_auth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_auth_proto_rawDesc), len(file_proto_auth_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   20,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_auth_proto_goTypes,
		DependencyIndexes: file_proto_auth_proto_depIdxs,
		MessageInfos:      file_proto_auth_proto_msgTypes,
	}.Build()
	File_proto_auth_proto = out.File
	file_proto_auth_proto_goTypes = nil
	file_proto_auth_proto_depIdxs = nil
}
