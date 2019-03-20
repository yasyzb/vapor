// Code generated by protoc-gen-go. DO NOT EDIT.
// source: bc.proto

/*
Package bc is a generated protocol buffer package.

It is generated from these files:
	bc.proto

It has these top-level messages:
	Hash
	Program
	AssetID
	AssetAmount
	AssetDefinition
	ValueSource
	ValueDestination
	Proof
	BytomBlockHeader
	BlockHeader
	TxHeader
	TxVerifyResult
	TransactionStatus
	Mux
	Coinbase
	Output
	Retirement
	Issuance
	Spend
	Claim
	Dpos
*/
package bc

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

type Hash struct {
	V0 uint64 `protobuf:"fixed64,1,opt,name=v0" json:"v0,omitempty"`
	V1 uint64 `protobuf:"fixed64,2,opt,name=v1" json:"v1,omitempty"`
	V2 uint64 `protobuf:"fixed64,3,opt,name=v2" json:"v2,omitempty"`
	V3 uint64 `protobuf:"fixed64,4,opt,name=v3" json:"v3,omitempty"`
}

func (m *Hash) Reset()                    { *m = Hash{} }
func (m *Hash) String() string            { return proto.CompactTextString(m) }
func (*Hash) ProtoMessage()               {}
func (*Hash) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Hash) GetV0() uint64 {
	if m != nil {
		return m.V0
	}
	return 0
}

func (m *Hash) GetV1() uint64 {
	if m != nil {
		return m.V1
	}
	return 0
}

func (m *Hash) GetV2() uint64 {
	if m != nil {
		return m.V2
	}
	return 0
}

func (m *Hash) GetV3() uint64 {
	if m != nil {
		return m.V3
	}
	return 0
}

type Program struct {
	VmVersion uint64 `protobuf:"varint,1,opt,name=vm_version,json=vmVersion" json:"vm_version,omitempty"`
	Code      []byte `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
}

func (m *Program) Reset()                    { *m = Program{} }
func (m *Program) String() string            { return proto.CompactTextString(m) }
func (*Program) ProtoMessage()               {}
func (*Program) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Program) GetVmVersion() uint64 {
	if m != nil {
		return m.VmVersion
	}
	return 0
}

func (m *Program) GetCode() []byte {
	if m != nil {
		return m.Code
	}
	return nil
}

// This message type duplicates Hash, above. One alternative is to
// embed a Hash inside an AssetID. But it's useful for AssetID to be
// plain old data (without pointers). Another alternative is use Hash
// in any protobuf types where an AssetID is called for, but it's
// preferable to have type safety.
type AssetID struct {
	V0 uint64 `protobuf:"fixed64,1,opt,name=v0" json:"v0,omitempty"`
	V1 uint64 `protobuf:"fixed64,2,opt,name=v1" json:"v1,omitempty"`
	V2 uint64 `protobuf:"fixed64,3,opt,name=v2" json:"v2,omitempty"`
	V3 uint64 `protobuf:"fixed64,4,opt,name=v3" json:"v3,omitempty"`
}

func (m *AssetID) Reset()                    { *m = AssetID{} }
func (m *AssetID) String() string            { return proto.CompactTextString(m) }
func (*AssetID) ProtoMessage()               {}
func (*AssetID) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *AssetID) GetV0() uint64 {
	if m != nil {
		return m.V0
	}
	return 0
}

func (m *AssetID) GetV1() uint64 {
	if m != nil {
		return m.V1
	}
	return 0
}

func (m *AssetID) GetV2() uint64 {
	if m != nil {
		return m.V2
	}
	return 0
}

func (m *AssetID) GetV3() uint64 {
	if m != nil {
		return m.V3
	}
	return 0
}

type AssetAmount struct {
	AssetId *AssetID `protobuf:"bytes,1,opt,name=asset_id,json=assetId" json:"asset_id,omitempty"`
	Amount  uint64   `protobuf:"varint,2,opt,name=amount" json:"amount,omitempty"`
}

func (m *AssetAmount) Reset()                    { *m = AssetAmount{} }
func (m *AssetAmount) String() string            { return proto.CompactTextString(m) }
func (*AssetAmount) ProtoMessage()               {}
func (*AssetAmount) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *AssetAmount) GetAssetId() *AssetID {
	if m != nil {
		return m.AssetId
	}
	return nil
}

func (m *AssetAmount) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type AssetDefinition struct {
	IssuanceProgram *Program `protobuf:"bytes,1,opt,name=issuance_program,json=issuanceProgram" json:"issuance_program,omitempty"`
	Data            *Hash    `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *AssetDefinition) Reset()                    { *m = AssetDefinition{} }
func (m *AssetDefinition) String() string            { return proto.CompactTextString(m) }
func (*AssetDefinition) ProtoMessage()               {}
func (*AssetDefinition) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *AssetDefinition) GetIssuanceProgram() *Program {
	if m != nil {
		return m.IssuanceProgram
	}
	return nil
}

func (m *AssetDefinition) GetData() *Hash {
	if m != nil {
		return m.Data
	}
	return nil
}

