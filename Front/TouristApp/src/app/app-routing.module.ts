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
const routes: Routes = [
  {path: 'login', component:LoginComponent},
  {path: 'register', component:SignupComponent},
  {path: 'userAccounts', component: UserAccountsComponent},
  {path: 'tourCreation', component: TourCreationComponent},
  {path: 'accounts', component:AccountsComponent},
  {path: 'toursList', component: ToursListComponent},
  {path: 'blogCreation', component: BlogCreationComponent},
  {path: 'blogsList', component: BlogsListComponent}, 
    { path: 'blogs/edit/:id', component: BlogEditComponent }, 
  {path: 'blogs/:id/comments', component: BlogCommentsComponent }

];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
