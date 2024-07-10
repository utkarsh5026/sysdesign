# 🏋️‍♂️ Weighted Least Connection Load Balancer Interview Questions 🖥️

1. 🤔 What is a weighted least connection load balancer?

A weighted least connection load balancer is a system that distributes incoming network traffic across multiple servers based on:

- The number of active connections each server currently has
- The weight (capacity) assigned to each server

It aims to send new requests to the server with the lowest ratio of active connections to its weight.

2. 💡 How does the 'weight' factor into the server selection process?

The weight represents a server's capacity or processing power. In the selection process:

- A server's load is calculated as (number of active connections) / (weight)
- The server with the lowest calculated load is chosen for the next request
- This allows more powerful servers (higher weight) to handle more connections

3. 🧮 If Server A has 10 connections and weight 5, and Server B has 15 connections and weight 10, which one would be selected for the next request?

Let's calculate:

- Server A: 10 connections / 5 weight = 2
- Server B: 15 connections / 10 weight = 1.5

Server B would be selected because it has a lower weighted connection ratio (1.5 < 2).

4. 🔧 What data structure is commonly used to efficiently implement a weighted least connection load balancer?

A min-heap (priority queue) is commonly used. This allows for:

- O(1) time complexity to find the server with the lowest weighted connection ratio
- O(log n) time complexity to add or remove a server, or update a server's load
  Where n is the number of servers.

5. 🚀 How does the load balancer handle a scenario where a server's capacity changes dynamically?

The load balancer should provide a method to update a server's weight dynamically. After updating the weight:

- The server's position in the priority queue (heap) is adjusted
- This ensures that future selections accurately reflect the server's new capacity

6. 🐞 What happens if a server's weight is set to zero?

If a server's weight is set to zero:

- It should be handled as a special case to avoid division by zero
- Typically, a server with zero weight is considered to have infinite load
- It should not be selected for new connections unless all servers have zero weight

7. 🔄 How does the load balancer maintain fairness when servers have different weights?

Fairness is maintained by:

- Always selecting the server with the lowest ratio of (connections / weight)
- This ensures that servers receive a number of connections proportional to their weight
- Over time, the load distribution will approach the ratio of the servers' weights

8. 📊 How would you monitor the effectiveness of a weighted least connection load balancer?

You could monitor:

- The distribution of connections across servers compared to their weights
- Response times from each server
- The frequency of weight adjustments needed
- Overall system throughput and latency
- Any occurrences of server overload despite the balancing
