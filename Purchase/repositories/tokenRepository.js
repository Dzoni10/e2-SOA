const Token = require('../models/Token');

async function createMany(tokens) {
  return Token.insertMany(tokens);
}
async function findByUser(userId) {
  return Token.find({ userId });
}

module.exports = { createMany, findByUser };
