
# Word of wisdom TCP server

TCP server with protection from DDOS based on Proof of Work


## Installation

Start server and client by docker compose

```bash
  make start
```
    
## Types of request

- Quit signal to close connection
- RequestChallenge request new challenge from server
- ResponseChallenge message with challenge for client
- RequestResource message with solved challenge
- ResponseResource message with useful information solution is correct or not


## Protocol

The solution uses TCP based protocol
header - integer to indicate which type of request was sent
payload - string that will be json depends on type of request

## Choice of the POW algorithm

Hashcash is considered to be a strong PoW algorithm due to its efficiency, security, flexibility, and popularity

Efficiency: Hashcash is a computationally efficient algorithm that uses a cryptographic hash function to generate a hash value that meets a specific criteria, and this can be done quickly and efficiently.

Security: Hashcash is a secure algorithm that makes it difficult for attackers to manipulate the system.
The algorithm requires the sender to perform a certain amount of computational work to generate a valid hash value, 
and this makes it difficult for attackers to generate large  cryptocurrency transactions.

Flexibility: Hashcash is a flexible algorithm that can be used in a wide range of applications.

Popularity: Hashcash has been widely adopted and is used in many different applications, including email and cryptocurrency.

## Disadvantages of other protocol
Merkle tree
Additional computational overhead
Potential centralization
Inefficiency for dynamic data
Limited scalability
Security risks if the Merkle root is compromised.

Guided tour puzzles may not be compatible with all devices and platforms, which can limit their effectiveness and accessibility for some users
