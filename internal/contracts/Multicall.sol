// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


//批量交易合约
contract Multicall {
    struct Call {
        address target;
        bytes callData;
        uint256 value;
    }

    struct Result {
        bool success;
        bytes returnData;
    }

    //批量执行只读调用
    function multicall(Call[] memory calls) public view returns (uint256 blockNumber, bytes[] memory returnData) {
        blockNumber = block.number;
        returnData = new bytes[](calls.length);
        for(uint256 i = 0; i < calls.length; i++) {
            (bool success, bytes memory ret) = calls[i].target.staticcall(calls[i].callData);
            require(success, "Multicall: call failed");
            returnData[i] = ret;
        }
    }

    // 批量执行，允许部分失败
    function tryBatch(Call[] memory calls) public payable returns (Result[] memory returnData) {
        returnData = new Result[](calls.length);
        
        for(uint256 i = 0; i < calls.length; i++) {
            (bool success, bytes memory ret) = calls[i].target.call{value: calls[i].value}(calls[i].callData);
            returnData[i] = Result(success, ret);
        }
    }

    receive() external payable {}

}
