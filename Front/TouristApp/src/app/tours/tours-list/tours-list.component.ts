import { Component, OnInit } from '@angular/core';
import { Tour, TourStatus } from '../model/tour.model';
import { ToursService } from '../tours.service';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';
import { CartService } from '../../cart/cart.service';

@Component({
  selector: 'app-tours-list',
  templateUrl: './tours-list.component.html',
  styleUrls: ['./tours-list.component.css']
})
export class ToursListComponent implements OnInit {


  tours: Tour[]=[];
  toursWithKeypoints: { tour: Tour, firstKeypoint?: any }[] = [];

  loading=true;
  isAllTours = false;

  constructor(private tourService: ToursService,private cartService: CartService, private authService: AuthService, private snackBar:MatSnackBar,private router: Router){}

  loggedUser = this.authService.getCurrentUser();
  purchasedTourIds: string[] = [];
  cartTourIds: string[] = [];

  ngOnInit(): void {
    const user = this.authService.getCurrentUser();

    if (!user || !user.userId) {
      this.snackBar.open("You must be logged in to see tours", "Close", {
        duration: 3000,
        horizontalPosition: "center"
      });
      return;
    }

    if (this.router.url.includes('allTours')) {
      this.loadAllTours();
      this.isAllTours = true;
    } else {
      this.loadMyTours(user.userId);
    }

    this.loadCartTours(user.userId.toString());
    this.loadPurchasedTours(user.userId.toString());
  }
  loadMyTours(userId: number): void {
    this.tourService.getAllTours().subscribe({
      next: (data) => {
        this.tours = data.filter(tour => Number(tour.creatorID) === Number(userId));
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading my tours', err);
        this.loading = false;
      }
    });
  }

loadAllTours(): void {
  this.tourService.getAllTours().subscribe({
    next: (data) => {
      // Filtriraj samo published ture
      const publishedTours = data.filter(tour => tour.status === "Published");
      
      // Učitaj prvi keypoint za svaku turu
      const tourRequests = publishedTours.map(tour => 
        this.tourService.getKeyPointsForTour(tour.id!).toPromise()
          .then(keypoints => ({
            tour: tour,
            firstKeypoint: keypoints && keypoints.length > 0 ? keypoints[0] : null
          }))
          .catch(() => ({
            tour: tour,
            firstKeypoint: null
          }))
      );

      Promise.all(tourRequests).then(results => {
        this.toursWithKeypoints = results;
        this.tours = publishedTours;
        this.loading = false;
      });
    },
    error: (err) => {
      console.error('Error loading all tours', err);
      this.loading = false;
    }
  });
}


  addReview(tourId: string | undefined): void {
    if (!tourId) {
      console.error("Tour ID is missing");
      return;
    }
    this.router.navigate(['/addReview', tourId]);
  }
  seeReview(tourId: string | undefined): void {
    if (!tourId) {
      console.error("Tour ID is missing");
      return;
    }
    this.router.navigate(['/tours', tourId, 'reviews']);
  }

  tourDetails(tourId: string | undefined) {
    this.router.navigate(['/tours', tourId]);
  }
  
    startTour(tour: Tour): void {
  // Ovde pozoveš startTour i otvoriš novu rutu
  this.tourService.startTour(tour.id!, 123, 44.8176, 20.4569).subscribe({
    next: (execution) => {
      this.router.navigate(['/tourExecution', execution.id]); 
    },
    error: (err) => console.error(err)
  });
}

 addToCart(tour: any) {
    const userId = this.loggedUser?.userId.toString(); // kasnije povuci iz AuthService ili LocalStorage

    if(userId)
    this.cartService.addToCart(userId, tour).subscribe({
      next: (res) => {
        console.log('Tour added to cart', res);
        alert(`Tour "${tour.name}" added to cart!`);
        if (tour.id && !this.cartTourIds.includes(tour.id)) {
          this.cartTourIds.push(tour.id);
        }
      },
      error: (err) => {
        console.error('Error adding tour to cart', err);
        alert('Failed to add tour to cart.');
      }
    });
  }
  loadCartTours(userId: string) {
    this.cartService.getCart(userId).subscribe({
      next: (res) => {
        this.cartTourIds = res.items.map((i: any) => i.tourId);
      }
    });
  }

  loadPurchasedTours(userId: string) {
    this.cartService.getPurchasedTours(userId).subscribe({
      next: (tokens) => {
        this.purchasedTourIds = tokens.map((t: any) => t.tourId);
      }
    });
  }

  isDisabled(tourId: string): boolean {
    return (
      this.cartTourIds.includes(tourId) ||
      this.purchasedTourIds.includes(tourId)
    );
  }
  isLoggedTourist(): boolean{
    return this.loggedUser?.role===1
  }
}
