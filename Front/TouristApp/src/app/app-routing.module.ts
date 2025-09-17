import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { LoginComponent } from './auth/login/login.component';
import { SignupComponent } from './auth/signup/signup.component';
import { UserAccountsComponent } from './stakeholders/user-accounts/user-accounts.component';
import { TourCreationComponent } from './tours/tour-creation/tour-creation.component';
import { AccountsComponent } from './auth/accounts/accounts.component';
import { ToursListComponent } from './tours/tours-list/tours-list.component';
import { BlogCreationComponent } from './blogs/blog-creation/blog-creation.component';
import { BlogsListComponent } from './blogs/blogs-list/blogs-list.component';
import { BlogCommentsComponent } from './blogs/blog-comments/blog-comments.component';
import { BlogEditComponent } from './blogs/blog-edit/blog-edit.component';
import { FollowRecommendationsComponent } from './followers/follow-recommendations/follow-recommendations.component';
import { AddReviewComponent } from './tours/add-review/add-review.component';
import { TourReviewComponent } from './tours/tour-review/tour-review.component';
import { TourDetailsComponent } from './tours/tour-details/tour-details.component';
import { PositionComponent } from './tours/position/position.component';
import { TourExecutionComponent } from './tours/tour-execution/tour-execution.component';
const routes: Routes = [
  {path: 'login', component:LoginComponent},
  {path: 'register', component:SignupComponent},
  {path: 'userAccounts', component: UserAccountsComponent},
  {path: 'tourCreation', component: TourCreationComponent},
  {path: 'accounts', component:AccountsComponent},
  {path: 'toursList', component: ToursListComponent},
  {path: 'allTours', component: ToursListComponent },
  {path: 'blogCreation', component: BlogCreationComponent},
  {path: 'blogsList', component: BlogsListComponent}, 
  {path: 'recommendations', component: FollowRecommendationsComponent},
  {path: 'blogs/edit/:id', component: BlogEditComponent }, 
  {path: 'blogs/:id/comments', component: BlogCommentsComponent },
  {path: 'addReview/:id', component: AddReviewComponent },
  {path: 'tours/:id/reviews', component: TourReviewComponent },
  {path: 'tours/:id', component: TourDetailsComponent },
  {path: 'position', component: PositionComponent },
  {path: 'tourExecution/:id', component: TourExecutionComponent}
  

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
