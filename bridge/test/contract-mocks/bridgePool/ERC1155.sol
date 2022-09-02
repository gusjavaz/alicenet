// SPDX-License-Identifier: MIT-open-group
pragma solidity ^0.8.16;
import "@openzeppelin/contracts/token/ERC1155/ERC1155.sol";

contract ERC1155Mock is ERC1155 {
    uint256 private constant _TOKEN_ID = 1;
    bytes private constant _DATA = abi.encodePacked("");

    constructor() ERC1155("https://ifps.io/{id}") {}

    function mint(address to, uint256 number) external {
        _mint(to, _TOKEN_ID, number, _DATA);
    }
}
