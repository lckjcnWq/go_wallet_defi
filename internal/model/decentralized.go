package model

import "github.com/gogf/gf/v2/os/gtime"

// Proposal 治理提案
type Proposal struct {
	Id            uint64      `json:"id"               description:"ID"`
	ProposalId    string      `json:"proposal_id"      description:"链上提案ID"`
	Title         string      `json:"title"            description:"标题"`
	Description   string      `json:"description"      description:"描述"`
	Proposer      string      `json:"proposer"         description:"提案人"`
	StartTime     *gtime.Time `json:"start_time"       description:"开始时间"`
	EndTime       *gtime.Time `json:"end_time"         description:"结束时间"`
	Status        int         `json:"status"           description:"状态"`
	ForVotes      string      `json:"for_votes"        description:"赞成票"`
	AgainstVotes  string      `json:"against_votes"    description:"反对票"`
	QuorumVotes   string      `json:"quorum_votes"     description:"法定票数"`
	ExecutionTime *gtime.Time `json:"execution_time"   description:"执行时间"`
	CreatedAt     *gtime.Time `json:"created_at"       description:"创建时间"`
}

// Vote 投票记录
type Vote struct {
	Id         uint64      `json:"id"           description:"ID"`
	ProposalId string      `json:"proposal_id"  description:"提案ID"`
	Voter      string      `json:"voter"        description:"投票人"`
	Support    int         `json:"support"      description:"投票选项 1:赞成 2:反对 3:弃权"`
	Weight     string      `json:"weight"       description:"票权重"`
	Reason     string      `json:"reason"       description:"投票理由"`
	CreatedAt  *gtime.Time `json:"created_at"   description:"创建时间"`
}

// IPFSFile 去中心化存储文件
type IPFSFile struct {
	Id        uint64      `json:"id"           description:"ID"`
	Hash      string      `json:"hash"         description:"IPFS哈希"`
	Name      string      `json:"name"         description:"文件名"`
	Size      int64       `json:"size"         description:"文件大小"`
	Type      string      `json:"type"         description:"文件类型"`
	Encrypted bool        `json:"encrypted"    description:"是否加密"`
	Owner     string      `json:"owner"        description:"所有者"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
}

// DIDDocument DID文档
type DIDDocument struct {
	Id             uint64      `json:"id"               description:"ID"`
	Did            string      `json:"did"              description:"DID标识"`
	PublicKey      string      `json:"public_key"       description:"公钥"`
	Authentication string      `json:"authentication"   description:"认证方式"`
	Controller     string      `json:"controller"       description:"控制者"`
	Service        string      `json:"service"          description:"服务终端"`
	CreatedAt      *gtime.Time `json:"created_at"       description:"创建时间"`
	UpdatedAt      *gtime.Time `json:"updated_at"       description:"更新时间"`
}

// P2PMessage P2P消息
type P2PMessage struct {
	Id        uint64      `json:"id"           description:"ID"`
	MessageId string      `json:"message_id"   description:"消息ID"`
	FromPeer  string      `json:"from_peer"    description:"发送方"`
	ToPeer    string      `json:"to_peer"      description:"接收方"`
	GroupId   string      `json:"group_id"     description:"群组ID"`
	Type      string      `json:"type"         description:"消息类型"`
	Content   string      `json:"content"      description:"加密内容"`
	Signature string      `json:"signature"    description:"签名"`
	Status    int         `json:"status"       description:"状态 0:待发送 1:已发送 2:已送达 3:已读"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
}

// P2PGroup P2P群组
type P2PGroup struct {
	Id        uint64      `json:"id"           description:"ID"`
	GroupId   string      `json:"group_id"     description:"群组ID"`
	Name      string      `json:"name"         description:"群组名称"`
	Owner     string      `json:"owner"        description:"群主"`
	Members   string      `json:"members"      description:"成员JSON"`
	CreatedAt *gtime.Time `json:"created_at"   description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updated_at"   description:"更新时间"`
}
