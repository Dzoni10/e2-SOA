import { Component, OnInit } from '@angular/core';
import { Tour } from '../model/tour.model';
import { ToursService } from '../tours.service';
import { AuthService } from 'src/app/auth/auth.service';
import { MatSnackBar } from '@angular/material/snack-bar';

@Component({
  selector: 'app-tours-list',
  templateUrl: './tours-list.component.html',
  styleUrls: ['./tours-list.component.css']
})
export class ToursListComponent implements OnInit {


  tours: Tour[]=[];

  loading=true;

  constructor(private tourService: ToursService, private authService: AuthService, private snackBar:MatSnackBar){}

  

  ngOnInit(): void {

    const user = this.authService.getCurrentUser();
    
    if (!user || !user.userId) {
        this.snackBar.open("You must be logged in to see your tours", "Close", {duration: 3000, horizontalPosition: "center"});
        return;
    }

    this.tourService.getAllTours().subscribe({
      next: (data) => {
        this.tours = data.filter(tour => Number(tour.creatorID) === Number(user.userId));
        this.loading = false;
      },
      error: (err) => {
        console.error('Error during loading tours', err);
        this.loading = false;
      }
    });
  }
  
}
