import { Component, OnInit,signal } from '@angular/core';
import { FormBuilder, FormGroup, Validators,FormControl } from '@angular/forms';
import { ToursService } from '../tours.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Tour } from '../model/tour.model';
import { MatChipInputEvent } from '@angular/material/chips';
import { AuthService } from 'src/app/auth/auth.service';

@Component({
  selector: 'app-tour-creation',
  templateUrl: './tour-creation.component.html',
  styleUrls: ['./tour-creation.component.css']
})
export class TourCreationComponent implements OnInit {


  tourCreationForm!: FormGroup;
  tagList: string[] = []//["Hiking", "Walk", "Run", "Summer", "Spring", "Fall", "Winter", "See", "Mountain", "City", "Village"]
  readonly templateKeywords = signal(this.tagList);

  constructor(private fb:FormBuilder, private tourService: ToursService, private router: Router, private snackBar:MatSnackBar, private authService: AuthService){}

  ngOnInit(): void {

    const user = this.authService.getCurrentUser();
    
    if (!user || !user.userId) {
        this.snackBar.open("You must be logged in to create a tour", "Close", {duration: 3000, horizontalPosition: "center"});
        return;
    }
    
    this.tourCreationForm=this.fb.group({
      name:['',Validators.required],
      description: ['',Validators.required],
      difficulty:[0,Validators.required],
      tags:[''],
      status: [0],
      cost:[0],
      tourLength:[0],
      creatorID: [Number(user.userId)]
    });
  }


  create(): void{

    if(this.tourCreationForm.invalid){
      this.snackBar.open("You must fill all fields!","Close",{duration:3000,horizontalPosition:"center"}) 
      return;
    }

    const newTour: Tour = this.tourCreationForm.value;

    this.tourService.createTour(newTour).subscribe({
      next:()=>{
        this.snackBar.open("Tour created successfully","Close",{duration:3000,horizontalPosition:"center"});
        this.tourCreationForm.reset();
      },
      error: (err)=>{
        console.error('Error creating tour',err);
        this.snackBar.open("Cannot create tour!","Close",{duration:3000,horizontalPosition:"center"}) 
      }
    })
  }

  addTag(event: MatChipInputEvent): void {
      const value = (event.value || '').trim();
  
      if (value) {
        this.templateKeywords.update(keywords => [...keywords, value]);
        this.tagList.push(value);

        this.tourCreationForm.patchValue({
          tags: this.tagList.join(' #')
        });
      }
      event.chipInput!.clear();
    }

     removeTag(keyword: string) {
      this.templateKeywords.update(keywords => {
        const index = keywords.indexOf(keyword);
        if (index >= 0) {
          keywords.splice(index,1);
          this.tagList=keywords;
          
          this.tourCreationForm.patchValue({
            tags:this.tagList.join(' #')
          });
        }
        return [...keywords];
      });
    }

}
