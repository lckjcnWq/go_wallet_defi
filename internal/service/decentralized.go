package service

import (
	"bytes"
	"context"
	"encoding/json"
	"go-wallet-defi/internal/dao"
	"go-wallet-defi/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
)

type DecentralizedService struct{}

var Decentralized = &DecentralizedService{}

// CreateProposal 创建治理提案
func (s *DecentralizedService) CreateProposal(ctx context.Context, chainId uint64, params *contracts.ProposalParams) (string, error) {
	// 获取治理合约
	governance, err := contracts.GetGovernanceContract(chainId)
	if err != nil {
		return "", err
	}

	// 创建提案
	tx, err := governance.Propose(
		params.Targets,
		params.Values,
		params.Calldatas,
		params.Description,
	)
	if err != nil {
		return "", err
	}

	// 等待交易确认
	receipt, err := ethclient.WaitForTxReceipt(ctx, chainId, tx.Hash())
	if err != nil {
		return "", err
	}

	// 解析提案ID
	proposalId := receipt.Logs[0].Topics[1].String()

	// 保存提案记录
	proposal := &model.Proposal{
		ProposalId:   proposalId,
		Title:        params.Title,
		Description:  params.Description,
		Proposer:     tx.From().Hex(),
		StartTime:    nil, // 将在区块确认后更新
		EndTime:      nil,
		Status:       0,
		ForVotes:     "0",
		AgainstVotes: "0",
		QuorumVotes:  "0",
	}

	err = dao.Decentralized.CreateProposal(ctx, proposal)
	if err != nil {
		return "", err
	}

	return proposalId, nil
}

// CastVote 投票
func (s *DecentralizedService) CastVote(ctx context.Context, chainId uint64, proposalId string, support uint8, reason string) error {
	governance, err := contracts.GetGovernanceContract(chainId)
	if err != nil {
		return err
	}

	// 投票
	tx, err := governance.CastVoteWithReason(proposalId, support, reason)
	if err != nil {
		return err
	}

	// 等待交易确认
	receipt, err := ethclient.WaitForTxReceipt(ctx, chainId, tx.Hash())
	if err != nil {
		return err
	}

	// 解析投票权重
	weight := receipt.Logs[0].Topics[3].Big()

	// 保存投票记录
	vote := &model.Vote{
		ProposalId: proposalId,
		Voter:      tx.From().Hex(),
		Support:    int(support),
		Weight:     weight.String(),
		Reason:     reason,
	}

	return dao.Decentralized.CreateVote(ctx, vote)
}

// UploadToIPFS 上传文件到IPFS
func (s *DecentralizedService) UploadToIPFS(ctx context.Context, file []byte, name string, encrypt bool) (string, error) {
	// 创建IPFS客户端
	sh := shell.NewShell("localhost:5001")

	// 如果需要加密
	if encrypt {
		key := generateEncryptionKey()
		encryptedFile, err := encryptFile(file, key)
		if err != nil {
			return "", err
		}
		file = encryptedFile
	}

	// 上传到IPFS
	hash, err := sh.Add(bytes.NewReader(file))
	if err != nil {
		return "", err
	}

	// 保存文件记录
	ipfsFile := &model.IPFSFile{
		Hash:      hash,
		Name:      name,
		Size:      int64(len(file)),
		Type:      getFileType(name),
		Encrypted: encrypt,
		Owner:     getCurrentUser(ctx),
	}

	err = dao.Decentralized.CreateIPFSFile(ctx, ipfsFile)
	if err != nil {
		return "", err
	}

	return hash, nil
}

// CreateDID 创建DID身份
func (s *DecentralizedService) CreateDID(ctx context.Context, params *model.DIDDocument) (string, error) {
	// 生成密钥对
	privateKey, publicKey := generateKeyPair()

	// 生成DID标识
	did := generateDID(publicKey)

	// 创建DID文档
	doc := &model.DIDDocument{
		Did:            did,
		PublicKey:      publicKey,
		Authentication: "secp256k1",
		Controller:     params.Controller,
		Service:        params.Service,
	}

	err := dao.Decentralized.CreateDIDDocument(ctx, doc)
	if err != nil {
		return "", err
	}

	return did, nil
}

// StartP2PNode 启动P2P节点
func (s *DecentralizedService) StartP2PNode(ctx context.Context) (*p2p.Node, error) {
	// 创建libp2p主机
	host, err := libp2p.New()
	if err != nil {
		return nil, err
	}

	// 创建P2P节点
	node := &p2p.Node{
		Host:     host,
		PeerID:   host.ID(),
		Messages: make(chan *model.P2PMessage, 100),
	}

	// 启动消息处理
	go s.handleP2PMessages(ctx, node)

	return node, nil
}

// SendP2PMessage 发送P2P消息
func (s *DecentralizedService) SendP2PMessage(ctx context.Context, node *p2p.Node, msg *model.P2PMessage) error {
	// 加密消息内容
	encryptedContent, err := encryptMessage(msg.Content, msg.ToPeer)
	if err != nil {
		return err
	}

	// 签名消息
	signature, err := signMessage(msg.MessageId, node.PrivateKey)
	if err != nil {
		return err
	}

	msg.Content = encryptedContent
	msg.Signature = signature
	msg.FromPeer = node.PeerID.String()

	// 保存消息记录
	err = dao.Decentralized.CreateP2PMessage(ctx, msg)
	if err != nil {
		return err
	}

	// 发送消息
	toPeerID, _ := peer.Decode(msg.ToPeer)
	return node.SendMessage(ctx, toPeerID, msg)
}

// CreateP2PGroup 创建P2P群组
func (s *DecentralizedService) CreateP2PGroup(ctx context.Context, name string, members []string) (string, error) {
	// 生成群组ID
	groupId := generateGroupId()

	// 创建群组记录
	membersJson, _ := json.Marshal(members)
	group := &model.P2PGroup{
		GroupId: groupId,
		Name:    name,
		Owner:   getCurrentUser(ctx),
		Members: string(membersJson),
	}

	err := dao.Decentralized.CreateP2PGroup(ctx, group)
	if err != nil {
		return "", err
	}

	return groupId, nil
}

// 处理P2P消息
func (s *DecentralizedService) handleP2PMessages(ctx context.Context, node *p2p.Node) {
	for {
		select {
		case msg := <-node.Messages:
			// 验证签名
			if !verifySignature(msg.MessageId, msg.Signature, msg.FromPeer) {
				continue
			}

			// 解密消息
			decryptedContent, err := decryptMessage(msg.Content, node.PrivateKey)
			if err != nil {
				continue
			}
			msg.Content = decryptedContent

			// 更新消息状态
			_ = dao.Decentralized.UpdateP2PMessage(ctx, msg.MessageId, g.Map{
				"status": 2, // 已送达
			})

			// 触发消息处理回调
			if node.MessageHandler != nil {
				node.MessageHandler(msg)
			}

		case <-ctx.Done():
			return
		}
	}
}