type ValueSource struct {
	Ref      *Hash        `protobuf:"bytes,1,opt,name=ref" json:"ref,omitempty"`
	Value    *AssetAmount `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Position uint64       `protobuf:"varint,3,opt,name=position" json:"position,omitempty"`
}

func (m *ValueSource) Reset()                    { *m = ValueSource{} }
func (m *ValueSource) String() string            { return proto.CompactTextString(m) }
func (*ValueSource) ProtoMessage()               {}
func (*ValueSource) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ValueSource) GetRef() *Hash {
	if m != nil {
		return m.Ref
	}
	return nil
}

func (m *ValueSource) GetValue() *AssetAmount {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ValueSource) GetPosition() uint64 {
	if m != nil {
		return m.Position
	}
	return 0
}

type ValueDestination struct {
	Ref      *Hash        `protobuf:"bytes,1,opt,name=ref" json:"ref,omitempty"`
	Value    *AssetAmount `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Position uint64       `protobuf:"varint,3,opt,name=position" json:"position,omitempty"`
}

func (m *ValueDestination) Reset()                    { *m = ValueDestination{} }
func (m *ValueDestination) String() string            { return proto.CompactTextString(m) }
func (*ValueDestination) ProtoMessage()               {}
func (*ValueDestination) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *ValueDestination) GetRef() *Hash {
	if m != nil {
		return m.Ref
	}
	return nil
}

func (m *ValueDestination) GetValue() *AssetAmount {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *ValueDestination) GetPosition() uint64 {
	if m != nil {
		return m.Position
	}
	return 0
}

type Proof struct {
	Sign           []byte `protobuf:"bytes,1,opt,name=sign,proto3" json:"sign,omitempty"`
	ControlProgram []byte `protobuf:"bytes,2,opt,name=controlProgram,proto3" json:"controlProgram,omitempty"`
	Address        []byte `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
}

func (m *Proof) Reset()                    { *m = Proof{} }
func (m *Proof) String() string            { return proto.CompactTextString(m) }
func (*Proof) ProtoMessage()               {}
func (*Proof) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *Proof) GetSign() []byte {
	if m != nil {
		return m.Sign
	}
	return nil
}

func (m *Proof) GetControlProgram() []byte {
	if m != nil {
		return m.ControlProgram
	}
	return nil
}

func (m *Proof) GetAddress() []byte {
	if m != nil {
		return m.Address
	}
	return nil
}

type BytomBlockHeader struct {
	Version               uint64             `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Height                uint64             `protobuf:"varint,2,opt,name=height" json:"height,omitempty"`
	PreviousBlockId       *Hash              `protobuf:"bytes,3,opt,name=previous_block_id,json=previousBlockId" json:"previous_block_id,omitempty"`
	Timestamp             uint64             `protobuf:"varint,4,opt,name=timestamp" json:"timestamp,omitempty"`
	TransactionsRoot      *Hash              `protobuf:"bytes,5,opt,name=transactions_root,json=transactionsRoot" json:"transactions_root,omitempty"`
	TransactionStatusHash *Hash              `protobuf:"bytes,6,opt,name=transaction_status_hash,json=transactionStatusHash" json:"transaction_status_hash,omitempty"`
	Nonce                 uint64             `protobuf:"varint,7,opt,name=nonce" json:"nonce,omitempty"`
	Bits                  uint64             `protobuf:"varint,8,opt,name=bits" json:"bits,omitempty"`
	TransactionStatus     *TransactionStatus `protobuf:"bytes,9,opt,name=transaction_status,json=transactionStatus" json:"transaction_status,omitempty"`
}

func (m *BytomBlockHeader) Reset()                    { *m = BytomBlockHeader{} }
func (m *BytomBlockHeader) String() string            { return proto.CompactTextString(m) }
func (*BytomBlockHeader) ProtoMessage()               {}
func (*BytomBlockHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *BytomBlockHeader) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *BytomBlockHeader) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BytomBlockHeader) GetPreviousBlockId() *Hash {
	if m != nil {
		return m.PreviousBlockId
	}
	return nil
}

func (m *BytomBlockHeader) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *BytomBlockHeader) GetTransactionsRoot() *Hash {
	if m != nil {
		return m.TransactionsRoot
	}
	return nil
}

func (m *BytomBlockHeader) GetTransactionStatusHash() *Hash {
	if m != nil {
		return m.TransactionStatusHash
	}
	return nil
}

func (m *BytomBlockHeader) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *BytomBlockHeader) GetBits() uint64 {
	if m != nil {
		return m.Bits
	}
	return 0
}

func (m *BytomBlockHeader) GetTransactionStatus() *TransactionStatus {
	if m != nil {
		return m.TransactionStatus
	}
	return nil
}

