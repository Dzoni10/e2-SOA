const mongoose = require('mongoose');

const TokenSchema = new mongoose.Schema({
  userId: { type: String, required: true },
  tourId: { type: String, required: true },
  issuedAt: { type: Date, default: Date.now }
});

module.exports = mongoose.model('Token', TokenSchema);