const { MongoClient } = require("mongodb");
require("dotenv").config();

const username = process.env.MONGO_USERNAME;
const password = process.env.MONGO_PASSWORD;
const cluster = process.env.MONGO_CLUSTERNAME;

const url = `mongodb+srv://${username}:${password}@${cluster}.i8qidli.mongodb.net/?retryWrites=true&w=majority&appName=${cluster}`;
const client = new MongoClient(url);

async function run() {
  try {
    await client.connect();
    const db = client.db("sample_guides");
    const collection = db.collection("planets");
    // const indexes = collection.listIndexes().toArray();

    // ---------------------------- INSERT DOCUMENTS --------------------------------- //
    // const newDocs = [
    //   {
    //     name: "Exoplanet-29",
    //     orderFromSun: 10,
    //     hasRings: false,
    //     mainAtmosphere: [],
    //     surfaceTemperatureC: {
    //       min: -20,
    //       max: 100,
    //       mean: 34,
    //     },
    //   },
    //   {
    //     name: "Exoplanet-30",
    //     orderFromSun: 11,
    //     hasRings: true,
    //     mainAtmosphere: [],
    //     surfaceTemperatureC: {
    //       min: -20,
    //       max: 100,
    //       mean: 34,
    //     },
    //   },
    // ];

    // const addDocuments = await collection.insertMany(newDocs);
    // console.log(
    //   `${addDocuments.insertedCount} documents successfully inserted.`
    // );

    // ---------------------------- FIND DOCUMENTS --------------------------------- //
    const findQuery = { "surfaceTemperatureC.min": { $eq: -20 } };
    const documents = await collection
      .find(findQuery)
      .sort({ name: 1 })
      .toArray();
    // console.log(`Found documents: ${JSON.stringify(documents)}`);

    let docIdArray = [];

    for (const document of documents) {
      const documentId = document._id;
      docIdArray.push(documentId);
    }
    // console.log("Array of document ids:", docIdArray);

    const specificDoc = documents.find((doc) => doc.orderFromSun === 11);
    if (specificDoc) {
      const docId = specificDoc._id;
      // console.log(`Found specific document: ${docId}`);
    } else {
      console.log("Could not find planet at 11th position from the Sun");
    }

    // ---------------------------- UPDATE DOCUMENTS --------------------------------- //
    const updateDocument = {
      $set: { mainAtmosphere: ["Oxygen", "Nitrogen", "Ammonia"] },
    };
    const updateOptions = {};
    const updateResult = await collection.findOneAndUpdate(
      findQuery,
      updateDocument,
      updateOptions
    );

    console.log(`Updated document: ${JSON.stringify(updateResult)}`);

    // ---------------------------- DELETE DOCUMENTS --------------------------------- //
    const deleteResult = await collection.deleteOne({ _id: specificDoc._id });
    console.log(`Deleted ${deleteResult.deletedCount} document.`);
  } catch (err) {
    console.error(err);
  }
}

run().catch(console.error);