type BlockHeader struct {
	Version               uint64             `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	Height                uint64             `protobuf:"varint,2,opt,name=height" json:"height,omitempty"`
	PreviousBlockId       *Hash              `protobuf:"bytes,3,opt,name=previous_block_id,json=previousBlockId" json:"previous_block_id,omitempty"`
	Timestamp             uint64             `protobuf:"varint,4,opt,name=timestamp" json:"timestamp,omitempty"`
	TransactionsRoot      *Hash              `protobuf:"bytes,5,opt,name=transactions_root,json=transactionsRoot" json:"transactions_root,omitempty"`
	TransactionStatusHash *Hash              `protobuf:"bytes,6,opt,name=transaction_status_hash,json=transactionStatusHash" json:"transaction_status_hash,omitempty"`
	Nonce                 uint64             `protobuf:"varint,7,opt,name=nonce" json:"nonce,omitempty"`
	Bits                  uint64             `protobuf:"varint,8,opt,name=bits" json:"bits,omitempty"`
	TransactionStatus     *TransactionStatus `protobuf:"bytes,9,opt,name=transaction_status,json=transactionStatus" json:"transaction_status,omitempty"`
	Proof                 *Proof             `protobuf:"bytes,10,opt,name=Proof" json:"Proof,omitempty"`
	Extra                 []byte             `protobuf:"bytes,11,opt,name=extra,proto3" json:"extra,omitempty"`
	Coinbase              []byte             `protobuf:"bytes,12,opt,name=coinbase,proto3" json:"coinbase,omitempty"`
}

func (m *BlockHeader) Reset()                    { *m = BlockHeader{} }
func (m *BlockHeader) String() string            { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()               {}
func (*BlockHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *BlockHeader) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *BlockHeader) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHeader) GetPreviousBlockId() *Hash {
	if m != nil {
		return m.PreviousBlockId
	}
	return nil
}

func (m *BlockHeader) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *BlockHeader) GetTransactionsRoot() *Hash {
	if m != nil {
		return m.TransactionsRoot
	}
	return nil
}

func (m *BlockHeader) GetTransactionStatusHash() *Hash {
	if m != nil {
		return m.TransactionStatusHash
	}
	return nil
}

func (m *BlockHeader) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *BlockHeader) GetBits() uint64 {
	if m != nil {
		return m.Bits
	}
	return 0
}

func (m *BlockHeader) GetTransactionStatus() *TransactionStatus {
	if m != nil {
		return m.TransactionStatus
	}
	return nil
}

func (m *BlockHeader) GetProof() *Proof {
	if m != nil {
		return m.Proof
	}
	return nil
}

func (m *BlockHeader) GetExtra() []byte {
	if m != nil {
		return m.Extra
	}
	return nil
}

func (m *BlockHeader) GetCoinbase() []byte {
	if m != nil {
		return m.Coinbase
	}
	return nil
}

type TxHeader struct {
	Version        uint64  `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	SerializedSize uint64  `protobuf:"varint,2,opt,name=serialized_size,json=serializedSize" json:"serialized_size,omitempty"`
	Data           *Hash   `protobuf:"bytes,3,opt,name=data" json:"data,omitempty"`
	TimeRange      uint64  `protobuf:"varint,4,opt,name=time_range,json=timeRange" json:"time_range,omitempty"`
	ResultIds      []*Hash `protobuf:"bytes,5,rep,name=result_ids,json=resultIds" json:"result_ids,omitempty"`
}

func (m *TxHeader) Reset()                    { *m = TxHeader{} }
func (m *TxHeader) String() string            { return proto.CompactTextString(m) }
func (*TxHeader) ProtoMessage()               {}
func (*TxHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *TxHeader) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *TxHeader) GetSerializedSize() uint64 {
	if m != nil {
		return m.SerializedSize
	}
	return 0
}

func (m *TxHeader) GetData() *Hash {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *TxHeader) GetTimeRange() uint64 {
	if m != nil {
		return m.TimeRange
	}
	return 0
}

func (m *TxHeader) GetResultIds() []*Hash {
	if m != nil {
		return m.ResultIds
	}
	return nil
}

type TxVerifyResult struct {
	StatusFail bool `protobuf:"varint,1,opt,name=status_fail,json=statusFail" json:"status_fail,omitempty"`
}

func (m *TxVerifyResult) Reset()                    { *m = TxVerifyResult{} }
func (m *TxVerifyResult) String() string            { return proto.CompactTextString(m) }
func (*TxVerifyResult) ProtoMessage()               {}
func (*TxVerifyResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *TxVerifyResult) GetStatusFail() bool {
	if m != nil {
		return m.StatusFail
	}
	return false
}

type TransactionStatus struct {
	Version      uint64            `protobuf:"varint,1,opt,name=version" json:"version,omitempty"`
	VerifyStatus []*TxVerifyResult `protobuf:"bytes,2,rep,name=verify_status,json=verifyStatus" json:"verify_status,omitempty"`
}

func (m *TransactionStatus) Reset()                    { *m = TransactionStatus{} }
func (m *TransactionStatus) String() string            { return proto.CompactTextString(m) }
func (*TransactionStatus) ProtoMessage()               {}
func (*TransactionStatus) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *TransactionStatus) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *TransactionStatus) GetVerifyStatus() []*TxVerifyResult {
	if m != nil {
		return m.VerifyStatus
	}
	return nil
}

