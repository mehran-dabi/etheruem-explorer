# Ethereum Explorer

A simple ethereum block chain explorer

## Structure

There are two microservices inside this project. `Indexer` and `REST`.

The `Indexer` service is in charge of scanning the blockchain and storing the blocks and transactions in the database.

The `REST` service includes a couple of REST APIs to be able to explore the blockchain data stored by the `Indexer`

## How to Run

To run the services you have to first run the following command:

```bash
go mod vendor
```

Then two run the services you have to run these two commands on separate terminals:
```bash
make indexer
make rest
```

After that you can access the APIs. A Postman collection is accessible [here](https://www.getpostman.com/collections/9f84d94ac1477db3a7d0) which contains all the APIs and examples on how to use them and how the response would look like.