package dao

import (
	"context"
	"go-wallet-defi/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type DecentralizedDao struct{}

var Decentralized = &DecentralizedDao{}

// CreateProposal 创建提案
func (d *DecentralizedDao) CreateProposal(ctx context.Context, proposal *model.Proposal) error {
	_, err := g.DB().Model("proposal").Data(proposal).Insert()
	return err
}

// UpdateProposal 更新提案
func (d *DecentralizedDao) UpdateProposal(ctx context.Context, proposalId string, data g.Map) error {
	_, err := g.DB().Model("proposal").Where("proposal_id", proposalId).Data(data).Update()
	return err
}

// GetProposal 获取提案
func (d *DecentralizedDao) GetProposal(ctx context.Context, proposalId string) (*model.Proposal, error) {
	var proposal *model.Proposal
	err := g.DB().Model("proposal").Where("proposal_id", proposalId).Scan(&proposal)
	return proposal, err
}

// CreateVote 创建投票
func (d *DecentralizedDao) CreateVote(ctx context.Context, vote *model.Vote) error {
	_, err := g.DB().Model("vote").Data(vote).Insert()
	return err
}

// GetProposalVotes 获取提案投票
func (d *DecentralizedDao) GetProposalVotes(ctx context.Context, proposalId string) ([]*model.Vote, error) {
	var votes []*model.Vote
	err := g.DB().Model("vote").Where("proposal_id", proposalId).Scan(&votes)
	return votes, err
}

// CreateIPFSFile 创建IPFS文件记录
func (d *DecentralizedDao) CreateIPFSFile(ctx context.Context, file *model.IPFSFile) error {
	_, err := g.DB().Model("ipfs_file").Data(file).Insert()
	return err
}

// GetIPFSFile 获取IPFS文件
func (d *DecentralizedDao) GetIPFSFile(ctx context.Context, hash string) (*model.IPFSFile, error) {
	var file *model.IPFSFile
	err := g.DB().Model("ipfs_file").Where("hash", hash).Scan(&file)
	return file, err
}

// CreateDIDDocument 创建DID文档
func (d *DecentralizedDao) CreateDIDDocument(ctx context.Context, doc *model.DIDDocument) error {
	_, err := g.DB().Model("did_document").Data(doc).Insert()
	return err
}

// UpdateDIDDocument 更新DID文档
func (d *DecentralizedDao) UpdateDIDDocument(ctx context.Context, did string, data g.Map) error {
	_, err := g.DB().Model("did_document").Where("did", did).Data(data).Update()
	return err
}

// GetDIDDocument 获取DID文档
func (d *DecentralizedDao) GetDIDDocument(ctx context.Context, did string) (*model.DIDDocument, error) {
	var doc *model.DIDDocument
	err := g.DB().Model("did_document").Where("did", did).Scan(&doc)
	return doc, err
}

// CreateP2PMessage 创建P2P消息
func (d *DecentralizedDao) CreateP2PMessage(ctx context.Context, msg *model.P2PMessage) error {
	_, err := g.DB().Model("p2p_message").Data(msg).Insert()
	return err
}

// UpdateP2PMessage 更新P2P消息
func (d *DecentralizedDao) UpdateP2PMessage(ctx context.Context, messageId string, data g.Map) error {
	_, err := g.DB().Model("p2p_message").Where("message_id", messageId).Data(data).Update()
	return err
}

// GetP2PMessages 获取P2P消息
func (d *DecentralizedDao) GetP2PMessages(ctx context.Context, peer string) ([]*model.P2PMessage, error) {
	var messages []*model.P2PMessage
	err := g.DB().Model("p2p_message").
		Where("to_peer", peer).
		Or("from_peer", peer).
		Order("created_at DESC").
		Scan(&messages)
	return messages, err
}

// CreateP2PGroup 创建P2P群组
func (d *DecentralizedDao) CreateP2PGroup(ctx context.Context, group *model.P2PGroup) error {
	_, err := g.DB().Model("p2p_group").Data(group).Insert()
	return err
}

// UpdateP2PGroup 更新P2P群组
func (d *DecentralizedDao) UpdateP2PGroup(ctx context.Context, groupId string, data g.Map) error {
	_, err := g.DB().Model("p2p_group").Where("group_id", groupId).Data(data).Update()
	return err
}

// GetP2PGroup 获取P2P群组
func (d *DecentralizedDao) GetP2PGroup(ctx context.Context, groupId string) (*model.P2PGroup, error) {
	var group *model.P2PGroup
	err := g.DB().Model("p2p_group").Where("group_id", groupId).Scan(&group)
	return group, err
}