type Mux struct {
	Sources             []*ValueSource      `protobuf:"bytes,1,rep,name=sources" json:"sources,omitempty"`
	Program             *Program            `protobuf:"bytes,2,opt,name=program" json:"program,omitempty"`
	WitnessDestinations []*ValueDestination `protobuf:"bytes,3,rep,name=witness_destinations,json=witnessDestinations" json:"witness_destinations,omitempty"`
	WitnessArguments    [][]byte            `protobuf:"bytes,4,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
}

func (m *Mux) Reset()                    { *m = Mux{} }
func (m *Mux) String() string            { return proto.CompactTextString(m) }
func (*Mux) ProtoMessage()               {}
func (*Mux) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *Mux) GetSources() []*ValueSource {
	if m != nil {
		return m.Sources
	}
	return nil
}

func (m *Mux) GetProgram() *Program {
	if m != nil {
		return m.Program
	}
	return nil
}

func (m *Mux) GetWitnessDestinations() []*ValueDestination {
	if m != nil {
		return m.WitnessDestinations
	}
	return nil
}

func (m *Mux) GetWitnessArguments() [][]byte {
	if m != nil {
		return m.WitnessArguments
	}
	return nil
}

type Coinbase struct {
	WitnessDestination *ValueDestination `protobuf:"bytes,1,opt,name=witness_destination,json=witnessDestination" json:"witness_destination,omitempty"`
	Arbitrary          []byte            `protobuf:"bytes,2,opt,name=arbitrary,proto3" json:"arbitrary,omitempty"`
}

func (m *Coinbase) Reset()                    { *m = Coinbase{} }
func (m *Coinbase) String() string            { return proto.CompactTextString(m) }
func (*Coinbase) ProtoMessage()               {}
func (*Coinbase) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *Coinbase) GetWitnessDestination() *ValueDestination {
	if m != nil {
		return m.WitnessDestination
	}
	return nil
}

func (m *Coinbase) GetArbitrary() []byte {
	if m != nil {
		return m.Arbitrary
	}
	return nil
}

type Output struct {
	Source         *ValueSource `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	ControlProgram *Program     `protobuf:"bytes,2,opt,name=control_program,json=controlProgram" json:"control_program,omitempty"`
	Ordinal        uint64       `protobuf:"varint,3,opt,name=ordinal" json:"ordinal,omitempty"`
}

func (m *Output) Reset()                    { *m = Output{} }
func (m *Output) String() string            { return proto.CompactTextString(m) }
func (*Output) ProtoMessage()               {}
func (*Output) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *Output) GetSource() *ValueSource {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Output) GetControlProgram() *Program {
	if m != nil {
		return m.ControlProgram
	}
	return nil
}

func (m *Output) GetOrdinal() uint64 {
	if m != nil {
		return m.Ordinal
	}
	return 0
}

type Retirement struct {
	Source  *ValueSource `protobuf:"bytes,1,opt,name=source" json:"source,omitempty"`
	Ordinal uint64       `protobuf:"varint,2,opt,name=ordinal" json:"ordinal,omitempty"`
}

