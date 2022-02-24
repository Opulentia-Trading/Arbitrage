# 2 pairs of (base, quote) from platform 1 and 2
# We essentially use "base" as the anchor of our arbitrage
# We buy x amount of quote token on the platform with the lower quote token price
# Then we exchange it in base tokens of amount y on platform 2 (as quote token is valued more on platform 2)
# We then give z amount of base token back to platform 1 since we only burrowed from there
# Altogether, we have a profit of y-z
 
from math import sqrt

# (BASE1, QUOTE1)
# (BASE2, QUOTE2)
# Conditions QUOTE2 < QUOTE1

# This returns the amount needed to be borrowed in terms of quote token 1
def calculate_arb(platform_1, platform_2):
    quote_1, base_1 = platform_1
    quote_2, base_2 = platform_2
    
    # We want to use the quadratic equation to maximize the potential profit

    a = quote_1*base_1 - quote_2*base_2
    b = 2*base_1*base_2*(quote_1 + quote_2)
    c = base_1*base_2*(quote_1*base_2 - quote_2*base_1)

    # Prove ax^2 + bx + c = 0
    m = b**2 - 4*a*c

    # We can't square root negatives
    if (m < 0):
        return -1

    # Our 2 possible solutions are here
    x_1 = (-b + sqrt(m))/(2*a)
    x_2 = (-b - sqrt(m))/(2*a)

    check_1 = x_1 > 0 and x_1 < base_1 and x_1 < base_2
    check_2 = x_2 > 0 and x_2 < base_2 and x_2 < base_2

    if (check_1 or check_2):
        return -1

    return x_1 if check_1 else x_2

if __name__ == "__main__":
    print(calculate_arb((10, 8), (10, 9)))