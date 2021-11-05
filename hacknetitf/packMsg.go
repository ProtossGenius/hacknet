// Code generated by smist. DO NOT EDIT.
package hacknetitf

import (
	"github.com/ProtossGenius/hacknet/pb/hmsg"
	"github.com/ProtossGenius/hacknet/pb/hnep"
	"github.com/ProtossGenius/hacknet/pb/hnp"
	"github.com/ProtossGenius/hacknet/pb/smn_dict"
	"github.com/ProtossGenius/hacknet/pb/udpfwd"
	"google.golang.org/protobuf/types/known/anypb"
)

/*@SMIST
setIgnoreInput(true);
include('parseProtos.js');
packMsgs("./protos/hnp.proto")
packMsgs("./protos/hnep.proto")
packMsgs("./protos/udpfwd.proto")
*/
// Pack_hnp_Register pack message hnp_Register.
func Pack_hnp_Register(email string, msg *hnp.Register) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnp_Register), Msg : any}, nil
}

// Pack_hnp_Result pack message hnp_Result.
func Pack_hnp_Result(email string, msg *hnp.Result) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnp_Result), Msg : any}, nil
}

// Pack_hnp_CheckEmail pack message hnp_CheckEmail.
func Pack_hnp_CheckEmail(email string, msg *hnp.CheckEmail) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnp_CheckEmail), Msg : any}, nil
}

// Pack_hnp_Forward pack message hnp_Forward.
func Pack_hnp_Forward(email string, msg *hnp.Forward) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnp_Forward), Msg : any}, nil
}

// Pack_hnp_ForwardMsg pack message hnp_ForwardMsg.
func Pack_hnp_ForwardMsg(email string, msg *hnp.ForwardMsg) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnp_ForwardMsg), Msg : any}, nil
}

// Pack_hnp_HeartJump pack message hnp_HeartJump.
func Pack_hnp_HeartJump(email string, msg *hnp.HeartJump) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnp_HeartJump), Msg : any}, nil
}

// Pack_hnep_HackAsk pack message hnep_HackAsk.
func Pack_hnep_HackAsk(email string, msg *hnep.HackAsk) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnep_HackAsk), Msg : any}, nil
}

// Pack_hnep_BeHackAns pack message hnep_BeHackAns.
func Pack_hnep_BeHackAns(email string, msg *hnep.BeHackAns) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnep_BeHackAns), Msg : any}, nil
}

// Pack_hnep_StrMsg pack message hnep_StrMsg.
func Pack_hnep_StrMsg(email string, msg *hnep.StrMsg) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_hnep_StrMsg), Msg : any}, nil
}

// Pack_udpfwd_UDPFwdMsg pack message udpfwd_UDPFwdMsg.
func Pack_udpfwd_UDPFwdMsg(email string, msg *udpfwd.UDPFwdMsg) (resp *hmsg.Message, err error) {
	var any *anypb.Any
	if any, err = anypb.New(msg); err != nil {
		return nil, err
	}

	return &hmsg.Message{Email: email, Enum: int32(smn_dict.EDict_udpfwd_UDPFwdMsg), Msg : any}, nil
}

