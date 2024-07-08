# Detailed Example of Weighted Round Robin Algorithm

Let's consider a scenario with three servers:

- Server A: Weight 5
- Server B: Weight 3
- Server C: Weight 2

## Initial Setup

```go
wrr := WeightedRoundRobinBalancer(3)
wrr.AddServer(Server{Weight: 5, Id: "A"})
wrr.AddServer(Server{Weight: 3, Id: "B"})
wrr.AddServer(Server{Weight: 2, Id: "C"})
```

After adding these servers:

- `MaxWeight` = 5 (highest weight)
- `CurrentWeight` = 5 (initialized to MaxWeight)
- `current` = -1 (index of the last selected server, starts at -1)

## Selection Process

Let's go through several rounds of server selection:

1. **First call to NextServer()**

   - `current` becomes 0 (A)
   - 5 >= 5, so A is selected
   - Return: A
2. **Second call**

   - `current` becomes 1 (B)
   - 3 < 5, move to next
   - `current` becomes 2 (C)
   - 2 < 5, move to next
   - `current` becomes 0 (A)
   - `CurrentWeight` decreases to 4
   - 5 >= 4, so A is selected
   - Return: A
3. **Third call**

   - `current` becomes 1 (B)
   - 3 < 4, move to next
   - `current` becomes 2 (C)
   - 2 < 4, move to next
   - `current` becomes 0 (A)
   - `CurrentWeight` decreases to 3
   - 5 >= 3, so A is selected
   - Return: A
4. **Fourth call**

   - `current` becomes 1 (B)
   - 3 >= 3, so B is selected
   - Return: B
5. **Fifth call**

   - `current` becomes 2 (C)
   - 2 < 3, move to next
   - `current` becomes 0 (A)
   - `CurrentWeight` decreases to 2
   - 5 >= 2, so A is selected
   - Return: A
6. **Sixth call**

   - `current` becomes 1 (B)
   - 3 >= 2, so B is selected
   - Return: B
7. **Seventh call**

   - `current` becomes 2 (C)
   - 2 >= 2, so C is selected
   - Return: C
8. **Eighth call**

   - `current` becomes 0 (A)
   - `CurrentWeight` decreases to 1
   - 5 >= 1, so A is selected
   - Return: A
9. **Ninth call**

   - `current` becomes 1 (B)
   - 3 >= 1, so B is selected
   - Return: B
10. **Tenth call**

    - `current` becomes 2 (C)
    - 2 >= 1, so C is selected
    - Return: C
11. **Eleventh call**

    - `current` becomes 0 (A)
    - `CurrentWeight` would decrease to 0, so it resets to 5 (MaxWeight)
    - 5 >= 5, so A is selected
    - Return: A

## Analysis

After these 11 selections, the distribution is:

- A: 5 times
- B: 3 times
- C: 2 times

This perfectly matches the weight ratios of the servers (5:3:2).

## Key Observations

1. Higher-weight servers (like A) are selected more frequently.
2. All servers get selected, preventing starvation of lower-weight servers.
3. The algorithm cycles through all servers before resetting `CurrentWeight`.
4. The distribution of selections closely follows the weight ratios over time.
5. The selection order maintains a degree of round-robin behavior while respecting weights.

This example demonstrates how the Weighted Round Robin algorithm balances the load across servers of different capacities, ensuring a fair distribution based on their assigned weights.
