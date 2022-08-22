// Script to fetch the UniswapV2Pair ABI from npm and save it to a file

const path = require('path');
const fs = require('fs').promises;
const {abi: IUniswapV2PairABI} = require("@uniswap/v2-core/build/IUniswapV2Pair.json");

const postScriptMsg = `\nIMPORTANT: The next step is to convert the generated ABI into an importable Go file
This can be automated into the script in the future
Run the following command after the ABI is saved:
    abigen --abi=uniswapV2Pair.abi --pkg=uniswapV2Pair --out=uniswapV2Pair.go`;

(async () => {
    try {
        const outputFilename = 'uniswapV2Pair.abi'
        await fs.writeFile(path.join(__dirname, outputFilename), JSON.stringify(IUniswapV2PairABI));
        console.log(`Success: ABI saved to '${outputFilename}'`)
        console.log(postScriptMsg)
    } catch (err) {
        console.error(err)
    }
})();
