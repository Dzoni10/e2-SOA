const mongoose = require('mongoose');

const OrderItemSchema = new mongoose.Schema({
  tourId: { type: String, required: true },
  name: { type: String, required: true },
  price: { type: Number, required: true }
});

const CartSchema = new mongoose.Schema({
  userId: { type: String, required: true },
  items: [OrderItemSchema],
  totalPrice: { type: Number, default: 0 }
});

CartSchema.pre('save', function (next) {
  this.totalPrice = this.items.reduce((sum, item) => sum + item.price, 0);
  next();
});

module.exports = mongoose.model('Cart', CartSchema);
