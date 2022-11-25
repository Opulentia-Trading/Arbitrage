// Script to fetch the UniswapV2Router02 ABI from npm and save it to a file

const path = require('path');
const fs = require('fs').promises;
const {abi: IUniswapV2Router02ABI} = require("@uniswap/v2-periphery/build/IUniswapV2Router02.json");

const postScriptMsg = `\nIMPORTANT: The next step is to convert the generated ABI into an importable Go file
This can be automated into the script in the future
Run the following command after the ABI is saved:
    abigen --abi=uniswapV2Router02.abi --pkg=uniswapV2Router02 --out=uniswapV2Router02.go`;

(async () => {
    try {
        const outputFilename = 'uniswapV2Router02.abi'
        await fs.writeFile(path.join(__dirname, outputFilename), JSON.stringify(IUniswapV2Router02ABI));
        console.log(`Success: ABI saved to '${outputFilename}'`)
        console.log(postScriptMsg)
    } catch (err) {
        console.error(err)
    }
})();
