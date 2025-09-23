import { Component, OnInit } from '@angular/core';
import { CartService } from './cart.service';
import { AuthService } from '../auth/auth.service';

@Component({
  selector: 'app-cart',
  templateUrl: './cart.component.html',
  styleUrls: ['./cart.component.css']
})
export class CartComponent implements OnInit {
  cart: any;
  loggedUser = this.authService.getCurrentUser();

  constructor(
    private cartService: CartService,
    private authService: AuthService
  ) {}

  ngOnInit() {
    if (this.loggedUser) {
      this.loadCart(this.loggedUser.userId.toString());
    }
  }

  loadCart(userId: string) {
    this.cartService.getCart(userId).subscribe({
      next: (res) => {
        this.cart = res;
      },
      error: (err) => {
        console.error('Error loading cart', err);
      }
    });
  }
  removeFromCart(tourId: string) {
    this.cartService.removeFromCart(this.loggedUser?.userId.toString()!, tourId).subscribe({
        next: (res) => {
        this.cart = res; // backend vraća ažuriran cart
        },
        error: (err) => {
        console.error("Error removing item", err);
        }
    });
    }
    checkout() {
    if (!this.loggedUser) return;

    this.cartService.checkout(this.loggedUser.userId.toString()).subscribe({
          next: (res) => {
          console.log('Checkout successful', res);
          alert('Checkout successful! Tokens generated.');
          if(this.loggedUser)
          this.loadCart(this.loggedUser.userId.toString()); 
          },
          error: (err) => {
          console.error('Error during checkout', err);
          alert('Checkout failed.');
          }
      });
      }
}
