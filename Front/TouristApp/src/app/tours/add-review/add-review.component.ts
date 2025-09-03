import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { ToursService } from '../tours.service';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Review } from '../model/review.model';
import { StakeholdersService } from '../../stakeholders/stakeholders.service'; // <-- dodato

@Component({
  selector: 'app-add-review',
  templateUrl: './add-review.component.html',
  styleUrls: ['./add-review.component.css']
})
export class AddReviewComponent implements OnInit {

  review: Review = {
    tourId: '',
    userId: 0,
    username: '',
    rating: 1,
    comment: '',
    image: ''
  };
  stars = [1, 2, 3, 4, 5];

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private tourService: ToursService,
    private authService: AuthService,
    private stakeholdersService: StakeholdersService, // <-- dodato
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    const tourId = this.route.snapshot.paramMap.get('id');
    if (tourId) {
      this.review.tourId = tourId;
    }

    const user = this.authService.getCurrentUser();
    if (user) {
      this.review.userId = user.userId;
      this.stakeholdersService.getUserById(user.userId).subscribe({
        next: (data) => {
          this.review.username = data.username; 
        },
        error: (err) => {
          console.error('Error fetching user details', err);
          this.snackBar.open('Could not load user info', 'Close', { duration: 3000 });
        }
      });
    }
  }

  submitReview(): void {
    this.tourService.addReview(this.review).subscribe({
      next: () => {
        this.snackBar.open('Review added successfully!', 'Close', { duration: 3000 });
        this.router.navigate(['/allTours']);
      },
      error: (err) => {
        console.error('Error adding review', err);
        this.snackBar.open('Failed to add review', 'Close', { duration: 3000 });
      }
    });
  }

   setRating(value: number): void {
        this.review.rating = value;
    }
}
