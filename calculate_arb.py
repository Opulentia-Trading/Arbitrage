# 2 pairs of (base, quote) from platform 1 and 2
# We essentially use "base" as the anchor of our arbitrage
# We buy x amount of quote token on the platform with the lower quote token price
# Then we exchange it in base tokens of amount y on platform 2 (as quote token is valued more on platform 2)
# We then give z amount of base token back to platform 1 since we only burrowed from there
# Altogether, we have a profit of y-z
 
from math import sqrt
from collections import namedtuple

# (BASE1, QUOTE1)
# (BASE2, QUOTE2)
# Conditions QUOTE2 < QUOTE1

from dataclasses import dataclass


@dataclass
class Platform:
    name: str
    tokens: list

@dataclass
class Tokens:
    reserves: int
    name: str


# FILL IN HERE
token1_1 = Tokens(, "eth")
token2_1 = Tokens(, "btc")
token1_2 = Tokens(, "eth")
token2_2 = Tokens(, "btc")

platform1 = Platform("PancakeSwap", [token1_1, token2_1])
platform2 = Platform("SushiSwap", [token1_2, token2_2])


# This returns the amount needed to be borrowed in terms of quote token 1
def calculate_arb(platform_1, platform_2):
    pool0token0, pool0token1 = platform_1.tokens
    pool1token0, pool1token1 = platform_2.tokens
    
    # The amount we will be buying on platform 1 denominated in token 2
    amount_to_buy = calc_max_to_buy(pool0token0.reserves, pool0token1.reserves, pool1token0.reserves, pool1token1.reserves)

    # The amount we will be buying on platform 1 denominated in token 1
    amount_to_pay_back = amount_in(amount_to_buy, pool0token0.reserves, pool0token1.reserves)

    # The amount we will be receiving on platform 2 once we sell the token 2 on platform 2
    amount_received = amount_out(amount_to_buy, pool1token1.reserves, pool1token0.reserves)

    if (amount_received < amount_to_pay_back):
        return None
    
    platform2_struct = create_struct("sell", amount_to_buy, pool1token0.name, pool0token1.name, platform_1.name, None)
    platform1_struct = create_struct("buy", amount_to_buy, None, pool0token1.name, platform_1.name, platform2_struct)

    return platform1_struct

def create_struct(action, amount, currency_received, token, platform, next):
    struct = {"Action": action, "Amount": amount, "To Receive": currency_received, "Token": token, "Platform": platform, "Action": next}
    return struct

def calc_max_to_buy(p0t0, p0t1, p1t0, p1t1):

    # Formulate a quadratic equation with their reserve amounts
    term1 = p0t0 * p0t1 - p1t0 * p1t1
    term2 = 2 * p0t1 * p1t1 * (p0t0 + p1t0)
    term3 = p0t1 * p1t1 * (p0t0 * p1t1 - p1t0 * p0t1)

    # Now we solve for quadratic equation: term1*x^2 + term2*x + term3 = 0
    discriminant = term2**2-4*term1*term3

    if discriminant < 0:
        print ("This equation has no real solution")
        return -1
    
    # This is the amount we want to buy
    x1 = (-term2+sqrt(discriminant))/(2*term1)
    x2 = (-term2-sqrt(discriminant))/(2*term1)

    if (x1 > 0 and x1 < p0t1 and x1 < p1t1):
        return x1
    elif (x2 > 0 and x2 < p0t1 and x2 < p1t1):
        return x2
    return -1


def amount_in(to_buy, p0t0, p0t1):
    numerator = p0t0 * to_buy * 100
    denominator = p0t1 * to_buy * 997
    return numerator / denominator + 1

def amount_out(to_buy, p0t0, p0t1):
    numerator = p0t0 * to_buy * 100
    denominator = p0t1 * to_buy * 997
    return numerator / denominator + 1

if __name__ == "__main__":
    print(calculate_arb(platform1, platform2))