const cartRepo = require('../repositories/cartRepository');
const tokenRepo = require('../repositories/tokenRepository');

async function addItem(userId, item) {
  let cart = await cartRepo.findByUser(userId);
  if (!cart) cart = await cartRepo.create(userId);

  cart.items.push(item);
  await cartRepo.save(cart);
  return cart;
}

async function removeItem(userId, tourId) {
  let cart = await cartRepo.findByUser(userId);
  if (!cart) throw new Error('Cart not found');

  cart.items = cart.items.filter(i => i.tourId !== tourId);
  await cartRepo.save(cart);
  return cart;
}

async function checkout(userId) {
  let cart = await cartRepo.findByUser(userId);
  if (!cart) throw new Error('Cart not found');

  const tokens = await tokenRepo.createMany(
    cart.items.map(i => ({ userId, tourId: i.tourId }))
  );

  cart.items = [];
  cart.totalPrice = 0;
  await cartRepo.save(cart);

  return tokens;
}

async function getCart(userId) {
  return (await cartRepo.findByUser(userId)) || { userId, items: [], totalPrice: 0 };
}

async function getPurchasedTours(userId) {
  return await tokenRepo.findByUser(userId);
}

module.exports = { addItem, removeItem, checkout, getCart, getPurchasedTours };
