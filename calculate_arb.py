# 2 pairs of (base, quote) from platform 1 and 2
# We essentially use "base" as the anchor of our arbitrage
# We buy x amount of quote token on the platform with the lower quote token price
# Then we exchange it in base tokens of amount y on platform 2 (as quote token is valued more on platform 2)
# We then give z amount of base token back to platform 1 since we only burrowed from there
# Altogether, we have a profit of y-z

# (BASE1, QUOTE1)
# (BASE2, QUOTE2)
# Conditions QUOTE2 < QUOTE1

from math import sqrt
from web3 import Web3
import json
import os
from dotenv import load_dotenv

load_dotenv("env\.env")

INFURA_ID = os.getenv('INFURA_PROJECT_ID')
ZERO_ADDRESS = '0x0000000000000000000000000000000000000000'
TOKEN_1_NAME, TOKEN_1 = ("DAI", "0x6b175474e89094c44da98b954eedeac495271d0f")
TOKEN_2_NAME, TOKEN_2 = ("ETH", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
PLATFORM_1, PLATFORM_2 = ("Uniswap", "Sushiswap")
PLATFORM_1_FACTORY, PLATFORM_2_FACTORY = ("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f", "0xC0AEe478e3658e2610c5F7A4A2E1777cE9e4f2Ac")

def create_struct(action, amount, currency_received, token, platform, next=None):
    struct = {"Action": action, "Amount": amount, "To Receive": currency_received, "Token": token, "Platform": platform, "Next": next}
    return struct

def calc_max_to_buy(p0t0, p0t1, p1t0, p1t1):

    # Formulate a quadratic equation with their reserve amounts
    minimum = min((p0t0, p0t1, p1t0, p1t1))
    if (minimum > 1e24):
        d = 1e20
    elif (minimum > 1e23):
        d = 1e19
    elif (minimum > 1e22):
        d = 1e18
    elif (minimum > 1e21):
        d = 1e17
    elif (minimum > 1e20):
        d = 1e16
    elif (minimum > 1e19):
        d = 1e15
    elif (minimum > 1e18):
        d = 1e14
    elif (minimum > 1e17):
        d = 1e13
    elif (minimum > 1e16):
        d = 1e12
    elif (minimum > 1e15):
        d = 1e11
    else:
        d = 1e10
    
    p0t0, p0t1, p1t0, p1t1 = (p0t0/d, p0t1/d, p1t0/d, p1t1/d)
    term1 = p0t0 * p0t1 - p1t0 * p1t1
    term2 = 2 * p0t1 * p1t1 * (p0t0 + p1t0)
    term3 = p0t1 * p1t1 * (p0t0 * p1t1 - p1t0 * p0t1)
    # Now we solve for quadratic equation: term1*x^2 + term2*x + term3 = 0
    discriminant = (term2**2) - (4*term1*term3)

    if discriminant < 0:
        print ("This equation has no real solution")
        return -1
    
    # This is the amount we want to buy
    x1 = (-term2+sqrt(discriminant))/(2*term1)
    x2 = (-term2-sqrt(discriminant))/(2*term1)

    if (x1 > 0 and x1 < p0t1 and x1 < p1t1):
        return x1 * d
    elif (x2 > 0 and x2 < p0t1 and x2 < p1t1):
        return x2 * d
    return -1

def amount_in(to_buy, p0t0, p0t1):
    numerator = p0t0 * to_buy * 1000
    denominator = (p0t1 - to_buy) * 997
    return (numerator / denominator) + 1

def amount_out(to_buy, p0t0, p0t1):
    amount_in_with_fee = to_buy*997
    numerator = amount_in_with_fee * p0t1
    denominator = (p0t0 * 1000) + amount_in_with_fee
    return numerator / denominator

# This returns the amount needed to be borrowed in terms of quote token 1
def calculate_arb(platform_1, platform_2):
    pool0token0, pool0token1 = platform_1["tokens"]
    pool1token0, pool1token1 = platform_2["tokens"]
    
    # The amount we will be buying on platform 1 denominated in token 2
    amount_to_buy = calc_max_to_buy(pool0token0["reserves"], pool0token1["reserves"], pool1token0["reserves"], pool1token1["reserves"])
    if amount_to_buy <= 0:
        return None
    
    # The amount we will be buying on platform 1 denominated in token 1
    amount_to_pay_back = amount_in(amount_to_buy, pool0token0["reserves"], pool0token1["reserves"])

    # The amount we will be receiving on platform 2 once we sell the token 2 on platform 2
    amount_received = amount_out(amount_to_buy, pool1token1["reserves"], pool1token0["reserves"])
    if (amount_received < amount_to_pay_back):
        return None
    
    last_step = create_struct("SELL", amount_to_pay_back, pool0token0["name"], pool1token0["name"], platform_1["name"])
    second_step = create_struct("BUY", amount_received, pool1token0["name"], pool0token1["name"], platform_2["name"], last_step)
    actions = create_struct("BUY", amount_to_buy, pool0token1["name"], "Arbitrary", platform_1["name"], second_step)
    return actions

def convert_to_platform(pair, platform):
    dic_token1 = {"name": TOKEN_1_NAME, "reserves": pair[0]}
    dic_token2 = {"name": TOKEN_2_NAME, "reserves": pair[1]}
    dic_platform = {"tokens": [dic_token1, dic_token2], "name": platform}
    return dic_platform
    
if __name__ == "__main__":
    infura_url = 'https://mainnet.infura.io/v3/' + INFURA_ID
    web3 = Web3(Web3.HTTPProvider(infura_url))

    #Load uniswap ABIs
    with open('json\IUniswapV2Factory.json') as f:
        factor_json = json.load(f)
    factory_abi = factor_json["abi"]

    with open('json\IUniswapV2Pair.json') as f:
        pair_json = json.load(f)
    pair_abi = pair_json["abi"]

    # uniswap factory
    uniswap_factory_address = PLATFORM_1_FACTORY
    uniswap_factory_contract = web3.eth.contract(address=uniswap_factory_address, abi=factory_abi)
    get_pair_contract = uniswap_factory_contract.functions.getPair(Web3.toChecksumAddress(TOKEN_1), Web3.toChecksumAddress(TOKEN_2)).call()
    if (get_pair_contract == ZERO_ADDRESS):
        exit(1)

    # sushiswap factory
    sushiswap_factory_address = PLATFORM_2_FACTORY
    sushiswap_factory_contract = web3.eth.contract(address=sushiswap_factory_address, abi=factory_abi)
    get_pair2_contract = sushiswap_factory_contract.functions.getPair(Web3.toChecksumAddress(TOKEN_1), Web3.toChecksumAddress(TOKEN_2)).call()
    if (get_pair_contract == ZERO_ADDRESS):
        exit(1)
    
    # get pair contract
    pair_contract = web3.eth.contract(address=get_pair_contract, abi=pair_abi)
    pair2_contract = web3.eth.contract(address=get_pair2_contract, abi=pair_abi)
    
    # get reserves
    pair = pair_contract.functions.getReserves().call()
    pair_2 = pair2_contract.functions.getReserves().call()

    #Convert to platform structs
    platform_1 = convert_to_platform(pair, PLATFORM_1)
    platform_2 = convert_to_platform(pair_2, PLATFORM_2)
    
    
    if(pair[0]/pair[1] < pair_2[0]/pair_2[1]):
        calculate_arb(platform_1, platform_2)
    else:
        calculate_arb(platform_2, platform_1)