# Arbitrage Smart Contracts

## Solidity Smart Contracts
To interact with a Solidity smart contract within Go, the contract must first be compiled to an ABI (application binary interface). The next step is to then convert the ABI into an importable Go file which exposes methods used to interact with the smart contract. Refer to https://goethereumbook.org/smart-contract-compile/ for further details.

### Smart Contract Compilation
Install the Solidity compiler by referring to this guide: https://docs.soliditylang.org/en/latest/installing-solidity.html

Once the Solidity compiler is installed, an ABI can be generated from a smart contract using the command below:
```
solc --abi <smart contract filename>.sol
```

### Converting an ABI to a GO file
The abigen tool from the go-ethereum package (https://github.com/ethereum/go-ethereum) must first be installed. Install abigen using the command below:
```
go install github.com/ethereum/go-ethereum/cmd/abigen@latest
```

After the abigen tool is installed, it can be used to generate a Go file based on an ABI using the command below:
```
abigen --abi=<abi filename>.abi --pkg=<package name> --out=<output filename>.go
```

Note that the "--pkg" flag is used to indicate the package of the generated Go file.

Once the Go output file is generated, it can be imported by another Go program and the methods exposed in the generated file can be called to interact with the smart contract. Refer to https://goethereumbook.org/smart-contract-read-erc20/ for a detailed example.

### UniswapV3 Smart Contracts
ABIs for UniswapV3 smart contracts are published to npm. As such, the compilation step can be skipped and ABIs can be pulled directly from npm using a Node.js script.

Refer to "contracts/uniswapV3Pool/fetchABI.js" for an example of such a script. This script fetches the ABI for the "IUniswapV3Pool.sol" smart contract and saves it to a file. After running the script, abigen must then be used to convert the ABI into an importable Go file.
