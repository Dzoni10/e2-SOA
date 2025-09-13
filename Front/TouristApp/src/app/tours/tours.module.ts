import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TourCreationComponent } from './tour-creation/tour-creation.component';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';
import {MatChipsModule } from '@angular/material/chips'
import { MatInputModule } from '@angular/material/input';
import { MatButtonModule } from '@angular/material/button';
import {MatCheckboxModule} from '@angular/material/checkbox'
import {MatIconModule} from '@angular/material/icon'
import { ReactiveFormsModule } from '@angular/forms';
import { ToursListComponent } from './tours-list/tours-list.component';
import { MatCardModule } from '@angular/material/card';
import { MatGridListModule } from '@angular/material/grid-list';
import { FormsModule } from '@angular/forms';
import { AddReviewComponent } from './add-review/add-review.component';
import { TourReviewComponent } from './tour-review/tour-review.component';
import { AddKeypointComponent } from './add-keypoint/add-keypoint.component';
import { MatDialogModule } from '@angular/material/dialog';
import { TourDetailsComponent } from './tour-details/tour-details.component';
import { PositionComponent } from './position/position.component';



@NgModule({
  declarations: [
    TourCreationComponent,
    ToursListComponent,
    AddReviewComponent,
    TourReviewComponent,
    AddKeypointComponent,
    TourDetailsComponent,
    PositionComponent
  ],
  imports: [
    CommonModule,
    MatFormFieldModule,
    MatSelectModule,
    MatChipsModule,
    MatInputModule,
    MatButtonModule,
    MatCheckboxModule,
    MatIconModule,
    ReactiveFormsModule,
    MatCardModule,
    MatGridListModule,
    FormsModule,
    MatDialogModule
  ]
})
export class ToursModule { }
