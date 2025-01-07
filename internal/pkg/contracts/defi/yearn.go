package defi

// Yearn Vault ABI
const YearnVaultABI = `[
    {
        "inputs": [{"name": "_amount", "type": "uint256"}],
        "name": "deposit",
        "outputs": [{"name": "", "type": "uint256"}],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {"name": "_shares", "type": "uint256"},
            {"name": "recipient", "type": "address"}
        ],
        "name": "withdraw",
        "outputs": [{"name": "", "type": "uint256"}],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "pricePerShare",
        "outputs": [{"name": "", "type": "uint256"}],
        "stateMutability": "view",
        "type": "function"
    }
]`
