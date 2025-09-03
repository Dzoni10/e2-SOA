import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { ToursService } from '../tours.service';
import { Review } from '../model/review.model';

@Component({
  selector: 'app-tour-reviews',
  templateUrl: './tour-review.component.html',
  styleUrls: ['./tour-review.component.css']
})
export class TourReviewComponent implements OnInit {

  reviews: Review[] = [];
  tourId: string = '';
  loading = true;

  constructor(
    private route: ActivatedRoute,
    private toursService: ToursService
  ) {}

  ngOnInit(): void {
    this.tourId = this.route.snapshot.paramMap.get('id') || '';
    if (this.tourId) {
      this.toursService.getReviewsForTour(this.tourId).subscribe({
        next: (data) => {
          this.reviews = data;
          this.loading = false;
        },
        error: (err) => {
          console.error('Error fetching reviews', err);
          this.loading = false;
        }
      });
    }
  }
}
