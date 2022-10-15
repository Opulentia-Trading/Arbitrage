# Inspiration: https://github.com/benleim/pathfinder
from web3 import Web3
from math import log
import json
import os
from dotenv import load_dotenv

load_dotenv("env\.env")

INFURA_ID = os.getenv('INFURA_PROJECT_ID')
ZERO_ADDRESS = '0x0000000000000000000000000000000000000000'
TOKEN_1_NAME, TOKEN_1 = ("DAI", "0x6b175474e89094c44da98b954eedeac495271d0f")
TOKEN_2_NAME, TOKEN_2 = ("USDC", "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
TOKEN_3_NAME, TOKEN_3 = ("ETH", "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2")
PLATFORM_1_FACTORY, PLATFORM_2_FACTORY = ("0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f", "0xC0AEe478e3658e2610c5F7A4A2E1777cE9e4f2Ac")

class Graph(object):

  def __init__(self):
    self.dist = {}
    self.edges = []
    self.previous_v = {}

  def add_edge(self, src, dst, wgt):
    self.edges.append([src, dst, wgt])
    self.dist[src] = float('Inf')
    self.dist[dst] = float('Inf')

  def bellman_ford(self, source):
    number_of_vrtx = len(self.dist)
    self.dist[source] = 0

    for _ in range(number_of_vrtx - 1):
      for u, v, w in self.edges:
        if self.dist[v] > self.dist[u] + w:
          self.dist[v] = self.dist[u] + w
          self.previous_v[v] = u

    cycle_detected = False
    for _ in range(number_of_vrtx - 1):
      for u, v, w in self.edges:
        if self.dist[v] > self.dist[u] + w:
          cycle_detected = True
          C = v
            # cyclePath = {}

            # curr = u
            # index = 1
            # cyclePath[u] = index
            # index += 1
            # print("HELLO")

            # while not cyclePath[curr]:
            #     cyclePath[curr] = index
            #     index += 1
            #     curr = self.previous_v[curr]
            
            # cyclePath[curr+'_'] = index
            
            # path = []
            # for k in cyclePath.keys():
            #     path += (k.replace('_',''))
            # path.reverse()
            # print(cyclePath)
            # for i in range(len(path)):
            #     if (i != 0 and path[0] == path[i]):
            #         path = path[0, i+1]
            #         break
            # cyclePaths += path
    cycle = []
    if cycle_detected:
      for i in range(number_of_vrtx - 1):      
            C = self.previous_v[C]
  
           
      v = C
         
      while (True):
          cycle.append(v)
          if (v == C and len(cycle) > 1):
              break
          v = self.previous_v[v]
  
      # Reverse cycle[]
      cycle.reverse()
    return cycle


def create_graph(pairs):
  g = Graph()
  for token1, token2, reserves in pairs:
    g.add_edge(token1, token2, -log(reserves[0]))
    g.add_edge(token2, token1, -log(reserves[1]))
  return g

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

    get_pair2_contract = uniswap_factory_contract.functions.getPair(Web3.toChecksumAddress(TOKEN_2), Web3.toChecksumAddress(TOKEN_3)).call()
    if (get_pair2_contract == ZERO_ADDRESS):
        exit(1)

    get_pair3_contract = uniswap_factory_contract.functions.getPair(Web3.toChecksumAddress(TOKEN_3), Web3.toChecksumAddress(TOKEN_1)).call()
    if (get_pair2_contract == ZERO_ADDRESS):
        exit(1)

    pair_contract = web3.eth.contract(address=get_pair_contract, abi=pair_abi)
    pair2_contract = web3.eth.contract(address=get_pair2_contract, abi=pair_abi)
    pair3_contract = web3.eth.contract(address=get_pair3_contract, abi=pair_abi)

    pair = [TOKEN_1_NAME, TOKEN_2_NAME, pair_contract.functions.getReserves().call()[:2]]
    pair2 = [TOKEN_2_NAME, TOKEN_3_NAME, pair2_contract.functions.getReserves().call()[:2]]
    pair3 = [TOKEN_3_NAME, TOKEN_1_NAME, pair3_contract.functions.getReserves().call()[:2]]

    pairs = [pair, pair2, pair3]
    g = create_graph(pairs)
    cycles = g.bellman_ford(TOKEN_1_NAME)
    print(cycles)