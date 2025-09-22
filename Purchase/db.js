const mongoose = require('mongoose');

async function connectDB() {
  const uri = process.env.MONGO_URI || "mongodb://purchase-db:27017/purchase";

  try {
    await mongoose.connect(uri, {
      useNewUrlParser: true,
      useUnifiedTopology: true
    });

    console.log("✅ MongoDB connected:", uri);
  } catch (err) {
    console.error("❌ Cannot connect to MongoDB:", err.message);
    process.exit(1);
  }
}

module.exports = connectDB;
