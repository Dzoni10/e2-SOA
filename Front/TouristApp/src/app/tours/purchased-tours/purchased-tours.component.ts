import { Component, OnInit } from '@angular/core';
import { Tour } from '../model/tour.model';
import { ToursService } from '../tours.service';
import { CartService } from '../../cart/cart.service';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

@Component({
  selector: 'app-purchased-tours',
  templateUrl: './purchased-tours.component.html',
  styleUrls: ['./purchased-tours.component.css']
})
export class PurchasedToursComponent implements OnInit {
  tours: Tour[] = [];
  toursWithKeypoints: { tour: Tour, firstKeypoint?: any }[] = [];
  loading = true;

  loggedUser = this.authService.getCurrentUser();

  constructor(
    private toursService: ToursService,
    private cartService: CartService,
    private authService: AuthService,
    private snackBar: MatSnackBar,
    private router: Router
  ) {}

  ngOnInit(): void {
    if (!this.loggedUser?.userId) {
      this.snackBar.open("You must be logged in to see purchased tours", "Close", {
        duration: 3000,
        horizontalPosition: "center"
      });
      return;
    } else {
      this.loadPurchasedTours(this.loggedUser.userId.toString());
    }
  }

  loadPurchasedTours(userId: string) {
    this.cartService.getPurchasedTours(userId).subscribe({
      next: (tokens) => {
        const requests = tokens.map((token: any) =>
          this.toursService.getTourByID(token.tourId).toPromise()
            .then(tour => {
              if (!tour) return null;
              return this.toursService.getKeyPointsForTour(tour.id!).toPromise()
                .then(keypoints => ({
                  tour,
                  firstKeypoint: keypoints && keypoints.length > 0 ? keypoints[0] : null
                }))
                .catch(() => ({ tour, firstKeypoint: null }));
            })
        );

        Promise.all(requests).then((results) => {
          this.toursWithKeypoints = results.filter(r => !!r) as { tour: Tour, firstKeypoint?: any }[];
          this.tours = this.toursWithKeypoints.map(r => r.tour);
          this.loading = false;
        });
      },
      error: (err) => {
        console.error('Error loading purchased tours', err);
        this.loading = false;
      }
    });
  }

  // 👇 sve funkcije iz ToursListComponent, osim addToCart
  addReview(tourId: string | undefined): void {
    if (!tourId) return;
    this.router.navigate(['/addReview', tourId]);
  }

  seeReview(tourId: string | undefined): void {
    if (!tourId) return;
    this.router.navigate(['/tours', tourId, 'reviews']);
  }

  tourDetails(tourId: string | undefined) {
    if (!tourId) return;
    this.router.navigate(['/tours', tourId]);
  }

  startTour(tour: Tour): void {
    this.toursService.startTour(tour.id!, 123, 44.8176, 20.4569).subscribe({
      next: (execution) => {
        this.router.navigate(['/tourExecution', execution.id]);
      },
      error: (err) => console.error(err)
    });
  }
}
