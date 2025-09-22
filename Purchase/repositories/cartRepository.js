const Cart = require('../models/Cart');

async function findByUser(userId) {
  return Cart.findOne({ userId });
}

async function create(userId) {
  return new Cart({ userId, items: [] });
}

async function save(cart) {
  return cart.save();
}

module.exports = { findByUser, create, save };
