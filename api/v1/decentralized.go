package v1

import "github.com/gogf/gf/v2/frame/g"

type CreateProposalReq struct {
	g.Meta      `path:"/dao/proposal/create" method:"post"`
	ChainId     uint64   `json:"chain_id"    v:"required"`
	Title       string   `json:"title"       v:"required"`
	Description string   `json:"description" v:"required"`
	Targets     []string `json:"targets"     v:"required"`
	Values      []string `json:"values"      v:"required"`
	Calldatas   []string `json:"calldatas"   v:"required"`
}

type CreateProposalRes struct {
	ProposalId string `json:"proposal_id"`
}

type CastVoteReq struct {
	g.Meta     `path:"/dao/vote" method:"post"`
	ChainId    uint64 `json:"chain_id"    v:"required"`
	ProposalId string `json:"proposal_id" v:"required"`
	Support    uint8  `json:"support"     v:"required|in:1,2,3"`
	Reason     string `json:"reason"`
}

type CastVoteRes struct{}

type UploadToIPFSReq struct {
	g.Meta  `path:"/ipfs/upload" method:"post"`
	File    []byte `json:"file"     v:"required"`
	Name    string `json:"name"     v:"required"`
	Encrypt bool   `json:"encrypt"`
}

type UploadToIPFSRes struct {
	Hash string `json:"hash"`
}

type CreateDIDReq struct {
	g.Meta     `path:"/did/create" method:"post"`
	Controller string `json:"controller" v:"required"`
	Service    string `json:"service"`
}

type CreateDIDRes struct {
	Did string `json:"did"`
}

type SendP2PMessageReq struct {
	g.Meta  `path:"/p2p/message/send" method:"post"`
	ToPeer  string `json:"to_peer"  v:"required"`
	GroupId string `json:"group_id"`
	Type    string `json:"type"     v:"required"`
	Content string `json:"content"  v:"required"`
}

type SendP2PMessageRes struct {
	MessageId string `json:"message_id"`
}

type CreateP2PGroupReq struct {
	g.Meta  `path:"/p2p/group/create" method:"post"`
	Name    string   `json:"name"    v:"required"`
	Members []string `json:"members" v:"required"`
}

type CreateP2PGroupRes struct {
	GroupId string `json:"group_id"`
}
