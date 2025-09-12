import { Component, OnInit,signal } from '@angular/core';
import { FormBuilder, FormGroup, Validators,FormControl } from '@angular/forms';
import { ToursService } from '../tours.service';
import { Router } from '@angular/router';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Tour } from '../model/tour.model';
import { MatChipInputEvent } from '@angular/material/chips';
import { AuthService } from 'src/app/auth/auth.service';
import { Keypoint } from '../model/keypoint.model';
import { AddKeypointComponent } from '../add-keypoint/add-keypoint.component';
import { MatDialog } from '@angular/material/dialog';
import { forkJoin, map, of, switchMap, throwError } from 'rxjs';

@Component({
  selector: 'app-tour-creation',
  templateUrl: './tour-creation.component.html',
  styleUrls: ['./tour-creation.component.css']
})
export class TourCreationComponent implements OnInit {
  keypoints: Keypoint[] = [];

  tourCreationForm!: FormGroup;
  tagList: string[] = []//["Hiking", "Walk", "Run", "Summer", "Spring", "Fall", "Winter", "See", "Mountain", "City", "Village"]
  readonly templateKeywords = signal(this.tagList);
  isCreating = false;
  
  constructor(private fb:FormBuilder, private tourService: ToursService, private router: Router, private snackBar:MatSnackBar, private authService: AuthService, private dialog: MatDialog){}

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

  openKeypointModal() {
    const dialogRef = this.dialog.open(AddKeypointComponent, {
      width: '800px',
      maxWidth: '95vw',
      maxHeight: '90vh',
      disableClose: false,
      data: {}
    });

    dialogRef.afterClosed().subscribe(result => {
      if (result) {
        this.addKeypointToList(result);
      }
    });
  }

  private addKeypointToList(keypointData: any) {
    const keypoint: Keypoint = {
      name: keypointData.name,
      description: keypointData.description,
      latitude: keypointData.latitude,
      longitude: keypointData.longitude,
      images: keypointData.images,
      formData: keypointData.formData,
    };

    this.keypoints.push(keypoint);
  }

  removeKeypoint(index: number) {
    // Clean up image preview URLs
    const keypoint = this.keypoints[index];
    keypoint.images.forEach(imageUrl => {
      if (imageUrl.startsWith('blob:')) {
        URL.revokeObjectURL(imageUrl);
      }
    });

    this.keypoints.splice(index, 1);
  }

   // Main method for creating tours with keypoints
  async createNew() {
  if (this.isCreating) return;
  if (this.tourCreationForm.invalid) return;

  this.isCreating = true;

  try {
    // 1️⃣ Create the tour first
    const newTour: Tour = {
      ...this.tourCreationForm.value,
      tags: this.templateKeywords(),
      authorId: this.tourCreationForm.get('creatorID')?.value
    };

    const createdTour = await this.tourService.createTour(newTour).toPromise();
    if (!createdTour?.id) throw new Error('Tour creation failed');

    const tourId = createdTour.id;

    // 2️⃣ Create keypoints with proper tourId
    const keypointRequests = this.keypoints.map((kp, index) => {
      const payload = new FormData();
      payload.set('name', kp.name);
      payload.set('description', kp.description);
      payload.set('latitude', kp.latitude.toString());
      payload.set('longitude', kp.longitude.toString());
      payload.set('order', (index + 1).toString());
      payload.set('tourId', tourId); // ✅ assign real tourId

      if (kp.formData) {
        const files = kp.formData.getAll('images') as File[];
        files.forEach(file => payload.append('images', file));
      }

      return this.tourService.createKeypoint(payload).toPromise();
    });

    const createdKeypoints = await Promise.all(keypointRequests);
    await this.tourService.updateTourLength(tourId).toPromise();

    this.snackBar.open('Tour and keypoints created successfully', 'Close', { duration: 3000 });
    
    this.resetForm();

  } catch (error) {
    console.error(error);
    this.snackBar.open('Cannot create tour: ' + error, 'Close', { duration: 5000 });
  } finally {
    this.isCreating = false;
  }
}




  private resetForm() {
    this.tourCreationForm.reset();
    this.tagList = [];
    this.templateKeywords.set([]);
    
    // Clean up keypoint image previews
    this.keypoints.forEach(keypoint => {
      if (keypoint.images) {
        keypoint.images.forEach(imageUrl => {
          if (typeof imageUrl === 'string' && imageUrl.startsWith('blob:')) {
            URL.revokeObjectURL(imageUrl);
          }
        });
      }
    });
    
    this.keypoints = [];

    // Reset form with initial values
    const user = this.authService.getCurrentUser();
    if (user && user.userId) {
      this.tourCreationForm.patchValue({
        creatorID: Number(user.userId),
        difficulty: 0,
        status: 0,
        cost: 0,
        tourLength: 0
      });
    }
  }

}
