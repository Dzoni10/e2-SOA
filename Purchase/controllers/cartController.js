const cartService = require('../services/cartService');

async function addItem(req, res) {
  try {
    const cart = await cartService.addItem(req.params.userId, req.body);
    res.json(cart);
  } catch (err) {
    res.status(500).send(err.message);
  }
}

async function removeItem(req, res) {
  try {
    const cart = await cartService.removeItem(req.params.userId, req.body.tourId);
    res.json(cart);
  } catch (err) {
    res.status(500).send(err.message);
  }
}

async function checkout(req, res) {
  try {
    const tokens = await cartService.checkout(req.params.userId);
    res.json(tokens);
  } catch (err) {
    res.status(500).send(err.message);
  }
}

async function getCart(req, res) {
  try {
    const cart = await cartService.getCart(req.params.userId);
    res.json(cart);
  } catch (err) {
    res.status(500).send(err.message);
  }
}

async function getPurchasedTours(req, res) {
  try {
    const { userId } = req.params;
    const tokens = await cartService.getPurchasedTours(userId);
    res.json(tokens);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}
module.exports = { addItem, removeItem, checkout, getCart, getPurchasedTours };
