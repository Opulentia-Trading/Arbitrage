# Inspiration: https://github.com/benleim/pathfinder

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

    for ctr in range(number_of_vrtx - 1):
      for u, v, w in self.edges:
        if self.dist[v] > self.dist[u] + w:
          self.dist[v] = self.dist[u] + w
          self.previous_v[v] = u


    cyclePaths = []
    foundCycles = {}
    for ctr in range(number_of_vrtx - 1):
      for u, v, w in self.edges:
        if self.dist[v] > self.dist[u] + w:
            cyclePath = {}

            curr = u
            index = 1
            cyclePath[u] = index
            index += 1

            while not cyclePath[curr]:
                cyclePath[curr.value] = index
                index += 1
                curr = self.previous_v[curr]
            
            cyclePath[curr.value+'_'] = index
            
            path = []
            for k in cyclePath.keys():
                path += (k.replace('_',''))
            path.reverse()

            for i in range(len(path)):
                if (i != 0 and path[0] == path[i]):
                    path = path[0, i+1]
                    break

            # TODO: ENSURE UNIQUENESS OF CYCLE
            cyclePaths += path
    return cyclePaths

def create_graph(pairs):
    g = Graph()
    for p in pairs:
        g.add_edge(p["name1"], p["name2"], p["price"])
    return g


# Pairs assumed to be in the form: (name1, name2, value between the 2)
def cyclic_arb(token_pairs):
    g = create_graph(token_pairs)
    results = g.bellman_ford(token_pairs[0]["name1"])
    # Calculating cycle scrapped since I'll be creating token vertex data types now.
    # This cycle is essentially enough to know what transactions to make
    # Calculate gas fee as well and we should be good