func (m *Retirement) Reset()                    { *m = Retirement{} }
func (m *Retirement) String() string            { return proto.CompactTextString(m) }
func (*Retirement) ProtoMessage()               {}
func (*Retirement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

func (m *Retirement) GetSource() *ValueSource {
	if m != nil {
		return m.Source
	}
	return nil
}

func (m *Retirement) GetOrdinal() uint64 {
	if m != nil {
		return m.Ordinal
	}
	return 0
}

type Issuance struct {
	NonceHash              *Hash             `protobuf:"bytes,1,opt,name=nonce_hash,json=nonceHash" json:"nonce_hash,omitempty"`
	Value                  *AssetAmount      `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	WitnessDestination     *ValueDestination `protobuf:"bytes,3,opt,name=witness_destination,json=witnessDestination" json:"witness_destination,omitempty"`
	WitnessAssetDefinition *AssetDefinition  `protobuf:"bytes,4,opt,name=witness_asset_definition,json=witnessAssetDefinition" json:"witness_asset_definition,omitempty"`
	WitnessArguments       [][]byte          `protobuf:"bytes,5,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
	Ordinal                uint64            `protobuf:"varint,6,opt,name=ordinal" json:"ordinal,omitempty"`
}

func (m *Issuance) Reset()                    { *m = Issuance{} }
func (m *Issuance) String() string            { return proto.CompactTextString(m) }
func (*Issuance) ProtoMessage()               {}
func (*Issuance) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

func (m *Issuance) GetNonceHash() *Hash {
	if m != nil {
		return m.NonceHash
	}
	return nil
}

func (m *Issuance) GetValue() *AssetAmount {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *Issuance) GetWitnessDestination() *ValueDestination {
	if m != nil {
		return m.WitnessDestination
	}
	return nil
}

func (m *Issuance) GetWitnessAssetDefinition() *AssetDefinition {
	if m != nil {
		return m.WitnessAssetDefinition
	}
	return nil
}

func (m *Issuance) GetWitnessArguments() [][]byte {
	if m != nil {
		return m.WitnessArguments
	}
	return nil
}

func (m *Issuance) GetOrdinal() uint64 {
	if m != nil {
		return m.Ordinal
	}
	return 0
}

type Spend struct {
	SpentOutputId      *Hash             `protobuf:"bytes,1,opt,name=spent_output_id,json=spentOutputId" json:"spent_output_id,omitempty"`
	WitnessDestination *ValueDestination `protobuf:"bytes,2,opt,name=witness_destination,json=witnessDestination" json:"witness_destination,omitempty"`
	WitnessArguments   [][]byte          `protobuf:"bytes,3,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
	Ordinal            uint64            `protobuf:"varint,4,opt,name=ordinal" json:"ordinal,omitempty"`
}

func (m *Spend) Reset()                    { *m = Spend{} }
func (m *Spend) String() string            { return proto.CompactTextString(m) }
func (*Spend) ProtoMessage()               {}
func (*Spend) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

func (m *Spend) GetSpentOutputId() *Hash {
	if m != nil {
		return m.SpentOutputId
	}
	return nil
}

func (m *Spend) GetWitnessDestination() *ValueDestination {
	if m != nil {
		return m.WitnessDestination
	}
	return nil
}

func (m *Spend) GetWitnessArguments() [][]byte {
	if m != nil {
		return m.WitnessArguments
	}
	return nil
}

func (m *Spend) GetOrdinal() uint64 {
	if m != nil {
		return m.Ordinal
	}
	return 0
}

type Claim struct {
	SpentOutputId      *Hash             `protobuf:"bytes,1,opt,name=spent_output_id,json=spentOutputId" json:"spent_output_id,omitempty"`
	WitnessDestination *ValueDestination `protobuf:"bytes,2,opt,name=witness_destination,json=witnessDestination" json:"witness_destination,omitempty"`
	WitnessArguments   [][]byte          `protobuf:"bytes,3,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
	Ordinal            uint64            `protobuf:"varint,4,opt,name=ordinal" json:"ordinal,omitempty"`
	Peginwitness       [][]byte          `protobuf:"bytes,5,rep,name=Peginwitness,proto3" json:"Peginwitness,omitempty"`
}

func (m *Claim) Reset()                    { *m = Claim{} }
func (m *Claim) String() string            { return proto.CompactTextString(m) }
func (*Claim) ProtoMessage()               {}
func (*Claim) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

func (m *Claim) GetSpentOutputId() *Hash {
	if m != nil {
		return m.SpentOutputId
	}
	return nil
}

func (m *Claim) GetWitnessDestination() *ValueDestination {
	if m != nil {
		return m.WitnessDestination
	}
	return nil
}

func (m *Claim) GetWitnessArguments() [][]byte {
	if m != nil {
		return m.WitnessArguments
	}
	return nil
}

func (m *Claim) GetOrdinal() uint64 {
	if m != nil {
		return m.Ordinal
	}
	return 0
}

func (m *Claim) GetPeginwitness() [][]byte {
	if m != nil {
		return m.Peginwitness
	}
	return nil
}

type Dpos struct {
	SpentOutputId      *Hash             `protobuf:"bytes,1,opt,name=spent_output_id,json=spentOutputId" json:"spent_output_id,omitempty"`
	WitnessDestination *ValueDestination `protobuf:"bytes,2,opt,name=witness_destination,json=witnessDestination" json:"witness_destination,omitempty"`
	WitnessArguments   [][]byte          `protobuf:"bytes,3,rep,name=witness_arguments,json=witnessArguments,proto3" json:"witness_arguments,omitempty"`
	Ordinal            uint64            `protobuf:"varint,4,opt,name=ordinal" json:"ordinal,omitempty"`
	Type               uint32            `protobuf:"varint,5,opt,name=type" json:"type,omitempty"`
	From               string            `protobuf:"bytes,6,opt,name=from" json:"from,omitempty"`
	To                 string            `protobuf:"bytes,7,opt,name=to" json:"to,omitempty"`
	Stake              uint64            `protobuf:"varint,8,opt,name=stake" json:"stake,omitempty"`
	Data               string            `protobuf:"bytes,9,opt,name=data" json:"data,omitempty"`
}

func (m *Dpos) Reset()                    { *m = Dpos{} }
func (m *Dpos) String() string            { return proto.CompactTextString(m) }
func (*Dpos) ProtoMessage()               {}
func (*Dpos) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

func (m *Dpos) GetSpentOutputId() *Hash {
	if m != nil {
		return m.SpentOutputId
	}
	return nil
}

func (m *Dpos) GetWitnessDestination() *ValueDestination {
	if m != nil {
		return m.WitnessDestination
	}
	return nil
}

func (m *Dpos) GetWitnessArguments() [][]byte {
	if m != nil {
		return m.WitnessArguments
	}
	return nil
}

func (m *Dpos) GetOrdinal() uint64 {
	if m != nil {
		return m.Ordinal
	}
	return 0
}

func (m *Dpos) GetType() uint32 {
	if m != nil {
		return m.Type
	}
	return 0
}

func (m *Dpos) GetFrom() string {
	if m != nil {
		return m.From
	}
	return ""
}

func (m *Dpos) GetTo() string {
	if m != nil {
		return m.To
	}
	return ""
}

func (m *Dpos) GetStake() uint64 {
	if m != nil {
		return m.Stake
	}
	return 0
}

func (m *Dpos) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*Hash)(nil), "bc.Hash")
	proto.RegisterType((*Program)(nil), "bc.Program")
	proto.RegisterType((*AssetID)(nil), "bc.AssetID")
	proto.RegisterType((*AssetAmount)(nil), "bc.AssetAmount")
	proto.RegisterType((*AssetDefinition)(nil), "bc.AssetDefinition")
	proto.RegisterType((*ValueSource)(nil), "bc.ValueSource")
	proto.RegisterType((*ValueDestination)(nil), "bc.ValueDestination")
	proto.RegisterType((*Proof)(nil), "bc.Proof")
	proto.RegisterType((*BytomBlockHeader)(nil), "bc.BytomBlockHeader")
	proto.RegisterType((*BlockHeader)(nil), "bc.BlockHeader")
	proto.RegisterType((*TxHeader)(nil), "bc.TxHeader")
	proto.RegisterType((*TxVerifyResult)(nil), "bc.TxVerifyResult")
	proto.RegisterType((*TransactionStatus)(nil), "bc.TransactionStatus")
	proto.RegisterType((*Mux)(nil), "bc.Mux")
	proto.RegisterType((*Coinbase)(nil), "bc.Coinbase")
	proto.RegisterType((*Output)(nil), "bc.Output")
	proto.RegisterType((*Retirement)(nil), "bc.Retirement")
	proto.RegisterType((*Issuance)(nil), "bc.Issuance")
	proto.RegisterType((*Spend)(nil), "bc.Spend")
	proto.RegisterType((*Claim)(nil), "bc.Claim")
	proto.RegisterType((*Dpos)(nil), "bc.Dpos")
}

