package bridge

// 跨链桥合约ABI
const BridgeABI = `[
    {
        "inputs": [
            {"name": "token", "type": "address"},
            {"name": "amount", "type": "uint256"},
            {"name": "toChainId", "type": "uint256"},
            {"name": "toAddress", "type": "address"},
            {"name": "nonce", "type": "uint256"}
        ],
        "name": "lock",
        "outputs": [],
        "stateMutability": "payable",
        "type": "function"
    },
    {
        "inputs": [
            {"name": "token", "type": "address"},
            {"name": "amount", "type": "uint256"},
            {"name": "toAddress", "type": "address"},
            {"name": "fromChainId", "type": "uint256"},
            {"name": "nonce", "type": "uint256"},
            {"name": "signature", "type": "bytes"}
        ],
        "name": "unlock",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "anonymous": false,
        "inputs": [
            {"indexed": true, "name": "token", "type": "address"},
            {"indexed": true, "name": "from", "type": "address"},
            {"indexed": false, "name": "amount", "type": "uint256"},
            {"indexed": false, "name": "toChainId", "type": "uint256"},
            {"indexed": false, "name": "toAddress", "type": "address"},
            {"indexed": false, "name": "nonce", "type": "uint256"}
        ],
        "name": "Lock",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {"indexed": true, "name": "token", "type": "address"},
            {"indexed": true, "name": "to", "type": "address"},
            {"indexed": false, "name": "amount", "type": "uint256"},
            {"indexed": false, "name": "fromChainId", "type": "uint256"},
            {"indexed": false, "name": "nonce", "type": "uint256"}
        ],
        "name": "Unlock",
        "type": "event"
    }
]`
