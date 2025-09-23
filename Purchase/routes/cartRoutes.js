const express = require('express');
const router = express.Router();
const cartController = require('../controllers/cartController');

router.post('/:userId/add', cartController.addItem);
router.post('/:userId/remove', cartController.removeItem);
router.post('/:userId/checkout', cartController.checkout);
router.get('/:userId', cartController.getCart);
router.get('/:userId/tokens', cartController.getPurchasedTours);


module.exports = router;
