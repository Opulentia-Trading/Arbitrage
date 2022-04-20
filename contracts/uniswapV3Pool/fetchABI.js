// Script to fetch the UniswapV3Pool ABI from npm and save it to a file

const path = require('path');
const fs = require('fs').promises;
const {abi: IUniswapV3PoolABI} = require("@uniswap/v3-core/artifacts/contracts/interfaces/IUniswapV3Pool.sol/IUniswapV3Pool.json");

const postScriptMsg = `\nIMPORTANT: The next step is to convert the generated ABI into an importable Go file
This can be automated into the script in the future
Run the following command after the ABI is saved:
    abigen --abi=uniswapV3Pool.abi --pkg=uniswapV3Pool --out=uniswapV3Pool.go`;

(async () => {
    try {
        const outputFilename = 'uniswapV3Pool.abi'
        await fs.writeFile(path.join(__dirname, outputFilename), JSON.stringify(IUniswapV3PoolABI));
        console.log(`Success: ABI saved to '${outputFilename}'`)
        console.log(postScriptMsg)
    } catch (err) {
        console.error(err)
    }
})();
