const express = require('express');
const connectDB = require('./db');
const cartRoutes = require('./routes/cartRoutes');

const app = express();
app.use(express.json());

// Mongo konekcija
connectDB();

app.use('/purchase/cart', cartRoutes);

const PORT = process.env.PORT || 5000;
app.listen(PORT, () => console.log(`Purchase-service running on port ${PORT}`));
