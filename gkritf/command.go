package gkritf

import "github.com/ProtossGenius/SureMoonNet/basis/smn_data"

// GeekerMsg geeker msg.
type GeekerMsg string

// NoticeC command notice from client.
type NoticeC struct {
	TargetSession string
	ExtraData     string
}

// NoticeS command notice from server.
type NoticeS struct {
	// sender's sessionID
	NodeSessionID string
	// sender's nodeInfo
	NodeInfo  GeekerNodeInfo
	ExtraData string
}

// SearchC search from client.
type SearchC struct {
	SessionID string
}

// SearchS search from server.
type SearchS struct {
	NodeInfo GeekerNodeInfo
}

// RegisterC register from client.
type RegisterC struct {
	SessionID string
}

// ReigsterS register from server.
type ReigsterS struct {
}

// Hole hole.
type Hole struct {
	Msg string
}

/*@SMIST
setIgnoreInput(true);
include('parseCmd.js');
parse();
*/

// NewNoticeC GeekerMsg to NoticeC.
func NewNoticeC(params GeekerMsg) (NoticeC, error){
	cmd := NoticeC{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message NoticeC to geeker message.
func (c NoticeC) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("NoticeC#" + str)
}


// NewNoticeS GeekerMsg to NoticeS.
func NewNoticeS(params GeekerMsg) (NoticeS, error){
	cmd := NoticeS{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message NoticeS to geeker message.
func (c NoticeS) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("NoticeS#" + str)
}


// NewSearchC GeekerMsg to SearchC.
func NewSearchC(params GeekerMsg) (SearchC, error){
	cmd := SearchC{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message SearchC to geeker message.
func (c SearchC) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("SearchC#" + str)
}


// NewSearchS GeekerMsg to SearchS.
func NewSearchS(params GeekerMsg) (SearchS, error){
	cmd := SearchS{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message SearchS to geeker message.
func (c SearchS) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("SearchS#" + str)
}


// NewRegisterC GeekerMsg to RegisterC.
func NewRegisterC(params GeekerMsg) (RegisterC, error){
	cmd := RegisterC{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message RegisterC to geeker message.
func (c RegisterC) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("RegisterC#" + str)
}


// NewReigsterS GeekerMsg to ReigsterS.
func NewReigsterS(params GeekerMsg) (ReigsterS, error){
	cmd := ReigsterS{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message ReigsterS to geeker message.
func (c ReigsterS) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("ReigsterS#" + str)
}


// NewHole GeekerMsg to Hole.
func NewHole(params GeekerMsg) (Hole, error){
	cmd := Hole{}
	err := smn_data.GetDataFromStr(string(params), &cmd)
	return cmd, err
}

// Message Hole to geeker message.
func (c Hole) Message() GeekerMsg {
	str, _ := smn_data.ValToJson(c)
	return GeekerMsg("Hole#" + str)
}

