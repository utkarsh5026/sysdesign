# Load Balancing Algorithms



## 1. Round Robin 🤗

Round Robin is one of the simplest load balancing algorithms. It works by cycling through a list of servers and sending each new request to the next server in the list.

**How it works:**

1. Maintain a list of servers.
2. For each incoming request, select the next server in the list.
3. After reaching the end of the list, start over from the beginning.

**Advantages:**

* Simple to implement and understand
* Fair distribution if all servers have equal capacity

**Disadvantages:**

* Doesn't consider server load or capacity
* Can lead to uneven distribution if requests vary significantly in processing time

**Use case:** Good for situations where servers have similar specifications and the requests are generally uniform in terms of processing requirements.


## 2. Weighted Round Robin 🏋️

This is an extension of the Round Robin algorithm that allows for servers with different capacities.

**How it works:**

1. Assign a weight to each server based on its capacity.
2. Distribute requests in proportion to the weights.

For example, if Server A has weight 3 and Server B has weight 1, Server A will receive 3 requests for every 1 request sent to Server B.


**Advantages:**

* Allows for heterogeneous server environments
* Can distribute load more evenly based on server capacity

**Disadvantages:**

* Requires manual assignment and maintenance of weights
* Still doesn't account for real-time server load

**Use case:** Suitable when you have servers with known, different capacities.



## 3. Least Connection 🤔

This algorithm directs traffic to the server with the fewest active connections.

**How it works:**

1. Maintain a count of active connections for each server.
2. When a new request arrives, send it to the server with the lowest number of active connections.

**Advantages:**

* Adapts to server load in real-time
* Works well for long-lived connections

**Disadvantages:**

* Doesn't consider the processing power required by different connections
* May not be optimal for short-lived connections due to overhead

**Use case:** Effective in environments where connections are long-lived and processing requirements are relatively uniform across requests.



## 4. Weighted Least Connection 🥺

This combines the Least Connection method with server weighting.

**How it works:**

1. Assign weights to servers based on their capacity.
2. For each server, calculate: (active connections) / (weight)
3. Send new requests to the server with the lowest result from step 2.

**Advantages:**

* Considers both current load and server capacity
* More adaptive than simple Weighted Round Robin

**Disadvantages:**

* Requires accurate weight assignment
* More complex to implement than simpler methods

**Use case:** Useful in heterogeneous environments where connections can be long-lived and server capacities differ.



## 5. Least Response Time ⏩

This algorithm selects the server with the fastest response time to handle new requests.

**How it works:**

1. Monitor the response time of each server.
2. Send new requests to the server with the lowest response time.

**Advantages:**

* Can lead to better overall system performance
* Adapts to both server load and network conditions

**Disadvantages:**

* Requires continuous monitoring of server response times
* Can be affected by momentary spikes in response time

**Use case:** Suitable for applications where minimizing response time is critical, such as real-time applications or APIs.



## 6. Resource-Based (Adaptive) 💡

This algorithm distributes load based on real-time analysis of server resources.

**How it works:**

1. Continuously monitor server resources (CPU, memory, disk I/O, network bandwidth).
2. Assign a "load score" to each server based on its current resource usage.
3. Distribute incoming requests to the server with the lowest load score.

**Advantages:**

* Highly adaptive to current server conditions
* Can prevent overloading of specific resources

**Disadvantages:**

* Requires sophisticated monitoring and analysis
* Can lead to rapid shifts in traffic patterns

**Use case:** Ideal for environments with varying request types that stress different system



## 7. IP Hash 🤷‍♂️

This algorithm uses the client's IP address to determine which server should receive the request.

**How it works:**

1. Calculate a hash value from the client's IP address.
2. Use this hash to select a server (often by modulo with the number of servers).

**Advantages:**

* Ensures requests from the same client always go to the same server (session persistence)
* Useful for applications that require client-server affinity

**Disadvantages:**

* Can lead to uneven distribution if client IP addresses are not well distributed
* Doesn't adapt to changing server loads

**Use case:** Suitable for applications where maintaining session state on the same server is important.


## 8. URL Hash 🥹

Similar to IP Hash, but uses the requested URL to determine the server.

**How it works:**

1. Calculate a hash value from the requested URL.
2. Use this hash to select a server.

**Advantages:**

* Ensures the same content is always served by the same server
* Useful for caching scenarios

**Disadvantages:**

* Can lead to uneven distribution if URL popularity varies significantly
* Doesn't adapt to server load

**Use case:** Effective for content delivery networks (CDNs) and caching systems.


## 9. Random 😂

This algorithm randomly selects a server for each incoming request.

**How it works:**

1. For each new request, randomly select a server from the pool.

**Advantages:**

* Very simple to implement
* No state needs to be maintained

**Disadvantages:**

* Can lead to uneven distribution, especially with a small number of requests
* Doesn't consider server load or capacity

**Use case:** Can be useful in situations where simplicity is valued over optimal distribution, or as part of more complex algorithms.


## 10. Least Bandwidth 😯

This method selects the server currently serving the least amount of traffic (measured in Mbps).

**How it works:**

1. Monitor the outgoing traffic for each server.
2. Direct new requests to the server with the lowest current bandwidth usage.

**Advantages:**

* Can prevent network bottlenecks
* Useful when network capacity is a primary concern

**Disadvantages:**

* Doesn't consider other resources like CPU or memory
* May not be optimal if bandwidth usage doesn't correlate with server load

**Use case:** Suitable for applications where network traffic is the primary bottleneck.


## 11. Least Packets 😎

This algorithm directs traffic to the server that has received the fewest packets over a given time period.

**How it works:**

1. Count the number of packets each server has received over a recent time window.
2. Send new requests to the server with the lowest packet count.

**Advantages:**

* Can be a good proxy for overall server load in some scenarios
* Simple to implement and understand

**Disadvantages:**

* Packet count may not always correlate with actual server load
* Doesn't consider packet size or processing requirements

**Use case:** Can be effective in environments where the number of requests (rather than size or complexity) is the primary concern.



## 12. Source IP Hash 👋

This algorithm is similar to IP Hash but focuses solely on the source IP address.

**How it works:**

1. Extract the source IP address from the incoming request.
2. Calculate a hash value from the source IP address.
3. Use this hash to select a server (often by modulo with the number of servers).

**Advantages:**

* Ensures requests from the same source IP always go to the same server
* Useful when client IP might be masked by intermediate proxies

**Disadvantages:**

* Can lead to uneven distribution if source IP addresses are not well distributed
* Doesn't adapt to changing server loads

**Use case:** Suitable for maintaining session persistence in scenarios where the full client IP might be obscured, such as when dealing with traffic from large corporate networks or mobile carriers.



## 14. Chained Failover 🔗

This method isn't strictly a load balancing algorithm, but rather a high-availability technique often used in conjunction with load balancing.

**How it works:**

1. Define a prioritized list of servers.
2. Send all requests to the highest-priority server.
3. If that server fails, automatically redirect all traffic to the next server in the list.

**Advantages:**

* Simple to understand and implement
* Provides a clear failover path

**Disadvantages:**

* Doesn't distribute load under normal conditions
* Can lead to sudden, large shifts in load during failover

**Use case:** Useful in scenarios where you have a preferred server and one or more backup servers, and where the primary concern is high availability rather than load distribution.