func init() { proto.RegisterFile("bc.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1098 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x57, 0xcd, 0x6e, 0xdb, 0x46,
	0x10, 0x86, 0x24, 0xea, 0x6f, 0xa4, 0x58, 0xf6, 0xda, 0x49, 0x89, 0x20, 0x45, 0x0c, 0x02, 0x89,
	0x53, 0x14, 0x30, 0xfc, 0x93, 0xb6, 0x97, 0x1e, 0xea, 0xc4, 0x4d, 0xa3, 0x83, 0x11, 0x63, 0x6d,
	0xf8, 0x56, 0x10, 0x2b, 0x72, 0x25, 0x2f, 0x22, 0x71, 0xd9, 0xdd, 0xa5, 0x6a, 0xfb, 0x96, 0xb7,
	0xe9, 0xbd, 0x87, 0x3e, 0x42, 0x4f, 0x45, 0x1f, 0xa4, 0x2f, 0x51, 0xec, 0x70, 0x29, 0x51, 0xb2,
	0xf2, 0x87, 0xa2, 0x28, 0x02, 0xf4, 0xc6, 0x99, 0x9d, 0x9d, 0x9f, 0x6f, 0x66, 0x67, 0x86, 0xd0,
	0x1a, 0x44, 0xbb, 0xa9, 0x92, 0x46, 0x92, 0xea, 0x20, 0x0a, 0x5e, 0x80, 0xf7, 0x92, 0xe9, 0x4b,
	0xb2, 0x06, 0xd5, 0xe9, 0x9e, 0x5f, 0xd9, 0xae, 0x3c, 0x69, 0xd0, 0xea, 0x74, 0x0f, 0xe9, 0x7d,
	0xbf, 0xea, 0xe8, 0x7d, 0xa4, 0x0f, 0xfc, 0x9a, 0xa3, 0x0f, 0x90, 0x3e, 0xf4, 0x3d, 0x47, 0x1f,
	0x06, 0xdf, 0x42, 0xf3, 0x54, 0xc9, 0x91, 0x62, 0x13, 0xf2, 0x39, 0xc0, 0x74, 0x12, 0x4e, 0xb9,
	0xd2, 0x42, 0x26, 0xa8, 0xd2, 0xa3, 0xed, 0xe9, 0xe4, 0x22, 0x67, 0x10, 0x02, 0x5e, 0x24, 0x63,
	0x8e, 0xba, 0xbb, 0x14, 0xbf, 0x83, 0x3e, 0x34, 0x8f, 0xb4, 0xe6, 0xa6, 0x7f, 0xfc, 0x8f, 0x1d,
	0x39, 0x81, 0x0e, 0xaa, 0x3a, 0x9a, 0xc8, 0x2c, 0x31, 0xe4, 0x31, 0xb4, 0x98, 0x25, 0x43, 0x11,
	0xa3, 0xd2, 0xce, 0x41, 0x67, 0x77, 0x10, 0xed, 0x3a, 0x6b, 0xb4, 0x89, 0x87, 0xfd, 0x98, 0xdc,
	0x83, 0x06, 0xc3, 0x1b, 0x68, 0xca, 0xa3, 0x8e, 0x0a, 0x46, 0xd0, 0x43, 0xd9, 0x63, 0x3e, 0x14,
	0x89, 0x30, 0x36, 0x80, 0xaf, 0x61, 0x5d, 0x68, 0x9d, 0xb1, 0x24, 0xe2, 0x61, 0x9a, 0xc7, 0x5c,
	0x56, 0xed, 0x60, 0xa0, 0xbd, 0x42, 0xa8, 0xc0, 0xe5, 0x01, 0x78, 0x31, 0x33, 0x0c, 0x0d, 0x74,
	0x0e, 0x5a, 0x56, 0xd6, 0x42, 0x4f, 0x91, 0x1b, 0x8c, 0xa1, 0x73, 0xc1, 0xc6, 0x19, 0x3f, 0x93,
	0x99, 0x8a, 0x38, 0xb9, 0x0f, 0x35, 0xc5, 0x87, 0x4e, 0xef, 0x5c, 0xd6, 0x32, 0xc9, 0x23, 0xa8,
	0x4f, 0xad, 0xa8, 0xd3, 0xd4, 0x9b, 0x05, 0x94, 0xc7, 0x4c, 0xf3, 0x53, 0x72, 0x1f, 0x5a, 0xa9,
	0xd4, 0xe8, 0x33, 0xe2, 0xe5, 0xd1, 0x19, 0x1d, 0xfc, 0x04, 0xeb, 0x68, 0xed, 0x98, 0x6b, 0x23,
	0x12, 0x86, 0x71, 0xfd, 0xcb, 0x26, 0x7f, 0x84, 0xfa, 0xa9, 0x92, 0x72, 0x68, 0x0b, 0x40, 0x8b,
	0x51, 0x5e, 0x19, 0x5d, 0x8a, 0xdf, 0xe4, 0x31, 0xac, 0x45, 0x32, 0x31, 0x4a, 0x8e, 0x1d, 0x5a,
	0xae, 0x3c, 0x96, 0xb8, 0xc4, 0x87, 0x26, 0x8b, 0x63, 0xc5, 0xb5, 0x46, 0xfd, 0x5d, 0x5a, 0x90,
	0xc1, 0x9b, 0x1a, 0xac, 0x3f, 0xbb, 0x36, 0x72, 0xf2, 0x6c, 0x2c, 0xa3, 0xd7, 0x2f, 0x39, 0x8b,
	0xb9, 0xb2, 0xe2, 0x8b, 0x75, 0x58, 0x90, 0x36, 0xdf, 0x97, 0x5c, 0x8c, 0x2e, 0x67, 0xf9, 0xce,
	0x29, 0xf2, 0x14, 0x36, 0x52, 0xc5, 0xa7, 0x42, 0x66, 0x3a, 0x1c, 0x58, 0x4d, 0xb6, 0x70, 0x6a,
	0x4b, 0x90, 0xf4, 0x0a, 0x11, 0xb4, 0xd5, 0x8f, 0xc9, 0x03, 0x68, 0x1b, 0x31, 0xe1, 0xda, 0xb0,
	0x49, 0x8a, 0xb5, 0xe8, 0xd1, 0x39, 0x83, 0x7c, 0x05, 0x1b, 0x46, 0xb1, 0x44, 0xb3, 0xc8, 0x02,
	0xa1, 0x43, 0x25, 0xa5, 0xf1, 0xeb, 0x4b, 0x3a, 0xd7, 0xcb, 0x22, 0x54, 0x4a, 0x43, 0xbe, 0x83,
	0xcf, 0x4a, 0xbc, 0x50, 0x1b, 0x66, 0x32, 0x1d, 0x5e, 0x32, 0x7d, 0xe9, 0x37, 0x96, 0x2e, 0xdf,
	0x2d, 0x09, 0x9e, 0xa1, 0x1c, 0x3e, 0xea, 0x2d, 0xa8, 0x27, 0x32, 0x89, 0xb8, 0xdf, 0x44, 0x97,
	0x72, 0xc2, 0xe2, 0x3f, 0x10, 0x46, 0xfb, 0x2d, 0x64, 0xe2, 0x37, 0x39, 0x06, 0x72, 0xdb, 0x96,
	0xdf, 0x46, 0x33, 0x77, 0xad, 0x99, 0xf3, 0x65, 0x03, 0x74, 0xe3, 0x96, 0xcd, 0xe0, 0xcf, 0x1a,
	0x74, 0xfe, 0x87, 0xff, 0xbf, 0x82, 0x9f, 0x3c, 0x74, 0x2f, 0xcc, 0x07, 0xbc, 0xd8, 0x76, 0xdd,
	0x48, 0x0e, 0xa9, 0x7b, 0x79, 0x5b, 0x50, 0xe7, 0x57, 0x46, 0x31, 0xbf, 0x83, 0x6f, 0x27, 0x27,
	0xec, 0xa3, 0x8d, 0xa4, 0x48, 0x06, 0x4c, 0x73, 0xbf, 0x8b, 0x07, 0x33, 0x3a, 0xf8, 0xb5, 0x02,
	0xad, 0xf3, 0xab, 0xf7, 0xa6, 0x73, 0x07, 0x7a, 0x9a, 0x2b, 0xc1, 0xc6, 0xe2, 0x86, 0xc7, 0xa1,
	0x16, 0x37, 0xdc, 0xe5, 0x75, 0x6d, 0xce, 0x3e, 0x13, 0x37, 0x7c, 0xd6, 0x03, 0x6b, 0xab, 0x7a,
	0xa0, 0x9d, 0x1c, 0x36, 0x6d, 0xa1, 0x62, 0xc9, 0x88, 0x97, 0x13, 0x49, 0x2d, 0x83, 0xec, 0x00,
	0x28, 0xae, 0xb3, 0xb1, 0x6d, 0xe6, 0xda, 0xaf, 0x6f, 0xd7, 0x16, 0x54, 0xb4, 0xf3, 0xb3, 0x7e,
	0xac, 0x83, 0x7d, 0x58, 0x3b, 0xbf, 0xba, 0xe0, 0x4a, 0x0c, 0xaf, 0x29, 0x32, 0xc9, 0x43, 0xe8,
	0xb8, 0x04, 0x0e, 0x99, 0x18, 0xa3, 0xfb, 0x2d, 0x0a, 0x39, 0xeb, 0x05, 0x13, 0xe3, 0x60, 0x08,
	0x1b, 0xb7, 0x30, 0x7e, 0x47, 0xc0, 0xdf, 0xc0, 0x9d, 0x29, 0xea, 0x2f, 0x72, 0x55, 0x45, 0x6f,
	0x08, 0xe6, 0x6a, 0xc1, 0x34, 0xed, 0xe6, 0x82, 0xee, 0x89, 0xfc, 0x51, 0x81, 0xda, 0x49, 0x76,
	0x45, 0xbe, 0x80, 0xa6, 0xc6, 0x4e, 0xaf, 0xfd, 0x0a, 0x5e, 0xc5, 0x96, 0x5a, 0x9a, 0x00, 0xb4,
	0x38, 0x27, 0x8f, 0xa0, 0x99, 0x96, 0x9a, 0xe2, 0xd2, 0x98, 0x29, 0xce, 0xc8, 0x0f, 0xb0, 0xf5,
	0xb3, 0x30, 0x09, 0xd7, 0x3a, 0x8c, 0xe7, 0x5d, 0xdd, 0xf6, 0x49, 0xab, 0x7e, 0x6b, 0xa6, 0xbe,
	0xd4, 0xf2, 0xe9, 0xa6, 0xbb, 0x51, 0xe2, 0x69, 0xf2, 0x25, 0x6c, 0x14, 0x8a, 0x98, 0x1a, 0x65,
	0x13, 0x9e, 0x18, 0xed, 0x7b, 0xdb, 0xb5, 0x27, 0x5d, 0xba, 0xee, 0x0e, 0x8e, 0x0a, 0x7e, 0x20,
	0xa1, 0xf5, 0xdc, 0x15, 0x0b, 0xf9, 0x1e, 0x36, 0x57, 0x78, 0xe0, 0x06, 0xca, 0x6a, 0x07, 0xc8,
	0x6d, 0x07, 0xec, 0x6b, 0x66, 0x6a, 0x20, 0x8c, 0x62, 0xea, 0xda, 0x8d, 0x81, 0x39, 0x23, 0x78,
	0x53, 0x81, 0xc6, 0xab, 0xcc, 0xa4, 0x99, 0x21, 0x3b, 0xd0, 0xc8, 0x31, 0x72, 0x26, 0x6e, 0x41,
	0xe8, 0x8e, 0xc9, 0x53, 0xe8, 0xb9, 0x39, 0x12, 0xbe, 0x03, 0xc9, 0x15, 0xb3, 0x46, 0xaa, 0x58,
	0x24, 0x6c, 0xec, 0x66, 0x59, 0x41, 0x06, 0xaf, 0x00, 0x28, 0x37, 0x42, 0x71, 0x8b, 0xc1, 0x87,
	0xbb, 0x51, 0x52, 0x58, 0x5d, 0x54, 0xf8, 0x5b, 0x15, 0x5a, 0x7d, 0xb7, 0x2e, 0xd8, 0x32, 0xc7,
	0x4e, 0x91, 0xf7, 0x9a, 0xe5, 0x71, 0xdc, 0xc6, 0x33, 0xec, 0x2f, 0x1f, 0x38, 0x94, 0xdf, 0x92,
	0x96, 0xda, 0x47, 0xa6, 0xe5, 0x04, 0xfc, 0x59, 0x59, 0xe0, 0x46, 0x15, 0xcf, 0x56, 0x22, 0x7c,
	0xaa, 0x9d, 0x83, 0xcd, 0x99, 0x03, 0xf3, 0x6d, 0x89, 0xde, 0x2b, 0x4a, 0x66, 0x69, 0x8b, 0x5a,
	0x59, 0x65, 0xf5, 0xd5, 0x55, 0x56, 0x46, 0xae, 0xb1, 0x88, 0xdc, 0xef, 0x15, 0xa8, 0x9f, 0xa5,
	0x3c, 0x89, 0xc9, 0x1e, 0xf4, 0x74, 0xca, 0x13, 0x13, 0x4a, 0xac, 0x8e, 0xf9, 0xc2, 0x37, 0xc7,
	0xee, 0x0e, 0x0a, 0xe4, 0xd5, 0xd3, 0x8f, 0xdf, 0x06, 0x4c, 0xf5, 0x23, 0x81, 0x59, 0x19, 0x49,
	0xed, 0xfd, 0x91, 0x78, 0x8b, 0x91, 0xfc, 0x55, 0x81, 0xfa, 0xf3, 0x31, 0x13, 0x93, 0x4f, 0x3d,
	0x12, 0x12, 0x40, 0xf7, 0x94, 0x8f, 0x44, 0xe2, 0xae, 0xb8, 0xac, 0x2e, 0xf0, 0x82, 0x5f, 0xaa,
	0xe0, 0x1d, 0xa7, 0x52, 0x7f, 0xf2, 0xc1, 0x12, 0xf0, 0xcc, 0x75, 0xca, 0x71, 0xa1, 0xb8, 0x43,
	0xf1, 0xdb, 0xf2, 0x86, 0x4a, 0x4e, 0xb0, 0x56, 0xdb, 0x14, 0xbf, 0xed, 0x7f, 0x8a, 0x91, 0xb8,
	0x09, 0xb4, 0x69, 0xd5, 0x48, 0x3b, 0x8b, 0xb5, 0x61, 0xaf, 0xb9, 0xdb, 0x03, 0x72, 0xc2, 0xde,
	0xc4, 0xf9, 0xd8, 0xce, 0x6f, 0xda, 0xef, 0x41, 0x03, 0xff, 0xd6, 0x0e, 0xff, 0x0e, 0x00, 0x00,
	0xff, 0xff, 0xbc, 0x6e, 0xdc, 0xf9, 0xb9, 0x0d, 0x00, 0x00,
}
