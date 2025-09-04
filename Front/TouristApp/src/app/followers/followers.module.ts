import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';

// Angular Material imports
import { MatCardModule } from '@angular/material/card';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatChipsModule } from '@angular/material/chips';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatListModule } from '@angular/material/list';
import { MatBadgeModule } from '@angular/material/badge';
import { MatTabsModule } from '@angular/material/tabs';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatProgressBarModule } from '@angular/material/progress-bar';

// Components
import { FollowRecommendationsComponent } from './follow-recommendations/follow-recommendations.component';
//import { FollowersListComponent } from './followers-list/followers-list.component';
//import { FollowingListComponent } from './following-list/following-list.component';
//import { UserStatsComponent } from './user-stats/user-stats.component';

// Routes
const routes = [
  {
    path: 'recommendations',
    component: FollowRecommendationsComponent
  }
];

@NgModule({
  declarations: [
    FollowRecommendationsComponent
    //FollowersListComponent,
    //FollowingListComponent,
    //UserStatsComponent
  ],
  imports: [
    CommonModule,
    RouterModule.forChild(routes),
    ReactiveFormsModule,
    FormsModule,
    
    // Angular Material modules
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatSnackBarModule,
    MatToolbarModule,
    MatListModule,
    MatBadgeModule,
    MatTabsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatProgressBarModule
  ],
  exports: [
    FollowRecommendationsComponent
    //FollowersListComponent,
   // FollowingListComponent,
   // UserStatsComponent
  ]
})
export class FollowersModule { }