package controller

import (
	"context"
	"go-wallet-defi/internal/model"
	"go-wallet-defi/internal/service"
)

type DecentralizedController struct{}

// CreateProposal 创建提案
func (c *DecentralizedController) CreateProposal(ctx context.Context, req *v1.CreateProposalReq) (res *v1.CreateProposalRes, err error) {
	params := &contracts.ProposalParams{
		Title:       req.Title,
		Description: req.Description,
		Targets:     req.Targets,
		Values:      req.Values,
		Calldatas:   req.Calldatas,
	}

	proposalId, err := service.Decentralized.CreateProposal(ctx, req.ChainId, params)
	if err != nil {
		return nil, err
	}

	return &v1.CreateProposalRes{
		ProposalId: proposalId,
	}, nil
}

// CastVote 投票
func (c *DecentralizedController) CastVote(ctx context.Context, req *v1.CastVoteReq) (res *v1.CastVoteRes, err error) {
	err = service.Decentralized.CastVote(ctx, req.ChainId, req.ProposalId, req.Support, req.Reason)
	if err != nil {
		return nil, err
	}

	return &v1.CastVoteRes{}, nil
}

// UploadToIPFS 上传文件到IPFS
func (c *DecentralizedController) UploadToIPFS(ctx context.Context, req *v1.UploadToIPFSReq) (res *v1.UploadToIPFSRes, err error) {
	hash, err := service.Decentralized.UploadToIPFS(ctx, req.File, req.Name, req.Encrypt)
	if err != nil {
		return nil, err
	}

	return &v1.UploadToIPFSRes{
		Hash: hash,
	}, nil
}

// CreateDID 创建DID身份
func (c *DecentralizedController) CreateDID(ctx context.Context, req *v1.CreateDIDReq) (res *v1.CreateDIDRes, err error) {
	params := &model.DIDDocument{
		Controller: req.Controller,
		Service:    req.Service,
	}

	did, err := service.Decentralized.CreateDID(ctx, params)
	if err != nil {
		return nil, err
	}

	return &v1.CreateDIDRes{
		Did: did,
	}, nil
}

// SendP2PMessage 发送P2P消息
func (c *DecentralizedController) SendP2PMessage(ctx context.Context, req *v1.SendP2PMessageReq) (res *v1.SendP2PMessageRes, err error) {
	msg := &model.P2PMessage{
		MessageId: generateMessageId(),
		ToPeer:    req.ToPeer,
		GroupId:   req.GroupId,
		Type:      req.Type,
		Content:   req.Content,
	}

	err = service.Decentralized.SendP2PMessage(ctx, p2p.CurrentNode, msg)
	if err != nil {
		return nil, err
	}

	return &v1.SendP2PMessageRes{
		MessageId: msg.MessageId,
	}, nil
}

// CreateP2PGroup 创建P2P群组
func (c *DecentralizedController) CreateP2PGroup(ctx context.Context, req *v1.CreateP2PGroupReq) (res *v1.CreateP2PGroupRes, err error) {
	groupId, err := service.Decentralized.CreateP2PGroup(ctx, req.Name, req.Members)
	if err != nil {
		return nil, err
	}

	return &v1.CreateP2PGroupRes{
		GroupId: groupId,
	}, nil
}
