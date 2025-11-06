// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract ProtectedNFT is ERC721URIStorage, Ownable {
    enum Category { Literature, Music, Art }

    // Global variables
    address public provider;
    string[] public authorIDs; // Array of 12-digit CCCD for multiple authors
    address[] public authorAddress; 
    Category public category;
    string[] public ownerIDs; // 12-digit CCCD
    uint256 public usageExpiry; // Unix timestamp for usage expiry

    struct NFTInfo {
        uint256 protectionExpiry; // Unix timestamp for expiry
        address[] owners; // List of owner addresses
        bool isDerivative;
        string[] platforms; // For derivatives
    }

    mapping(uint256 => NFTInfo) public nftInfo;
    uint256 private _tokenIdCounter;

    constructor(
        string memory _name,
        string memory _code,
        string[] memory _authorIDs,
        uint8 _category,
        string[] memory _currentOwnerIDs,
        uint256 _usageExpiry
    ) ERC721(_name, _code) Ownable(msg.sender) {
        require(bytes(_name).length > 0, "Name cannot be empty");
        require(bytes(_code).length > 0, "Code cannot be empty");
        require(_authorIDs.length > 0, "At least one author ID required");
        for (uint256 i = 0; i < _authorIDs.length; i++) {
            require(bytes(_authorIDs[i]).length == 12, "Each author ID must be 12 digits");
        }
        require(_category <= 2, "Invalid category");

        for (uint256 i = 0; i < _currentOwnerIDs.length; i++) {
            require(bytes(_currentOwnerIDs[i]).length == 12, "Each owner ID must be 12 digits");
        }
        
        authorIDs = _authorIDs;
        ownerIDs = _currentOwnerIDs;
        category = Category(_category);
        usageExpiry = _usageExpiry;
        provider = msg.sender;
    }


    // Create derivative NFT based on an original
    function createDerivative(
        string[] memory _platforms,
        string memory _uri,
        address[] memory owners,
        uint256 protectionExpiry
    ) public onlyOwner {
        
        _tokenIdCounter++;
        uint256 tokenId = _tokenIdCounter;
        _safeMint(msg.sender, tokenId);
        _setTokenURI(tokenId, _uri);

        nftInfo[tokenId] = NFTInfo({
            owners: owners,
            isDerivative: true,
            platforms: _platforms,
            protectionExpiry: protectionExpiry
        });
    }

    // Update author IDs (only by issuer)
    function updateAuthorIDs(string[] memory _newAuthorIDs) public onlyOwner {
        require(_newAuthorIDs.length > 0, "At least one author ID required");
        for (uint256 i = 0; i < _newAuthorIDs.length; i++) {
            require(bytes(_newAuthorIDs[i]).length == 12, "Each author ID must be 12 digits");
        }
        authorIDs = _newAuthorIDs;
    }

    // Update owner IDs (only by issuer)
    function updateOwnerIDs(string[] memory _newOwnerIDs) public onlyOwner {
        require(_newOwnerIDs.length > 0, "At least one owner ID required");
        for (uint256 i = 0; i < _newOwnerIDs.length; i++) {
            require(bytes(_newOwnerIDs[i]).length == 12, "Each owner ID must be 12 digits");
        }
        ownerIDs = _newOwnerIDs;
    }

}