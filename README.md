This program was fun to write. I thought I would write up a quick summary to explain my thought process.

# Usage
1. Installation:

```
go mod init test/blockchain
go install test/blockchain
go build
```

2. Running the program:

```
./blockchain
```

# Assumptions

- This program scans every single block from the latest to the earliest per address, so it is very primitive and not optimized.
- If there is a new block, this program still scans the older blocks again. To improve the speed and not do repetitive work, a cache needs to be implemented to not rescan older blocks. (See Caveats)


# Caveats
Below are some caveats of the code:

- Limited Error Handling: While the code handles some errors (e.g., retrieving block number, empty receiver address), it could benefit from more comprehensive error handling. Consider including checks for nil pointers, transaction parsing errors, and potential issues during subscription.

- Repetitive Block Processing: The GetTransactions function iterates through all blocks from the latest parsed block downwards. This can be inefficient for frequently updated addresses, as it reprocesses previously seen blocks. Consider implementing a mechanism to track processed transactions by hash or block number to avoid redundant processing.

- Limited Information: The Transaction struct only captures basic transaction details (from, to, value, hash, block number). I'm not sure exactly what to retrieve for the transaction so I just retrieve the basic fields. Additional information like transaction type (e.g., ERC-20 transfer), gas price, or gas used can also be included.


# Future Improvements and Productionalization

Without knowing how this application is used, i.e. how many blocks we want to scan, how many addresses there will be, to improve this program a few things come to my mind:

- Event Monitoring: Implement a mechanism to listen for new block events from the Ethereum node. This allows capturing transactions as they occur without manually polling for updates. By actively querying the block events, we reduce future latency in case many users query for transactions at once.

- Filtering: Allow filtering transactions based on specific criteria (e.g., transaction type, minimum value). This can be helpful for focusing on relevant transactions.

- Caching: Implement a caching mechanism to store retrieved transactions and reduce the load on the Ethereum node for frequently accessed addresses.

- Concurrency: Explore using concurrency for tasks like block processing and transaction retrieval to improve performance, especially when dealing with a large number of addresses.

- Logging: Enhance logging to capture more details about errors, warnings, and program execution for better debugging and monitoring.

- Single node dependency: we can also consider subscribing to different Eth node in case one node goes down. Maybe a node pool can be used for reliability and load balancing (I was afraid of being throttled while writing this program but luckily I wasn't!).



