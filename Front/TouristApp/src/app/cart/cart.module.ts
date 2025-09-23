import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CartComponent } from './cart.component';
import { RouterModule } from '@angular/router';
import { HttpClientModule } from '@angular/common/http';

@NgModule({
  declarations: [
    CartComponent
  ],
  imports: [
    CommonModule,         // 👈 zbog *ngIf, *ngFor
    HttpClientModule,     // 👈 da radi HttpClient u CartService
    RouterModule
  ],
  exports: [
    CartComponent
  ]
})
export class CartModule {}
