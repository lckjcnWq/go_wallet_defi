// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/governance/Governor.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorSettings.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorCountingSimple.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorVotes.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorVotesQuorumFraction.sol";
import "@openzeppelin/contracts/governance/extensions/GovernorTimelockControl.sol";

contract GovernanceContract is 
    Governor, 
    GovernorSettings,
    GovernorCountingSimple,
    GovernorVotes,
    GovernorVotesQuorumFraction,
    GovernorTimelockControl 
{
    constructor(
        IVotes _token,
        TimelockController _timelock,
        uint256 _votingDelay,
        uint256 _votingPeriod,
        uint256 _proposalThreshold,
        uint256 _quorumPercentage
    )
        Governor("GovernanceContract")
        GovernorSettings(_votingDelay, _votingPeriod, _proposalThreshold)
        GovernorVotes(_token)
        GovernorVotesQuorumFraction(_quorumPercentage)
        GovernorTimelockControl(_timelock)
    {}

    // 提案状态变更事件
    event ProposalStatusChanged(uint256 indexed proposalId, ProposalState status);
    
    // 投票事件
    event VoteCast(
        address indexed voter,
        uint256 indexed proposalId,
        uint8 support,
        uint256 weight,
        string reason
    );

    // 提案执行事件
    event ProposalExecuted(uint256 indexed proposalId);

    // 提供提案详细信息
    function getProposalDetails(uint256 proposalId) 
        public 
        view 
        returns (
            address[] memory targets,
            uint256[] memory values,
            bytes[] memory calldatas,
            string memory description
        ) 
    {
        return proposalDetails[proposalId];
    }

    // 获取提案投票统计
    function getProposalVotes(uint256 proposalId)
        public
        view
        returns (
            uint256 againstVotes,
            uint256 forVotes,
            uint256 abstainVotes
        )
    {
        ProposalVote storage proposalVote = _proposalVotes[proposalId];
        return (
            proposalVote.againstVotes,
            proposalVote.forVotes,
            proposalVote.abstainVotes
        );
    }

    // 检查地址是否可以投票
    function hasVoted(uint256 proposalId, address account) 
        public 
        view 
        override 
        returns (bool) 
    {
        return _proposalVotes[proposalId].hasVoted[account];
    }

    // 提案状态变更时的钩子
    function _afterProposalStatusChange(uint256 proposalId, ProposalState status)
        internal
        override
    {
        emit ProposalStatusChanged(proposalId, status);
    }

    // 投票时的钩子
    function _countVote(
        uint256 proposalId,
        address account,
        uint8 support,
        uint256 weight,
        bytes memory params
    ) internal override {
        super._countVote(proposalId, account, support, weight, params);
        emit VoteCast(account, proposalId, support, weight, "");
    }

    // 提案执行时的钩子
    function _execute(
        uint256 proposalId,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        bytes32 descriptionHash
    ) internal override {
        super._execute(proposalId, targets, values, calldatas, descriptionHash);
        emit ProposalExecuted(proposalId);
    }
}