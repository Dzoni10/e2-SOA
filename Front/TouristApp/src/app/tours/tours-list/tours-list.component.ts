import { Component, OnInit } from '@angular/core';
import { Tour } from '../model/tour.model';
import { ToursService } from '../tours.service';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Router } from '@angular/router';

@Component({
  selector: 'app-tours-list',
  templateUrl: './tours-list.component.html',
  styleUrls: ['./tours-list.component.css']
})
export class ToursListComponent implements OnInit {


  tours: Tour[]=[];

  loading=true;
  isAllTours = false;

  constructor(private tourService: ToursService, private authService: AuthService, private snackBar:MatSnackBar,private router: Router){}

  

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
        this.tours = data;
        this.loading = false;
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
  
}
