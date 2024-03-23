const { MongoClient } = require("mongodb");
require("dotenv").config();

const username = process.env.MONGO_USERNAME;
const password = process.env.MONGO_PASSWORD;
const cluster = process.env.MONGO_CLUSTERNAME;

// const url = `mongodb+srv://${username}:${password}@${cluster}.mongodb.net/?retryWrites=true&w=majority`;

// Had to downgrade to node 2.12 or later conn str - not sure why the latest node version doesn't work
const url = `mongodb://${username}:${password}@ac-11nsdvj-shard-00-00.i8qidli.mongodb.net:27017,ac-11nsdvj-shard-00-01.i8qidli.mongodb.net:27017,ac-11nsdvj-shard-00-02.i8qidli.mongodb.net:27017/?ssl=true&replicaSet=atlas-27vp6s-shard-0&authSource=admin&retryWrites=true&w=majority&appName=${cluster}`;

// Connect to your Atlas cluster
const client = new MongoClient(url);
async function run() {
  try {
    await client.connect();
    console.log("Successfully connected to Atlas");
  } catch (err) {
    console.log(err.stack);
  } finally {
    await client.close();
  }
}
run().catch(console.dir);
